package logoper

import (
	"bytes"
	"doc-notifier/internal/pkg/ocr/processing"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Service struct {
	mu             *sync.Mutex
	Address        string
	timeout        time.Duration
	ProcessingJobs map[string]*processing.ProcessJob
}

func New(address string, timeout time.Duration) *Service {
	service := &Service{
		mu:             &sync.Mutex{},
		Address:        address,
		timeout:        timeout,
		ProcessingJobs: make(map[string]*processing.ProcessJob),
	}

	time.AfterFunc(1*time.Hour, service.clearSuccessfulTasks)
	return service
}

const RecognitionURL = "/api/v2/text/create_extraction"
const GetResultURL = "/api/v2/text/get/"

type OcrJobErrorType int

const (
	Processing OcrJobErrorType = iota
	FailedResponse
)

type OcrJobError struct {
	Type    OcrJobErrorType
	Message string
}

type OcrJob struct {
	JobId string `json:"job_id"`
}

func (do *Service) RecognizeFile(document *reader.Document) (string, error) {
	filePath := document.DocumentPath
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

	targetURL := do.Address + RecognitionURL
	log.Printf("Sending file %s to recognize", filePath)

	mimeType := writer.FormDataContentType()
	respData, err := sender.SendRequest(&reqBody, &targetURL, &mimeType, do.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "", err
	}

	var ocrJob = &OcrJob{}
	_ = json.Unmarshal(respData, ocrJob)

	jaba := &processing.ProcessJob{
		JobId:    ocrJob.JobId,
		Status:   false,
		Document: document,
	}

	do.mu.Lock()
	do.ProcessingJobs[ocrJob.JobId] = jaba
	do.mu.Unlock()

	waitCh := make(chan *reader.OcrResult)
	go do.awaitOcrResult(ocrJob.JobId, waitCh)
	result := <-waitCh

	do.mu.Lock()
	do.ProcessingJobs[ocrJob.JobId].Status = true
	do.mu.Unlock()

	document.OcrMetadata = result
	return result.Text, nil
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

	targetURL := do.Address + RecognitionURL
	log.Printf("Sending file to recognize file data")

	mimeType := writer.FormDataContentType()
	respData, err := sender.SendRequest(&reqBody, &targetURL, &mimeType, do.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "", err
	}

	var resTest = &OcrJob{}
	_ = json.Unmarshal(respData, resTest)

	waitCh := make(chan *reader.OcrResult)
	go do.awaitOcrResult(resTest.JobId, waitCh)
	result := <-waitCh

	return result.Text, nil
}

func (do *Service) awaitOcrResult(jobId string, waitCh chan *reader.OcrResult) {
	getURLAddress := do.Address + GetResultURL + jobId
	for {
		time.Sleep(5 * time.Second)
		res, err := do.checkOcrJobStatus(getURLAddress, jobId)
		if err != nil {
			log.Println(err.Message)
			switch err.Type {
			case Processing:
				continue
			case FailedResponse:
				waitCh <- res
				break
			}
		}

		waitCh <- res
		break
	}
}

func (do *Service) checkOcrJobStatus(targetURL string, jobId string) (*reader.OcrResult, *OcrJobError) {
	var ocrResult = &reader.OcrResult{}
	response, err := http.Get(targetURL)
	if err != nil {
		msg := fmt.Sprintf("Error while creating request: %s", err)
		return ocrResult, &OcrJobError{Type: FailedResponse, Message: msg}
	}

	if response.StatusCode == 202 {
		msg := fmt.Sprintf("Job '%s' are processing...", jobId)
		return ocrResult, &OcrJobError{Type: Processing, Message: msg}
	}

	if response.StatusCode > 200 {
		msg := fmt.Sprintf("Error response for job '%s'", jobId)
		return ocrResult, &OcrJobError{Type: FailedResponse, Message: msg}
	}

	log.Printf("Successful response for job '%s': %s", jobId, response.Status)
	defer func() { _ = response.Body.Close() }()
	respData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed while reading response reqBody: ", err)
		msg := fmt.Sprintf("Failed while reading response reqBody: %s", err)
		return ocrResult, &OcrJobError{Type: FailedResponse, Message: msg}
	}

	_ = json.Unmarshal(respData, ocrResult)
	return ocrResult, nil
}

func (do *Service) GetProcessingJobs() map[string]*processing.ProcessJob {
	return do.ProcessingJobs
}

func (do *Service) GetProcessingJob(jobId string) *processing.ProcessJob {
	do.mu.Lock()
	value, ok := do.ProcessingJobs[jobId]
	do.mu.Unlock()

	if !ok {
		return nil
	}

	return value
}

func (do *Service) clearSuccessfulTasks() {
	var collectedJobs []string
	for key, value := range do.ProcessingJobs {
		if value.Status {
			collectedJobs = append(collectedJobs, key)
		}
	}

	do.mu.Lock()
	for _, jobId := range collectedJobs {
		delete(do.ProcessingJobs, jobId)
	}
	do.mu.Unlock()
}
