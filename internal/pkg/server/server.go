package server

import (
	"context"
	_ "doc-notifier/docs"
	"doc-notifier/internal/pkg/server/endpoints"
	"doc-notifier/internal/pkg/watcher"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	s.server.Use(middleware.Logger())
	s.server.Use(middleware.CORS())

	// Just store watcher service ptr to get functionality access.
	s.server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("Watcher", s.watcher)
			return next(c)
		}
	})

	s.server.GET("/hello/", endpoints.Hello)

	s.server.POST("/watcher/attach", endpoints.AttachDirectories)
	s.server.POST("/watcher/detach", endpoints.DetachDirectories)
	s.server.GET("/watcher/all", endpoints.WatchedDirsList)
	s.server.GET("/watcher/pause", endpoints.PauseWatchers)
	s.server.GET("/watcher/run", endpoints.RunWatchers)

	s.server.POST("/files/download", endpoints.DownloadFile)
	s.server.POST("/files/upload", endpoints.UploadFiles)
	s.server.GET("/files/upload", endpoints.UploadFileForm)
	s.server.POST("/files/analyse", endpoints.AnalyseFiles)
	s.server.POST("/files/move", endpoints.MoveFile)
	s.server.POST("/files/remove", endpoints.RemoveFile)
	s.server.GET("/files/unrecognized", endpoints.GetUnrecognized)

	s.server.GET("/swagger/*", echoSwagger.WrapHandler)

	address := fmt.Sprintf("%s:%d", s.options.hostAddress, s.options.portNumber)
	_ = s.server.Start(address)
}

func (s *EchoServer) StopServer() {
	_ = s.server.Shutdown(context.Background())
}
