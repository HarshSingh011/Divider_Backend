package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"stocktrack-backend/internal/domain"
)

// JWTProvider implements TokenProvider interface
type JWTProvider struct {
	secretKey string
	expiryTime time.Duration
}

// NewJWTProvider creates a new JWT provider
func NewJWTProvider(secretKey string, expiryTime time.Duration) *JWTProvider {
	return &JWTProvider{
		secretKey: secretKey,
		expiryTime: expiryTime,
	}
}

// GenerateToken creates a new JWT token
func (jp *JWTProvider) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(jp.expiryTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jp.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken verifies and parses a JWT token
func (jp *JWTProvider) ValidateToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jp.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	userID, ok := (*claims)["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		return nil, errors.New("invalid email in token")
	}

	username, ok := (*claims)["username"].(string)
	if !ok {
		return nil, errors.New("invalid username in token")
	}

	exp, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid expiry in token")
	}

	return &domain.Claims{
		UserID:    userID,
		Email:     email,
		Username:  username,
		ExpiresAt: int64(exp),
	}, nil
}
