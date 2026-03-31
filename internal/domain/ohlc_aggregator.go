package domain

import (
	"fmt"
	"sync"
	"time"
)

type OHLCAggregator struct {
	candleRepo     CandleRepository
	currentCandles map[string]*CandleBuffer
	mu             sync.RWMutex
	ticker         *time.Ticker
	stopChan       chan struct{}
}

type CandleBuffer struct {
	Symbol    string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int
	Count     int
	StartTime time.Time
}

func NewOHLCAggregator(candleRepo CandleRepository) *OHLCAggregator {
	return &OHLCAggregator{
		candleRepo:     candleRepo,
		currentCandles: make(map[string]*CandleBuffer),
		ticker:         time.NewTicker(1 * time.Minute),
		stopChan:       make(chan struct{}),
	}
}

func (o *OHLCAggregator) Start() {
	go func() {
		for {
			select {
			case <-o.ticker.C:
				o.aggregateAndSave()
			case <-o.stopChan:
				return
			}
		}
	}()
}

func (o *OHLCAggregator) Stop() {
	o.ticker.Stop()
	o.stopChan <- struct{}{}
}

func (o *OHLCAggregator) UpdatePriceTick(symbol string, price float64) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, exists := o.currentCandles[symbol]; !exists {
		o.currentCandles[symbol] = &CandleBuffer{
			Symbol:    symbol,
			Open:      price,
			High:      price,
			Low:       price,
			Close:     price,
			StartTime: time.Now(),
		}
	} else {
		buffer := o.currentCandles[symbol]
		buffer.Close = price
		if price > buffer.High {
			buffer.High = price
		}
		if price < buffer.Low {
			buffer.Low = price
		}
		buffer.Count++
	}

	return nil
}

func (o *OHLCAggregator) aggregateAndSave() {
	o.mu.Lock()
	defer o.mu.Unlock()

	now := time.Now()
	timeframeKey := now.Format("2006-01-02-15:04")

	for symbol, buffer := range o.currentCandles {
		candle := &Candle{
			ID:           fmt.Sprintf("%s_%s", symbol, timeframeKey),
			Symbol:       symbol,
			Open:         buffer.Open,
			High:         buffer.High,
			Low:          buffer.Low,
			Close:        buffer.Close,
			Volume:       buffer.Count,
			TimeframeKey: timeframeKey,
			Timestamp:    now,
		}

		if err := o.candleRepo.SaveCandle(candle); err != nil {
			fmt.Printf("[OHLC] Error saving candle for %s: %v\n", symbol, err)
		}

		o.currentCandles[symbol] = &CandleBuffer{
			Symbol:    symbol,
			Open:      buffer.Close,
			High:      buffer.Close,
			Low:       buffer.Close,
			Close:     buffer.Close,
			StartTime: now,
		}

		fmt.Printf("[OHLC] Saved candle for %s: O=%.2f H=%.2f L=%.2f C=%.2f\n",
			symbol, candle.Open, candle.High, candle.Low, candle.Close)
	}
}

func (o *OHLCAggregator) GetCandles(symbol string, limit int) ([]Candle, error) {
	return o.candleRepo.GetCandles(symbol, limit)
}
