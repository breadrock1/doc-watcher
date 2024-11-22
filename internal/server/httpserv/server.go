package httpserv

import (
	"context"

	"doc-watcher/internal/server"
	"doc-watcher/internal/watcher"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "doc-watcher/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Service struct {
	config  *server.Config
	server  *echo.Echo
	watcher *watcher.Service
}

func New(servConf *server.Config, nw *watcher.Service) *server.Server {
	httpServer := &Service{
		config:  servConf,
		watcher: nw,
	}

	return &server.Server{
		Server: httpServer,
	}
}

func (s *Service) setupServer() {
	s.server = echo.New()

	s.server.Use(middleware.CORS())
	s.server.Use(middleware.Recover())
	s.server.Use(InitLogger(s.config))

	_ = s.CreateWatcherGroup()

	s.server.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (s *Service) Start(_ context.Context) error {
	s.setupServer()
	return s.server.Start(s.config.Address)
}

func (s *Service) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
