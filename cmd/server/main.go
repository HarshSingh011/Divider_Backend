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
	http.HandleFunc("/auth/register", c.AuthHandler.Register)
	http.HandleFunc("/auth/login", c.AuthHandler.Login)

	http.HandleFunc("/health", c.AuthHandler.Health)

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

	http.HandleFunc("/user/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			c.AuthMiddleware.Protect(c.ProfileHandler.GetProfile)(w, r)
		} else if r.Method == "PUT" {
			c.AuthMiddleware.Protect(c.ProfileHandler.UpdateProfile)(w, r)
		}
	})
	http.HandleFunc("/user/sessions", c.AuthMiddleware.Protect(c.ProfileHandler.GetSessions))
	http.HandleFunc("/user/logout", c.AuthMiddleware.Protect(c.ProfileHandler.Logout))

	http.HandleFunc("/ws", c.AuthMiddleware.ProtectWebSocket(c.WebSocketHub.Handler))
}
