package httpserv

import (
	"fmt"
	"strings"

	"doc-watcher/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitLogger(servConfig *server.Config) echo.MiddlewareFunc {
	config := middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			uri := c.Path()
			return strings.Contains(uri, "swagger")
		},

		Format: fmt.Sprintf(
			"%s  %s %s request{%s}: %s %s ms %s\n",
			"${time_rfc3339}",
			servConfig.LoggerLevel,
			"${id}",
			"method=${method} uri=${path}",
			"latency=${latency}",
			"status=${status}",
			"error=\"${error}\"",
		),
	}

	return middleware.LoggerWithConfig(config)
}
