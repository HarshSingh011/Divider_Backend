package http

import (
	"fmt"
	"net/http"
	"strings"

	"stocktrack-backend/internal/domain"
)

type AuthMiddleware struct {
	authService domain.AuthService
}

func NewAuthMiddleware(authService domain.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Email", claims.Email)
		r.Header.Set("X-Username", claims.Username)

		next(w, r)
	}
}

func (m *AuthMiddleware) ProtectWebSocket(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenParam := r.URL.Query().Get("token")
		if tokenParam == "" {
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

		claims, err := m.authService.ValidateToken(tokenParam)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Email", claims.Email)
		r.Header.Set("X-Username", claims.Username)

		next(w, r)
	}
}
