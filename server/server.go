package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server represents the HTTP server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates a new server instance.
func New(address string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}
}

// Start starts the server.
func (s *Server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server failed to shut down gracefully: %v", err)
	}
}
