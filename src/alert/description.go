package alert

import (
	"fmt"

	"github.com/tatamiya/gcp-cost-alert/src/data"
)

type AlertDescription struct {
	Charged *Cost
	Budget  *Cost
	AlertLevel
}

type Cost struct {
	Amount       float64
	CurrencyCode string
}

func (c *Cost) String() string {
	return fmt.Sprintf("%.2f %s", c.Amount, c.CurrencyCode)
}

func (d *AlertDescription) AsMessage() string {
	headLine := d.headline()

	messageBody := fmt.Sprintf(":money_with_wings: GCP 利用額が %s を超過しました！（現在 %s）", d.Budget, d.Charged)

	return fmt.Sprintf("%s\n%s", headLine, messageBody)

}

func NewAlertDescription(payload *data.PubSubPayload) *AlertDescription {
	exceededThreshold := *payload.AlertThresholdExceeded
	level := newAlertLevel(exceededThreshold)
	unit := payload.CurrencyCode
	charged := payload.CostAmount
	budget := payload.BudgetAmount * exceededThreshold

	return &AlertDescription{
		Charged:    &Cost{charged, unit},
		Budget:     &Cost{budget, unit},
		AlertLevel: level,
	}
}
