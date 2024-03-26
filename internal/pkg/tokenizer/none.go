package tokenizer

type NoneTokenizer struct {
}

func NewNone() *NoneTokenizer {
	return &NoneTokenizer{}
}

func (nt *NoneTokenizer) TokenizeTextData(text string) (*ComputedTokens, error) {
	return &ComputedTokens{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}, nil
}
