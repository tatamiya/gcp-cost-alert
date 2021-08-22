package gcp_cost_alert

import (
	"fmt"
	"log"

	"github.com/tatamiya/gcp-cost-alert/src/alert"
	"github.com/tatamiya/gcp-cost-alert/src/data"
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
