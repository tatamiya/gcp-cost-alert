package gcp_cost_alert

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/tatamiya/gcp-cost-alert/src/data"
	"github.com/tatamiya/gcp-cost-alert/src/notification"
)

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
