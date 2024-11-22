package tests

import (
	"context"
	"doc-watcher/tests/mock"
	"path"
	"testing"
	"time"

	"doc-watcher/internal/searcher"
	"doc-watcher/internal/watcher"
	"github.com/stretchr/testify/assert"
)

func TestStoreDocument(t *testing.T) {
	const Resources = "resources"
	const TestFile = "test_file_1.txt"

	searchService := searcher.New(&searcher.Config{
		Address:   "localhost:3451",
		EnableSSL: false,
	})

	mockServer := mock.InitMockServer()
	go func() {
		_ = mockServer.Start("localhost:3451")
		time.AfterFunc(5*time.Second, func() {
			_ = mockServer.Shutdown(context.Background())
		})
	}()

	t.Run("Store Document", func(t *testing.T) {
		filePath := path.Join(Resources, "directory", TestFile)
		document, parseErr := watcher.ParseFile(filePath)
		storeErr := searchService.StoreDocument(document)

		assert.NoError(t, parseErr, "error while parsing file")
		assert.NoError(t, storeErr, "error while storing document")
	})

	t.Run("Caught error while storing file", func(t *testing.T) {
		filePath := path.Join(Resources, "directory", "any")
		_, parseErr := watcher.ParseFile(filePath)

		assert.Error(t, parseErr, "non null error pointer")
		assert.Error(t, parseErr, "not equal file data")
	})

	t.Run("Caught error with service denied", func(t *testing.T) {
		sService := searcher.New(&searcher.Config{
			Address:   "localhost:4444",
			EnableSSL: false,
		})

		es := mock.InitMockServer()
		go func() {
			_ = mockServer.Start("localhost:3451")
			time.AfterFunc(2*time.Second, func() {
				_ = es.Shutdown(context.Background())
			})
		}()

		filePath := path.Join(Resources, "directory", TestFile)
		document, parseErr := watcher.ParseFile(filePath)
		storeErr := sService.StoreDocument(document)

		assert.NoError(t, parseErr, "non null error pointer")
		assert.Error(t, storeErr, "not equal file data")
	})
}
