package watcher

type Options struct {
	WatcherServiceAddress string
	WatchedDirectories    []string

	OcrServiceAddress string
	OcrServiceMode    string

	DocSearchAddress string

	TokenizerServiceAddress string
	TokenizerServiceMode    string
	TokenizerChunkSize      int
	TokenizerChunkOverlap   int
	TokenizerReturnChunks   bool
	TokenizerChunkBySelf    bool
	TokenizerTimeout        uint
}
