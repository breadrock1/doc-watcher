package searcher

import (
	"bytes"
	"doc-notifier/internal/config"
	"doc-notifier/internal/reader"
	"doc-notifier/internal/sender"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type Service struct {
	Address    string
	Timeout    time.Duration
	FoldersMap map[string]*Folder
}

func New(config *config.SearcherConfig) *Service {
	return &Service{
		Address:    config.Address,
		Timeout:    config.Timeout,
		FoldersMap: make(map[string]*Folder, 0),
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

type Folder struct {
	Health       string `json:"health"`
	Status       string `json:"status"`
	ID           string `json:"id"`
	UUID         string `json:"uuid"`
	Pri          string `json:"pri"`
	Rep          string `json:"rep"`
	DocsCount    string `json:"docs_count"`
	DocsDeleted  string `json:"docs_deleted"`
	StoreSize    string `json:"store_size"`
	PriStoreSize string `json:"pri_store_size"`
	Name         string `json:"name"`
}

func (s *Service) GetFolderID(folderName string) (string, error) {
	folder, exists := s.FoldersMap[folderName]
	if exists {
		return folder.ID, nil
	}

	method := "GET"
	mimeType := "application/json"
	targetURL := "/storage/folders"
	respData, err := sender.SendRequest(nil, &targetURL, &method, &mimeType, s.Timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return "unrecognized", err
	}

	var folders []*Folder
	if err := json.Unmarshal(respData, &folders); err != nil {
		log.Println("Failed while reading response reqBody: ", err)
		return "unrecognized", err
	}

	for _, fold := range folders {
		s.FoldersMap[fold.Name] = fold
	}

	folder, exists = s.FoldersMap[folderName]
	if exists {
		return folder.ID, nil
	}

	return "unrecognized", errors.New("not exists")
}
