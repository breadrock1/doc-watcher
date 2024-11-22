package watcher

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	documentMimes = []string{
		"csv", "msword", "html", "json", "pdf",
		"rtf", "plain", "vnd.ms-excel", "xml",
		"vnd.ms-powerpoint", "vnd.oasis.opendocument.text",
		"vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"vnd.openxmlformats-officedocument.wordprocessingml.document",
		"vnd.openxmlformats-officedocument.presentationml.presentation",
	}
)

func ParseDocumentType(extension string) string {
	mimeType := mime.TypeByExtension(extension)
	attributes := strings.Split(mimeType, "/")
	switch attributes[0] {
	case "audio":
		return "audio"
	case "image":
		return "image"
	case "video":
		return "video"
	case "text":
		return "document"
	case "application":
		return extractApplicationMimeType(attributes[1])
	default:
		return "document"
	}
}

func extractApplicationMimeType(attribute string) string {
	for _, mimeType := range documentMimes {
		if mimeType == attribute {
			return "document"
		}
	}

	return "document"
}

func GetEntityFiles(filePath string) []string {
	<-time.After(time.Second)

	var files []string
	if err := filepath.Walk(filePath, VisitEntity(&files)); err != nil {
		log.Println("files walking error: ", err)
	}
	return files
}

func VisitEntity(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to visit entity: %w", err)
		}

		if !info.IsDir() {
			*files = append(*files, path)
		}

		return nil
	}
}
