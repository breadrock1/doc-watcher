package watcher

import (
	"bytes"
	"doc-notifier/internal/models"
)

type Service struct {
	Watcher CombinedInterface
}

type CombinedInterface interface {
	WatcherService
	StorageService
}

type WatcherService interface {
	GetAddress() string
	RunWatchers()
	PauseWatchers(flag bool)
	IsPausedWatchers() bool
	TerminateWatchers()
	AppendDirectories(directories []string) error
	RemoveDirectories(directories []string) error
}

type StorageService interface {
	GetBuckets() []string
	CopyFile(bucket, srcPath, dstPath string) error
	MoveFile(bucket, srcPath, dstPath string) error
	GetListFiles(bucket, dirName string) []*models.StorageItem
	CreateBucket(dirName string) error
	RemoveBucket(dirName string) error
	UploadFile(bucket string, fileName string, fileData bytes.Buffer) error
	DownloadFile(bucket string, objName string) (bytes.Buffer, error)
	RemoveFile(bucket string, fileName string) error
}
