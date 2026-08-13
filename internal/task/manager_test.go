package task

import (
	"errors"
	"testing"
)

func TestTaskManager_AddAndListTasks(t *testing.T) {
	memStore := NewMemoryStore()
	tm := NewTaskManager(memStore)

	t.Run("List empty task list", func(t *testing.T) {
		tasks, err := tm.ListTasks()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 0 {
			t.Errorf("expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("Add valid tasks", func(t *testing.T) {
		t1, err := tm.AddTask("Write unit tests")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if t1.ID != 1 {
			t.Errorf("expected ID 1, got %d", t1.ID)
		}
		if t1.Description != "Write unit tests" {
			t.Errorf("expected description 'Write unit tests', got '%s'", t1.Description)
		}

		t2, err := tm.AddTask("Implement CLI handlers")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if t2.ID != 2 {
			t.Errorf("expected ID 2, got %d", t2.ID)
		}
	})

	t.Run("Add task with empty description fails", func(t *testing.T) {
		_, err := tm.AddTask("")
		if !errors.Is(err, ErrEmptyDescription) {
			t.Errorf("expected ErrEmptyDescription, got %v", err)
		}

		_, err = tm.AddTask("   ")
		if !errors.Is(err, ErrEmptyDescription) {
			t.Errorf("expected ErrEmptyDescription for whitespace, got %v", err)
		}
	})

	t.Run("List tasks returns created tasks", func(t *testing.T) {
		tasks, err := tm.ListTasks()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(tasks))
		}
		if tasks[0].Description != "Write unit tests" || tasks[1].Description != "Implement CLI handlers" {
			t.Errorf("unexpected task list output: %+v", tasks)
		}
	})
}

func TestTaskManager_CompleteTask(t *testing.T) {
	memStore := NewMemoryStore()
	tm := NewTaskManager(memStore)

	t1, err := tm.AddTask("Fix bug in parser")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	t.Run("Complete existing incomplete task", func(t *testing.T) {
		completedTask, err := tm.CompleteTask(t1.ID)
		if err != nil {
			t.Fatalf("expected successful completion, got %v", err)
		}
		if !completedTask.Completed {
			t.Errorf("expected task to be completed")
		}
		if completedTask.CompletedAt == nil {
			t.Errorf("expected CompletedAt timestamp to be set")
		}
		if completedTask.StatusTag() != "[x]" {
			t.Errorf("expected status tag '[x]', got '%s'", completedTask.StatusTag())
		}
	})

	t.Run("Complete already completed task returns ErrTaskAlreadyCompleted", func(t *testing.T) {
		_, err := tm.CompleteTask(t1.ID)
		if !errors.Is(err, ErrTaskAlreadyCompleted) {
			t.Errorf("expected ErrTaskAlreadyCompleted, got %v", err)
		}
	})

	t.Run("Complete non-existent task returns ErrTaskNotFound", func(t *testing.T) {
		_, err := tm.CompleteTask(999)
		if !errors.Is(err, ErrTaskNotFound) {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("Complete task with invalid ID returns ErrInvalidTaskID", func(t *testing.T) {
		_, err := tm.CompleteTask(0)
		if !errors.Is(err, ErrInvalidTaskID) {
			t.Errorf("expected ErrInvalidTaskID for 0, got %v", err)
		}

		_, err = tm.CompleteTask(-1)
		if !errors.Is(err, ErrInvalidTaskID) {
			t.Errorf("expected ErrInvalidTaskID for negative ID, got %v", err)
		}
	})
}
