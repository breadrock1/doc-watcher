package ocr

import (
	"doc-notifier/internal/pkg/ocr/assistant"
	"doc-notifier/internal/pkg/ocr/dedoc"
	"doc-notifier/internal/pkg/ocr/raw"
	"doc-notifier/internal/pkg/ocr/tesseract"
)

type OcrService struct {
	Ocr Recognizer
}

type Recognizer interface {
	RecognizeFile(filePath string) (string, error)
	RecognizeFileData(data []byte) (string, error)
}

func New(options *Options) *OcrService {
	service := &OcrService{}

	switch options.Mode {
	case ReadRawFile:
		service.Ocr = raw.New()
	case LocalTesseract:
		service.Ocr = tesseract.New()
	case RemoteDedoc:
		service.Ocr = dedoc.New(options.Address)
	case RemoteTesseract:
		service.Ocr = assistant.New(options.Address)
	default:
		service.Ocr = raw.New()
	}

	return service
}
