# Tasks API - Go Serverless

A serverless API for managing tasks, built with Go and AWS Lambda.

## Project Structure

```
tasks-go/
├── .gitignore
├── README.md
├── go.mod
├── go.sum
├── package.json
├── serverless.yml
├── Makefile
├── .github/
│   └── workflows/
│       └── api.yml        # CI/CD workflow definition
└── cmd/
    └── api/
        ├── main.go         # Lambda handler and API routes
        ├── models.go       # Task struct and related types
        ├── store.go        # DynamoDB operations
        ├── handlers.go     # API handlers
        ├── models_test.go  # Tests for models
        ├── store_test.go   # Tests for store
        ├── store_mock.go   # Mock store for testing
        └── handlers_test.go # Tests for handlers
└── resources/
    └── dynamodb.yml       # DynamoDB table definition
```

## Prerequisites

- Go 1.x or later
- AWS CLI configured with appropriate credentials
- Serverless Framework

## Setup

1. Initialize the Go module:

```bash
go mod init github.com/user/tasks-api
go mod tidy
```

2. Create a `.env.development` file with your environment variables:

```
ALLOWED_ORIGINS=*
```

## Build

To build the Lambda function:

```bash
make build
```

This will create a binary in the `bin/` directory.

## Deploy

To deploy to AWS:

```bash
make deploy
```

To deploy to a specific stage:

```bash
make deploy-stage stage=production
```

## CI/CD Pipeline

This project includes a GitHub Actions workflow for continuous integration and deployment:

### Workflow Overview

The workflow is triggered on pushes to the repository that affect files in the `tasks-go` directory or the workflow file itself. It consists of three main jobs:

1. **Test**: Runs all Go tests and uploads coverage reports to Codecov
2. **Code Quality**: Runs golangci-lint to ensure code quality
3. **Deploy Development**: Builds and deploys the application to the development environment

### Required Secrets

The following secrets need to be configured in your GitHub repository:

- `SERVERLESS_ACCESS_KEY`: Your Serverless Framework access key
- `AWS_ACCESS_KEY_ID`: AWS access key with deployment permissions
- `AWS_SECRET_ACCESS_KEY`: Corresponding AWS secret key
- `CODECOV_TOKEN`: Token for uploading coverage reports to Codecov (optional)

### Workflow File

The workflow is defined in `.github/workflows/api.yml`. You can customize the workflow by editing this file.

### Running Tests Locally

To run the same tests that the CI pipeline runs:

```bash
go test -v -coverprofile=coverage.out ./...
```

To run linting locally:

```bash
# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

# Run linting
golangci-lint run ./...
```

## API Endpoints

- `GET /api/health-check/`: Health check endpoint
- `GET /api/tasks/?owner={owner}&status={status}`: List tasks for an owner (status is optional, defaults to OPEN)
- `POST /api/tasks/`: Create a new task
- `GET /api/tasks/{taskId}?owner={owner}`: Get a task by ID

## Example Requests

### Create a Task

```bash
curl -X POST https://your-api-url/api/tasks/ \
  -H "Content-Type: application/json" \
  -d '{"title": "Clean your office", "owner": "john@doe.com"}'
```

### List Open Tasks

```bash
curl https://your-api-url/api/tasks/?owner=john@doe.com
```

### List Closed Tasks

```bash
curl https://your-api-url/api/tasks/?owner=john@doe.com&status=CLOSED
```

### Get a Task by ID

```bash
curl https://your-api-url/api/tasks/123e4567-e89b-12d3-a456-426614174000?owner=john@doe.com
```
