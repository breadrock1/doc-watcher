package server

import (
	"github.com/labstack/echo/v4"
)

func (s *Service) CreateFilesGroup() error {
	group := s.server.Group("/watcher/files")

	group.POST("/upload-from-office", s.UploadFileFromOnlyOffice)
	group.POST("/download", s.DownloadFile)

	return nil
}

// UploadFileFromOnlyOffice
// @Summary Upload files to analyse from onlyoffice
// @Description Upload files to analyse from onlyoffice
// @ID files-upload-office
// @Tags files
// @Produce  json
// @Param fileName query string true "File name to download from office"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /watcher/files/upload-from-office [post]
func (s *Service) UploadFileFromOnlyOffice(c echo.Context) error {
	fileName := c.QueryParam("fileName")
	if err := s.office.DownloadDocument(fileName); err != nil {
		return c.JSON(400, createStatusResponse(400, err.Error()))
	}
	return c.JSON(200, createStatusResponse(200, "Done"))
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
