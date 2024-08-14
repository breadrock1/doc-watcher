package minio

import (
	"context"
	"log"

	"doc-notifier/internal/config"
	"doc-notifier/internal/ocr"
	"doc-notifier/internal/searcher"
	"doc-notifier/internal/summarizer"
	"doc-notifier/internal/tokenizer"
	"doc-notifier/internal/watcher"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioWatcher struct {
	stopCh chan bool

	Address string

	mc *minio.Client

	Ocr        *ocr.Service
	Searcher   *searcher.Service
	Tokenizer  *tokenizer.Service
	Summarizer *summarizer.Service
}

func New(
	config *config.MinioConfig,
	ocrService *ocr.Service,
	searcherService *searcher.Service,
	tokenService *tokenizer.Service,
	summarizeService *summarizer.Service,
) *watcher.Service {
	client, _ := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioRootUser, config.MinioRootPassword, ""),
		Secure: config.MinioUseSSL,
	})

	watcherInst := &MinioWatcher{
		stopCh:  make(chan bool),
		Address: config.Address,

		mc: client,

		Ocr:        ocrService,
		Searcher:   searcherService,
		Tokenizer:  tokenService,
		Summarizer: summarizeService,
	}

	return &watcher.Service{Watcher: watcherInst}
}

func (mw *MinioWatcher) GetAddress() string {
	return mw.Address
}

func (mw *MinioWatcher) RunWatchers() {
	ctx := context.Background()
	buckets, err := mw.mc.ListBuckets(ctx)
	if err != nil {
		log.Println(err)
	}

	log.Println(buckets)

	go mw.launchProcessEventLoop()
	<-mw.stopCh
}

func (mw *MinioWatcher) TerminateWatchers() {

}

func (mw *MinioWatcher) launchProcessEventLoop() {
	var (
		prefix       = ""
		suffix       = ""
		eventsFilter = []string{
			"s3:ObjectCreated:*",
			"s3:ObjectAccessed:*",
			"s3:ObjectRemoved:*",
		}
	)

	ctx := context.Background()
	for event := range mw.mc.ListenNotification(ctx, prefix, suffix, eventsFilter) {
		if event.Err != nil {
			log.Fatalln(event.Err)
		}
		log.Println(event)
	}
}
