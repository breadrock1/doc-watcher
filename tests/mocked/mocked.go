package mocked

import (
	"doc-notifier/internal/models"
	ocr "doc-notifier/internal/ocr/assistant"
	tokenizer "doc-notifier/internal/tokenizer/assistant"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
)

type DocumentForm struct {
	Context string
}

type TokenizerForm struct {
}

func CreateMockedServer() *echo.Echo {
	e := echo.New()
	e.POST(ocr.RecognitionURL, RecognizeFile)
	e.POST(tokenizer.EmbeddingsAssistantURL, ComputeTokens)
	e.PUT("/storage/folders/common-folder/documents/c31964293145484954679b19a114188e", StoreDocument)

	return e
}

func RecognizeFile(c echo.Context) error {
	document := DocumentForm{Context: "test_file_1"}
	log.Println("Got request to recognize file: ", document.Context)
	return c.JSON(200, document)
}

func StoreDocument(c echo.Context) error {
	document := &models.Document{}
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
	tokensForm := &models.EmbedAllForm{}
	decoder := json.NewDecoder(c.Request().Body)
	_ = decoder.Decode(tokensForm)

	log.Println("Got request to tokenize doc: ", tokensForm.Inputs)
	if tokensForm.Inputs != "test_file_1" {
		log.Println("Non correct doc: ", tokensForm.Inputs)
		return c.JSON(403, tokensForm)
	}

	tokenizedVector := models.ComputedTokens{
		Chunks:      1,
		ChunkedText: []string{"test_file_1"},
		Vectors:     [][]float64{{0.345, 0.045}},
	}

	return c.JSON(200, tokenizedVector)
}
