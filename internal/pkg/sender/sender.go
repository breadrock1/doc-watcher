package sender

import (
	"bytes"
	"doc-notifier/internal/pkg/reader"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const SearcherURL = "/searcher/document/new"
const EmbeddingsURL = "/api/v1/get_text_vectors"
const RecognitionURL = "/api/assistant/extract-file/"

type FileSender struct {
	OrcServiceAddress string
	SearcherAddress   string
	LlmServiceAddress string
}

type TokenizedVectors struct {
	Chunks      int         `json:"chunks"`
	ChunkedText [][]string  `json:"chunked_text"`
	Vectors     [][]float64 `json:"vectors"`
}

type DocumentForm struct {
	Context string `json:"context"`
}

type TokenizerForm struct {
	Text              string `json:"text"`
	ChunkSize         int    `json:"chunk_size"`
	ReturnChunkedText bool   `json:"return_chunked_text"`
}

func New(searcherAddr string, assistantAddr string, llmAddress string) *FileSender {
	return &FileSender{
		OrcServiceAddress: assistantAddr,
		SearcherAddress:   searcherAddr,
		LlmServiceAddress: llmAddress,
	}
}

func (s *FileSender) StoreDocument(document *reader.Document) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		log.Println(err)
	}

	body := bytes.NewBuffer(jsonData)
	targetAddress := s.SearcherAddress + SearcherURL
	req, err := http.NewRequest("POST", targetAddress, body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		log.Println("Error making request:", err)
		return nil
	}

	log.Println(resp.StatusCode, resp.Body)
	return nil
}

func (s *FileSender) ComputeContentTokens(document *reader.Document) *TokenizedVectors {
	textVectors := &TokenizerForm{
		Text:              document.Content,
		ChunkSize:         len(document.Content),
		ReturnChunkedText: true,
	}

	jsonData, err := json.Marshal(textVectors)
	if err != nil {
		log.Println(err)
		return nil
	}

	body := bytes.NewBuffer(jsonData)
	targetAddress := s.LlmServiceAddress + EmbeddingsURL
	req, err := http.NewRequest("POST", targetAddress, body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		log.Println("Error making request:", err)
		return nil
	}

	bodyData, err := io.ReadAll(resp.Body)
	var resTest = &TokenizedVectors{}
	_ = json.Unmarshal(bodyData, resTest)
	return resTest
}

func (s *FileSender) RecognizeFileData(filePath string) (string, error) {
	fileHandle, err := os.Open(filePath)
	defer func() { _ = fileHandle.Close() }()
	if err != nil {
		msg := "Failed while opening file: "
		log.Println(msg, err)
		return "", err
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		msg := "Failed while creating form file: "
		log.Println(msg, err)
		return "", err
	}

	_, err = io.Copy(part, fileHandle)
	if err != nil {
		msg := "Failed while coping file form part to file handle: "
		log.Println(msg, err)
		return "", err
	}
	_ = writer.Close()

	resp, err := s.sendRequest(&body, writer)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		return "", err
	}

	bodyData, err := io.ReadAll(resp.Body)
	var resTest = &DocumentForm{}
	_ = json.Unmarshal(bodyData, resTest)
	return resTest.Context, nil
}

func (s *FileSender) sendRequest(body *bytes.Buffer, writer *multipart.Writer) (*http.Response, error) {
	targetUrl := s.OrcServiceAddress + RecognitionURL
	req, err := http.NewRequest("POST", targetUrl, body)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		log.Println("Error making request:", err)
		return nil, err
	}

	return resp, nil
}
