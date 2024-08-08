package models

import (
	"crypto/md5"
	"fmt"

	"github.com/glaslos/ssdeep"
	"github.com/google/uuid"
)

type Document struct {
	FolderID            string        `json:"folder_id"`
	FolderPath          string        `json:"folder_path"`
	Content             string        `json:"content"`
	DocumentID          string        `json:"document_id"`
	DocumentSSDEEP      string        `json:"document_ssdeep"`
	DocumentName        string        `json:"document_name"`
	DocumentPath        string        `json:"document_path"`
	DocumentSize        int64         `json:"document_size"`
	DocumentType        string        `json:"document_type"`
	DocumentExtension   string        `json:"document_extension"`
	DocumentPermissions int32         `json:"document_permissions"`
	DocumentClass       string        `json:"document_class"`
	DocumentCreated     string        `json:"document_created"`
	DocumentModified    string        `json:"document_modified"`
	QualityRecognized   int32         `json:"quality_recognition"`
	OcrMetadata         *OcrMetadata  `json:"ocr_metadata"`
	Embeddings          []*Embeddings `json:"embeddings"`
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

type Embeddings struct {
	ChunkID   string    `json:"chunk_id"`
	TextChunk string    `json:"text_chunk"`
	Vector    []float64 `json:"vector"`
}

func DefaultOcr() *OcrMetadata {
	return &OcrMetadata{
		JobId:      "",
		Text:       "",
		PagesCount: 1,
		DocType:    "",
		Artifacts:  make([]*Artifacts, 0),
	}
}

func (d *Document) SetFolderID(folderID string) {
	d.FolderID = folderID
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

func (d *Document) SetEmbeddings(embeddings []*Embeddings) {
	d.Embeddings = embeddings
}

func (d *Document) SetQuality(quality int32) {
	d.QualityRecognized = quality
}

func (d *Document) SetDocumentClass(class string) {
	if len(class) > 0 {
		d.DocumentClass = class
	} else {
		d.DocumentClass = "неизвестно"
	}
}

func (d *Document) SetOcrMetadata(ocr *OcrMetadata) {
	d.OcrMetadata = ocr
}

func (d *Document) GetDocType() string {
	if d.OcrMetadata == nil {
		return ""
	}

	return d.OcrMetadata.DocType
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

func (d *Document) AppendContentVector(text string, tokens []float64) {
	embeddings := &Embeddings{
		ChunkID:   uuid.New().String(),
		Vector:    tokens,
		TextChunk: text,
	}

	d.Embeddings = append(d.Embeddings, embeddings)
}

func (d *Document) ComputeMd5Hash() {
	data := []byte(d.Content)
	d.ComputeMd5HashData(data)
}

func (d *Document) ComputeMd5HashData(data []byte) {
	d.DocumentID = fmt.Sprintf("%x", md5.Sum(data))
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
