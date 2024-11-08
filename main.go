package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/kn-lim/slackingway-bot/internal/utils"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	// Set debug mode
	DEBUG := os.Getenv("DEBUG") == "true"

	// Log the request
	if DEBUG {
		log.Printf("Request Headers: %v", request.Headers)
		log.Printf("Request Body: %v", request.Body)
	}

	// Validate the request
	if err := slackingway.VerifyRequest(request); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}
	}

	// Create a new SlackRequestBody
	var slackRequestBody slackingway.SlackRequestBody
	switch request.Headers["Content-Type"] {
	// Slash Command & Interactive Components
	case "application/x-www-form-urlencoded":
		if err := slackRequestBody.ParseTimestamp(request.Headers["X-Slack-Request-Timestamp"]); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}

		// Parse the form data
		formData, err := url.ParseQuery(request.Body)
		if err != nil {
			log.Printf("Error parsing request body: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}

		// Check for payload
		if formData.Get("payload") != "" {
			// Parse the payload
			if err := slackRequestBody.ParsePayload(formData.Get("payload")); err != nil {
				log.Printf("Error parsing payload: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}
		} else {
			// Parse the slash command
			if err := slackRequestBody.ParseSlashCommand(formData); err != nil {
				log.Printf("Error parsing slash command: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}
		}
	// Any other Slack request
	case "application/json":
		if err := json.Unmarshal([]byte(request.Body), &slackRequestBody); err != nil {
			log.Printf("Error parsing request body: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}
	default:
		log.Printf("Unknown content type: %s", request.Headers["Content-Type"])
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}
	}

	// Initialize Slackingway
	s := slackingway.NewSlackingway(&slackRequestBody)

	// Handle the request
	var message slack.Msg
	switch s.SlackRequestBody.Type {
	// Challenge request
	case "url_verification":
		headers := make(map[string]string)
		headers["Content-Type"] = "text/plain"
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: http.StatusOK,
			Body:       s.SlackRequestBody.Challenge,
		}
	case "slash_command":
		switch s.SlackRequestBody.Command {
		case "/ping":
			err := s.WriteToHistory()
			if err != nil {
				log.Printf("Error writing to history: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}

			message, err = slackingway.Ping()
			if err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}
		case "/delayed-ping":
			err := s.WriteToHistory()
			if err != nil {
				log.Printf("Error writing to history: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}

			message, err = slackingway.DelayedPing(s)
			if err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}
		case "/echo":
			err := slackingway.Echo(s)
			if err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}
		default:
			log.Printf("Unknown command: %v", s.SlackRequestBody.Command)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}
		}
	case "view_submission":
		if DEBUG {
			// Log the view
			viewString, err := utils.GetStructFields(slackRequestBody.View)
			if err != nil {
				log.Printf("Error parsing view: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}
			log.Printf("Slack View: %v", viewString)
		}

		switch s.SlackRequestBody.View.CallbackID {
		case "/echo":
			err := s.WriteToHistory()
			if err != nil {
				log.Printf("Error writing to history: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}

			message, err = slackingway.ReceivedEcho(s)
			if err != nil {
				log.Printf("Error with %s: %v", s.SlackRequestBody.Command, err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
			}
		default:
			log.Printf("Unknown CallbackID: %v", s.SlackRequestBody.View.CallbackID)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}
		}
	default:
		log.Printf("Unknown request type: %s", s.SlackRequestBody.Type)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}
	}

	// Send the response to Slack if there is a message
	if message.Text == "" {
		// Create the response
		response, err := s.NewResponse(message)
		if err != nil {
			log.Printf("Error creating response: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}

		// Send the response
		if err := s.SendResponse(response); err != nil {
			log.Printf("Error sending response: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}
		}
	}

	// Return an empty response
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: http.StatusOK,
	}
}

func main() {
	lambda.Start(handler)
}
