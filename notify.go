package main

import (
	"fmt"

	"github.com/lytics/slackhook"
)

func notifySlack(msg string, tkn string, context string) {
	fmt.Println("Notifying Slack")
	c := slackhook.New(tkn)

	m := &slackhook.Message{
		Text:      fmt.Sprintf("```%s```", msg),
		Channel:   context,
		IconEmoji: ":shipit:",
	}
	c.Send(m)
}
