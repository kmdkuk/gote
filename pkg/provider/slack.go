package provider

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kmdkuk/gote/pkg/option"
	"go.uber.org/zap"
)

func NewSlack(slackOpts option.Slack) Provider {
	return slack{
		channel:    slackOpts.Channel,
		webhookURL: slackOpts.WebhookURL,
	}
}

type slack struct {
	channel    string
	webhookURL string
}

func (s slack) Post(message string) error {
	logger := zap.L()
	name := "HealthChecker"
	payload := "{\"channel\": \"" + s.channel + "\", \"username\": \"" + name + "\", \"text\": \"" + message + "\", \"icon_emoji\": \":ghost:\"}"
	logger.Info(fmt.Sprintf("payload: %v", payload))
	resp, err := http.PostForm(s.webhookURL, url.Values{"payload": {payload}})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("%v", resp))
	defer resp.Body.Close()
	return nil
}
