package slackingway

import (
	"log"
	"time"

	"github.com/slack-go/slack"
)

const (
	PING_DELAY = 5 * time.Second
)

// DelayedPing sends a delayed ping response to the user
func DelayedPing(responseURL string) error {
	log.Printf("/delayed-ping received")

	// Create Slack message
	message := slack.Msg{
		Text:            "...",
		ReplaceOriginal: true,
	}

	// Create the response
	response, err := NewResponse(responseURL, message)
	if err != nil {
		return err
	}

	// Send the response to Slack
	if err := SendResponse(response); err != nil {
		return err
	}

	// Wait for the delay
	time.Sleep(PING_DELAY)

	// Create Slack message
	message = slack.Msg{
		Text:            "Delayed Pong!",
		ReplaceOriginal: true,
	}

	// Create the response
	response, err = NewResponse(responseURL, message)
	if err != nil {
		return err
	}

	// Send the response to Slack
	if err := SendResponse(response); err != nil {
		return err
	}

	return nil
}
