package sender

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"os"
)

const RecognitionURL = "/api/v1/extract_text"

// const RecognitionURL = "/api/assistant/extract-file/"

type DocumentForm struct {
	// Context string `json:"text"`
	Context string `json:"context"`
}

func (fs *FileSender) ReadRawFileData(filePath string) (string, error) {
	bytesData, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Failed while reading file: ", err)
		return "", err
	}

	return string(bytesData), nil
}

func (fs *FileSender) RecognizeFileData(filePath string) (string, error) {
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

	targetURL := fs.OrcServiceAddress + RecognitionURL
	log.Printf("Sending file %s to recognize", filePath)
	respData, err := fs.sendRequest(&reqBody, &targetURL)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "", err
	}

	var resTest = &DocumentForm{}
	_ = json.Unmarshal(respData, resTest)

	return resTest.Context, nil
}
