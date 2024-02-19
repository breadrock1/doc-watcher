package server

import (
	"doc-notifier/internal/pkg/server/endpoints"
	"doc-notifier/internal/pkg/watcher"
	"fmt"
	"github.com/labstack/echo/v4"
)

type Server struct {
	options *ServerOptions
	watcher *watcher.NotifyWatcher
}

func New(options *ServerOptions, nw *watcher.NotifyWatcher) *Server {
	return &Server{
		options: options,
		watcher: nw,
	}
}

func (s *Server) RunServer() {
	e := echo.New()

	e.POST("/download", endpoints.DownloadFile)

	address := fmt.Sprintf("%s:%d", s.options.hostAddress, s.options.portNumber)
	_ = e.Start(address)
}
