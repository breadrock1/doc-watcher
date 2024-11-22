package httpserv

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
)

func (s *Service) CreateWatcherGroup() error {
	group := s.server.Group("/watcher")

	group.GET("/run", s.RunWatchers)
	group.GET("/stop", s.StopWatchers)
	group.PUT("/attach", s.AttachDirectory)
	group.DELETE("/:bucket/detach", s.DetachDirectory)
	group.POST("/processing/fetch", s.FetchProcessingDocuments)
	group.POST("/processing/clean", s.CleanProcessingDocuments)

	return nil
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
	ctx := c.Request().Context()
	s.watcher.Watcher.RunWatchers(ctx)
	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// StopWatchers
// @Summary Stop all watchers
// @Description Stop all watchers
// @ID watcher-stop
// @Tags watcher
// @Produce  json
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/stop [get]
func (s *Service) StopWatchers(c echo.Context) error {
	ctx := c.Request().Context()
	s.watcher.Watcher.TerminateWatchers(ctx)
	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// AttachDirectory
// @Summary Attach new directory to watcher
// @Description Attach new directory to watcher
// @ID folders-attach
// @Tags watcher
// @Accept  json
// @Produce json
// @Param jsonQuery body AttachDirectoryForm true "File entity"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/attach [put]
func (s *Service) AttachDirectory(c echo.Context) error {
	jsonForm := &AttachDirectoryForm{}
	decoder := json.NewDecoder(c.Request().Body)
	err := decoder.Decode(jsonForm)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	err = s.watcher.Watcher.AttachDirectory(ctx, jsonForm.BucketName)
	if err != nil {
		return err
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// DetachDirectory
// @Summary Attach new directory to watcher
// @Description Attach new directory to watcher
// @ID folders-detach
// @Tags watcher
// @Accept  json
// @Produce json
// @Param bucket path string true "Folder ids"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/{bucket}/detach [delete]
func (s *Service) DetachDirectory(c echo.Context) error {
	bucket := c.Param("bucket")

	ctx := c.Request().Context()
	if err := s.watcher.Watcher.DetachDirectory(ctx, bucket); err != nil {
		return err
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// FetchProcessingDocuments
// @Summary Fetch processing documents
// @Description Load processing/unrecognized/done documents by names list
// @ID fetch-documents
// @Tags watcher
// @Accept  json
// @Produce json
// @Param jsonQuery body FetchDocumentsList true "File names to fetch processing status"
// @Success 200 {object} watcher.ProcessingDocuments "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/processing/fetch [post]
func (s *Service) FetchProcessingDocuments(c echo.Context) error {
	jsonForm := &FetchDocumentsList{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return err
	}

	ctx := c.Request().Context()
	documents := s.watcher.Watcher.FetchProcessingDocuments(ctx, jsonForm.FileNames)
	return c.JSON(200, documents)
}

// CleanProcessingDocuments
// @Summary Clean processing documents
// @Description Clean processing documents
// @ID clean-documents
// @Tags watcher
// @Accept  json
// @Param jsonQuery body FetchDocumentsList true "File names to clean processing status"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/processing/clean [post]
func (s *Service) CleanProcessingDocuments(c echo.Context) error {
	jsonForm := &FetchDocumentsList{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := s.watcher.Watcher.CleanProcessingDocuments(ctx, jsonForm.FileNames); err != nil {
		return err
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}
