package server

import (
	"doc-notifier/internal/pkg/server/endpoints"
	"doc-notifier/internal/pkg/watcher"
	"fmt"
	"github.com/labstack/echo/v4"
)

type EchoServer struct {
	options *ServerOptions
	watcher *watcher.NotifyWatcher
}

func New(options *ServerOptions, nw *watcher.NotifyWatcher) *EchoServer {
	return &EchoServer{
		options: options,
		watcher: nw,
	}
}

func (s *EchoServer) RunServer() {
	e := echo.New()

	// Just store watcher service ptr to get functionality access.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("Watcher", s.watcher)
			return next(c)
		}
	})

	e.POST("/attach", endpoints.AttachDirectories)
	e.POST("/detach", endpoints.DetachDirectories)
	e.GET("/watcher", endpoints.WatchedDirsList)

	e.POST("/download", endpoints.DownloadFile)
	e.POST("/upload", endpoints.UploadFile)
	e.GET("/upload", endpoints.UploadFileForm)

	address := fmt.Sprintf("%s:%d", s.options.hostAddress, s.options.portNumber)
	_ = e.Start(address)
}
