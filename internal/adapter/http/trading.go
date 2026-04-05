package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

func (h *TradingHandler) sendJSONError(w http.ResponseWriter, statusCode int, errorCode string, message string, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    errorCode,
			"message": message,
		},
	}
	if details != nil {
		response["error"].(map[string]interface{})["details"] = details
	}
	json.NewEncoder(w).Encode(response)
}

func (h *TradingHandler) ExecuteTrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST method is allowed", nil)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.sendJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated. X-User-ID header is required", nil)
		return
	}

	var req struct {
		Symbol   string  `json:"symbol"`
		Quantity float64 `json:"quantity"`
		Price    float64 `json:"price"`
		Type     string  `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid JSON request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := h.walletService.ExecuteTrade(userID, req.Symbol, req.Quantity, req.Price, req.Type); err != nil {
		// Determine error type based on error message
		errorCode := "TRADE_FAILED"
		details := map[string]interface{}{
			"symbol":   req.Symbol,
			"quantity": req.Quantity,
			"price":    req.Price,
			"type":     req.Type,
		}

		if strings.Contains(err.Error(), "insufficient holdings") {
			errorCode = "INSUFFICIENT_HOLDINGS"
		} else if strings.Contains(err.Error(), "insufficient cash") {
			errorCode = "INSUFFICIENT_CASH"
		} else if strings.Contains(err.Error(), "insufficient stock") {
			errorCode = "INSUFFICIENT_STOCK"
		} else if strings.Contains(err.Error(), "invalid trade type") {
			errorCode = "INVALID_TRADE_TYPE"
		} else if strings.Contains(err.Error(), "quantity and price must be positive") {
			errorCode = "INVALID_QUANTITY_OR_PRICE"
		} else if strings.Contains(err.Error(), "maximum quantity exceeded") {
			errorCode = "MAXIMUM_QUANTITY_EXCEEDED"
		}

		h.sendJSONError(w, http.StatusBadRequest, errorCode, err.Error(), details)
		return
	}

	total := req.Quantity * req.Price
	fee := total * 0.005
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"status":   "Trade executed successfully",
			"total":    total,
			"fee":      fee,
			"symbol":   req.Symbol,
			"quantity": req.Quantity,
			"price":    req.Price,
			"type":     req.Type,
		},
	})
}

func (h *TradingHandler) GetWalletSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET method is allowed", nil)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.sendJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated. X-User-ID header is required", nil)
		return
	}

	snapshot, err := h.walletService.GetWalletSnapshot(userID)
	if err != nil {
		h.sendJSONError(w, http.StatusInternalServerError, "WALLET_SNAPSHOT_ERROR", "Failed to get wallet snapshot", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    snapshot,
	})
}

func (h *TradingHandler) DepositCash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST method is allowed", nil)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.sendJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated. X-User-ID header is required", nil)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid JSON request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := h.walletService.DepositCash(userID, req.Amount); err != nil {
		errorCode := "DEPOSIT_FAILED"
		details := map[string]interface{}{
			"amount": req.Amount,
		}

		if strings.Contains(err.Error(), "deposit amount must be positive") {
			errorCode = "INVALID_AMOUNT"
		}

		h.sendJSONError(w, http.StatusBadRequest, errorCode, err.Error(), details)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"status": "Deposit successful",
			"amount": req.Amount,
		},
	})
}

func (h *TradingHandler) WithdrawCash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST method is allowed", nil)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		h.sendJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated. X-User-ID header is required", nil)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid JSON request body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := h.walletService.WithdrawCash(userID, req.Amount); err != nil {
		errorCode := "WITHDRAWAL_FAILED"
		details := map[string]interface{}{
			"amount": req.Amount,
		}

		if strings.Contains(err.Error(), "withdrawal amount must be positive") {
			errorCode = "INVALID_AMOUNT"
		} else if strings.Contains(err.Error(), "insufficient funds") {
			errorCode = "INSUFFICIENT_FUNDS"
		}

		h.sendJSONError(w, http.StatusBadRequest, errorCode, err.Error(), details)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"status": "Withdrawal successful",
			"amount": req.Amount,
		},
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
