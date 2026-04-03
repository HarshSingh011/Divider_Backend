package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stocktrack-backend/internal/domain"
)

type TradingHandler struct {
	walletService         domain.WalletService
	alertService          domain.AlertService
	ohlcService           domain.OHLCService
	transactionRepository domain.TransactionRepository
}

func NewTradingHandler(
	walletService domain.WalletService,
	alertService domain.AlertService,
	ohlcService domain.OHLCService,
	transactionRepository domain.TransactionRepository,
) *TradingHandler {
	return &TradingHandler{
		walletService:         walletService,
		alertService:          alertService,
		ohlcService:           ohlcService,
		transactionRepository: transactionRepository,
	}
}

func (h *TradingHandler) ExecuteTrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		Symbol   string  `json:"symbol"`
		Quantity float64 `json:"quantity"`
		Price    float64 `json:"price"`
		Type     string  `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.walletService.ExecuteTrade(userID, req.Symbol, req.Quantity, req.Price, req.Type); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	total := req.Quantity * req.Price
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "Trade executed successfully",
		"total":    total,
		"symbol":   req.Symbol,
		"quantity": req.Quantity,
		"price":    req.Price,
		"type":     req.Type,
	})
}

func (h *TradingHandler) GetWalletSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	snapshot, err := h.walletService.GetWalletSnapshot(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

func (h *TradingHandler) DepositCash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.walletService.DepositCash(userID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Deposit successful",
	})
}

func (h *TradingHandler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		Symbol    string  `json:"symbol"`
		Price     float64 `json:"price"`
		Condition string  `json:"condition"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	alert, err := h.alertService.CreateAlert(userID, req.Symbol, req.Price, req.Condition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(alert)
}

func (h *TradingHandler) GetUserAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	alerts, err := h.alertService.GetUserAlerts(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

func (h *TradingHandler) GetCandles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Missing symbol parameter", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	candles, err := h.ohlcService.GetCandles(symbol, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candles)
}

func (h *TradingHandler) GetMarketData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all candles for all symbols to compute market data
	symbols := []string{"RELIANCE-CE-2900", "HDFC-PE-1400", "INFY-CE-1500", "TCS-PE-3500"}

	var stocks []map[string]interface{}

	for _, symbol := range symbols {
		candles, err := h.ohlcService.GetCandles(symbol, 1)
		if err != nil || len(candles) == 0 {
			continue
		}

		latest := candles[len(candles)-1]

		// Get previous candle for comparison if available
		var previousClose float64 = latest.Open
		allCandles, _ := h.ohlcService.GetCandles(symbol, 2)
		if len(allCandles) > 1 {
			previousClose = allCandles[0].Close
		}

		change := latest.Close - previousClose
		changePercent := (change / previousClose) * 100

		stock := map[string]interface{}{
			"symbol":         symbol,
			"current_price":  latest.Close,
			"change":         change,
			"change_percent": changePercent,
			"high":           latest.High,
			"low":            latest.Low,
			"volume":         latest.Volume,
		}

		stocks = append(stocks, stock)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]interface{}{
			"stocks": stocks,
		},
	})
}

func (h *TradingHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	transactions, err := h.transactionRepository.FindTransactionsByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transactions": transactions,
		"count":        len(transactions),
	})
}
