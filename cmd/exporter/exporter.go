package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"doc-notifier/internal/sender"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
)

const CloudBucketName = "common-folder"

func main() {
	_ = godotenv.Load()

	exportDirPath := loadString("EXPORTER_DIRECTORY_PATH")
	serviceAddress := loadString("DOC_NOTIFIER_ADDRESS")

	//serviceConfig := cmd.Execute()
	//ocrService := ocr.New(&serviceConfig.Ocr)
	//tokenService := tokenizer.New(&serviceConfig.Tokenizer)
	//searchService := searcher.New(&serviceConfig.Searcher)
	//allDocuments := watcher.ParseCaughtFiles(exportDirPath)
	//
	//log.Printf("Caught %d files into directory %s", len(allDocuments), exportDirPath)
	//for _, document := range watcher.ParseCaughtFiles(exportDirPath) {
	//	log.Println("indexing file: ", document.DocumentPath)
	//	if err := ocrService.Ocr.RecognizeFile(document, document.DocumentPath); err != nil {
	//		log.Println("failed to recognize file: ", document.DocumentPath)
	//		continue
	//	}
	//
	//	document.ComputeMd5Hash()
	//	document.ComputeSsdeepHash()
	//	document.SetEmbeddings([]*models.Embeddings{})
	//
	//	log.Println("computing tokens for extracted text: ", document.DocumentName)
	//	tokenVectors, _ := tokenService.Tokenizer.TokenizeTextData(document.Content)
	//	for chunkID, chunkData := range tokenVectors.Vectors {
	//		text := tokenVectors.ChunkedText[chunkID]
	//		document.AppendContentVector(text, chunkData)
	//	}
	//
	//	log.Println("storing document to searcher: ", document.DocumentName)
	//	if err := searchService.StoreDocument(document); err != nil {
	//		log.Println("failed while storing document: ", err)
	//	}
	//}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go launchProcessEventLoop(exportDirPath, serviceAddress)
	<-sigs
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

func launchProcessEventLoop(directory string, uploadAddr string) {
	var (
		mu            = &sync.RWMutex{}
		timers        = make(map[string]*time.Timer)
		waitFor       = 100 * time.Millisecond
		uploadFileURL = fmt.Sprintf("%s/storage/%s/file/upload", uploadAddr, CloudBucketName)

		processFileCallback = func(e fsnotify.Event) {
			execProcessingPipeline(&e, uploadFileURL)

			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	notifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Stopped watching: ", err.Error())
	}

	if err := notifyWatcher.Add(directory); err != nil {
		log.Fatal("failed to append directory: ", err.Error())
	}

	for {
		select {
		case err, ok := <-notifyWatcher.Errors:
			if !ok {
				return
			}
			log.Println("Caught error: ", err.Error())

		case event, ok := <-notifyWatcher.Events:
			if !ok {
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

func execProcessingPipeline(event *fsnotify.Event, uploadAddr string) {
	absFilePath, err := filepath.Abs(event.Name)
	if err != nil {
		log.Println("failed while getting abs path of file: ", err.Error())
		return
	}

	if err := sendFileToCloud(absFilePath, uploadAddr); err != nil {
		log.Println("failed to send file to cloud: ", err.Error())
		return
	}
}

func sendFileToCloud(filePath string, targetURL string) error {
	var recErr error
	var fileHandle *os.File
	if fileHandle, recErr = os.Open(filePath); recErr != nil {
		return fmt.Errorf("file %s not found: %e", filePath, recErr)
	}
	defer func() { _ = fileHandle.Close() }()

	var reqBody bytes.Buffer
	var writer *multipart.Writer
	if writer, recErr = sender.CreateFormFile(fileHandle, &reqBody, "files"); recErr != nil {
		return fmt.Errorf("failed create form file: %e", recErr.Error())
	}

	log.Printf("sending file %s to recognize", filePath)

	mimeType := writer.FormDataContentType()
	if _, recErr = sender.POST(&reqBody, targetURL, mimeType, 300*time.Second); recErr != nil {
		return fmt.Errorf("failed send request: %e", recErr.Error())
	}

	time.Sleep(5 * time.Second)
	return nil
}
