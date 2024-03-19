package raw

import (
	"log"
	"os"
)

type RawExractor struct {
}

func New() *RawExractor {
	return &RawExractor{}
}

func (re *RawExractor) RecognizeFile(filePath string) (string, error) {
	bytesData, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Failed while reading file: ", err)
		return "", err
	}

	return string(bytesData), nil
}

func (re *RawExractor) RecognizeFileData(data []byte) (string, error) {
	return string(data), nil
}
