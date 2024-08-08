package watcher

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"doc-notifier/internal/models"
)

const MaxQualityValue = 10000

var (
	bucketName    = "common-folder"
	timeFormat    = time.RFC3339
	documentMimes = []string{
		"csv", "msword", "html", "json", "pdf",
		"rtf", "plain", "vnd.ms-excel", "xml",
		"vnd.ms-powerpoint", "vnd.oasis.opendocument.text",
		"vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"vnd.openxmlformats-officedocument.wordprocessingml.document",
		"vnd.openxmlformats-officedocument.presentationml.presentation",
	}
)

func ParseCaughtFiles(filePath string) []*models.Document {
	mu := &sync.Mutex{}
	var customList []*models.Document

	wg := &sync.WaitGroup{}
	for _, filePath := range getEntityFiles(filePath) {
		wg.Add(1)
		filePath := filePath

		go func() {
			defer wg.Done()

			if doc, err := ParseFile(filePath); err == nil {
				log.Println("Caught parsed document: ", doc.DocumentName)
				mu.Lock()
				customList = append(customList, doc)
				mu.Unlock()
				return
			}

			log.Println("Failed parsing document: ", filePath)
		}()
	}

	wg.Wait()

	return customList
}

func ParseFile(filePath string) (*models.Document, error) {
	absFilePath, _ := filepath.Abs(filePath)
	fileInfo, err := os.Stat(absFilePath)
	if err != nil {
		err = fmt.Errorf("failed while getting stat of file: %e", err)
		return nil, err
	}

	modifiedTime := time.Now().UTC()
	createdTime := fileInfo.ModTime().UTC()
	fileExt := filepath.Ext(filePath)
	data, _ := os.ReadFile(absFilePath)

	document := &models.Document{}
	document.FolderID = bucketName
	document.FolderPath = parseBucketName(absFilePath)
	document.DocumentPath = absFilePath
	document.DocumentName = fileInfo.Name()
	document.DocumentSize = fileInfo.Size()
	document.DocumentType = parseDocumentType(fileExt)
	document.DocumentExtension = fileExt
	document.DocumentPermissions = int32(fileInfo.Mode().Perm())
	document.DocumentModified = modifiedTime.Format(timeFormat)
	document.DocumentCreated = createdTime.Format(timeFormat)
	document.QualityRecognized = -1

	document.ComputeMd5HashData(data)
	document.ComputeSsdeepHashData(data)

	return document, nil
}

func parseBucketName(filePath string) string {
	currPath := os.Getenv("PWD")
	relPath, err := filepath.Rel(currPath, filePath)
	relPath2, err := filepath.Rel("indexer", relPath)
	if err != nil {
		log.Printf("Failed while parsing bucket name")
		return bucketName
	}

	bucketNameRes, _ := filepath.Split(relPath2)
	bucketNameRes2 := strings.ReplaceAll(bucketNameRes, "/", "")
	if bucketNameRes2 == "" {
		return bucketName
	}

	return bucketNameRes2
}

func parseDocumentType(extension string) string {
	mimeType := mime.TypeByExtension(extension)
	attributes := strings.Split(mimeType, "/")
	switch attributes[0] {
	case "audio":
		return "audio"
	case "image":
		return "image"
	case "video":
		return "video"
	case "text":
		return "document"
	case "application":
		return extractApplicationMimeType(attributes[1])
	default:
		return "document"
	}
}

func extractApplicationMimeType(attribute string) string {
	for _, mimeType := range documentMimes {
		if mimeType == attribute {
			return "document"
		}
	}

	return "document"
}

func getEntityFiles(filePath string) []string {
	<-time.After(time.Second)

	var files []string
	err := filepath.Walk(filePath, visitEntity(&files))
	if err != nil {
		log.Println("Error while walking: ", err)
	}
	return files
}

func visitEntity(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		if !info.IsDir() {
			*files = append(*files, path)
		}

		return nil
	}
}
