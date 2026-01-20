package server

import (
	"net/http"

	"github.com/paperos-labs/libreoffice-as-a-service/config"
)

// NewMux creates a new HTTP request multiplexer with all routes
func NewMux(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	// API routes with logging middleware
	mux.HandleFunc("GET /api/versions", LoggingMiddleware(VersionsHandler(cfg)))
	mux.HandleFunc("POST /api/convert/{format}", LoggingMiddleware(AuthMiddleware(cfg, ConvertHandler(cfg))))

	return mux
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer that captures the status code
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(lrw, r)

		// Log request details
		// Using a simpler format than fastify's default logger
		// Could be enhanced with a proper logging library if needed
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
