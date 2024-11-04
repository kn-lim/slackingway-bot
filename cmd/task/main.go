package main

import (
	"context"
	"errors"
	"log"
	"time"

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
	case "view_submission":
		// Update the modal
		var updatedModal slack.ModalViewRequest
		switch slackRequestBody.CallbackID {
		case "echo":
			updatedModal = slackingway.UpdateEchoModal()
		default:
			log.Printf("Unknown CallbackID: %v", slackRequestBody.CallbackID)
			return errors.New("Unknown CallbackID")
		}

		err := s.WriteToHistory()
		if err != nil {
			return err
		}

		_, err = s.APIClient.UpdateView(updatedModal, slackRequestBody.View.ExternalID, slackRequestBody.View.Hash, slackRequestBody.View.ID)
		time.Sleep(time.Second * 2) // Delay to see the updated modal
		if err != nil {
			return err
		}
	default:
		log.Printf("Unknown command: %v", slackRequestBody.Command)
		return errors.New("Unknown command")
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
