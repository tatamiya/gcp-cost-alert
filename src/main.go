// The core function of the notification system.
//
// This is called in the CostAlert function
// after a JSON payload of a Pub/Sub message is parsed.
package src

import (
	"fmt"
	"log"

	"github.com/tatamiya/gcp-cost-alert/src/alert"
	"github.com/tatamiya/gcp-cost-alert/src/data"
)

// Notifier interface receives a notification message
// and send it to a channel such as Slack.
type Notifier interface {
	Send(message string) error
}

// AlertNotification is the core function of this notification system.
// It receives a parsed Pub/Sub message payload
// and send a notification message via designated notifier.
func AlertNotification(payload *data.PubSubPayload, notifier Notifier) error {

	alertDescription := alert.NewAlertDescription(payload)
	if alertDescription.AlertLevel == alert.Unexpected {
		log.Printf("Unexpected AlertLevel! Input payload: %v", payload)
		return fmt.Errorf("Unexpected AlertLevel with AlertThresholdExceeded=%.2f!", *payload.AlertThresholdExceeded)
	}
	message := alertDescription.AsMessage()

	err := notifier.Send(message)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil

}
