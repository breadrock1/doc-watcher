package reader

const MaxQualityValue = 10000

type DocumentPreview struct {
	DocumentID        string               `json:"id"`
	DocumentName      string               `json:"name"`
	CreatedAt         string               `json:"created_at"`
	QualityOCR        int                  `json:"quality_recognition"`
	FileSize          int64                `json:"file_size"`
	Location          string               `json:"location"`
	PreviewProperties []*PreviewProperties `json:"preview_properties"`
}

type PreviewProperties struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func From(document *Document) *DocumentPreview {
	return &DocumentPreview{
		DocumentID:        document.DocumentMD5,
		DocumentName:      document.DocumentName,
		CreatedAt:         document.DocumentCreated,
		QualityOCR:        0,
		FileSize:          document.DocumentSize,
		Location:          document.DocumentPath,
		PreviewProperties: make([]*PreviewProperties, 0),
	}
}

func (dp *DocumentPreview) SetQuality(quality int) {
	dp.QualityOCR = quality
}
