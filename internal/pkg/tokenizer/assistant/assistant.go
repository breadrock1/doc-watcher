package assistant

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"doc-notifier/internal/pkg/tokenizer"
	"encoding/json"
	"log"
	"strings"
)

const EmbeddingsURL = "/api/v1/vectorizer/get_vectors"

type AssistantTokenizer struct {
	address string
	options *tokenizer.Options
}

func New(options *tokenizer.Options) *AssistantTokenizer {
	return &AssistantTokenizer{
		address: options.Address,
		options: options,
	}
}

func (at *AssistantTokenizer) TokenizeTextData(content string) (*tokenizer.ComputedTokens, error) {
	computedTokens := tokenizer.CreateComputedTokesForm()

	contentData := strings.ReplaceAll(content, "\n", " ")
	textVectors := tokenizer.CreateTokenizerForm(contentData, at.options)

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		log.Println("Failed while marshaling text vectors: ", err)
		return computedTokens, err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := at.address + EmbeddingsURL
	log.Printf("Sending file to extract tokens")
	respData, err := sender.SendRequest(reqBody, &targetURL)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return computedTokens, err
	}

	_ = json.Unmarshal(respData, computedTokens)
	return computedTokens, nil
}
