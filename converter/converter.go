package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Convert converts a file from one format to another
// It determines the appropriate converter based on the source and target format
func Convert(sourcePath, targetFormat string) (string, error) {
	ext := strings.ToLower(filepath.Ext(sourcePath))

	// Check if converting from PDF to text
	if ext == ".pdf" && targetFormat == "txt" {
		return ConvertPDFToText(sourcePath)
	}

	// Otherwise use LibreOffice
	return ConvertWithSoffice(sourcePath, targetFormat, nil)
}

// ConvertWithSoffice converts a file using LibreOffice (soffice)
func ConvertWithSoffice(inputPath, format string, options map[string]string) (string, error) {
	// Create a temporary directory for LibreOffice
	tmpDir, err := os.MkdirTemp("", "soffice-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Build format string with optional filter
	formatStr := format
	if filter, ok := options["filter"]; ok {
		formatStr = fmt.Sprintf("%s:\"%s\"", format, filter)
	}

	// Build LibreOffice command
	args := []string{
		fmt.Sprintf("-env:UserInstallation=file://%s", tmpDir),
		"--headless",
		"--convert-to",
		formatStr,
		"--outdir",
		tmpDir,
		inputPath,
	}

	// Execute LibreOffice conversion
	cmd := exec.Command("soffice", args...)
	cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + os.Getenv("PATH"),
		"SHELL=" + os.Getenv("SHELL"),
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("soffice conversion failed: %w\nOutput: %s", err, string(output))
	}

	// Determine output file path
	inputExt := filepath.Ext(inputPath)
	basename := filepath.Base(inputPath)
	basename = strings.TrimSuffix(basename, inputExt)
	
	// For filters with special chars, use just the format name
	outFormat := format
	if strings.Contains(format, ":") {
		outFormat = strings.Split(format, ":")[0]
	}
	
	outFile := filepath.Join(tmpDir, fmt.Sprintf("%s.%s", basename, outFormat))

	// Verify output file exists
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		return "", fmt.Errorf("output file not created: %s", outFile)
	}

	return outFile, nil
}

// ConvertPDFToText converts a PDF to text using pdftotext
func ConvertPDFToText(inputPath string) (string, error) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "pdftotext-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Determine output file path
	inputExt := filepath.Ext(inputPath)
	basename := filepath.Base(inputPath)
	basename = strings.TrimSuffix(basename, inputExt)
	outPath := filepath.Join(tmpDir, basename+".txt")

	// Build pdftotext command
	args := []string{
		"-eol", "unix",
		inputPath,
		outPath,
	}

	cmd := exec.Command("pdftotext", args...)
	cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + os.Getenv("PATH"),
		"SHELL=" + os.Getenv("SHELL"),
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pdftotext conversion failed: %w\nOutput: %s", err, string(output))
	}

	// Verify output file exists
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		return "", fmt.Errorf("output file not created: %s", outPath)
	}

	return outPath, nil
}
