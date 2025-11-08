package server

import (
	"log/slog"
	"net/http"

	"github.com/junghwan16/test-server/internal/health"
)

type Server struct {
	router        *http.ServeMux
	logger        *slog.Logger
	healthHandler *health.Handler
}

func New(logger *slog.Logger, healthHandler *health.Handler) *Server {
	s := &Server{
		router:        http.NewServeMux(),
		logger:        logger,
		healthHandler: healthHandler,
	}
	s.setupHealthRoutes()
	return s
}

func (s *Server) setupHealthRoutes() {
	s.router.HandleFunc("GET /healthz", s.healthHandler.HandleLiveness)
	s.router.HandleFunc("GET /readyz", s.healthHandler.HandleReadiness)
}

// ServeHTTP allows Server to implement the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Router exposes the router so additional routes can be registered
func (s *Server) Router() *http.ServeMux {
	return s.router
}
