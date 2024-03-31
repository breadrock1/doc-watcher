package tokenizer

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"log"
	"strings"
	"time"
)

const ServiceURL = "/api/v1/get_vectors"

type LangChainTokenizer struct {
	address           string
	timeout           time.Duration
	ChunkSize         int
	ChunkOverlap      int
	ReturnChunkedText bool
}

func NewLangChain(options *Options) *LangChainTokenizer {
	return &LangChainTokenizer{
		address:           options.Address,
		timeout:           options.Timeout,
		ChunkSize:         options.ChunkSize,
		ChunkOverlap:      options.ChunkOverlap,
		ReturnChunkedText: options.ChunkedFlag,
	}
}

func (lt *LangChainTokenizer) TokenizeTextData(content string) (*ComputedTokens, error) {
	computedTokens := &ComputedTokens{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}

	contentData := strings.ReplaceAll(content, "\n", " ")
	textVectors := &GetTokensForm{
		Text:              contentData,
		ChunkSize:         lt.ChunkSize,
		ChunkOverlap:      lt.ChunkOverlap,
		ReturnChunkedText: lt.ReturnChunkedText,
	}

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		log.Println("Failed while marshaling text vectors: ", err)
		return computedTokens, err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := lt.address + ServiceURL
	log.Printf("Sending file to extract tokens")

	mimeType := "application/json"
	respData, err := sender.SendRequest(reqBody, &targetURL, &mimeType, lt.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return computedTokens, err
	}

	_ = json.Unmarshal(respData, computedTokens)
	return computedTokens, nil
}
