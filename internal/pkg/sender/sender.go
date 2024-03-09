package sender

import (
	"bytes"
	"errors"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"time"
)

type FileSender struct {
	ReadRawFileFlag   bool
	OrcServiceAddress string
	SearcherAddress   string
	LlmServiceAddress string
}

func New(searcherAddr string, assistantAddr string, llmAddress string, readRawFlag bool) *FileSender {
	return &FileSender{
		ReadRawFileFlag:   readRawFlag,
		OrcServiceAddress: assistantAddr,
		SearcherAddress:   searcherAddr,
		LlmServiceAddress: llmAddress,
	}
}

func (fs *FileSender) sendRequest(body *bytes.Buffer, url *string) ([]byte, error) {
	req, err := http.NewRequest("POST", *url, body)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := &http.Client{Timeout: 120 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error while creating request:", err)
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode > 200 {
		log.Println("Non Ok response status: ", response.Status)
		return nil, errors.New("non 200 response code status")
	}

	respData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed while reading response reqBody: ", err)
		return nil, err
	}

	return respData, nil
}
