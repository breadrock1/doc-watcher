package tokenizer

type TokenizerService struct {
	Tokenizer        Tokenizer
	TokenizerOptions *Options
}

type Tokenizer interface {
	TokenizeTextData(text string) (*ComputedTokens, error)
}

func New(options *Options) *TokenizerService {
	service := &TokenizerService{
		TokenizerOptions: options,
	}

	switch options.Mode {
	case Assistant:
		service.Tokenizer = NewAssistant(options)
	case LangChain:
		service.Tokenizer = NewLangChain(options)
	case None:
		service.Tokenizer = NewNone()
	default:
		service.Tokenizer = NewNone()
	}

	return service
}
