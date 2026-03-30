package http

import (
	"fmt"
	"net/http"
	"strings"

	"stocktrack-backend/internal/domain"
)

// AuthMiddleware validates JWT tokens from Authorization header
type AuthMiddleware struct {
	authService domain.AuthService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService domain.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Protect wraps an HTTP handler with JWT authentication
func (m *AuthMiddleware) Protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Parse "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// Store claims in request context for downstream handlers
		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Email", claims.Email)
		r.Header.Set("X-Username", claims.Username)

		next(w, r)
	}
}

// ProtectWebSocket validates JWT tokens for WebSocket upgrades
func (m *AuthMiddleware) ProtectWebSocket(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for token in query parameter (for WebSocket)
		tokenParam := r.URL.Query().Get("token")
		if tokenParam == "" {
			// Check Authorization header as fallback
			authHeader := r.Header.Get("Authorization")
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenParam = parts[1]
			}
		}

		if tokenParam == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := m.authService.ValidateToken(tokenParam)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		// Store claims in request context
		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Email", claims.Email)
		r.Header.Set("X-Username", claims.Username)

		next(w, r)
	}
}
