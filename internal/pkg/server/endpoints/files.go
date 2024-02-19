package endpoints

import "github.com/labstack/echo/v4"

func DownloadFile(c echo.Context) error {
	filePath := c.FormValue("file_path")
	return c.File(filePath)
}
