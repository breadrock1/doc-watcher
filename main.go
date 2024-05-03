package main

import (
	"doc-notifier/cmd"
	"doc-notifier/internal/pkg/options"
	"doc-notifier/internal/pkg/server"
	"doc-notifier/internal/pkg/watcher"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func enableFileLogTranslating() {
	log.SetFlags(log.Ldate | log.Ltime)
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)
}

func main() {
	serviceOptions := cmd.Execute()

	if serviceOptions.EnableFileLog {
		enableFileLogTranslating()
	}

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

	go watcherService.RunWatchers()
	defer watcherService.TerminateWatchers()

	serverOptions := options.ParseServerAddress(serviceOptions.WatcherServiceAddress)
	httpServer := server.New(serverOptions, watcherService)
	go httpServer.RunServer()
	defer httpServer.StopServer()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	<-killSignal
}
