package watcher

import (
	"doc-notifier/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/ocr"
	"doc-notifier/internal/searcher"
	"doc-notifier/internal/summarizer"
	"doc-notifier/internal/tokenizer"
	"doc-notifier/internal/watcher"
	"github.com/stretchr/testify/assert"
)

const TestcaseDirPath = "../testcases/"
const IndexerDirPath = "../../indexer/"

func TestWatcherManager(t *testing.T) {
	timeoutDuration := time.Duration(10) * time.Second
	ocrService := ocr.New(&config.OcrConfig{
		Mode:    "raw",
		Address: "http://localhost:3451",
		Timeout: timeoutDuration,
	})
	searcherService := searcher.New(&config.SearcherConfig{
		Address: "http://localhost:3451",
		Timeout: timeoutDuration,
	})
	tokenizerService := tokenizer.New(&config.TokenizerConfig{
		Address:      "http://localhost:3451",
		Mode:         "assistant",
		ChunkSize:    500,
		ChunkOverlap: 1,
		ReturnChunks: false,
		ChunkBySelf:  false,
		Timeout:      timeoutDuration,
	})
	storeService, _ := summarizer.New(&config.StorageConfig{
		DriverName: "postgres",
		User:       "postgres",
		Password:   "postgres",
		Address:    "localhost",
		Port:       5432,
		DbName:     "postgres",
		EnableSSL:  "disable",
		AddressLLM: "http://localhost:8081",
	})

	watcherConf := &config.WatcherConfig{
		Address:            "0.0.0.0:2893",
		WatchedDirectories: []string{IndexerDirPath},
	}

	watch := watcher.New(watcherConf, ocrService, searcherService, tokenizerService, storeService)

	t.Run("Append directory to watch", func(t *testing.T) {
		err := watch.AppendDirectories([]string{TestcaseDirPath})
		assert.NoError(t, err, "Failed while appending dir to watch")

		dirs := watch.GetWatchedDirectories()
		assert.Equal(t, len(dirs), 1, "Not equal appended dirs")

		err = watch.RemoveDirectories([]string{TestcaseDirPath})
		assert.NoError(t, err, "Failed while detach dir to watch")
	})

	t.Run("Append multiple dirs to watch", func(t *testing.T) {
		dirs := []string{TestcaseDirPath, IndexerDirPath}
		err := watch.AppendDirectories(dirs)
		assert.NoError(t, err, "Failed while appending dir to watch")

		attached := watch.GetWatchedDirectories()
		assert.Equal(t, len(dirs), len(attached), "Not equal appended dirs")

		err = watch.RemoveDirectories(dirs)
		assert.NoError(t, err, "Failed while detach dir to watch")
	})

	t.Run("Caught error while attach and detach non existing firs", func(t *testing.T) {
		err := watch.AppendDirectories([]string{TestcaseDirPath + "any"})
		assert.Error(t, err, "Failed while catching error to append")

		err = watch.RemoveDirectories([]string{TestcaseDirPath + "any"})
		assert.Error(t, err, "Failed while catching error to detach")
	})

	t.Run("Parse complex structure with OcrMetadata", func(t *testing.T) {
		file, _ := os.ReadFile(TestcaseDirPath + "ocr_result.json")
		var previews []models.Document
		if err := json.Unmarshal(file, &previews); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%v", previews)
	})
}
