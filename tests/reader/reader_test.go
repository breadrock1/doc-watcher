package reader

import (
	"doc-notifier/internal/reader"
	"github.com/stretchr/testify/require"
	"testing"
)

const TestcaseDirPath = "../testcases/"

func TestParseCaughtFiles(t *testing.T) {
	fileReader := &reader.Service{}

	t.Run("Parse file by path", func(t *testing.T) {
		documents := fileReader.ParseCaughtFiles(TestcaseDirPath + "directory/test_file_1.txt")

		require.NotEmpty(t, documents, "Empty list")
		require.Equal(t, len(documents), 1, "Empty list")

		firstDocument := documents[0]
		require.Equal(t, firstDocument.DocumentName, "test_file_1.txt")
	})

	t.Run("Parse entity files by dir path", func(t *testing.T) {
		documents := fileReader.ParseCaughtFiles(TestcaseDirPath + "directory")

		require.NotEmpty(t, documents, "Empty list")
		require.Equal(t, len(documents), 9, "Empty list")
	})

	t.Run("Parse unexist file", func(t *testing.T) {
		documents := fileReader.ParseCaughtFiles(TestcaseDirPath + "any")

		require.Empty(t, documents, "Empty list")
		require.Equal(t, len(documents), 0, "Not correct docs length")
	})
}
