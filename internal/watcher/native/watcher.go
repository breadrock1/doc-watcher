package native

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/ocr"
	"doc-notifier/internal/searcher"
	"doc-notifier/internal/summarizer"
	"doc-notifier/internal/tokenizer"
	"doc-notifier/internal/watcher"
	"github.com/fsnotify/fsnotify"
)

type NotifyWatcher struct {
	stopCh chan bool

	Address       string
	pauseWatchers bool
	directories   []string
	Watcher       *fsnotify.Watcher

	Ocr        *ocr.Service
	Searcher   *searcher.Service
	Tokenizer  *tokenizer.Service
	Summarizer *summarizer.Service
}

func New(
	config *config.WatcherConfig,
	ocrService *ocr.Service,
	searcherService *searcher.Service,
	tokenService *tokenizer.Service,
	summarizeService *summarizer.Service,
) *watcher.Service {
	notifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Stopped watching: ", err)
	}

	watcherInst := &NotifyWatcher{
		stopCh:        make(chan bool),
		Address:       config.Address,
		pauseWatchers: false,
		directories:   config.WatchedDirectories,
		Ocr:           ocrService,
		Watcher:       notifyWatcher,
		Searcher:      searcherService,
		Tokenizer:     tokenService,
		Summarizer:    summarizeService,
	}

	return &watcher.Service{Watcher: watcherInst}
}

func (nw *NotifyWatcher) GetAddress() string {
	return nw.Address
}

func (nw *NotifyWatcher) RunWatchers() {
	defer func() { _ = nw.Watcher.Close() }()
	go nw.launchProcessEventLoop()
	go func() { _ = nw.AppendDirectories(nw.directories) }()
	<-nw.stopCh
}

func (nw *NotifyWatcher) IsPausedWatchers() bool {
	return nw.pauseWatchers
}

func (nw *NotifyWatcher) PauseWatchers(flag bool) {
	nw.pauseWatchers = flag
}

func (nw *NotifyWatcher) TerminateWatchers() {
	dirs := nw.Watcher.WatchList()
	_ = nw.RemoveDirectories(dirs)
	nw.stopCh <- true
}

func (nw *NotifyWatcher) CreateDirectory(dirName string) error {
	folderPath := path.Join("./indexer", dirName)
	return os.Mkdir(folderPath, os.ModePerm)
}

func (nw *NotifyWatcher) RemoveDirectory(dirName string) error {
	folderPath := path.Join("./indexer", dirName)
	return os.RemoveAll(folderPath)
}

func (nw *NotifyWatcher) GetWatchedDirectories() []string {
	return nw.Watcher.WatchList()
}

func (nw *NotifyWatcher) UploadFile(bucket string, fileName string, fileData bytes.Buffer) error {
	filePath := fmt.Sprintf("./indexer/watcher/%s/%s", bucket, fileName)
	return os.WriteFile(filePath, fileData.Bytes(), os.ModePerm)
}

func (nw *NotifyWatcher) DownloadFile(bucket string, objName string) (bytes.Buffer, error) {
	var fileBuffer bytes.Buffer
	filePath := path.Join(bucket, objName)
	fileHandler, err := os.Open(filePath)
	if err != nil {
		return fileBuffer, err
	}

	_, err = fileHandler.Read(fileBuffer.Bytes())
	if err != nil {
		return fileBuffer, err
	}

	return fileBuffer, nil
}

func (nw *NotifyWatcher) AppendDirectories(directories []string) error {
	return consumeWatcherDirectories(directories, nw.Watcher.Add)
}

func (nw *NotifyWatcher) RemoveDirectories(directories []string) error {
	return consumeWatcherDirectories(directories, nw.Watcher.Remove)
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
		case err, ok := <-nw.Watcher.Errors:
			if !ok {
				return
			}
			log.Println("Caught error: ", err)

		case event, ok := <-nw.Watcher.Events:
			if !ok {
				return
			}

			if nw.pauseWatchers {
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

	triggeredFiles := watcher.ParseCaughtFiles(absFilePath)
	nw.recognizeTriggeredDoc(triggeredFiles)
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
