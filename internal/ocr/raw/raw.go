package raw

import (
	"doc-notifier/internal/ocr/artifacts"
	"fmt"
	"log"
	"os"

	"doc-notifier/internal/ocr/processing"
	"doc-notifier/internal/reader"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) RecognizeFile(document *reader.Document) error {
	bytesData, err := os.ReadFile(document.DocumentPath)
	if err != nil {
		log.Println("Failed while reading file: ", err)
		return err
	}

	stringData := string(bytesData)
	if len(stringData) == 0 {
		return fmt.Errorf("returned empty content data for: %s", document.DocumentID)
	}

	document.SetContentData(stringData)
	return nil
}

func (s *Service) GetProcessingJobs() map[string]*processing.ProcessJob {
	return make(map[string]*processing.ProcessJob)
}

func (s *Service) GetProcessingJob(jobId string) *processing.ProcessJob {
	return nil
}

func (s *Service) GetArtifacts() *artifacts.OcrArtifacts {
	return nil
}
