// The JSON payload of a Pub/Sub message.
package data

// PubSubPayload is the JSON object contained in the Data parts of a Pub/Sub message.
// See also https://cloud.google.com/billing/docs/how-to/budgets-programmatic-notifications#notification_format
//
// `AlertThresholdExceeded` is a float pointer
// so that you get nil if the payload does not have this field.
type PubSubPayload struct {
	AlertThresholdExceeded *float64 `json:"alertThresholdExceeded"`
	CostAmount             float64  `json:"costAmount"`
	BudgetAmount           float64  `json:"budgetAmount"`
	CurrencyCode           string   `json:"currencyCode"`
}
