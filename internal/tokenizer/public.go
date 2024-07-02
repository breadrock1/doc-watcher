package tokenizer

import (
	"doc-notifier/internal/config"
	"doc-notifier/internal/models"
	"doc-notifier/internal/tokenizer/assistant"
	"doc-notifier/internal/tokenizer/none"
)

type Service struct {
	Tokenizer        Tokenizer
	TokenizerOptions *config.TokenizerConfig
}

type Tokenizer interface {
	TokenizeTextData(text string) (*models.ComputedTokens, error)
}

func New(config *config.TokenizerConfig) *Service {
	service := &Service{}

	switch config.Mode {
	case "assistant":
		service.Tokenizer = assistant.New(config)
	case "none":
		service.Tokenizer = none.New()
	default:
		service.Tokenizer = none.New()
	}

	return service
}
