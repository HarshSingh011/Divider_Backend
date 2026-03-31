package domain

import (
	"fmt"
	"sync"
	"time"
)

type WalletServiceImpl struct {
	transactionRepo TransactionRepository
	mu              sync.RWMutex
}

func NewWalletService(transactionRepo TransactionRepository) *WalletServiceImpl {
	return &WalletServiceImpl{
		transactionRepo: transactionRepo,
	}
}

func (ws *WalletServiceImpl) ExecuteTrade(userID, symbol string, quantity, price float64, tradeType string) error {
	if tradeType != "BUY" && tradeType != "SELL" {
		return fmt.Errorf("invalid trade type: %s", tradeType)
	}

	if quantity <= 0 || price <= 0 {
		return fmt.Errorf("quantity and price must be positive")
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	amount := quantity * price
	fee := amount * 0.005

	transaction := &Transaction{
		ID:        fmt.Sprintf("txn_%d", time.Now().UnixNano()),
		UserID:    userID,
		Type:      tradeType,
		Symbol:    symbol,
		Quantity:  quantity,
		Price:     price,
		Fee:       fee,
		Amount:    amount,
		Status:    "COMPLETED",
		CreatedAt: time.Now(),
	}

	if err := ws.transactionRepo.SaveTransaction(transaction); err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	fmt.Printf("[TRANSACTION] %s %s: %d units of %s @ %.2f (Fee: %.2f)\n",
		userID, tradeType, int(quantity), symbol, price, fee)

	return nil
}

func (ws *WalletServiceImpl) DepositCash(userID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	transaction := &Transaction{
		ID:        fmt.Sprintf("txn_%d", time.Now().UnixNano()),
		UserID:    userID,
		Type:      "DEPOSIT",
		Symbol:    "CASH",
		Quantity:  1,
		Price:     amount,
		Amount:    amount,
		Fee:       0,
		Status:    "COMPLETED",
		Remarks:   "Cash deposit to wallet",
		CreatedAt: time.Now(),
	}

	if err := ws.transactionRepo.SaveTransaction(transaction); err != nil {
		return fmt.Errorf("failed to save deposit transaction: %w", err)
	}

	fmt.Printf("[WALLET] Deposit: %.2f to %s\n", amount, userID)
	return nil
}

func (ws *WalletServiceImpl) GetWalletSnapshot(userID string) (*WalletSnapshot, error) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	transactions, err := ws.transactionRepo.FindTransactionsByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	availableCash := 100000.0
	positions := make(map[string]Position)

	for _, txn := range transactions {
		switch txn.Type {
		case "DEPOSIT":
			availableCash += txn.Amount

		case "BUY":
			availableCash -= (txn.Amount + txn.Fee)

			if pos, exists := positions[txn.Symbol]; exists {
				totalCost := (pos.AverageCost * pos.Quantity) + txn.Amount
				totalQty := pos.Quantity + txn.Quantity
				pos.AverageCost = totalCost / totalQty
				pos.Quantity = totalQty
				positions[txn.Symbol] = pos
			} else {
				positions[txn.Symbol] = Position{
					Symbol:      txn.Symbol,
					Quantity:    txn.Quantity,
					AverageCost: txn.Price,
				}
			}

		case "SELL":
			availableCash += (txn.Amount - txn.Fee)

			if pos, exists := positions[txn.Symbol]; exists {
				pos.Quantity -= txn.Quantity
				if pos.Quantity <= 0 {
					delete(positions, txn.Symbol)
				} else {
					positions[txn.Symbol] = pos
				}
			}

		case "WITHDRAWAL":
			availableCash -= txn.Amount

		case "BROKERAGE_FEE":
			availableCash -= txn.Fee
		}
	}

	totalBalance := availableCash
	for symbol := range positions {
		pos := positions[symbol]
		totalBalance += (pos.AverageCost * pos.Quantity)
	}

	snapshot := &WalletSnapshot{
		UserID:        userID,
		TotalBalance:  totalBalance,
		AvailableCash: availableCash,
		Positions:     positions,
		LastUpdated:   time.Now(),
	}

	return snapshot, nil
}
