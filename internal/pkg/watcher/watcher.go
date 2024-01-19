package watcher

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
)

type NotifyWatcher struct {
	directories []string
	watcher     *fsnotify.Watcher
	sender      *sender.FileSender
	reader      *reader.FileReader
}

func New(cmdOpts *cmd.Options) *NotifyWatcher {
	fileReader := reader.New()
	fileSender := sender.New(cmdOpts.SearcherAddress, cmdOpts.AssistantAddress)
	notifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Stopped watching: ", err)
	}

	return &NotifyWatcher{
		directories: cmdOpts.WatchDirectories,
		watcher:     notifyWatcher,
		sender:      fileSender,
		reader:      fileReader,
	}
}

func (nw *NotifyWatcher) RunWatcher() {
	defer func() { _ = nw.watcher.Close() }()
	go nw.parseEventSlot()
	nw.appendDirectories()
	<-make(chan interface{})
}

func (nw *NotifyWatcher) parseEventSlot() {
	for {
		select {
		case event, ok := <-nw.watcher.Events:
			if !ok {
				return
			}
			nw.switchEventCase(event)
		case err, ok := <-nw.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Caught error: ", err)
		}
	}
}

func (nw *NotifyWatcher) appendDirectories() {
	for _, watchDir := range nw.directories {
		err := nw.watcher.Add(watchDir)
		if err != nil {
			msg := "Failed while append directory to watcher: "
			log.Println(msg, err)
		}
	}
}

func (nw *NotifyWatcher) switchEventCase(event fsnotify.Event) {
	absFilePath, err := filepath.Abs(event.Name)
	if err != nil {
		msg := "Failed while getting abs path of file: "
		log.Println(msg, err)
		return
	}

	if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
		log.Println("Caught event: ", event.Op)
		triggeredFiles := nw.reader.ParseCaughtFiles(absFilePath)
		for _, document := range triggeredFiles {
			document := document
			go func() { _ = nw.processTriggeredFile(document) }()
		}
	}
}

func (nw *NotifyWatcher) processTriggeredFile(document *reader.Document) error {
	if entity, err := nw.sender.RecognizeFileData(document.DocumentPath); err == nil {
		nw.reader.SetEntityData(document, entity)
		nw.reader.ComputeMd5Hash(document)
		nw.reader.ComputeSsdeepHash(document)
		if err = nw.sender.StoreDocument(document); err != nil {
			log.Println("Failed while storing document: ", err)
			return err
		}
	}
	return nil
}
