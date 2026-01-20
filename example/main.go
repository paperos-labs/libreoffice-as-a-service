package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	baseURL := os.Getenv("LAAS_BASE_URL")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:5227"
	}

	apiToken := os.Getenv("LAAS_API_TOKEN")

	// Read input file
	inputFile := "fixtures/Writing1.docx"
	outputFile := "Writing1.pdf"

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Create request
	url := fmt.Sprintf("%s/api/convert/pdf?filename=Writing1.docx", baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	if apiToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("******", apiToken))
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error: Server returned %d\n%s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	// Write output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}
