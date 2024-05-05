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
	"sync"
	"time"
)

const CheckJobStatusTimeout = 5 * time.Second
const RecognitionURL = "/api/v2/text/create_extraction"
const GetResultURL = "/api/v2/text/get"

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
	targetURL := s.Address + RecognitionURL
	mimeType := writer.FormDataContentType()
	respData, recErr = sender.SendRequest(&reqBody, &targetURL, &mimeType, s.timeout)
	if recErr != nil {
		return fmt.Errorf("failed send request: %e", recErr)
	}

	var ocrJob = &OcrJob{}
	_ = json.Unmarshal(respData, ocrJob)
	document.OcrMetadata = s.launchAndAwait(ocrJob.JobId, document)

	return nil
}

func (s *Service) launchAndAwait(jobID string, document *reader.Document) *reader.OcrMetadata {
	var ocrResult *reader.OcrMetadata
	processingJob := &processing.ProcessJob{
		JobId:    jobID,
		Status:   false,
		Document: document,
	}

	s.mu.Lock()
	s.ProcessingJobs[jobID] = processingJob
	s.mu.Unlock()

	waitCh := make(chan *reader.OcrMetadata)
	go s.checkOcrJobStatus(jobID, waitCh)
	ocrResult = <-waitCh

	s.mu.Lock()
	s.ProcessingJobs[jobID].Status = true
	s.mu.Unlock()

	return ocrResult
}

func (s *Service) checkOcrJobStatus(jobId string, waitCh chan *reader.OcrMetadata) {
	getURLAddress := fmt.Sprintf("%s%s/%s", s.Address, GetResultURL, jobId)
	for {
		time.Sleep(CheckJobStatusTimeout)
		res, err := s.sendCheckJobStatus(getURLAddress, jobId)
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

func (s *Service) sendCheckJobStatus(targetURL string, jobId string) (*reader.OcrMetadata, *OcrJobError) {
	var ocrResult = &reader.OcrMetadata{}
	response, err := http.Get(targetURL)
	if err != nil {
		msg := fmt.Sprintf("Error while creating request: %s", err)
		return ocrResult, &OcrJobError{Type: FailedResponse, Message: msg}
	}

	if response.StatusCode == 202 {
		msg := fmt.Sprintf("Job '%s' are processing...", jobId)
		return ocrResult, &OcrJobError{Type: Processing, Message: msg}
	}

	if response.StatusCode > 210 {
		msg := fmt.Sprintf("Error response for job '%s'", jobId)
		return ocrResult, &OcrJobError{Type: FailedResponse, Message: msg}
	}

	log.Printf("Successful response for job '%s': %s", jobId, response.Status)
	defer func() { _ = response.Body.Close() }()

	var respData []byte
	if respData, err = io.ReadAll(response.Body); err != nil {
		msg := fmt.Sprintf("Failed while reading response reqBody: %s", err)
		return ocrResult, &OcrJobError{Type: FailedResponse, Message: msg}
	}

	_ = json.Unmarshal(respData, ocrResult)
	return ocrResult, nil
}

func (s *Service) GetProcessingJobs() map[string]*processing.ProcessJob {
	return s.ProcessingJobs
}

func (s *Service) GetProcessingJob(jobId string) *processing.ProcessJob {
	s.mu.Lock()
	value, ok := s.ProcessingJobs[jobId]
	s.mu.Unlock()

	if !ok {
		return nil
	}

	return value
}

func (s *Service) clearSuccessfulTasks() {
	var collectedJobs []string
	for key, value := range s.ProcessingJobs {
		if value.Status {
			collectedJobs = append(collectedJobs, key)
		}
	}

	s.mu.Lock()
	for _, jobId := range collectedJobs {
		delete(s.ProcessingJobs, jobId)
	}
	s.mu.Unlock()
}
