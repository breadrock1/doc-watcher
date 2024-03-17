package watcher

type Options struct {
	OcrAddress string
	OcrMode    string

	SearcherAddress string

	TokenizerMode         string
	TokenizerAddress      string
	TokenizerChunkedFlag  bool
	TokenizerChunkSize    int
	TokenizerChunkOverlap int

	WatcherDirectories []string
}
