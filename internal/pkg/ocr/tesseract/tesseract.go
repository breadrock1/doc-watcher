package tesseract

import (
	"github.com/otiai10/gosseract/v2"
)

type TesseractOCR struct {
	client *gosseract.Client
}

func New() *TesseractOCR {
	client := gosseract.NewClient()
	return &TesseractOCR{
		client: client,
	}
}

func (gc *TesseractOCR) RecognizeFile(filePath string) (string, error) {
	_ = gc.client.SetImage(filePath)
	return gc.client.Text()
}

func (gc *TesseractOCR) RecognizeFileData(data []byte) (string, error) {
	_ = gc.client.SetImageFromBytes(data)
	return gc.client.Text()
}
