package api

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server represents the HTTP server
type Server struct {
	server  *http.Server
	handler *Handler
}

// NewServer creates a new HTTP server
func NewServer(addr string, handler *Handler) *Server {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/v1/job/submit", handler.SubmitJob)
	mux.HandleFunc("/api/v1/job/status", handler.GetJobStatus)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return &Server{
		server:  server,
		handler: handler,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Server starting on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
