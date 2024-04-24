package server

import (
	"context"
	"doc-notifier/internal/pkg/server/endpoints"
	"doc-notifier/internal/pkg/watcher"
	"fmt"
	"github.com/labstack/echo/v4"
)

type EchoServer struct {
	options *Options
	server  *echo.Echo
	watcher *watcher.NotifyWatcher
}

func New(options *Options, nw *watcher.NotifyWatcher) *EchoServer {
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

	s.server.GET("/hello/", endpoints.Hello)

	s.server.POST("/watcher/create", endpoints.CreateWatchDirectory)
	s.server.DELETE("/watcher/remove", endpoints.RemoveWatchDirectory)
	s.server.GET("/watcher/all", endpoints.GetWatchedDirectories)
	s.server.GET("/watcher/unrecognized", endpoints.GetUnrecognizedFiles)

	s.server.POST("/file/download", endpoints.DownloadFile)
	s.server.POST("/file/upload", endpoints.UploadFile)
	s.server.GET("/file/upload", endpoints.UploadFileForm)

	s.server.GET("/watcher/stop", endpoints.PauseWatchers)
	s.server.GET("/watcher/run", endpoints.RunWatchers)

	address := fmt.Sprintf("%s:%d", s.options.hostAddress, s.options.portNumber)
	_ = s.server.Start(address)
}

func (s *EchoServer) StopServer() {
	_ = s.server.Shutdown(context.Background())
}
