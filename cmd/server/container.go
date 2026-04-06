package main

import (
	"fmt"

	"stocktrack-backend/config"
	"stocktrack-backend/internal/adapter/auth"
	"stocktrack-backend/internal/adapter/db"
	"stocktrack-backend/internal/adapter/http"
	"stocktrack-backend/internal/adapter/storage"
	"stocktrack-backend/internal/adapter/ws"
	"stocktrack-backend/internal/domain"
)

type Container struct {
	Config                *config.Config
	Database              *db.Database
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
	DashboardHandler      *http.DashboardHandler
	AuthMiddleware        *http.AuthMiddleware
}

func NewContainer() *Container {
	cfg := config.NewDefaultConfig()

	var database *db.Database
	var userRepo domain.UserRepository
	var candleRepo domain.CandleRepository
	var alertRepo domain.AlertRepository
	var transactionRepo domain.TransactionRepository

	if cfg.Database.Driver == "postgres" {
		postgresDB, err := db.NewDatabase(db.Config{
			URL:      cfg.Database.URL,
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			DBName:   cfg.Database.DBName,
			SSLMode:  cfg.Database.SSLMode,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to connect to PostgreSQL: %v", err))
		}

		if err := postgresDB.Migrate(); err != nil {
			_ = postgresDB.Close()
			panic(fmt.Sprintf("failed to run PostgreSQL migrations: %v", err))
		}

		database = postgresDB
		userRepo = storage.NewPostgresUserRepository(postgresDB.GetConn())
		candleRepo = storage.NewPostgresCandleRepository(postgresDB.GetConn())
		alertRepo = storage.NewPostgresAlertRepository(postgresDB.GetConn())
		transactionRepo = storage.NewPostgresTransactionRepository(postgresDB.GetConn())
		fmt.Println("[DB] Using PostgreSQL repositories")
	} else {
		userRepo = storage.NewInMemoryUserRepository()
		candleRepo = storage.NewInMemoryCandleRepository()
		alertRepo = storage.NewInMemoryAlertRepository()
		transactionRepo = storage.NewInMemoryTransactionRepository()
		fmt.Println("[DB] Using in-memory repositories")
	}

	tokenProvider := auth.NewJWTProvider(cfg.Auth.JWTSecretKey, cfg.Auth.TokenExpiry)

	authService := domain.NewAuthService(userRepo, tokenProvider)

	walletService := domain.NewWalletService(transactionRepo)

	alertService := domain.NewAlertService(alertRepo)

	ohlcAggregator := domain.NewOHLCAggregator(candleRepo)

	marketEngine := domain.NewMarketEngine()
	marketEngine.SetOHLCAggregator(ohlcAggregator)
	marketEngine.SetAlertService(alertService)
	marketEngine.SetTransactionRepository(transactionRepo)

	// Set market engine on wallet service for price lookups
	walletService.SetMarketEngine(marketEngine)

	wsHub := ws.NewHub()

	authHandler := http.NewAuthHandler(authService)
	tradingHandler := http.NewTradingHandler(walletService, alertService, ohlcAggregator, transactionRepo)
	profileHandler := http.NewProfileHandler(authService)
	dashboardHandler := http.NewDashboardHandler(walletService)

	authMiddleware := http.NewAuthMiddleware(authService)

	return &Container{
		Config:                cfg,
		Database:              database,
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
		DashboardHandler:      dashboardHandler,
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

	if c.Database != nil {
		_ = c.Database.Close()
	}
}
