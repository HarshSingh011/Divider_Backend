package domain

import "time"

type MarketTick struct {
	Symbol           string    `json:"symbol"`
	CurrentPrice     float64   `json:"currentPrice"`
	PercentageChange float64   `json:"percentageChange"`
	Timestamp        time.Time `json:"timestamp"`
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	ExpiresAt int64  `json:"exp"`
}
