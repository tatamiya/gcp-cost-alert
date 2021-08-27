package src

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/gcp-cost-alert/src/data"
)

type notifierStub struct {
	SentMessage string
	Err         error
}

func (n *notifierStub) Send(message string) error {
	if n.Err != nil {
		return n.Err
	}
	n.SentMessage = message
	return nil
}

func TestRunWholeProcessCorrectly(t *testing.T) {
	exceededThreshold := 1.0
	inputAlertPayload := data.PubSubPayload{
		AlertThresholdExceeded: &exceededThreshold,
		CostAmount:             100.01,
		BudgetAmount:           100.00,
		CurrencyCode:           "JPY",
	}
	testNotifier := notifierStub{}
	expectedMessage := ":fire::fire::fire: 警告 :fire::fire::fire:\n:money_with_wings: GCP 利用額が 100.00 JPY を超過しました！（現在 100.01 JPY）"
	err := AlertNotification(&inputAlertPayload, &testNotifier)
	actualMessage := testNotifier.SentMessage

	assert.Nil(t, err)
	assert.EqualValues(t, expectedMessage, actualMessage)
}

func TestReturnErrorWhenZeroAlertThresholdIsInput(t *testing.T) {
	exceededThreshold := 0.0
	inputAlertPayload := data.PubSubPayload{
		AlertThresholdExceeded: &exceededThreshold,
		CostAmount:             0.00,
		BudgetAmount:           0.00,
		CurrencyCode:           "JPY",
	}
	testNotifier := notifierStub{}
	err := AlertNotification(&inputAlertPayload, &testNotifier)

	assert.NotNil(t, err)
	assert.EqualValues(
		t,
		"Unexpected AlertLevel with AlertThresholdExceeded=0.00!",
		err.Error(),
	)
}

func TestReturnErrorWhenSlackNotificationFailed(t *testing.T) {
	exceededThreshold := 1.0
	inputAlertPayload := data.PubSubPayload{
		AlertThresholdExceeded: &exceededThreshold,
		CostAmount:             100.01,
		BudgetAmount:           100.00,
		CurrencyCode:           "JPY",
	}
	testNotifier := notifierStub{Err: fmt.Errorf("Something Wrong!")}
	err := AlertNotification(&inputAlertPayload, &testNotifier)

	assert.NotNil(t, err)
}
