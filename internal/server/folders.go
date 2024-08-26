package server

import (
	"bytes"
	"encoding/json"
	"log"
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

func (s *Service) CreateFoldersGroup() error {
	group := s.server.Group("/watcher/folders")

	group.GET("/", s.GetWatchedDirectories)
	group.POST("/create", s.CreateFolder)
	group.POST("/delete", s.RemoveFolder)
	group.POST("/attach", s.AttachDirectories)
	group.POST("/detach", s.DetachDirectories)
	group.POST("/upload", s.UploadFilesToWatcher)
	group.POST("/download", s.LoadFileFromWatcher)
	group.POST("/update", s.UpdateWatchedFile)
	group.POST("/hierarchy", s.GetUserHierarchy)

	return nil
}

// GetWatchedDirectories
// @Summary Get watched directories list
// @Description Get watched directories list
// @ID folders-all
// @Tags watcher
// @Produce  json
// @Success 200 {array} string "Ok"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/ [get]
func (s *Service) GetWatchedDirectories(c echo.Context) error {
	watcherDirs := s.watcher.Watcher.GetWatchedDirectories()
	return c.JSON(200, watcherDirs)
}

// CreateFolder
// @Summary Create folder to store documents
// @Description Create folder to store documents
// @ID folder-create
// @Tags files
// @Produce  json
// @Param jsonQuery body FolderNameForm true "Folder name to create"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/create [post]
func (s *Service) CreateFolder(c echo.Context) error {
	jsonForm := &FolderNameForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	if err := s.watcher.Watcher.CreateDirectory(jsonForm.FolderName); err != nil {
		respErr := createStatusResponse(208, err.Error())
		return c.JSON(208, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// RemoveFolder
// @Summary Remove folder
// @Description Remove folder
// @ID folder-remove
// @Tags files
// @Produce  json
// @Param jsonQuery body FolderNameForm true "Folder name to remove"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/remove [post]
func (s *Service) RemoveFolder(c echo.Context) error {
	jsonForm := &FolderNameForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	if err := s.watcher.Watcher.RemoveDirectory(jsonForm.FolderName); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// AttachDirectories
// @Summary Attach new directory to watcher
// @Description Attach new directory to watcher
// @ID folders-attach
// @Tags watcher
// @Accept  json
// @Produce  json
// @Param jsonQuery body WatcherDirectoriesForm true "File entity"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/attach [post]
func (s *Service) AttachDirectories(c echo.Context) error {
	jsonForm := &WatcherDirectoriesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	if err := s.watcher.Watcher.AppendDirectories(jsonForm.Paths); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// DetachDirectories
// @Summary Attach new directory to watcher
// @Description Attach new directory to watcher
// @ID folders-detach
// @Tags watcher
// @Accept  json
// @Produce  json
// @Param jsonQuery body WatcherDirectoriesForm true "Folder ids"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/detach [post]
func (s *Service) DetachDirectories(c echo.Context) error {
	jsonForm := &WatcherDirectoriesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	if err := s.watcher.Watcher.RemoveDirectories(jsonForm.Paths); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// UploadFilesToWatcher
// @Summary Upload file to watcher directory
// @Description Upload file to watcher directory
// @ID watcher-upload
// @Tags watcher
// @Accept  multipart/form
// @Produce  json
// @Param files formData file true "Files multipart form"
// @Param bucket query string true "Bucket name to upload"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/upload [post]
func (s *Service) UploadFilesToWatcher(c echo.Context) error {
	var uploadErr error
	var fileData bytes.Buffer
	var multipartForm *multipart.Form

	if multipartForm, uploadErr = c.MultipartForm(); uploadErr != nil {
		respErr := createStatusResponse(400, uploadErr.Error())
		return c.JSON(400, respErr)
	}

	bucketName := c.QueryParam("bucket")
	for _, fileForm := range multipartForm.File["files"] {

		fileName := fileForm.Filename
		fileHandler, uploadErr := fileForm.Open()
		if uploadErr != nil {
			log.Println(uploadErr)
			continue
		}

		_, uploadErr = fileData.ReadFrom(fileHandler)
		if uploadErr != nil {
			log.Println(uploadErr)
			continue
		}

		uploadErr = s.watcher.Watcher.UploadFile(bucketName, fileName, fileData)
		if uploadErr != nil {
			log.Println(uploadErr)
			continue
		}
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// LoadFileFromWatcher
// @Summary Load file from watcher directory
// @Description Load file from watcher directory
// @ID watcher-download
// @Tags watcher
// @Accept  multipart/form
// @Produce  json
// @Param jsonQuery body DownloadFile true "Download file form"
// @Success 200 {object} models.Document "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/download [post]
func (s *Service) LoadFileFromWatcher(c echo.Context) error {
	jsonForm := &DownloadFile{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	fileData, err := s.watcher.Watcher.DownloadFile(jsonForm.Bucket, jsonForm.FileName)
	if err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	return c.Blob(200, echo.MIMEMultipartForm, fileData.Bytes())
}

// UpdateWatchedFile
// @Summary Update file into watcher directory
// @Description Update file into watcher directory
// @ID watcher-update
// @Tags watcher
// @Accept  multipart/form
// @Produce  json
// @Param files formData file true "Files (multipart/form) to updated"
// @Param bucket query string true "Bucket name to upload"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/update [post]
func (s *Service) UpdateWatchedFile(c echo.Context) error {
	return s.UploadFilesToWatcher(c)
}

// GetUserHierarchy
// @Summary Get bucket fs hierarchy
// @Description Get bucket fs hierarchy
// @ID watcher-hierarchy
// @Tags watcher
// @Accept  multipart/form
// @Produce  json
// @Param jsonQuery body HierarchyForm true "Hierarchy form"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/hierarchy [post]
func (s *Service) GetUserHierarchy(c echo.Context) error {
	jsonForm := &HierarchyForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	listObjects := s.watcher.Watcher.GetHierarchy(jsonForm.Bucket, jsonForm.DirectoryName)
	return c.JSON(200, listObjects)
}
