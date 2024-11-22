package sovavec

import (
	"bytes"
	"doc-watcher/internal/watcher"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"doc-watcher/internal/embeddings"
	"doc-watcher/internal/sender"
	"github.com/labstack/echo/v4"
)

const EmbeddingsAssistantURL = "/embed"

type Service struct {
	config *embeddings.Config
}

func New(config *embeddings.Config) *embeddings.Service {
	servClient := &Service{
		config: config,
	}

	return &embeddings.Service{
		Tokenizer: servClient,
	}
}

func (s *Service) Tokenize(doc *watcher.Document) (*embeddings.ComputeTokens, error) {
	computedTokens := &embeddings.ComputeTokens{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}

	contentData := strings.ReplaceAll(doc.Content, "\n", " ")
	chunkedText := s.splitContent(contentData, s.config.ChunkSize)
	for _, textData := range chunkedText {
		tokens, err := s.loadTextDataTokens(textData)
		if err != nil {
			log.Printf("failed to load embeddings for %s: %v", doc.DocumentName, err)
			continue
		}

		computedTokens.Chunks++
		computedTokens.Vectors = append(computedTokens.Vectors, tokens)
		computedTokens.ChunkedText = append(computedTokens.ChunkedText, textData)
	}

	return computedTokens, nil
}

func (s *Service) loadTextDataTokens(content string) ([]float64, error) {
	textVectors := &embeddings.EmbedAllForm{
		Inputs:    content,
		Truncate:  false,
		Normalize: true,
	}

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		return []float64{}, fmt.Errorf("failed to marshal tokens: %w", err)
	}

	log.Printf("sending chunk to generate embeddings")

	reqBody := bytes.NewBuffer(jsonData)
	mimeType := echo.MIMEApplicationJSON
	timeoutReq := time.Duration(300) * time.Second
	targetURL := sender.BuildTargetURL(s.config.EnableSSL, s.config.Address, EmbeddingsAssistantURL)
	respData, err := sender.POST(reqBody, targetURL, mimeType, timeoutReq)
	if err != nil {
		return []float64{}, fmt.Errorf("failed to load embeddings: %w", err)
	}

	tokens := make([][]float64, 0)
	_ = json.Unmarshal(respData, &tokens)

	if len(tokens) < 1 {
		log.Println("returned empty tokens from tokenizer service")
		return make([]float64, 0), nil
	}

	return tokens[0], nil
}

func (s *Service) splitContent(content string, chunkSize int) []string {
	strLength := len(content)
	splitLength := int(math.Ceil(float64(strLength) / float64(chunkSize)))
	splitString := make([]string, splitLength)
	var start, stop int
	for i := 0; i < splitLength; i++ {
		start = i * chunkSize
		stop = start + chunkSize
		if stop > strLength {
			stop = strLength
		}

		splitString[i] = content[start:stop]
	}

	return splitString
}
