package searcher

import (
	"bytes"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"time"
)

const CreateDocumentURL = "/document/new"
const CreateBucketURL = "/bucket/new"
const DeleteBucketURL = "/bucket/"

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

func (ss *Service) StoreDocument(document *reader.Document) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := ss.address + CreateDocumentURL
	log.Printf("Storing document %s to elastic", document.DocumentPath)

	mimeType := "application/json"
	_, err = sender.SendRequest(reqBody, &targetURL, &mimeType, ss.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	return nil
}

type BucketForm struct {
	BucketName string `json:"bucket_name"`
}

func (ss *Service) CreateBucket(bucketName string) error {
	bucketForm := &BucketForm{BucketName: bucketName}
	jsonData, err := json.Marshal(bucketForm)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := ss.address + CreateBucketURL
	log.Printf("Creating bucket %s into elastic", bucketName)

	mimeType := "application/json"
	_, err = sender.SendRequest(reqBody, &targetURL, &mimeType, ss.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	return nil
}

func (ss *Service) DeleteBucket(bucketName string) error {
	bucketForm := &BucketForm{BucketName: bucketName}
	jsonData, err := json.Marshal(bucketForm)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := ss.address + DeleteBucketURL + bucketName

	req, err := http.NewRequest("DELETE", targetURL, reqBody)
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	mimeType := "application/json"
	req.Header.Set(echo.HeaderContentType, mimeType)
	client := &http.Client{Timeout: ss.timeout}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error while creating request:", err)
		return err
	}
	defer func() { _ = response.Body.Close() }()

	respData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed while reading response reqBody: ", err)
		return err
	}

	if response.StatusCode > 200 {
		log.Printf("Non Ok response status %s: %s", response.Status, string(respData))
		return errors.New("non 200 response code status")
	}

	return nil
}
