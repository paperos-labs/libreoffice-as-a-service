package server

import (
	"log"
	"net/http"

	"github.com/therootcompany/golib/http/middleware"
)

// Config holds server configuration
type Config struct {
	Port     int
	Bind     string
	APIToken string
}

// Server represents the HTTP server
type Server struct {
	apiToken string
}

// NewServer creates a new server instance
func NewServer(cfg Config) http.Handler {
	s := &Server{
		apiToken: cfg.APIToken,
	}

	mux := http.NewServeMux()
	m := middleware.WithMux(mux, s.LogPanics)
	
	// Public routes
	m.HandleFunc("GET /api/versions", s.VersionsHandler)
	
	// Authenticated routes
	authM := m.With(s.RequireAuth)
	authM.HandleFunc("POST /api/convert/{format}", s.ConvertHandler)

	// Static files
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	return mux
}

// LogPanics is middleware to recover from panics and log them
func (s *Server) LogPanics(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}
