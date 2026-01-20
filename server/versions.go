package server

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/paperos-labs/libreoffice-as-a-service/config"
)

// VersionsHandler returns version information
func VersionsHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		versions := map[string]string{
			"api": "1.2.0",
			"go":  runtime.Version(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versions)
	}
}
