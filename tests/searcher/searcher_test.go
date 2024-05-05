package searcher

import (
	"context"
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/searcher"
	"doc-notifier/tests/mocked"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TestcaseOtherDirPath = "../testcases/directory/"
const TestcaseFilePath = TestcaseOtherDirPath + "test_file_1.txt"
const TestcaseNonExistingFilePath = TestcaseOtherDirPath + "any_file.txt"

func TestStoreDocument(t *testing.T) {
	timeoutDuration := time.Duration(10) * time.Second
	searcherService := searcher.New("http://localhost:3451", timeoutDuration)

	t.Run("Store Document", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseFilePath)
		storeErr := searcherService.StoreDocument(document)

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

		_, parseErr := reader.ParseFile(TestcaseNonExistingFilePath)
		//storeErr := searcherService.StoreDocument(document)

		assert.Error(t, parseErr, "Returned non null error pointer")
		assert.Error(t, parseErr, "Returned non equal file data")

		time.AfterFunc(20*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error with service denied", func(t *testing.T) {
		searcherService := searcher.New("http://localhost:4444", timeoutDuration)

		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		document, parseErr := reader.ParseFile(TestcaseFilePath)
		storeErr := searcherService.StoreDocument(document)

		assert.NoError(t, parseErr, "Returned non null error pointer")
		assert.Error(t, storeErr, "Returned non equal file data")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})
}
