package assistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/models"
	"doc-notifier/internal/sender"
)

const RecognitionURL = "/ocr_extract_text"

type Service struct {
	address string
	timeout time.Duration
}

func New(config *config.OcrConfig) *Service {
	timeoutReq := config.Timeout * time.Second
	return &Service{
		address: config.Address,
		timeout: timeoutReq,
	}
}

func (s *Service) RecognizeFile(document *models.Document, filePath string) error {
	var recErr error
	var fileHandle *os.File
	if fileHandle, recErr = os.Open(filePath); recErr != nil {
		return fmt.Errorf("file %s not found: %e", filePath, recErr)
	}
	defer func() { _ = fileHandle.Close() }()

	var reqBody bytes.Buffer
	var writer *multipart.Writer
	if writer, recErr = sender.CreateFormFile(fileHandle, &reqBody); recErr != nil {
		return fmt.Errorf("failed create form file: %e", recErr)
	}

	log.Printf("Sending file %s to recognize", filePath)

	var respData []byte
	targetURL := s.address + RecognitionURL
	mimeType := writer.FormDataContentType()
	if respData, recErr = sender.POST(&reqBody, targetURL, mimeType, s.timeout); recErr != nil {
		return fmt.Errorf("failed send request: %e", recErr)
	}

	var resTest = &DocumentForm{}
	_ = json.Unmarshal(respData, resTest)
	document.SetContentData(resTest.Context)
	document.SetOcrMetadata(models.DefaultOcr())

	if len(resTest.Context) == 0 {
		return fmt.Errorf("returned empty content data")
	}

	document.QualityRecognized = 10000
	return nil
}
