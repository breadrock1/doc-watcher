package assistant

import (
	"bytes"
	"encoding/json"
	"log"
	"math"
	"strings"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/sender"
	"doc-notifier/internal/tokenizer/forms"
)

const EmbeddingsAssistantURL = "/embed"

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
		Chunks:      1,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}

	contentData := strings.ReplaceAll(content, "\n", " ")
	chunkedText := s.splitContent(contentData, s.ChunkSize)
	for _, textData := range chunkedText {
		if tokens, err := s.loadTextDataTokens(textData); err == nil {
			computedTokens.Chunks++
			computedTokens.Vectors = append(computedTokens.Vectors, tokens)
			computedTokens.ChunkedText = append(computedTokens.ChunkedText, textData)
		}
	}

	return computedTokens, nil
}

func (s *Service) loadTextDataTokens(content string) ([]float64, error) {
	textVectors := &EmbedAllForm{
		Inputs:   content,
		Truncate: false,
	}

	var tokenErr error
	var jsonData []byte
	if jsonData, tokenErr = json.Marshal(textVectors); tokenErr != nil {
		log.Println("Failed while marshaling text vectors: ", tokenErr)
		return []float64{}, tokenErr
	}

	log.Printf("Sending file to extract tokens")
	reqBody := bytes.NewBuffer(jsonData)

	method := "PUT"
	mimeType := "application/json"
	targetURL := s.address + EmbeddingsAssistantURL
	respData, err := sender.SendRequest(reqBody, &targetURL, &method, &mimeType, s.timeout)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return []float64{}, err
	}

	tokensDense := &[][]float64{}
	_ = json.Unmarshal(respData, tokensDense)

	tmp := *tokensDense
	return tmp[0], nil
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
