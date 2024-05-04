package sender

import (
	"doc-notifier/internal/pkg/ocr"
	"doc-notifier/internal/pkg/reader"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TestcaseDirPath = "../../testcases/"
const TestcaseOtherDirPath = "../../testcases/directory/"
const TestcaseFilePath = TestcaseOtherDirPath + "test_file_1.txt"
const TestcaseNonExistingFilePath = TestcaseOtherDirPath + "any_file.txt"

func TestReadRawFileData(t *testing.T) {
	timeoutDuration := time.Duration(10) * time.Second
	readerService := reader.New()
	ocrService := ocr.New(&ocr.Options{
		Mode:    ocr.GetModeFromString("read-raw-file"),
		Address: "http://localhost:3451",
		Timeout: timeoutDuration,
	})

	t.Run("Read existing file", func(t *testing.T) {
		document := readerService.ParseCaughtFiles(TestcaseDirPath + "directory/test_file_1.txt")[0]
		data, err := ocrService.Ocr.RecognizeFile(document)

		assert.NoError(t, err, "Returned non null error pointer")
		assert.Equal(t, data, "test_file_1", "Returned non equal file data")
	})

	//t.Run("Read non existing file", func(t *testing.T) {
	//	document := readerService.ParseCaughtFiles(TestcaseNonExistingFilePath)[0]
	//	data, err := ocrService.Ocr.RecognizeFile(document)
	//
	//	assert.Error(t, err, "Returned null error pointer for non existing file")
	//	assert.Empty(t, data, "Returned non null file data")
	//})

	//t.Run("Caught error while read directory as file", func(t *testing.T) {
	//	document := readerService.ParseCaughtFiles(TestcaseNonExistingFilePath)[0]
	//	data, err := ocrService.Ocr.RecognizeFile(document)
	//
	//	assert.Error(t, err, "Returned null error pointer for non file ptr")
	//	assert.Empty(t, data, "Returned non null file data")
	//})
}

//func TestRecognizeFileData(t *testing.T) {
//	timeoutDuration := time.Duration(10) * time.Second
//	readerService := reader.New()
//	ocrService := ocr.New(&ocr.Options{
//		Mode:    ocr.GetModeFromString("read-raw-file"),
//		Address: "http://localhost:3451",
//		Timeout: timeoutDuration,
//	})
//
//	t.Run("Recognize file data", func(t *testing.T) {
//		e := mocked.CreateMockedServer()
//		go func() {
//			_ = e.Start("localhost:3451")
//		}()
//
//		data, err := ocrService.Ocr.RecognizeFile(TestcaseFilePath)
//
//		assert.NoError(t, err, "Returned non null error pointer")
//		assert.Equal(t, data, "test_file_1", "Returned non equal file data")
//
//		time.AfterFunc(2*time.Second, func() {
//			_ = e.Shutdown(context.Background())
//		})
//	})
//
//	t.Run("Caught error while recognize non existing file", func(t *testing.T) {
//		e := mocked.CreateMockedServer()
//		go func() {
//			_ = e.Start("localhost:3451")
//		}()
//
//		data, err := ocrService.Ocr.RecognizeFile(TestcaseNonExistingFilePath)
//
//		assert.Error(t, err, "Returned non null error pointer")
//		assert.Empty(t, data, "Returned non equal file data")
//
//		time.AfterFunc(2*time.Second, func() {
//			_ = e.Shutdown(context.Background())
//		})
//	})
//
//	t.Run("Caught error while recognize directory as file", func(t *testing.T) {
//		e := mocked.CreateMockedServer()
//		go func() {
//			_ = e.Start("localhost:3451")
//		}()
//
//		data, err := ocrService.Ocr.RecognizeFile(TestcaseOtherDirPath)
//
//		assert.Error(t, err, "Returned non null error pointer")
//		assert.Empty(t, data, "Returned non equal file data")
//
//		time.AfterFunc(2*time.Second, func() {
//			_ = e.Shutdown(context.Background())
//		})
//	})
//
//	t.Run("Caught error with service denied", func(t *testing.T) {
//		ocrService := ocr.New(&ocr.Options{
//			Mode:    ocr.GetModeFromString("assistant"),
//			Address: "http://localhost:4444",
//			Timeout: timeoutDuration,
//		})
//
//		e := mocked.CreateMockedServer()
//		go func() {
//			_ = e.Start("localhost:3451")
//		}()
//
//		data, err := ocrService.Ocr.RecognizeFile(TestcaseFilePath)
//
//		assert.Error(t, err, "Returned non null error pointer")
//		assert.Empty(t, data, "Returned non equal file data")
//
//		time.AfterFunc(2*time.Second, func() {
//			_ = e.Shutdown(context.Background())
//		})
//	})
//}
