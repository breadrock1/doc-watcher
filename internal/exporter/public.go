package exporter

import "doc-notifier/internal/models"

type Service struct {
	client Exporter
}

type Exporter interface {
	Close() error
	GetListDirs() []string
	WalkFiles(entryDir string) []*models.Document
}
