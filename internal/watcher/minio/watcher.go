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
	"github.com/patrickmn/go-cache"
)

type MinioWatcher struct {
	stopCh   chan bool
	recFiles *cache.Cache

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
		stopCh:   make(chan bool),
		recFiles: cache.New(10*time.Minute, 30*time.Minute),

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

func (mw *MinioWatcher) GetAddress() string {
	return mw.Address
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

func (mw *MinioWatcher) FetchProcessingDocuments(files []string) *models.ProcessingDocuments {
	procDocs := &models.ProcessingDocuments{}

	for _, file := range files {
		obj, ok := mw.recFiles.Get(file)
		if !ok {
			continue
		}

		document := obj.(*models.Document)
		switch document.QualityRecognized {
		case -1:
			procDocs.Processing = append(procDocs.Processing, file)
		case 0:
			procDocs.Unrecognized = append(procDocs.Unrecognized, file)
		default:
			procDocs.Done = append(procDocs.Done, file)
		}
	}

	return procDocs
}

func (mw *MinioWatcher) CleanProcessingDocuments(files []string) error {
	// TODO: Add RWLock to escape data race!
	for _, file := range files {
		mw.recFiles.Delete(file)
	}

	return nil
}

func (mw *MinioWatcher) GetBuckets() ([]string, error) {
	ctx := context.Background()
	buckets, err := mw.mc.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	bucketNames := make([]string, len(buckets))
	for _, bucketInfo := range buckets {
		bucketNames = append(bucketNames, bucketInfo.Name)
	}

	return bucketNames, nil
}

func (mw *MinioWatcher) GetListFiles(bucket, dirName string) ([]*models.StorageItem, error) {
	ctx := context.Background()
	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Prefix:    dirName,
		Recursive: false,
	}

	if mw.mc.IsOffline() {
		return nil, errors.New("cloud is offline")
	}

	dirObjects := make([]*models.StorageItem, 0)
	for obj := range mw.mc.ListObjects(ctx, bucket, opts) {
		if obj.Err != nil {
			log.Println("failed to get object: ", obj.Err)
			continue
		}

		dirObjects = append(dirObjects, &models.StorageItem{
			FileName:      obj.Key,
			DirectoryName: dirName,
			IsDirectory:   len(obj.ETag) == 0,
		})
	}

	return dirObjects, nil
}

func (mw *MinioWatcher) CreateBucket(dirName string) error {
	ctx := context.Background()
	opts := minio.MakeBucketOptions{}
	return mw.mc.MakeBucket(ctx, dirName, opts)
}

func (mw *MinioWatcher) RemoveBucket(dirName string) error {
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

func (mw *MinioWatcher) CopyFile(bucket, srcPath, dstPath string) error {
	ctx := context.Background()
	srcOpts := minio.CopySrcOptions{Bucket: bucket, Object: srcPath}
	dstOpts := minio.CopyDestOptions{Bucket: bucket, Object: dstPath}
	_, err := mw.mc.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		log.Println("failed to copy object: ", err)
		return err
	}

	return nil
}

func (mw *MinioWatcher) MoveFile(bucket, srcPath, dstPath string) error {
	copyErr := mw.CopyFile(bucket, srcPath, dstPath)
	if copyErr != nil {
		log.Println("failed to copy file: ", copyErr)
		return copyErr
	}

	removeErr := mw.RemoveFile(bucket, srcPath)
	if removeErr != nil {
		log.Println("failed to remove old file: ", removeErr)
		return removeErr
	}

	return nil
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

func (mw *MinioWatcher) GetShareURL(bucket string, fileName string) (string, error) {
	ctx := context.Background()

	url, err := mw.mc.PresignedGetObject(ctx, bucket, fileName, 900*24*time.Hour, map[string][]string{})
	if err != nil {
		return "", err
	}

	return url.String(), nil
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

		createdAt := time.Now().UTC().Format(time.RFC3339)
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

		mw.recFiles.Set(fileName, document, 10*time.Minute)
		data, err := mw.DownloadFile(bucketName, fileName)
		if err != nil {
			document.QualityRecognized = 0
			log.Println("failed to load file data: ", err.Error())
			continue
		}

		document.ComputeMd5HashData(data.Bytes())
		document.ComputeSsdeepHashData(data.Bytes())
		tmpFilePath := fmt.Sprintf("%s/%s", "./uploads", fileName)
		document.Content = tmpFilePath
		err = os.WriteFile(tmpFilePath, data.Bytes(), os.ModePerm)
		if err != nil {
			document.QualityRecognized = 0
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
