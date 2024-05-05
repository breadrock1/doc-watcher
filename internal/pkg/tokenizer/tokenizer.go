package tokenizer

import (
	"doc-notifier/internal/pkg/tokenizer/assistant"
	"doc-notifier/internal/pkg/tokenizer/forms"
	"doc-notifier/internal/pkg/tokenizer/langchain"
	"doc-notifier/internal/pkg/tokenizer/none"
	"doc-notifier/internal/pkg/tokenizer/tokoptions"
)

type Service struct {
	Tokenizer        Tokenizer
	TokenizerOptions *tokoptions.Options
}

type Tokenizer interface {
	TokenizeTextData(text string) (*forms.ComputedTokens, error)
}

func New(options *tokoptions.Options) *Service {
	service := &Service{
		TokenizerOptions: options,
	}

	switch options.Mode {
	case tokoptions.None:
		service.Tokenizer = none.New()
	case tokoptions.Assistant:
		service.Tokenizer = assistant.New(options)
	case tokoptions.LangChain:
		service.Tokenizer = langchain.New(options)
	default:
		service.Tokenizer = none.New()
	}

	return service
}
