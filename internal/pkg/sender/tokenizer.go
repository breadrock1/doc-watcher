package sender

import (
	"bytes"
	"doc-notifier/internal/pkg/reader"
	"encoding/json"
	"log"
	"strings"
)

const EmbeddingsURL = "/api/v1/vectorizer/get_vectors"

type TokenizedVectors struct {
	Chunks      int         `json:"chunks"`
	ChunkedText []string    `json:"chunked_text"`
	Vectors     [][]float64 `json:"vectors"`
}

func CreateTokenizedVectors() *TokenizedVectors {
	return &TokenizedVectors{
		Chunks:      0,
		ChunkedText: []string{},
		Vectors:     [][]float64{},
	}
}

type TokenizerForm struct {
	Text              string `json:"text"`
	ChunkSize         int    `json:"chunk_size"`
	ChunkOverlap      int    `json:"chunk_overlap"`
	ReturnChunkedText bool   `json:"return_chunked_text"`
}

func CreateTokenizerForm(orData string) *TokenizerForm {
	return &TokenizerForm{
		Text:              orData,
		ChunkSize:         800,
		ChunkOverlap:      100,
		ReturnChunkedText: false,
	}
}

func (fs *FileSender) ComputeContentTokens(document *reader.Document) (*TokenizedVectors, error) {
	tokenizedVector := CreateTokenizedVectors()

	orData := strings.ReplaceAll(document.Content, "\n", " ")
	textVectors := CreateTokenizerForm(orData)

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		log.Println("Failed while marshaling text vectors: ", err)
		return tokenizedVector, err
	}

	reqBody := bytes.NewBuffer(jsonData)
	targetUrl := fs.LlmServiceAddress + EmbeddingsURL
	log.Printf("Sending file %s to extract tokens", document.DocumentPath)
	respData, err := fs.sendRequest(reqBody, &targetUrl)
	if err != nil {
		log.Println("Failed while sending request: ", err)
		return tokenizedVector, err
	}

	_ = json.Unmarshal(respData, tokenizedVector)
	return tokenizedVector, nil
}
