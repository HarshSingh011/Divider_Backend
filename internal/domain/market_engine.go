package domain

import (
	"math/rand"
	"time"
)

type MarketEngine struct {
	prices         []MarketTick
	broadcast      chan []MarketTick
	ticker         *time.Ticker
	stopChan       chan struct{}
	ohlcAggregator *OHLCAggregator
	alertService   AlertService
}

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

	for i := range m.prices {
		oldPrice := m.prices[i].CurrentPrice
		change := (rand.Float64() - 0.5) * 2
		m.prices[i].CurrentPrice += change
		if m.prices[i].CurrentPrice < 0 {
			m.prices[i].CurrentPrice = 0.01
		}
		m.prices[i].PercentageChange = ((m.prices[i].CurrentPrice - oldPrice) / oldPrice) * 100

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
