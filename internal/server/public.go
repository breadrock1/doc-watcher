package server

import "context"

type Server struct {
	Server Session
}

type Session interface {
	Start(_ context.Context) error
	Shutdown(ctx context.Context) error
}
