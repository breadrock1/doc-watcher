package searcher

import (
	"bytes"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"log"
)

const SearcherURL = "/document/new"

type Service struct {
	address string
}

func New(address string) *Service {
	return &Service{
		address: address,
	}
}

func (ss *Service) StoreDocument(document *reader.Document) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := ss.address + SearcherURL
	log.Printf("Storing document %s to elastic", document.DocumentPath)
	_, err = sender.SendRequest(reqBody, &targetURL, "application/json")
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	return nil
}
