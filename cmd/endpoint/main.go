package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaSvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/slack-go/slack"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
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

	// Validate the request
	verifier, err := slack.NewSecretsVerifier(headers, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		log.Printf("Error creating verifier: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	if _, err := verifier.Write([]byte(request.Body)); err != nil {
		log.Printf("Error writing body to verifier: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
	if err = verifier.Ensure(); err != nil {
		log.Printf("Invalid request: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, nil
	}

	// Parse the request body
	var slackRequestBody slackingway.SlackRequestBody

	// Check content type
	switch request.Headers["Content-Type"] {
	// Slash Command & Interactive Components
	case "application/x-www-form-urlencoded":
		// log.Printf("Found application/x-www-form-urlencoded request")

		slackRequestBody.Type = "interaction"

		// Get Time Stamp
		timestamp, err := strconv.ParseInt(request.Headers["X-Slack-Request-Timestamp"], 10, 64)
		if err != nil {
			log.Printf("Error parsing timestamp: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
		}
		slackRequestBody.Timestamp = time.Unix(timestamp, 0).UTC().Format(time.RFC822)

		// Parse the form data
		formData, err := url.ParseQuery(request.Body)
		if err != nil {
			log.Printf("Error parsing request body: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
		}

		// Convert form data to JSON string
		formDataJSON, err := json.Marshal(formData)
		if err != nil {
			log.Printf("Error marshalling form data to JSON: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
		}

		// Unmarshal JSON string into slackRequestBody
		if err := json.Unmarshal(formDataJSON, &slackRequestBody); err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
		}

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
	// User interaction request
	case "interaction":
		s := slackingway.NewSlackingway(&slackRequestBody)
		switch slackRequestBody.Command {
		// Add all slash commands that involves trigger_id
		case "/echo":
			err := s.WriteToHistory()
			if err != nil {
				log.Printf("Error writing to history: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
			}

			err = slackingway.Echo(s)
			if err != nil {
				log.Printf("Error echoing: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
			}

			return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
		// For all other slash commands
		default:
			// Create a new AWS Lambda client
			cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
			if err != nil {
				log.Printf("Error loading AWS config: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
			}
			client := lambdaSvc.NewFromConfig(cfg)

			// Create the payload for the task function
			payload, err := json.Marshal(slackRequestBody)
			if err != nil {
				log.Printf("Error marshalling payload: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
			}

			// Create the input for the task function
			input := &lambdaSvc.InvokeInput{
				FunctionName:   aws.String(os.Getenv("TASK_FUNCTION_NAME")),
				Payload:        payload,
				InvocationType: types.InvocationTypeEvent,
			}

			// Invoke the task function
			if _, err = client.Invoke(ctx, input); err != nil {
				log.Printf("Error invoking task function: %v", err)
				return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
			}

			// Return an empty response
			headers := make(map[string]string)
			headers["Content-Type"] = "application/json"
			return events.APIGatewayProxyResponse{
				Headers:    headers,
				StatusCode: http.StatusOK,
			}, nil
		}
	// Unknown request type
	default:
		log.Printf("Unknown request type: %s", slackRequestBody.Type)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}
}

func main() {
	lambda.Start(handler)
}
