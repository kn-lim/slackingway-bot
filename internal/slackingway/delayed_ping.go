package slackingway

import (
	"log"
	"time"

	"github.com/slack-go/slack"
)

var (
	PingDelay = 5 * time.Second
)

// DelayedPing sends a delayed ping response to the user
func DelayedPing(s Slackingway) (slack.Msg, error) {
	log.Printf("/delayed-ping received")

	// Create Slack message
	message := slack.Msg{
		Text:            "...",
		ReplaceOriginal: true,
	}

	// Create the response
	response, err := s.NewResponse(message)
	if err != nil {
		return slack.Msg{}, err
	}

	// Send the response to Slack
	if err := s.SendResponse(response); err != nil {
		return slack.Msg{}, err
	}

	// Wait for the delay
	time.Sleep(PingDelay)

	return slack.Msg{
		Text:            "Delayed Pong!",
		ReplaceOriginal: true,
	}, nil
}
