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
