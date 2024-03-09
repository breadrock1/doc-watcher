package sender

import (
	"context"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"doc-notifier/test/mocked"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TestcaseDirPath = "../../testcases/"

func TestReadRawFileData(t *testing.T) {
	fileSender := &sender.FileSender{
		ReadRawFileFlag: true,
	}

	t.Run("Read existing file", func(t *testing.T) {
		data, err := fileSender.ReadRawFileData(TestcaseDirPath + "directory/test_file_1.txt")

		assert.NoError(t, err, "Returned non null error pointer")
		assert.Equal(t, data, "test_file_1", "Returned non equal file data")
	})

	t.Run("Read non existing file", func(t *testing.T) {
		data, err := fileSender.ReadRawFileData(TestcaseDirPath + "directory/any.txt")

		assert.Error(t, err, "Returned null error pointer for non existing file")
		assert.Empty(t, data, "Returned non null file data")
	})

	t.Run("Caught error while read directory as file", func(t *testing.T) {
		data, err := fileSender.ReadRawFileData(TestcaseDirPath + "directory")

		assert.Error(t, err, "Returned null error pointer for non file ptr")
		assert.Empty(t, data, "Returned non null file data")
	})
}

func TestRecognizeFileData(t *testing.T) {
	fileSender := &sender.FileSender{
		ReadRawFileFlag:   false,
		OrcServiceAddress: "http://localhost:3451",
	}

	t.Run("Recognize file data", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		data, err := fileSender.RecognizeFileData(TestcaseDirPath + "directory/test_file_1.txt")

		assert.NoError(t, err, "Returned non null error pointer")
		assert.Equal(t, data, "test_file_1", "Returned non equal file data")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error while recognize non existing file", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		data, err := fileSender.RecognizeFileData(TestcaseDirPath + "directory/any.txt")

		assert.Error(t, err, "Returned non null error pointer")
		assert.Empty(t, data, "Returned non equal file data")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error while recognize directory as file", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		data, err := fileSender.RecognizeFileData(TestcaseDirPath + "directory")

		assert.Error(t, err, "Returned non null error pointer")
		assert.Empty(t, data, "Returned non equal file data")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error with service denied", func(t *testing.T) {
		_fileSender := &sender.FileSender{
			ReadRawFileFlag:   false,
			OrcServiceAddress: "http://localhost:8080",
		}

		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		data, err := _fileSender.RecognizeFileData(TestcaseDirPath + "directory/test_file_1.txt")

		assert.Error(t, err, "Returned non null error pointer")
		assert.Empty(t, data, "Returned non equal file data")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})
}

func TestStoreDocument(t *testing.T) {
	fileSender := &sender.FileSender{
		ReadRawFileFlag: false,
		SearcherAddress: "http://localhost:3451",
	}

	t.Run("Store Document", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseDirPath + "directory/test_file_1.txt")
		storeErr := fileSender.StoreDocument(document)

		assert.NoError(t, parseErr, "Returned error while parsing file")
		assert.NoError(t, storeErr, "Returned error while storing document")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error while storing file", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseDirPath + "directory/test_file_2.txt")
		storeErr := fileSender.StoreDocument(document)

		assert.NoError(t, parseErr, "Returned non null error pointer")
		assert.Error(t, storeErr, "Returned non equal file data")

		time.AfterFunc(20*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error with service denied", func(t *testing.T) {
		_fileSender := &sender.FileSender{
			ReadRawFileFlag:   false,
			OrcServiceAddress: "http://localhost:8080",
		}

		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseDirPath + "directory/test_file_1.txt")
		storeErr := _fileSender.StoreDocument(document)

		assert.NoError(t, parseErr, "Returned non null error pointer")
		assert.Error(t, storeErr, "Returned non equal file data")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})
}

func TestComputeContentTokens(t *testing.T) {
	fileSender := &sender.FileSender{
		ReadRawFileFlag:   false,
		LlmServiceAddress: "http://localhost:3451",
	}

	t.Run("Compute tokens", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseDirPath + "directory/test_file_1.txt")
		document.Content = "test_file_1"
		document.DocumentName = "test_file_1.txt"
		tokens, computeErr := fileSender.ComputeContentTokens(document)

		assert.NoError(t, parseErr, "Returned error while parsing file")
		assert.NoError(t, computeErr, "Returned error while storing document")
		assert.Equal(t, tokens.Chunks, 1, "Non correct returned chunks size")
		assert.Equal(t, tokens.ChunkedText[0], "test_file_1", "Non correct returned chunks data")
		assert.Equal(t, tokens.Vectors[0], []float64{0.345, 0.045}, "Non correct returned vectors")

		time.AfterFunc(200*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error while computing tokens", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseDirPath + "directory/test_file_2.txt")
		_, computeErr := fileSender.ComputeContentTokens(document)

		assert.NoError(t, parseErr, "Returned error while parsing file")
		assert.Error(t, computeErr, "Returned error while storing document")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error with service denied", func(t *testing.T) {
		_fileSender := &sender.FileSender{
			ReadRawFileFlag:   false,
			LlmServiceAddress: "http://localhost:3451",
		}

		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseDirPath + "directory/test_file_2.txt")
		_, computeErr := _fileSender.ComputeContentTokens(document)

		assert.NoError(t, parseErr, "Returned error while parsing file")
		assert.Error(t, computeErr, "Returned error while storing document")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})
}
