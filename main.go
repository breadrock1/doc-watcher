package main

import (
	"doc-notifier/cmd"
	"doc-notifier/internal/pkg/options"
	"doc-notifier/internal/pkg/server"
	"doc-notifier/internal/pkg/watcher"
)

func main() {
	serviceOptions := cmd.Execute()

	watcherService := watcher.New(
		serviceOptions.ReadRawFileFlag,
		serviceOptions.StoreChunksFlag,
		serviceOptions.DocSearchAddress,
		serviceOptions.OcrServiceAddress,
		serviceOptions.LlmServiceAddress,
		serviceOptions.WatchDirectories,
	)
	go watcherService.RunWatcher()

	serverOptions := options.ParseServerAddress(serviceOptions.ServerAddress)
	httpServer := server.New(serverOptions, watcherService)
	go httpServer.RunServer()

	<-make(chan interface{})
}
