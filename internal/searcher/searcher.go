package searcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"doc-watcher/internal/sender"
	"doc-watcher/internal/watcher"
	"github.com/labstack/echo/v4"
)

const DocumentJsonMime = echo.MIMEApplicationJSON

type Service struct {
	config *Config
}

func New(config *Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) StoreDocument(doc *watcher.Document) error {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed while marshaling doc: %w", err)
	}

	buildURL := strings.Builder{}
	buildURL.WriteString(sender.GetHttpSchema(s.config.EnableSSL))
	buildURL.WriteString("://")
	buildURL.WriteString(s.config.Address)
	buildURL.WriteString("/storage/folders/")
	buildURL.WriteString(doc.FolderID)
	buildURL.WriteString("/documents/")
	buildURL.WriteString(doc.DocumentID)
	targetURL := buildURL.String()

	log.Printf("storing document %s to index %s", doc.DocumentID, doc.FolderID)

	reqBody := bytes.NewBuffer(jsonData)
	timeoutReq := time.Duration(300) * time.Second
	_, err = sender.PUT(reqBody, targetURL, DocumentJsonMime, timeoutReq)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) StoreVector(doc *watcher.Document) error {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed while marshaling doc: %w", err)
	}

	folderID := fmt.Sprintf("%s-vector", doc.FolderID)

	buildURL := strings.Builder{}
	buildURL.WriteString(sender.GetHttpSchema(s.config.EnableSSL))
	buildURL.WriteString("://")
	buildURL.WriteString(s.config.Address)
	buildURL.WriteString("/storage/folders/")
	buildURL.WriteString(folderID)
	buildURL.WriteString("/documents/")
	buildURL.WriteString(doc.DocumentID)
	buildURL.WriteString("?document_type=vectors")
	targetURL := buildURL.String()

	log.Printf("storing document %s to index %s", doc.DocumentID, doc.FolderID)

	reqBody := bytes.NewBuffer(jsonData)
	timeoutReq := time.Duration(300) * time.Second
	_, err = sender.PUT(reqBody, targetURL, DocumentJsonMime, timeoutReq)
	if err != nil {
		return err
	}

	return nil
}
