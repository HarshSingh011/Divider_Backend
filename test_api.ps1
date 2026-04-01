$baseURL = "https://divider-backend.onrender.com"

Write-Host "`n╔════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║         StockTrack API - Deployment Testing              ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════╝`n" -ForegroundColor Cyan

# Test 1: Health Check
Write-Host "🔍 Test 1: Health Check Endpoint" -ForegroundColor Yellow
Write-Host "URL: GET $baseURL/health" -ForegroundColor Gray
try {
    $response = Invoke-WebRequest -Uri "$baseURL/health" -Method GET -TimeoutSec 10
    Write-Host "✅ Status: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "Response: $($response.Content)" -ForegroundColor Green
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 2: Register User
Write-Host "🔍 Test 2: Register User" -ForegroundColor Yellow
Write-Host "URL: POST $baseURL/auth/register" -ForegroundColor Gray
$registerBody = @{
    email = "john.doe@example.com"
    username = "john_doe"
    password = "SecurePass123!"
} | ConvertTo-Json

Write-Host "Request: $registerBody" -ForegroundColor Gray
try {
    $response = Invoke-WebRequest -Uri "$baseURL/auth/register" -Method POST `
        -Headers @{"Content-Type"="application/json"} `
        -Body $registerBody -TimeoutSec 10
    Write-Host "✅ Status: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "Response: $($response.Content)" -ForegroundColor Green
    $content = $response.Content | ConvertFrom-Json
    $script:token = $content.token
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 3: Login User
Write-Host "🔍 Test 3: Login User" -ForegroundColor Yellow
Write-Host "URL: POST $baseURL/auth/login" -ForegroundColor Gray
$loginBody = @{
    email = "john.doe@example.com"
    password = "SecurePass123!"
} | ConvertTo-Json

Write-Host "Request: $loginBody" -ForegroundColor Gray
try {
    $response = Invoke-WebRequest -Uri "$baseURL/auth/login" -Method POST `
        -Headers @{"Content-Type"="application/json"} `
        -Body $loginBody -TimeoutSec 10
    Write-Host "✅ Status: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "Response: $($response.Content)" -ForegroundColor Green
    $content = $response.Content | ConvertFrom-Json
    $script:token = $content.token
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 4: Get Profile (with auth token)
if ($script:token) {
    Write-Host "🔍 Test 4: Get User Profile (Protected Route)" -ForegroundColor Yellow
    Write-Host "URL: GET $baseURL/user/profile" -ForegroundColor Gray
    Write-Host "Auth: Bearer $($script:token.Substring(0, 20))..." -ForegroundColor Gray
    try {
        $headers = @{
            "Content-Type" = "application/json"
            "Authorization" = "Bearer $($script:token)"
        }
        $response = Invoke-WebRequest -Uri "$baseURL/user/profile" -Method GET `
            -Headers $headers -TimeoutSec 10
        Write-Host "✅ Status: $($response.StatusCode)" -ForegroundColor Green
        Write-Host "Response: $($response.Content)" -ForegroundColor Green
    } catch {
        Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    }
    Write-Host ""
}

# Test 5: Invalid Email (Regex Validation)
Write-Host "🔍 Test 5: Invalid Email (Testing Regex Validation)" -ForegroundColor Yellow
Write-Host "URL: POST $baseURL/auth/register" -ForegroundColor Gray
$invalidEmailBody = @{
    email = "invalid-email"
    username = "test_user"
    password = "SecurePass123!"
} | ConvertTo-Json

Write-Host "Request with invalid email: $invalidEmailBody" -ForegroundColor Gray
try {
    $response = Invoke-WebRequest -Uri "$baseURL/auth/register" -Method POST `
        -Headers @{"Content-Type"="application/json"} `
        -Body $invalidEmailBody -TimeoutSec 10
    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Yellow
    Write-Host "Response: $($response.Content)" -ForegroundColor Yellow
} catch {
    Write-Host "✅ Validation Working! Error: $($_.Exception.Response.StatusCode)" -ForegroundColor Green
    Write-Host "Message: $($_.Exception.Message)" -ForegroundColor Green
}
Write-Host ""

# Test 6: Invalid Password (Regex Validation)
Write-Host "🔍 Test 6: Invalid Password (Testing Regex Validation)" -ForegroundColor Yellow
Write-Host "URL: POST $baseURL/auth/register" -ForegroundColor Gray
$invalidPasswordBody = @{
    email = "valid@example.com"
    username = "test_user"
    password = "weak"
} | ConvertTo-Json

Write-Host "Request with weak password: $invalidPasswordBody" -ForegroundColor Gray
try {
    $response = Invoke-WebRequest -Uri "$baseURL/auth/register" -Method POST `
        -Headers @{"Content-Type"="application/json"} `
        -Body $invalidPasswordBody -TimeoutSec 10
    Write-Host "Status: $($response.StatusCode)" -ForegroundColor Yellow
    Write-Host "Response: $($response.Content)" -ForegroundColor Yellow
} catch {
    Write-Host "✅ Validation Working! Error: $($_.Exception.Response.StatusCode)" -ForegroundColor Green
    Write-Host "Message: $($_.Exception.Message)" -ForegroundColor Green
}
Write-Host ""

Write-Host "╔════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║              API Testing Complete!                        ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════════╝`n" -ForegroundColor Cyan
