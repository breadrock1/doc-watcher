package tokenizer

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"log"
	"strings"
)

const EmbeddingsLCURL = "/api/v1/vectorizer/get_vectors"

type LangChainTokenizer struct {
	address           string
	ChunkSize         int
	ChunkOverlap      int
	ReturnChunkedText bool
}

func NewLangChain(options *Options) *LangChainTokenizer {
	return &LangChainTokenizer{
		address:           options.Address,
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
	textVectors := &TokenizerForm{
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
	targetURL := lt.address + EmbeddingsLCURL
	log.Printf("Sending file to extract tokens")
	respData, err := sender.SendRequest(reqBody, &targetURL)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return computedTokens, err
	}

	_ = json.Unmarshal(respData, computedTokens)
	return computedTokens, nil
}
