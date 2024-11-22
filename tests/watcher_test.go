package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"doc-watcher/internal/embeddings"
	"doc-watcher/internal/embeddings/sovavec"
	"doc-watcher/internal/ocr"
	"doc-watcher/internal/ocr/sovaocr"
	"doc-watcher/internal/searcher"
	"doc-watcher/internal/watcher"
	"doc-watcher/internal/watcher/native"
	"github.com/stretchr/testify/assert"
)

func TestWatcherManager(t *testing.T) {
	const IndexerDir = "../indexer"
	const Resources = "resources"
	const TestFile = "ocr_result.json"

	timeoutDuration := time.Duration(10) * time.Second
	ocrService := sovaocr.New(&ocr.Config{
		Address:   "localhost:3451",
		EnableSSL: false,
		Timeout:   timeoutDuration,
	})

	searchService := searcher.New(&searcher.Config{
		Address:   "localhost:3451",
		EnableSSL: false,
	})

	embedService := sovavec.New(&embeddings.Config{
		Address:      "localhost:3451",
		EnableSSL:    false,
		ChunkSize:    500,
		ChunkOverlap: 1,
		ReturnChunks: false,
		ChunkBySelf:  false,
	})

	watcherConf := &watcher.Config{
		Address:            "0.0.0.0:2893",
		WatchedDirectories: []string{IndexerDir},
	}

	watcherService := native.New(watcherConf, ocrService, searchService, embedService)

	t.Run("Append directory to watch", func(t *testing.T) {
		dirPath := path.Join(Resources)
		err := watcherService.Watcher.AppendDirectories([]string{dirPath})
		assert.NoError(t, err, "failed while appending dir to watch")

		dirs := watcherService.Watcher.GetWatchedDirs()
		assert.Equal(t, len(dirs), 1, "not equal appended dirs")
		assert.Empty(t, err, "non null error to get buckets")

		err = watcherService.Watcher.RemoveDirectories([]string{dirPath})
		assert.NoError(t, err, "failed while detach dir to watch")
	})

	t.Run("Append multiple dirs to watch", func(t *testing.T) {
		dirs := []string{Resources, IndexerDir}
		err := watcherService.Watcher.AppendDirectories(dirs)
		assert.NoError(t, err, "failed while appending dir to watch")

		attached := watcherService.Watcher.GetWatchedDirs()
		assert.Equal(t, len(dirs), len(attached), "not equal appended dirs")
		assert.Empty(t, err, "not null error to get buckets")

		err = watcherService.Watcher.RemoveDirectories(dirs)
		assert.NoError(t, err, "failed while detach dir to watch")
	})

	t.Run("Caught error while attach and detach non existing firs", func(t *testing.T) {
		filePath := path.Join(Resources, "any")
		err := watcherService.Watcher.AppendDirectories([]string{filePath})
		assert.Error(t, err, "Failed while catching error to append")

		err = watcherService.Watcher.RemoveDirectories([]string{filePath})
		assert.Error(t, err, "Failed while catching error to detach")
	})

	t.Run("Parse complex structure with OcrMetadata", func(t *testing.T) {
		filePath := path.Join(Resources, TestFile)
		file, _ := os.ReadFile(filePath)
		var previews []watcher.Document
		if err := json.Unmarshal(file, &previews); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%v", previews)
	})
}
