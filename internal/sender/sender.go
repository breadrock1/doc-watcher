package sender

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

func GET(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("failed to send GET request:", err)
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func PUT(body *bytes.Buffer, url, mime string, timeout time.Duration) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		log.Println("failed to create PUT request:", err)
		return nil, err
	}
	req.Header.Set(echo.HeaderContentType, mime)

	client := &http.Client{Timeout: timeout}
	return sendRequest(client, req)
}

func POST(body *bytes.Buffer, url, mime string, timeout time.Duration) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Println("failed to create POST request:", err)
		return nil, err
	}
	req.Header.Set(echo.HeaderContentType, mime)

	client := &http.Client{Timeout: timeout}
	return sendRequest(client, req)
}

func sendRequest(client *http.Client, req *http.Request) ([]byte, error) {
	response, err := client.Do(req)
	if err != nil {
		log.Println("failed to send request:", err)
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	respData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("failed to read response body: ", err)
		return nil, err
	}

	if response.StatusCode > 200 {
		msg := fmt.Sprintf("failed response %s: %s", response.Status, string(respData))
		return nil, errors.New(msg)
	}

	return respData, nil
}

func CreateFormFile(fileHandle *os.File, reqBody *bytes.Buffer) (*multipart.Writer, error) {
	writer := multipart.NewWriter(reqBody)
	filePath := filepath.Base(fileHandle.Name())
	formFile, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		log.Println("Failed while creating form file: ", err)
		return nil, err
	}

	if _, err = io.Copy(formFile, fileHandle); err != nil {
		log.Println("Failed while coping file form part to file handle: ", err)
		return nil, err
	}

	if err = writer.Close(); err != nil {
		log.Println("Failed while closing req body writer: ", err)
		return nil, err
	}

	return writer, nil
}
