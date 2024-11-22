package embeddings

type ComputeTokens struct {
	Chunks      int
	ChunkedText []string
	Vectors     [][]float64
}

type EmbedAllForm struct {
	Inputs    string `json:"inputs"`
	Truncate  bool   `json:"truncate"`
	Normalize bool   `json:"normalize"`
}
