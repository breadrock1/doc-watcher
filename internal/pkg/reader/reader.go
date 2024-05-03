package reader

import "sync"

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

func (s *Service) PopDoneDocument(documentID string) {
	s.mu.Lock()
	delete(s.AwaitingDocs, documentID)
	s.mu.Unlock()
}

func (s *Service) GetAwaitDocument(documentID string) *Document {
	var awaitDoc *Document
	s.mu.RLock()
	if document, ok := s.AwaitingDocs[documentID]; ok {
		awaitDoc = document
	}
	s.mu.RUnlock()
	return awaitDoc
}

func (s *Service) GetAwaitDocuments() []*Document {
	var awaitDocs []*Document
	for _, document := range s.AwaitingDocs {
		awaitDocs = append(awaitDocs, document)
	}
	return awaitDocs
}
