package sender

import (
	"context"
	"doc-notifier/internal/pkg/ocr"
	"doc-notifier/tests/mocked"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TestcaseDirPath = "../../testcases/"
const TestcaseOtherDirPath = "../../testcases/directory/"
const TestcaseFilePath = TestcaseOtherDirPath + "test_file_1.txt"
const TestcaseNonExistingFilePath = TestcaseOtherDirPath + "any_file.txt"

func TestReadRawFileData(t *testing.T) {
	ocrService := ocr.New(&ocr.Options{
		Mode:    ocr.GetModeFromString("read-raw-file"),
		Address: "http://localhost:3451",
	})

	t.Run("Read existing file", func(t *testing.T) {
		data, err := ocrService.Ocr.RecognizeFile(TestcaseDirPath + "directory/test_file_1.txt")

		assert.NoError(t, err, "Returned non null error pointer")
		assert.Equal(t, data, "test_file_1", "Returned non equal file data")
	})

	t.Run("Read non existing file", func(t *testing.T) {
		data, err := ocrService.Ocr.RecognizeFile(TestcaseNonExistingFilePath)

		assert.Error(t, err, "Returned null error pointer for non existing file")
		assert.Empty(t, data, "Returned non null file data")
	})

	t.Run("Caught error while read directory as file", func(t *testing.T) {
		data, err := ocrService.Ocr.RecognizeFile(TestcaseDirPath + "directory")

		assert.Error(t, err, "Returned null error pointer for non file ptr")
		assert.Empty(t, data, "Returned non null file data")
	})
}

func TestRecognizeFileData(t *testing.T) {
	ocrService := ocr.New(&ocr.Options{
		Mode:    ocr.GetModeFromString("read-raw-file"),
		Address: "http://localhost:3451",
	})

	t.Run("Recognize file data", func(t *testing.T) {
		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		data, err := ocrService.Ocr.RecognizeFile(TestcaseFilePath)

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

		data, err := ocrService.Ocr.RecognizeFile(TestcaseNonExistingFilePath)

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

		data, err := ocrService.Ocr.RecognizeFile(TestcaseOtherDirPath)

		assert.Error(t, err, "Returned non null error pointer")
		assert.Empty(t, data, "Returned non equal file data")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})

	t.Run("Caught error with service denied", func(t *testing.T) {
		ocrService := ocr.New(&ocr.Options{
			Mode:    ocr.GetModeFromString("read-raw-file"),
			Address: "http://localhost:3451",
		})

		e := mocked.CreateMockedServer()
		go func() {
			_ = e.Start("localhost:3451")
		}()

		data, err := ocrService.Ocr.RecognizeFile(TestcaseFilePath)

		assert.Error(t, err, "Returned non null error pointer")
		assert.Empty(t, data, "Returned non equal file data")

		time.AfterFunc(2*time.Second, func() {
			_ = e.Shutdown(context.Background())
		})
	})
}
