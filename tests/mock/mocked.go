package mock

import (
	"doc-watcher/internal/embeddings/sovavec"
	"doc-watcher/internal/ocr/sovaocr"
	"encoding/json"
	"log"

	"doc-watcher/internal/embeddings"
	"doc-watcher/internal/watcher"
	"github.com/labstack/echo/v4"
)

type DocumentForm struct {
	Content string `json:"text"`
}

type TokenizerForm struct {
}

func InitMockServer() *echo.Echo {
	e := echo.New()
	e.POST(sovaocr.RecognitionURL, RecognizeFile)
	e.POST(sovavec.EmbeddingsAssistantURL, ComputeTokens)
	e.PUT("/storage/folders/common-folder/documents/c31964293145484954679b19a114188e", StoreDocument)

	return e
}

func RecognizeFile(c echo.Context) error {
	document := DocumentForm{Content: "test_file_1"}
	log.Println("got request to recognize file: ", document.Content)
	return c.JSON(200, document)
}

func StoreDocument(c echo.Context) error {
	document := &watcher.Document{}
	decoder := json.NewDecoder(c.Request().Body)
	_ = decoder.Decode(document)

	log.Println("Got request to store doc: ", document.DocumentName)
	if document.DocumentName != "test_file_1.txt" {
		log.Println("Non correct doc: ", document.DocumentName)
		return c.JSON(403, document)
	}

	return c.JSON(200, document)
}

func ComputeTokens(c echo.Context) error {
	tokensForm := &embeddings.EmbedAllForm{}
	decoder := json.NewDecoder(c.Request().Body)
	_ = decoder.Decode(tokensForm)

	log.Println("Got request to tokenize doc: ", tokensForm.Inputs)
	if tokensForm.Inputs != "test_file_1" {
		log.Println("Non correct doc: ", tokensForm.Inputs)
		return c.JSON(403, tokensForm)
	}

	return c.JSON(200, [][]float64{{0.345, 0.045}})
}
