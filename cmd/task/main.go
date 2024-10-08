package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

func handler(ctx context.Context, slackRequestBody slackingway.SlackRequestBody) error {
	// Log the request
	log.Printf("Slack Request Body: %v", slackRequestBody)

	switch slackRequestBody.Type {
	case "slash_command":
		switch slackRequestBody.Command {
		case "/delayed-ping":
			return slackingway.ReturnDelayedPing(slackRequestBody.ResponseURL)
		default:
			log.Printf("Unknown command: %v", slackRequestBody.Command)
			return errors.New("Unknown command")
		}
	default:
		log.Printf("Unknown request type: %v", slackRequestBody.Type)
		return errors.New("Unknown request type")
	}
}

func main() {
	lambda.Start(handler)
}
