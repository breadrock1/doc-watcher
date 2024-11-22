package tests

import (
	"context"
	"path"
	"testing"
	"time"

	"doc-watcher/internal/embeddings"
	"doc-watcher/internal/embeddings/sovavec"
	"doc-watcher/internal/watcher"
	"doc-watcher/tests/mock"
	"github.com/stretchr/testify/assert"
)

func TestComputeContentTokens(t *testing.T) {
	const Resources = "resources"
	const TestFile = "test_file_1.txt"
	const TestFileName = "test_file_1"

	tokenizerService := sovavec.New(&embeddings.Config{
		Address:      "localhost:3451",
		EnableSSL:    false,
		ChunkSize:    500,
		ChunkOverlap: 1,
		ReturnChunks: false,
		ChunkBySelf:  false,
	})

	mockServer := mock.InitMockServer()
	go func() {
		_ = mockServer.Start("localhost:3451")
		time.AfterFunc(5*time.Second, func() {
			_ = mockServer.Shutdown(context.Background())
		})
	}()

	t.Run("Compute tokens", func(t *testing.T) {

		filePath := path.Join(Resources, "directory", TestFile)
		document, parseErr := watcher.ParseFile(filePath)
		document.Content = TestFileName
		document.DocumentName = TestFile
		tokens, computeErr := tokenizerService.Tokenizer.Tokenize(document.Content)

		assert.NoError(t, parseErr, "Returned error while parsing file")
		assert.NoError(t, computeErr, "Returned error while storing document")
		assert.Equal(t, 1, tokens.Chunks, "Non correct returned chunks size")
		assert.Equal(t, "test_file_1", tokens.ChunkedText[0], "Non correct returned chunks data")
		assert.Equal(t, []float64{0.345, 0.045}, tokens.Vectors[0], "Non correct returned vectors")
	})

	t.Run("Caught error while computing tokens", func(t *testing.T) {
		filePath := path.Join(Resources, "directory", "any")
		_, parseErr := watcher.ParseFile(filePath)
		assert.Error(t, parseErr, "error while parsing file")
	})

	t.Run("Caught error with service denied", func(t *testing.T) {
		eMockServ := mock.InitMockServer()
		go func() {
			_ = eMockServ.Start("localhost:3451")
			time.AfterFunc(2*time.Second, func() {
				_ = eMockServ.Shutdown(context.Background())
			})

		}()

		filePath := path.Join(Resources, "directory", "any")
		_, parseErr := watcher.ParseFile(filePath)
		assert.Error(t, parseErr, "error while parsing file")
	})
}
