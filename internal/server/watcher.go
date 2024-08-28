package server

import (
	"encoding/json"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Service) CreateWatcherGroup() error {
	group := s.server.Group("/watcher")

	group.GET("/run", s.RunWatchers)
	group.GET("/pause", s.PauseWatchers)
	group.POST("/attach", s.AttachDirectories)
	group.POST("/detach", s.DetachDirectories)

	return nil
}

// IsWatcherPaused
// @Summary Check does watcher has been paused
// @Description Check does watcher has been paused
// @ID is-watcher-pause
// @Tags watcher
// @Produce  json
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/is-pause [get]
func (s *Service) IsWatcherPaused(c echo.Context) error {
	isPaused := s.watcher.Watcher.IsPausedWatchers()
	isPausedStr := strconv.FormatBool(isPaused)
	return c.JSON(200, createStatusResponse(200, isPausedStr))
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
	if s.watcher.Watcher.IsPausedWatchers() {
		okResp := createStatusResponse(204, "Already paused")
		return c.JSON(204, okResp)
	}

	s.watcher.Watcher.PauseWatchers(true)

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
	s.watcher.Watcher.PauseWatchers(false)
	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// AttachDirectories
// @Summary Attach new directory to watcher
// @Description Attach new directory to watcher
// @ID folders-attach
// @Tags watcher
// @Accept  json
// @Produce json
// @Param jsonQuery body WatcherDirectoriesForm true "File entity"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/attach [post]
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
// @Produce json
// @Param jsonQuery body WatcherDirectoriesForm true "Folder ids"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/detach [post]
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
