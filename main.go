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

type SlackNotifier struct {
	webhookURL string
}

func (s *SlackNotifier) Send(message string) error {
	msg := slack.WebhookMessage{
		Text: message,
	}
	err := slack.PostWebhook(s.webhookURL, &msg)
	return err
}

func NewSlackNotifier() *SlackNotifier {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	return &SlackNotifier{webhookURL: webhookURL}
}

type Notifier interface {
	Send(message string) error
}

func alertNotification(payload *PubSubData, notifier Notifier) error {

	alertDescription := NewAlertDescription(payload)
	if alertDescription.AlertLevel == Unexpected {
		log.Printf("Unexpected AlertLevel! Input payload: %v", payload)
		return fmt.Errorf("Unexpected AlertLevel with charged cost %s!", alertDescription.Charged)
	}
	message := alertDescription.AsMessage()

	err := notifier.Send(message)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil

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
	slackNotifier := NewSlackNotifier()
	err := alertNotification(&alertData, slackNotifier)

	return err
}
