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

func SendRequest(body *bytes.Buffer, url *string, formData string) ([]byte, error) {
	req, err := http.NewRequest("POST", *url, body)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set(echo.HeaderContentType, formData)

	client := &http.Client{Timeout: 120 * time.Second}
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
