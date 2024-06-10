package reader

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func getEntityFiles(filePath string) []string {
	<-time.After(time.Second)

	var files []string
	err := filepath.Walk(filePath, visitEntity(&files))
	if err != nil {
		log.Println("Error while walking: ", err)
	}
	return files
}

func visitEntity(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		if !info.IsDir() {
			*files = append(*files, path)
		}

		return nil
	}
}

func moveFileTo(src, dst string) error {
	var moveErr error
	var srcFile, dstFile *os.File

	if srcFile, moveErr = os.Open(src); moveErr != nil {
		return fmt.Errorf("failed while open file %s: %e", src, moveErr)
	}

	if dstFile, moveErr = os.Create(dst); moveErr != nil {
		_ = srcFile.Close()
		return fmt.Errorf("failed while create file %s: %e", dst, moveErr)
	}
	defer func() { _ = dstFile.Close() }()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		_ = srcFile.Close()
		return fmt.Errorf("failed while coping data: %e", err)
	}

	_ = srcFile.Close()
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed while remove src file %s: %e", src, err)
	}

	return nil
}
