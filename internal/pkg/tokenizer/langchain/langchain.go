package langchain

import (
	"bytes"
	"doc-notifier/internal/pkg/sender"
	"doc-notifier/internal/pkg/tokenizer"
	"encoding/json"
	"log"
	"strings"
)

const EmbeddingsURL = "/api/v1/vectorizer/get_vectors"

type LangChainTokenizer struct {
	address string
	options *tokenizer.Options
}

func New(options *tokenizer.Options) *LangChainTokenizer {
	return &LangChainTokenizer{
		address: options.Address,
		options: options,
	}
}

func (lt *LangChainTokenizer) TokenizeTextData(content string) (*tokenizer.ComputedTokens, error) {
	computedTokens := tokenizer.CreateComputedTokesForm()

	contentData := strings.ReplaceAll(content, "\n", " ")
	textVectors := tokenizer.CreateTokenizerForm(contentData, lt.options)

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		log.Println("Failed while marshaling text vectors: ", err)
		return computedTokens, err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetURL := lt.address + EmbeddingsURL
	log.Printf("Sending file to extract tokens")
	respData, err := sender.SendRequest(reqBody, &targetURL)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return computedTokens, err
	}

	_ = json.Unmarshal(respData, computedTokens)
	return computedTokens, nil
}
