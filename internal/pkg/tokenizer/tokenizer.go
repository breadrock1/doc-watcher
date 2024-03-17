package tokenizer

import (
	"doc-notifier/internal/pkg/tokenizer/assistant"
	"doc-notifier/internal/pkg/tokenizer/langchain"
	"doc-notifier/internal/pkg/tokenizer/local"
	"doc-notifier/internal/pkg/tokenizer/none"
)

type TokenizerService struct {
	ocr interface{}
}

type Tokenizer interface {
	TokenizeTextData(text string) ([][]float64, error)
}

func New(options *Options) *TokenizerService {
	service := &TokenizerService{}

	switch options.Mode {
	case Assistant:
		service.ocr = assistant.New(options)
	case LangChain:
		service.ocr = langchain.New(options)
	case Local:
		service.ocr = local.New(options)
	case None:
		service.ocr = none.New(options)
	default:
		service.ocr = none.New(options)
	}

	return service
}
