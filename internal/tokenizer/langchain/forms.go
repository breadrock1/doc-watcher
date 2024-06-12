package langchain

type GetTokensForm struct {
	Text              string `json:"text"`
	ChunkSize         int    `json:"chunk_size"`
	ChunkOverlap      int    `json:"chunk_overlap"`
	ReturnChunkedText bool   `json:"return_chunked_text"`
}
