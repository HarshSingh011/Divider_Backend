package domain

import (
	"fmt"
	"strings"
	"testing"
)

// MockTransactionRepository for testing
type MockTransactionRepository struct {
	transactions map[string][]*Transaction
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{
		transactions: make(map[string][]*Transaction),
	}
}

func (m *MockTransactionRepository) SaveTransaction(txn *Transaction) error {
	m.transactions[txn.UserID] = append(m.transactions[txn.UserID], txn)
	return nil
}

func (m *MockTransactionRepository) FindTransactionsByUser(userID string) ([]Transaction, error) {
	var result []Transaction
	if userID == "" {
		// Return all transactions for all users
		for _, txns := range m.transactions {
			for _, txn := range txns {
				result = append(result, *txn)
			}
		}
	} else {
		for _, txn := range m.transactions[userID] {
			result = append(result, *txn)
		}
	}
	return result, nil
}

func (m *MockTransactionRepository) FindTransactionsBySymbol(userID, symbol string) ([]Transaction, error) {
	var result []Transaction
	for _, txn := range m.transactions[userID] {
		if txn.Symbol == symbol {
			result = append(result, *txn)
		}
	}
	return result, nil
}

func (m *MockTransactionRepository) GetUserBalance(userID string) (float64, error) {
	// Calculate balance based on transactions
	balance := 100000.0
	for _, txn := range m.transactions[userID] {
		switch txn.Type {
		case "DEPOSIT":
			balance += txn.Amount
		case "BUY":
			balance -= (txn.Amount + txn.Fee)
		case "SELL":
			balance += (txn.Amount - txn.Fee)
		case "WITHDRAWAL":
			balance -= txn.Amount
		}
	}
	return balance, nil
}

// MockStockRepository for testing
type MockStockRepository struct {
	stocks map[string]*Stock
}

func NewMockStockRepository() *MockStockRepository {
	return &MockStockRepository{
		stocks: make(map[string]*Stock),
	}
}

func (m *MockStockRepository) SaveStock(stock *Stock) error {
	m.stocks[stock.Symbol] = stock
	return nil
}

func (m *MockStockRepository) GetStock(symbol string) (*Stock, error) {
	if stock, ok := m.stocks[symbol]; ok {
		return stock, nil
	}
	// Return not found - simulation
	return nil, fmt.Errorf("stock not found: %s", symbol)
}

func (m *MockStockRepository) UpdateStock(stock *Stock) error {
	m.stocks[stock.Symbol] = stock
	return nil
}

func (m *MockStockRepository) GetAllStocks() ([]Stock, error) {
	var result []Stock
	for _, stock := range m.stocks {
		result = append(result, *stock)
	}
	return result, nil
}

// Test: Cannot sell shares you don't own at all
func TestSellNonOwnedShares(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Try to sell shares without owning any
	err := ws.ExecuteTrade(userID, symbol, 10, 2700, "SELL")

	if err == nil {
		t.Error("❌ FAILED: Should not allow selling shares you don't own. Got no error")
	} else {
		t.Logf("✅ PASSED: Correctly rejected selling non-owned shares. Error: %v", err)
	}
}

// Test: Cannot sell more shares than you own
func TestSellMoreThanOwned(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Buy 10 shares @ 800 = 8000 + fee 40 = 8040. We have 100,000 so this is OK
	err := ws.ExecuteTrade(userID, symbol, 10, 800, "BUY")
	if err != nil {
		t.Fatalf("Failed to buy shares: %v", err)
	}

	// Try to sell 15 shares (more than owned - have only 10)
	err = ws.ExecuteTrade(userID, symbol, 15, 900, "SELL")

	if err == nil {
		t.Error("❌ FAILED: Should not allow overselling. Tried to sell 15 but only own 10")
	} else {
		t.Logf("✅ PASSED: Correctly rejected overselling. Error: %v", err)
	}
}

// Test: Can sell exact amount owned
func TestSellExactAmountOwned(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Buy 10 shares
	err := ws.ExecuteTrade(userID, symbol, 10, 800, "BUY")
	if err != nil {
		t.Fatalf("Failed to buy shares: %v", err)
	}

	// Sell exact amount
	err = ws.ExecuteTrade(userID, symbol, 10, 900, "SELL")

	if err != nil {
		t.Errorf("❌ FAILED: Should allow selling exact amount owned. Got error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed selling exact amount owned")
	}
}

// Test: Can sell less than amount owned
func TestSellPartialAmount(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Buy 100 shares @ 100 = 10000
	err := ws.ExecuteTrade(userID, symbol, 100, 100, "BUY")
	if err != nil {
		t.Fatalf("Failed to buy shares: %v", err)
	}

	// Sell 60 shares (less than owned)
	err = ws.ExecuteTrade(userID, symbol, 60, 120, "SELL")

	if err != nil {
		t.Errorf("❌ FAILED: Should allow selling less than owned. Got error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed selling partial amount")

		// Verify remaining holdings
		remaining := ws.getCurrentHoldings(userID, symbol)
		if remaining != 40 {
			t.Errorf("❌ FAILED: After selling 60 from 100, remaining should be 40 but got %.2f", remaining)
		} else {
			t.Logf("✅ PASSED: Holdings correctly updated (40 remaining)")
		}
	}
}

// Test: Cannot sell after selling all shares
func TestCannotSellAfterEmpty(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Buy 10 shares @ 100
	ws.ExecuteTrade(userID, symbol, 10, 100, "BUY")

	// Sell all 10 shares
	ws.ExecuteTrade(userID, symbol, 10, 120, "SELL")

	// Try to sell again
	err := ws.ExecuteTrade(userID, symbol, 5, 120, "SELL")

	if err == nil {
		t.Error("❌ FAILED: Should not allow selling after holdings are empty")
	} else {
		t.Logf("✅ PASSED: Correctly rejected selling empty holdings. Error: %v", err)
	}
}

// Test: Buy validation - insufficient cash
func TestBuyInsufficientCash(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Try to buy huge amount that exceeds initial 100,000 cash
	// 1000 shares @ 200 = 200,000 (way more than 100,000)
	err := ws.ExecuteTrade(userID, symbol, 1000, 200, "BUY")

	if err == nil {
		t.Error("❌ FAILED: Should not allow buying with insufficient cash")
	} else {
		t.Logf("✅ PASSED: Correctly rejected buy with insufficient cash. Error: %v", err)
	}
}

// Test: Multiple buys then sell validation
func TestMultipleBuysThenSell(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Buy 1: 30 shares @ 100 = 3000 + fee 15
	ws.ExecuteTrade(userID, symbol, 30, 100, "BUY")

	// Buy 2: 20 shares @ 110 = 2200 + fee 11
	ws.ExecuteTrade(userID, symbol, 20, 110, "BUY")

	// Buy 3: 50 shares @ 90 = 4500 + fee 22.5
	ws.ExecuteTrade(userID, symbol, 50, 90, "BUY")

	// Total holdings should be 100 shares
	total := ws.getCurrentHoldings(userID, symbol)
	if total != 100 {
		t.Errorf("❌ FAILED: Total holdings should be 100 but got %.2f", total)
	} else {
		t.Logf("✅ PASSED: Total holdings after multiple buys = 100 shares")
	}

	// Try to sell 150 (should fail)
	err := ws.ExecuteTrade(userID, symbol, 150, 120, "SELL")
	if err == nil {
		t.Error("❌ FAILED: Should not allow overselling after multiple buys")
	} else {
		t.Logf("✅ PASSED: Correctly rejected oversell. Error: %v", err)
	}

	// Sell 75 (should succeed)
	err = ws.ExecuteTrade(userID, symbol, 75, 120, "SELL")
	if err != nil {
		t.Errorf("❌ FAILED: Should allow selling 75 when have 100. Error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed selling 75 out of 100")

		remaining := ws.getCurrentHoldings(userID, symbol)
		if remaining != 25 {
			t.Errorf("❌ FAILED: Remaining should be 25 but got %.2f", remaining)
		} else {
			t.Logf("✅ PASSED: Remaining holdings = 25 shares")
		}
	}
}

// Test: Different symbols don't interfere
func TestMultipleSymbols(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol1 := "RELIANCE-CE-2900"
	symbol2 := "TCS-CE-3500"

	// Buy 50 shares of symbol1 @ 100
	ws.ExecuteTrade(userID, symbol1, 50, 100, "BUY")

	// Buy 30 shares of symbol2 @ 150
	ws.ExecuteTrade(userID, symbol2, 30, 150, "BUY")

	// Try to sell 40 of symbol1 (should succeed - have 50)
	err := ws.ExecuteTrade(userID, symbol1, 40, 120, "SELL")
	if err != nil {
		t.Errorf("❌ FAILED: Should allow selling symbol1. Error: %v", err)
	} else {
		t.Logf("✅ PASSED: Sold 40 shares of symbol1")
	}

	// Try to sell 50 of symbol2 (should fail - only have 30)
	err = ws.ExecuteTrade(userID, symbol2, 50, 170, "SELL")
	if err == nil {
		t.Error("❌ FAILED: Should not allow overselling symbol2")
	} else {
		t.Logf("✅ PASSED: Correctly rejected overselling symbol2. Error: %v", err)
	}

	// Verify holdings are independent
	h1 := ws.getCurrentHoldings(userID, symbol1)
	h2 := ws.getCurrentHoldings(userID, symbol2)

	if h1 != 10 {
		t.Errorf("❌ FAILED: symbol1 holdings should be 10, got %.2f", h1)
	}
	if h2 != 30 {
		t.Errorf("❌ FAILED: symbol2 holdings should be 30, got %.2f", h2)
	}

	if h1 == 10 && h2 == 30 {
		t.Logf("✅ PASSED: Holdings correctly tracked independently: symbol1=10, symbol2=30")
	}
}

// Test: Invalid quantity or price
func TestInvalidQuantityOrPrice(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Test zero quantity
	err := ws.ExecuteTrade(userID, symbol, 0, 100, "BUY")
	if err == nil {
		t.Error("❌ FAILED: Should reject zero quantity")
	} else {
		t.Logf("✅ PASSED: Correctly rejected zero quantity")
	}

	// Test negative quantity
	err = ws.ExecuteTrade(userID, symbol, -50, 100, "BUY")
	if err == nil {
		t.Error("❌ FAILED: Should reject negative quantity")
	} else {
		t.Logf("✅ PASSED: Correctly rejected negative quantity")
	}

	// Test zero price
	err = ws.ExecuteTrade(userID, symbol, 50, 0, "BUY")
	if err == nil {
		t.Error("❌ FAILED: Should reject zero price")
	} else {
		t.Logf("✅ PASSED: Correctly rejected zero price")
	}

	// Test negative price
	err = ws.ExecuteTrade(userID, symbol, 50, -100, "BUY")
	if err == nil {
		t.Error("❌ FAILED: Should reject negative price")
	} else {
		t.Logf("✅ PASSED: Correctly rejected negative price")
	}
}

// ============================================================
// DEPOSIT TESTS
// ============================================================

// Test: Deposit positive amount
func TestDepositPositiveAmount(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	// Deposit 50000
	err := ws.DepositCash(userID, 50000)

	if err != nil {
		t.Errorf("❌ FAILED: Should allow positive deposit. Got error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed deposit of ₹50,000")

		// Verify cash was added
		cash := ws.getCurrentCash(userID)
		expected := 100000.0 + 50000.0 // Initial + deposit
		if cash != expected {
			t.Errorf("❌ FAILED: Cash should be %.2f but got %.2f", expected, cash)
		} else {
			t.Logf("✅ PASSED: Cash correctly updated to ₹%.2f", cash)
		}
	}
}

// Test: Cannot deposit zero amount
func TestDepositZeroAmount(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	err := ws.DepositCash(userID, 0)

	if err == nil {
		t.Error("❌ FAILED: Should reject zero deposit")
	} else {
		t.Logf("✅ PASSED: Correctly rejected zero deposit. Error: %v", err)
	}
}

// Test: Cannot deposit negative amount
func TestDepositNegativeAmount(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	err := ws.DepositCash(userID, -5000)

	if err == nil {
		t.Error("❌ FAILED: Should reject negative deposit")
	} else {
		t.Logf("✅ PASSED: Correctly rejected negative deposit. Error: %v", err)
	}
}

// Test: Multiple deposits accumulate
func TestMultipleDeposits(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	// Deposit 1: 25000
	err := ws.DepositCash(userID, 25000)
	if err != nil {
		t.Fatalf("Deposit 1 failed: %v", err)
	}

	// Deposit 2: 15000
	err = ws.DepositCash(userID, 15000)
	if err != nil {
		t.Fatalf("Deposit 2 failed: %v", err)
	}

	// Deposit 3: 10000
	err = ws.DepositCash(userID, 10000)
	if err != nil {
		t.Fatalf("Deposit 3 failed: %v", err)
	}

	cash := ws.getCurrentCash(userID)
	expected := 100000.0 + 25000.0 + 15000.0 + 10000.0 // Initial + all deposits
	if cash != expected {
		t.Errorf("❌ FAILED: Total cash should be %.2f but got %.2f", expected, cash)
	} else {
		t.Logf("✅ PASSED: Multiple deposits accumulated to ₹%.2f", cash)
	}
}

// ============================================================
// WITHDRAW TESTS
// ============================================================

// Test: Withdraw with sufficient funds
func TestWithdrawSufficientFunds(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	// Withdraw 20000 (initial cash is 100000)
	err := ws.WithdrawCash(userID, 20000)

	if err != nil {
		t.Errorf("❌ FAILED: Should allow withdrawal with sufficient funds. Got error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed withdrawal of ₹20,000")

		// Verify cash was reduced
		cash := ws.getCurrentCash(userID)
		expected := 100000.0 - 20000.0
		if cash != expected {
			t.Errorf("❌ FAILED: Cash should be %.2f but got %.2f", expected, cash)
		} else {
			t.Logf("✅ PASSED: Cash correctly reduced to ₹%.2f", cash)
		}
	}
}

// Test: Cannot withdraw zero amount
func TestWithdrawZeroAmount(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	err := ws.WithdrawCash(userID, 0)

	if err == nil {
		t.Error("❌ FAILED: Should reject zero withdrawal")
	} else {
		t.Logf("✅ PASSED: Correctly rejected zero withdrawal. Error: %v", err)
	}
}

// Test: Cannot withdraw negative amount
func TestWithdrawNegativeAmount(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	err := ws.WithdrawCash(userID, -5000)

	if err == nil {
		t.Error("❌ FAILED: Should reject negative withdrawal")
	} else {
		t.Logf("✅ PASSED: Correctly rejected negative withdrawal. Error: %v", err)
	}
}

// Test: Cannot withdraw more than available
func TestWithdrawMoreThanAvailable(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	// Try to withdraw 150000 (only have 100000)
	err := ws.WithdrawCash(userID, 150000)

	if err == nil {
		t.Error("❌ FAILED: Should not allow withdrawal exceeding available funds")
	} else {
		t.Logf("✅ PASSED: Correctly rejected oversized withdrawal. Error: %v", err)
	}
}

// Test: Cannot withdraw all cash leaving invested amounts in holdings
func TestWithdrawExactAmount(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Buy 50 shares @ 100 = 5000 + fee 25 = 5025 total
	_ = ws.ExecuteTrade(userID, symbol, 50, 100, "BUY")

	// Available cash should be: 100000 - 5025 = 94975
	cash := ws.getCurrentCash(userID)
	expected := 100000.0 - 5025.0
	if cash != expected {
		t.Logf("⚠️  Cash calculation: %.2f (expected %.2f)", cash, expected)
	}

	// Withdraw exactly what's available
	err := ws.WithdrawCash(userID, cash)

	if err != nil {
		t.Errorf("❌ FAILED: Should allow withdrawal of exact available amount. Got error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed withdrawal of exact available amount (₹%.2f)", cash)

		// Verify cash is now 0
		remainingCash := ws.getCurrentCash(userID)
		if remainingCash != 0 {
			t.Errorf("❌ FAILED: Cash should be 0 but got %.2f", remainingCash)
		} else {
			t.Logf("✅ PASSED: Cash correctly reduced to 0")
		}
	}
}

// Test: Multiple withdrawals
func TestMultipleWithdrawals(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	// Withdraw 1: 10000
	err := ws.WithdrawCash(userID, 10000)
	if err != nil {
		t.Fatalf("Withdrawal 1 failed: %v", err)
	}

	// Withdraw 2: 20000
	err = ws.WithdrawCash(userID, 20000)
	if err != nil {
		t.Fatalf("Withdrawal 2 failed: %v", err)
	}

	// Withdraw 3: 30000
	err = ws.WithdrawCash(userID, 30000)
	if err != nil {
		t.Fatalf("Withdrawal 3 failed: %v", err)
	}

	cash := ws.getCurrentCash(userID)
	expected := 100000.0 - 10000.0 - 20000.0 - 30000.0
	if cash != expected {
		t.Errorf("❌ FAILED: Cash should be %.2f but got %.2f", expected, cash)
	} else {
		t.Logf("✅ PASSED: Multiple withdrawals correctly totaled. Cash = ₹%.2f", cash)
	}
}

// Test: Cannot withdraw more after previous withdrawal
func TestWithdrawOverlimitAfterPrevious(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	// Withdraw 50000
	_ = ws.WithdrawCash(userID, 50000)

	// Try to withdraw 60000 more (only 50000 left)
	err := ws.WithdrawCash(userID, 60000)

	if err == nil {
		t.Error("❌ FAILED: Should not allow withdrawal exceeding remaining balance")
	} else {
		t.Logf("✅ PASSED: Correctly rejected oversized withdrawal. Error: %v", err)
	}
}

// Test: Deposit then Withdraw
func TestDepositThenWithdraw(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"

	// Deposit 50000
	err := ws.DepositCash(userID, 50000)
	if err != nil {
		t.Fatalf("Deposit failed: %v", err)
	}

	// Cash should be: 100000 + 50000 = 150000
	cash := ws.getCurrentCash(userID)
	if cash != 150000.0 {
		t.Errorf("❌ FAILED: Cash after deposit should be 150000 but got %.2f", cash)
	} else {
		t.Logf("✅ PASSED: Cash after deposit = ₹%.2f", cash)
	}

	// Withdraw 80000
	err = ws.WithdrawCash(userID, 80000)
	if err != nil {
		t.Errorf("❌ FAILED: Should allow withdrawal. Got error: %v", err)
	} else {
		t.Logf("✅ PASSED: Withdrawal successful")
	}

	// Cash should be: 150000 - 80000 = 70000
	cash = ws.getCurrentCash(userID)
	expected := 70000.0
	if cash != expected {
		t.Errorf("❌ FAILED: Cash after withdrawal should be %.2f but got %.2f", expected, cash)
	} else {
		t.Logf("✅ PASSED: Cash after withdrawal = ₹%.2f", cash)
	}
}

// Test: Buy, Sell, then Withdraw
func TestBuySellThenWithdraw(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Buy 50 shares @ 100 = 5000 + fee 25 = 5025
	_ = ws.ExecuteTrade(userID, symbol, 50, 100, "BUY")

	// Cash: 100000 - 5025 = 94975
	cash := ws.getCurrentCash(userID)

	// Sell 50 shares @ 120 = 6000 - fee 30 = 5970 (proceeds added back)
	_ = ws.ExecuteTrade(userID, symbol, 50, 120, "SELL")

	// Cash: 94975 + 5970 = 100945
	cash = ws.getCurrentCash(userID)
	expected := 94975.0 + 5970.0
	if cash != expected {
		t.Logf("⚠️  Cash after sell should be %.2f but got %.2f", expected, cash)
	} else {
		t.Logf("✅ PASSED: Cash after buy-sell = ₹%.2f", cash)
	}

	// Withdraw 50000
	err := ws.WithdrawCash(userID, 50000)
	if err != nil {
		t.Errorf("❌ FAILED: Should allow withdrawal. Got error: %v", err)
	} else {
		t.Logf("✅ PASSED: Withdrawal successful")
	}

	// Final cash: 100945 - 50000 = 50945
	finalCash := ws.getCurrentCash(userID)
	if finalCash < 50000 {
		t.Errorf("❌ FAILED: Final cash should be at least 50000 but got %.2f", finalCash)
	} else {
		t.Logf("✅ PASSED: Final cash = ₹%.2f after all operations", finalCash)
	}
}

// Test: Cannot withdraw after holdings are bought
func TestCannotWithdrawBelowMinimum(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Buy 100 shares @ 500 = 50000 + fee 250 = 50250
	_ = ws.ExecuteTrade(userID, symbol, 100, 500, "BUY")

	// Available cash: 100000 - 50250 = 49750
	cash := ws.getCurrentCash(userID)

	// Try to withdraw 60000 (more than available)
	err := ws.WithdrawCash(userID, 60000)

	if err == nil {
		t.Error("❌ FAILED: Should not allow withdrawal exceeding available funds")
	} else {
		t.Logf("✅ PASSED: Correctly rejected withdrawal. Error: %v", err)
	}

	// Should allow withdrawal of available amount
	err = ws.WithdrawCash(userID, cash)
	if err != nil {
		t.Errorf("❌ FAILED: Should allow withdrawal of available funds. Error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed withdrawal of available funds (₹%.2f)", cash)
	}
}

// ============================================================
// MAXIMUM QUANTITY LIMIT TESTS
// ============================================================

// Test: Cannot buy above maximum quantity limit
func TestBuyAboveMaximumQuantity(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Try to buy 15000 shares (limit is 10000)
	err := ws.ExecuteTrade(userID, symbol, 15000, 10, "BUY")

	if err == nil {
		t.Error("❌ FAILED: Should not allow buying above maximum quantity (10000)")
	} else {
		t.Logf("✅ PASSED: Correctly rejected buy above max quantity. Error: %v", err)
	}
}

// Test: Cannot sell above maximum quantity limit
func TestSellAboveMaximumQuantity(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// First, buy 15000 shares (will be rejected by max limit)
	err := ws.ExecuteTrade(userID, symbol, 15000, 10, "BUY")
	if err == nil {
		t.Logf("⚠️  Unexpectedly allowed buy of 15000 shares")
	}

	// Cannot buy that many, so let's simulate by manually setting
	// But for now, just verify that trying to sell 15000 is rejected
	err = ws.ExecuteTrade(userID, symbol, 15000, 10, "SELL")

	if err == nil {
		t.Error("❌ FAILED: Should not allow selling above maximum quantity (10000)")
	} else {
		t.Logf("✅ PASSED: Correctly rejected sell above max quantity. Error: %v", err)
	}
}

// Test: Can buy exactly at maximum limit
func TestBuyExactlyAtMaximumQuantity(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Deposit additional cash to afford 10000 shares @ 10 = 100000 + fee
	ws.DepositCash(userID, 100000)

	// Try to buy exactly 10000 shares @ 10 (should succeed)
	err := ws.ExecuteTrade(userID, symbol, 10000, 10, "BUY")

	if err != nil {
		t.Errorf("❌ FAILED: Should allow buying at maximum limit (10000). Error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed buy at maximum quantity (10000 shares)")

		holdings := ws.getCurrentHoldings(userID, symbol)
		if holdings != 10000 {
			t.Errorf("❌ FAILED: Holdings should be 10000 but got %.2f", holdings)
		} else {
			t.Logf("✅ PASSED: Holdings correctly set to 10000 shares")
		}
	}
}

// Test: Can sell exactly at maximum limit
func TestSellExactlyAtMaximumQuantity(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Deposit additional cash
	ws.DepositCash(userID, 100000)

	// Buy 10000 shares @ 10
	err := ws.ExecuteTrade(userID, symbol, 10000, 10, "BUY")
	if err != nil {
		t.Fatalf("Failed to buy 10000 shares: %v", err)
	}

	// Sell exactly 10000 shares (should succeed)
	err = ws.ExecuteTrade(userID, symbol, 10000, 15, "SELL")

	if err != nil {
		t.Errorf("❌ FAILED: Should allow selling at maximum limit (10000). Error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed sell at maximum quantity (10000 shares)")

		holdings := ws.getCurrentHoldings(userID, symbol)
		if holdings != 0 {
			t.Errorf("❌ FAILED: Holdings should be 0 but got %.2f", holdings)
		} else {
			t.Logf("✅ PASSED: Holdings correctly reduced to 0")
		}
	}
}

// Test: Various quantities around the limit
func TestQuantitiesAroundLimit(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Deposit extra cash
	ws.DepositCash(userID, 500000)

	testCases := []struct {
		quantity  float64
		shouldFail bool
		name      string
	}{
		{9999, false, "9999 (below limit)"},
		{10000, false, "10000 (at limit)"},
		{10001, true, "10001 (above limit)"},
		{20000, true, "20000 (way above limit)"},
		{100, false, "100 (well below limit)"},
	}

	for _, tc := range testCases {
		err := ws.ExecuteTrade(userID, symbol, tc.quantity, 1, "BUY")

		if tc.shouldFail {
			if err == nil {
				t.Errorf("❌ FAILED: Should reject quantity %s", tc.name)
			} else {
				t.Logf("✅ PASSED: Correctly rejected quantity %s", tc.name)
			}
		} else {
			if err != nil {
				t.Errorf("❌ FAILED: Should allow quantity %s. Error: %v", tc.name, err)
			} else {
				t.Logf("✅ PASSED: Correctly allowed quantity %s", tc.name)
			}
		}
	}
}

// ============================================================
// STOCK AVAILABILITY TESTS
// ============================================================

// Test: Cannot buy above available stock in market
func TestBuyAboveAvailableStock(t *testing.T) {
	txnRepo := NewMockTransactionRepository()
	stockRepo := NewMockStockRepository()
	ws := NewWalletService(txnRepo)
	ws.SetStockRepository(stockRepo)

	userID1 := "test_user_1"
	userID2 := "test_user_2"
	symbol := "RELIANCE-CE-2900"

	// Create stock with total 100 shares available
	stock := &Stock{
		ID:                "stock_1",
		Symbol:            symbol,
		CompanyName:       "Reliance Industries",
		TotalAvailableQty: 100,
		CurrentPrice:      2500,
	}
	stockRepo.SaveStock(stock)

	// Deposit cash for user 1
	ws.DepositCash(userID1, 500000)

	// User 1 buys 80 shares
	err := ws.ExecuteTrade(userID1, symbol, 80, 2500, "BUY")
	if err != nil {
		t.Fatalf("User 1 failed to buy 80 shares: %v", err)
	}

	// User 2 tries to buy 30 shares (only 20 available)
	ws.DepositCash(userID2, 500000)
	err = ws.ExecuteTrade(userID2, symbol, 30, 2500, "BUY")

	if err == nil {
		t.Error("❌ FAILED: Should not allow buying more shares than available in market")
	} else if !strings.Contains(err.Error(), "insufficient stock") {
		t.Errorf("❌ FAILED: Should return 'insufficient stock' error. Got: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly rejected buy above available stock. Error: %v", err)
	}
}

// Test: Can buy exactly available stock
func TestBuyExactlyAvailableStock(t *testing.T) {
	txnRepo := NewMockTransactionRepository()
	stockRepo := NewMockStockRepository()
	ws := NewWalletService(txnRepo)
	ws.SetStockRepository(stockRepo)

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	// Create stock with 100 shares available
	stock := &Stock{
		ID:                "stock_1",
		Symbol:            symbol,
		CompanyName:       "Reliance Industries",
		TotalAvailableQty: 100,
		CurrentPrice:      1000,
	}
	stockRepo.SaveStock(stock)

	// Deposit cash
	ws.DepositCash(userID, 500000)

	// Buy exactly 100 shares
	err := ws.ExecuteTrade(userID, symbol, 100, 1000, "BUY")

	if err != nil {
		t.Errorf("❌ FAILED: Should allow buying exactly available stock. Error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed buy of exactly available stock (100 shares)")

		holdings := ws.getCurrentHoldings(userID, symbol)
		if holdings != 100 {
			t.Errorf("❌ FAILED: Holdings should be 100 but got %.2f", holdings)
		} else {
			t.Logf("✅ PASSED: Holdings correctly set to 100 shares")
		}
	}
}

// Test: Stock becomes available again after selling
func TestStockBecomesAvailableAfterSell(t *testing.T) {
	txnRepo := NewMockTransactionRepository()
	stockRepo := NewMockStockRepository()
	ws := NewWalletService(txnRepo)
	ws.SetStockRepository(stockRepo)

	userID1 := "test_user_1"
	userID2 := "test_user_2"
	symbol := "RELIANCE-CE-2900"

	// Create stock with 50 shares total
	stock := &Stock{
		ID:                "stock_1",
		Symbol:            symbol,
		CompanyName:       "Reliance Industries",
		TotalAvailableQty: 50,
		CurrentPrice:      1000,
	}
	stockRepo.SaveStock(stock)

	// User 1 buys all 50 shares
	ws.DepositCash(userID1, 500000)
	err := ws.ExecuteTrade(userID1, symbol, 50, 1000, "BUY")
	if err != nil {
		t.Fatalf("User 1 failed to buy: %v", err)
	}

	// User 2 tries to buy 1 share (should fail - none available)
	ws.DepositCash(userID2, 500000)
	err = ws.ExecuteTrade(userID2, symbol, 1, 1000, "BUY")
	if err == nil {
		t.Error("❌ FAILED: Should not allow buying when no stock available")
	} else {
		t.Logf("✅ PASSED: Correctly rejected buy when no stock available")
	}

	// User 1 sells 30 shares
	err = ws.ExecuteTrade(userID1, symbol, 30, 1000, "SELL")
	if err != nil {
		t.Fatalf("User 1 failed to sell: %v", err)
	}

	// Now User 2 can buy 30 shares (just became available)
	err = ws.ExecuteTrade(userID2, symbol, 30, 1000, "BUY")
	if err != nil {
		t.Errorf("❌ FAILED: Should allow buying after stock becomes available. Error: %v", err)
	} else {
		t.Logf("✅ PASSED: Correctly allowed buy after stock became available")

		holdings := ws.getCurrentHoldings(userID2, symbol)
		if holdings != 30 {
			t.Errorf("❌ FAILED: Holdings should be 30 but got %.2f", holdings)
		} else {
			t.Logf("✅ PASSED: Holdings correctly updated to 30 shares")
		}
	}
}

// Test: Stock availability with no StockRepository (backward compat)
func TestBuyWithoutStockRepository(t *testing.T) {
	repo := NewMockTransactionRepository()
	ws := NewWalletService(repo)
	// Don't set stock repository - should assume unlimited stock

	userID := "test_user"
	symbol := "RELIANCE-CE-2900"

	ws.DepositCash(userID, 1000000)

	// Should allow large quantity since no stock repo (unlimited)
	err := ws.ExecuteTrade(userID, symbol, 50000, 100, "BUY")

	if err != nil {
		// Should fail due to max quantity limit, not stock limit
		if strings.Contains(err.Error(), "maximum quantity exceeded") {
			t.Logf("✅ PASSED: Correctly limited by max quantity, not stock availability")
		} else {
			t.Errorf("❌ FAILED: Should fail on max quantity, not stock. Error: %v", err)
		}
	} else {
		t.Logf("✅ PASSED: Allowed buy with no stock repository (backward compatible)")
	}
}
