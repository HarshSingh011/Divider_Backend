package http

import (
	"encoding/json"
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Invalid authorization header format. Expected 'Bearer TOKEN', got %d parts", len(parts))})
			return
		}

		if parts[0] != "Bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Invalid authorization header format. Expected 'Bearer' prefix, got '%s'", parts[0])})
			return
		}

		token := parts[1]
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Token is empty"})
			return
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Invalid token: %v", err)})
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
