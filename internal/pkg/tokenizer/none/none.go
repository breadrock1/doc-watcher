package none

import "doc-notifier/internal/pkg/tokenizer"

type NoneTokenizer struct {
}

func New(_ *tokenizer.Options) *NoneTokenizer {
	return &NoneTokenizer{}
}

func (nt *NoneTokenizer) TokenizeTextData(text string) (string, error) {
	return "", nil
}
