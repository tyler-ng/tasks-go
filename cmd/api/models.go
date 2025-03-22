package main

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	// TaskStatusOpen represents an open task
	TaskStatusOpen TaskStatus = "OPEN"
	// TaskStatusClosed represents a closed task
	TaskStatusClosed TaskStatus = "CLOSED"
)

// Task represents a task in the system
type Task struct {
	ID     uuid.UUID  `json:"id"`
	Title  string     `json:"title"`
	Status TaskStatus `json:"status"`
	Owner  string     `json:"owner"`
}

// NewTask creates a new task with the given ID, title, and owner
func NewTask(id uuid.UUID, title, owner string) Task {
	return Task{
		ID:     id,
		Title:  title,
		Status: TaskStatusOpen,
		Owner:  owner,
	}
}

// DynamoDBTask represents a task in DynamoDB
type DynamoDBTask struct {
	PK    string     `json:"PK"`
	SK    string     `json:"SK"`
	GS1PK string     `json:"GS1PK"`
	GS1SK string     `json:"GS1SK"`
	ID    string     `json:"id"`
	Title string     `json:"title"`
	Owner string     `json:"owner"`
	Status TaskStatus `json:"status"`
}

// ToTask converts a DynamoDBTask to a Task
func (dt *DynamoDBTask) ToTask() (Task, error) {
	id, err := uuid.Parse(dt.ID)
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:     id,
		Title:  dt.Title,
		Status: dt.Status,
		Owner:  dt.Owner,
	}, nil
}

// ToDynamoDBTask converts a Task to a DynamoDBTask
func ToDynamoDBTask(task Task) DynamoDBTask {
	now := time.Now().UTC().Format(time.RFC3339)
	return DynamoDBTask{
		PK:    "#" + task.Owner,
		SK:    "#" + task.ID.String(),
		GS1PK: "#" + task.Owner + "#" + string(task.Status),
		GS1SK: "#" + now,
		ID:    task.ID.String(),
		Title: task.Title,
		Owner: task.Owner,
		Status: task.Status,
	}
}
