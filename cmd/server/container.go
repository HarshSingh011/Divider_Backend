package main

import (
	"fmt"

	"stocktrack-backend/config"
	"stocktrack-backend/internal/adapter/auth"
	"stocktrack-backend/internal/adapter/http"
	"stocktrack-backend/internal/adapter/storage"
	"stocktrack-backend/internal/adapter/ws"
	"stocktrack-backend/internal/domain"
)

// Container holds all application dependencies
type Container struct {
	Config            *config.Config
	UserRepository    domain.UserRepository
	TokenProvider     domain.TokenProvider
	AuthService       domain.AuthService
	MarketEngine      *domain.MarketEngine
	WebSocketHub      *ws.Hub
	AuthHandler       *http.AuthHandler
	AuthMiddleware    *http.AuthMiddleware
}

// NewContainer initializes and wires all dependencies
func NewContainer() *Container {
	cfg := config.NewDefaultConfig()

	// Initialize repositories
	userRepo := storage.NewInMemoryUserRepository()

	// Initialize token provider
	tokenProvider := auth.NewJWTProvider(cfg.Auth.JWTSecretKey, cfg.Auth.TokenExpiry)

	// Initialize auth service
	authService := domain.NewAuthService(userRepo, tokenProvider)

	// Initialize market engine
	marketEngine := domain.NewMarketEngine()

	// Initialize WebSocket hub
	wsHub := ws.NewHub()

	// Initialize HTTP handlers
	authHandler := http.NewAuthHandler(authService)

	// Initialize middleware
	authMiddleware := http.NewAuthMiddleware(authService)

	return &Container{
		Config:         cfg,
		UserRepository: userRepo,
		TokenProvider:  tokenProvider,
		AuthService:    authService,
		MarketEngine:   marketEngine,
		WebSocketHub:   wsHub,
		AuthHandler:    authHandler,
		AuthMiddleware: authMiddleware,
	}
}

// Start initializes all background services
func (c *Container) Start() error {
	fmt.Println("Starting application...")

	// Start market engine
	c.MarketEngine.Start()
	fmt.Println("✓ Market engine started")

	// Start WebSocket hub
	c.WebSocketHub.Start()
	fmt.Println("✓ WebSocket hub started")

	// Start publishing market data to WebSocket clients
	go func() {
		for marketData := range c.MarketEngine.Subscribe() {
			c.WebSocketHub.Broadcast(marketData)
		}
	}()
	fmt.Println("✓ Market publishing started")

	return nil
}

// Stop gracefully shuts down all services
func (c *Container) Stop() {
	fmt.Println("\nShutting down...")
	c.MarketEngine.Stop()
	c.WebSocketHub.Stop()
	fmt.Println("✓ Services stopped")
}
