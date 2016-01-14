package main

import (
	"fmt"

	"github.com/lytics/slackhook"
)

const (
	slackDefaultEmoji = ":rocket:"
)

func notifySlack(msg string, config slackConfig) {

	if config.Emoji == "" {
		config.Emoji = slackDefaultEmoji
	}

	if config.WebhookURL == "" {
		fmt.Println("WebhookUrl is missing, Slack notification will not be sent.")
		return
	}

	if config.Channel == "" {
		fmt.Println("Channel is missing, Slack notification will not be sent.")
		return
	}

	fmt.Println("Notifying Slack")
	c := slackhook.New(config.WebhookURL)

	m := &slackhook.Message{
		Text:      fmt.Sprintf("```%s```", msg),
		Channel:   config.Channel,
		IconEmoji: config.Emoji,
	}
	c.Send(m)
}
