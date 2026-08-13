package task

import (
	"errors"
	"strings"
	"time"
)

// Common domain errors.
var (
	ErrEmptyDescription     = errors.New("task description cannot be empty")
	ErrTaskNotFound         = errors.New("task not found")
	ErrTaskAlreadyCompleted = errors.New("task is already completed")
	ErrCorruptedStorage     = errors.New("task storage file is corrupted or unreadable")
	ErrInvalidTaskID        = errors.New("task ID must be a positive integer")
)

// TaskStatus represents the completion state of a task.
type TaskStatus string

const (
	StatusIncomplete TaskStatus = "incomplete"
	StatusCompleted  TaskStatus = "completed"
)

// Task represents a developer task item.
type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// NewTask creates a new Task entity after validating the description.
func NewTask(id int, description string) (*Task, error) {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return nil, ErrEmptyDescription
	}
	return &Task{
		ID:          id,
		Description: trimmed,
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		CompletedAt: nil,
	}, nil
}

// StatusTag returns a visual representation of the task status for CLI display.
func (t *Task) StatusTag() string {
	if t.Completed {
		return "[x]"
	}
	return "[ ]"
}
