package tokenizer

import (
	"context"
	"testing"
	"time"

	"doc-notifier/internal/config"
	"doc-notifier/internal/reader"
	"doc-notifier/internal/tokenizer"
	"doc-notifier/tests/mocked"
	"github.com/stretchr/testify/assert"
)

const TestcaseOtherDirPath = "../testcases/directory/"
const TestcaseFilePath = TestcaseOtherDirPath + "test_file_1.txt"
const TestcaseNonExistingFilePath = TestcaseOtherDirPath + "any_file.txt"

func TestComputeContentTokens(t *testing.T) {
	timeoutDuration := time.Duration(10) * time.Second
	tokenizerService := tokenizer.New(&config.TokenizerConfig{
		Address:      "http://localhost:3451",
		Mode:         "assistant",
		ChunkSize:    500,
		ChunkOverlap: 1,
		ReturnChunks: false,
		ChunkBySelf:  false,
		Timeout:      timeoutDuration,
	})

	t.Run("Compute tokens", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseFilePath)
		document.Content = "test_file_1"
		document.DocumentName = "test_file_1.txt"
		tokens, computeErr := tokenizerService.Tokenizer.TokenizeTextData(document.Content)

		assert.NoError(t, parseErr, "Returned error while parsing file")
		assert.NoError(t, computeErr, "Returned error while storing document")
		assert.Equal(t, 1, tokens.Chunks, "Non correct returned chunks size")
		assert.Equal(t, "test_file_1", tokens.ChunkedText[0], "Non correct returned chunks data")
		assert.Equal(t, []float64{0.345, 0.045}, tokens.Vectors[0], "Non correct returned vectors")

		time.AfterFunc(200*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error while computing tokens", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		_, parseErr := reader.ParseFile(TestcaseNonExistingFilePath)
		//_, computeErr := tokenizerService.Tokenizer.TokenizeTextData(document.Content)

		assert.Error(t, parseErr, "Returned error while parsing file")
		//assert.Error(t, computeErr, "Returned error while storing document")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error with service denied", func(t *testing.T) {
		_ = tokenizer.New(&config.TokenizerConfig{
			Address:      "http://localhost:4444",
			Mode:         "assistant",
			ChunkSize:    0,
			ChunkOverlap: 0,
			ReturnChunks: false,
			ChunkBySelf:  false,
			Timeout:      timeoutDuration,
		})

		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		_, parseErr := reader.ParseFile(TestcaseNonExistingFilePath)
		//_, computeErr := tokenizerService.Tokenizer.TokenizeTextData(document.Content)

		assert.Error(t, parseErr, "Returned error while parsing file")
		//assert.Error(t, computeErr, "Returned error while storing document")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})
}
