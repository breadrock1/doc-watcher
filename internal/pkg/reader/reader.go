package reader

import (
	"crypto/md5"
	"fmt"
	"github.com/glaslos/ssdeep"
	"github.com/google/uuid"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"
)

const DateTimeFormatting = "2004-01-01T00:00:00Z"

type FileReader struct {
}

type Document struct {
	BucketUuid          string    `json:"bucket_uuid"`
	BucketPath          string    `json:"bucket_path"`
	ContentUuid         string    `json:"content_uuid"`
	ContentMd5          string    `json:"content_md5"`
	Content             string    `json:"content"`
	ContentVector       []float64 `json:"content_vector"`
	DocumentMd5         string    `json:"document_md5"`
	DocumentSsdeep      string    `json:"document_ssdeep"`
	DocumentName        string    `json:"document_name"`
	DocumentPath        string    `json:"document_path"`
	DocumentSize        int64     `json:"document_size"`
	DocumentType        string    `json:"document_type"`
	DocumentExtension   string    `json:"document_extension"`
	DocumentPermissions int32     `json:"document_permissions"`
	DocumentCreated     string    `json:"document_created"`
	DocumentModified    string    `json:"document_modified"`
}

func New() *FileReader {
	return &FileReader{}
}

func (f *FileReader) ParseCaughtFiles(path string) []*Document {
	var documents []*Document

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println("Failed to get stat for file: ", err)
		return documents
	}

	if fileInfo.IsDir() {
		log.Println(path, " is a directory. Extracting entities...")
		return f.parseDirectory(path)
	}

	doc, err := f.parseFile(path)
	if err == nil {
		log.Println("Appended parsed document: ", doc.DocumentName)
		documents = append(documents, doc)
	}

	return documents
}

func (f *FileReader) parseDirectory(dirPath string) []*Document {
	var collectedDocuments []*Document
	dirEntry, err := os.ReadDir(dirPath)
	if err != nil {
		log.Println("Failed while reading directory: ", err)
		return collectedDocuments
	}

	for _, entry := range dirEntry {
		if entry.IsDir() {
			log.Println(entry.Name(), " is a directory. Extracting entities...")
			parsedDocuments := f.parseDirectory(entry.Name())
			collectedDocuments = append(collectedDocuments, parsedDocuments...)
		} else {
			log.Println(entry.Name(), " is a file. Parsing content...")
			if document, err := f.parseFile(entry.Name()); err != nil {
				collectedDocuments = append(collectedDocuments, document)
			}
		}
	}

	return collectedDocuments
}

func (f *FileReader) parseFile(filePath string) (*Document, error) {
	absFilePath, _ := filepath.Abs(filePath)
	fileInfo, err := os.Stat(absFilePath)
	if err != nil {
		log.Println("Failed while getting stat of file: ", err)
		return nil, err
	}

	modifiedTime := time.Now()
	createdTime := fileInfo.ModTime()
	modifiedTimeNew := modifiedTime.Format(DateTimeFormatting)
	createdTimeNew := createdTime.Format(DateTimeFormatting)

	document := Document{}
	document.BucketUuid = "common_data"
	document.BucketPath = "/"
	document.DocumentName = fileInfo.Name()
	document.DocumentPath = absFilePath
	document.DocumentSize = fileInfo.Size()
	document.DocumentType = "document"
	document.DocumentExtension = filepath.Ext(filePath)
	document.DocumentPermissions = 777
	document.ContentUuid = uuid.NewString()
	//document.ContentMd5 = uuid.NewString()
	//document.Content = ""
	//document.ContentVector = []string{}
	//document.DocumentMd5 = fileInfo.Name()
	//document.DocumentSsdeep = ""
	document.DocumentModified = modifiedTimeNew
	document.DocumentCreated = createdTimeNew

	return &document, nil
}

func (f *FileReader) SetContentData(document *Document, data string) {
	document.Content = data
}

func (f *FileReader) SetContentVector(document *Document, data []float64) {
	document.ContentVector = data
}

func (f *FileReader) ComputeMd5Hash(document *Document) {
	data := []byte(document.Content)
	document.DocumentMd5 = fmt.Sprintf("%x", md5.Sum(data))
}

func (f *FileReader) ComputeContentMd5Hash(document *Document) {
	data := []byte(document.Content)
	document.ContentMd5 = fmt.Sprintf("%x", md5.Sum(data))
}

func (f *FileReader) ComputeSsdeepHash(document *Document) {
	data := []byte(document.Content)
	if hashData, err := ssdeep.FuzzyBytes(data); err == nil {
		document.DocumentSsdeep = hashData
	}
}

func (f *FileReader) ComputeUuid(document *Document) {
	data := []byte(document.Content)
	if uuidToken, err := uuid.FromBytes(data); err == nil {
		document.ContentUuid = uuidToken.String()
	}
}

func (f *FileReader) SplitContent(content string, chunkSize int) []string {
	strLength := len(content)
	splitLength := int(math.Ceil(float64(strLength) / float64(chunkSize)))
	splitString := make([]string, splitLength)
	var start, stop int
	for i := 0; i < splitLength; i += 1 {
		start = i * chunkSize
		stop = start + chunkSize
		if stop > strLength {
			stop = strLength
		}

		splitString[i] = content[start:stop]
	}

	return splitString
}
