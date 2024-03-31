package searcher

import (
	"bytes"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"log"
	"time"
)

const ServiceURL = "/document/new"

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
	targetURL := ss.address + ServiceURL
	log.Printf("Storing document %s to elastic", document.DocumentPath)

	mimeType := "application/json"
	_, err = sender.SendRequest(reqBody, &targetURL, &mimeType, ss.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	return nil
}
