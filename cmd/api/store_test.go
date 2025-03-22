package main

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestMockTaskStore_Add(t *testing.T) {
	// Arrange
	store := NewMockTaskStore()
	ctx := context.Background()
	task := NewTask(uuid.New(), "Test Task", "test@example.com")

	// Act
	err := store.Add(ctx, task)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMockTaskStore_GetByID(t *testing.T) {
	// Arrange
	store := NewMockTaskStore()
	ctx := context.Background()
	task := NewTask(uuid.New(), "Test Task", "test@example.com")
	_ = store.Add(ctx, task)

	// Act
	retrievedTask, err := store.GetByID(ctx, task.ID, task.Owner)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedTask.ID != task.ID {
		t.Errorf("Expected task ID to be %v, got %v", task.ID, retrievedTask.ID)
	}
	if retrievedTask.Title != task.Title {
		t.Errorf("Expected task title to be %s, got %s", task.Title, retrievedTask.Title)
	}
	if retrievedTask.Owner != task.Owner {
		t.Errorf("Expected task owner to be %s, got %s", task.Owner, retrievedTask.Owner)
	}
	if retrievedTask.Status != task.Status {
		t.Errorf("Expected task status to be %s, got %s", task.Status, retrievedTask.Status)
	}
}

func TestMockTaskStore_GetByID_NotFound(t *testing.T) {
	// Arrange
	store := NewMockTaskStore()
	ctx := context.Background()
	id := uuid.New()
	owner := "test@example.com"

	// Act
	_, err := store.GetByID(ctx, id, owner)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestMockTaskStore_ListOpen(t *testing.T) {
	// Arrange
	store := NewMockTaskStore()
	ctx := context.Background()
	owner := "test@example.com"
	openTask := NewTask(uuid.New(), "Open Task", owner)
	closedTask := Task{
		ID:     uuid.New(),
		Title:  "Closed Task",
		Status: TaskStatusClosed,
		Owner:  owner,
	}
	_ = store.Add(ctx, openTask)
	_ = store.Add(ctx, closedTask)

	// Act
	tasks, err := store.ListOpen(ctx, owner)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if len(tasks) > 0 && tasks[0].ID != openTask.ID {
		t.Errorf("Expected task ID to be %v, got %v", openTask.ID, tasks[0].ID)
	}
}

func TestMockTaskStore_ListClosed(t *testing.T) {
	// Arrange
	store := NewMockTaskStore()
	ctx := context.Background()
	owner := "test@example.com"
	openTask := NewTask(uuid.New(), "Open Task", owner)
	closedTask := Task{
		ID:     uuid.New(),
		Title:  "Closed Task",
		Status: TaskStatusClosed,
		Owner:  owner,
	}
	_ = store.Add(ctx, openTask)
	_ = store.Add(ctx, closedTask)

	// Act
	tasks, err := store.ListClosed(ctx, owner)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if len(tasks) > 0 && tasks[0].ID != closedTask.ID {
		t.Errorf("Expected task ID to be %v, got %v", closedTask.ID, tasks[0].ID)
	}
}
