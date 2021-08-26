// alert package implements objects to represent the
// information extracted from the alert massage from Cloud Billing.
//
// AlertDescription summarizes the alert message and
// converts itself into the notification massage sent to Slack.
package alert

import (
	"fmt"

	"github.com/tatamiya/gcp-cost-alert/src/data"
)

// AlertDescription is a summary of the alert Pub/Sub message.
// It consists of the charged cost, exceeded cost,
// and the alert level defined in this system.
type AlertDescription struct {
	Charged  *Cost
	Exceeded *Cost
	AlertLevel
}

// Cost object consists of the cost amount and its currency unit.
type Cost struct {
	Amount       float64
	CurrencyCode string
}

// String method displays its amount with its unit.
func (c *Cost) String() string {
	return fmt.Sprintf("%.2f %s", c.Amount, c.CurrencyCode)
}

// AsMessage creates a notification message sent to Slack.
func (d *AlertDescription) AsMessage() string {
	headLine := d.headline()

	messageBody := fmt.Sprintf(":money_with_wings: GCP 利用額が %s を超過しました！（現在 %s）", d.Exceeded, d.Charged)

	return fmt.Sprintf("%s\n%s", headLine, messageBody)

}

// NewAlertDesctiption constructs a AlertDescription from a parsed Pub/Sub message payload.
func NewAlertDescription(payload *data.PubSubPayload) *AlertDescription {
	exceededThreshold := *payload.AlertThresholdExceeded
	level := newAlertLevel(exceededThreshold)
	unit := payload.CurrencyCode
	charged := payload.CostAmount
	exceeded := payload.BudgetAmount * exceededThreshold

	return &AlertDescription{
		Charged:    &Cost{charged, unit},
		Exceeded:   &Cost{exceeded, unit},
		AlertLevel: level,
	}
}
