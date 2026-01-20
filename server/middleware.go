package server

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"
)

// RequireAuth is middleware that validates the Bearer token
func (s *Server) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip token check if API_TOKEN is "*" or empty
		if s.apiToken == "*" || s.apiToken == "" {
			next(w, r)
			return
		}

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED")
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED")
			return
		}

		token := parts[1]

		// Constant-time comparison to prevent timing attacks
		if !secureCompare(token, s.apiToken) {
			sendJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED")
			return
		}

		next(w, r)
	}
}

// secureCompare performs constant-time string comparison
func secureCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// sendJSONError sends a JSON error response
func sendJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
