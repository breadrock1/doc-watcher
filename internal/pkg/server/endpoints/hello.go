package endpoints

import "github.com/labstack/echo/v4"

func Hello(c echo.Context) error {
	responseForm := &ResponseForm{
		Status:  200,
		Message: "Done",
	}

	return c.JSON(200, responseForm)
}
