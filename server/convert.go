package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/paperos-labs/libreoffice-as-a-service/converter"
)

// ConvertHandler handles file conversion requests
func (s *Server) ConvertHandler(w http.ResponseWriter, r *http.Request) {
	// Get target format from URL path parameter
	format := r.PathValue("format")
	if format == "" {
		sendJSONError(w, http.StatusBadRequest, "Missing format parameter")
		return
	}

	// Get filename from query parameter
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		sendJSONError(w, http.StatusBadRequest, "Missing filename query parameter")
		return
	}

	// Sanitize filename to prevent directory traversal
	filename = filepath.Base(filename)

	// Create temporary directory for incoming file
	tmpDir, err := os.MkdirTemp("", "laas-source-")
	if err != nil {
		log.Printf("Failed to create temp dir: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Write uploaded file to temporary location
	sourcePath := filepath.Join(tmpDir, filename)
	sourceFile, err := os.Create(sourcePath)
	if err != nil {
		log.Printf("Failed to create source file: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Copy request body to file
	_, err = io.Copy(sourceFile, r.Body)
	sourceFile.Close()
	if err != nil {
		log.Printf("Failed to write uploaded file: %v", err)
		os.Remove(sourcePath)
		sendJSONError(w, http.StatusInternalServerError, "Failed to receive file")
		return
	}

	log.Printf("[convert] Received file: %s (format: %s)", filename, format)

	// Convert the file
	convertedPath, err := converter.Convert(sourcePath, format)
	if err != nil {
		log.Printf("Conversion failed: %v", err)
		os.Remove(sourcePath)
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Conversion failed: %v", err))
		return
	}

	// Clean up source file
	defer os.Remove(sourcePath)
	defer os.Remove(convertedPath)

	// Open converted file
	convertedFile, err := os.Open(convertedPath)
	if err != nil {
		log.Printf("Failed to open converted file: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer convertedFile.Close()

	// Get file info for content-length
	fileInfo, err := convertedFile.Stat()
	if err != nil {
		log.Printf("Failed to stat converted file: %v", err)
		sendJSONError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Set headers
	suggestedName := filepath.Base(convertedPath)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, suggestedName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Stream the converted file back to client
	_, err = io.Copy(w, convertedFile)
	if err != nil {
		log.Printf("Failed to send converted file: %v", err)
		return
	}

	log.Printf("[convert] Success: %s -> %s", filename, suggestedName)
}
