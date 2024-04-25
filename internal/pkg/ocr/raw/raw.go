package raw

import (
	"doc-notifier/internal/pkg/ocr/processing"
	"doc-notifier/internal/pkg/reader"
	"log"
	"os"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (re *Service) RecognizeFile(document *reader.Document) (string, error) {
	filePath := document.DocumentPath
	bytesData, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Failed while reading file: ", err)
		return "", err
	}

	stringData := string(bytesData)
	if len(stringData) == 0 {
		log.Println("Failed: returned empty string data...")
		return "", err
	}

	return stringData, nil
}

func (re *Service) RecognizeFileData(data []byte) (string, error) {
	return string(data), nil
}

func (re *Service) GetProcessingJobs() map[string]*processing.ProcessJob {
	return make(map[string]*processing.ProcessJob)
}

func (re *Service) GetProcessingJob(jobId string) *processing.ProcessJob {
	return nil
}
