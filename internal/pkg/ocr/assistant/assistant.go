package assistant

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"os"
)

type AssistantOCR struct {
	address string
}

func New(address string) *AssistantOCR {
	return &AssistantOCR{
		address: address,
	}
}

const RecognitionURL = "/api/assistant/extract-file/"

type documentForm struct {
	Context string `json:"context"`
}

func (ro *AssistantOCR) RecognizeFile(filePath string) (string, error) {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		log.Println("Failed while opening file: ", err)
		return "", err
	}
	defer func() { _ = fileHandle.Close() }()

	var reqBody bytes.Buffer
	writer := multipart.NewWriter(&reqBody)
	part, err := writer.CreateFormFile("file", filePath)
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

	targetURL := ro.address + RecognitionURL
	log.Printf("Sending file %s to recognize", filePath)
	respData, err := sender.SendRequest(&reqBody, &targetURL, writer.FormDataContentType())
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "", err
	}

	var resTest = &documentForm{}
	_ = json.Unmarshal(respData, resTest)

	return resTest.Context, nil
}

func (ro *AssistantOCR) RecognizeFileData(data []byte) (string, error) {
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

	targetURL := ro.address + RecognitionURL
	log.Printf("Sending file to recognize file data")
	respData, err := sender.SendRequest(&reqBody, &targetURL, writer.FormDataContentType())
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "", err
	}

	var resTest = &documentForm{}
	_ = json.Unmarshal(respData, resTest)

	return resTest.Context, nil
}
