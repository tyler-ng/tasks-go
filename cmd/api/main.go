package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// getTableName gets the DynamoDB table name from the environment
func getTableName() string {
	// Get the stage from the environment
	stage := os.Getenv("APP_ENVIRONMENT")
	if stage == "" {
		stage = "development"
	}

	// Return the table name
	return fmt.Sprintf("%s-tasks-api", stage)
}

// handleRequest is the Lambda handler
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get the table name
	tableName := getTableName()

	// Create the API
	api, err := NewAPI(tableName)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"message": "Failed to create API: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Handle the request
	return api.HandleRequest(ctx, request)
}

func main() {
	// Start the Lambda handler
	lambda.Start(handleRequest)
}
