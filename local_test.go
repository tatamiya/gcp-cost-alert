package gcp_cost_alert

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/gcp-cost-notification/src/data"
)

func TestCostAlert(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping")
	}
	sampleData := struct {
		BudgetDisplayName      string    `json:"budgetDisplayName"`
		AlertThresholdExceeded float64   `json:"alertThresholdExceeded"`
		CostAmount             float64   `json:"costAmount"`
		CostIntervalStart      time.Time `json:"costIntervalStart"`
		BudgetAmount           float64   `json:"budgetAmount"`
		BudgetAmountType       string    `json:"budgetAmountType"`
		CurrencyCode           string    `json:"currencyCode"`
	}{
		BudgetDisplayName:      "name-of-budget",
		AlertThresholdExceeded: 1.0,
		CostAmount:             100.01,
		CostIntervalStart:      time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
		BudgetAmount:           100.00,
		BudgetAmountType:       "SPECIFIED_AMOUNT",
		CurrencyCode:           "USD",
	}

	s, err := json.Marshal(sampleData)
	m := pubsub.Message{
		Data: s,
	}
	err = CostAlert(context.Background(), m)
	assert.Nil(t, err)
}

func TestNotificationIsNotSentWhenPayloadIsEmpty(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping")
	}

	sampleData := data.PubSubPayload{}
	s, err := json.Marshal(sampleData)
	m := pubsub.Message{
		Data: s,
	}
	err = CostAlert(context.Background(), m)
	assert.Nil(t, err)
}
