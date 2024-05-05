package none

import (
	"doc-notifier/internal/pkg/tokenizer/forms"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (nt *Service) TokenizeTextData(_ string) (*forms.ComputedTokens, error) {
	return &forms.ComputedTokens{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}, nil
}
