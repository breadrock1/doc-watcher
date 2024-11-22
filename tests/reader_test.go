package tests

import (
	"path"
	"testing"

	"doc-watcher/internal/watcher"
	"github.com/stretchr/testify/require"
)

func TestParseCaughtFiles(t *testing.T) {
	const Resources = "resources"
	const TestFile = "test_file_1.txt"

	t.Run("Parse file by path", func(t *testing.T) {
		docPath := path.Join(Resources, "directory", TestFile)
		documents := watcher.ParseCaughtFiles(docPath)

		require.NotEmpty(t, documents, "empty list")
		require.Equal(t, len(documents), 1, "empty list")

		firstDocument := documents[0]
		require.Equal(t, firstDocument.DocumentName, TestFile)
	})

	t.Run("Parse entity files by dir path", func(t *testing.T) {
		dirPath := path.Join(Resources, "directory")
		documents := watcher.ParseCaughtFiles(dirPath)

		require.NotEmpty(t, documents, "empty list")
		require.Equal(t, len(documents), 9, "empty list")
	})

	t.Run("Parse non existing file", func(t *testing.T) {
		docPath := path.Join(Resources, "directory", "any")
		documents := watcher.ParseCaughtFiles(docPath)

		require.Empty(t, documents, "empty list")
		require.Equal(t, len(documents), 0, "incorrect list length")
	})
}
