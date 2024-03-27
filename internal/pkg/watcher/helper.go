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
	contentData, recognizeErr := nw.ocr.Ocr.RecognizeFile(document.DocumentPath)
	if recognizeErr == nil {
		nw.reader.SetContentData(document, contentData)
		if nw.tokenizer.TokenizerOptions.ChunkedFlag {
			return nw.loadChunkedDocument(document)
		}

		return nw.loadFullDocument(document)
	}
	return nil
}

func (nw *NotifyWatcher) loadFullDocument(document *reader.Document) error {
	nw.reader.ComputeMd5Hash(document)
	nw.reader.ComputeSsdeepHash(document)
	nw.reader.ComputeUUID(document)
	nw.reader.ComputeContentMd5Hash(document)
	nw.reader.SetContentVector(document, []float64{})

	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := nw.tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for _, chunkData := range tokenVectors.Vectors {
		nw.reader.AppendContentVector(document, chunkData)
	}

	log.Println("Storing document object: ", document.DocumentName)
	if err := nw.searcher.StoreDocument(document); err != nil {
		log.Println("Failed while storing document: ", err)
		return err
	}

	return nil
}

func (nw *NotifyWatcher) loadChunkedDocument(document *reader.Document) error {
	nw.reader.ComputeMd5Hash(document)
	nw.reader.ComputeSsdeepHash(document)
	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := nw.tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for chunkIndex, chunkData := range tokenVectors.ChunkedText {
		nw.reader.SetContentData(document, chunkData)

		contentVector := tokenVectors.Vectors[chunkIndex]
		nw.reader.SetContentVector(document, contentVector)

		nw.reader.ComputeUUID(document)
		nw.reader.ComputeContentMd5Hash(document)
		log.Println("Storing computed chunk data: ", document.ContentMD5)
		if err := nw.searcher.StoreDocument(document); err != nil {
			log.Println("Failed while storing document: ", err)
			continue
		}
	}

	log.Println("Storing document chunks done for: ", document.DocumentName)
	return nil
}
