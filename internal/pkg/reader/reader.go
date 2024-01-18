package reader

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glaslos/ssdeep"
)

const DateTimeFormatting = "2004-01-01T00:00:00Z"

type FileReader struct {
}

type Document struct {
	BucketUuid          string   `json:"bucket_uuid"`
	BucketPath          string   `json:"bucket_path"`
	DocumentName        string   `json:"document_name"`
	DocumentPath        string   `json:"document_path"`
	DocumentSize        int64    `json:"document_size"`
	DocumentType        string   `json:"document_type"`
	DocumentExtension   string   `json:"document_extension"`
	DocumentPermissions int32    `json:"document_permissions"`
	DocumentMd5Hash     string   `json:"document_md5_hash"`
	DocumentSsdeepHash  string   `json:"document_ssdeep_hash"`
	EntityData          string   `json:"entity_data"`
	EntityKeywords      []string `json:"entity_keywords"`
	DocumentCreated     string   `json:"document_created"`
	DocumentModified    string   `json:"document_modified"`
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
		return f.parseDirectory(path)
	}

	doc, err := f.parseFile(path)
	if err == nil {
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
			parsedDocuments := f.parseDirectory(entry.Name())
			collectedDocuments = append(collectedDocuments, parsedDocuments...)
		} else {
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
	document.DocumentMd5Hash = fileInfo.Name()
	document.DocumentSsdeepHash = ""
	document.DocumentModified = modifiedTimeNew
	document.DocumentCreated = createdTimeNew
	document.EntityKeywords = make([]string, 0)

	return &document, nil
}

func (f *FileReader) SetEntityData(document *Document, data string) {
	document.EntityData = data
}

func (f *FileReader) ComputeMd5Hash(document *Document) {
	data := []byte(document.EntityData)
	document.DocumentMd5Hash = fmt.Sprintf("%x", md5.Sum(data))
}

func (f *FileReader) ComputeSsdeepHash(document *Document) {
	data := []byte(document.EntityData)
	if hashData, err := ssdeep.FuzzyBytes(data); err == nil {
		document.DocumentSsdeepHash = hashData
	}
}
