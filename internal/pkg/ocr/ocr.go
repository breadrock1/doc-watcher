package ocr

import (
	"doc-notifier/internal/pkg/ocr/assistant"
	"doc-notifier/internal/pkg/ocr/dedoc"
	"doc-notifier/internal/pkg/ocr/raw"
)

type Service struct {
	Ocr Recognizer
}

type Recognizer interface {
	RecognizeFile(filePath string) (string, error)
	RecognizeFileData(data []byte) (string, error)
}

func New(options *Options) *Service {
	service := &Service{}

	switch options.Mode {
	case ReadRawFile:
		service.Ocr = raw.New()
	case DedocWrapper:
		service.Ocr = dedoc.New(options.Address, options.Timeout)
	case AssistantMode:
		service.Ocr = assistant.New(options.Address, options.Timeout)
	default:
		service.Ocr = raw.New()
	}

	return service
}
