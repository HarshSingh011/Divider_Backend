package domain

import (
	"fmt"
	"sync"
	"time"
)

// Constants for trading limits
const (
	MAX_QUANTITY_PER_TRADE = 10000  // Maximum 10,000 shares per transaction
	MAX_SELL_QUANTITY      = 10000  // Maximum 10,000 shares for SELL
)

type WalletServiceImpl struct {
	transactionRepo TransactionRepository
	stockRepo       StockRepository
	marketEngine    *MarketEngine
	mu              sync.RWMutex
}

func NewWalletService(transactionRepo TransactionRepository) *WalletServiceImpl {
	return &WalletServiceImpl{
		transactionRepo: transactionRepo,
	}
}

func (ws *WalletServiceImpl) SetStockRepository(stockRepo StockRepository) {
	ws.stockRepo = stockRepo
}

func (ws *WalletServiceImpl) SetMarketEngine(engine *MarketEngine) {
	ws.marketEngine = engine
}

// Helper function to calculate current holdings for a symbol
func (ws *WalletServiceImpl) getCurrentHoldings(userID, symbol string) float64 {
	transactions, err := ws.transactionRepo.FindTransactionsByUser(userID)
	if err != nil {
		return 0
	}

	quantity := 0.0
	for _, txn := range transactions {
		if txn.Symbol == symbol {
			if txn.Type == "BUY" {
				quantity += txn.Quantity
			} else if txn.Type == "SELL" {
				quantity -= txn.Quantity
			}
		}
	}
	return quantity
}

// Helper function to calculate current available cash
func (ws *WalletServiceImpl) getCurrentCash(userID string) float64 {
	transactions, err := ws.transactionRepo.FindTransactionsByUser(userID)
	if err != nil {
		return 100000.0 // Default initial cash
	}

	availableCash := 100000.0
	for _, txn := range transactions {
		switch txn.Type {
		case "DEPOSIT":
			availableCash += txn.Amount
		case "BUY":
			availableCash -= (txn.Amount + txn.Fee)
		case "SELL":
			availableCash += (txn.Amount - txn.Fee)
		case "WITHDRAWAL":
			availableCash -= txn.Amount
		case "BROKERAGE_FEE":
			availableCash -= txn.Fee
		}
	}
	return availableCash
}

// Helper function to calculate available stock quantity (not held by any user)
func (ws *WalletServiceImpl) getAvailableStockQuantity(symbol string) (float64, error) {
	// If no stock repo, assume unlimited stock (backward compatibility)
	if ws.stockRepo == nil {
		return 999999999, nil
	}

	stock, err := ws.stockRepo.GetStock(symbol)
	if err != nil {
		// Stock not found, assume unlimited
		return 999999999, nil
	}

	// Get all transactions for this symbol to calculate total held
	transactions, err := ws.transactionRepo.FindTransactionsByUser("")
	if err != nil {
		return stock.TotalAvailableQty, nil
	}

	totalHeld := 0.0
	for _, txn := range transactions {
		if txn.Symbol == symbol {
			if txn.Type == "BUY" {
				totalHeld += txn.Quantity
			} else if txn.Type == "SELL" {
				totalHeld -= txn.Quantity
			}
		}
	}

	// Available = Total - Currently held by users
	available := stock.TotalAvailableQty - totalHeld
	if available < 0 {
		available = 0
	}

	return available, nil
}

func (ws *WalletServiceImpl) ExecuteTrade(userID, symbol string, quantity, price float64, tradeType string) error {
	if tradeType != "BUY" && tradeType != "SELL" {
		return fmt.Errorf("invalid trade type: %s", tradeType)
	}

	if quantity <= 0 || price <= 0 {
		return fmt.Errorf("quantity and price must be positive")
	}

	// Validate maximum quantity limit
	if quantity > float64(MAX_QUANTITY_PER_TRADE) {
		return fmt.Errorf("maximum quantity exceeded: you can only buy/sell up to %.0f shares per transaction, but trying to %s %.0f", float64(MAX_QUANTITY_PER_TRADE), tradeType, quantity)
	}

	// Validate SELL: Check if user has sufficient shares
	if tradeType == "SELL" {
		holdings := ws.getCurrentHoldings(userID, symbol)
		if holdings < quantity {
			return fmt.Errorf("insufficient holdings: you have %.2f shares of %s but trying to sell %.2f", holdings, symbol, quantity)
		}
		if holdings == 0 {
			return fmt.Errorf("you do not own any shares of %s", symbol)
		}
	}

	// Validate BUY: Check if user has sufficient cash
	if tradeType == "BUY" {
		requiredCash := (quantity * price) + (quantity * price * 0.005) // amount +fee
		availableCash := ws.getCurrentCash(userID)
		if availableCash < requiredCash {
			return fmt.Errorf("insufficient cash: you have ₹%.2f but need ₹%.2f", availableCash, requiredCash)
		}

		// Check if stock quantity is available
		availableStock, err := ws.getAvailableStockQuantity(symbol)
		if err == nil && availableStock < quantity {
			return fmt.Errorf("insufficient stock: only %.2f shares of %s available in market, but trying to buy %.2f", availableStock, symbol, quantity)
		}
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
		CreatedAt: time.Now().In(ISTLocation),
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
		CreatedAt: time.Now().In(ISTLocation),
	}

	if err := ws.transactionRepo.SaveTransaction(transaction); err != nil {
		return fmt.Errorf("failed to save deposit transaction: %w", err)
	}

	fmt.Printf("[WALLET] Deposit: %.2f to %s\n", amount, userID)
	return nil
}

func (ws *WalletServiceImpl) WithdrawCash(userID string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive")
	}

	availableCash := ws.getCurrentCash(userID)
	if availableCash < amount {
		return fmt.Errorf("insufficient funds: you have ₹%.2f but trying to withdraw ₹%.2f", availableCash, amount)
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	transaction := &Transaction{
		ID:        fmt.Sprintf("txn_%d", time.Now().UnixNano()),
		UserID:    userID,
		Type:      "WITHDRAWAL",
		Symbol:    "CASH",
		Quantity:  1,
		Price:     amount,
		Amount:    amount,
		Fee:       0,
		Status:    "COMPLETED",
		Remarks:   "Cash withdrawal from wallet",
		CreatedAt: time.Now().In(ISTLocation),
	}

	if err := ws.transactionRepo.SaveTransaction(transaction); err != nil {
		return fmt.Errorf("failed to save withdrawal transaction: %w", err)
	}

	fmt.Printf("[WALLET] Withdrawal: %.2f from %s\n", amount, userID)
	return nil
}

func (ws *WalletServiceImpl) GetWalletSnapshot(userID string) (*WalletSnapshot, error) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	transactions, err := ws.transactionRepo.FindTransactionsByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Get current market prices
	var currentPrices map[string]float64
	if ws.marketEngine != nil {
		currentPrices = make(map[string]float64)
		prices := ws.marketEngine.GetCurrentPrices()
		for _, tick := range prices {
			currentPrices[tick.Symbol] = tick.CurrentPrice
		}
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

	// Calculate current prices, PnL, and percentages for all positions
	if currentPrices != nil {
		for symbol, pos := range positions {
			if currentPrice, ok := currentPrices[symbol]; ok {
				pos.CurrentPrice = currentPrice
				// Calculate unrealized PnL
				costBasis := pos.AverageCost * pos.Quantity
				currentValue := currentPrice * pos.Quantity
				pos.UnrealizedPnL = currentValue - costBasis
				// Calculate percentage change
				if costBasis > 0 {
					pos.Percentage = (pos.UnrealizedPnL / costBasis) * 100
				}
				positions[symbol] = pos
			}
		}
	} else {
		// If market engine not available, use average cost as current price
		for symbol, pos := range positions {
			pos.CurrentPrice = pos.AverageCost
			pos.UnrealizedPnL = 0
			pos.Percentage = 0
			positions[symbol] = pos
		}
	}

	totalBalance := availableCash
	for symbol := range positions {
		pos := positions[symbol]
		if pos.CurrentPrice > 0 {
			totalBalance += (pos.CurrentPrice * pos.Quantity)
		} else {
			totalBalance += (pos.AverageCost * pos.Quantity)
		}
	}

	// Calculate invested amount (holdings value at current prices)
	investedAmount := totalBalance - availableCash

	if investedAmount < 0 {
		investedAmount = 0
	}

	snapshot := &WalletSnapshot{
		UserID:         userID,
		TotalBalance:   totalBalance,
		AvailableCash:  availableCash,
		InvestedAmount: investedAmount,
		Positions:      positions,
		LastUpdated:    time.Now().In(ISTLocation),
	}

	return snapshot, nil
}
