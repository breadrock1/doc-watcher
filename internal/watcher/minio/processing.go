package minio

import (
	"bytes"
	"context"
	"log"
	"path"
	"time"

	"doc-watcher/internal/watcher"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/notification"
)

func (mw *S3Minio) extractAndStoreDocument(ctx context.Context, event notification.Info) {
	for _, record := range event.Records {
		s3Object := record.S3

		bucketName := s3Object.Bucket.Name
		filePath := s3Object.Object.Key
		fileSize := s3Object.Object.Size
		fileName := path.Base(filePath)
		fileExt := path.Ext(fileName)
		folderPath := path.Dir(filePath)
		filePath = path.Join(folderPath, fileName)
		fileType := watcher.ParseDocumentType(s3Object.Object.ContentType)

		log.Printf("caught event with file %s into bucket %s", filePath, bucketName)

		createdAt := time.Now().UTC().Format(time.RFC3339)
		modifiedAt := createdAt

		document := &watcher.Document{}
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

		mw.cacher.Set(fileName, document, mw.config.CacheExpire*time.Minute)

		data, err := mw.downloadFile(ctx, bucketName, fileName)
		if err != nil {
			document.QualityRecognized = 0
			log.Printf("failed to load file data: %w", err)
			continue
		}
		defer data.Reset()

		document.SetQuality(0)
		err = mw.ocrServ.Ocr.RecognizeFile(document, data.Bytes())
		if err != nil {
			log.Printf("failed to recognize file: %w", err)
			return
		}

		mw.recognizeDocument(document)
	}
}

func (mw *S3Minio) downloadFile(ctx context.Context, bucket, filePath string) (bytes.Buffer, error) {
	var objBody bytes.Buffer

	opts := minio.GetObjectOptions{}
	obj, err := mw.mc.GetObject(ctx, bucket, filePath, opts)
	if err != nil {
		return objBody, err
	}

	_, err = objBody.ReadFrom(obj)
	if err != nil {
		return objBody, err
	}

	return objBody, nil
}

func (mw *S3Minio) recognizeDocument(doc *watcher.Document) {
	doc.ComputeMd5Hash()
	doc.ComputeSsdeepHash()
	doc.SetEmbeddings([]*watcher.Embeddings{})

	log.Printf("loading embeddings for doc %s: ", doc.DocumentName)
	tokenVectors, _ := mw.tokenServ.Tokenizer.Tokenize(doc)
	for chunkID, chunkData := range tokenVectors.Vectors {
		text := tokenVectors.ChunkedText[chunkID]
		doc.AppendContentVector(text, chunkData)
	}

	log.Println("storing doc to searcher: ", doc.DocumentName)
	if err := mw.searchServ.StoreDocument(doc); err != nil {
		log.Printf("failed to store doc %s: %w", doc.DocumentName, err)
		if obj, ok := mw.cacher.Get(doc.DocumentName); ok {
			obj.(*watcher.Document).QualityRecognized = 0
		}
	}
}
