package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"

	"github.com/kn-lim/slackingway-bot/internal/dice"
	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/kn-lim/slackingway-bot/internal/utils"
)

func handler(ctx context.Context, request slackingway.SlackRequestBody) error {
	// Set debug mode
	DEBUG := os.Getenv("DEBUG") == "true"

	// Log the request
	if DEBUG {
		requestString, err := utils.PrintStructFields(request)
		if err != nil {
			log.Printf("Error parsing form data: %v", err)
			return err
		}
		log.Printf("Slack Request Body: %v", requestString)
	}

	// Initialize Slackingway
	s := slackingway.NewSlackingway(&request)

	// Handle the request
	var message slack.Msg
	var blocks []slack.Block
	var err error
	switch request.Type {
	// Slash Command
	case "slash_command":
		switch request.Command {
		case "/delayed-ping":
			if err := s.WriteToHistory(); err != nil {
				log.Printf("Error writing to history: %v", err)
				return err
			}

			message, err = slackingway.DelayedPing(s)
			if err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return err
			}
		case "/ping":
			if err := s.WriteToHistory(); err != nil {
				log.Printf("Error writing to history: %v", err)
				return err
			}

			message, err = slackingway.Ping()
			if err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return err
			}
		case "/roll":
			if err := s.WriteToHistory(); err != nil {
				log.Printf("Error writing to history: %v", err)
				return err
			}

			resultString, resultInt, err := dice.Roll(s.SlackRequestBody.Text)
			if err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return err
			}

			message.Text = fmt.Sprintf("%s: %s = **%d**", s.SlackRequestBody.Text, resultString, resultInt)
		default:
			log.Printf("Unknown command: %v", request.Command)
			return errors.New("unknown command")
		}

		// Send the response to Slack if there is a message
		if message.Text != "" && len(blocks) == 0 {
			// Create the response
			response, err := s.NewResponse(message)
			if err != nil {
				log.Printf("Error creating response: %v", err)
				return err
			}

			// Send the response
			if err := s.SendResponse(response); err != nil {
				log.Printf("Error sending response: %v", err)
				return err
			}
		}
	// Modal Submission
	case "view_submission":
		switch s.SlackRequestBody.View.CallbackID {
		case "/echo":
			if err := s.WriteToHistory(); err != nil {
				log.Printf("Error writing to history: %v", err)
				return err
			}

			message, err = slackingway.ReceivedEcho(s)
			if err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return err
			}
		case "/menu":
			if err := s.WriteToHistory(); err != nil {
				log.Printf("Error writing to history: %v", err)
				return err
			}

			blocks, err = slackingway.ReceivedMenu(s)
			if err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return err
			}
		default:
			log.Printf("Unknown CallbackID: %v", s.SlackRequestBody.View.CallbackID)
			return errors.New("unknown CallbackID")
		}

		// Send the response to Slack if there is a message
		if message.Text != "" && len(blocks) == 0 {
			if err := s.SendTextMessage(message.Text, os.Getenv("SLACK_OUTPUT_CHANNEL_ID")); err != nil {
				log.Printf("Error sending message: %v", err)
				return err
			}
		} else {
			if err := s.SendBlockMessage(blocks, os.Getenv("SLACK_OUTPUT_CHANNEL_ID")); err != nil {
				log.Printf("Error sending message: %v", err)
				return err
			}
		}
	default:
		log.Printf("Unknown request type: %v", request.Type)
		return errors.New("unknown request type")
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
