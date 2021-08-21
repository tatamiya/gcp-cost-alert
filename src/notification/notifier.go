package notification

import (
	"os"

	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	webhookURL string
}

func (s *SlackNotifier) Send(message string) error {
	msg := slack.WebhookMessage{
		Text: message,
	}
	err := slack.PostWebhook(s.webhookURL, &msg)
	return err
}

func NewSlackNotifier() *SlackNotifier {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	return &SlackNotifier{webhookURL: webhookURL}
}
