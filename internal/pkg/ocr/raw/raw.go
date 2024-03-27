package raw

import (
	"log"
	"os"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (re *Service) RecognizeFile(filePath string) (string, error) {
	bytesData, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Failed while reading file: ", err)
		return "", err
	}

	return string(bytesData), nil
}

func (re *Service) RecognizeFileData(data []byte) (string, error) {
	return string(data), nil
}
