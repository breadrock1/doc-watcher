package tokenizer

type LocalTokenizer struct{}

func NewLocal() *LocalTokenizer {
	return &LocalTokenizer{}
}

func (lt *LocalTokenizer) TokenizeTextData(_ string) (*ComputedTokens, error) {
	return &ComputedTokens{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}, nil
}
