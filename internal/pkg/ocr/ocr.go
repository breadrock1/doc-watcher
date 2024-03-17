package ocr

import (
	"doc-notifier/internal/pkg/ocr/assistant"
	"doc-notifier/internal/pkg/ocr/dedoc"
	"doc-notifier/internal/pkg/ocr/raw"
	"doc-notifier/internal/pkg/ocr/tesseract"
)

type OcrService struct {
	Ocr *Recognizer
}

type Recognizer interface {
	RecognizeFile(filePath string) (string, error)
	RecognizeFileData(data []byte) (string, error)
}

func New(options *Options) *OcrService {
	service := &OcrService{}

	switch options.Mode {
	case ReadRawFile:
		service.Ocr = raw.New(options)
	case LocalTesseract:
		service.Ocr = tesseract.New(options)
	case RemoteTesseract:
		service.Ocr = assistant.New(options)
	case RemoteDedoc:
		service.Ocr = dedoc.New(options)
	default:
		service.Ocr = raw.New(options)
	}

	return service
}
