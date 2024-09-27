package raw

import (
	"fmt"
	"log"
	"os"

	"doc-notifier/internal/models"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) RecognizeFile(document *models.Document, filePath string) error {
	bytesData, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Failed while reading file: ", err)
		return err
	}

	stringData := string(bytesData)
	if len(stringData) == 0 {
		return fmt.Errorf("returned empty content data for: %s", document.DocumentID)
	}

	document.SetContentData(stringData)
	document.SetQuality(10000)
	return nil
}
