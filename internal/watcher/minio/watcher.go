package minio

import (
	"context"
	"errors"
	"log"
	"slices"
	"sync"

	"doc-watcher/internal/embeddings"
	"doc-watcher/internal/ocr"
	"doc-watcher/internal/searcher"
	"doc-watcher/internal/watcher"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/patrickmn/go-cache"
)

var (
	prefix       = ""
	suffix       = ""
	eventsFilter = []string{
		"s3:ObjectCreated:*",
		"s3:ObjectRemoved:*",
	}
)

type S3Minio struct {
	stopCh      chan bool
	config      *watcher.Config
	bindBuckets *sync.Map

	mc *minio.Client

	cacher     *cache.Cache
	ocrServ    *ocr.Service
	searchServ *searcher.Service
	tokenServ  *embeddings.Service
}

func New(
	config *watcher.Config,
	ocrServ *ocr.Service,
	searchServ *searcher.Service,
	tokenServ *embeddings.Service,
) *watcher.Service {
	minioCreds := credentials.NewStaticV4(config.Username, config.Password, "")
	minioOpts := &minio.Options{
		Creds:  minioCreds,
		Secure: config.EnableSSL,
	}

	mc, err := minio.New(config.Address, minioOpts)
	if err != nil {
		log.Fatalln("failed to connect to minio cloud: ", err)
	}

	cacheServ := cache.New(config.CacheExpire, config.CacheCleanInterval)
	bindBuckets := &sync.Map{}

	watcherInst := &S3Minio{
		stopCh:      make(chan bool),
		config:      config,
		bindBuckets: bindBuckets,

		mc: mc,

		cacher:     cacheServ,
		ocrServ:    ocrServ,
		searchServ: searchServ,
		tokenServ:  tokenServ,
	}

	return &watcher.Service{Watcher: watcherInst}
}

func (mw *S3Minio) RunWatchers(_ context.Context) {
	go mw.launchProcessEventLoop()
	<-mw.stopCh
}

func (mw *S3Minio) TerminateWatchers(_ context.Context) {
	mw.stopCh <- true
}

func (mw *S3Minio) GetWatchedDirs(ctx context.Context) ([]string, error) {
	buckets, err := mw.mc.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	bucketNames := make([]string, len(buckets))
	for index, bucketInfo := range buckets {
		bucketNames[index] = bucketInfo.Name
	}

	return bucketNames, nil
}

func (mw *S3Minio) AttachDirectory(ctx context.Context, dir string) error {
	_, ok := mw.bindBuckets.Load(dir)
	if ok {
		return errors.New("directory already attached")
	}

	watcherDirs, err := mw.GetWatchedDirs(ctx)
	if err != nil {
		return err
	}

	if !slices.Contains(watcherDirs, dir) {
		return errors.New("there is no such bucket to launch")
	}

	go func() {
		defer func() {
			mw.bindBuckets.Delete(dir)
		}()

		cCtx, cancel := context.WithCancel(context.Background())
		mw.bindBuckets.Store(dir, cancel)

		for event := range mw.mc.ListenBucketNotification(cCtx, dir, prefix, suffix, eventsFilter) {
			if event.Err == nil {
				mw.extractAndStoreDocument(cCtx, event)
			}
		}
	}()

	return nil
}

func (mw *S3Minio) DetachDirectory(_ context.Context, dir string) error {
	ch, ok := mw.bindBuckets.Load(dir)
	if !ok {
		return errors.New("there is no such bucket to detach")
	}

	cancel := ch.(context.CancelFunc)
	cancel()

	return nil
}

func (mw *S3Minio) FetchProcessingDocuments(_ context.Context, files []string) *watcher.ProcessingDocuments {
	procDocs := &watcher.ProcessingDocuments{}

	for _, file := range files {
		obj, ok := mw.cacher.Get(file)
		if !ok {
			continue
		}

		document := obj.(*watcher.Document)
		switch document.QualityRecognized {
		case -1:
			procDocs.Processing = append(procDocs.Processing, file)
		case 0:
			procDocs.Unrecognized = append(procDocs.Unrecognized, file)
		default:
			procDocs.Done = append(procDocs.Done, file)
		}
	}

	return procDocs
}

func (mw *S3Minio) CleanProcessingDocuments(_ context.Context, files []string) error {
	// TODO: Add RWLock to escape data race!
	for _, file := range files {
		mw.cacher.Delete(file)
	}

	return nil
}

func (mw *S3Minio) launchProcessEventLoop() {
	wg := &sync.WaitGroup{}
	for _, bucketName := range mw.config.WatchedDirectories {
		wg.Add(1)
		go func() {
			defer func() {
				mw.bindBuckets.Delete(bucketName)
				wg.Done()
			}()

			cCtx, cancel := context.WithCancel(context.Background())
			mw.bindBuckets.Store(bucketName, cancel)

			bind := mw.mc.ListenBucketNotification(cCtx, bucketName, prefix, suffix, eventsFilter)
			for event := range bind {
				if event.Err == nil {
					mw.extractAndStoreDocument(cCtx, event)
				}
			}
		}()
	}

	wg.Wait()
}
