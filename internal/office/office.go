package office

import (
	"fmt"
	"log"
	"os"

	"doc-notifier/internal/config"
	"doc-notifier/internal/sender"
)

type Service struct {
	address string
}

func New(config *config.OfficeConfig) *Service {
	return &Service{
		address: config.Address,
	}
}

func (s *Service) DownloadDocument(fileName string) error {
	targetUrl := fmt.Sprintf("%s/download?fileName=%s", s.address, fileName)
	fileData, err := sender.SendGETRequest(targetUrl)
	if err != nil {
		log.Println("failed to download file from office: ")
		return err
	}

	filePath := fmt.Sprintf("./indexer/watcher/%s", fileName)
	err = os.WriteFile(filePath, fileData, os.ModePerm)
	if err != nil {
		log.Println("failed to write file: ", err)
	}

	return err
}
