package gcp_cost_alert

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/slack-go/slack"
)

type PubSubData struct {
	AlertThresholdExceeded float64 `json:"alertThresholdExceeded"`
	CostAmount             float64 `json:"costAmount"`
	BudgetAmount           float64 `json:"budgetAmount"`
	CurrencyCode           string  `json:"currencyCode"`
}

type Cost struct {
	Amount       float64
	CurrencyCode string
}

func (c *Cost) String() string {
	return fmt.Sprintf("%.2f %s", c.Amount, c.CurrencyCode)
}

type AlertLevel int

const (
	Unexpected AlertLevel = iota
	Low
	Middle
	High
)

func newAlertLevel(threshold float64) AlertLevel {
	switch {
	case threshold >= 0.2 && threshold < 0.5:
		return Low
	case threshold >= 0.5 && threshold < 1.0:
		return Middle
	case threshold >= 1.0:
		return High
	default:
		return Unexpected
	}
}

func (l *AlertLevel) headline() string {
	switch *l {
	case Low:
		return ":warning: 注意 :warning:"
	case Middle:
		return ":rotating_light: 警報 :rotating_light:"
	case High:
		return ":fire::fire::fire: 警告 :fire::fire::fire:"
	default:
		return ":asyncparrot:"
	}
}

type AlertDescription struct {
	Charged *Cost
	Budget  *Cost
	AlertLevel
}

func NewAlertDescription(payload *PubSubData) *AlertDescription {
	level := newAlertLevel(payload.AlertThresholdExceeded)
	unit := payload.CurrencyCode
	charged := payload.CostAmount
	budget := payload.BudgetAmount * payload.AlertThresholdExceeded

	return &AlertDescription{
		Charged:    &Cost{charged, unit},
		Budget:     &Cost{budget, unit},
		AlertLevel: level,
	}
}

func (d *AlertDescription) AsMessage() string {
	headLine := d.headline()

	messageBody := fmt.Sprintf(":money_with_wings: GCP 利用額が %s を超過しました！（現在 %s）", d.Budget, d.Charged)

	return fmt.Sprintf("%s\n%s", headLine, messageBody)

}

func generateHeadLine(threshold float64) string {
	switch {
	case threshold <= 0.5:
		return ":warning: 注意 :warning:"
	case threshold < 1.0:
		return ":rotating_light: 警報 :rotating_light:"
	case threshold >= 1.0:
		return ":fire::fire::fire: 警告 :fire::fire::fire:"
	default:
		return ":asyncparrot:"
	}
}

func createNotificationString(alertMessage PubSubData) string {

	headLine := generateHeadLine(alertMessage.AlertThresholdExceeded)

	budget := alertMessage.BudgetAmount * alertMessage.AlertThresholdExceeded
	amount := alertMessage.CostAmount
	currency := alertMessage.CurrencyCode
	messageBody := fmt.Sprintf(":money_with_wings: GCP 利用額が %.2f %s を超過しました！（現在 %.2f %s）", budget, currency, amount, currency)

	return fmt.Sprintf("%s\n%s", headLine, messageBody)
}

func sendMessageToSlack(webhookURL string, messageText string) error {
	msg := slack.WebhookMessage{
		Text: messageText,
	}
	err := slack.PostWebhook(webhookURL, &msg)
	return err
}

func CostAlert(ctx context.Context, m pubsub.Message) error {

	var alertData PubSubData
	if err := json.Unmarshal(m.Data, &alertData); err != nil {
		panic(err)
	}
	if alertData.AlertThresholdExceeded == 0.0 {
		// NOTE:
		// When the amount does not exceed the threshold,
		// Pub/Sub message does not have this key.
		return nil
	}
	messageString := createNotificationString(alertData)
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	err := sendMessageToSlack(webhookURL, messageString)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
