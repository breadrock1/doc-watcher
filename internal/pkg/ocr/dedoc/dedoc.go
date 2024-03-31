package dedoc

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

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

const RecognitionURL = "/api/v1/extract_text"

type DocumentForm struct {
	Context string `json:"text"`
}

func (do *Service) RecognizeFile(filePath string) (string, error) {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		log.Println("Failed while opening file: ", err)
		return "", err
	}
	defer func() { _ = fileHandle.Close() }()

	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		log.Println("Failed while creating form file: ", err)
		return "", err
	}

	if _, err = io.Copy(part, fileHandle); err != nil {
		log.Println("Failed while coping file form part to file handle: ", err)
		return "", err
	}

	if err := writer.Close(); err != nil {
		log.Println("Failed while closing req body writer: ", err)
		return "", err
	}

	targetURL := do.address + RecognitionURL
	log.Printf("Sending file %s to recognize", filePath)

	mimeType := writer.FormDataContentType()
	respData, err := sender.SendRequest(&reqBody, &targetURL, &mimeType, do.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "", err
	}

	var resTest = &DocumentForm{}
	_ = json.Unmarshal(respData, resTest)

	if len(resTest.Context) == 0 {
		log.Println("Failed: returned empty string data...")
		return "", err
	}

	return resTest.Context, nil
}

func (do *Service) RecognizeFileData(data []byte) (string, error) {
	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)

	part, _ := writer.CreateFormField("file")
	_, err := part.Write(data)
	if err != nil {
		log.Println("Failed while creating form file: ", err)
		return "", err
	}

	if err := writer.Close(); err != nil {
		log.Println("Failed while closing req body writer: ", err)
		return "", err
	}

	targetURL := do.address + RecognitionURL
	log.Printf("Sending file to recognize file data")

	mimeType := writer.FormDataContentType()
	respData, err := sender.SendRequest(&reqBody, &targetURL, &mimeType, do.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "", err
	}

	var resTest = &DocumentForm{}
	_ = json.Unmarshal(respData, resTest)

	return resTest.Context, nil
}
