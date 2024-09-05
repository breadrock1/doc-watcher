package searcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/models"
	"doc-notifier/internal/sender"
	"github.com/labstack/echo/v4"
)

type Service struct {
	Address string
	Timeout time.Duration
}

func New(config *config.SearcherConfig) *Service {
	timeout := config.Timeout * time.Second
	return &Service{
		Address: config.Address,
		Timeout: timeout,
	}
}

func (s *Service) StoreDocument(document *models.Document) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := fmt.Sprintf("%s/storage/folders/%s/documents/%s", s.Address, document.FolderID, document.DocumentID)
	log.Printf("Storing document %s to elastic", document.FolderID)

	mimeType := echo.MIMEApplicationJSON
	_, err = sender.PUT(reqBody, targetURL, mimeType, s.Timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return err
	}

	reqBody = bytes.NewBuffer(jsonData)
	folderID := fmt.Sprintf("%s-vector", document.FolderID)
	targetURL = fmt.Sprintf("%s/storage/folders/%s/documents/%s?document_type=vectors", s.Address, folderID, document.DocumentID)
	_, err = sender.PUT(reqBody, targetURL, mimeType, s.Timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
	}

	return nil
}
