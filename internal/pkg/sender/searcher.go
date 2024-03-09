package sender

import (
	"bytes"
	"doc-notifier/internal/pkg/reader"
	"encoding/json"
	"log"
)

const SearcherURL = "/document/new"

type DocumentForm struct {
	// Context string `json:"text"`
	Context string `json:"context"`
}

func (fs *FileSender) StoreDocument(document *reader.Document) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := fs.SearcherAddress + SearcherURL
	log.Printf("Storing document %s to elastic", document.DocumentPath)
	_, err = fs.sendRequest(reqBody, &targetURL)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	return nil
}
