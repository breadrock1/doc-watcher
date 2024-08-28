package server

import (
	"context"
	_ "doc-notifier/docs"
	"doc-notifier/internal/office"
	"doc-notifier/internal/watcher"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Service struct {
	server  *echo.Echo
	office  *office.Service
	watcher *watcher.Service
}

func New(nw *watcher.Service, officeService *office.Service) *Service {
	server := &Service{
		office:  officeService,
		watcher: nw,
	}

	server.setupServer()
	return server
}

func (s *Service) setupServer() {
	s.server = echo.New()

	s.server.Use(middleware.CORS())
	s.server.Use(middleware.Logger())
	s.server.Use(middleware.Recover())

	_ = s.CreateHelloGroup()
	_ = s.CreateStorageGroup()
	_ = s.CreateWatcherGroup()

	s.server.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (s *Service) RunServer(_ context.Context) error {
	return s.server.Start(s.watcher.Watcher.GetAddress())
}

func (s *Service) StopServer(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
