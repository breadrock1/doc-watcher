package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"doc-notifier/cmd"
	"doc-notifier/internal/logger"
	"doc-notifier/internal/ocr"
	"doc-notifier/internal/office"
	"doc-notifier/internal/reader"
	"doc-notifier/internal/searcher"
	"doc-notifier/internal/server"
	"doc-notifier/internal/storage"
	"doc-notifier/internal/tokenizer"
	"doc-notifier/internal/watcher"
)

func main() {
	serviceConfig := cmd.Execute()

	if serviceConfig.Logger.EnableFileLog {
		logger.EnableFileLogTranslating()
	}

	readService := reader.New()
	officeService := office.New(&serviceConfig.Office)
	ocrService := ocr.New(&serviceConfig.Ocr)
	searchService := searcher.New(&serviceConfig.Searcher)
	tokenService := tokenizer.New(&serviceConfig.Tokenizer)
	storeService := storage.New(&serviceConfig.Storage)
	watchService := watcher.New(
		&serviceConfig.Watcher,
		readService,
		ocrService,
		searchService,
		tokenService,
		storeService,
		officeService,
	)

	ctx, cancel := context.WithCancel(context.Background())
	go awaitSystemSignals(cancel)

	httpServer := server.New(watchService)
	go func() {
		err := httpServer.RunServer(ctx)
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()

	go watchService.RunWatchers()

	<-ctx.Done()
	cancel()
	shutdownServices(ctx, httpServer, watchService)
}

func awaitSystemSignals(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	cancel()
}

func shutdownServices(ctx context.Context, httpServ *server.Service, watchServ *watcher.NotifyWatcher) {
	watchServ.TerminateWatchers()
	if err := httpServ.StopServer(ctx); err != nil {
		log.Println(err)
	}
}
