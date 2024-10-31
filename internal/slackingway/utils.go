package slackingway

import (
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func ConveryProxyRequestToHTTPRequest(request events.APIGatewayProxyRequest) *http.Request {
	// Convert the request to an HTTP request
	httpRequest, err := http.NewRequest(request.HTTPMethod, request.Path, strings.NewReader(request.Body))
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return nil
	}

	// Copy the headers
	for key, value := range request.Headers {
		httpRequest.Header.Add(key, value)
	}

	return httpRequest
}
