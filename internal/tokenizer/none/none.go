package none

import (
	"doc-notifier/internal/models"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (nt *Service) TokenizeTextData(_ string) (*models.ComputedTokens, error) {
	return &models.ComputedTokens{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}, nil
}
