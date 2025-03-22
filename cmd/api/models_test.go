package main

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTask(t *testing.T) {
	// Arrange
	id := uuid.New()
	title := "Test Task"
	owner := "test@example.com"

	// Act
	task := NewTask(id, title, owner)

	// Assert
	if task.ID != id {
		t.Errorf("Expected task ID to be %v, got %v", id, task.ID)
	}
	if task.Title != title {
		t.Errorf("Expected task title to be %s, got %s", title, task.Title)
	}
	if task.Owner != owner {
		t.Errorf("Expected task owner to be %s, got %s", owner, task.Owner)
	}
	if task.Status != TaskStatusOpen {
		t.Errorf("Expected task status to be %s, got %s", TaskStatusOpen, task.Status)
	}
}

func TestToDynamoDBTask(t *testing.T) {
	// Arrange
	id := uuid.New()
	task := Task{
		ID:     id,
		Title:  "Test Task",
		Status: TaskStatusOpen,
		Owner:  "test@example.com",
	}

	// Act
	dbTask := ToDynamoDBTask(task)

	// Assert
	if dbTask.PK != "#"+task.Owner {
		t.Errorf("Expected PK to be #%s, got %s", task.Owner, dbTask.PK)
	}
	if dbTask.SK != "#"+task.ID.String() {
		t.Errorf("Expected SK to be #%s, got %s", task.ID.String(), dbTask.SK)
	}
	if dbTask.GS1PK != "#"+task.Owner+"#"+string(task.Status) {
		t.Errorf("Expected GS1PK to be #%s#%s, got %s", task.Owner, task.Status, dbTask.GS1PK)
	}
	if dbTask.ID != task.ID.String() {
		t.Errorf("Expected ID to be %s, got %s", task.ID.String(), dbTask.ID)
	}
	if dbTask.Title != task.Title {
		t.Errorf("Expected Title to be %s, got %s", task.Title, dbTask.Title)
	}
	if dbTask.Owner != task.Owner {
		t.Errorf("Expected Owner to be %s, got %s", task.Owner, dbTask.Owner)
	}
	if dbTask.Status != task.Status {
		t.Errorf("Expected Status to be %s, got %s", task.Status, dbTask.Status)
	}
}

func TestDynamoDBTaskToTask(t *testing.T) {
	// Arrange
	id := uuid.New()
	dbTask := DynamoDBTask{
		PK:     "#test@example.com",
		SK:     "#" + id.String(),
		GS1PK:  "#test@example.com#OPEN",
		GS1SK:  "#2023-01-01T00:00:00Z",
		ID:     id.String(),
		Title:  "Test Task",
		Owner:  "test@example.com",
		Status: TaskStatusOpen,
	}

	// Act
	task, err := dbTask.ToTask()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if task.ID != id {
		t.Errorf("Expected task ID to be %v, got %v", id, task.ID)
	}
	if task.Title != dbTask.Title {
		t.Errorf("Expected task title to be %s, got %s", dbTask.Title, task.Title)
	}
	if task.Owner != dbTask.Owner {
		t.Errorf("Expected task owner to be %s, got %s", dbTask.Owner, task.Owner)
	}
	if task.Status != dbTask.Status {
		t.Errorf("Expected task status to be %s, got %s", dbTask.Status, task.Status)
	}
}
