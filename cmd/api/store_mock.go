package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// MockTaskStore is a mock implementation of the TaskStore for testing
type MockTaskStore struct {
	tasks map[string]map[string]Task // map[owner]map[taskID]Task
}

// NewMockTaskStore creates a new MockTaskStore
func NewMockTaskStore() *MockTaskStore {
	return &MockTaskStore{
		tasks: make(map[string]map[string]Task),
	}
}

// Add adds a task to the mock store
func (m *MockTaskStore) Add(ctx context.Context, task Task) error {
	// Initialize the owner's map if it doesn't exist
	if _, ok := m.tasks[task.Owner]; !ok {
		m.tasks[task.Owner] = make(map[string]Task)
	}

	// Add the task
	m.tasks[task.Owner][task.ID.String()] = task

	return nil
}

// GetByID gets a task by ID and owner
func (m *MockTaskStore) GetByID(ctx context.Context, taskID uuid.UUID, owner string) (Task, error) {
	// Check if the owner exists
	ownerTasks, ok := m.tasks[owner]
	if !ok {
		return Task{}, fmt.Errorf("task not found")
	}

	// Check if the task exists
	task, ok := ownerTasks[taskID.String()]
	if !ok {
		return Task{}, fmt.Errorf("task not found")
	}

	return task, nil
}

// ListOpen lists open tasks for an owner
func (m *MockTaskStore) ListOpen(ctx context.Context, owner string) ([]Task, error) {
	return m.listByStatus(ctx, owner, TaskStatusOpen)
}

// ListClosed lists closed tasks for an owner
func (m *MockTaskStore) ListClosed(ctx context.Context, owner string) ([]Task, error) {
	return m.listByStatus(ctx, owner, TaskStatusClosed)
}

// listByStatus lists tasks by status for an owner
func (m *MockTaskStore) listByStatus(ctx context.Context, owner string, status TaskStatus) ([]Task, error) {
	// Check if the owner exists
	ownerTasks, ok := m.tasks[owner]
	if !ok {
		return []Task{}, nil
	}

	// Filter tasks by status
	var tasks []Task
	for _, task := range ownerTasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}
