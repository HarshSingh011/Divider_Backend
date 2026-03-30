# StockTrack Backend - Real-Time Trading Engine

A high-performance, concurrent trading backend built with Go. Uses Pub/Sub pattern with WebSockets to stream live market data to multiple clients.

## Features

- **Real-Time Market Data Stream** - Price updates every 500ms
- **Concurrent WebSocket Support** - Handles thousands of simultaneous connections
- **Pub/Sub Architecture** - Clean separation between market engine and network delivery
- **Race Condition Protection** - Uses `sync.Mutex` for thread-safe operations
- **Mock Trading Data** - Simulates real market behavior with random price changes

## Project Structure

```
.
├── main.go                 # WebSocket server, Hub, Market Engine, Broadcasting logic
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── cmd/
│   └── client/
│       └── main.go         # WebSocket test client
└── README.md               # This file
```

## Tech Stack

- **Language**: Go 1.19+
- **WebSocket**: [gorilla/websocket](https://github.com/gorilla/websocket)
- **Concurrency**: Go Goroutines and Channels

## Quick Start

### 1. **Start the Server**

```bash
cd C:\Users\harsh\Python\StockTrack
go run main.go
```

You should see:
```
=== Market Server running on :8080 ===
WebSocket endpoint: ws://localhost:8080/ws
Health check: http://localhost:8080/health

Waiting for connections...
```

### 2. **Connect a Test Client** (in a new terminal)

```bash
cd C:\Users\harsh\Python\StockTrack\cmd\client
go run main.go
```

The client will display live market ticks for 10 seconds.

## API Endpoints

### WebSocket
- **URL**: `ws://localhost:8080/ws`
- **Data Format**: JSON array of MarketTick objects
- **Update Frequency**: Every 500ms

### HTTP
- **Health Check**: `http://localhost:8080/health`
- **Response**: JSON status and server time

## Market Data Format

Each market tick contains:
```json
{
  "symbol": "RELIANCE-CE-2900",
  "currentPrice": 45.50,
  "percentageChange": 2.34,
  "timestamp": "2026-03-30T12:00:00Z"
}
```

## How It Works

### Architecture Diagram

1. **Market Engine** (Goroutine)
   - Runs infinite loop every 500ms
   - Generates random price changes for mock stocks
   - Sends updated data to broadcast channel

2. **Broadcasting Service** (Goroutine)
   - Listens to broadcast channel
   - Manages client connections/disconnections
   - Sends data to all connected clients simultaneously

3. **WebSocket Handler** (HTTP Handler)
   - Upgrades HTTP connection to WebSocket
   - Registers/unregisters clients with the Hub
   - Detects client disconnections

### Concurrency Model

- **Market Engine**: 1 Goroutine generating prices
- **Broadcaster**: 1 Goroutine distributing to all clients
- **Client Listeners**: 1 Goroutine per connected client
- **Synchronization**: `sync.Mutex` protects shared `clients` map

This ensures even if a client has slow connection, the market engine keeps running at full speed.

## Scalability

- Tested with simultaneous connections
- Handles thousands of concurrent WebSocket clients
- Non-blocking broadcasting using Go channels
- Efficient memory usage with Goroutines

## Next Steps: React Native Integration

Once you connect your React Native app, it will:
1. Connect to `ws://localhost:8080/ws` (or your production URL)
2. Receive real-time market updates
3. Render price charts with live data

## Build & Deploy

### Local Build
```bash
go build -o stocktrack-backend.exe
```

### Run Compiled Binary
```bash
./stocktrack-backend.exe
```

### Docker Deployment (Optional)
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o app main.go
EXPOSE 8080
CMD ["./app"]
```

## Troubleshooting

**Port 8080 already in use?**
```bash
netstat -ano | findstr :8080  # Check what's using the port
```

**WebSocket connection refused?**
- Ensure server is running (`go run main.go`)
- Check firewall isn't blocking port 8080
- Verify client URL: `ws://localhost:8080/ws` (not `http://`)

**Build errors?**
```bash
go mod tidy       # Clean up dependencies
go get ./...      # Re-download modules
go build ./...    # Rebuild
```

## Interview Points

✨ **Key Takeaways to Mention**:

1. **Separation of Concerns**: Market engine runs independently from network layer
2. **Race Conditions**: `sync.Mutex` prevents crashes from concurrent access
3. **Goroutines**: Lightweight concurrency without thread overhead
4. **Channels**: Elegant producer-consumer pattern between engine and broadcaster
5. **Non-blocking Broadcasts**: All clients receive data simultaneously
6. **Error Handling**: Gracefully removes disconnected clients

---

**Ready to connect React Native? Let me know! 🚀**
