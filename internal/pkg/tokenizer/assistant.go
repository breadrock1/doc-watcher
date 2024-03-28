package tokenizer

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"log"
	"math"
	"strings"
)

const EmbeddingsAssistantURL = "/embed"

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

type EmbedAllForm struct {
	Inputs   string `json:"inputs"`
	Truncate bool   `json:"truncate"`
}

func (at *AssistantTokenizer) TokenizeTextData(content string) (*ComputedTokens, error) {
	computedTokens := &ComputedTokens{
		Chunks:      1,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}

	contentData := strings.ReplaceAll(content, "\n", " ")
	chunkedText := at.splitContent(contentData, at.ChunkSize)
	for _, textData := range chunkedText {
		if tokens, err := at.loadTextDataTokens(textData); err == nil {
			computedTokens.Chunks++
			computedTokens.Vectors = append(computedTokens.Vectors, tokens)
			computedTokens.ChunkedText = append(computedTokens.ChunkedText, textData)
		}
	}

	return computedTokens, nil
}

func (at *AssistantTokenizer) loadTextDataTokens(content string) ([]float64, error) {
	textVectors := &EmbedAllForm{
		Inputs:   content,
		Truncate: false,
	}

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		log.Println("Failed while marshaling text vectors: ", err)
		return []float64{}, err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := at.address + EmbeddingsAssistantURL
	log.Printf("Sending file to extract tokens")
	respData, err := sender.SendRequest(reqBody, &targetURL, "application/json")
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return []float64{}, err
	}

	tokensDense := &[][]float64{}
	_ = json.Unmarshal(respData, tokensDense)
	tmp := *tokensDense
	return tmp[0], nil
}

func (at *AssistantTokenizer) splitContent(content string, chunkSize int) []string {
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
