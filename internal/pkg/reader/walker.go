package reader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func (f *ReaderService) ParseCaughtFiles(filePath string) []*Document {
	mu := &sync.Mutex{}
	var customList []*Document

	wg := &sync.WaitGroup{}
	for _, filePath := range getEntityFiles(filePath) {
		wg.Add(1)
		filePath := filePath

		go func() {
			defer wg.Done()

			if doc, err := ParseFile(filePath); err == nil {
				log.Println("Caught parsed document: ", doc.DocumentName)
				mu.Lock()
				customList = append(customList, doc)
				mu.Unlock()
				return
			}

			log.Println("Failed parsing document: ", filePath)
		}()
	}

	wg.Wait()

	return customList
}

func getEntityFiles(filePath string) []string {
	<-time.After(time.Second)

	var files []string
	err := filepath.Walk(filePath, visitEntity(&files))
	if err != nil {
		log.Println("Error while walking: ", err)
	}
	return files
}

func visitEntity(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		if !info.IsDir() {
			*files = append(*files, path)
		}

		return nil
	}
}
