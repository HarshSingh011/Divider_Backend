package http

import (
	"encoding/json"
	"net/http"

	"stocktrack-backend/internal/domain"
)

type DashboardHandler struct {
	walletService domain.WalletService
}

func NewDashboardHandler(walletService domain.WalletService) *DashboardHandler {
	return &DashboardHandler{
		walletService: walletService,
	}
}

// HomeScreenResponse contains only the data needed for home screen UI
type HomeScreenResponse struct {
	Success bool       `json:"success"`
	Data    HomeScreen `json:"data"`
}

type HomeScreen struct {
	TotalBalance    float64 `json:"total_balance"`     // Total portfolio value
	AvailableCash   float64 `json:"available_cash"`    // Cash ready to trade
	InvestedAmount  float64 `json:"invested_amount"`   // Money in stocks
	TotalPnL        float64 `json:"total_pnl"`         // Total profit/loss
	TotalPnLPercent float64 `json:"total_pnl_percent"` // Profit/loss percentage
	HoldingCount    int     `json:"holding_count"`     // Number of stocks held
	TopGainer       string  `json:"top_gainer"`        // Best performing stock
	TopLoser        string  `json:"top_loser"`         // Worst performing stock
	LastUpdated     string  `json:"last_updated"`      // Timestamp
}

// GetHomeScreen returns minimized home screen data
func (h *DashboardHandler) GetHomeScreen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "UNAUTHORIZED",
				"message": "User not authenticated. X-User-ID header is required",
			},
		})
		return
	}

	// Get full wallet snapshot
	snapshot, err := h.walletService.GetWalletSnapshot(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]string{
				"code":    "WALLET_ERROR",
				"message": "Failed to get wallet data",
			},
		})
		return
	}

	// Calculate total P&L from positions
	totalPnL := 0.0
	holdingCount := 0
	topGainer := ""
	topLoserName := ""
	maxGain := 0.0
	maxLoss := 0.0

	for symbol, position := range snapshot.Positions {
		totalPnL += position.UnrealizedPnL
		holdingCount++

		// Find top gainer
		if position.UnrealizedPnL > maxGain {
			maxGain = position.UnrealizedPnL
			topGainer = symbol
		}

		// Find top loser
		if position.UnrealizedPnL < maxLoss {
			maxLoss = position.UnrealizedPnL
			topLoserName = symbol
		}
	}

	// Calculate P&L percentage
	totalPnLPercent := 0.0
	if snapshot.InvestedAmount > 0 {
		totalPnLPercent = (totalPnL / snapshot.InvestedAmount) * 100
	}

	// Build home screen response with only necessary data
	homeScreen := HomeScreen{
		TotalBalance:    snapshot.TotalBalance,
		AvailableCash:   snapshot.AvailableCash,
		InvestedAmount:  snapshot.InvestedAmount,
		TotalPnL:        totalPnL,
		TotalPnLPercent: totalPnLPercent,
		HoldingCount:    holdingCount,
		TopGainer:       topGainer,
		TopLoser:        topLoserName,
		LastUpdated:     snapshot.LastUpdated.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HomeScreenResponse{
		Success: true,
		Data:    homeScreen,
	})
}
