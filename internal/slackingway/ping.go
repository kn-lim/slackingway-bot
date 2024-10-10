package slackingway

import (
	"log"

	"github.com/slack-go/slack"
)

// Ping responds to a ping request
func (s *Slackingway) Ping() (slack.Msg, error) {
	log.Printf("/ping received")

	return slack.Msg{Text: "Pong!"}, nil
}
