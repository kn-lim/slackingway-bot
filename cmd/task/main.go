package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) {
	// TODO
}

func main() {
	lambda.Start(handler)
}
