package gcp_cost_alert

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/tatamiya/gcp-cost-alert/src"
	"github.com/tatamiya/gcp-cost-alert/src/data"
	"github.com/tatamiya/gcp-cost-alert/src/notification"
)

func CostAlert(ctx context.Context, m pubsub.Message) error {

	var payload data.PubSubPayload
	if err := json.Unmarshal(m.Data, &payload); err != nil {
		panic(err)
	}
	if payload.AlertThresholdExceeded == nil {
		// NOTE:
		// Even if the charged amount does not exceed the threshold,
		// Pub/Sub message is sent 2~3 times per hour.
		// In this case, the payload does not have `alertThresholdExceeded` field.
		return nil
	}
	slackNotifier := notification.NewSlackNotifier()
	err := src.AlertNotification(&payload, slackNotifier)

	return err
}
