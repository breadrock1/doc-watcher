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
	"doc-notifier/internal/searcher"
	"doc-notifier/internal/server"
	"doc-notifier/internal/summarizer"
	"doc-notifier/internal/tokenizer"
	"doc-notifier/internal/watcher"
	"doc-notifier/internal/watcher/minio"
)

func main() {
	serviceConfig := cmd.Execute()

	if serviceConfig.Logger.EnableFileLog {
		logger.EnableFileLogTranslating()
	}

	summarizeService, err := summarizer.New(&serviceConfig.Storage)
	if err != nil {
		log.Fatalln("failed to init summarize: ", err.Error())
	}

	ocrService := ocr.New(&serviceConfig.Ocr)
	searchService := searcher.New(&serviceConfig.Searcher)
	tokenService := tokenizer.New(&serviceConfig.Tokenizer)
	watchService := minio.New(
		&serviceConfig.Minio,
		ocrService,
		searchService,
		tokenService,
		summarizeService,
	)

	ctx, cancel := context.WithCancel(context.Background())
	go awaitSystemSignals(cancel)

	officeService := office.New(&serviceConfig.Office)
	httpServer := server.New(watchService, officeService)
	go func() {
		err := httpServer.RunServer(ctx)
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()

	go watchService.Watcher.RunWatchers()

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

func shutdownServices(ctx context.Context, httpServ *server.Service, watchServ *watcher.Service) {
	watchServ.Watcher.TerminateWatchers()
	if err := httpServ.StopServer(ctx); err != nil {
		log.Println(err)
	}
}
