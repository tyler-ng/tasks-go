.PHONY: build clean deploy

# Binary output directory
BIN_DIR := bin
# Lambda function name
LAMBDA_FUNCTION := api

# Go build flags
GOFLAGS := -ldflags="-s -w"

# Default target
all: build

# Build the Lambda function
build:
	@echo "Building Lambda function..."
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BIN_DIR)/$(LAMBDA_FUNCTION) ./cmd/api

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	rm -rf .serverless

# Deploy to AWS using Serverless Framework
deploy: build
	@echo "Deploying to AWS..."
	serverless deploy --verbose

# Deploy to a specific stage
deploy-stage: build
	@echo "Deploying to stage: $(stage)"
	serverless deploy --stage $(stage) --verbose

# Remove the service from AWS
remove:
	@echo "Removing service from AWS..."
	serverless remove

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
