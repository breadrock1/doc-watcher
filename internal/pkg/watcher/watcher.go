package watcher

import (
	"doc-notifier/internal/pkg/ocr"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/searcher"
	"doc-notifier/internal/pkg/tokenizer"
	"doc-notifier/internal/pkg/tokenizer/tokoptions"
	"errors"
	"github.com/fsnotify/fsnotify"
	"log"
	"math"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type NotifyWatcher struct {
	mu                  *sync.RWMutex
	stopCh              chan bool
	AppendCh            chan *reader.Document
	ReturnCh            chan []*reader.Document
	RecognizedDocuments map[string]*reader.Document

	PauseWatchers bool
	directories   []string
	Watcher       *fsnotify.Watcher

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

	tokenizerService := tokenizer.New(&tokoptions.Options{
		Mode:         tokoptions.GetModeFromString(options.TokenizerServiceMode),
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
		mu:                  &sync.RWMutex{},
		stopCh:              make(chan bool),
		AppendCh:            make(chan *reader.Document),
		ReturnCh:            make(chan []*reader.Document),
		RecognizedDocuments: make(map[string]*reader.Document),
		PauseWatchers:       false,
		directories:         options.WatchedDirectories,
		Ocr:                 ocrService,
		Watcher:             notifyWatcher,
		Reader:              readerService,
		Searcher:            searcherService,
		Tokenizer:           tokenizerService,
	}
}

func (nw *NotifyWatcher) RunWatchers() {
	defer func() { _ = nw.Watcher.Close() }()
	go nw.launchProcessEventLoop()
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

func (nw *NotifyWatcher) IsRecognizedDocument(documentID string) bool {
	nw.mu.RLock()
	_, ok := nw.RecognizedDocuments[documentID]
	nw.mu.RUnlock()
	return ok
}

func (nw *NotifyWatcher) PopRecognizedDocument(documentID string) *reader.Document {
	var document *reader.Document
	nw.mu.Lock()
	document, _ = nw.RecognizedDocuments[documentID]
	delete(nw.RecognizedDocuments, documentID)
	nw.mu.Unlock()
	return document
}

func (nw *NotifyWatcher) AppendRecognizedDocument(document *reader.Document) {
	nw.mu.Lock()
	nw.RecognizedDocuments[document.DocumentMD5] = document
	nw.mu.Unlock()
}

func (nw *NotifyWatcher) launchProcessEventLoop() {
	var (
		mu      = &sync.RWMutex{}
		timers  = make(map[string]*time.Timer)
		waitFor = 100 * time.Millisecond

		processFileCallback = func(e fsnotify.Event) {
			nw.execProcessingPipeline(&e)

			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	for {
		select {
		case newDocument := <-nw.AppendCh:
			nw.execDocumentProcessing(newDocument)

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

			mu.RLock()
			t, ok := timers[event.Name]
			mu.RUnlock()

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

func (nw *NotifyWatcher) execProcessingPipeline(event *fsnotify.Event) {
	absFilePath, err := filepath.Abs(event.Name)
	if err != nil {
		log.Println("Failed while getting abs path of file: ", err)
		return
	}

	triggeredFiles := nw.Reader.ParseCaughtFiles(absFilePath)
	nw.processTriggeredDocument(triggeredFiles)
}

func (nw *NotifyWatcher) execDocumentProcessing(document *reader.Document) {
	if recognizeErr := nw.Ocr.Ocr.RecognizeFile(document); recognizeErr != nil {
		document.SetQuality(1)
		log.Println(recognizeErr)
		return
	}

	srcDocPath := document.DocumentPath
	targetDirPath := document.OcrMetadata.DocType
	folderPath := path.Join("./indexer/", targetDirPath)
	_ = nw.Reader.MoveFileToDir(srcDocPath, folderPath)
	dstDocPath := path.Join(folderPath, document.DocumentName)

	document.SetFolderPath(folderPath)
	document.SetDocumentPath(dstDocPath)
	document.SetContentVector([]float64{})
	document.SetQuality(reader.MaxQualityValue)
	document.SetContentMd5Hash(document.DocumentMD5)

	nw.AppendRecognizedDocument(document)
	_ = nw.Searcher.StoreDocument(document)
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
