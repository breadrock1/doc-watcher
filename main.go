package main

import (
	"doc-notifier/cmd"
	"doc-notifier/internal/pkg/options"
	"doc-notifier/internal/pkg/server"
	"doc-notifier/internal/pkg/watcher"
	"os"
	"os/signal"
	"syscall"
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
	defer watcherService.StopWatcher()

	serverOptions := options.ParseServerAddress(serviceOptions.ServerAddress)
	httpServer := server.New(serverOptions, watcherService)
	go httpServer.RunServer()
	defer httpServer.StopServer()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGKILL, syscall.SIGABRT)
	<-killSignal
}
