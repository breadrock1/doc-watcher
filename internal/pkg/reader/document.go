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
	bucketName    = "common_folder"
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
	FolderID            string       `json:"folder_id"`
	FolderPath          string       `json:"folder_path"`
	ContentUUID         string       `json:"content_uuid"`
	ContentMD5          string       `json:"content_md5"`
	Content             string       `json:"content"`
	ContentVector       []float64    `json:"content_vector"`
	DocumentMD5         string       `json:"document_md5"`
	DocumentSSDEEP      string       `json:"document_ssdeep"`
	DocumentName        string       `json:"document_name"`
	DocumentPath        string       `json:"document_path"`
	DocumentSize        int64        `json:"document_size"`
	DocumentType        string       `json:"document_type"`
	DocumentExtension   string       `json:"document_extension"`
	DocumentPermissions int32        `json:"document_permissions"`
	DocumentCreated     string       `json:"document_created"`
	DocumentModified    string       `json:"document_modified"`
	QualityRecognized   int32        `json:"quality_recognition"`
	OcrMetadata         *OcrMetadata `json:"ocr_metadata"`
}

type OcrMetadata struct {
	JobId      string       `json:"job_id"`
	Text       string       `json:"text"`
	PagesCount int          `json:"pages_count"`
	DocType    string       `json:"doc_type"`
	Artifacts  []*Artifacts `json:"artifacts"`
}

type Artifacts struct {
	GroupName     string `json:"group_name"`
	GroupJsonName string `json:"group_json_name"`
	GroupValues   []struct {
		Name     string `json:"name"`
		JsonName string `json:"json_name"`
		Type     string `json:"type"`
		Value    string `json:"value"`
	} `json:"group_values"`
}

func ParseFile(filePath string) (*Document, error) {
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

	document := &Document{}
	document.FolderID = bucketName
	document.FolderPath = parseBucketName(absFilePath)
	document.DocumentPath = absFilePath
	document.DocumentName = fileInfo.Name()
	document.DocumentSize = fileInfo.Size()
	document.DocumentType = parseDocumentType(fileExt)
	document.DocumentExtension = fileExt
	document.DocumentPermissions = int32(fileInfo.Mode().Perm())
	document.ContentUUID = uuid.NewString()
	document.DocumentModified = modifiedTime.Format(timeFormat)
	document.DocumentCreated = createdTime.Format(timeFormat)
	document.QualityRecognized = -1

	document.ComputeMd5HashData(data)
	document.ComputeSsdeepHashData(data)

	return document, nil
}

func (d *Document) SetFolderPath(path string) {
	d.FolderPath = path
}

func (d *Document) SetDocumentPath(path string) {
	d.DocumentPath = path
}

func (d *Document) SetContentData(data string) {
	d.Content = data
}

func (d *Document) SetContentVector(data []float64) {
	d.ContentVector = data
}

func (d *Document) AppendContentVector(data []float64) {
	d.ContentVector = append(d.ContentVector, data...)
}

func (d *Document) ComputeMd5Hash() {
	data := []byte(d.Content)
	d.ComputeMd5HashData(data)
}

func (d *Document) ComputeMd5HashData(data []byte) {
	d.DocumentMD5 = fmt.Sprintf("%x", md5.Sum(data))
}

func (d *Document) ComputeContentMd5Hash() {
	data := []byte(d.Content)
	d.ContentMD5 = fmt.Sprintf("%x", md5.Sum(data))
}

func (d *Document) SetContentMd5Hash(contentID string) {
	d.ContentMD5 = contentID
}

func (d *Document) ComputeSsdeepHash() {
	data := []byte(d.Content)
	d.ComputeSsdeepHashData(data)
}

func (d *Document) ComputeSsdeepHashData(data []byte) {
	if hashData, err := ssdeep.FuzzyBytes(data); err == nil {
		d.DocumentSSDEEP = hashData
	}
}

func (d *Document) SetContentUUID(contentUUID string) {
	d.ContentUUID = contentUUID
}

func (d *Document) ComputeContentUUID() {
	data := []byte(d.Content)
	if uuidToken, err := uuid.FromBytes(data); err == nil {
		d.ContentUUID = uuidToken.String()
	}
}

func (d *Document) SetQuality(quality int32) {
	d.QualityRecognized = quality
}

func (d *Document) SetOcrMetadata(ocr *OcrMetadata) {
	d.OcrMetadata = ocr
}

func (d *Document) GetArtifacts() []*Artifacts {
	if d.OcrMetadata == nil {
		return make([]*Artifacts, 0)
	}

	if d.OcrMetadata.Artifacts == nil {
		return make([]*Artifacts, 0)
	}

	return d.OcrMetadata.Artifacts
}

func (d *Document) MoveMetadataTextToContent() {
	if d.OcrMetadata == nil {
		return
	}

	if len(d.OcrMetadata.Text) == 0 {
		return
	}

	d.Content = d.OcrMetadata.Text
	d.OcrMetadata.Text = ""
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
