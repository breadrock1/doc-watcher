package tokenizer

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"log"
	"strings"
)

const EmbeddingsAssistantURL = "/api/v1/vectorizer/get_vectors"

type AssistantTokenizer struct {
	address           string
	ChunkSize         int
	ChunkOverlap      int
	ReturnChunkedText bool
}

func NewAssistant(options *Options) *AssistantTokenizer {
	return &AssistantTokenizer{
		address:           options.Address,
		ChunkSize:         options.ChunkSize,
		ChunkOverlap:      options.ChunkOverlap,
		ReturnChunkedText: options.ChunkedFlag,
	}
}

func (at *AssistantTokenizer) TokenizeTextData(content string) (*ComputedTokens, error) {
	computedTokens := &ComputedTokens{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}

	contentData := strings.ReplaceAll(content, "\n", " ")
	textVectors := &GetTokensForm{
		Text:              contentData,
		ChunkSize:         at.ChunkSize,
		ChunkOverlap:      at.ChunkOverlap,
		ReturnChunkedText: at.ReturnChunkedText,
	}

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		log.Println("Failed while marshaling text vectors: ", err)
		return computedTokens, err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := at.address + EmbeddingsAssistantURL
	log.Printf("Sending file to extract tokens")
	respData, err := sender.SendRequest(reqBody, &targetURL, "application/json")
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return computedTokens, err
	}

	_ = json.Unmarshal(respData, computedTokens)
	return computedTokens, nil
}
