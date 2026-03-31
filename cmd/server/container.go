package main

import (
	"stocktrack-backend/config"
	"stocktrack-backend/internal/adapter/auth"
	"stocktrack-backend/internal/adapter/http"
	"stocktrack-backend/internal/adapter/storage"
	"stocktrack-backend/internal/adapter/ws"
	"stocktrack-backend/internal/domain"
)

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

func NewContainer() *Container {
	cfg := config.NewDefaultConfig()

	userRepo := storage.NewInMemoryUserRepository()
	candleRepo := storage.NewInMemoryCandleRepository()
	alertRepo := storage.NewInMemoryAlertRepository()
	transactionRepo := storage.NewInMemoryTransactionRepository()

	tokenProvider := auth.NewJWTProvider(cfg.Auth.JWTSecretKey, cfg.Auth.TokenExpiry)

	authService := domain.NewAuthService(userRepo, tokenProvider)

	walletService := domain.NewWalletService(transactionRepo)

	alertService := domain.NewAlertService(alertRepo)

	ohlcAggregator := domain.NewOHLCAggregator(candleRepo)

	marketEngine := domain.NewMarketEngine()
	marketEngine.SetOHLCAggregator(ohlcAggregator)
	marketEngine.SetAlertService(alertService)

	wsHub := ws.NewHub()

	authHandler := http.NewAuthHandler(authService)
	tradingHandler := http.NewTradingHandler(walletService, alertService, ohlcAggregator)
	profileHandler := http.NewProfileHandler(authService)

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

func (c *Container) Start() error {
	c.OHLCService.(*domain.OHLCAggregator).Start()

	c.MarketEngine.Start()

	c.WebSocketHub.Start()

	go func() {
		for marketData := range c.MarketEngine.Subscribe() {
			c.WebSocketHub.Broadcast(marketData)
		}
	}()

	return nil
}

func (c *Container) Stop() {
	c.MarketEngine.Stop()
	ohlc := c.OHLCService.(*domain.OHLCAggregator)
	ohlc.Stop()
	c.WebSocketHub.Stop()
}
