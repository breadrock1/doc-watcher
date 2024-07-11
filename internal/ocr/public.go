package ocr

import (
	"doc-notifier/internal/config"
	"doc-notifier/internal/models"
	"doc-notifier/internal/ocr/assistant"
	"doc-notifier/internal/ocr/raw"
)

type Service struct {
	Ocr Recognizer
}

type Recognizer interface {
	RecognizeFile(document *models.Document) error
}

func New(config *config.OcrConfig) *Service {
	service := &Service{}

	switch config.Mode {
	case "assistant":
		service.Ocr = assistant.New(config)
	case "raw":
		service.Ocr = raw.New()
	default:
		service.Ocr = raw.New()
	}

	return service
}
