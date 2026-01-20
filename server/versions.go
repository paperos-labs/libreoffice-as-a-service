package server

import (
	"encoding/json"
	"net/http"
	"runtime"
)

// VersionsHandler returns version information
func (s *Server) VersionsHandler(w http.ResponseWriter, r *http.Request) {
	versions := map[string]string{
		"api": "1.2.0",
		"go":  runtime.Version(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)
}
