package watcher

import (
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"errors"
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
	"strings"
)

type NotifyWatcher struct {
	storeChunksFlag bool
	readRawFileFlag bool
	directories     []string
	watcher         *fsnotify.Watcher
	sender          *sender.FileSender
	reader          *reader.FileReader
}

func New(rawFlag, storeFlag bool, searchAddr, ocrAddr, llmAddr string, watchDirs []string) *NotifyWatcher {
	fileReader := reader.New()
	fileSender := sender.New(
		searchAddr,
		ocrAddr,
		llmAddr,
		rawFlag,
	)

	notifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Stopped watching: ", err)
	}

	return &NotifyWatcher{
		readRawFileFlag: rawFlag,
		storeChunksFlag: storeFlag,
		directories:     watchDirs,
		watcher:         notifyWatcher,
		sender:          fileSender,
		reader:          fileReader,
	}
}

func (nw *NotifyWatcher) RunWatcher() {
	defer func() { _ = nw.watcher.Close() }()
	go nw.parseEventSlot()
	go func() { _ = nw.AppendDirectories(nw.directories) }()
	<-make(chan interface{})
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
	for {
		select {
		case event, ok := <-nw.watcher.Events:
			if !ok {
				return
			}

			nw.switchEventCase(&event)

		case err, ok := <-nw.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Caught error: ", err)
		}
	}
}

func (nw *NotifyWatcher) switchEventCase(event *fsnotify.Event) {
	if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
		absFilePath, err := filepath.Abs(event.Name)
		if err != nil {
			log.Println("Failed while getting abs path of file: ", err)
			return
		}

		triggeredFiles := nw.reader.ParseCaughtFiles(absFilePath)
		nw.storeExtractedDocuments(triggeredFiles)
	}
}
