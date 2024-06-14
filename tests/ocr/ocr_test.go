package sender

import (
	"doc-notifier/internal/config"
	ocr2 "doc-notifier/internal/ocr"
	"doc-notifier/internal/reader"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TestcaseDirPath = "../testcases/"
const TestcaseOtherDirPath = "../testcases/directory/"
const TestcaseNonExistingFilePath = TestcaseOtherDirPath + "any_file.txt"

func TestReadRawFileData(t *testing.T) {
	readerService := reader.New()
	timeoutDuration := time.Duration(10) * time.Second
	ocrService := ocr2.New(&config.OcrConfig{
		Mode:    "raw",
		Address: "http://localhost:3451",
		Timeout: timeoutDuration,
	})

	t.Run("Read existing file", func(t *testing.T) {
		document := readerService.ParseCaughtFiles(TestcaseDirPath + "directory/test_file_1.txt")[0]
		err := ocrService.Ocr.RecognizeFile(document)

		assert.NoError(t, err, "Returned non null error pointer")
		assert.Equal(t, document.Content, "test_file_1", "Returned non equal file data")
	})

	t.Run("Read non existing file", func(t *testing.T) {
		documents := readerService.ParseCaughtFiles(TestcaseNonExistingFilePath)
		assert.Empty(t, documents, "Returned non null file data")
	})

	t.Run("Caught error while read directory as file", func(t *testing.T) {
		documents := readerService.ParseCaughtFiles(TestcaseOtherDirPath)
		assert.NotEmpty(t, documents, "Returned non-null file data")

		err := ocrService.Ocr.RecognizeFile(documents[0])
		assert.NoError(t, err, "Returned null error pointer for non file ptr")
		assert.NotEmpty(t, documents[0].Content, "Returned non-null file data")
	})
}
