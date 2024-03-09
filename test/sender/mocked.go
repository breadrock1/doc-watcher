package sender

import (
	"doc-notifier/internal/pkg/reader"
	"doc-notifier/internal/pkg/sender"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
)

func CreateMockedServer() *echo.Echo {
	e := echo.New()
	e.POST(sender.RecognitionURL, RecognizeFile)
	e.POST(sender.EmbeddingsURL, ComputeTokens)
	e.POST(sender.SearcherURL, StoreDocument)

	return e
}

func RecognizeFile(c echo.Context) error {
	document := sender.DocumentForm{
		Context: "test_file_1",
	}
	log.Println("Got request to recognize file: ", document.Context)
	return c.JSON(200, document)
}

func StoreDocument(c echo.Context) error {
	document := &reader.Document{}
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
	tokensForm := &sender.TokenizerForm{}
	decoder := json.NewDecoder(c.Request().Body)
	_ = decoder.Decode(tokensForm)

	log.Println("Got request to tokenize doc: ", tokensForm.Text)
	if tokensForm.Text != "test_file_1" {
		log.Println("Non correct doc: ", tokensForm.Text)
		return c.JSON(403, tokensForm)
	}

	tokenizedVector := sender.CreateTokenizedVectors()
	tokenizedVector.Chunks = 1
	tokenizedVector.ChunkedText = []string{"test_file_1"}
	tokenizedVector.Vectors = [][]float64{{0.345, 0.045}}

	return c.JSON(200, tokenizedVector)
}
