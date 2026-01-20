package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/paperos-labs/libreoffice-as-a-service/config"
)

// NewMux creates a new HTTP request multiplexer with all routes
func NewMux(cfg *config.Config) http.Handler {
	// Create a simple router that dispatches based on path
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for API routes first
		if r.Method == "GET" && r.URL.Path == "/api/versions" {
			LoggingMiddleware(VersionsHandler(cfg))(w, r)
			return
		}
		
		// Check for convert API route
		if r.Method == "POST" && strings.HasPrefix(r.URL.Path, "/api/convert/") {
			// Extract format from path
			format := strings.TrimPrefix(r.URL.Path, "/api/convert/")
			// Validate format is non-empty and doesn't contain path separators
			if format == "" || strings.Contains(format, "/") {
				http.Error(w, "Invalid format parameter", http.StatusBadRequest)
				return
			}
			// Set the format as a path value
			r.SetPathValue("format", format)
			LoggingMiddleware(AuthMiddleware(cfg, ConvertHandler(cfg)))(w, r)
			return
		}
		
		// Serve static files for everything else
		fs := http.FileServer(http.Dir("./public"))
		fs.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer that captures the status code
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(lrw, r)

		// Log request details
		duration := time.Since(start)
		log.Printf("%s %s - %d (%v)", r.Method, r.URL.Path, lrw.statusCode, duration)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
