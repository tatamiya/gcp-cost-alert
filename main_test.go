package gcp_cost_alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNotificationString(t *testing.T) {
	inputAlertPayload := PubSubData{
		AlertThresholdExceeded: 1.0,
		CostAmount:             100.01,
		BudgetAmount:           100.00,
		CurrencyCode:           "USD",
	}

	expectedOutput := ":fire::fire::fire: 警告 :fire::fire::fire:\n:money_with_wings: GCP 利用額が 100.00 USD を超過しました！（現在 100.01 USD）"

	actualOutput := createNotificationString(inputAlertPayload)
	assert.EqualValues(t, expectedOutput, actualOutput)
}

func TestCreateNotificationStringAlertMiddle(t *testing.T) {
	inputAlertPayload := PubSubData{
		AlertThresholdExceeded: 0.9,
		CostAmount:             90.24,
		BudgetAmount:           100.00,
		CurrencyCode:           "JPY",
	}

	expectedOutput := ":rotating_light: 警報 :rotating_light:\n:money_with_wings: GCP 利用額が 90.00 JPY を超過しました！（現在 90.24 JPY）"

	actualOutput := createNotificationString(inputAlertPayload)
	assert.EqualValues(t, expectedOutput, actualOutput)
}

func TestCreateNotificationStringAlertLow(t *testing.T) {
	inputAlertPayload := PubSubData{
		AlertThresholdExceeded: 0.40,
		CostAmount:             1012.01,
		BudgetAmount:           2500.00,
		CurrencyCode:           "JPY",
	}

	expectedOutput := ":warning: 注意 :warning:\n:money_with_wings: GCP 利用額が 1000.00 JPY を超過しました！（現在 1012.01 JPY）"

	actualOutput := createNotificationString(inputAlertPayload)
	assert.EqualValues(t, expectedOutput, actualOutput)
}

// func TestSlackPost(t *testing.T) {
// 	inputURL := os.Getenv("SLACK_WEBHOOK_URL")
// 	inputMessage := "test\nこれはテスト投稿です。"
//
// 	err := sendMessageToSlack(inputURL, inputMessage)
// 	assert.Nil(t, err)
// }
