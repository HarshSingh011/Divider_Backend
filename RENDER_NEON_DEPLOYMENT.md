# Deploy to Render + Neon PostgreSQL

Perfect choice! Neon offers a better free tier than Render's PostgreSQL.

---

## STEP 1: Set Up Neon PostgreSQL (Free)

### A. Create Neon Account
1. Go to https://neon.tech
2. Click "Sign Up" 
3. Sign in with GitHub (recommended)
4. Create account

### B. Create PostgreSQL Project on Neon

1. Click "Create Project"
2. Fill in:
   - **Project Name**: `stocktrack`
   - **Database Name**: `stocktrack_db`
   - **Region**: Select closest to you
   - **Compute Size**: Shared (free)
3. Click "Create Project"

### C. Get Connection String from Neon

1. After project created, go to "Connection Details"
2. You'll see a connection string like:
```
postgresql://user:password@ep-xyz.us-east-1.neon.tech/stocktrack_db?sslmode=require
```

3. **Extract these details:**
   - **Host**: `ep-xyz.us-east-1.neon.tech`
   - **Port**: `5432`
   - **Database**: `stocktrack_db`
   - **User**: `neondb_owner` (or the user shown)
   - **Password**: The password shown
   - **SSL Mode**: `require`

---

## STEP 2: Prepare Your Code for Deployment

### A. Update container.go to Use Neon

Edit `cmd/server/container.go`:

```go
package main

import (
	"fmt"
	"os"
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
	DB                    *db.Database
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

	database, err := db.NewDatabase(db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	userRepo := storage.NewPostgresUserRepository(database.GetConn())
	candleRepo := storage.NewPostgresCandleRepository(database.GetConn())
	alertRepo := storage.NewPostgresAlertRepository(database.GetConn())
	transactionRepo := storage.NewPostgresTransactionRepository(database.GetConn())

	tokenProvider := auth.NewJWTProvider(cfg.Auth.JWTSecretKey, cfg.Auth.TokenExpiry)

	authService := domain.NewAuthService(userRepo, tokenProvider)
	walletService := domain.NewWalletService(transactionRepo)
	alertService := domain.NewAlertService(alertRepo)
	ohlcAggregator := domain.NewOHLCAggregator(candleRepo)

	marketEngine := domain.NewMarketEngine()
	marketEngine.SetOHLCAggregator(ohlcAggregator)
	marketEngine.SetAlertService(alertService)

	webSocketHub := ws.NewHub()
	webSocketHub.SetMarketEngine(marketEngine)

	authHandler := http.NewAuthHandler(authService)
	tradingHandler := http.NewTradingHandler(walletService, alertService, ohlcAggregator)
	profileHandler := http.NewProfileHandler(authService)
	authMiddleware := http.NewAuthMiddleware(tokenProvider)

	return &Container{
		Config:                cfg,
		DB:                    database,
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
		WebSocketHub:          webSocketHub,
		AuthHandler:           authHandler,
		TradingHandler:        tradingHandler,
		ProfileHandler:        profileHandler,
		AuthMiddleware:        authMiddleware,
	}
}

func (c *Container) Start() error {
	c.OHLCService.Start()
	c.MarketEngine.Start()
	c.WebSocketHub.Start()
	return nil
}

func (c *Container) Stop() {
	c.MarketEngine.Stop()
	c.WebSocketHub.Stop()
	if c.DB != nil {
		c.DB.Close()
	}
}
```

### B. Add Server Port Configuration

Update your `config/config.go`:

```go
package config

import (
	"os"
	"time"
)

type Config struct {
	Server ServerConfig
	Auth   AuthConfig
}

type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

type AuthConfig struct {
	JWTSecretKey string
	TokenExpiry  time.Duration
}

func NewDefaultConfig() *Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "0.0.0.0:8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key-change-in-production"
	}

	return &Config{
		Server: ServerConfig{
			Port:           port,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		Auth: AuthConfig{
			JWTSecretKey: jwtSecret,
			TokenExpiry:  24 * time.Hour,
		},
	}
}
```

### C. Update main.go to Use Port from Environment

```go
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
		fmt.Printf("Server running on %s\n", server.Addr)
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
```

---

## STEP 3: Deploy to Render

### A. Push Code to GitHub

```powershell
cd c:\Users\harsh\Python\StockTrack
git add .
git commit -m "Add PostgreSQL support and prepare for Render deployment"
git push
```

### B. Create Render Web Service

1. Go to https://render.com
2. Sign up/Sign in with GitHub
3. Click "New +" → "Web Service"
4. Connect your GitHub repository
5. Fill in:
   - **Name**: `stocktrack-api`
   - **Environment**: `Go`
   - **Build Command**: `go build -o stocktrack ./cmd/server`
   - **Start Command**: `./stocktrack`
   - **Instance Type**: Free

### C. Add Environment Variables on Render

In the Web Service settings, go to "Environment" and add:

```
DB_HOST=<neon-host-from-step-1>
DB_PORT=5432
DB_USER=<neon-user-from-step-1>
DB_PASSWORD=<neon-password-from-step-1>
DB_NAME=stocktrack_db
DB_SSLMODE=require
SERVER_PORT=0.0.0.0:10000
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long-here
```

**Example values:**
```
DB_HOST=ep-xyz.us-east-1.neon.tech
DB_PORT=5432
DB_USER=neondb_owner
DB_PASSWORD=abc123xyz456
DB_NAME=stocktrack_db
DB_SSLMODE=require
SERVER_PORT=0.0.0.0:10000
JWT_SECRET=super-secret-key-at-least-32-chars
```

### D. Deploy

1. **After adding environment variables**, click "Deploy"
2. **Wait 3-5 minutes** for build and deployment
3. Check the build logs for any errors
4. Once deployed, you'll get a URL like:
```
https://stocktrack-api.onrender.com
```

---

## STEP 4: Test Your Deployed API

### Health Check
```powershell
curl https://stocktrack-api.onrender.com/health
```

Expected response:
```json
{"status":"OK"}
```

### Register User
```powershell
curl -X POST https://stocktrack-api.onrender.com/auth/register `
  -H "Content-Type: application/json" `
  -d '{
    "email":"test@example.com",
    "username":"testuser",
    "password":"Password123"
  }'
```

---

## STEP 5: Connect Frontend to Deployed API

Update `app.html` (find and replace):

```javascript
const API_BASE = "https://stocktrack-api.onrender.com";

// All API calls now use deployed backend
```

---

## Troubleshooting

**"Database connection failed"**
- Verify environment variables are correct in Render dashboard
- Check DB_HOST, DB_PORT, DB_USER, DB_PASSWORD from Neon

**"SSL certificate verification failed"**
- Make sure DB_SSLMODE=require (not disable)

**"Service crashes immediately"**
- Check Render logs
- Verify all environment variables are set
- Try `/health` endpoint to debug

**"Neon project suspended"**
- Free tier has usage limits
- Upgrade plan or recreate project

---

## Neon Free Tier Benefits

✅ Free PostgreSQL database
✅ 3 projects per account
✅ 10 GB storage
✅ Auto-suspend after inactivity
✅ Perfect for development

---

## Final Checklist

- [ ] Neon account created
- [ ] PostgreSQL project created on Neon
- [ ] Connection details extracted
- [ ] Code updated for database
- [ ] Environment variables added to Render
- [ ] Code pushed to GitHub
- [ ] Render Web Service created
- [ ] Deployment completed
- [ ] Health endpoint tested
- [ ] Frontend connected to deployed API

**Your live API will be at:**
```
https://stocktrack-api.onrender.com
```
