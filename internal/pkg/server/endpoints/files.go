package endpoints

import (
	"doc-notifier/internal/pkg/reader"
	watcher2 "doc-notifier/internal/pkg/watcher"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
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

	var processedDocuments []*reader.DocumentPreview
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	for _, documentID := range jsonForm.DocumentIDs {

		document := watcher.Reader.GetAwaitDocument(documentID)
		if document == nil {
			continue
		}

		oldLocation := document.DocumentPath
		prevDocument := reader.From(document)
		if err := watcher.ProcessTriggeredFile(document); err == nil {
			prevDocument.SetQuality(reader.MaxQualityValue)
		} else {
			prevDocument.SetQuality(1)
		}

		if prevDocument.QualityOCR > 1 {
			_ = watcher.Reader.MoveFileTo(oldLocation, prevDocument.Location)
		}

		watcher.Reader.PopDoneDocument(prevDocument.DocumentID)
		processedDocuments = append(processedDocuments, prevDocument)
	}

	return c.JSON(200, processedDocuments)
}

type MoveFileForm struct {
	TargetDirectory string   `json:"target_directory"`
	DocumentPaths   []string `json:"document_paths"`
}

// MoveFile
// @Summary Moving files to target directory
// @Description Moving files to target directory
// @ID moving
// @Tags files
// @Accept  json
// @Produce  json
// @Param jsonQuery body MoveFileForm true "Document ids to move"
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /files/move [post]
func MoveFile(c echo.Context) error {
	jsonForm := &MoveFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	var collectedErrors []string
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	for _, documentPath := range jsonForm.DocumentPaths {
		if err := watcher.Reader.MoveFileTo(documentPath, jsonForm.TargetDirectory); err != nil {
			log.Println(err)
			collectedErrors = append(collectedErrors, documentPath)
			continue
		}
	}

	if len(collectedErrors) > 0 {
		files := strings.Join(collectedErrors, ", ")
		msg := fmt.Sprintf("Failed while moving files: %s", files)
		return c.JSON(403, returnStatusResponse(403, msg))
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

type RemoveFileForm struct {
	FilePath string `json:"file_path"`
}

func RemoveFile(c echo.Context) error {
	jsonForm := &RemoveFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	documents := watcher.Reader.ParseCaughtFiles(jsonForm.FilePath)
	res, err := os.ReadFile(jsonForm.FilePath)
	if err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}
	watcher.Reader.ComputeMd5HashByData(documents[0], res)

	targetURL := fmt.Sprintf("%s/document/%s/%s", watcher.Searcher.Address, documents[0].BucketUUID, documents[0].DocumentMD5)
	response, err := http.Get(targetURL)
	if err != nil || response.StatusCode != 200 {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	if err = os.RemoveAll(jsonForm.FilePath); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

type ParsingJobForm struct {
	JobId string `json:"job_id"`
}

func GetUploadingStatus(c echo.Context) error {
	jsonForm := &ParsingJobForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	jobRes := watcher.Ocr.Ocr.GetProcessingJob(jsonForm.JobId)
	return c.JSON(200, jobRes)
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

	return c.JSON(200, UnrecognizedDocuments{Unrecognized: collectedPreviews})
}
