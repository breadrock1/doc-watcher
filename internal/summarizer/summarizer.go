package summarizer

import (
	"bytes"
	"context"
	"doc-notifier/internal/models"
	"doc-notifier/internal/sender"
	"doc-notifier/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"strings"
	"time"

	"doc-notifier/internal/config"
	_ "github.com/lib/pq"
)

type Service struct {
	dbAddress  string
	llmAddress string
	db         storage.ServiceStorage
}

func New(config *config.StorageConfig) (*Service, error) {
	address := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Address,
		config.Port,
		config.User,
		config.Password,
		config.DbName,
		config.EnableSSL,
	)

	dbStore := storage.New(config)
	if err := dbStore.Connect(context.Background()); err != nil {
		log.Println("failed to connect to db: ", err.Error())
		return nil, err
	}

	return &Service{
		dbAddress:  address,
		llmAddress: config.AddressLLM,
		db:         dbStore,
	}, nil
}

type SummaryWrapper struct {
	Content string `json:"content"`
}

type SummaryResponse struct {
	Summary string `json:"summary"`
	Class   string `json:"thematic"`
}

func (s *Service) LoadSummary(document *models.Document) {
	summaryRequest := models.NewLLM(document.Content)

	jsonData, err := json.Marshal(summaryRequest)
	if err != nil {
		log.Println("Failed while marshaling doc: ", err)
		return
	}

	reqBody := bytes.NewBuffer(jsonData)

	targetURL := fmt.Sprintf("%s/completion", s.llmAddress)
	mimeType := echo.MIMEApplicationJSON

	respData, recErr := sender.POST(reqBody, targetURL, mimeType, 300*time.Second)
	if recErr != nil {
		log.Println("failed send request: ", recErr)
		return
	}

	var summaryWrapper *SummaryWrapper
	if err := json.Unmarshal(respData, &summaryWrapper); err != nil {
		log.Println("Failed while reading response reqBody: ", err)
		return
	}

	var summaryResponse *SummaryResponse
	summaryWrapper.Content = strings.ReplaceAll(summaryWrapper.Content, "\t", "")
	summaryWrapper.Content = strings.ReplaceAll(summaryWrapper.Content, "\n", "")
	summaryWrapper.Content = strings.ReplaceAll(summaryWrapper.Content, "`", "")

	if err := json.Unmarshal([]byte(summaryWrapper.Content), &summaryResponse); err != nil {
		log.Println("Failed while reading response reqBody: ", err)
		return
	}

	document.Content = summaryResponse.Summary
	document.SetDocumentClass(summaryResponse.Class)

	ctx := context.Background()
	if err = s.StoreSummary(ctx, document); err != nil {
		log.Println("Failed while storing metadata to psql: ", err)
	}
}

func (s *Service) StoreSummary(ctx context.Context, document *models.Document) error {
	_, err := s.db.Create(ctx, document)
	return err
}
