package slackingway

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaSvc "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/slack-go/slack"
)

const (
	DELAY = 5 * time.Second
)

// StartDelayedPing responds to a delayed ping request and sends a payload to the task function
func StartDelayedPing(ctx context.Context, slackRequestBody SlackRequestBody) (events.APIGatewayProxyResponse, error) {
	log.Printf("/delayed-ping received")

	// Create a new AWS Lambda client
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
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

	// Send a message acknowledging the request
	response := slack.Msg{Text: "..."}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: http.StatusOK,
		Body:       response.Text,
	}, nil
}

// ReturnDelayedPing sends a delayed ping response to the user
func ReturnDelayedPing(responseURL string) error {
	log.Printf("Returning delayed ping")

	// Wait for the delay
	time.Sleep(DELAY)

	// Create Slack message
	message := slack.Msg{
		Text:            "Delayed Pong!",
		ReplaceOriginal: true,
	}

	// Convert the response to JSON
	responseBody, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		return err
	}

	// Send the response to the response URL
	req, err := http.NewRequest("POST", responseURL, bytes.NewBuffer(responseBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client and send request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return err
	}
	defer response.Body.Close()

	// Check for non-OK status
	if response.StatusCode != http.StatusOK {
		log.Printf("Non-OK HTTP status: %v", response.StatusCode)
		return err
	}

	return nil
}
