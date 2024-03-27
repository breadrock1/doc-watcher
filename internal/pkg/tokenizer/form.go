package tokenizer

type GetTokensForm struct {
	Text              string `json:"text"`
	ChunkSize         int    `json:"chunk_size"`
	ChunkOverlap      int    `json:"chunk_overlap"`
	ReturnChunkedText bool   `json:"return_chunked_text"`
}

type ComputedTokens struct {
	Chunks      int         `json:"chunks"`
	ChunkedText []string    `json:"chunked_text"`
	Vectors     [][]float64 `json:"vectors"`
}
