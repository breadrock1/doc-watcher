package reader

import (
	"crypto/md5"
	"fmt"
	"github.com/fatih/structs"
	"github.com/glaslos/ssdeep"
	"github.com/google/uuid"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	bucketName    = "common_folder"
	timeFormat    = time.RFC3339
	documentMimes = []string{
		"csv", "msword", "html", "json", "pdf",
		"rtf", "plain", "vnd.ms-excel", "xml",
		"vnd.ms-powerpoint", "vnd.oasis.opendocument.text",
		"vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"vnd.openxmlformats-officedocument.wordprocessingml.document",
		"vnd.openxmlformats-officedocument.presentationml.presentation",
	}
)

type Document struct {
	FolderID            string       `json:"folder_id"`
	FolderPath          string       `json:"folder_path"`
	ContentUUID         string       `json:"content_uuid"`
	ContentMD5          string       `json:"content_md5"`
	Content             string       `json:"content"`
	ContentVector       []float64    `json:"content_vector"`
	DocumentMD5         string       `json:"document_md5"`
	DocumentSSDEEP      string       `json:"document_ssdeep"`
	DocumentName        string       `json:"document_name"`
	DocumentPath        string       `json:"document_path"`
	DocumentSize        int64        `json:"document_size"`
	DocumentType        string       `json:"document_type"`
	DocumentExtension   string       `json:"document_extension"`
	DocumentPermissions int32        `json:"document_permissions"`
	DocumentCreated     string       `json:"document_created"`
	DocumentModified    string       `json:"document_modified"`
	QualityRecognized   int32        `json:"quality_recognition"`
	OcrMetadata         *OcrMetadata `json:"ocr_metadata"`
}

type OcrMetadata struct {
	JobId      string     `json:"job_id"`
	Text       string     `json:"text"`
	PagesCount int        `json:"pages_count"`
	DocType    string     `json:"doc_type"`
	Artifacts  *Artifacts `json:"artifacts"`
}

type Artifacts struct {
	TransportInvoiceDate      string `json:"date_of_transport_invoice" name:"Дата транспортной накладной"`
	TransportInvoiceNumber    string `json:"number_of_transport_invoice" name:"Номер транспортной накладной"`
	TransferCompany           string `json:"transfer_company" name:"Трансферная Компания"`
	OrderNumber               string `json:"order_number" name:"Номер заказа"`
	Carrier                   string `json:"carrier" name:"Перевозчик"`
	VehicleNumber             string `json:"vehicle_number" name:"Номер автомобиля"`
	CargoDateArrival          string `json:"arrival_of_cargo_date_time" name:"Дата прибытия груза"`
	CargoDateDeparture        string `json:"departure_of_cargo_date_time" name:"Дата отправления груза"`
	AddressRedirection        string `json:"redirection_address" name:"Адрес перенаправление"`
	DateRedirection           string `json:"redirection_date_time" name:"Дата перенаправления"`
	CargoIssueAddress         string `json:"cargo_issue_address" name:"Адрес выдачи груза"`
	CargoIssueDate            string `json:"cargo_issue_date" name:"Дата выдачи груза"`
	CargoWeight               string `json:"cargo_weight" name:"Вес груза"`
	CargoPlacesNumber         string `json:"number_of_cargo_places" name:"Номер места для автомобиля"`
	ContainerReceiptActNumber string `json:"container_receipt_act_number" name:"Номер акта получения контейнера"`
	ContainerReceiptActDate   string `json:"container_receipt_act_date_time" name:"Дата акта приема контейнера"`
	ContainerNumber           string `json:"container_number" name:"Номер контейнера"`
	TerminalName              string `json:"terminal_name" name:"Имя терминала"`
	KtkName                   string `json:"ktk_state" name:"Имя ктк"`
	DriverFullName            string `json:"driver_full_name" name:"Полное имя водителя"`
	DocumentNumber            string `json:"document_number" name:"Номер документа"`
	ShipName                  string `json:"ship_name" name:"Название корабля"`
	FlightNumber              string `json:"flight_number" name:"Номер рейса"`
	ShipDate                  string `json:"ship_date" name:"Дата отправки"`
	DocumentType              string `json:"document_type" name:"Тип документа"`
	Seals                     bool   `json:"seals" name:"Морские котики?"`
}

func ParseFile(filePath string) (*Document, error) {
	absFilePath, _ := filepath.Abs(filePath)
	fileInfo, err := os.Stat(absFilePath)
	if err != nil {
		err = fmt.Errorf("failed while getting stat of file: %e", err)
		return nil, err
	}

	modifiedTime := time.Now().UTC()
	createdTime := fileInfo.ModTime().UTC()
	fileExt := filepath.Ext(filePath)
	data, _ := os.ReadFile(absFilePath)

	document := &Document{}
	document.FolderID = bucketName
	document.FolderPath = parseBucketName(absFilePath)
	document.DocumentMD5 = fmt.Sprintf("%x", md5.Sum(data))
	document.DocumentPath = absFilePath
	document.DocumentName = fileInfo.Name()
	document.DocumentSize = fileInfo.Size()
	document.DocumentType = parseDocumentType(fileExt)
	document.DocumentExtension = fileExt
	document.DocumentPermissions = int32(fileInfo.Mode().Perm())
	document.ContentUUID = uuid.NewString()
	document.DocumentModified = modifiedTime.Format(timeFormat)
	document.DocumentCreated = createdTime.Format(timeFormat)
	document.QualityRecognized = -1

	return document, nil
}

func (d *Document) SetContentData(data string) {
	d.Content = data
}

func (d *Document) SetContentVector(data []float64) {
	d.ContentVector = data
}

func (d *Document) AppendContentVector(data []float64) {
	d.ContentVector = append(d.ContentVector, data...)
}

func (d *Document) ComputeMd5Hash() {
	data := []byte(d.Content)
	d.DocumentMD5 = fmt.Sprintf("%x", md5.Sum(data))
}

func (d *Document) ComputeContentMd5Hash() {
	data := []byte(d.Content)
	d.ContentMD5 = fmt.Sprintf("%x", md5.Sum(data))
}

func (d *Document) ComputeSsdeepHash() {
	data := []byte(d.Content)
	if hashData, err := ssdeep.FuzzyBytes(data); err == nil {
		d.DocumentSSDEEP = hashData
	}
}

func (d *Document) SetContentUUID(contentUUID string) {
	d.ContentUUID = contentUUID
}

func (d *Document) ComputeContentUUID() {
	data := []byte(d.Content)
	if uuidToken, err := uuid.FromBytes(data); err == nil {
		d.ContentUUID = uuidToken.String()
	}
}

func (d *Document) SetQuality(quality int32) {
	d.QualityRecognized = quality
}

func (d *Document) GetGroupedProperties() []*PreviewProperties {
	properties := make([]*PreviewProperties, 0)
	if d.OcrMetadata.Artifacts == nil {
		return properties
	}

	for _, field := range structs.Fields(d.OcrMetadata.Artifacts) {
		if field.Tag("json") == "seals" {
			continue
		}

		fieldData := field.Value()
		if fieldData == nil {
			continue
		}

		value := fieldData.(string)
		if len(value) == 0 {
			continue
		}

		key := field.Tag("json")
		name := field.Tag("name")
		properties = append(properties, &PreviewProperties{
			Key:   key,
			Name:  name,
			Value: value,
		})
	}

	return properties
}

func parseBucketName(filePath string) string {
	currPath := os.Getenv("PWD")
	relPath, err := filepath.Rel(currPath, filePath)
	relPath2, err := filepath.Rel("indexer", relPath)
	if err != nil {
		log.Printf("Failed while parsing bucket name")
		return bucketName
	}

	bucketNameRes, _ := filepath.Split(relPath2)
	bucketNameRes2 := strings.ReplaceAll(bucketNameRes, "/", "")
	if bucketNameRes2 == "" {
		return bucketName
	}

	return bucketNameRes2
}

func parseDocumentType(extension string) string {
	mimeType := mime.TypeByExtension(extension)
	attributes := strings.Split(mimeType, "/")
	switch attributes[0] {
	case "audio":
		return "audio"
	case "image":
		return "image"
	case "video":
		return "video"
	case "text":
		return "document"
	case "application":
		return extractApplicationMimeType(attributes[1])
	default:
		return "unknown"
	}
}

func extractApplicationMimeType(attribute string) string {
	for _, mimeType := range documentMimes {
		if mimeType == attribute {
			return "document"
		}
	}

	return "unknown"
}
