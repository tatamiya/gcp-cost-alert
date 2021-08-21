package gcp_cost_alert

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/tatamiya/gcp-cost-notification/src/data"
	"github.com/tatamiya/gcp-cost-notification/src/notification"
)

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

func NewAlertDescription(payload *data.PubSubPayload) *AlertDescription {
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

type Notifier interface {
	Send(message string) error
}

func alertNotification(payload *data.PubSubPayload, notifier Notifier) error {

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

	var payload data.PubSubPayload
	if err := json.Unmarshal(m.Data, &payload); err != nil {
		panic(err)
	}
	if payload.AlertThresholdExceeded == 0.0 {
		// NOTE:
		// When the amount does not exceed the threshold,
		// Pub/Sub message does not have this key.
		return nil
	}
	slackNotifier := notification.NewSlackNotifier()
	err := alertNotification(&payload, slackNotifier)

	return err
}
