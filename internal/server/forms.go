package server

import (
	"doc-notifier/internal/models"
)

// ResponseForm example
type ResponseForm struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Done"`
}

// BadRequestForm example
type BadRequestForm struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad Request message"`
}

// ServerErrorForm example
type ServerErrorForm struct {
	Status  int    `json:"status" example:"503"`
	Message string `json:"message" example:"Server Error message"`
}

func createStatusResponse(status int, msg string) *ResponseForm {
	return &ResponseForm{Status: status, Message: msg}
}

// WatcherDirectoriesForm example
type WatcherDirectoriesForm struct {
	Paths []string `json:"paths" example:"./indexer/test_folder"`
}

// FolderNameForm example
type FolderNameForm struct {
	FolderName string `json:"folder_name" example:"test_folder"`
}

// AnalyseFilesForm example
type AnalyseFilesForm struct {
	DocumentIDs []string `json:"document_ids" example:"886f7e11874040ca8b8461fb4cd1aa2c"`
}

// UnrecognizedDocuments example
type UnrecognizedDocuments struct {
	Unrecognized []*models.Document `json:"unrecognized"`
}

// MoveFilesForm example
type MoveFilesForm struct {
	TargetDirectory string   `json:"location" example:"common-folder"`
	SourceDirectory string   `json:"src_folder_id" example:"unrecognized"`
	DocumentPaths   []string `json:"document_ids" example:"./indexer/upload/test.txt"`
}

// RemoveFilesForm example
type RemoveFilesForm struct {
	DocumentPaths []string `json:"document_paths" example:"./indexer/upload/test.txt"`
}

// RemoveFilesError example
type RemoveFilesError struct {
	Code      int      `json:"code" example:"403"`
	Message   string   `json:"message" example:"File not found"`
	FilePaths []string `json:"file_paths" example:"./indexer/upload/test.txt"`
}

// DownloadFile example
type DownloadFile struct {
	Bucket   string `json:"bucket" example:"test-bucket"`
	FileName string `json:"file_name" example:"test-file.docx"`
}

// HierarchyForm example
type HierarchyForm struct {
	Bucket        string `json:"bucket" example:"test-bucket"`
	DirectoryName string `json:"directory" example:"test-folder/"`
}
