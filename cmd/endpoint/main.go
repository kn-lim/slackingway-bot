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
	"github.com/kn-lim/slackingway-bot/internal/slackingway"
	"github.com/slack-go/slack"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log the request
	log.Printf("Request Headers: %v", request.Headers)
	log.Printf("Request Body: %v", request.Body)

	// Copy the headers
	headers := http.Header{}
	for key, value := range request.Headers {
		headers.Add(key, value)
	}

	// Create a new Slack secrets verifier
	verifier, err := slack.NewSecretsVerifier(headers, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		log.Printf("Error creating verifier: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	// Parse the request body
	var slackRequestBody slackingway.SlackRequestBody

	// Check content type
	switch request.Headers["Content-Type"] {
	// Slash Command
	case "application/x-www-form-urlencoded":
		// log.Printf("Found application/x-www-form-urlencoded request")

		formData, err := url.ParseQuery(request.Body)
		if err != nil {
			log.Printf("Error parsing request body: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
		}
		slackRequestBody.Type = "slash_command"
		slackRequestBody.Token = formData.Get("token")
		slackRequestBody.Command = formData.Get("command")
		slackRequestBody.Text = formData.Get("text")
		slackRequestBody.ResponseURL = formData.Get("response_url")
		slackRequestBody.UserID = formData.Get("user_id")
		slackRequestBody.ChannelID = formData.Get("channel_id")
		slackRequestBody.TeamID = formData.Get("team_id")

		log.Printf("Parsed form data: %v", slackRequestBody)
	// Any other Slack request
	case "application/json":
		// log.Printf("Found application/json request")

		if err := json.Unmarshal([]byte(request.Body), &slackRequestBody); err != nil {
			log.Printf("Error parsing request body: %v", err)

			headers := make(map[string]string)
			headers["Content-Type"] = "text/plain"
			return events.APIGatewayProxyResponse{
				Headers:    headers,
				StatusCode: http.StatusInternalServerError,
				Body:       "Error parsing request body",
			}, nil
		}
	default:
		log.Printf("Unknown content type: %s", request.Headers["Content-Type"])
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	// Validate the request
	if _, err := verifier.Write([]byte(request.Body)); err != nil {
		log.Printf("Error writing body to verifier: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	if err = verifier.Ensure(); err != nil {
		log.Printf("Invalid request: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, nil
	}

	switch slackRequestBody.Type {
	// Challenge request
	case "url_verification":
		headers := make(map[string]string)
		headers["Content-Type"] = "text/plain"
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: http.StatusOK,
			Body:       slackRequestBody.Challenge,
		}, nil
	// Slash command request
	case "slash_command":
		// Process the SlackRequestBody as needed
		switch slackRequestBody.Command {
		case "/ping":
			return slackingway.Ping()
		case "/delayed-ping":
			return slackingway.StartDelayedPing(ctx, slackRequestBody)
		default:
			headers := make(map[string]string)
			headers["Content-Type"] = "text/plain"
			return events.APIGatewayProxyResponse{
				Headers:    headers,
				StatusCode: http.StatusOK,
				Body:       "Unknown command",
			}, nil
		}
	default:
		log.Printf("Unknown request type: %s", slackRequestBody.Type)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
}

func main() {
	lambda.Start(handler)
}
