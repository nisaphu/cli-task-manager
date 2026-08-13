package task

import (
	"errors"
	"testing"
)

func TestNewTask(t *testing.T) {
	t.Run("Valid task creation", func(t *testing.T) {
		task, err := NewTask(1, "Write unit tests")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.ID != 1 {
			t.Errorf("expected ID 1, got %d", task.ID)
		}
		if task.Description != "Write unit tests" {
			t.Errorf("expected description 'Write unit tests', got '%s'", task.Description)
		}
		if task.Completed {
			t.Errorf("expected new task to be incomplete")
		}
		if task.CompletedAt != nil {
			t.Errorf("expected CompletedAt to be nil")
		}
		if task.StatusTag() != "[ ]" {
			t.Errorf("expected status tag '[ ]', got '%s'", task.StatusTag())
		}
	})

	t.Run("Empty description validation", func(t *testing.T) {
		_, err := NewTask(1, "")
		if !errors.Is(err, ErrEmptyDescription) {
			t.Errorf("expected ErrEmptyDescription, got %v", err)
		}

		_, err = NewTask(1, "   ")
		if !errors.Is(err, ErrEmptyDescription) {
			t.Errorf("expected ErrEmptyDescription for whitespace, got %v", err)
		}
	})

	t.Run("Whitespace trimming", func(t *testing.T) {
		task, err := NewTask(2, "  Clean codebase  ")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if task.Description != "Clean codebase" {
			t.Errorf("expected trimmed description 'Clean codebase', got '%s'", task.Description)
		}
	})
}
