package sender

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"notifier/internal/pkg/reader"
)

const SearcherURL = "/searcher/document/new"
const AssistantURL = "/api/assistant/extract-file/"

type FileSender struct {
	AssistantAddress string
	SearcherAddress  string
}

type DocumentText struct {
	Context string `json:"context"`
}

func New(searcherAddr string, assistantAddr string) *FileSender {
	return &FileSender{
		AssistantAddress: assistantAddr,
		SearcherAddress:  searcherAddr,
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
	var resTest = struct {
		Context string `json:"context"`
	}{}
	_ = json.Unmarshal(bodyData, &resTest)
	return resTest.Context, nil
}

func (s *FileSender) sendRequest(body *bytes.Buffer, writer *multipart.Writer) (*http.Response, error) {
	targetUrl := s.AssistantAddress + AssistantURL
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
