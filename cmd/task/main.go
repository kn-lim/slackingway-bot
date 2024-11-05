package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/kn-lim/slackingway-bot/internal/utils"
)

func handler(ctx context.Context, slackRequestBody slackingway.SlackRequestBody) error {
	// Log the request
	requestString, err := utils.GetStructFields(slackRequestBody)
	if err != nil {
		log.Printf("Error parsing form data: %v", err)
		return err
	}
	log.Printf("Slack Request Body: %v", requestString)

	// Initialize Slackingway
	s := slackingway.NewSlackingway(&slackRequestBody)

	// Parse the request
	var message slack.Msg
	switch slackRequestBody.Type {
	case "slash_command":
		switch slackRequestBody.Command {
		case "/ping":
			err := s.WriteToHistory()
			if err != nil {
				return err
			}

			message, err = slackingway.Ping()
			if err != nil {
				return err
			}
		case "/delayed-ping":
			err := s.WriteToHistory()
			if err != nil {
				return err
			}

			message, err = slackingway.DelayedPing(s)
			if err != nil {
				return err
			}
		default:
			log.Printf("Unknown command: %v", slackRequestBody.Command)
			return errors.New("Unknown command")
		}
	case "view_submission":
		// Log the view
		viewString, err := utils.GetStructFields(slackRequestBody.View)
		if err != nil {
			log.Printf("Error parsing view: %v", err)
			return err
		}
		log.Printf("Slack View: %v", viewString)

		switch slackRequestBody.View.CallbackID {
		case "echo":
			message, err = slackingway.ReturnEcho(s)
			if err != nil {
				return err
			}
		default:
			log.Printf("Unknown CallbackID: %v", slackRequestBody.View.CallbackID)
			return errors.New("Unknown CallbackID")
		}

		err = s.WriteToHistory()
		if err != nil {
			return err
		}
	default:
		log.Printf("Unknown request type: %v", slackRequestBody.Type)
		return errors.New("Unknown request type")
	}

	// Check if message is not empty
	if message.Text != "" {
		// Create the response
		response, err := s.NewResponse(message)
		if err != nil {
			return err
		}

		// Send the response to Slack
		if err := s.SendResponse(response); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
