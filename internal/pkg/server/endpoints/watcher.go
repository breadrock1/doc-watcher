package endpoints

import (
	watcher2 "doc-notifier/internal/pkg/watcher"
	"encoding/json"
	"github.com/labstack/echo/v4"
)

// WatcherDirectoriesForm example
type WatcherDirectoriesForm struct {
	Paths []string `json:"paths" example:"./indexer"`
}

func returnStatusResponse(status int, msg string) *ResponseForm {
	return &ResponseForm{
		Status:  status,
		Message: msg,
	}
}

// WatchedDirsList
// @Summary Get watcher dirs list
// @Description Get watcher dirs list
// @ID watched-dirs
// @Tags watcher
// @Accept  json
// @Produce  json
// @Success 200 {array} string "Watched dirs"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /watcher/all [get]
func WatchedDirsList(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	watcherDirs := watcher.GetWatchedDirectories()
	return c.JSON(200, watcherDirs)
}

// AttachDirectories
// @Summary Attach new directory to watcher
// @Description Attach new directory to watcher
// @ID attach
// @Tags watcher
// @Accept  json
// @Produce  json
// @Param jsonQuery body WatcherDirectoriesForm true "File entity"
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /watcher/attach [post]
func AttachDirectories(c echo.Context) error {
	jsonForm := &WatcherDirectoriesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if err := watcher.AppendDirectories(jsonForm.Paths); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

// DetachDirectories
// @Summary Attach new directory to watcher
// @Description Attach new directory to watcher
// @ID attach
// @Tags watcher
// @Accept  json
// @Produce  json
// @Param jsonQuery body WatcherDirectoriesForm true "File entity"
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /watcher/attach [post]
func DetachDirectories(c echo.Context) error {
	jsonForm := &WatcherDirectoriesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if err := watcher.RemoveDirectories(jsonForm.Paths); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

// PauseWatchers
// @Summary Pause all watchers
// @Description Pause all watchers
// @ID pause
// @Tags watcher
// @Produce  json
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /watcher/pause [get]
func PauseWatchers(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if watcher.PauseWatchers {
		return c.JSON(200, returnStatusResponse(400, "Already paused"))
	}

	watcher.PauseWatchers = true
	return c.JSON(200, returnStatusResponse(200, "Done"))
}

// RunWatchers
// @Summary Run all watchers
// @Description Run all watchers
// @ID run
// @Tags watcher
// @Produce  json
// @Success 200 {object} ResponseForm "Done"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Router /watcher/run [get]
func RunWatchers(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	watcher.PauseWatchers = false
	return c.JSON(200, returnStatusResponse(200, "Done"))
}
