package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	container := NewContainer()

	if err := container.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start services: %v\n", err)
		os.Exit(1)
	}

	setupRoutes(container)

	server := &http.Server{
		Addr:           container.Config.Server.Port,
		ReadTimeout:    container.Config.Server.ReadTimeout,
		WriteTimeout:   container.Config.Server.WriteTimeout,
		IdleTimeout:    container.Config.Server.IdleTimeout,
		MaxHeaderBytes: container.Config.Server.MaxHeaderBytes,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	container.Stop()
}

func setupRoutes(c *Container) {
	// CORS middleware wrapper
	corsHandler := func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-Email, X-Username")
			w.Header().Set("Access-Control-Max-Age", "3600")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			handler(w, r)
		}
	}

	http.HandleFunc("/auth/register", corsHandler(c.AuthHandler.Register))
	http.HandleFunc("/auth/login", corsHandler(c.AuthHandler.Login))
	http.HandleFunc("/health", corsHandler(c.AuthHandler.Health))

	http.HandleFunc("/trading/trade", corsHandler(c.AuthMiddleware.Protect(c.TradingHandler.ExecuteTrade)))
	http.HandleFunc("/trading/wallet", corsHandler(c.AuthMiddleware.Protect(c.TradingHandler.GetWalletSnapshot)))
	http.HandleFunc("/trading/deposit", corsHandler(c.AuthMiddleware.Protect(c.TradingHandler.DepositCash)))
	http.HandleFunc("/trading/alerts", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			c.AuthMiddleware.Protect(c.TradingHandler.CreateAlert)(w, r)
		} else if r.Method == "GET" {
			c.AuthMiddleware.Protect(c.TradingHandler.GetUserAlerts)(w, r)
		}
	}))
	http.HandleFunc("/trading/candles", corsHandler(c.TradingHandler.GetCandles))

	http.HandleFunc("/user/profile", corsHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			c.AuthMiddleware.Protect(c.ProfileHandler.GetProfile)(w, r)
		} else if r.Method == "PUT" {
			c.AuthMiddleware.Protect(c.ProfileHandler.UpdateProfile)(w, r)
		}
	}))
	http.HandleFunc("/user/sessions", corsHandler(c.AuthMiddleware.Protect(c.ProfileHandler.GetSessions)))
	http.HandleFunc("/user/logout", corsHandler(c.AuthMiddleware.Protect(c.ProfileHandler.Logout)))

	http.HandleFunc("/ws", c.AuthMiddleware.ProtectWebSocket(c.WebSocketHub.Handler))
}
