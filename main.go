package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/paperos-labs/libreoffice-as-a-service/server"
)

const (
	name         = "libreoffice-as-a-service"
	licenseYear  = "2024"
	licenseOwner = "Paperos Labs"
	licenseType  = "MPL-2.0"
)

// set by GoReleaser via ldflags
var (
	version = "1.2.0"
	commit  = "0000000"
	date    = "0001-01-01T00:00:00Z"
)

// MainConfig holds main configuration
type MainConfig struct {
	Port     int
	Bind     string
	APIToken string
}

// printVersion displays the version, commit, and build date.
func printVersion() {
	if len(commit) > 7 {
		commit = commit[:7]
	}
	fmt.Fprintf(os.Stderr, "%s v%s %s (%s)\n", name, version, commit, date)
	fmt.Fprintf(os.Stderr, "Copyright (C) %s %s\n", licenseYear, licenseOwner)
	fmt.Fprintf(os.Stderr, "Licensed under the %s license\n", licenseType)
}

func main() {
	// Load .env files
	homeDir, _ := os.UserHomeDir()
	_ = godotenv.Load(filepath.Join(homeDir, ".env"))
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env.secret")

	// Get defaults from environment
	cfg := MainConfig{
		Port:     getEnvInt("PORT", 5227),
		Bind:     getEnv("BIND", "0.0.0.0"),
		APIToken: getAPIToken(),
	}

	var showVersion bool
	mainFlags := flag.NewFlagSet("", flag.ContinueOnError)
	mainFlags.BoolVar(&showVersion, "version", false, "Print version and exit")
	mainFlags.IntVar(&cfg.Port, "port", cfg.Port, "Port to listen on")
	mainFlags.StringVar(&cfg.Bind, "bind", cfg.Bind, "Address to bind to")

	mainFlags.Usage = func() {
		printVersion()
		out := mainFlags.Output()
		_, _ = fmt.Fprintf(out, "\n")
		_, _ = fmt.Fprintf(out, "USAGE\n")
		_, _ = fmt.Fprintf(out, "   %s [options]\n", name)
		_, _ = fmt.Fprintf(out, "\n")
		mainFlags.PrintDefaults()
	}

	// Handle version and help flags out-of-band
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-V", "version", "-version", "--version":
			printVersion()
			return
		case "help", "-help", "--help":
			mainFlags.SetOutput(os.Stdout)
			mainFlags.Usage()
			return
		}
	}

	if err := mainFlags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		mainFlags.SetOutput(os.Stderr)
		mainFlags.Usage()
		os.Exit(1)
		return
	}

	// Handle --version flag after parsing
	if showVersion {
		printVersion()
		return
	}

	// Create server with configuration
	srv := server.NewServer(server.Config{
		Port:     cfg.Port,
		Bind:     cfg.Bind,
		APIToken: cfg.APIToken,
	})

	addr := fmt.Sprintf("%s:%d", cfg.Bind, cfg.Port)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           srv,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       120 * time.Second, // Allow time for file uploads
		WriteTimeout:      120 * time.Second, // Allow time for conversions
		MaxHeaderBytes:    1024 * 1024,       // 1MiB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("\nListening on http://%s\n\n", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("\nShutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func getAPIToken() string {
	// Check LAAS_API_TOKEN first, then API_TOKEN
	if token := os.Getenv("LAAS_API_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("API_TOKEN")
}
