package embeddings

type Config struct {
	Address      string
	EnableSSL    bool
	ChunkSize    int
	ChunkOverlap int
	ReturnChunks bool
	ChunkBySelf  bool
}
