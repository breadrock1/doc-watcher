package local

import (
	"doc-notifier/internal/pkg/tokenizer"
)

type LocalTokenizer struct {
	options *tokenizer.Options
}

func New(options *tokenizer.Options) *LocalTokenizer {
	return &LocalTokenizer{
		options: options,
	}
}

func (lt *LocalTokenizer) TokenizeTextData(_ string) (*tokenizer.ComputedTokens, error) {
	computedTokens := tokenizer.CreateComputedTokesForm()
	return computedTokens, nil
}
