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
		fmt.Printf("\n🔐 Authentication (Public):\n")
		fmt.Printf("  POST   /auth/register\n")
		fmt.Printf("  POST   /auth/login\n")
		fmt.Printf("\n💹 Trading (Protected):\n")
		fmt.Printf("  POST   /trading/trade\n")
		fmt.Printf("  GET    /trading/wallet\n")
		fmt.Printf("  POST   /trading/deposit\n")
		fmt.Printf("  POST   /trading/alerts\n")
		fmt.Printf("  GET    /trading/alerts\n")
		fmt.Printf("  GET    /trading/candles?symbol=SYMBOL&limit=100\n")
		fmt.Printf("\n📊 WebSocket (Protected):\n")
		fmt.Printf("  WS     ws://localhost:8080/ws?token=<JWT>\n")
		fmt.Printf("\n🏥 Health:\n")
		fmt.Printf("  GET    /health\n\n")

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

	// Trading routes (protected)
	http.HandleFunc("/trading/trade", c.AuthMiddleware.Protect(c.TradingHandler.ExecuteTrade))
	http.HandleFunc("/trading/wallet", c.AuthMiddleware.Protect(c.TradingHandler.GetWalletSnapshot))
	http.HandleFunc("/trading/deposit", c.AuthMiddleware.Protect(c.TradingHandler.DepositCash))
	http.HandleFunc("/trading/alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			c.AuthMiddleware.Protect(c.TradingHandler.CreateAlert)(w, r)
		} else if r.Method == "GET" {
			c.AuthMiddleware.Protect(c.TradingHandler.GetUserAlerts)(w, r)
		}
	})
	http.HandleFunc("/trading/candles", c.TradingHandler.GetCandles)

	// User profile routes (protected)
	http.HandleFunc("/user/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			c.AuthMiddleware.Protect(c.ProfileHandler.GetProfile)(w, r)
		} else if r.Method == "PUT" {
			c.AuthMiddleware.Protect(c.ProfileHandler.UpdateProfile)(w, r)
		}
	})
	http.HandleFunc("/user/sessions", c.AuthMiddleware.Protect(c.ProfileHandler.GetSessions))
	http.HandleFunc("/user/logout", c.AuthMiddleware.Protect(c.ProfileHandler.Logout))

	// WebSocket route (protected)
	http.HandleFunc("/ws", c.AuthMiddleware.ProtectWebSocket(c.WebSocketHub.Handler))
}
