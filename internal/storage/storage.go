package storage

import (
	"context"
	"database/sql"
	"doc-notifier/internal/config"
	"doc-notifier/internal/reader"
	"fmt"
	"log"
)

type Service struct {
	address string
	db      *sql.DB
}

func New(config *config.StorageConfig) *Service {
	address := fmt.Sprintf("%s://%s:%s@%s:%d/%s", config.DriverName, config.User, config.Password, config.Address, config.Port, config.DbName)
	return &Service{
		address: address,
	}
}

func (s *Service) Connect(ctx context.Context, connect string) error {
	db, err := sql.Open("pgx", connect)
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

//  FolderID            string        `json:"folder_id"`
//	FolderPath          string        `json:"folder_path"`
//	Content             string        `json:"content"`
//	DocumentMD5         string        `json:"document_md5"`
//	DocumentSSDEEP      string        `json:"document_ssdeep"`
//	DocumentName        string        `json:"document_name"`
//	DocumentPath        string        `json:"document_path"`
//	DocumentSize        int64         `json:"document_size"`
//	DocumentType        string        `json:"document_type"`
//	DocumentExtension   string        `json:"document_extension"`
//	DocumentPermissions int32         `json:"document_permissions"`
//	DocumentCreated     string        `json:"document_created"`
//	DocumentModified    string        `json:"document_modified"`
//	QualityRecognized   int32         `json:"quality_recognition"`
//	OcrMetadata         *OcrMetadata  `json:"ocr_metadata"`
//	Embeddings          []*Embeddings `json:"embeddings"`

func (s *Service) Create(ctx context.Context, document *reader.Document) (int, error) {
	query := `
		INSERT INTO event (folder_id, folder_path, content, document_id, document_ssdeep, 
		                   document_name, document_path, document_size, document_type,
		                   document_ext, document_perm, document_created, document_modified, class)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING document_id;
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
		document.DocumentModified,
		document.DocumentClass,
	}

	var id int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(id)
	if err != nil {
		return 0, fmt.Errorf("db exec: %w", err)
	}
	return id, nil
}
