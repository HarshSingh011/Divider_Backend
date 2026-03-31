# Render + Neon Deployment Checklist

## Status: ✅ BACKEND READY TO DEPLOY

Your StockTrack backend has been compiled successfully and is ready for production deployment.

---

## Quick Deployment Steps

### Step 1: Create Neon PostgreSQL Project (2 minutes)
- [ ] Go to https://neon.tech
- [ ] Sign up with GitHub
- [ ] Create project named `stocktrack`
- [ ] Note down: **Host**, **User**, **Password**, **Database name**

**Important**: Keep SSL mode as `require` for Neon

---

### Step 2: Push Code to GitHub (1 minute)
```powershell
cd c:\Users\harsh\Python\StockTrack
git add .
git commit -m "Deploy to Render with Neon PostgreSQL"
git push
```

---

### Step 3: Deploy on Render (3 minutes)

1. Go to https://render.com
2. Sign up with GitHub (connect repo)
3. Click "New +" → "Web Service"
4. Select your StockTrack repository
5. Configure:
   - **Name**: `stocktrack-api`
   - **Environment**: `Go`
   - **Build Command**: `go build -o stocktrack ./cmd/server`
   - **Start Command**: `./stocktrack`
   - **Instance Type**: Free

---

### Step 4: Add Environment Variables (2 minutes)

Copy these into Render's "Environment" section:

```
DB_HOST=<neon-host>
DB_PORT=5432
DB_USER=<neon-user>
DB_PASSWORD=<neon-password>
DB_NAME=stocktrack_db
DB_SSLMODE=require
SERVER_PORT=0.0.0.0:10000
JWT_SECRET=your-super-secret-key-min-32-chars-here
```

Example:
```
DB_HOST=ep-xyz.us-east-1.neon.tech
DB_PORT=5432
DB_USER=neondb_owner
DB_PASSWORD=abc123xyz456
DB_NAME=stocktrack_db
DB_SSLMODE=require
SERVER_PORT=0.0.0.0:10000
JWT_SECRET=super-secret-jwt-key-min-32-characters-long-here
```

---

### Step 5: Deploy & Test (5 minutes)

1. Click "Deploy" on Render
2. Wait 3-5 minutes for build
3. Check logs for errors
4. You'll get a URL like: `https://stocktrack-api.onrender.com`

---

### Step 6: Test Your API

**Health Check:**
```powershell
curl https://stocktrack-api.onrender.com/health
```

Expected response:
```json
{"status":"OK"}
```

**Register User:**
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

### Step 7: Update Frontend (Optional)

Edit your `app.html`:

```javascript
// Find this line:
const API_BASE = "http://localhost:8080";

// Change to:
const API_BASE = "https://stocktrack-api.onrender.com";
```

---

## Environment Variables Explanation

| Variable | Example | Purpose |
|----------|---------|---------|
| `DB_HOST` | `ep-xyz.us-east-1.neon.tech` | Neon server location |
| `DB_PORT` | `5432` | Standard PostgreSQL port |
| `DB_USER` | `neondb_owner` | Neon database user |
| `DB_PASSWORD` | `abc123xyz456` | Neon password |
| `DB_NAME` | `stocktrack_db` | Database name |
| `DB_SSLMODE` | `require` | Use SSL for cloud DB |
| `SERVER_PORT` | `0.0.0.0:10000` | Render assigns port dynamically |
| `JWT_SECRET` | `any-32+-char-secret` | Token signing key |

---

## Troubleshooting

**"Database connection refused"**
- Verify all DB_* vars match Neon project info
- Check DB_SSLMODE=require (not disable)
- Ensure Neon project is active (not suspended)

**"Service crashes on startup"**
- Check Render logs for specific error
- Verify JWT_SECRET is at least 32 characters
- Try health endpoint: `/health`

**"Build failed"**
- Ensure code was pushed to GitHub
- Check if `go.mod` has all dependencies
- Render logs show full error

**"504 Gateway Timeout"**
- Render free tier may be cold-starting
- Try again after 30 seconds
- First request takes 30-60 seconds

---

## Cost Breakdown

| Service | Tier | Cost |
|---------|------|------|
| Render Web Service | Free | $0/month |
| Neon PostgreSQL | Hobby | Free (shared compute) |
| **Total** | | **$0/month** |

Upgrade to paid plans only if you exceed free tier limits.

---

## Next Steps After Deployment

1. ✅ Backend deployed and running
2. ✅ PostgreSQL database connected
3. **→ Build React Native app** (if needed)
4. **→ Add more features** (2FA, email notifications, etc.)
5. **→ Scale to paid tier** (when traffic increases)

---

## Your Deployed API

Once deployed, your API will be at:
```
https://stocktrack-api.onrender.com
```

All endpoints from your local machine (`http://localhost:8080`) will work at this URL!

---

**Questions?** Refer to [RENDER_NEON_DEPLOYMENT.md](./RENDER_NEON_DEPLOYMENT.md) for detailed instructions.
