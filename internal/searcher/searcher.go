package searcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/reader"
	"doc-notifier/internal/sender"
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
		FoldersMap: make(map[string]*Folder),
	}
}

func (s *Service) StoreDocument(document *reader.Document) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := fmt.Sprintf("%s/storage/folders/%s/documents/%s", s.Address, document.FolderID, document.DocumentMD5)
	log.Printf("Storing document %s to elastic", document.FolderID)

	method := "PUT"
	mimeType := "application/json"
	_, err = sender.SendRequest(reqBody, &targetURL, &method, &mimeType, s.Timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	reqBody = bytes.NewBuffer(jsonData)
	targetURL = fmt.Sprintf("%s/storage/folders/%s/documents/%s", s.Address, "history", document.DocumentMD5)
	_, err = sender.SendRequest(reqBody, &targetURL, &method, &mimeType, s.Timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
	}

	reqBody = bytes.NewBuffer(jsonData)
	targetURL = fmt.Sprintf("%s/storage/folders/%s/documents/%s?document_type=vectors", s.Address, document.FolderID, document.DocumentMD5)
	_, err = sender.SendRequest(reqBody, &targetURL, &method, &mimeType, s.Timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
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
	folder, exists := s.FoldersMap[strings.ToLower(folderName)]
	if exists {
		return folder.ID, nil
	}

	targetURL := fmt.Sprintf("%s%s", s.Address, "/storage/folders")
	respData, err := sender.SendGETRequest(targetURL)
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
		s.FoldersMap[strings.ToLower(fold.Name)] = fold
	}

	folder, exists = s.FoldersMap[folderName]
	if exists {
		return folder.ID, nil
	}

	return "unrecognized", errors.New("not exists")
}
