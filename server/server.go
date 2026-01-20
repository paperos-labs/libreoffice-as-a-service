package server

import (
	"log"
	"net/http"
	"time"
)

// Config holds server configuration
type Config struct {
	Port     int
	Bind     string
	APIToken string
}

// Server represents the HTTP server
type Server struct {
	mux      *http.ServeMux
	apiToken string
}

// NewServer creates a new server instance
func NewServer(cfg Config) *Server {
	s := &Server{
		mux:      http.NewServeMux(),
		apiToken: cfg.APIToken,
	}

	// Register routes
	s.mux.HandleFunc("GET /api/versions", LogPanics(s.VersionsHandler))
	s.mux.HandleFunc("POST /api/convert/{format}", LogPanics(s.RequireAuth(s.ConvertHandler)))

	// Static files
	fs := http.FileServer(http.Dir("./public"))
	s.mux.Handle("/", fs)

	return s
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lrw := &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	s.mux.ServeHTTP(lrw, r)

	duration := time.Since(start)
	log.Printf("%s %s - %d (%v)", r.Method, r.URL.Path, lrw.statusCode, duration)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LogPanics is middleware to recover from panics and log them
func LogPanics(next http.HandlerFunc) http.HandlerFunc {
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
