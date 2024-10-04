package postqress

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"doc-notifier/internal/config"
	"doc-notifier/internal/models"
)

type Service struct {
	dbAddress string
	db        *sql.DB
}

func New(config *config.StorageConfig) Service {
	address := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Address,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
		config.EnableSSL,
	)

	return Service{
		dbAddress: address,
	}
}

func (s *Service) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Service) Connect(ctx context.Context) error {
	db, err := sql.Open("postgres", s.dbAddress)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	s.db = db
	return s.db.PingContext(ctx)
}

func (s *Service) Create(ctx context.Context, document *models.Document) (int, error) {
	query := `
		INSERT INTO documents (folder_id, folder_path, content, document_id, document_ssdeep, 
		                       document_name, document_path, document_size, document_type, 
		                       document_ext, document_perm, document_created, document_modified, class)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, DATE($12), DATE($13), $14);
	`

	args := []interface{}{
		document.FolderID,
		document.FolderPath,
		[]byte(document.Content),
		document.DocumentID,
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
	row := s.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return 0, fmt.Errorf("db exec: %s", row.Err())
	}
	return id, nil
}
