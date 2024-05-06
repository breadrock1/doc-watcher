package reader

const MaxQualityValue = 10000

type DocumentPreview struct {
	DocumentID        string       `json:"id" example:"886f7e11874040ca8b8461fb4cd1aa2c"`
	DocumentName      string       `json:"name" example:"document_name.pdf"`
	CreatedAt         string       `json:"created_at" example:"2024-05-04T22:53:57Z"`
	QualityOCR        int          `json:"quality_recognition" example:"10000"`
	FileSize          int64        `json:"file_size" example:"311652"`
	Location          string       `json:"location" example:"unrecognized"`
	PreviewProperties []*Artifacts `json:"preview_properties"`
}

type PreviewProperties struct {
	Key   string `json:"key" example:"field_date_transaction"`
	Name  string `json:"name" example:"Date and time of transaction"`
	Value string `json:"value" example:"18.03.2024, 23:59"`
}

func FromDocument(document *Document) *DocumentPreview {
	var location string
	var previewProperties []*Artifacts
	ocrQuality := 1
	if document.OcrMetadata != nil {
		ocrQuality = MaxQualityValue
		location = document.OcrMetadata.DocType
		previewProperties = document.GetArtifacts()
	}

	return &DocumentPreview{
		DocumentID:        document.DocumentMD5,
		DocumentName:      document.DocumentName,
		CreatedAt:         document.DocumentCreated,
		QualityOCR:        ocrQuality,
		FileSize:          document.DocumentSize,
		Location:          location,
		PreviewProperties: previewProperties,
	}
}
