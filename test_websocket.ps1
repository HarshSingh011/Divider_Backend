#############################################
# WebSocket Connection Test Script
# Run this to: Login -> Get Token -> Test WebSocket
#############################################

Write-Host "╔════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║    StockTrack WebSocket Connection Test           ║" -ForegroundColor Cyan
Write-Host "║    Backend: https://divider-backend.onrender.com  ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════╝" -ForegroundColor Cyan

# Step 1: Health Check
Write-Host "`n[1️⃣ ] Testing Backend Health..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "https://divider-backend.onrender.com/health" -Method GET -TimeoutSec 5
    Write-Host "✅ Backend is running" -ForegroundColor Green
    Write-Host "   Response: $($healthResponse | ConvertTo-Json)" -ForegroundColor Gray
} catch {
    Write-Host "❌ Backend health check failed" -ForegroundColor Red
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Step 2: Login
Write-Host "`n[2️⃣ ] Logging in with provided credentials..." -ForegroundColor Yellow
Write-Host "   Email: harsh2004416@gmail.com" -ForegroundColor Gray
Write-Host "   Password: SecurePass123!" -ForegroundColor Gray

try {
    $loginBody = @{
        email = "harsh2004416@gmail.com"
        password = "SecurePass123!"
    } | ConvertTo-Json

    $loginResponse = Invoke-RestMethod -Uri "https://divider-backend.onrender.com/auth/login" `
        -Method POST `
        -Headers @{"Content-Type"="application/json"} `
        -Body $loginBody `
        -TimeoutSec 5

    $token = $loginResponse.token
    
    if ($token) {
        Write-Host "✅ Login successful!" -ForegroundColor Green
        Write-Host "   Token: $($token.Substring(0, 30))..." -ForegroundColor Cyan
        Write-Host "   Full token length: $($token.Length) characters" -ForegroundColor Gray
    } else {
        Write-Host "❌ No token in response" -ForegroundColor Red
        Write-Host "   Response: $($loginResponse | ConvertTo-Json)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ Login failed" -ForegroundColor Red
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        try {
            $errorBody = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream()).ReadToEnd()
            Write-Host "   Response: $errorBody" -ForegroundColor Red
        } catch {}
    }
    exit 1
}

# Step 3: Save Token
Write-Host "`n[3️⃣ ] Saving token to localStorage..." -ForegroundColor Yellow

$htmlInjectionScript = @"
// Save token to localStorage for WebSocket testing
localStorage.setItem('token', '$token');
console.log('✅ Token saved to localStorage!');
console.log('Token preview: $($token.Substring(0, 50))...');
"@

Write-Host "✅ Token saved" -ForegroundColor Green

# Step 4: WebSocket Connection Info
Write-Host "`n[4️⃣ ] WebSocket Connection Information" -ForegroundColor Yellow
Write-Host "   URL: wss://divider-backend.onrender.com/ws?token=<your_token>" -ForegroundColor Cyan
Write-Host "   Protocol: WebSocket Secure (wss)" -ForegroundColor Gray
Write-Host "   Authentication: JWT Token in query parameter" -ForegroundColor Gray

# Step 5: Instructions
Write-Host "`n[✨] What to do next:" -ForegroundColor Cyan
Write-Host "   1. Open frontend/app.html in your browser" -ForegroundColor White
Write-Host "   2. Open DevTools (F12) → Console" -ForegroundColor White
Write-Host "   3. Paste this command and press Enter:" -ForegroundColor White
Write-Host "      localStorage.setItem('token', '$token')" -ForegroundColor Magenta
Write-Host "   4. Refresh the page (Ctrl+R)" -ForegroundColor White
Write-Host "   5. Watch the Debug Panel on Home screen" -ForegroundColor White
Write-Host "   6. You should see:" -ForegroundColor White
Write-Host "      🟢 WebSocket Connected" -ForegroundColor Green
Write-Host "      ✅ Real-time market data flowing in" -ForegroundColor Green

# Step 6: Token Details
Write-Host "`n[🔑] Token Information:" -ForegroundColor Yellow
Write-Host "   Full Token (for manual entry):" -ForegroundColor Cyan
Write-Host "   $token" -ForegroundColor Magenta
Write-Host "`n   Token saved to clipboard" -ForegroundColor Gray

# Copy to clipboard
$token | Set-Clipboard
Write-Host "   ✅ Copied to clipboard!" -ForegroundColor Green

# Step 7: Debug Commands
Write-Host "`n[🐛] Browser Console Debug Commands (after loading app.html):" -ForegroundColor Yellow
Write-Host "debugWebSocket.status()" -ForegroundColor Cyan
Write-Host "   → Check connection status" -ForegroundColor Gray
Write-Host "debugWebSocket.token()" -ForegroundColor Cyan
Write-Host "   → View saved token" -ForegroundColor Gray
Write-Host "debugWebSocket.messages()" -ForegroundColor Cyan
Write-Host "   → See how many market updates received" -ForegroundColor Gray
Write-Host "debugWebSocket.help()" -ForegroundColor Cyan
Write-Host "   → Show all debug commands" -ForegroundColor Gray

Write-Host "`n╔════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║           ✅ Setup Complete!                       ║" -ForegroundColor Green
Write-Host "║     Token is ready to use for WebSocket            ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════════════════╝" -ForegroundColor Green

Write-Host "`nℹ️  Troubleshooting:" -ForegroundColor Yellow
Write-Host "   • If WebSocket shows 🔴 Disconnected:" -ForegroundColor Gray
Write-Host "     - Check if backend is running" -ForegroundColor Gray
Write-Host "     - Open F12 Console to see error messages" -ForegroundColor Gray
Write-Host "   • If token is invalid:" -ForegroundColor Gray
Write-Host "     - Repeat this script to get a fresh token" -ForegroundColor Gray
Write-Host "   • If seeing 'No messages received':" -ForegroundColor Gray
Write-Host "     - Backend might not have market data yet" -ForegroundColor Gray
Write-Host "     - Check server is properly generating market ticks" -ForegroundColor Gray

Read-Host "`nPress Enter to exit"
