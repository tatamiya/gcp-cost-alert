package data

// The JSON object contained in the Data parts of a Pub/Sub message.
// See also https://cloud.google.com/billing/docs/how-to/budgets-programmatic-notifications#notification_format
type PubSubPayload struct {
	AlertThresholdExceeded *float64 `json:"alertThresholdExceeded"`
	CostAmount             float64  `json:"costAmount"`
	BudgetAmount           float64  `json:"budgetAmount"`
	CurrencyCode           string   `json:"currencyCode"`
}
