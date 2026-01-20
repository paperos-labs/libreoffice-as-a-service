package server

import (
	"net/http"

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
		
		if r.Method == "POST" && len(r.URL.Path) > 13 && r.URL.Path[:13] == "/api/convert/" {
			// Extract format from path
			format := r.URL.Path[13:]
			// Create a new request with PathValue support
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
