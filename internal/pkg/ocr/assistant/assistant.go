package assistant

import (
	"bytes"
	"doc-notifier/internal/pkg/ocr/processing"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"
)

const RecognitionURL = "/ocr_extract_text"

type Service struct {
	address string
	timeout time.Duration
}

func New(address string, timeout time.Duration) *Service {
	return &Service{
		address: address,
		timeout: timeout,
	}
}

func (s *Service) RecognizeFile(document *reader.Document) error {
	filePath := document.DocumentPath

	var recErr error
	var fileHandle *os.File
	if fileHandle, recErr = os.Open(filePath); recErr != nil {
		return fmt.Errorf("file %s not found: %e", filePath, recErr)
	}
	defer func() { _ = fileHandle.Close() }()

	var reqBody bytes.Buffer
	var writer *multipart.Writer
	if writer, recErr = sender.CreateFormFile(fileHandle, &reqBody); recErr != nil {
		return fmt.Errorf("failed create forl file: %e", recErr)
	}

	log.Printf("Sending file %s to recognize", filePath)

	var respData []byte
	targetURL := s.address + RecognitionURL
	mimeType := writer.FormDataContentType()
	respData, recErr = sender.SendRequest(&reqBody, &targetURL, &mimeType, s.timeout)
	if recErr != nil {
		return fmt.Errorf("failed send request: %e", recErr)
	}

	var resTest = &DocumentForm{}
	_ = json.Unmarshal(respData, resTest)

	if len(resTest.Context) == 0 {
		return fmt.Errorf("returned empty content data")
	}

	return nil
}

func (s *Service) GetProcessingJobs() map[string]*processing.ProcessJob {
	return make(map[string]*processing.ProcessJob)
}

func (s *Service) GetProcessingJob(_ string) *processing.ProcessJob {
	return nil
}
