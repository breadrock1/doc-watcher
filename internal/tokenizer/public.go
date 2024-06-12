package tokenizer

import (
	"doc-notifier/internal/config"
	"doc-notifier/internal/tokenizer/assistant"
	"doc-notifier/internal/tokenizer/forms"
	"doc-notifier/internal/tokenizer/langchain"
	"doc-notifier/internal/tokenizer/none"
)

type Service struct {
	Tokenizer        Tokenizer
	TokenizerOptions *config.TokenizerConfig
}

type Tokenizer interface {
	TokenizeTextData(text string) (*forms.ComputedTokens, error)
}

func New(config *config.TokenizerConfig) *Service {
	service := &Service{}

	switch config.Mode {
	case "none":
		service.Tokenizer = none.New()
	case "assistant":
		service.Tokenizer = assistant.New(config)
	case "langchain":
		service.Tokenizer = langchain.New(config)
	default:
		service.Tokenizer = none.New()
	}

	return service
}
