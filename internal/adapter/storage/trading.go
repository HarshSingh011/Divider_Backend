package storage

import (
	"errors"
	"sync"
	"time"

	"stocktrack-backend/internal/domain"
)

type InMemoryCandleRepository struct {
	candles map[string][]domain.Candle
	mu      sync.RWMutex
}

func NewInMemoryCandleRepository() *InMemoryCandleRepository {
	return &InMemoryCandleRepository{
		candles: make(map[string][]domain.Candle),
	}
}

func (r *InMemoryCandleRepository) SaveCandle(candle *domain.Candle) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if candle == nil || candle.Symbol == "" {
		return errors.New("invalid candle")
	}

	r.candles[candle.Symbol] = append(r.candles[candle.Symbol], *candle)
	return nil
}

func (r *InMemoryCandleRepository) GetCandles(symbol string, limit int) ([]domain.Candle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	candles, exists := r.candles[symbol]
	if !exists {
		return []domain.Candle{}, nil
	}

	start := len(candles) - limit
	if start < 0 {
		start = 0
	}

	result := make([]domain.Candle, len(candles[start:]))
	copy(result, candles[start:])
	return result, nil
}

func (r *InMemoryCandleRepository) GetCandlesByTimeRange(symbol string, from, to time.Time) ([]domain.Candle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	candles, exists := r.candles[symbol]
	if !exists {
		return []domain.Candle{}, nil
	}

	var result []domain.Candle
	for _, candle := range candles {
		if candle.Timestamp.After(from) && candle.Timestamp.Before(to) {
			result = append(result, candle)
		}
	}

	return result, nil
}

type InMemoryAlertRepository struct {
	alerts map[string]*domain.Alert
	mu     sync.RWMutex
}

func NewInMemoryAlertRepository() *InMemoryAlertRepository {
	return &InMemoryAlertRepository{
		alerts: make(map[string]*domain.Alert),
	}
}

func (r *InMemoryAlertRepository) SaveAlert(alert *domain.Alert) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if alert == nil || alert.ID == "" {
		return errors.New("invalid alert")
	}

	r.alerts[alert.ID] = alert
	return nil
}

func (r *InMemoryAlertRepository) FindAlertsByUser(userID string) ([]domain.Alert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Alert
	for _, alert := range r.alerts {
		if alert.UserID == userID {
			result = append(result, *alert)
		}
	}

	return result, nil
}

func (r *InMemoryAlertRepository) FindActiveAlerts() ([]domain.Alert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Alert
	for _, alert := range r.alerts {
		if alert.IsActive {
			result = append(result, *alert)
		}
	}

	return result, nil
}

func (r *InMemoryAlertRepository) UpdateAlert(alert *domain.Alert) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.alerts[alert.ID]; !exists {
		return errors.New("alert not found")
	}

	r.alerts[alert.ID] = alert
	return nil
}

func (r *InMemoryAlertRepository) DeleteAlert(alertID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.alerts, alertID)
	return nil
}

type InMemoryTransactionRepository struct {
	transactions map[string][]domain.Transaction
	mu           sync.RWMutex
}

func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
	return &InMemoryTransactionRepository{
		transactions: make(map[string][]domain.Transaction),
	}
}

func (r *InMemoryTransactionRepository) SaveTransaction(transaction *domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if transaction == nil || transaction.UserID == "" {
		return errors.New("invalid transaction")
	}

	r.transactions[transaction.UserID] = append(r.transactions[transaction.UserID], *transaction)
	return nil
}

func (r *InMemoryTransactionRepository) FindTransactionsByUser(userID string) ([]domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactions, exists := r.transactions[userID]
	if !exists {
		return []domain.Transaction{}, nil
	}

	result := make([]domain.Transaction, len(transactions))
	copy(result, transactions)
	return result, nil
}

func (r *InMemoryTransactionRepository) FindTransactionsBySymbol(userID, symbol string) ([]domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transactions, exists := r.transactions[userID]
	if !exists {
		return []domain.Transaction{}, nil
	}

	var result []domain.Transaction
	for _, txn := range transactions {
		if txn.Symbol == symbol {
			result = append(result, txn)
		}
	}

	return result, nil
}

func (r *InMemoryTransactionRepository) GetUserBalance(userID string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	balance := 100000.0

	transactions, exists := r.transactions[userID]
	if !exists {
		return balance, nil
	}

	for _, txn := range transactions {
		switch txn.Type {
		case "DEPOSIT":
			balance += txn.Amount
		case "WITHDRAWAL", "BUY":
			balance -= (txn.Amount + txn.Fee)
		case "SELL":
			balance += (txn.Amount - txn.Fee)
		case "BROKERAGE_FEE":
			balance -= txn.Fee
		}
	}

	return balance, nil
}
