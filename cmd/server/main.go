package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize container (dependency injection)
	container := NewContainer()

	// Start all services
	if err := container.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start services: %v\n", err)
		os.Exit(1)
	}

	// Setup HTTP routes
	setupRoutes(container)

	// Create HTTP server
	server := &http.Server{
		Addr:           container.Config.Server.Port,
		ReadTimeout:    container.Config.Server.ReadTimeout,
		WriteTimeout:   container.Config.Server.WriteTimeout,
		IdleTimeout:    container.Config.Server.IdleTimeout,
		MaxHeaderBytes: container.Config.Server.MaxHeaderBytes,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("\n=== StockTrack Backend Server ===\n")
		fmt.Printf("Running on http://localhost%s\n\n", container.Config.Server.Port)
		fmt.Printf("📋 API Endpoints:\n")
		fmt.Printf("  POST   http://localhost:8080/auth/register\n")
		fmt.Printf("  POST   http://localhost:8080/auth/login\n")
		fmt.Printf("  GET    http://localhost:8080/health\n")
		fmt.Printf("  WS     ws://localhost:8080/ws?token=<JWT>\n\n")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	container.Stop()
}

// setupRoutes configures HTTP routes
func setupRoutes(c *Container) {
	// Auth routes (unprotected)
	http.HandleFunc("/auth/register", c.AuthHandler.Register)
	http.HandleFunc("/auth/login", c.AuthHandler.Login)

	// Health check (unprotected)
	http.HandleFunc("/health", c.AuthHandler.Health)

	// WebSocket route (protected)
	http.HandleFunc("/ws", c.AuthMiddleware.ProtectWebSocket(c.WebSocketHub.Handler))
}
