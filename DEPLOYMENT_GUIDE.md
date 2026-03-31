# StockTrack Deployment Guide

## Overview
This guide covers deploying StockTrack to Render with PostgreSQL database.

---

## STEP 1: Set Up PostgreSQL Locally (Recommended)

### A. Install PostgreSQL
**Windows:**
1. Download from https://www.postgresql.org/download/windows/
2. Run installer, remember the password
3. Port: 5432 (default)
4. Add to PATH during installation

**Verify Installation:**
```powershell
psql --version
```

### B. Create Local Database

```powershell
psql -U postgres
```

Enter password when prompted, then:

```sql
CREATE DATABASE stocktrack_local;
CREATE USER stocktrack_user WITH PASSWORD 'your_secure_password';
ALTER ROLE stocktrack_user SET client_encoding TO 'utf8';
ALTER ROLE stocktrack_user SET default_transaction_isolation TO 'read committed';
ALTER ROLE stocktrack_user SET default_transaction_deferrable TO on;
ALTER ROLE stocktrack_user SET timezone TO 'UTC';
GRANT ALL PRIVILEGES ON DATABASE stocktrack_local TO stocktrack_user;
ALTER DATABASE stocktrack_local OWNER TO stocktrack_user;
\q
```

### C. Update Backend Config

Create/Edit `config/config.go` or add environment variables:

```go
type DbConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// Load from environment
dbConfig := DbConfig{
    Host:     os.Getenv("DB_HOST"),      // localhost
    Port:     os.Getenv("DB_PORT"),      // 5432
    User:     os.Getenv("DB_USER"),      // stocktrack_user
    Password: os.Getenv("DB_PASSWORD"),  // your_secure_password
    DBName:   os.Getenv("DB_NAME"),      // stocktrack_local
    SSLMode:  os.Getenv("DB_SSLMODE"),   // disable (for local)
}
```

### D. Set Environment Variables (Windows PowerShell)

```powershell
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "stocktrack_user"
$env:DB_PASSWORD = "your_secure_password"
$env:DB_NAME = "stocktrack_local"
$env:DB_SSLMODE = "disable"
$env:SERVER_PORT = ":8080"
$env:JWT_SECRET = "your-super-secret-jwt-key-min-32-chars-here"
```

### E. Update Main.go to Use Database

Modify `cmd/server/container.go` to use PostgreSQL instead of in-memory:

```go
func NewContainer() *Container {
    cfg := config.NewDefaultConfig()

    db, err := db.NewDatabase(db.Config{
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
        SSLMode:  os.Getenv("DB_SSLMODE"),
    })
    if err != nil {
        panic(err)
    }

    userRepo := storage.NewPostgresUserRepository(db.GetConn())
    // Switch other repos to PostgreSQL too
    ...
}
```

### F. Test Locally

```powershell
cd c:\Users\harsh\Python\StockTrack
go build -o stocktrack.exe ./cmd/server
.\stocktrack.exe
```

Expected output: Server should connect and start without errors

---

## STEP 2: Deploy to Render

### A. Create Render Account
1. Go to https://render.com
2. Sign up with GitHub account (recommended)
3. Verify email

### B. Connect GitHub Repository

1. Push your code to GitHub:
```powershell
cd c:\Users\harsh\Python\StockTrack
git add .
git commit -m "Add database support and prepare for deployment"
git push
```

2. On Render dashboard: Click "New +" → "Web Service"
3. Connect your GitHub repo

### C. Configure Web Service

**Build Settings:**
- Environment: `Go`
- Build Command: `go build -o stocktrack ./cmd/server`
- Start Command: `./stocktrack`

**Environment Variables:**
Add in Render dashboard:
```
DB_HOST: (will get from PostgreSQL)
DB_PORT: 5432
DB_USER: postgres
DB_PASSWORD: (will generate)
DB_NAME: stocktrack
DB_SSLMODE: require
SERVER_PORT: 0.0.0.0:10000
JWT_SECRET: your-super-secret-jwt-key-min-32-chars-here
```

### D. Create PostgreSQL Database on Render

1. Render dashboard → "New +" → "PostgreSQL"
2. Name: `stocktrack-db`
3. Database: `stocktrack`
4. User: `postgres`
5. Generate password (copy it!)
6. Region: Select closest to you

**Get Connection Info:**
- Copy Host, Port, User, Password from Render PostgreSQL page
- Update environment variables in Web Service

### E. Deploy

1. Click "Deploy" on the Web Service
2. Wait for build (~2-3 minutes)
3. Check logs for errors
4. Once deployed, you'll get a URL like: `https://stocktrack.onrender.com`

---

## STEP 3: Connect Frontend to Deployed Backend

Update `app.html`:

```javascript
const API_BASE = "https://stocktrack.onrender.com";

// All fetch calls use API_BASE instead of localhost
```

---

## STEP 4: Monitor & Maintain

### View Logs
```
Render Dashboard → Select Web Service → Logs tab
```

### Database Connection
```
psql "postgresql://postgres:PASSWORD@HOST:5432/stocktrack"
```

### Auto-deploy on Push
- Render automatically redeploys when you push to GitHub
- To disable: Settings → Auto-deploy: Off

---

## Troubleshooting

**Error: "could not connect to database"**
- Check DB_HOST, DB_PORT, DB_USER, DB_PASSWORD are correct
- Verify PostgreSQL is running
- Check SSL mode (use `require` for cloud, `disable` for local)

**Error: "relation does not exist"**
- Database tables not created
- Run migrations or createSchema() function

**Error: "SSL certificate verification failed"**
- Use `DB_SSLMODE: require` for cloud
- Use `DB_SSLMODE: disable` for local

**Service crashes after deploy**
- Check logs in Render dashboard
- Verify all environment variables are set
- Test locally first before deploying

---

## Cost Estimates (as of 2024)

- **Render Web Service**: Free tier available ($0-7/month paid)
- **PostgreSQL**: ~$15/month for hobby tier
- **Total**: Free tier or ~$15-25/month

---

## Next Steps

1. ✅ Choose: Local testing first OR direct cloud deploy
2. ✅ Set up PostgreSQL (local or on Render)
3. ✅ Update environment variables
4. ✅ Test backend connection to database
5. ✅ Push to GitHub
6. ✅ Deploy on Render
7. ✅ Test via deployed API
8. ✅ Update frontend to use deployed backend URL
