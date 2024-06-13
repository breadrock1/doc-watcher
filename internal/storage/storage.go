package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"doc-notifier/internal/config"
	"doc-notifier/internal/reader"
	_ "github.com/lib/pq"
)

type Service struct {
	address    string
	LLMAddress string
	db         *sql.DB
}

func New(config *config.StorageConfig) *Service {
	address := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Address,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
		config.EnableSSL,
	)

	return &Service{
		address:    address,
		LLMAddress: config.AddressLLM,
	}
}

func (s *Service) Connect(ctx context.Context) error {
	db, err := sql.Open("postgres", s.address)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	s.db = db
	return s.db.PingContext(ctx)
}

func (s *Service) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Service) Create(ctx context.Context, document *reader.Document) (int, error) {
	query := `
		INSERT INTO documents (folder_id, folder_path, content, document_id, document_ssdeep, 
		                       document_name, document_path, document_size, document_type, 
		                       document_ext, document_perm, document_created, document_modified, class)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, DATE($12), DATE($13), $14);
	`

	args := []interface{}{
		document.FolderID,
		document.FolderPath,
		document.Content,
		document.DocumentMD5,
		document.DocumentSSDEEP,
		document.DocumentName,
		document.DocumentPath,
		document.DocumentSize,
		document.DocumentType,
		document.DocumentExtension,
		document.DocumentPermissions,
		document.DocumentCreated,
		document.DocumentModified,
		document.DocumentClass,
	}

	var id int
	err := s.db.QueryRowContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("db exec: %v", err)
	}
	return id, nil
}
