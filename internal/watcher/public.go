package watcher

import (
	"bytes"

	"doc-notifier/internal/models"
)

type Service struct {
	Watcher WatcherService
}

type WatcherService interface {
	RunWatchers()
	PauseWatchers(flag bool)
	IsPausedWatchers() bool
	TerminateWatchers()
	GetAddress() string
	GetWatchedDirectories() []string
	GetHierarchy(bucket, dirName string) []*models.StorageItem
	CreateDirectory(dirName string) error
	RemoveDirectory(dirName string) error
	UploadFile(bucket string, fileName string, fileData bytes.Buffer) error
	DownloadFile(bucket string, objName string) (bytes.Buffer, error)
	RemoveFile(bucket string, fileName string) error
	AppendDirectories(directories []string) error
	RemoveDirectories(directories []string) error
}
