package main

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHealthCheck(t *testing.T) {
	// Arrange
	api := &API{}
	ctx := context.Background()

	// Act
	response, err := api.healthCheck(ctx)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Parse the response body
	var body map[string]string
	if err := json.Unmarshal([]byte(response.Body), &body); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// Check the message
	if message, ok := body["message"]; !ok || message != "OK" {
		t.Errorf("Expected message to be 'OK', got '%s'", message)
	}
}

func TestHandleRequestHealthCheck(t *testing.T) {
	// Arrange
	api := &API{}
	ctx := context.Background()
	request := events.APIGatewayProxyRequest{
		Path:       "/api/health-check/",
		HTTPMethod: http.MethodGet,
	}

	// Act
	response, err := api.HandleRequest(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Parse the response body
	var body map[string]string
	if err := json.Unmarshal([]byte(response.Body), &body); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// Check the message
	if message, ok := body["message"]; !ok || message != "OK" {
		t.Errorf("Expected message to be 'OK', got '%s'", message)
	}
}

func TestHandleRequestNotFound(t *testing.T) {
	// Arrange
	api := &API{}
	ctx := context.Background()
	request := events.APIGatewayProxyRequest{
		Path:       "/api/not-found/",
		HTTPMethod: http.MethodGet,
	}

	// Act
	response, err := api.HandleRequest(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.StatusCode)
	}

	// Parse the response body
	var body map[string]string
	if err := json.Unmarshal([]byte(response.Body), &body); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// Check the message
	if message, ok := body["message"]; !ok || message != "Not Found" {
		t.Errorf("Expected message to be 'Not Found', got '%s'", message)
	}
}
