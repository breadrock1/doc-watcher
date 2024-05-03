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
	stopCh        chan bool
	AppendCh      chan string
	ReturnCh      chan []*reader.Document
	PauseWatchers bool

	directories []string
	Watcher     *fsnotify.Watcher

	Ocr       *ocr.Service
	Reader    *reader.Service
	Searcher  *searcher.Service
	Tokenizer *tokenizer.Service
}

func New(options *Options) *NotifyWatcher {
	readerService := reader.New()
	timeoutDuration := time.Duration(options.TokenizerTimeout) * time.Second
	searcherService := searcher.New(options.DocSearchAddress, timeoutDuration)

	ocrService := ocr.New(&ocr.Options{
		Mode:    ocr.GetModeFromString(options.OcrServiceMode),
		Address: options.OcrServiceAddress,
		Timeout: timeoutDuration,
	})

	tokenizerService := tokenizer.New(&tokenizer.Options{
		Mode:         tokenizer.GetModeFromString(options.TokenizerServiceMode),
		Address:      options.TokenizerServiceAddress,
		Timeout:      timeoutDuration,
		ChunkSize:    options.TokenizerChunkSize,
		ChunkedFlag:  options.TokenizerReturnChunks,
		ChunkOverlap: options.TokenizerChunkOverlap,
	})

	notifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Stopped watching: ", err)
	}

	return &NotifyWatcher{
		stopCh:        make(chan bool),
		AppendCh:      make(chan string),
		PauseWatchers: false,
		directories:   options.WatchedDirectories,
		Ocr:           ocrService,
		Watcher:       notifyWatcher,
		Reader:        readerService,
		Searcher:      searcherService,
		Tokenizer:     tokenizerService,
	}
}

func (nw *NotifyWatcher) RunWatchers() {
	defer func() { _ = nw.Watcher.Close() }()
	go nw.parseEventSlot()
	go func() { _ = nw.AppendDirectories(nw.directories) }()
	<-nw.stopCh
}

func (nw *NotifyWatcher) TerminateWatchers() {
	dirs := nw.Watcher.WatchList()
	_ = nw.RemoveDirectories(dirs)
	nw.stopCh <- true
}

func (nw *NotifyWatcher) GetWatchedDirectories() []string {
	return nw.Watcher.WatchList()
}

func (nw *NotifyWatcher) AppendDirectories(directories []string) error {
	return consumeWatcherDirectories(directories, nw.Watcher.Add)
}

func (nw *NotifyWatcher) RemoveDirectories(directories []string) error {
	return consumeWatcherDirectories(directories, nw.Watcher.Remove)
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

		processFileCallback = func(e fsnotify.Event) {
			nw.runProcessingPipeline(&e)

			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	for {
		select {
		case docID := <-nw.AppendCh:
			document := nw.Reader.GetAwaitDocument(docID)
			if document == nil {
				nw.ReturnCh <- nil
			}
			absFilePath, _ := filepath.Abs(document.DocumentPath)
			parsedDoc := nw.Reader.ParseCaughtFiles(absFilePath)
			nw.ReturnCh <- parsedDoc

		case err, ok := <-nw.Watcher.Errors:
			if !ok {
				return
			}
			log.Println("Caught error: ", err)

		case event, ok := <-nw.Watcher.Events:
			if !ok {
				return
			}

			if nw.PauseWatchers {
				return
			}

			if !event.Has(fsnotify.Write) && !event.Has(fsnotify.Create) {
				continue
			}

			mu.Lock()
			t, ok := timers[event.Name]
			mu.Unlock()

			if !ok {
				t = time.AfterFunc(math.MaxInt64, func() { processFileCallback(event) })
				t.Stop()

				mu.Lock()
				timers[event.Name] = t
				mu.Unlock()
			}

			t.Reset(waitFor)
		}
	}
}

func (nw *NotifyWatcher) runProcessingPipeline(event *fsnotify.Event) {
	absFilePath, err := filepath.Abs(event.Name)
	if err != nil {
		log.Println("Failed while getting abs path of file: ", err)
		return
	}

	triggeredFiles := nw.Reader.ParseCaughtFiles(absFilePath)
	nw.storeExtractedDocuments(triggeredFiles)
}
