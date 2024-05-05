package reader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Service struct {
	mu           *sync.RWMutex
	AwaitingDocs map[string]*Document
}

func New() *Service {
	return &Service{
		mu:           &sync.RWMutex{},
		AwaitingDocs: make(map[string]*Document),
	}
}

func (s *Service) AddAwaitDocument(document *Document) {
	s.mu.Lock()
	s.AwaitingDocs[document.DocumentMD5] = document
	s.mu.Unlock()
}

func (s *Service) PopUnrecognizedDocument(documentID string) *Document {
	var document *Document
	s.mu.Lock()
	document, _ = s.AwaitingDocs[documentID]
	delete(s.AwaitingDocs, documentID)
	s.mu.Unlock()
	return document
}

func (s *Service) IsUnrecognizedDocument(documentID string) bool {
	s.mu.RLock()
	_, ok := s.AwaitingDocs[documentID]
	s.mu.RUnlock()
	return ok
}

func (s *Service) GetAwaitDocuments() []*Document {
	var awaitDocs []*Document
	for _, document := range s.AwaitingDocs {
		awaitDocs = append(awaitDocs, document)
	}
	return awaitDocs
}

func (s *Service) ParseCaughtFiles(filePath string) []*Document {
	mu := &sync.Mutex{}
	var customList []*Document

	wg := &sync.WaitGroup{}
	for _, filePath := range getEntityFiles(filePath) {
		wg.Add(1)
		filePath := filePath

		go func() {
			defer wg.Done()

			if doc, err := ParseFile(filePath); err == nil {
				log.Println("Caught parsed document: ", doc.DocumentName)
				mu.Lock()
				customList = append(customList, doc)
				mu.Unlock()
				return
			}

			log.Println("Failed parsing document: ", filePath)
		}()
	}

	wg.Wait()

	return customList
}

func (s *Service) MoveFileToUnrecognized(document *Document) error {
	outputFilePath := fmt.Sprintf("%s/%s", "./indexer/unrecognized", document.DocumentName)
	if err := moveFileTo(document.DocumentPath, outputFilePath); err != nil {
		return err
	}

	return nil
}

func (s *Service) MoveFileToDir(filePath string, targetDir string) error {
	_, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(targetDir, os.ModePerm)
	}

	_, fileName := filepath.Split(filePath)
	outputFilePath := fmt.Sprintf("%s/%s", targetDir, fileName)
	if err = moveFileTo(filePath, outputFilePath); err != nil {
		return err
	}

	return nil
}
