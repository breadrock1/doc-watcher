package ocr

import (
	"doc-notifier/internal/config"
	"doc-notifier/internal/ocr/assistant"
	"doc-notifier/internal/ocr/logoper"
	"doc-notifier/internal/reader"
)

type Service struct {
	Ocr Recognizer
}

type Recognizer interface {
	RecognizeFile(document *reader.Document) error
}

func New(config *config.OcrConfig) *Service {
	service := &Service{}

	switch config.Mode {
	case "assistant":
		service.Ocr = assistant.New(config.Address, config.Timeout)
	case "logoper":
		service.Ocr = logoper.New(config.Address, config.Timeout)
	default:
		service.Ocr = assistant.New(config.Address, config.Timeout)
	}

	return service
}
