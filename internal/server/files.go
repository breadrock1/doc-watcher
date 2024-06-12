package server

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"os"
	"path"

	"doc-notifier/internal/reader"
	"github.com/labstack/echo/v4"
)

func (s *Service) CreateFilesGroup() error {
	group := s.server.Group("/watcher/files")

	//timeoutMW := middleware.TimeoutWithConfig(middleware.TimeoutConfig{
	//	Timeout: s.watcher.Searcher.Timeout,
	//})

	group.POST("/upload", s.UploadFilesToUnrecognized)
	//group.POST("/analyse", s.AnalyseFiles, timeoutMW)
	group.POST("/analyse", s.AnalyseFiles)
	group.POST("/download", s.DownloadFile)
	group.POST("/move", s.MoveFiles)
	group.POST("/remove", s.RemoveFiles)
	group.GET("/unrecognized", s.GetUnrecognized)

	return nil
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
func (s *Service) UploadFilesToUnrecognized(c echo.Context) error {
	var uploadErr error
	var multipartForm *multipart.Form

	if multipartForm, uploadErr = c.MultipartForm(); uploadErr != nil {
		respErr := createStatusResponse(400, uploadErr.Error())
		return c.JSON(400, respErr)
	}

	uploadFiles := make([]*reader.Document, 0)
	for _, fileForm := range multipartForm.File["files"] {
		filePath := path.Join("./indexer/uploads/", fileForm.Filename)
		if uploadErr = s.writeMultipart(fileForm, filePath); uploadErr != nil {
			log.Println(uploadErr)
			continue
		}

		document := s.watcher.Reader.ParseCaughtFiles(filePath)[0]
		uploadFiles = append(uploadFiles, document)
		s.watcher.Reader.AddAwaitDocument(document)
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
func (s *Service) AnalyseFiles(c echo.Context) error {
	jsonForm := &AnalyseFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	processedDocuments := make([]*reader.Document, 0)
	for _, documentID := range jsonForm.DocumentIDs {
		// 1. Check does document id stored into recognized documents?!
		// 2.1. YES?! Stored this document to processedDocuments array;
		// 2.2. NO?!  Get document by id from unrecognized documents and send to channel;
		// 3. Return processedDocuments as Response.

		if s.watcher.IsRecognizedDocument(documentID) {
			recognizedDoc := s.watcher.PopRecognizedDocument(documentID)
			processedDocuments = append(processedDocuments, recognizedDoc)
			continue
		}

		if s.watcher.Reader.IsUnrecognizedDocument(documentID) {
			unrecognizedDoc := s.watcher.Reader.PopUnrecognizedDocument(documentID)
			s.watcher.AppendCh <- unrecognizedDoc
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
func (s *Service) GetUnrecognized(c echo.Context) error {
	return c.JSON(200, s.watcher.Reader.GetAwaitDocuments())
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
func (s *Service) MoveFiles(c echo.Context) error {
	jsonForm := &MoveFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	var collectedErrors []string
	targetDir := jsonForm.TargetDirectory
	sourceDir := jsonForm.SourceDirectory

	for _, documentName := range jsonForm.DocumentPaths {
		srcFilePath := path.Join("./indexer", sourceDir, documentName)
		targetDirPath := path.Join("./indexer", targetDir)

		err := s.watcher.Reader.MoveFileToDir(srcFilePath, targetDirPath)
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
func (s *Service) RemoveFiles(c echo.Context) error {
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
func (s *Service) DownloadFile(c echo.Context) error {
	filePath := c.FormValue("file_path")
	return c.File(filePath)
}
