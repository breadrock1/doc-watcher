package watcher

import (
	"context"
	"log"
	"sync"
	"time"

	"doc-notifier/internal/reader"
)

func (nw *NotifyWatcher) recognizeTriggeredDoc(documents []*reader.Document) {
	wg := sync.WaitGroup{}
	for _, document := range documents {
		wg.Add(1)
		document := document
		go func() {
			defer wg.Done()
			nw.recognizeCallback(document)
			<-time.After(2 * time.Second)
		}()
	}

	wg.Wait()
}

func (nw *NotifyWatcher) recognizeCallback(document *reader.Document) {
	document.SetQuality(0)
	if err := nw.Ocr.Ocr.RecognizeFile(document); err != nil {
		log.Println(err)
		return
	}

	document.ComputeMd5Hash()
	document.ComputeSsdeepHash()
	document.SetEmbeddings([]*reader.Embeddings{})

	log.Println("Computing tokens for extracted text: ", document.DocumentName)
	tokenVectors, _ := nw.Tokenizer.Tokenizer.TokenizeTextData(document.Content)
	for chunkID, chunkData := range tokenVectors.Vectors {
		text := tokenVectors.ChunkedText[chunkID]
		document.AppendContentVector(text, chunkData)
	}

	log.Println("Storing document to searcher: ", document.DocumentName)
	if err := nw.Searcher.StoreDocument(document); err != nil {
		log.Println("Failed while storing document: ", err)
	}

	ctx := context.Background()
	nw.Storage.LoadSummary(document)
	if _, err := nw.Storage.Create(ctx, document); err != nil {
		log.Println("Failed while storing metadata to psql: ", err)
	}
}
