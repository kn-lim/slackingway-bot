package slackingway

import (
	"log"

	"github.com/slack-go/slack"
)

// Ping responds to a ping request
func Ping(responseURL string) error {
	log.Printf("/ping received")

	// Create the response
	message := slack.Msg{Text: "Pong!"}
	response, err := NewResponse(responseURL, message)
	if err != nil {
		return err
	}

	// Send the response to Slack
	if err := SendResponse(response); err != nil {
		return err
	}

	return nil
}
