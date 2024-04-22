package logoper

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
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

const RecognitionURL = "/api/v2/text/create_extraction"
const GetResultURL = "/api/v2/text/get/"

type OcrJob struct {
	JobId string `json:"job_id"`
}

type OcrResult struct {
	DocType string `json:"doc_type"`
	Text    string `json:"text"`
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

	var resTest = &OcrJob{}
	_ = json.Unmarshal(respData, resTest)

	waitCh := make(chan string)
	go do.awaitOcrResult(resTest.JobId, waitCh)
	result := <-waitCh

	return result, nil
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

	var resTest = &OcrJob{}
	_ = json.Unmarshal(respData, resTest)

	waitCh := make(chan string)
	go do.awaitOcrResult(resTest.JobId, waitCh)
	result := <-waitCh

	return result, nil
}

func (do *Service) awaitOcrResult(jobId string, waitCh chan string) {
	getURLAddress := do.address + RecognitionURL + jobId
	for {
		res, err := do.checkOcrJobStatus(getURLAddress)

		if err != nil {
			waitCh <- res
			break
		}

		time.Sleep(5 * time.Second)
	}
}

func (do *Service) checkOcrJobStatus(targetURL string) (string, error) {
	response, err := http.Get(targetURL)
	if err != nil {
		log.Println("Error while creating request:", err)
		return "", err
	}

	if response.StatusCode > 200 {
		log.Printf("Non Ok response status: %s", response.Status)
		defer func() { _ = response.Body.Close() }()
		//respData, err := io.ReadAll(response.Body)
		//if err != nil {
		//	log.Println("Failed while reading response reqBody: ", err)
		//	return "", err
		//}
		return "", err
	}

	return "", nil
}
