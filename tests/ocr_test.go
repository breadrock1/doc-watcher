package tests

import (
	"context"
	"doc-watcher/tests/mock"
	"path"
	"testing"
	"time"

	"doc-watcher/internal/ocr"
	"doc-watcher/internal/ocr/sovaocr"
	"doc-watcher/internal/watcher"
	"github.com/stretchr/testify/assert"
)

func TestReadRawFileData(t *testing.T) {
	const Resources = "resources"
	const TestFile = "test_file_1.txt"
	const TestFileName = "test_file_1"

	timeoutDuration := time.Duration(10) * time.Second
	ocrService := sovaocr.New(&ocr.Config{
		Address:   "localhost:3451",
		EnableSSL: false,
		Timeout:   timeoutDuration,
	})

	mockServer := mock.InitMockServer()
	go func() {
		_ = mockServer.Start("localhost:3451")
		time.AfterFunc(5*time.Second, func() {
			_ = mockServer.Shutdown(context.Background())
		})
	}()

	t.Run("Read existing file", func(t *testing.T) {
		filePath := path.Join(Resources, "directory", TestFile)
		document := watcher.ParseCaughtFiles(filePath)[0]
		err := ocrService.Ocr.RecognizeFile(document, filePath)

		assert.NoError(t, err, "non null error pointer")
		assert.Equal(t, TestFileName, document.Content, "not equal file data")
	})

	t.Run("Read non existing file", func(t *testing.T) {
		filePath := path.Join(Resources, "directory", "any")
		documents := watcher.ParseCaughtFiles(filePath)
		assert.Empty(t, documents, "non null file data")
	})

	t.Run("Caught error while read directory as file", func(t *testing.T) {
		dirPath := path.Join(Resources, "directory")
		documents := watcher.ParseCaughtFiles(dirPath)
		assert.NotEmpty(t, documents, "non null file data")

		filePath := path.Join(dirPath, documents[0].DocumentName)
		err := ocrService.Ocr.RecognizeFile(documents[0], filePath)
		assert.NoError(t, err, "null error pointer for null file ptr")
		assert.NotEmpty(t, documents[0].Content, "non null file data")
	})
}
