package domain

import (
	"fmt"
	"sync"
	"time"
)

type AlertServiceImpl struct {
	alertRepo AlertRepository
	mu        sync.RWMutex
}

func NewAlertService(alertRepo AlertRepository) *AlertServiceImpl {
	return &AlertServiceImpl{
		alertRepo: alertRepo,
	}
}

func (as *AlertServiceImpl) CreateAlert(userID, symbol string, price float64, condition string) (*Alert, error) {
	if condition != "ABOVE" && condition != "BELOW" {
		return nil, fmt.Errorf("invalid condition: must be ABOVE or BELOW")
	}

	alert := &Alert{
		ID:             fmt.Sprintf("alert_%d", time.Now().UnixNano()),
		UserID:         userID,
		Symbol:         symbol,
		ThresholdPrice: price,
		Condition:      condition,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := as.alertRepo.SaveAlert(alert); err != nil {
		return nil, fmt.Errorf("failed to save alert: %w", err)
	}

	fmt.Printf("[ALERT] Created alert: %s %s $%.2f (%s)\n", userID, symbol, price, condition)
	return alert, nil
}

func (as *AlertServiceImpl) GetUserAlerts(userID string) ([]Alert, error) {
	return as.alertRepo.FindAlertsByUser(userID)
}

func (as *AlertServiceImpl) CheckAndTriggerAlerts(currentPrices map[string]float64) error {
	activeAlerts, err := as.alertRepo.FindActiveAlerts()
	if err != nil {
		return fmt.Errorf("failed to fetch active alerts: %w", err)
	}

	now := time.Now()

	for _, alert := range activeAlerts {
		currentPrice, exists := currentPrices[alert.Symbol]
		if !exists {
			continue
		}

		triggered := false

		if alert.Condition == "ABOVE" && currentPrice >= alert.ThresholdPrice {
			triggered = true
		} else if alert.Condition == "BELOW" && currentPrice <= alert.ThresholdPrice {
			triggered = true
		}

		if triggered {
			fmt.Printf("[ALERT TRIGGERED] %s: %s crossed %.2f at %.2f\n",
				alert.UserID, alert.Symbol, alert.ThresholdPrice, currentPrice)

			alert.IsActive = false
			alert.TriggeredAt = &now
			alert.UpdatedAt = now

			if err := as.alertRepo.UpdateAlert(&alert); err != nil {
				fmt.Printf("[ALERT] Error updating alert: %v\n", err)
			}
		}
	}

	return nil
}
