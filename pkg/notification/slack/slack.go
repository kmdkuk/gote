package slack

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kmdkuk/gote/cmd/option"
	"go.uber.org/zap"
)

func Post(message string) error {
	logger := zap.L()
	name := "HealthChecker"
	payload := "{\"channel\": \"" + option.Opt.Slack.Channel + "\", \"username\": \"" + name + "\", \"text\": \"" + message + "\", \"icon_emoji\": \":ghost:\"}"
	logger.Info(fmt.Sprintf("payload: %v", payload))
	webhookURL := option.Opt.Slack.WebhookURL
	resp, err := http.PostForm(webhookURL, url.Values{"payload": {payload}})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("%v", resp))
	defer resp.Body.Close()
	return nil
}
