package searcher

import (
	"bytes"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"log"
	"time"
)

const ServiceURL = "/documents/create"

type Service struct {
	Address string
	timeout time.Duration
}

func New(address string, timeout time.Duration) *Service {
	return &Service{
		Address: address,
		timeout: timeout,
	}
}

func (s *Service) StoreDocument(document *reader.Document) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := s.Address + ServiceURL
	log.Printf("Storing document %s to elastic", document.FolderID)

	mimeType := "application/json"
	_, err = sender.SendRequest(reqBody, &targetURL, &mimeType, s.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	return nil
}
