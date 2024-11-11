package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/kn-lim/slackingway-bot/internal/utils"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Set debug mode
	DEBUG := os.Getenv("DEBUG") == "true"

	// Log the request
	if DEBUG {
		log.Printf("Request Headers: %v", request.Headers)
		log.Printf("Request Body: %v", request.Body)
	}

	// Validate the request
	if err := slackingway.ValidateRequest(request); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, err
	}

	// Create a new SlackRequestBody
	var slackRequestBody slackingway.SlackRequestBody
	switch request.Headers["Content-Type"] {
	// Slash Command & Interactive Components
	case "application/x-www-form-urlencoded":
		if err := slackRequestBody.ParseTimestamp(request.Headers["X-Slack-Request-Timestamp"]); err != nil {
			log.Printf("Error parsing timestamp: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}

		// Parse the form data
		formData, err := url.ParseQuery(request.Body)
		if err != nil {
			log.Printf("Error parsing request body: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}

		// Check for payload
		if formData.Get("payload") != "" {
			if DEBUG {
				// Log the payload
				payloadString, err := utils.PrintPayloadFields(formData.Get("payload"))
				if err != nil {
					log.Printf("Error parsing payload: %v", err)
					return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
				}
				log.Printf("Slack Payload: %v", payloadString)
			}

			// Parse the payload
			if err := slackRequestBody.ParsePayload(formData.Get("payload")); err != nil {
				log.Printf("Error parsing payload: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
		} else {
			// Parse the slash command
			if err := slackRequestBody.ParseSlashCommand(formData); err != nil {
				log.Printf("Error parsing slash command: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
		}
	// Any other Slack request
	case "application/json":
		if err := json.Unmarshal([]byte(request.Body), &slackRequestBody); err != nil {
			log.Printf("Error parsing request body: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
	default:
		log.Printf("Unknown content type: %s", request.Headers["Content-Type"])
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, errors.New("Unknown content type")
	}

	// Initialize Slackingway
	s := slackingway.NewSlackingway(&slackRequestBody)

	// Handle the request
	switch s.SlackRequestBody.Type {
	// Challenge request
	case "url_verification":
		headers := make(map[string]string)
		headers["Content-Type"] = "text/plain"
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: http.StatusOK,
			Body:       s.SlackRequestBody.Challenge,
		}, nil
	// Events
	case "event_callback":
		switch s.SlackRequestBody.Event.Type {
		case "app_home_opened":
			if err := slackingway.HomeTab(s); err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Event.Type, err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
		default:
			log.Printf("Unknown event type: %s", s.SlackRequestBody.Event.Type)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, errors.New("Unknown event type")
		}
	// Slash command
	case "slash_command":
		switch s.SlackRequestBody.Command {
		// Check for slash commands with modals
		case "/echo":
			if err := slackingway.Echo(s); err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
		case "/menu":
			if !slackingway.ValidateAdminRole(s.SlackRequestBody.UserID) {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusOK,
					Body:       fmt.Sprintf("You do not have permission to use `%s`", s.SlackRequestBody.Command),
				}, nil
			}

			if err := slackingway.Menu(s); err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
		// All other commands
		default:
			// Invoke the task function with the SlackRequestBody as the payload
			if err := utils.InvokeTaskFunction(ctx, *s.SlackRequestBody); err != nil {
				log.Printf("Error invoking task function: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
			}
		}
	// Modal submission
	case "view_submission":
		// Invoke the task function with the SlackRequestBody as the payload
		if err := utils.InvokeTaskFunction(ctx, *s.SlackRequestBody); err != nil {
			log.Printf("Error invoking task function: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
	// Block Actions
	case "block_actions":
		// Do nothing for now
	default:
		log.Printf("Unknown request type: %s", s.SlackRequestBody.Type)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, errors.New("Unknown request type")
	}

	// Return an empty response
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
