package endpoints

import (
	watcher2 "doc-notifier/internal/pkg/watcher"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
)

// WatcherDirectoriesForm example
type WatcherDirectoriesForm struct {
	Paths []string `json:"paths" example:"./indexer/test_folder"`
}

// FolderNameForm example
type FolderNameForm struct {
	FolderName string `json:"folder_name" example:"test_folder"`
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
func GetWatchedDirectories(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	watcherDirs := watcher.GetWatchedDirectories()
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
func CreateFolder(c echo.Context) error {
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
func RemoveFolder(c echo.Context) error {
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
func AttachDirectories(c echo.Context) error {
	jsonForm := &WatcherDirectoriesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if err := watcher.AppendDirectories(jsonForm.Paths); err != nil {
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
func DetachDirectories(c echo.Context) error {
	jsonForm := &WatcherDirectoriesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if err := watcher.RemoveDirectories(jsonForm.Paths); err != nil {
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
func PauseWatchers(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if watcher.PauseWatchers {
		okResp := createStatusResponse(204, "Already paused")
		return c.JSON(204, okResp)
	}

	watcher.PauseWatchers = true

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
func RunWatchers(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	watcher.PauseWatchers = false
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
// @Param directory formData string true "Directory to upload"
// @Success 200 {array} reader.DocumentPreview "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/folders/upload [post]
func UploadFilesToWatcher(c echo.Context) error {
	var uploadErr error
	var multipartForm *multipart.Form

	if multipartForm, uploadErr = c.MultipartForm(); uploadErr != nil {
		respErr := createStatusResponse(400, uploadErr.Error())
		return c.JSON(400, respErr)
	}

	targetDir := c.FormValue("directory")
	for _, fileForm := range multipartForm.File["files"] {
		filePath := fmt.Sprintf("./indexer/%s/%s", targetDir, fileForm.Filename)
		if uploadErr = writeMultipart(fileForm, filePath); uploadErr != nil {
			log.Println(uploadErr)
			continue
		}
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

func writeMultipart(fileForm *multipart.FileHeader, filePath string) error {
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
