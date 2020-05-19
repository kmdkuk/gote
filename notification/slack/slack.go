package slack

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/kmdkuk/gote/cmd/option"
)

func Post(message string) error {
	name := "HealthChecker"
	payload := "{\"channel\": \"" + option.Opt.Slack.Channel + "\", \"username\": \"" + name + "\", \"text\": \"" + message + "\", \"icon_emoji\": \":ghost:\"}"
	log.Println("payload:", payload)
	webhookURL := option.Opt.Slack.WebhookURL
	log.Println("webhookURL:", webhookURL)
	resp, err := http.PostForm(webhookURL, url.Values{"payload": {payload}})
	if err != nil {
		return err
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	return nil
}
