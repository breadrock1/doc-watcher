package watcher

import (
	"context"
	"doc-notifier/internal/office"
	"errors"
	"log"
	"math"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/ocr"
	"doc-notifier/internal/reader"
	"doc-notifier/internal/searcher"
	"doc-notifier/internal/storage"
	"doc-notifier/internal/tokenizer"
	"github.com/fsnotify/fsnotify"
)

type NotifyWatcher struct {
	mu                  *sync.RWMutex
	stopCh              chan bool
	AppendCh            chan *reader.Document
	ReturnCh            chan []*reader.Document
	RecognizedDocuments map[string]*reader.Document

	Address       string
	PauseWatchers bool
	directories   []string
	Watcher       *fsnotify.Watcher

	Ocr       *ocr.Service
	Reader    *reader.Service
	Searcher  *searcher.Service
	Tokenizer *tokenizer.Service
	Storage   *storage.Service
	Office    *office.Service
}

func New(
	config *config.WatcherConfig,
	readerService *reader.Service,
	ocrService *ocr.Service,
	searcherService *searcher.Service,
	tokenService *tokenizer.Service,
	storageService *storage.Service,
	officeService *office.Service,
) *NotifyWatcher {
	notifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Stopped watching: ", err)
	}

	_ = storageService.Connect(context.Background())

	return &NotifyWatcher{
		mu:                  &sync.RWMutex{},
		stopCh:              make(chan bool),
		AppendCh:            make(chan *reader.Document, 10),
		ReturnCh:            make(chan []*reader.Document),
		RecognizedDocuments: make(map[string]*reader.Document),
		Address:             config.Address,
		PauseWatchers:       false,
		directories:         config.WatchedDirectories,
		Ocr:                 ocrService,
		Watcher:             notifyWatcher,
		Reader:              readerService,
		Searcher:            searcherService,
		Tokenizer:           tokenService,
		Storage:             storageService,
		Office:              officeService,
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
	nw.RecognizedDocuments[document.DocumentID] = document
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
	nw.recognizeTriggeredDoc(triggeredFiles)
}

func (nw *NotifyWatcher) execDocumentProcessing(document *reader.Document) {
	if recognizeErr := nw.Ocr.Ocr.RecognizeFile(document); recognizeErr != nil {
		document.SetQuality(1)
		log.Println(recognizeErr)
		return
	}

	targetDirPath := strings.ToLower(document.GetDocType())
	folderID, _ := nw.Searcher.GetFolderID(targetDirPath)
	if folderID == "unrecognized" {
		document.SetQuality(0)
	}

	folderPath := path.Join("./indexer/", folderID)
	_ = nw.Reader.MoveFileToDir(document.DocumentPath, folderPath)
	dstDocPath := path.Join(folderPath, document.DocumentName)

	document.SetFolderID(folderID)
	document.SetFolderPath(folderPath)
	document.SetDocumentPath(dstDocPath)
	document.SetEmbeddings([]*reader.Embeddings{})

	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := nw.Tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for chunkID, chunkData := range tokenVectors.Vectors {
		text := tokenVectors.ChunkedText[chunkID]
		document.AppendContentVector(text, chunkData)
	}

	nw.AppendRecognizedDocument(document)
	log.Printf("Store successful document %s to %s: ", document.DocumentName, targetDirPath)
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
