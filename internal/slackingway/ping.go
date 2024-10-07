package slackingway

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

// Ping responds to a ping request
func Ping() (events.APIGatewayProxyResponse, error) {
	log.Printf("Ping request received")

	response := slack.Msg{Text: "Pong!"}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: http.StatusOK,
		Body:       response.Text,
	}, nil
}
