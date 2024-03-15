package server

import (
	"context"
	"doc-notifier/internal/pkg/server/endpoints"
	"doc-notifier/internal/pkg/watcher"
	"fmt"
	"github.com/labstack/echo/v4"
)

type EchoServer struct {
	options *ServerOptions
	server  *echo.Echo
	watcher *watcher.NotifyWatcher
}

func New(options *ServerOptions, nw *watcher.NotifyWatcher) *EchoServer {
	return &EchoServer{
		options: options,
		watcher: nw,
	}
}

func (s *EchoServer) RunServer() {
	s.server = echo.New()

	// Just store watcher service ptr to get functionality access.
	s.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("Watcher", s.watcher)
			return next(c)
		}
	})

	s.server.POST("/attach", endpoints.AttachDirectories)
	s.server.POST("/detach", endpoints.DetachDirectories)
	s.server.GET("/watcher", endpoints.WatchedDirsList)

	s.server.POST("/download", endpoints.DownloadFile)
	s.server.POST("/upload", endpoints.UploadFile)
	s.server.GET("/upload", endpoints.UploadFileForm)

	address := fmt.Sprintf("%s:%d", s.options.hostAddress, s.options.portNumber)
	_ = s.server.Start(address)
}

func (s *EchoServer) StopServer() {
	_ = s.server.Shutdown(context.Background())
}
