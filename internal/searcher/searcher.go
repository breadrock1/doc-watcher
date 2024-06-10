package searcher

import (
	"bytes"
	"doc-notifier/internal/config"
	"doc-notifier/internal/reader"
	"doc-notifier/internal/sender"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Service struct {
	Address string
	Timeout time.Duration
}

func New(config *config.SearcherConfig) *Service {
	return &Service{
		Address: config.Address,
		Timeout: config.Timeout,
	}
}

func (s *Service) StoreDocument(document *reader.Document) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := fmt.Sprintf("/storage/folders/%s/documents/%s", document.FolderID, document.DocumentMD5)
	log.Printf("Storing document %s to elastic", document.FolderID)

	method := "PUT"
	mimeType := "application/json"
	_, err = sender.SendRequest(reqBody, &targetURL, &method, &mimeType, s.Timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	targetURL = fmt.Sprintf("/storage/folders/%s/documents/%s", "history", document.DocumentMD5)
	_, err = sender.SendRequest(reqBody, &targetURL, &method, &mimeType, s.Timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	targetURL = fmt.Sprintf("/storage/folders/%s/documents/%s?document_type=vectors", document.FolderID, document.DocumentMD5)
	_, err = sender.SendRequest(reqBody, &targetURL, &method, &mimeType, s.Timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	return nil
}
