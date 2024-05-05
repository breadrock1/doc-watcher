package langchain

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"doc-notifier/internal/pkg/tokenizer/forms"
	"doc-notifier/internal/pkg/tokenizer/tokoptions"
	"encoding/json"
	"log"
	"strings"
	"time"
)

const ServiceURL = "/api/v1/get_vectors"

type Service struct {
	address           string
	timeout           time.Duration
	ChunkSize         int
	ChunkOverlap      int
	ReturnChunkedText bool
}

func New(options *tokoptions.Options) *Service {
	return &Service{
		address:           options.Address,
		timeout:           options.Timeout,
		ChunkSize:         options.ChunkSize,
		ChunkOverlap:      options.ChunkOverlap,
		ReturnChunkedText: options.ChunkedFlag,
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
		ReturnChunkedText: s.ReturnChunkedText,
	}

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		log.Println("Failed while marshaling text vectors: ", err)
		return computedTokens, err
	}

	log.Printf("Sending file to extract tokens")
	reqBody := bytes.NewBuffer(jsonData)

	mimeType := "application/json"
	targetURL := s.address + ServiceURL
	respData, err := sender.SendRequest(reqBody, &targetURL, &mimeType, s.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return computedTokens, err
	}

	_ = json.Unmarshal(respData, computedTokens)
	return computedTokens, nil
}
