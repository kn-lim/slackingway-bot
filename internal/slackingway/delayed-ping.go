package slackingway

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

// DelayedPing responds to a ping request
func DelayedPing(responseURL string) (events.APIGatewayProxyResponse, error) {
	log.Printf("/delayed-ping received")

	response := slack.Msg{Text: "Delayed Pong!"}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: http.StatusOK,
		Body:       response.Text,
	}, nil
}
