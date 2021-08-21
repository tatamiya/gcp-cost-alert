package gcp_cost_alert

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/tatamiya/gcp-cost-notification/src/alert"
	"github.com/tatamiya/gcp-cost-notification/src/data"
	"github.com/tatamiya/gcp-cost-notification/src/notification"
)

type Notifier interface {
	Send(message string) error
}

func alertNotification(payload *data.PubSubPayload, notifier Notifier) error {

	alertDescription := alert.NewAlertDescription(payload)
	if alertDescription.AlertLevel == alert.Unexpected {
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
