package models

type StorageItem struct {
	FileName      string `json:"file_name"`
	DirectoryName string `json:"directory_name"`
	IsDirectory   bool   `json:"is_directory"`
}
