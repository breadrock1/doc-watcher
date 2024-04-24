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
	contentData, recognizeErr := nw.Ocr.Ocr.RecognizeFile(document)
	if recognizeErr == nil {
		nw.Reader.SetContentData(document, contentData)
		if nw.Tokenizer.TokenizerOptions.ChunkedFlag {
			return nw.loadChunkedDocument(document)
		}
		return nw.loadFullDocument(document)
	}
	nw.Reader.MoveFileToUnrecognized(document)
	return nil
}

func (nw *NotifyWatcher) loadFullDocument(document *reader.Document) error {
	nw.Reader.ComputeMd5Hash(document)
	nw.Reader.ComputeSsdeepHash(document)
	nw.Reader.ComputeUUID(document)
	nw.Reader.ComputeContentMd5Hash(document)
	nw.Reader.SetContentVector(document, []float64{})

	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := nw.Tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for _, chunkData := range tokenVectors.Vectors {
		nw.Reader.AppendContentVector(document, chunkData)
	}

	log.Println("Storing document object: ", document.DocumentName)
	if err := nw.Searcher.StoreDocument(document); err != nil {
		log.Println("Failed while storing document: ", err)
		return err
	}

	return nil
}

func (nw *NotifyWatcher) loadChunkedDocument(document *reader.Document) error {
	nw.Reader.ComputeMd5Hash(document)
	nw.Reader.ComputeSsdeepHash(document)
	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := nw.Tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for chunkIndex, chunkData := range tokenVectors.ChunkedText {
		nw.Reader.SetContentData(document, chunkData)

		contentVector := tokenVectors.Vectors[chunkIndex]
		nw.Reader.SetContentVector(document, contentVector)

		nw.Reader.ComputeUUID(document)
		nw.Reader.ComputeContentMd5Hash(document)
		log.Println("Storing computed chunk data: ", document.ContentMD5)
		if err := nw.Searcher.StoreDocument(document); err != nil {
			log.Println("Failed while storing document: ", err)
			continue
		}
	}

	log.Println("Storing document chunks done for: ", document.DocumentName)
	return nil
}
