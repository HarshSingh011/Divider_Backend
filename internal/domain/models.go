package domain

import "time"

// MarketTick represents a single price update for a stock or option
type MarketTick struct {
	Symbol           string    `json:"symbol"`
	CurrentPrice     float64   `json:"currentPrice"`
	PercentageChange float64   `json:"percentageChange"`
	Timestamp        time.Time `json:"timestamp"`
}

// User represents a trading platform user
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never expose password in JSON
	CreatedAt time.Time `json:"created_at"`
}

// AuthRequest is the input for login/register
type AuthRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse is the output after successful auth
type AuthResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// Claims represents JWT token claims
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	ExpiresAt int64 `json:"exp"`
}
