package sovaocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"doc-watcher/internal/ocr"
	"doc-watcher/internal/sender"
	"doc-watcher/internal/watcher"
)

const RecognitionURL = "/ocr_extract_text"

type Service struct {
	config *ocr.Config
}

func New(config *ocr.Config) *ocr.Service {
	servClient := &Service{
		config: config,
	}

	return &ocr.Service{
		Ocr: servClient,
	}
}

func (s *Service) RecognizeFile(document *watcher.Document, data []byte) error {
	var buf bytes.Buffer

	mpw := multipart.NewWriter(&buf)
	fileForm, err := mpw.CreateFormFile("file", document.DocumentName)
	if err != nil {
		return err
	}

	if _, err = fileForm.Write(data); err != nil {
		return err
	}

	if err = mpw.Close(); err != nil {
		return err
	}

	log.Printf("sending file %s to recognize", document.DocumentName)

	mimeType := mpw.FormDataContentType()
	timeoutReq := s.config.Timeout * time.Second
	targetURL := sender.BuildTargetURL(s.config.EnableSSL, s.config.Address, RecognitionURL)

	respData, err := sender.POST(&buf, targetURL, mimeType, timeoutReq)
	if err != nil {
		return fmt.Errorf("failed send request: %w", err)
	}

	var resTest ocr.DocumentForm
	_ = json.Unmarshal(respData, &resTest)
	document.SetContentData(resTest.Content)
	document.SetOcrMetadata(watcher.DefaultOcr())

	if len(resTest.Content) == 0 {
		return fmt.Errorf("returned empty content data")
	}

	document.SetQuality(10000)
	return nil
}
