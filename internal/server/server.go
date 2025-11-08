package server

import (
	"log/slog"
	"net/http"
)

type Server struct {
	router *http.ServeMux
	logger *slog.Logger
}

func New(logger *slog.Logger) *Server {
	return &Server{
		router: http.NewServeMux(),
		logger: logger,
	}
}

// ServeHTTP allows Server to implement the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Router exposes the router so additional routes can be registered
func (s *Server) Router() *http.ServeMux {
	return s.router
}
