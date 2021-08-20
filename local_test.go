package gcp_cost_alert

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestCostAlert(t *testing.T) {
	sampleData := PubSubData{
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
