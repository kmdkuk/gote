package slack

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func Post(message string) error {
	name := "HealthChecker"
	// TODO: select channel
	channel := os.Getenv("SLACK_CHANNEL")
	payload := "{\"channel\": \"" + channel + "\", \"username\": \"" + name + "\", \"text\": \"" + message + "\", \"icon_emoji\": \":ghost:\"}"
	log.Println("payload:", payload)
	webhookURL := os.Getenv("WEBHOOK_URL")
	log.Println("webhookURL:", webhookURL)
	resp, err := http.PostForm(webhookURL, url.Values{"payload": {payload}})
	if err != nil {
		return err
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	return nil
}
