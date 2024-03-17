package tokenizer

type TokenizerForm struct {
	Text              string `json:"text"`
	ChunkSize         int    `json:"chunk_size"`
	ChunkOverlap      int    `json:"chunk_overlap"`
	ReturnChunkedText bool   `json:"return_chunked_text"`
}

func CreateTokenizerForm(text string, options *Options) *TokenizerForm {
	return &TokenizerForm{
		Text:              text,
		ChunkSize:         options.ChunkSize,
		ChunkOverlap:      options.ChunkOverlap,
		ReturnChunkedText: options.ChunkedFlag,
	}
}

type ComputedTokens struct {
	Chunks      int         `json:"chunks"`
	ChunkedText []string    `json:"chunked_text"`
	Vectors     [][]float64 `json:"vectors"`
}

func CreateComputedTokesForm() *ComputedTokens {
	return &ComputedTokens{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}
}
