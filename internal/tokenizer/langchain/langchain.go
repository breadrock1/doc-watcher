package langchain

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/sender"
	"doc-notifier/internal/tokenizer/forms"
)

const ServiceURL = "/api/v1/get_vectors"

type Service struct {
	address      string
	timeout      time.Duration
	ChunkSize    int
	ChunkOverlap int
	ReturnChunks bool
}

func New(config *config.TokenizerConfig) *Service {
	return &Service{
		address:      config.Address,
		timeout:      config.Timeout,
		ChunkSize:    config.ChunkSize,
		ChunkOverlap: config.ChunkOverlap,
		ReturnChunks: config.ReturnChunks,
	}
}

func (s *Service) TokenizeTextData(content string) (*forms.ComputedTokens, error) {
	computedTokens := &forms.ComputedTokens{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}

	contentData := strings.ReplaceAll(content, "\n", " ")
	textVectors := &GetTokensForm{
		Text:              contentData,
		ChunkSize:         s.ChunkSize,
		ChunkOverlap:      s.ChunkOverlap,
		ReturnChunkedText: s.ReturnChunks,
	}

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		log.Println("Failed while marshaling text vectors: ", err)
		return computedTokens, err
	}

	log.Printf("Sending file to extract tokens")
	reqBody := bytes.NewBuffer(jsonData)

	method := "PUT"
	mimeType := "application/json"
	targetURL := s.address + ServiceURL
	respData, err := sender.SendRequest(reqBody, &targetURL, &method, &mimeType, s.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return computedTokens, err
	}

	_ = json.Unmarshal(respData, computedTokens)
	return computedTokens, nil
}
