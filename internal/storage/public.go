package storage

import (
	"context"

	"doc-notifier/internal/config"
	"doc-notifier/internal/models"
	"doc-notifier/internal/storage/postgress"
)

type Storage struct {
	ss ServiceStorage
}

type ServiceStorage interface {
	Close(ctx context.Context) error
	Connect(ctx context.Context) error
	Create(ctx context.Context, document *models.Document) (int, error)
}

func New(config *config.StorageConfig) ServiceStorage {
	storeService := postgress.New(config)
	return &storeService
}
