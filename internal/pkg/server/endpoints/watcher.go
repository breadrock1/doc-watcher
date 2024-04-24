package endpoints

import (
	watcher2 "doc-notifier/internal/pkg/watcher"
	"encoding/json"
	"github.com/labstack/echo/v4"
)

type CreateWatchDirectoryForm struct {
	DirectoryName string `json:"directory_name"`
}

func GetWatchedDirectories(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	watcherDirs := watcher.WatchedDirsList()
	return c.JSON(200, watcherDirs)
}

func CreateWatchDirectory(c echo.Context) error {
	jsonForm := &CreateWatchDirectoryForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if err := watcher.AppendDirectory(jsonForm.DirectoryName); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	if err := watcher.Searcher.CreateBucket(jsonForm.DirectoryName); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

func RemoveWatchDirectory(c echo.Context) error {
	jsonForm := &CreateWatchDirectoryForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if err := watcher.RemoveDirectory(jsonForm.DirectoryName); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	if err := watcher.Searcher.DeleteBucket(jsonForm.DirectoryName); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

func returnStatusResponse(status int, msg string) *ResponseForm {
	return &ResponseForm{
		Status:  status,
		Message: msg,
	}
}

func PauseWatchers(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if watcher.PauseWatchers {
		return c.JSON(200, returnStatusResponse(400, "Already paused"))
	} else {
		watcher.PauseWatchers = true
		return c.JSON(200, returnStatusResponse(200, "Done"))
	}
}

func RunWatchers(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	watcher.PauseWatchers = false
	return c.JSON(200, returnStatusResponse(200, "Done"))
}

func GetUnrecognizedFiles(c echo.Context) error {
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	files := watcher.UnrecognizedFiles()
	return c.JSON(200, files)
}
