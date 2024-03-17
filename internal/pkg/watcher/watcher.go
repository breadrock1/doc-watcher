package watcher

import (
	"doc-notifier/internal/pkg/ocr"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/searcher"
	"doc-notifier/internal/pkg/tokenizer"
	"errors"
	"github.com/fsnotify/fsnotify"
	"log"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type NotifyWatcher struct {
	stopCh chan bool

	directories []string
	watcher     *fsnotify.Watcher

	ocr       *ocr.OcrService
	reader    *reader.ReaderService
	searcher  *searcher.SearcherService
	tokenizer *tokenizer.TokenizerService
}

func New(options *Options) *NotifyWatcher {
	ocrService := ocr.New(&ocr.Options{
		Mode:    ocr.GetModeFromString(options.OcrMode),
		Address: options.OcrAddress,
	})
	readerService := reader.New()
	searcherService := searcher.New(options.SearcherAddress)
	tokenizerService := tokenizer.New(&tokenizer.Options{
		Mode:         tokenizer.GetModeFromString(options.TokenizerMode),
		Address:      options.TokenizerAddress,
		ChunkSize:    options.TokenizerChunkSize,
		ChunkedFlag:  options.TokenizerChunkedFlag,
		ChunkOverlap: options.TokenizerChunkOverlap,
	})

	notifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Stopped watching: ", err)
	}

	return &NotifyWatcher{
		stopCh:      make(chan bool),
		directories: options.WatcherDirectories,
		ocr:         ocrService,
		watcher:     notifyWatcher,
		reader:      readerService,
		searcher:    searcherService,
		tokenizer:   tokenizerService,
	}
}

func (nw *NotifyWatcher) RunWatcher() {
	defer func() { _ = nw.watcher.Close() }()
	go nw.parseEventSlot()
	go func() { _ = nw.AppendDirectories(nw.directories) }()
	<-nw.stopCh
}

func (nw *NotifyWatcher) StopWatcher() {
	dirs := nw.watcher.WatchList()
	_ = nw.RemoveDirectories(dirs)
	nw.stopCh <- true
}

func (nw *NotifyWatcher) WatchedDirsList() []string {
	return nw.watcher.WatchList()
}

func (nw *NotifyWatcher) AppendDirectories(directories []string) error {
	return consumeWatcherDirectories(directories, nw.watcher.Add)
}

func (nw *NotifyWatcher) RemoveDirectories(directories []string) error {
	return consumeWatcherDirectories(directories, nw.watcher.Remove)
}

func consumeWatcherDirectories(directories []string, consumer func(name string) error) error {
	var collectedErrs []string
	for _, watchDir := range directories {
		if err := consumer(watchDir); err != nil {
			collectedErrs = append(collectedErrs, err.Error())
			continue
		}
	}

	if len(collectedErrs) > 0 {
		msg := strings.Join(collectedErrs, "\n")
		return errors.New(msg)
	}

	return nil
}

func (nw *NotifyWatcher) parseEventSlot() {
	var (
		mu      sync.Mutex
		timers  = make(map[string]*time.Timer)
		waitFor = 100 * time.Millisecond

		testFunc = func(e fsnotify.Event) {
			nw.switchEventCase(&e)

			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	for {
		select {

		case err, ok := <-nw.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Caught error: ", err)

		case event, ok := <-nw.watcher.Events:
			if !ok {
				return
			}

			if !event.Has(fsnotify.Write) && !event.Has(fsnotify.Create) {
				continue
			}

			mu.Lock()
			t, ok := timers[event.Name]
			mu.Unlock()

			if !ok {
				t = time.AfterFunc(math.MaxInt64, func() { testFunc(event) })
				t.Stop()

				mu.Lock()
				timers[event.Name] = t
				mu.Unlock()
			}

			t.Reset(waitFor)

		}
	}
}

func (nw *NotifyWatcher) switchEventCase(event *fsnotify.Event) {
	absFilePath, err := filepath.Abs(event.Name)
	if err != nil {
		log.Println("Failed while getting abs path of file: ", err)
		return
	}

	triggeredFiles := nw.reader.ParseCaughtFiles(absFilePath)
	nw.storeExtractedDocuments(triggeredFiles)
}
