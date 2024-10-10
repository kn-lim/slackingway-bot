package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

func handler(ctx context.Context, slackRequestBody slackingway.SlackRequestBody) error {
	// Log the request
	log.Printf("Slack Request Body: %v", slackRequestBody)

	// Initialize Slackingway
	s := slackingway.NewSlackingway(&slackRequestBody)

	// Parse the request
	var message slack.Msg
	var err error
	switch slackRequestBody.Type {
	case "slash_command":
		switch slackRequestBody.Command {
		case "/ping":
			message, err = slackingway.Ping()
			if err != nil {
				return err
			}
		case "/delayed-ping":
			message, err = slackingway.DelayedPing(s)
			if err != nil {
				return err
			}
		default:
			log.Printf("Unknown command: %v", slackRequestBody.Command)
			return errors.New("Unknown command")
		}
	default:
		log.Printf("Unknown request type: %v", slackRequestBody.Type)
		return errors.New("Unknown request type")
	}

	// Create the response
	response, err := s.NewResponse(message)
	if err != nil {
		return err
	}

	// Send the response to Slack
	if err := s.SendResponse(response); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
