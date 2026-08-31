package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pictl/pictl/internal/handlers"
	"github.com/pictl/pictl/internal/services"
	"github.com/pictl/pictl/web"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	socketPath := flag.String("socket", "/run/pictl.sock", "path to the pictl helper's Unix socket")
	flag.Parse()

	// Initialize services
	sysService, err := services.NewRemoteSystemService(*socketPath)
	if err != nil {
		log.Fatalf("connect to pictl: %v\n", err)
	}
	defer sysService.Close()

	// Initialize handlers
	h := handlers.NewHandlers(sysService, web.StaticFS)

	// Setup routes
	mux := http.NewServeMux()

	// Serve static files
	mux.HandleFunc("/", h.ServeHTML())
	mux.HandleFunc("/static/", h.ServeStatic())

	// API endpoints
	mux.HandleFunc("/api/system/info", h.GetSystemInfo)
	mux.HandleFunc("/api/system/update", h.UpdateSystem)
	mux.HandleFunc("/api/system/restart", h.RestartSystem)
	mux.HandleFunc("/api/system/shutdown", h.ShutdownSystem)

	// Create HTTP server
	server := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s\n", *addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	fmt.Printf("\nReceived signal: %v\n", sig)
	fmt.Println("Shutting down gracefully...")

	if err := server.Close(); err != nil {
		log.Fatalf("Server close error: %v\n", err)
	}

	fmt.Println("Server stopped")
}
