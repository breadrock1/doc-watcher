package ocr

import (
	"doc-watcher/internal/watcher"
)

type Service struct {
	Ocr Recognizer
}

type Recognizer interface {
	RecognizeFile(document *watcher.Document, data []byte) error
}
