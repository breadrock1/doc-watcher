package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"doc-watcher/cmd"
	"doc-watcher/internal/embeddings/sovavec"
	"doc-watcher/internal/ocr/sovaocr"
	"doc-watcher/internal/searcher"
	"doc-watcher/internal/server"
	"doc-watcher/internal/server/httpserv"
	"doc-watcher/internal/watcher"
	"doc-watcher/internal/watcher/minio"
)

func main() {
	servConfig := cmd.Execute()

	ocrService := sovaocr.New(&servConfig.Ocr)
	searchService := searcher.New(&servConfig.Searcher)
	embedService := sovavec.New(&servConfig.Embeddings)
	watchService := minio.New(
		&servConfig.Watcher,
		ocrService,
		searchService,
		embedService,
	)

	ctx, cancel := context.WithCancel(context.Background())
	go awaitSystemSignals(cancel)

	httpServer := httpserv.New(&servConfig.Server, watchService)
	go func() {
		if err := httpServer.Server.Start(ctx); err != nil {
			log.Printf("failed to start server: %w", err)
		}
	}()

	go watchService.Watcher.RunWatchers(ctx)

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

func shutdownServices(ctx context.Context, httpServ *server.Server, watchServ *watcher.Service) {
	watchServ.Watcher.TerminateWatchers(ctx)
	if err := httpServ.Server.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
