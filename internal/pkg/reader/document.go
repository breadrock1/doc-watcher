package reader

import (
	"crypto/md5"
	"fmt"
	"github.com/glaslos/ssdeep"
	"github.com/google/uuid"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	bucketPath    = "/"
	bucketName    = "common_bucket"
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

type Document struct {
	BucketUUID          string    `json:"bucket_uuid"`
	BucketPath          string    `json:"bucket_path"`
	ContentUUID         string    `json:"content_uuid"`
	ContentMD5          string    `json:"content_md5"`
	Content             string    `json:"content"`
	ContentVector       []float64 `json:"content_vector"`
	DocumentMD5         string    `json:"document_md5"`
	DocumentSSDEEP      string    `json:"document_ssdeep"`
	DocumentName        string    `json:"document_name"`
	DocumentPath        string    `json:"document_path"`
	DocumentSize        int64     `json:"document_size"`
	DocumentType        string    `json:"document_type"`
	DocumentExtension   string    `json:"document_extension"`
	DocumentPermissions int32     `json:"document_permissions"`
	DocumentCreated     string    `json:"document_created"`
	DocumentModified    string    `json:"document_modified"`
}

func ParseFile(filePath string) (*Document, error) {
	absFilePath, _ := filepath.Abs(filePath)
	fileInfo, err := os.Stat(absFilePath)
	if err != nil {
		log.Println("Failed while getting stat of file: ", err)
		return nil, err
	}

	modifiedTime := time.Now().UTC()
	createdTime := fileInfo.ModTime().UTC()
	modifiedTimeNew := modifiedTime.Format(timeFormat)
	createdTimeNew := createdTime.Format(timeFormat)

	fileExt := filepath.Ext(filePath)
	filePerms := int32(fileInfo.Mode().Perm())

	document := Document{}
	document.BucketPath = bucketPath
	document.BucketUUID = bucketName
	document.DocumentPath = absFilePath
	document.DocumentName = fileInfo.Name()
	document.DocumentSize = fileInfo.Size()
	document.DocumentType = ParseDocumentType(fileExt)
	document.DocumentExtension = fileExt
	document.DocumentPermissions = filePerms
	document.ContentUUID = uuid.NewString()
	document.DocumentModified = modifiedTimeNew
	document.DocumentCreated = createdTimeNew

	return &document, nil
}

func ParseDocumentType(extension string) string {
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
		return "unknown"
	}
}

func extractApplicationMimeType(attribute string) string {
	for _, mimeType := range documentMimes {
		if mimeType == attribute {
			return "document"
		}
	}

	return "unknown"
}

func (f *Service) SetContentData(document *Document, data string) {
	document.Content = data
}

func (f *Service) SetContentVector(document *Document, data []float64) {
	document.ContentVector = data
}

func (f *Service) AppendContentVector(document *Document, data []float64) {
	document.ContentVector = append(document.ContentVector, data...)
}

func (f *Service) ComputeMd5Hash(document *Document) {
	data := []byte(document.Content)
	document.DocumentMD5 = fmt.Sprintf("%x", md5.Sum(data))
}

func (f *Service) ComputeContentMd5Hash(document *Document) {
	if len(document.DocumentMD5) == 0 {
		f.ComputeMd5Hash(document)
	}
	document.ContentMD5 = document.DocumentMD5
}

func (f *Service) ComputeSsdeepHash(document *Document) {
	data := []byte(document.Content)
	if hashData, err := ssdeep.FuzzyBytes(data); err == nil {
		document.DocumentSSDEEP = hashData
	}
}

func (f *Service) ComputeUUID(document *Document) {
	data := []byte(document.Content)
	if uuidToken, err := uuid.FromBytes(data); err == nil {
		document.ContentUUID = uuidToken.String()
	}
}
