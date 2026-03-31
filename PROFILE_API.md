# Profile API Implementation - March 30, 2026

## New Endpoints Created

### 1. GET /user/profile
**Protected:** Yes (requires JWT token)

**Response:**
```json
{
  "id": "user123",
  "username": "harsh",
  "email": "harsh@example.com",
  "phone": "+91 98765 43210",
  "bank_account": "HDFC Bank - Savings | ***1234",
  "bank_account_status": "Verified",
  "member_since": "2025-03-30T00:00:00Z",
  "is_verified": true,
  "theme": "Light",
  "notification_alerts": true,
  "notification_trades": true,
  "notification_news": false,
  "two_factor_enabled": true
}
```

### 2. PUT /user/profile
**Protected:** Yes (requires JWT token)

**Request:**
```json
{
  "theme": "Dark",
  "notification_alerts": true,
  "notification_trades": true,
  "notification_news": false
}
```

**Response:**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "theme": "Dark",
    "notification_alerts": true,
    "notification_trades": true,
    "notification_news": false,
    "two_factor_enabled": true
  }
}
```

### 3. GET /user/sessions
**Protected:** Yes (requires JWT token)

**Response:**
```json
[
  {
    "device_name": "iPhone 12 - This Device",
    "last_active": "2026-03-30T10:30:00Z",
    "ip_address": "192.168.1.100"
  },
  {
    "device_name": "Chrome Browser - Desktop",
    "last_active": "2026-03-29T10:30:00Z",
    "ip_address": "192.168.1.101"
  }
]
```

### 4. POST /user/logout
**Protected:** Yes (requires JWT token)

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

## Key Features

### Fake Bank Account Generation
- **Deterministic**: Same user always gets the same fake bank account (based on user ID hash)
- **Format**: `HDFC Bank - Savings | ***{last 4 digits}`
- **Safe**: No real account numbers used
- **Example**: `HDFC Bank - Savings | ***5234`

### Fake Phone Number Generation
- **Deterministic**: Same user always gets the same fake phone
- **Format**: `+91 {area code} {number}`
- **Example**: `+91 98765 43210`

### Auto-Generated Profile Data
- **Phone & Bank Account**: Generated on-the-fly based on user ID
- **Member Since**: Set to 1 year before current date
- **Verified Status**: Always true
- **Default Settings**: Theme=Light, Alerts=On, Trades=On, News=Off, 2FA=On

## Files Modified

1. **`internal/adapter/http/profile.go`** (NEW)
   - ProfileHandler struct
   - GetProfile, UpdateProfile, GetSessions, Logout handlers
   - Fake data generators: generateFakeBankAccount(), generateFakePhoneNumber()

2. **`cmd/server/container.go`**
   - Added ProfileHandler field to Container struct
   - Initialize ProfileHandler in NewContainer()

3. **`cmd/server/main.go`**
   - Added routes for `/user/profile`, `/user/sessions`, `/user/logout`

4. **`frontend/app.html`**
   - Added loadUserProfile() function to fetch data from backend
   - Added updateProfileUI() function to display real data on the screen
   - Integrated profile load when user switches to Profile tab

## How It Works

1. **User registers/logs in** → Gets JWT token and stores in localStorage
2. **User clicks Profile tab** → Frontend calls loadUserProfile()
3. **loadUserProfile()** → Makes GET request to `/user/profile` with Bearer token
4. **Backend generates data** → Creates consistent fake bank account, phone, etc.
5. **Frontend updates UI** → Displays real user data from backend response

## Backend Data Flow

```
User Credentials → JWT Token → /user/profile endpoint
                    ↓
              X-User-ID header extracted
                    ↓
            Generate consistent fake data
            (bank account, phone, sessions)
                    ↓
            Return UserProfile as JSON
                    ↓
            Frontend updates UI with real data
```

## Security Notes

- **All profile endpoints protected** with JWT authentication
- **No sensitive data exposed** (only fake bank account last 4 digits)
- **Consistent generation** ensures same user always gets same fake data
- **No database persistence needed yet** (can be added later)

## Testing

To test the profile API:

```bash
# 1. Start the server
go run ./cmd/server

# 2. Register a user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"harsh@example.com","username":"harsh","password":"password123"}'

# 3. Extract the token from response

# 4. Get profile
curl -X GET http://localhost:8080/user/profile \
  -H "Authorization: Bearer {token}"
```

## Next Steps (Optional)

- [ ] Add database persistence for profile settings
- [ ] Implement real 2FA functionality
- [ ] Add password change endpoint
- [ ] Implement session management with logout for specific devices
- [ ] Add email verification flow
- [ ] Implement profile picture upload
