package endpoints

import "github.com/labstack/echo/v4"

// Hello
// @Summary Hello
// @Tags hello
// @Description Check service is available
// @ID hello
// @Accept  json
// @Produce  json
// @Success 200 {object} ResponseForm "Done"
// @Failure	503 {object} ServerErrorForm "Server does not available""
// @Router /hello/ [get]
func Hello(c echo.Context) error {
	responseForm := &ResponseForm{
		Status:  200,
		Message: "Done",
	}

	return c.JSON(200, responseForm)
}
