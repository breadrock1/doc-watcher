package watcher

import (
	"doc-notifier/internal/pkg/reader"
	"log"
	"sync"
	"time"
)

func (nw *NotifyWatcher) processTriggeredDocument(documents []*reader.Document) {
	wg := sync.WaitGroup{}
	for _, document := range documents {
		wg.Add(1)
		document := document
		go func() {
			defer wg.Done()
			nw.recognizedAndLoad(document)
			<-time.After(1 * time.Second)
		}()
	}

	wg.Wait()
}

func (nw *NotifyWatcher) recognizedAndLoad(document *reader.Document) {
	document.SetQuality(0)
	if err := nw.Ocr.Ocr.RecognizeFile(document); err != nil {
		log.Println(err)
		return
	}

	if nw.Tokenizer.TokenizerOptions.ChunkedFlag {
		if err := nw.loadChunkedDocument(document); err != nil {
			log.Println(err)
		}
		return
	}

	if err := nw.loadFullDocument(document); err != nil {
		log.Println(err)
		return
	}
}

func (nw *NotifyWatcher) loadFullDocument(document *reader.Document) error {
	document.ComputeMd5Hash()
	document.ComputeSsdeepHash()
	document.ComputeContentUUID()
	document.SetContentVector([]float64{})
	document.SetContentMd5Hash(document.DocumentMD5)

	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := nw.Tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for _, chunkData := range tokenVectors.Vectors {
		document.AppendContentVector(chunkData)
	}

	log.Println("Storing document to searcher: ", document.DocumentName)
	if err := nw.Searcher.StoreDocument(document); err != nil {
		log.Println("Failed while storing document: ", err)
		return err
	}

	return nil
}

func (nw *NotifyWatcher) loadChunkedDocument(document *reader.Document) error {
	document.ComputeMd5Hash()
	document.ComputeSsdeepHash()

	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := nw.Tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for chunkIndex, chunkData := range tokenVectors.ChunkedText {
		document.SetContentData(chunkData)

		contentVector := tokenVectors.Vectors[chunkIndex]
		document.SetContentVector(contentVector)
		document.ComputeContentUUID()
		document.ComputeContentMd5Hash()

		log.Println("Storing computed chunk data: ", document.ContentMD5)
		if err := nw.Searcher.StoreDocument(document); err != nil {
			log.Println("Failed while storing document: ", err)
		}
	}

	log.Println("Storing document chunks done for: ", document.DocumentName)
	return nil
}
