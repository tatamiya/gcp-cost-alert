// notification implements a struct to send an alert notification to Slack.
package notification

import (
	"os"

	"github.com/slack-go/slack"
)

// SlackNotifier is an object to send an alert message to Slack.
// It holds a webhook URL to send the message.
type SlackNotifier struct {
	webhookURL string
}

// Send method sends an alert notification message to Slack.
func (s *SlackNotifier) Send(message string) error {
	msg := slack.WebhookMessage{
		Text: message,
	}
	err := slack.PostWebhook(s.webhookURL, &msg)
	return err
}

// NewSlackNotifier constructs a SlackNotifier object
// using webhook URL set as an environment variable.
func NewSlackNotifier() *SlackNotifier {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	return &SlackNotifier{webhookURL: webhookURL}
}
