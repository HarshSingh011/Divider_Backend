package domain

import (
	"math/rand"
	"time"
)

// MarketEngine handles core market simulation logic (no external dependencies)
type MarketEngine struct {
	prices    []MarketTick
	broadcast chan []MarketTick
	ticker    *time.Ticker
	stopChan  chan struct{}
}

// NewMarketEngine creates and returns a new market engine
func NewMarketEngine() *MarketEngine {
	return &MarketEngine{
		prices: []MarketTick{
			{Symbol: "RELIANCE-CE-2900", CurrentPrice: 45.50, PercentageChange: 0},
			{Symbol: "HDFC-PE-1400", CurrentPrice: 12.20, PercentageChange: 0},
			{Symbol: "INFY-CE-1500", CurrentPrice: 62.75, PercentageChange: 0},
			{Symbol: "TCS-PE-3500", CurrentPrice: 85.30, PercentageChange: 0},
		},
		broadcast: make(chan []MarketTick, 10),
		ticker:    time.NewTicker(500 * time.Millisecond),
		stopChan:  make(chan struct{}),
	}
}

// Start runs the market simulation
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

// Stop halts the market engine
func (m *MarketEngine) Stop() {
	m.ticker.Stop()
	m.stopChan <- struct{}{}
}

// GetCurrentPrices returns a copy of current prices
func (m *MarketEngine) GetCurrentPrices() []MarketTick {
	result := make([]MarketTick, len(m.prices))
	copy(result, m.prices)
	for i := range result {
		result[i].Timestamp = time.Now()
	}
	return result
}

// Subscribe returns a channel for market updates
func (m *MarketEngine) Subscribe() chan []MarketTick {
	return m.broadcast
}

// updatePrices applies random changes to stock prices
func (m *MarketEngine) updatePrices() {
	for i := range m.prices {
		oldPrice := m.prices[i].CurrentPrice
		change := (rand.Float64() - 0.5) * 2 // -1.0 to +1.0
		m.prices[i].CurrentPrice += change
		if m.prices[i].CurrentPrice < 0 {
			m.prices[i].CurrentPrice = 0.01 // Prevent negative prices
		}
		m.prices[i].PercentageChange = ((m.prices[i].CurrentPrice - oldPrice) / oldPrice) * 100
	}
}
