package searcher

import (
	"bytes"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"log"
)

const SearcherURL = "/document/new"

type SearcherService struct {
	address string
}

func New(address string) *SearcherService {
	return &SearcherService{
		address: address,
	}
}

func (ss *SearcherService) StoreDocument(document *reader.Document) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := ss.address + SearcherURL
	log.Printf("Storing document %s to elastic", document.DocumentPath)
	_, err = sender.SendRequest(reqBody, &targetURL)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	return nil
}
