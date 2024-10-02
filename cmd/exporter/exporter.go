package main

import (
	"fmt"
	"log"
	"os"

	"doc-notifier/cmd"
	"doc-notifier/internal/models"
	"doc-notifier/internal/ocr"
	"doc-notifier/internal/searcher"
	"doc-notifier/internal/tokenizer"
	"doc-notifier/internal/watcher"
	"github.com/joho/godotenv"
)

func main() {
	serviceConfig := cmd.Execute()
	ocrService := ocr.New(&serviceConfig.Ocr)
	tokenService := tokenizer.New(&serviceConfig.Tokenizer)
	searchService := searcher.New(&serviceConfig.Searcher)

	_ = godotenv.Load()
	exportDirPath := loadString("EXPORTER_DIRECTORY_PATH")

	allDocuments := watcher.ParseCaughtFiles(exportDirPath)
	log.Printf("Caught %d files into directory %s", len(allDocuments), exportDirPath)
	for _, document := range watcher.ParseCaughtFiles(exportDirPath) {
		log.Println("indexing file: ", document.DocumentPath)
		if err := ocrService.Ocr.RecognizeFile(document, document.DocumentPath); err != nil {
			log.Println("failed to recognize file: ", document.DocumentPath)
			continue
		}

		document.ComputeMd5Hash()
		document.ComputeSsdeepHash()
		document.SetEmbeddings([]*models.Embeddings{})

		log.Println("computing tokens for extracted text: ", document.DocumentName)
		tokenVectors, _ := tokenService.Tokenizer.TokenizeTextData(document.Content)
		for chunkID, chunkData := range tokenVectors.Vectors {
			text := tokenVectors.ChunkedText[chunkID]
			document.AppendContentVector(text, chunkData)
		}

		log.Println("storing document to searcher: ", document.DocumentName)
		if err := searchService.StoreDocument(document); err != nil {
			log.Println("failed while storing document: ", err)
		}
	}
}

func loadString(envName string) string {
	value, exists := os.LookupEnv(envName)
	if !exists {
		msg := fmt.Sprintf("failed to extract %s env var: %s", envName, value)
		log.Println(msg)
		return ""
	}
	return value
}
