package server

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Service) CreateWatcherGroup() error {
	group := s.server.Group("/watcher")

	group.GET("/run", s.RunWatchers)
	group.GET("/pause", s.PauseWatchers)

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
