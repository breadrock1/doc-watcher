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

	watcherService := watcher.New(&watcher.Options{
		WatcherServiceAddress: serviceOptions.WatcherServiceAddress,
		WatchedDirectories:    serviceOptions.WatchedDirectories,

		OcrServiceAddress: serviceOptions.OcrServiceAddress,
		OcrServiceMode:    serviceOptions.OcrServiceMode,

		DocSearchAddress: serviceOptions.DocSearchAddress,

		TokenizerServiceAddress: serviceOptions.TokenizerServiceAddress,
		TokenizerServiceMode:    serviceOptions.TokenizerServiceMode,
		TokenizerChunkSize:      serviceOptions.TokenizerChunkSize,
		TokenizerChunkOverlap:   serviceOptions.TokenizerChunkOverlap,
		TokenizerReturnChunks:   serviceOptions.TokenizerReturnChunks,
		TokenizerChunkBySelf:    serviceOptions.TokenizerChunkBySelf,
		TokenizerTimeout:        serviceOptions.TokenizerTimeout,
	})

	go watcherService.RunWatcher()
	defer watcherService.StopWatcher()

	serverOptions := options.ParseServerAddress(serviceOptions.WatcherServiceAddress)
	httpServer := server.New(serverOptions, watcherService)
	go httpServer.RunServer()
	defer httpServer.StopServer()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	<-killSignal
}
