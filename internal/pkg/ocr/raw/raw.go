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

	stringData := string(bytesData)
	if len(stringData) == 0 {
		log.Println("Failed: returned empty string data...")
		return "", err
	}

	return stringData, nil
}

func (re *Service) RecognizeFileData(data []byte) (string, error) {
	return string(data), nil
}
