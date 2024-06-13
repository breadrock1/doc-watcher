package sender

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

func SendRequest(body *bytes.Buffer, url, method, mime *string, timeout time.Duration) ([]byte, error) {
	req, err := http.NewRequest(*method, *url, body)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set(echo.HeaderContentType, *mime)

	//client := &http.Client{Timeout: timeout}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error while creating request:", err)
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	respData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed while reading response reqBody: ", err)
		return nil, err
	}

	if response.StatusCode > 200 {
		log.Printf("Non Ok response status %s: %s", response.Status, string(respData))
		return nil, errors.New("non 200 response code status")
	}

	return respData, nil
}

func SendGETRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed while reading response reqBody: ", err)
		return nil, err
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
