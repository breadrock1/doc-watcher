package reader

import (
	"fmt"
	"testing"
)

func TestParseCaughtFiles(t *testing.T) {
	testDirPath := "/Users/breadrock/Projects/internal/testcases"

	t.Run("dir parsing", func(t *testing.T) {
		docs := ParseCaughtFiles(testDirPath)
		fmt.Println("Done", docs)
	})

}
