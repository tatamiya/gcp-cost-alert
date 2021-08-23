package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/gcp-cost-alert/src/data"
)

func TestCreateHighAlertDescriptionCorrectly(t *testing.T) {
	alertThreshold := 1.0
	inputAlertPayload := data.PubSubPayload{
		AlertThresholdExceeded: &alertThreshold,
		CostAmount:             100.01,
		BudgetAmount:           100.00,
		CurrencyCode:           "JPY",
	}
	expectedAlertDescription := &AlertDescription{
		Charged:    &Cost{100.01, "JPY"},
		Budget:     &Cost{100.00, "JPY"},
		AlertLevel: High,
	}
	actualAlertDescription := NewAlertDescription(&inputAlertPayload)
	assert.EqualValues(t, expectedAlertDescription, actualAlertDescription)
}

func TestCreateHighAlertMessageCorrectly(t *testing.T) {
	inputAlertDescription := &AlertDescription{
		Charged:    &Cost{100.01, "JPY"},
		Budget:     &Cost{100.00, "JPY"},
		AlertLevel: High,
	}

	expectedMessage := ":fire::fire::fire: 警告 :fire::fire::fire:\n:money_with_wings: GCP 利用額が 100.00 JPY を超過しました！（現在 100.01 JPY）"
	actualOutput := inputAlertDescription.AsMessage()
	assert.EqualValues(t, expectedMessage, actualOutput)
}

func TestCreateMiddleAlertDescriptionCorrectly(t *testing.T) {
	alertThreshold := 0.5
	inputAlertPayload := data.PubSubPayload{
		AlertThresholdExceeded: &alertThreshold,
		CostAmount:             60.24,
		BudgetAmount:           100.00,
		CurrencyCode:           "JPY",
	}
	expectedAlertDescription := &AlertDescription{
		Charged:    &Cost{60.24, "JPY"},
		Budget:     &Cost{50.00, "JPY"},
		AlertLevel: Middle,
	}
	actualAlertDescription := NewAlertDescription(&inputAlertPayload)
	assert.EqualValues(t, expectedAlertDescription, actualAlertDescription)
}

func TestCreateMiddleAlertMessageCorrectly(t *testing.T) {
	inputAlertDescription := &AlertDescription{
		Charged:    &Cost{60.24, "JPY"},
		Budget:     &Cost{50.00, "JPY"},
		AlertLevel: Middle,
	}

	expectedMessage := ":rotating_light: 警報 :rotating_light:\n:money_with_wings: GCP 利用額が 50.00 JPY を超過しました！（現在 60.24 JPY）"

	actualOutput := inputAlertDescription.AsMessage()
	assert.EqualValues(t, expectedMessage, actualOutput)
}

func TestCreateLowAlertDescriptionCorrectly(t *testing.T) {
	alertThreshold := 0.2
	inputAlertPayload := data.PubSubPayload{
		AlertThresholdExceeded: &alertThreshold,
		CostAmount:             31.42,
		BudgetAmount:           100.00,
		CurrencyCode:           "JPY",
	}
	expectedAlertDescription := &AlertDescription{
		Charged:    &Cost{31.42, "JPY"},
		Budget:     &Cost{20.00, "JPY"},
		AlertLevel: Low,
	}
	actualAlertDescription := NewAlertDescription(&inputAlertPayload)
	assert.EqualValues(t, expectedAlertDescription, actualAlertDescription)
}

func TestCreateLowAlertMessageCorrectly(t *testing.T) {
	inputAlertDescription := &AlertDescription{
		Charged:    &Cost{31.42, "JPY"},
		Budget:     &Cost{20.00, "JPY"},
		AlertLevel: Low,
	}

	expectedMessage := ":warning: 注意 :warning:\n:money_with_wings: GCP 利用額が 20.00 JPY を超過しました！（現在 31.42 JPY）"

	actualOutput := inputAlertDescription.AsMessage()
	assert.EqualValues(t, expectedMessage, actualOutput)
}

func TestCreateUnexpectedAlertDescriptionCorrectly(t *testing.T) {
	alertThreshold := 0.0
	inputAlertPayload := data.PubSubPayload{
		AlertThresholdExceeded: &alertThreshold,
		CostAmount:             0.00,
		BudgetAmount:           100.00,
		CurrencyCode:           "JPY",
	}
	expectedAlertDescription := &AlertDescription{
		Charged:    &Cost{0.00, "JPY"},
		Budget:     &Cost{0.00, "JPY"},
		AlertLevel: Unexpected,
	}
	actualAlertDescription := NewAlertDescription(&inputAlertPayload)
	assert.EqualValues(t, expectedAlertDescription, actualAlertDescription)
}
