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
)

type DedocOCR struct {
	address string
}

func New(address string) *DedocOCR {
	return &DedocOCR{
		address: address,
	}
}

const RecognitionURL = "/api/v1/extract_text"

type DocumentForm struct {
	Context string `json:"text"`
}

func (do *DedocOCR) RecognizeFile(filePath string) (string, error) {
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
	respData, err := sender.SendRequest(&reqBody, &targetURL, writer.FormDataContentType())
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "", err
	}

	var resTest = &DocumentForm{}
	_ = json.Unmarshal(respData, resTest)

	return resTest.Context, nil
}

func (do *DedocOCR) RecognizeFileData(data []byte) (string, error) {
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
	respData, err := sender.SendRequest(&reqBody, &targetURL, writer.FormDataContentType())
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "", err
	}

	var resTest = &DocumentForm{}
	_ = json.Unmarshal(respData, resTest)

	return resTest.Context, nil
}
