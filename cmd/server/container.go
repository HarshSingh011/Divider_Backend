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
	Config                *config.Config
	UserRepository        domain.UserRepository
	CandleRepository      domain.CandleRepository
	AlertRepository       domain.AlertRepository
	TransactionRepository domain.TransactionRepository
	TokenProvider         domain.TokenProvider
	AuthService           domain.AuthService
	WalletService         domain.WalletService
	AlertService          domain.AlertService
	OHLCService           domain.OHLCService
	MarketEngine          *domain.MarketEngine
	WebSocketHub          *ws.Hub
	AuthHandler           *http.AuthHandler
	TradingHandler        *http.TradingHandler
	ProfileHandler        *http.ProfileHandler
	AuthMiddleware        *http.AuthMiddleware
}

// NewContainer initializes and wires all dependencies
func NewContainer() *Container {
	cfg := config.NewDefaultConfig()

	// Initialize repositories (using in-memory for now)
	userRepo := storage.NewInMemoryUserRepository()
	candleRepo := storage.NewInMemoryCandleRepository()
	alertRepo := storage.NewInMemoryAlertRepository()
	transactionRepo := storage.NewInMemoryTransactionRepository()

	fmt.Println("✓ Using in-memory storage")

	// Initialize token provider
	tokenProvider := auth.NewJWTProvider(cfg.Auth.JWTSecretKey, cfg.Auth.TokenExpiry)

	// Initialize auth service
	authService := domain.NewAuthService(userRepo, tokenProvider)

	// Initialize wallet service
	walletService := domain.NewWalletService(transactionRepo)

	// Initialize alert service
	alertService := domain.NewAlertService(alertRepo)

	// Initialize OHLC aggregator
	ohlcAggregator := domain.NewOHLCAggregator(candleRepo)

	// Initialize market engine with integrations
	marketEngine := domain.NewMarketEngine()
	marketEngine.SetOHLCAggregator(ohlcAggregator)
	marketEngine.SetAlertService(alertService)

	// Initialize WebSocket hub
	wsHub := ws.NewHub()

	// Initialize HTTP handlers
	authHandler := http.NewAuthHandler(authService)
	tradingHandler := http.NewTradingHandler(walletService, alertService, ohlcAggregator)
	profileHandler := http.NewProfileHandler(authService)

	// Initialize middleware
	authMiddleware := http.NewAuthMiddleware(authService)

	return &Container{
		Config:                cfg,
		UserRepository:        userRepo,
		CandleRepository:      candleRepo,
		AlertRepository:       alertRepo,
		TransactionRepository: transactionRepo,
		TokenProvider:         tokenProvider,
		AuthService:           authService,
		WalletService:         walletService,
		AlertService:          alertService,
		OHLCService:           ohlcAggregator,
		MarketEngine:          marketEngine,
		WebSocketHub:          wsHub,
		AuthHandler:           authHandler,
		TradingHandler:        tradingHandler,
		ProfileHandler:        profileHandler,
		AuthMiddleware:        authMiddleware,
	}
}

// Start initializes all background services
func (c *Container) Start() error {
	fmt.Println("Starting application...")

	// Start OHLC aggregator
	c.OHLCService.(*domain.OHLCAggregator).Start()
	fmt.Println("✓ OHLC aggregator started")

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
	ohlc := c.OHLCService.(*domain.OHLCAggregator)
	ohlc.Stop()
	c.WebSocketHub.Stop()
	fmt.Println("✓ Services stopped")
}
