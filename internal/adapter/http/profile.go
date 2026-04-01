package http

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"time"

	"stocktrack-backend/internal/domain"
)

type ProfileHandler struct {
	authService domain.AuthService
}

func NewProfileHandler(authService domain.AuthService) *ProfileHandler {
	return &ProfileHandler{
		authService: authService,
	}
}

type UserProfile struct {
	ID                 string    `json:"id"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	BankAccount        string    `json:"bank_account"`
	BankAccountStatus  string    `json:"bank_account_status"`
	MemberSince        time.Time `json:"member_since"`
	IsVerified         bool      `json:"is_verified"`
	Theme              string    `json:"theme"`
	NotificationAlerts bool      `json:"notification_alerts"`
	NotificationTrades bool      `json:"notification_trades"`
	NotificationNews   bool      `json:"notification_news"`
	TwoFactorEnabled   bool      `json:"two_factor_enabled"`
}

func generateFakeBankAccount(userID string) string {
	h := fnv.New32a()
	h.Write([]byte(userID))
	accountNum := h.Sum32()

	lastFour := (accountNum % 10000)

	return fmt.Sprintf("HDFC Bank - Savings | ***%04d", lastFour)
}

func generateFakePhoneNumber(userID string) string {
	h := fnv.New32a()
	h.Write([]byte(userID + "phone"))
	phoneNum := h.Sum32()

	part1 := (phoneNum / 100000) % 100000
	part2 := phoneNum % 10000

	return fmt.Sprintf("+91 %d %05d", 98000+part1%1000, part2)
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	username := r.Header.Get("X-Username")
	email := r.Header.Get("X-Email")

	if email == "" {
		email = username + "@stocktrack.com"
	}

	profile := UserProfile{
		ID:                 userID,
		Username:           username,
		Email:              email,
		Phone:              generateFakePhoneNumber(userID),
		BankAccount:        generateFakeBankAccount(userID),
		BankAccountStatus:  "Verified",
		MemberSince:        time.Now().AddDate(-1, 0, 0),
		IsVerified:         true,
		Theme:              "Light",
		NotificationAlerts: true,
		NotificationTrades: true,
		NotificationNews:   false,
		TwoFactorEnabled:   true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

type UpdateProfileRequest struct {
	Theme              string `json:"theme"`
	NotificationAlerts bool   `json:"notification_alerts"`
	NotificationTrades bool   `json:"notification_trades"`
	NotificationNews   bool   `json:"notification_news"`
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
		"data": map[string]interface{}{
			"theme":               req.Theme,
			"notification_alerts": req.NotificationAlerts,
			"notification_trades": req.NotificationTrades,
			"notification_news":   req.NotificationNews,
			"two_factor_enabled":  true,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type ActiveSession struct {
	DeviceName string    `json:"device_name"`
	LastActive time.Time `json:"last_active"`
	IPAddress  string    `json:"ip_address"`
}

func (h *ProfileHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	sessions := []ActiveSession{
		{
			DeviceName: "iPhone 12 - This Device",
			LastActive: time.Now(),
			IPAddress:  "192.168.1.100",
		},
		{
			DeviceName: "Chrome Browser - Desktop",
			LastActive: time.Now().Add(-24 * time.Hour),
			IPAddress:  "192.168.1.101",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

func (h *ProfileHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}
