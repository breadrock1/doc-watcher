package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"

	"github.com/labstack/echo/v4"
)

func (s *Service) CreateFoldersGroup() error {
	group := s.server.Group("/watcher/folders")

	group.GET("/", s.GetWatchedDirectories)
	group.POST("/attach", s.AttachDirectories)
	group.POST("/detach", s.DetachDirectories)
	group.POST("/upload", s.UploadFilesToWatcher)
	group.POST("/create", s.CreateFolder)
	group.POST("/delete", s.RemoveFolder)

	return nil
}

func (s *Service) CreateWatcherGroup() error {
	group := s.server.Group("/watcher")

	group.GET("/run", s.RunWatchers)
	group.GET("/pause", s.PauseWatchers)

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
	watcherDirs := s.watcher.GetWatchedDirectories()
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

	folderPath := path.Join("./indexer", jsonForm.FolderName)
	if err := os.Mkdir(folderPath, os.ModePerm); err != nil {
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

	folderPath := path.Join("./indexer", jsonForm.FolderName)
	if err := os.RemoveAll(folderPath); err != nil {
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

	if err := s.watcher.AppendDirectories(jsonForm.Paths); err != nil {
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

	if err := s.watcher.RemoveDirectories(jsonForm.Paths); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// PauseWatchers
// @Summary Pause all watchers
// @Description Pause all watchers
// @ID watcher-pause
// @Tags watcher
// @Produce  json
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/pause [get]
func (s *Service) PauseWatchers(c echo.Context) error {
	if s.watcher.PauseWatchers {
		okResp := createStatusResponse(204, "Already paused")
		return c.JSON(204, okResp)
	}

	s.watcher.PauseWatchers = true

	okResp := createStatusResponse(200, "Ok")
	return c.JSON(200, okResp)
}

// RunWatchers
// @Summary Run all watchers
// @Description Run all watchers
// @ID watcher-run
// @Tags watcher
// @Produce  json
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/run [get]
func (s *Service) RunWatchers(c echo.Context) error {
	s.watcher.PauseWatchers = false
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
// @Success 200 {array} []reader.Document "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/upload [post]
func (s *Service) UploadFilesToWatcher(c echo.Context) error {
	var uploadErr error
	var multipartForm *multipart.Form

	if multipartForm, uploadErr = c.MultipartForm(); uploadErr != nil {
		respErr := createStatusResponse(400, uploadErr.Error())
		return c.JSON(400, respErr)
	}

	for _, fileForm := range multipartForm.File["files"] {
		filePath := fmt.Sprintf("./indexer/watcher/%s", fileForm.Filename)
		if uploadErr = s.writeMultipart(fileForm, filePath); uploadErr != nil {
			log.Println(uploadErr)
			continue
		}
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

func (s *Service) writeMultipart(fileForm *multipart.FileHeader, filePath string) error {
	var respErr error

	srcStream, respErr := fileForm.Open()
	if respErr != nil {
		return fmt.Errorf("failed to open opens file header: %e", respErr)
	}
	defer func() { _ = srcStream.Close() }()

	dstStream, respErr := os.Create(filePath)
	if respErr != nil {
		return fmt.Errorf("failed to create file: %e", respErr)
	}
	defer func() { _ = dstStream.Close() }()

	if _, respErr = io.Copy(dstStream, srcStream); respErr != nil {
		return fmt.Errorf("failed to write data: %e", respErr)
	}

	return nil
}
