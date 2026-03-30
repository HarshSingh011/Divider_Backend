# StockTrack Backend - Real-Time Trading Engine

A production-grade, concurrent trading backend built with Go following **Hexagonal Architecture (Ports & Adapters)**. Features JWT authentication, WebSocket real-time market streaming, and clean separation of concerns.

## 🏗️ Architecture Overview

This project follows **Hexagonal Architecture** (also called Ports & Adapters), ensuring:
- **Core domain logic** is completely independent of external frameworks
- **Ports (Interfaces)** define contracts without implementation details
- **Adapters** implement those interfaces for specific technologies
- **Easy testing** with mock implementations
- **Framework agnostic** - swap databases, auth, or web servers without touching domain logic

### Layered Structure

```
internal/
├── domain/                    # 🎯 Core business logic (NO external dependencies)
│   ├── models.go             # Data structures (User, MarketTick, Claims)
│   ├── services.go           # Interfaces (ports) - what the app needs
│   ├── market_engine.go      # Market simulation logic
│   └── auth_service.go       # Authentication logic
│
├── port/                      # Interface definitions (abstraction)
│   └── [external dependencies defined here]
│
└── adapter/                   # 🔌 External implementations (adapters)
    ├── auth/
    │   └── jwt.go           # JWT token implementation
    ├── storage/
    │   └── memory.go        # In-memory user repository
    ├── http/
    │   ├── auth.go          # REST API handlers
    │   └── middleware.go    # JWT validation middleware
    └── ws/
        └── handler.go       # WebSocket broadcast manager

cmd/
├── server/
│   ├── main.go              # Entry point, route setup
│   └── container.go         # Dependency injection (wiring)
└── client/
    └── main_auth.go         # Test client with auth

config/
└── config.go                # Configuration management
```

## Features

✨ **Production-Ready**:
- Hexagonal Architecture for clean code
- JWT Authentication (register/login)
- Password hashing with bcrypt
- Thread-safe operations with sync.Mutex
- Real-time WebSocket streaming
- Comprehensive error handling

📊 **Real-Time Market Data**:
- 4 mock stocks with live price updates
- Updates every 500ms
- Percentage change calculations
- Timestamp tracking

🔒 **Security**:
- JWT token-based authentication
- Protected WebSocket endpoint
- Password hashing
- Bearer token validation

## Quick Start

### 1. Install Dependencies
```bash
cd C:\Users\harsh\Python\StockTrack
go mod download
```

### 2. Start the Server
```bash
cd cmd/server
go run .
```

Output:
```
=== StockTrack Backend Server ===
Running on http://localhost:8080

📋 API Endpoints:
  POST   http://localhost:8080/auth/register
  POST   http://localhost:8080/auth/login
  GET    http://localhost:8080/health
  WS     ws://localhost:8080/ws?token=<JWT>
```

### 3. Test with Authentication
```bash
cd cmd/client
go run main_auth.go
```

## API Endpoints

### Register New User
```bash
POST /auth/register
Content-Type: application/json

{
  "email": "trader@example.com",
  "username": "trader123",
  "password": "securepass123"
}
```

Response:
```json
{
  "id": "user_1711778844923456789",
  "email": "trader@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "email": "trader@example.com",
  "password": "securepass123"
}
```

Response: Same as register

### Health Check
```bash
GET /health
```

Response:
```json
{
  "status": "OK"
}
```

### WebSocket Subscribe (Protected)
```
ws://localhost:8080/ws?token=<JWT_TOKEN>
```

Message format:
```json
[
  {
    "symbol": "RELIANCE-CE-2900",
    "currentPrice": 45.50,
    "percentageChange": 2.34,
    "timestamp": "2026-03-30T12:00:00Z"
  }
]
```

## How It Works

### Dependency Injection (Container Pattern)

The `Container` in `cmd/server/container.go` wires all dependencies:

```go
container := NewContainer()
// ✓ Creates repositories
// ✓ Creates token providers
// ✓ Creates services
// ✓ Creates handlers
// ✓ Wires everything together
```

### Authentication Flow

1. **Registration**
   - Validate input
   - Check if user exists
   - Hash password with bcrypt
   - Save user to repository
   - Generate JWT token
   - Return token

2. **Login**
   - Find user by email
   - Verify password
   - Generate JWT token
   - Return token

3. **WebSocket Protection**
   - Extract token from query param or Authorization header
   - Validate JWT signature and expiry
   - Extract claims (user_id, email, username)
   - Allow connection if valid

### Market Engine Flow

1. **Market Engine** (runs every 500ms)
   - Generates random price changes
   - Calculates percentage changes
   - Sends update to broadcast channel

2. **Broadcasting Service**
   - Listens to market engine's broadcast channel
   - Sends data to all connected WebSocket clients
   - Handles client disconnections gracefully

3. **WebSocket Clients**
   - Send channel receives market data
   - Write to client connection
   - Automatic disconnect on error

## Hexagonal Architecture Benefits

### ✓ Testability
```go
// Easy to test - just mock the interfaces
type MockUserRepo struct { ... }
type MockTokenProvider struct { ... }

authService := domain.NewAuthService(mockRepo, mockProvider)
// Test without database or JWT library!
```

### ✓ Flexibility
- Swap in-memory storage for PostgreSQL (just implement UserRepository)
- Replace JWT with OAuth (just implement TokenProvider)
- Change WebSocket library (just reimplement Handler)

### ✓ Maintainability
- Domain logic has zero framework dependencies
- Clear separation: business logic vs infrastructure
- Easy to understand data flow

### ✓ Scalability
- Stateless services - easy to horizontally scale
- Pub/Sub pattern prevents bottlenecks
- Non-blocking broadcasts

## Production Checklist

- [ ] Change JWT secret key in `config/config.go`
- [ ] Use PostgreSQL/MySQL instead of in-memory storage
- [ ] Add proper logging (zerolog, zap)
- [ ] Add metrics (Prometheus)
- [ ] Add request validation
- [ ] Add rate limiting
- [ ] Enable HTTPS/TLS
- [ ] Add CORS middleware
- [ ] Add request tracing
- [ ] Add unit tests for domain layer
- [ ] Add integration tests

## Key Interview Points

When explaining this architecture:

1. **Separation of Concerns**
   - Domain logic doesn't know about HTTP, JWT, or databases
   - Each layer has a single responsibility

2. **SOLID Principles**
   - **S**ingle Responsibility: Each adapter does one thing
   - **O**pen/Closed: Open for extension, closed for modification
   - **L**iskov Substitution: Any implementation of port works
   - **I**nterface Segregation: Small, focused interfaces
   - **D**ependency Inversion: Depend on abstractions, not implementations

3. **Testing Strategy**
   ```go
   // No mocks needed - just swap implementations
   testAuthService := domain.NewAuthService(
       storage.NewInMemoryUserRepository(),
       &testTokenProvider{},
   )
   // Test pure business logic!
   ```

4. **Real-Time Capabilities**
   - Go channels for elegant producer-consumer
   - Goroutines for concurrent clients
   - Non-blocking broadcasts
   - Graceful connection handling

## Troubleshooting

**Build Error: "cannot find module"**
```bash
go mod download
go mod tidy
```

**WebSocket auth fails**
- Ensure token is passed: `ws://localhost:8080/ws?token=YOUR_TOKEN`
- Check token hasn't expired (24 hours default)
- Try logging in again to get fresh token

**Port 8080 in use**
```bash
netstat -ano | findstr :8080
```

## Dependencies

- `github.com/gorilla/websocket` - WebSocket support
- `github.com/golang-jwt/jwt/v5` - JWT tokens
- `golang.org/x/crypto` - Password hashing (bcrypt)

---

**Built with clean architecture principles for production-grade trading systems.** 🚀
