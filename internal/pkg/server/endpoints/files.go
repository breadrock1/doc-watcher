package endpoints

import (
	"doc-notifier/internal/pkg/reader"
	watcher2 "doc-notifier/internal/pkg/watcher"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

const UploadHTMLForm = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Single file upload</title>
</head>
<body>
	<h1>Upload single file</h1>
	<form action="/files/upload" method="post" enctype="multipart/form-data">
        <input type="file" name="files" multiple>
        <input type="submit" value="Upload">
    </form>
</form>
</body>
</html>
`

// UploadFileForm
// @Summary Get upload file form
// @Description Get upload file form
// @ID upload-form
// @Tags files
// @Produce  html
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /files/upload [get]
func UploadFileForm(c echo.Context) error {
	return c.HTML(http.StatusOK, UploadHTMLForm)
}

// UploadFiles
// @Summary Upload file to server
// @Description Upload file to server
// @ID upload
// @Tags files
// @Accept  multipart/form
// @Produce  json
// @Param files formData file true "File entity"
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /files/upload [post]
func UploadFiles(c echo.Context) error {
	var uploadErr error
	var multipartForm *multipart.Form
	if multipartForm, uploadErr = c.MultipartForm(); uploadErr != nil {
		return uploadErr
	}

	var dstStream *os.File
	var srcStream multipart.File
	var uploadFiles []*reader.DocumentPreview

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)

	for _, fileForm := range multipartForm.File["files"] {
		if srcStream, uploadErr = fileForm.Open(); uploadErr != nil {
			_ = srcStream.Close()
			return uploadErr
		}

		filePath := path.Join("./indexer/unrecognized/", fileForm.Filename)
		if dstStream, uploadErr = os.Create(filePath); uploadErr != nil {
			_ = srcStream.Close()
			_ = dstStream.Close()
			return uploadErr
		}

		if _, uploadErr = io.Copy(dstStream, srcStream); uploadErr != nil {
			_ = srcStream.Close()
			_ = dstStream.Close()
			return uploadErr
		}

		_ = srcStream.Close()
		_ = dstStream.Close()

		docs := watcher.Reader.ParseCaughtFiles(filePath)
		prevDoc := reader.From(docs[0])
		uploadFiles = append(uploadFiles, prevDoc)
		watcher.Reader.AddAwaitDocument(docs[0])
	}

	return c.JSON(200, uploadFiles)
}

type AnalyseFilesForm struct {
	DocumentIDs []string `json:"document_ids"`
}

// AnalyseFiles
// @Summary Analyse uploaded files by ids
// @Description Analyse uploaded files by ids
// @ID analyse
// @Tags files
// @Accept  json
// @Produce  json
// @Param jsonQuery body AnalyseFilesForm true "Document ids to analyse"
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /files/analyse [post]
func AnalyseFiles(c echo.Context) error {
	jsonForm := &AnalyseFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	processedDocuments := make([]*reader.DocumentPreview, 0)
	for _, documentID := range jsonForm.DocumentIDs {
		// 1. Check does document id stored into recognized documents?!
		// 2.1. YES?! Stored this document to processedDocuments array
		// 2.2. NO?!  Get document by id from unrecognized documents and send to channel
		// 3. Return processedDocuments as Response

		if watcher.IsRecognizedDocument(documentID) {
			recognizedDoc := watcher.PopRecognizedDocument(documentID)
			recPrevDoc := reader.From(recognizedDoc)
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

type UnrecognizedDocuments struct {
	Unrecognized []*reader.DocumentPreview `json:"unrecognized"`
}

// GetUnrecognized
// @Summary Get unrecognized documents
// @Description Get unrecognized documents
// @ID unrecognized
// @Tags files
// @Produce  json
// @Success 200 {object} UnrecognizedDocuments "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /files/unrecognized [get]
func GetUnrecognized(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	var collectedPreviews []*reader.DocumentPreview
	for _, document := range watcher.Reader.GetAwaitDocuments() {
		collectedPreviews = append(collectedPreviews, reader.From(document))
	}

	return c.JSON(200, collectedPreviews)
}

type MoveFilesForm struct {
	TargetDirectory string   `json:"target_directory"`
	DocumentPaths   []string `json:"document_paths"`
}

type RemoveFilesForm struct {
	DocumentPaths []string `json:"document_paths"`
}

type RemoveFilesError struct {
	Code      int      `json:"code" example:"403"`
	Message   string   `json:"message" example:"File not found"`
	FilePaths []string `json:"file_paths" example:"[]"`
}

// MoveFiles
// @Summary Moving files to target directory
// @Description Moving files to target directory
// @ID moving
// @Tags files
// @Accept  json
// @Produce  json
// @Param jsonQuery body MoveFilesForm true "Document ids to move"
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /files/move [post]
func MoveFiles(c echo.Context) error {
	jsonForm := &MoveFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	var collectedErrors []string
	targetDir := jsonForm.TargetDirectory
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	for _, documentPath := range jsonForm.DocumentPaths {
		if err := watcher.Reader.MoveFileTo(documentPath, targetDir); err != nil {
			log.Println(err)
			collectedErrors = append(collectedErrors, documentPath)
			continue
		}
	}

	if len(collectedErrors) > 0 {
		return c.JSON(206, RemoveFilesError{
			Code:      206,
			Message:   "These files hasn't been moved",
			FilePaths: collectedErrors,
		})
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

// RemoveFiles
// @Summary Remove files from directory
// @Description Remove files from directory
// @ID remove
// @Tags files
// @Accept  json
// @Produce  json
// @Param jsonQuery body RemoveFilesForm true "Document paths to remove"
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} RemoveFilesError "Bad Request message"
// @Router /files/remove [post]
func RemoveFiles(c echo.Context) error {
	jsonForm := &RemoveFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	var collectedErrors []string
	for _, documentPath := range jsonForm.DocumentPaths {
		if err := os.RemoveAll(documentPath); err != nil {
			collectedErrors = append(collectedErrors, documentPath)
		}
	}

	if len(collectedErrors) > 0 {
		return c.JSON(206, RemoveFilesError{
			Code:      206,
			Message:   "These files hasn't been removed",
			FilePaths: collectedErrors,
		})
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

// DownloadFile
// @Summary Download file from server
// @Description Download file from server
// @ID download
// @Tags files
// @Produce  multipart/form
// @Param file_path formData string true "Path to file on server"
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /files/download [post]
func DownloadFile(c echo.Context) error {
	filePath := c.FormValue("file_path")
	return c.File(filePath)
}
