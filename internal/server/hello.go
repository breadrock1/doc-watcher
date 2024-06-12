package server

import (
	"github.com/labstack/echo/v4"
)

func (s *Service) CreateHelloGroup() error {
	group := s.server.Group("/hello")
	group.GET("/", s.Hello)
	return nil
}

// Hello
// @Summary Hello
// @Tags hello
// @Description Check service is available
// @ID hello
// @Produce  json
// @Success 200 {object} ResponseForm "Ok"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /hello/ [get]
func (s *Service) Hello(c echo.Context) error {
	okResp := createStatusResponse(200, "Ok")
	return c.JSON(200, okResp)
}
