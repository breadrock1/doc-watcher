package watcher

import (
	"doc-notifier/internal/pkg/llm"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
	"strings"
)

type Options struct {
	DocSearchAddress  string
	OcrServiceAddress string
	LlmServiceAddress string
	WatchDirectories  []string
}

type NotifyWatcher struct {
	directories []string
	watcher     *fsnotify.Watcher
	sender      *sender.FileSender
	reader      *reader.FileReader
	llm         *llm.Tokenizer
}

func New(cmdOpts *Options) *NotifyWatcher {
	llmService := &llm.Tokenizer{}
	fileReader := reader.New()
	fileSender := sender.New(
		cmdOpts.DocSearchAddress,
		cmdOpts.OcrServiceAddress,
		cmdOpts.LlmServiceAddress,
	)

	notifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Stopped watching: ", err)
	}

	return &NotifyWatcher{
		directories: cmdOpts.WatchDirectories,
		watcher:     notifyWatcher,
		sender:      fileSender,
		reader:      fileReader,
		llm:         llmService,
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
			continue
		}
		log.Println("Added dir to watch: ", watchDir)
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
		log.Println("Caught fs-event: ", event.Op)
		triggeredFiles := nw.reader.ParseCaughtFiles(absFilePath)
		for _, document := range triggeredFiles {
			document := document
			go func() { _ = nw.processTriggeredFile(document) }()
		}
	}
}

func (nw *NotifyWatcher) processTriggeredFile(document *reader.Document) error {
	if entity, err := nw.sender.RecognizeFileData(document.DocumentPath); err == nil {
		nw.reader.SetContentData(document, entity)
		nw.reader.ComputeMd5Hash(document)
		nw.reader.ComputeSsdeepHash(document)

		tokenVectors := nw.sender.ComputeContentTokens(document)
		for chunkIndex, chunkData := range tokenVectors.ChunkedText {
			contentData := strings.Join(chunkData, " ")
			nw.reader.SetContentData(document, contentData)

			contentVector := tokenVectors.Vectors[chunkIndex]
			nw.reader.SetContentVector(document, contentVector)

			nw.reader.ComputeUuid(document)
			nw.reader.ComputeContentMd5Hash(document)
			if err = nw.sender.StoreDocument(document); err != nil {
				log.Println("Failed while storing document: ", err)
				continue
			}
		}

		// TODO: Split text to chunks myself.
		//for _, content := range nw.reader.SplitContent(entity, 1000) {
		//	nw.reader.SetContentData(document, content)
		//	contentTokens := nw.sender.ComputeContentTokens(document)
		//
		//	nw.reader.SetContentVector(document, contentTokens)
		//	nw.reader.ComputeUuid(document)
		//	nw.reader.ComputeContentMd5Hash(document)
		//	if err = nw.sender.StoreDocument(document); err != nil {
		//		log.Println("Failed while storing document: ", err)
		//		continue
		//	}
		//}
	}
	return nil
}
