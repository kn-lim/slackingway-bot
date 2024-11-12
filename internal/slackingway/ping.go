package slackingway

import (
	"github.com/slack-go/slack"
)

// Ping responds to a ping request
func Ping() (slack.Msg, error) {
	return slack.Msg{Text: "Pong!"}, nil
}
