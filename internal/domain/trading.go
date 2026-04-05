package domain

import "time"

type Candle struct {
	ID           string    `json:"id"`
	Symbol       string    `json:"symbol"`
	Open         float64   `json:"open"`
	High         float64   `json:"high"`
	Low          float64   `json:"low"`
	Close        float64   `json:"close"`
	Volume       int       `json:"volume"`
	TimeframeKey string    `json:"timeframe_key"`
	Timestamp    time.Time `json:"timestamp"`
}

type CandleRepository interface {
	SaveCandle(candle *Candle) error
	GetCandles(symbol string, limit int) ([]Candle, error)
	GetCandlesByTimeRange(symbol string, from, to time.Time) ([]Candle, error)
}

type Alert struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Symbol         string     `json:"symbol"`
	ThresholdPrice float64    `json:"threshold_price"`
	Condition      string     `json:"condition"`
	IsActive       bool       `json:"is_active"`
	TriggeredAt    *time.Time `json:"triggered_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type AlertRepository interface {
	SaveAlert(alert *Alert) error
	FindAlertsByUser(userID string) ([]Alert, error)
	FindActiveAlerts() ([]Alert, error)
	UpdateAlert(alert *Alert) error
	DeleteAlert(alertID string) error
}

// Stock represents a security listed on the market
type Stock struct {
	ID                 string    `json:"id"`
	Symbol             string    `json:"symbol"`
	CompanyName        string    `json:"company_name"`
	TotalAvailableQty  float64   `json:"total_available_qty"`   // Total shares available in market
	CurrentPrice       float64   `json:"current_price"`
	MarketCap          float64   `json:"market_cap"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type StockRepository interface {
	SaveStock(stock *Stock) error
	GetStock(symbol string) (*Stock, error)
	UpdateStock(stock *Stock) error
	GetAllStocks() ([]Stock, error)
}

type Transaction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Symbol    string    `json:"symbol"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	Fee       float64   `json:"fee"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	Remarks   string    `json:"remarks"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionRepository interface {
	SaveTransaction(transaction *Transaction) error
	FindTransactionsByUser(userID string) ([]Transaction, error)
	FindTransactionsBySymbol(userID string, symbol string) ([]Transaction, error)
	GetUserBalance(userID string) (float64, error)
}

type WalletSnapshot struct {
	UserID         string              `json:"user_id"`
	TotalBalance   float64             `json:"total_balance"`
	AvailableCash  float64             `json:"available_cash"`
	InvestedAmount float64             `json:"invested_amount"`
	Positions      map[string]Position `json:"positions"`
	LastUpdated    time.Time           `json:"last_updated"`
}

type Position struct {
	Symbol        string  `json:"symbol"`
	Quantity      float64 `json:"quantity"`
	AverageCost   float64 `json:"average_cost"`
	CurrentPrice  float64 `json:"current_price"`
	UnrealizedPnL float64 `json:"unrealized_pnl"`
	Percentage    float64 `json:"percentage"`
}

type WalletService interface {
	ExecuteTrade(userID, symbol string, quantity, price float64, tradeType string) error
	GetWalletSnapshot(userID string) (*WalletSnapshot, error)
	DepositCash(userID string, amount float64) error
	WithdrawCash(userID string, amount float64) error
}

type AlertService interface {
	CreateAlert(userID, symbol string, price float64, condition string) (*Alert, error)
	GetUserAlerts(userID string) ([]Alert, error)
	CheckAndTriggerAlerts(currentPrices map[string]float64) error
}

type OHLCService interface {
	UpdatePriceTick(symbol string, price float64) error
	GetCandles(symbol string, limit int) ([]Candle, error)
}
