package watcher

import "context"

type Service struct {
	Watcher IWatcher
}

type IWatcher interface {
	IDirectories
	ILaunch
	IProcessing
}

type ILaunch interface {
	RunWatchers(ctx context.Context)
	TerminateWatchers(ctx context.Context)
}

type IDirectories interface {
	GetWatchedDirs(ctx context.Context) ([]string, error)
	AttachDirectory(ctx context.Context, dir string) error
	DetachDirectory(ctx context.Context, dir string) error
}

type IProcessing interface {
	CleanProcessingDocuments(ctx context.Context, files []string) error
	FetchProcessingDocuments(ctx context.Context, files []string) *ProcessingDocuments
}
