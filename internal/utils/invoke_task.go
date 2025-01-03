package utils

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/kn-lim/slackingway-bot/internal/slackingway"
)

// InvokeTaskFunction invokes the task function with the given Slack request
func InvokeTaskFunction(ctx context.Context, request slackingway.SlackRequestBody, region string, taskFunction string) error {
	// Create a new AWS Lambda client
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
		return err
	}
	client := lambda.NewFromConfig(cfg)

	// Create the payload for the task function
	payload, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshalling payload: %v", err)
		return err
	}

	// Create the input for the task function
	input := &lambda.InvokeInput{
		FunctionName:   aws.String(taskFunction),
		Payload:        payload,
		InvocationType: types.InvocationTypeEvent,
	}

	// Invoke the task function
	if _, err = client.Invoke(ctx, input); err != nil {
		log.Printf("Error invoking task function: %v", err)
		return err
	}

	return nil
}
