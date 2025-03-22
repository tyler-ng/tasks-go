package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// CreateTaskRequest represents a request to create a task
type CreateTaskRequest struct {
	Title string `json:"title"`
	Owner string `json:"owner"`
}

// API handles API requests
type API struct {
	store *TaskStore
}

// NewAPI creates a new API
func NewAPI(tableName string) (*API, error) {
	store, err := NewTaskStore(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to create task store: %w", err)
	}

	return &API{
		store: store,
	}, nil
}

// HandleRequest handles API Gateway proxy requests
func (api *API) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract the path and method
	path := request.Path
	method := request.HTTPMethod

	// Handle health check
	if path == "/api/health-check/" && method == http.MethodGet {
		return api.healthCheck(ctx)
	}

	// Handle tasks
	if strings.HasPrefix(path, "/api/tasks/") {
		// Extract the task ID if present
		parts := strings.Split(path, "/")
		if len(parts) > 3 {
			taskID := parts[3]
			if taskID != "" {
				return api.handleTaskByID(ctx, method, taskID, request)
			}
		}

		// Handle tasks collection
		return api.handleTasks(ctx, method, request)
	}

	// Handle unknown paths
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"message": "Not Found"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// healthCheck handles health check requests
func (api *API) healthCheck(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message": "OK"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// handleTasks handles requests to the tasks collection
func (api *API) handleTasks(ctx context.Context, method string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch method {
	case http.MethodGet:
		return api.listTasks(ctx, request)
	case http.MethodPost:
		return api.createTask(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       `{"message": "Method Not Allowed"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}

// handleTaskByID handles requests to a specific task
func (api *API) handleTaskByID(ctx context.Context, method, taskID string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch method {
	case http.MethodGet:
		return api.getTask(ctx, taskID, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       `{"message": "Method Not Allowed"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}

// listTasks lists tasks
func (api *API) listTasks(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get the owner from the query parameters
	owner := request.QueryStringParameters["owner"]
	if owner == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Owner is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Get the status from the query parameters
	status := request.QueryStringParameters["status"]
	var tasks []Task
	var err error

	// List tasks by status
	if status == string(TaskStatusClosed) {
		tasks, err = api.store.ListClosed(ctx, owner)
	} else {
		// Default to open tasks
		tasks, err = api.store.ListOpen(ctx, owner)
	}

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message": "Failed to list tasks: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Marshal the tasks to JSON
	body, err := json.Marshal(tasks)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message": "Failed to marshal tasks: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// getTask gets a task by ID
func (api *API) getTask(ctx context.Context, taskIDStr string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the task ID
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"message": "Invalid task ID: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Get the owner from the query parameters
	owner := request.QueryStringParameters["owner"]
	if owner == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Owner is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Get the task
	task, err := api.store.GetByID(ctx, taskID, owner)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       fmt.Sprintf(`{"message": "Task not found: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Marshal the task to JSON
	body, err := json.Marshal(task)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message": "Failed to marshal task: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

// createTask creates a new task
func (api *API) createTask(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the request body
	var createRequest CreateTaskRequest
	if err := json.Unmarshal([]byte(request.Body), &createRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"message": "Invalid request body: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Validate the request
	if createRequest.Title == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Title is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	if createRequest.Owner == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Owner is required"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Create the task
	task := NewTask(uuid.New(), createRequest.Title, createRequest.Owner)

	// Add the task to the store
	if err := api.store.Add(ctx, task); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message": "Failed to create task: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Marshal the task to JSON
	body, err := json.Marshal(task)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message": "Failed to marshal task: %s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
