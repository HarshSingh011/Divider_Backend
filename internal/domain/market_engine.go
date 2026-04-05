package domain

import (
	"math/rand"
	"time"
)

type MarketEngine struct {
	prices              []MarketTick
	broadcast           chan []MarketTick
	ticker              *time.Ticker
	stopChan            chan struct{}
	ohlcAggregator      *OHLCAggregator
	alertService        AlertService
	stockRepo           StockRepository
	transactionRepo     TransactionRepository
	initialStockQty     map[string]float64  // Total stock quantity for each symbol
}

func NewMarketEngine() *MarketEngine {
	// Initialize with 6 companies: 4 existing + 2 new (TCS-CE-4500, ITC-PE-2800)
	stockQtyMap := map[string]float64{
		"RELIANCE-CE-2900": 100000,
		"HDFC-PE-1400":     50000,
		"INFY-CE-1500":     75000,
		"TCS-PE-3500":      60000,
		"TCS-CE-4500":      80000,    // NEW: TCS-CE-4500
		"ITC-PE-2800":      90000,    // NEW: ITC-PE-2800
	}

	return &MarketEngine{
		prices: []MarketTick{
			{Symbol: "RELIANCE-CE-2900", CurrentPrice: 45.50, PercentageChange: 0, TotalQuantity: 100000, AvailableQuantity: 100000, HeldQuantity: 0},
			{Symbol: "HDFC-PE-1400", CurrentPrice: 12.20, PercentageChange: 0, TotalQuantity: 50000, AvailableQuantity: 50000, HeldQuantity: 0},
			{Symbol: "INFY-CE-1500", CurrentPrice: 62.75, PercentageChange: 0, TotalQuantity: 75000, AvailableQuantity: 75000, HeldQuantity: 0},
			{Symbol: "TCS-PE-3500", CurrentPrice: 85.30, PercentageChange: 0, TotalQuantity: 60000, AvailableQuantity: 60000, HeldQuantity: 0},
			{Symbol: "TCS-CE-4500", CurrentPrice: 55.20, PercentageChange: 0, TotalQuantity: 80000, AvailableQuantity: 80000, HeldQuantity: 0},   // NEW
			{Symbol: "ITC-PE-2800", CurrentPrice: 38.75, PercentageChange: 0, TotalQuantity: 90000, AvailableQuantity: 90000, HeldQuantity: 0},  // NEW
		},
		broadcast:         make(chan []MarketTick, 10),
		ticker:            time.NewTicker(500 * time.Millisecond),
		stopChan:          make(chan struct{}),
		initialStockQty:   stockQtyMap,
	}
}

func (m *MarketEngine) SetStockRepository(repo StockRepository) {
	m.stockRepo = repo
}

func (m *MarketEngine) SetTransactionRepository(repo TransactionRepository) {
	m.transactionRepo = repo
}

func (m *MarketEngine) SetOHLCAggregator(aggregator *OHLCAggregator) {
	m.ohlcAggregator = aggregator
}

func (m *MarketEngine) SetAlertService(alertService AlertService) {
	m.alertService = alertService
}

func (m *MarketEngine) Start() {
	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.updatePrices()
				m.broadcast <- m.GetCurrentPrices()
			case <-m.stopChan:
				return
			}
		}
	}()
}

func (m *MarketEngine) Stop() {
	m.ticker.Stop()
	m.stopChan <- struct{}{}
}

func (m *MarketEngine) GetCurrentPrices() []MarketTick {
	result := make([]MarketTick, len(m.prices))
	copy(result, m.prices)
	for i := range result {
		result[i].Timestamp = time.Now()
	}
	return result
}

func (m *MarketEngine) Subscribe() chan []MarketTick {
	return m.broadcast
}

func (m *MarketEngine) updatePrices() {
	currentPriceMap := make(map[string]float64)

	// First, calculate held quantities for each symbol from transactions
	heldQtyMap := make(map[string]float64)
	if m.transactionRepo != nil {
		// Get all transactions (pass empty userID to get all)
		allTransactions, err := m.transactionRepo.FindTransactionsByUser("")
		if err == nil {
			for _, txn := range allTransactions {
				if txn.Type == "BUY" {
					heldQtyMap[txn.Symbol] += txn.Quantity
				} else if txn.Type == "SELL" {
					heldQtyMap[txn.Symbol] -= txn.Quantity
				}
			}
		}
	}

	for i := range m.prices {
		oldPrice := m.prices[i].CurrentPrice
		change := (rand.Float64() - 0.5) * 2
		m.prices[i].CurrentPrice += change
		if m.prices[i].CurrentPrice < 0 {
			m.prices[i].CurrentPrice = 0.01
		}
		m.prices[i].PercentageChange = ((m.prices[i].CurrentPrice - oldPrice) / oldPrice) * 100

		// Update held and available quantities
		totalQty := m.initialStockQty[m.prices[i].Symbol]
		heldQty := heldQtyMap[m.prices[i].Symbol]
		if heldQty < 0 {
			heldQty = 0
		}
		availableQty := totalQty - heldQty
		if availableQty < 0 {
			availableQty = 0
		}

		m.prices[i].TotalQuantity = totalQty
		m.prices[i].HeldQuantity = heldQty
		m.prices[i].AvailableQuantity = availableQty

		currentPriceMap[m.prices[i].Symbol] = m.prices[i].CurrentPrice

		if m.ohlcAggregator != nil {
			m.ohlcAggregator.UpdatePriceTick(m.prices[i].Symbol, m.prices[i].CurrentPrice)
		}
	}

	if m.alertService != nil {
		if err := m.alertService.CheckAndTriggerAlerts(currentPriceMap); err != nil {
		}
	}
}
