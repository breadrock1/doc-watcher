package embeddings

import "doc-watcher/internal/watcher"

type Service struct {
	Tokenizer Tokenizer
}

type Tokenizer interface {
	Tokenize(doc *watcher.Document) (*ComputeTokens, error)
}
