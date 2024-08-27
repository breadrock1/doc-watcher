package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/models"
	"doc-notifier/internal/ocr"
	"doc-notifier/internal/searcher"
	"doc-notifier/internal/summarizer"
	"doc-notifier/internal/tokenizer"
	"doc-notifier/internal/watcher"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/notification"
)

type MinioWatcher struct {
	stopCh chan bool

	Address       string
	pauseWatchers bool

	mc *minio.Client

	Ocr        *ocr.Service
	Searcher   *searcher.Service
	Tokenizer  *tokenizer.Service
	Summarizer *summarizer.Service
}

func New(
	config *config.MinioConfig,
	ocrService *ocr.Service,
	searcherService *searcher.Service,
	tokenService *tokenizer.Service,
	summarizeService *summarizer.Service,
) *watcher.Service {
	client, _ := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioRootUser, config.MinioRootPassword, ""),
		Secure: config.MinioUseSSL,
	})

	watcherInst := &MinioWatcher{
		stopCh:        make(chan bool),
		Address:       config.Address,
		pauseWatchers: false,

		mc: client,

		Ocr:        ocrService,
		Searcher:   searcherService,
		Tokenizer:  tokenService,
		Summarizer: summarizeService,
	}

	return &watcher.Service{Watcher: watcherInst}
}

func (mw *MinioWatcher) RunWatchers() {
	go mw.launchProcessEventLoop()
	<-mw.stopCh
}

func (mw *MinioWatcher) IsPausedWatchers() bool {
	return mw.pauseWatchers
}

func (mw *MinioWatcher) PauseWatchers(flag bool) {
	mw.pauseWatchers = flag
}

func (mw *MinioWatcher) TerminateWatchers() {
	mw.stopCh <- true
}

func (mw *MinioWatcher) GetAddress() string {
	return mw.Address
}

func (mw *MinioWatcher) GetWatchedDirectories() []string {
	ctx := context.Background()
	buckets, err := mw.mc.ListBuckets(ctx)
	if err != nil {
		return make([]string, 0)
	}

	bucketNames := make([]string, len(buckets))
	for _, bucketInfo := range buckets {
		bucketNames = append(bucketNames, bucketInfo.Name)
	}

	return bucketNames
}

func (mw *MinioWatcher) GetHierarchy(bucket, dirName string) []*models.StorageItem {
	ctx := context.Background()
	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Prefix:    dirName,
		Recursive: false,
	}

	dirObjects := make([]*models.StorageItem, 0)
	for obj := range mw.mc.ListObjects(ctx, bucket, opts) {
		if obj.Err != nil {
			log.Println(obj.Err)
			continue
		}

		dirObjects = append(dirObjects, &models.StorageItem{
			FileName:      obj.Key,
			DirectoryName: dirName,
			IsDirectory:   len(obj.ETag) == 0,
		})
	}

	return dirObjects
}

func (mw *MinioWatcher) CreateDirectory(dirName string) error {
	ctx := context.Background()
	opts := minio.MakeBucketOptions{}
	return mw.mc.MakeBucket(ctx, dirName, opts)
}

func (mw *MinioWatcher) RemoveDirectory(dirName string) error {
	ctx := context.Background()
	return mw.mc.RemoveBucket(ctx, dirName)
}

func (mw *MinioWatcher) RemoveFile(bucket string, fileName string) error {
	ctx := context.Background()
	opts := minio.RemoveObjectOptions{}
	return mw.mc.RemoveObject(ctx, bucket, fileName, opts)
}

func (mw *MinioWatcher) UploadFile(bucket string, fileName string, fileData bytes.Buffer) error {
	ctx := context.Background()
	dataLen := int64(fileData.Len())
	opts := minio.PutObjectOptions{}
	_, err := mw.mc.PutObject(ctx, bucket, fileName, &fileData, dataLen, opts)
	return err
}

func (mw *MinioWatcher) DownloadFile(bucket string, objName string) (bytes.Buffer, error) {
	var objBody bytes.Buffer

	ctx := context.Background()
	opts := minio.GetObjectOptions{}
	obj, err := mw.mc.GetObject(ctx, bucket, objName, opts)
	if err != nil {
		return objBody, err
	}

	_, rErr := objBody.ReadFrom(obj)
	if rErr != nil {
		return objBody, err
	}

	return objBody, nil
}

func (mw *MinioWatcher) AppendDirectories(directories []string) error {
	ctx := context.Background()

	var collectedErrs []string
	for _, dir := range directories {
		opts := minio.MakeBucketOptions{}
		if err := mw.mc.MakeBucket(ctx, dir, opts); err != nil {
			collectedErrs = append(collectedErrs, err.Error())
		}
	}

	if len(collectedErrs) > 0 {
		msg := strings.Join(collectedErrs, "\n")
		return errors.New(msg)
	}

	return nil
}

func (mw *MinioWatcher) RemoveDirectories(directories []string) error {
	ctx := context.Background()

	var collectedErrs []string
	for _, dir := range directories {
		if err := mw.mc.RemoveBucket(ctx, dir); err != nil {
			collectedErrs = append(collectedErrs, err.Error())
		}
	}

	if len(collectedErrs) > 0 {
		msg := strings.Join(collectedErrs, "\n")
		return errors.New(msg)
	}

	return nil
}

func (mw *MinioWatcher) launchProcessEventLoop() {
	var (
		prefix       = ""
		suffix       = ""
		eventsFilter = []string{
			"s3:ObjectCreated:*",
			"s3:ObjectRemoved:*",
		}
	)

	ctx := context.Background()
	for event := range mw.mc.ListenNotification(ctx, prefix, suffix, eventsFilter) {
		if mw.pauseWatchers {
			log.Println("catch event but listener has been paused")
			continue
		}

		if event.Err != nil {
			log.Println("failed event: ", event.Err)
			continue
		}

		mw.extractDocumentFromEvent(event)
	}
}

func (mw *MinioWatcher) extractDocumentFromEvent(event notification.Info) {
	for _, record := range event.Records {
		s3Object := record.S3
		bucketName := s3Object.Bucket.Name
		folderPath := s3Object.Object.Key

		fileName := s3Object.Object.Key
		filePath := path.Join(folderPath, fileName)
		fileSize := s3Object.Object.Size
		fileType := watcher.ParseDocumentType(s3Object.Object.ContentType)

		fileExt := path.Ext(fileName)

		createdAt := time.Now().Format(time.RFC3339)
		modifiedAt := createdAt

		document := &models.Document{}
		document.FolderID = bucketName
		document.FolderPath = folderPath
		document.DocumentPath = filePath
		document.DocumentName = fileName
		document.DocumentSize = fileSize
		document.DocumentType = fileType
		document.DocumentExtension = fileExt
		document.DocumentPermissions = int32(777)
		document.DocumentModified = modifiedAt
		document.DocumentCreated = createdAt
		document.QualityRecognized = -1

		data, err := mw.DownloadFile(bucketName, fileName)
		if err != nil {
			log.Println("failed to load file data: ", err.Error())
			continue
		}

		document.ComputeMd5HashData(data.Bytes())
		document.ComputeSsdeepHashData(data.Bytes())
		tmpFilePath := fmt.Sprintf("%s/%s", "./uploads", fileName)
		document.Content = tmpFilePath
		err = os.WriteFile(tmpFilePath, data.Bytes(), os.ModePerm)
		if err != nil {
			log.Println("failed to write file: ", err)
			continue
		}

		mw.recognizeDocument(document)
	}
}

func (mw *MinioWatcher) recognizeDocument(document *models.Document) {
	document.SetQuality(0)
	if err := mw.Ocr.Ocr.RecognizeFile(document, document.Content); err != nil {
		log.Println(err)
		return
	}

	document.ComputeMd5Hash()
	document.ComputeSsdeepHash()
	document.SetEmbeddings([]*models.Embeddings{})

	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := mw.Tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for chunkID, chunkData := range tokenVectors.Vectors {
		text := tokenVectors.ChunkedText[chunkID]
		document.AppendContentVector(text, chunkData)
	}

	log.Println("Storing document to searcher: ", document.DocumentName)
	if err := mw.Searcher.StoreDocument(document); err != nil {
		log.Println("Failed while storing document: ", err)
	}

	mw.Summarizer.LoadSummary(document)
}
