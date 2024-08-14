package watcher

type Service struct {
	Watcher WatcherService
}

type WatcherService interface {
	GetAddress() string
	RunWatchers()
	TerminateWatchers()
}
