package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

// TaskStore handles operations on tasks in DynamoDB
type TaskStore struct {
	client    *dynamodb.Client
	tableName string
}

// NewTaskStore creates a new TaskStore
func NewTaskStore(tableName string) (*TaskStore, error) {
	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create a DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	return &TaskStore{
		client:    client,
		tableName: tableName,
	}, nil
}

// Add adds a task to DynamoDB
func (ts *TaskStore) Add(ctx context.Context, task Task) error {
	// Convert the task to a DynamoDB item
	item := ToDynamoDBTask(task)

	// Marshal the item to a map
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// Put the item in DynamoDB
	_, err = ts.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(ts.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put task in DynamoDB: %w", err)
	}

	return nil
}

// GetByID gets a task by ID and owner
func (ts *TaskStore) GetByID(ctx context.Context, taskID uuid.UUID, owner string) (Task, error) {
	// Create the key
	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "#" + owner},
		"SK": &types.AttributeValueMemberS{Value: "#" + taskID.String()},
	}

	// Get the item from DynamoDB
	result, err := ts.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(ts.tableName),
		Key:       key,
	})
	if err != nil {
		return Task{}, fmt.Errorf("failed to get task from DynamoDB: %w", err)
	}

	// Check if the item exists
	if result.Item == nil {
		return Task{}, fmt.Errorf("task not found")
	}

	// Unmarshal the item
	var dbTask DynamoDBTask
	if err := attributevalue.UnmarshalMap(result.Item, &dbTask); err != nil {
		return Task{}, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	// Convert to a Task
	task, err := dbTask.ToTask()
	if err != nil {
		return Task{}, fmt.Errorf("failed to convert to task: %w", err)
	}

	return task, nil
}

// ListOpen lists open tasks for an owner
func (ts *TaskStore) ListOpen(ctx context.Context, owner string) ([]Task, error) {
	return ts.listByStatus(ctx, owner, TaskStatusOpen)
}

// ListClosed lists closed tasks for an owner
func (ts *TaskStore) ListClosed(ctx context.Context, owner string) ([]Task, error) {
	return ts.listByStatus(ctx, owner, TaskStatusClosed)
}

// listByStatus lists tasks by status for an owner
func (ts *TaskStore) listByStatus(ctx context.Context, owner string, status TaskStatus) ([]Task, error) {
	// Create the query input
	input := &dynamodb.QueryInput{
		TableName:              aws.String(ts.tableName),
		IndexName:              aws.String("GS1"),
		KeyConditionExpression: aws.String("GS1PK = :gspk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gspk": &types.AttributeValueMemberS{Value: "#" + owner + "#" + string(status)},
		},
	}

	// Query DynamoDB
	var tasks []Task
	var lastKey map[string]types.AttributeValue

	for {
		// If there's a last evaluated key, use it
		if lastKey != nil {
			input.ExclusiveStartKey = lastKey
		}

		// Execute the query
		result, err := ts.client.Query(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to query tasks: %w", err)
		}

		// Unmarshal the items
		var dbTasks []DynamoDBTask
		if err := attributevalue.UnmarshalListOfMaps(result.Items, &dbTasks); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
		}

		// Convert to Tasks
		for _, dbTask := range dbTasks {
			task, err := dbTask.ToTask()
			if err != nil {
				return nil, fmt.Errorf("failed to convert to task: %w", err)
			}
			tasks = append(tasks, task)
		}

		// Check if there are more items
		lastKey = result.LastEvaluatedKey
		if lastKey == nil {
			break
		}
	}

	return tasks, nil
}
