package tokenizer

type TokenizerService struct {
	Tokenizer Tokenizer
}

type Tokenizer interface {
	TokenizeTextData(text string) (*ComputedTokens, error)
}

func New(options *Options) *TokenizerService {
	service := &TokenizerService{}

	switch options.Mode {
	case Assistant:
		service.Tokenizer = NewAssistant(options)
	case LangChain:
		service.Tokenizer = NewLangChain(options)
	case Local:
		service.Tokenizer = NewLocal()
	case None:
		service.Tokenizer = NewNone()
	default:
		service.Tokenizer = NewNone()
	}

	return service
}
