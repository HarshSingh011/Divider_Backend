# PostgreSQL Setup Guide

This guide covers setting up PostgreSQL for StockTrack backend, both locally and using Neon (serverless).

## Option 1: Local PostgreSQL Setup

### 1. Install PostgreSQL

**Windows:**
- Download from https://www.postgresql.org/download/windows/
- Or use Chocolatey: `choco install postgresql`

**macOS:**
```bash
brew install postgresql
brew services start postgresql
```

**Linux (Ubuntu):**
```bash
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### 2. Create Database and User

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE stocktrack;

# Create user
CREATE USER stocktrack_user WITH PASSWORD 'securepassword123';

# Grant permissions
ALTER ROLE stocktrack_user SET client_encoding TO 'utf8';
ALTER ROLE stocktrack_user SET default_transaction_isolation TO 'read committed';
ALTER ROLE stocktrack_user SET default_transaction_deferrable TO on;
ALTER ROLE stocktrack_user SET default_transaction_read_only TO off;
GRANT ALL PRIVILEGES ON DATABASE stocktrack TO stocktrack_user;

# Exit psql
\q
```

### 3. Configure Environment

Create `.env` file in project root:
```
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=stocktrack_user
DB_PASSWORD=securepassword123
DB_NAME=stocktrack
DB_SSLMODE=disable
```

### 4. Start Application

```bash
cd cmd/server
go run .
```

The app will automatically create tables on startup.

---

## Option 2: Neon (Serverless PostgreSQL) - RECOMMENDED FOR PRODUCTION

Neon is perfect for cloud deployments! Free tier includes 5GB storage.

### 1. Create Neon Account & Project

1. Go to https://neon.tech/
2. Sign up with GitHub or email
3. Create a new project
4. Choose region closest to users

### 2. Get Connection String

After creating project, Neon shows a connection string like:
```
postgresql://neon_user:password@ep-silent-moon-12345.neon.tech/neon_db?sslmode=require
```

### 3. Configure Environment

Create `.env` file:
```
DB_DRIVER=postgres
DB_HOST=ep-silent-moon-12345.neon.tech
DB_PORT=5432
DB_USER=neon_user
DB_PASSWORD=your_password
DB_NAME=neon_db
DB_SSLMODE=require
```

Or use the full connection string approach (alternative):
```bash
# Parse Neon connection string
postgresql://<USER>:<PASSWORD>@<HOST>/dbname?sslmode=require
```

### 4. Test Connection

```bash
cd cmd/server
go run .
```

Watch for "✓ PostgreSQL connected" message.

---

## Database Schema

The application automatically creates the schema on startup:

```sql
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

No manual migrations needed - tables are created automatically!

---

## Local Testing

### 1. Connect to Local Database

```bash
psql -U stocktrack_user -d stocktrack

# List tables
\dt

# View users
SELECT * FROM users;
```

### 2. Test with API

Register a user:
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "trader@example.com",
    "username": "trader123",
    "password": "pass123"
  }'
```

Check database:
```bash
psql -U stocktrack_user -d stocktrack -c "SELECT email, username, created_at FROM users;"
```

---

## Neon Tips & Best Practices

### Connection Pool
- Neon supports up to 100 concurrent connections on free tier
- Each Go application can have up to 25 open connections (configured in code)
- Multiple app instances should share connection pool

### Cost Optimization
- **Free tier**: 5GB storage, 100 connections per month limit
- **Scale with Pro tier**: Pay-as-you-go, unlimited connections
- Storage is billed per GB-hour

### Development Workflow

```bash
# Local development
DB_DRIVER=memory go run cmd/server/main.go

# Testing with real database
DB_DRIVER=postgres \
DB_HOST=localhost \
DB_USER=stocktrack_user \
DB_PASSWORD=securepassword123 \
go run cmd/server/main.go

# Production (Neon)
DB_DRIVER=postgres \
DB_HOST=ep-xxxx.neon.tech \
DB_USER=neon_user \
DB_PASSWORD=xxxx \
DB_SSLMODE=require \
go run cmd/server/main.go
```

---

## Troubleshooting

### "Failed to initialize PostgreSQL: role does not exist"
```bash
# Create user with proper permissions
psql -U postgres -c "CREATE ROLE stocktrack_user WITH LOGIN PASSWORD 'pass';"
psql -U postgres -c "CREATE DATABASE stocktrack OWNER stocktrack_user;"
```

### "SSL verification failed" (Neon connection)
Ensure `DB_SSLMODE=require` in environment. Neon requires SSL.

### "dial tcp: connection refused"
- Check if PostgreSQL is running locally: `psql --version`
- For Neon: Verify internet connection and connection string

### "user already exists"
```bash
psql -U postgres -c "DROP ROLE IF EXISTS stocktrack_user;"
# Then recreate
```

### Test Connection

```bash
# Local
psql -U stocktrack_user -h localhost -d stocktrack -c "SELECT 1;"

# Neon
psql postgresql://user:pass@host/dbname -c "SELECT 1;"
```

---

## Switching Storage Backends

### In-Memory (Development)
```bash
DB_DRIVER=memory go run cmd/server/main.go
# Data lost on restart - perfect for quick testing
```

### PostgreSQL (Production)
```bash
DB_DRIVER=postgres \
DB_HOST=your.host \
DB_USER=user \
DB_PASSWORD=pass \
DB_NAME=dbname \
go run cmd/server/main.go
# Data persists
```

### Fallback Logic
If PostgreSQL fails to connect, app automatically falls back to in-memory storage. This prevents crashes during database issues.

---

## Performance Tips

1. **Connection Pool Tuning**
   - MaxOpenConns: 25 (configured in code)
   - MaxIdleConns: 5
   - ConnMaxLifetime: 5 minutes

2. **Query Optimization**
   - Indexes on `email` and `username` for fast lookups
   - Parameterized queries prevent SQL injection

3. **Neon Specifics**
   - Use Neon's serverless endpoints for scaling
   - Monitor CPU/storage in Neon console
   - Set up auto-scaling on Pro plan

---

## Next Steps

Once PostgreSQL is working:

1. ✅ Replace in-memory user repository with PostgreSQL
2. ⏭️ Add audit logging table
3. ⏭️ Add portfolio/positions tracking
4. ⏭️ Add trade history persistence
5. ⏭️ Add analytics/reports

All powered by Neon! 🚀
