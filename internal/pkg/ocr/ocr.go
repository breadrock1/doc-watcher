package ocr

import (
	"doc-notifier/internal/pkg/ocr/assistant"
	"doc-notifier/internal/pkg/ocr/dedoc"
	"doc-notifier/internal/pkg/ocr/logoper"
	"doc-notifier/internal/pkg/ocr/processing"
	"doc-notifier/internal/pkg/ocr/raw"
	"doc-notifier/internal/pkg/reader"
)

type Service struct {
	Ocr Recognizer
}

type Recognizer interface {
	RecognizeFile(document *reader.Document) (string, error)
	RecognizeFileData(data []byte) (string, error)
	GetProcessingJobs() map[string]*processing.ProcessJob
	GetProcessingJob(jobId string) *processing.ProcessJob
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
	case Logoper:
		service.Ocr = logoper.New(options.Address, options.Timeout)
	default:
		service.Ocr = raw.New()
	}

	return service
}
