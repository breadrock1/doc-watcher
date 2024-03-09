package watcher

import (
	"doc-notifier/internal/pkg/reader"
	"log"
	"sync"
	"time"
)

func (nw *NotifyWatcher) storeExtractedDocuments(documents []*reader.Document) {
	wg := sync.WaitGroup{}
	for _, document := range documents {
		wg.Add(1)
		document := document
		go func() {
			defer wg.Done()
			_ = nw.processTriggeredFile(document)
			<-time.After(1 * time.Second)
		}()
	}

	wg.Wait()
}

func (nw *NotifyWatcher) processTriggeredFile(document *reader.Document) error {
	var recognizeErr error
	var contentData string

	if nw.readRawFileFlag {
		contentData, recognizeErr = nw.sender.ReadRawFileData(document.DocumentPath)
	} else {
		contentData, recognizeErr = nw.sender.RecognizeFileData(document.DocumentPath)
	}

	if recognizeErr == nil {
		nw.reader.SetContentData(document, contentData)
		if nw.storeChunksFlag {
			return nw.loadChunkedDocument(document)
		} else {
			return nw.loadFullDocument(document)
		}
	}
	return nil
}

func (nw *NotifyWatcher) loadFullDocument(document *reader.Document) error {
	nw.reader.ComputeMd5Hash(document)
	nw.reader.ComputeSsdeepHash(document)
	nw.reader.ComputeUuid(document)
	nw.reader.ComputeContentMd5Hash(document)
	nw.reader.SetContentVector(document, []float64{})

	tokenVectors, _ := nw.sender.ComputeContentTokens(document)
	for _, chunkData := range tokenVectors.Vectors {
		nw.reader.AppendContentVector(document, chunkData)
	}

	if err := nw.sender.StoreDocument(document); err != nil {
		log.Println("Failed while storing document: ", err)
		return err
	}

	return nil
}

func (nw *NotifyWatcher) loadChunkedDocument(document *reader.Document) error {
	nw.reader.ComputeMd5Hash(document)
	nw.reader.ComputeSsdeepHash(document)
	tokenVectors, _ := nw.sender.ComputeContentTokens(document)
	for chunkIndex, chunkData := range tokenVectors.ChunkedText {
		nw.reader.SetContentData(document, chunkData)

		contentVector := tokenVectors.Vectors[chunkIndex]
		nw.reader.SetContentVector(document, contentVector)

		nw.reader.ComputeUuid(document)
		nw.reader.ComputeContentMd5Hash(document)
		if err := nw.sender.StoreDocument(document); err != nil {
			log.Println("Failed while storing document: ", err)
			continue
		}
	}

	return nil
}