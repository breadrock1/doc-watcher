package endpoints

import (
	"doc-notifier/internal/pkg/reader"
	watcher2 "doc-notifier/internal/pkg/watcher"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"mime/multipart"
	"os"
	"path"
)

// AnalyseFilesForm example
type AnalyseFilesForm struct {
	DocumentIDs []string `json:"document_ids" example:"886f7e11874040ca8b8461fb4cd1aa2c"`
}

// UnrecognizedDocuments example
type UnrecognizedDocuments struct {
	Unrecognized []*reader.DocumentPreview `json:"unrecognized"`
}

// MoveFilesForm example
type MoveFilesForm struct {
	TargetDirectory string   `json:"location" example:"common_folder"`
	SourceDirectory string   `json:"src_folder_id" example:"unrecognized"`
	DocumentPaths   []string `json:"document_ids" example:"./indexer/upload/test.txt"`
}

// RemoveFilesForm example
type RemoveFilesForm struct {
	DocumentPaths []string `json:"document_paths" example:"./indexer/upload/test.txt"`
}

// RemoveFilesError example
type RemoveFilesError struct {
	Code      int      `json:"code" example:"403"`
	Message   string   `json:"message" example:"File not found"`
	FilePaths []string `json:"file_paths" example:"./indexer/upload/test.txt"`
}

// UploadFilesToUnrecognized
// @Summary Upload files to analyse
// @Description Upload files to analyse
// @ID files-upload
// @Tags files
// @Accept  multipart/form
// @Produce  json
// @Param files formData file true "Files multipart form"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/files/upload [post]
func UploadFilesToUnrecognized(c echo.Context) error {
	var uploadErr error
	var multipartForm *multipart.Form

	if multipartForm, uploadErr = c.MultipartForm(); uploadErr != nil {
		respErr := createStatusResponse(400, uploadErr.Error())
		return c.JSON(400, respErr)
	}

	uploadFiles := make([]*reader.DocumentPreview, 0)
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	for _, fileForm := range multipartForm.File["files"] {
		filePath := path.Join("./indexer/unrecognized/", fileForm.Filename)
		if uploadErr = writeMultipart(fileForm, filePath); uploadErr != nil {
			log.Println(uploadErr)
			continue
		}

		document := watcher.Reader.ParseCaughtFiles(filePath)[0]
		prevDoc := reader.FromDocument(document)
		uploadFiles = append(uploadFiles, prevDoc)
		watcher.Reader.AddAwaitDocument(document)
	}

	return c.JSON(200, uploadFiles)
}

// AnalyseFiles
// @Summary Analyse uploaded files by ids
// @Description Analyse uploaded files by ids
// @ID files-analyse
// @Tags files
// @Accept  json
// @Produce  json
// @Param jsonQuery body AnalyseFilesForm true "Document ids to analyse"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/files/analyse [post]
func AnalyseFiles(c echo.Context) error {
	jsonForm := &AnalyseFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	processedDocuments := make([]*reader.DocumentPreview, 0)
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	for _, documentID := range jsonForm.DocumentIDs {
		// 1. Check does document id stored into recognized documents?!
		// 2.1. YES?! Stored this document to processedDocuments array;
		// 2.2. NO?!  Get document by id from unrecognized documents and send to channel;
		// 3. Return processedDocuments as Response.

		if watcher.IsRecognizedDocument(documentID) {
			recognizedDoc := watcher.PopRecognizedDocument(documentID)
			recPrevDoc := reader.FromDocument(recognizedDoc)
			processedDocuments = append(processedDocuments, recPrevDoc)
			continue
		}

		if watcher.Reader.IsUnrecognizedDocument(documentID) {
			unrecognizedDoc := watcher.Reader.PopUnrecognizedDocument(documentID)
			watcher.AppendCh <- unrecognizedDoc
			continue
		}
	}

	if len(processedDocuments) > 0 {
		return c.JSON(200, processedDocuments)
	}

	return c.JSON(102, processedDocuments)
}

// GetUnrecognized
// @Summary Get unrecognized documents
// @Description Get unrecognized documents
// @ID files-unrecognized
// @Tags files
// @Produce  json
// @Success 200 {object} UnrecognizedDocuments "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/files/unrecognized [get]
func GetUnrecognized(c echo.Context) error {
	collectedPreviews := make([]*reader.DocumentPreview, 0)
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	for _, document := range watcher.Reader.GetAwaitDocuments() {
		collectedPreviews = append(collectedPreviews, reader.FromDocument(document))
	}

	return c.JSON(200, collectedPreviews)
}

// MoveFiles
// @Summary Moving files to target directory
// @Description Moving files to target directory
// @ID moving
// @Tags files
// @Accept  json
// @Produce  json
// @Param jsonQuery body MoveFilesForm true "Document ids to move"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/files/move [post]
func MoveFiles(c echo.Context) error {
	jsonForm := &MoveFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	var collectedErrors []string
	targetDir := jsonForm.TargetDirectory
	sourceDir := jsonForm.SourceDirectory

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	for _, documentName := range jsonForm.DocumentPaths {
		srcFilePath := path.Join("./indexer", sourceDir, documentName)
		targetDirPath := path.Join("./indexer", targetDir)

		err := watcher.Reader.MoveFileToDir(srcFilePath, targetDirPath)
		if err != nil {
			collectedErrors = append(collectedErrors, documentName)
			log.Println(err)
			continue
		}
	}

	if len(collectedErrors) > 0 {
		msg := "These files hasn't been moved"
		respErr := RemoveFilesError{Code: 400, Message: msg, FilePaths: collectedErrors}
		return c.JSON(400, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// RemoveFiles
// @Summary Remove files from directory
// @Description Remove files from directory
// @ID files-remove
// @Tags files
// @Accept  json
// @Produce  json
// @Param jsonQuery body RemoveFilesForm true "Document paths to remove"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} RemoveFilesError "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/files/remove [post]
func RemoveFiles(c echo.Context) error {
	jsonForm := &RemoveFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	var collectedErrors []string
	for _, documentPath := range jsonForm.DocumentPaths {
		if err := os.RemoveAll(documentPath); err != nil {
			collectedErrors = append(collectedErrors, documentPath)
		}
	}

	if len(collectedErrors) > 0 {
		msg := "These files hasn't been removed"
		resErr := RemoveFilesError{Code: 400, Message: msg, FilePaths: collectedErrors}
		return c.JSON(400, resErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// DownloadFile
// @Summary Download file by path
// @Description Download file by path
// @ID files-download
// @Tags files
// @Produce  multipart/form
// @Param file_path formData string true "Path to file"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/files/download [post]
func DownloadFile(c echo.Context) error {
	filePath := c.FormValue("file_path")
	return c.File(filePath)
}
