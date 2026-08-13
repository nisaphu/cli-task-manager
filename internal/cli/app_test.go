package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cli-task-manager/internal/store"
	"cli-task-manager/internal/task"
)

// MockTimerRunner provides a fast timer runner for testing the CLI focus handler.
type MockTimerRunner struct {
	ShouldCancel bool
}

func (m *MockTimerRunner) Start() error {
	if m.ShouldCancel {
		return nil
	}
	return nil
}

func setupTestApp(t *testing.T) (*App, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "tasks.json")
	jsonStore := store.NewJSONStore(filePath)
	tm := task.NewTaskManager(jsonStore)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	mockTimer := &MockTimerRunner{}

	app := NewApp(tm, mockTimer, &outBuf, &errBuf)
	return app, &outBuf, &errBuf
}

func TestCLI_AddAndList(t *testing.T) {
	app, out, errOut := setupTestApp(t)

	t.Run("Empty task list message", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"list"})
		if code != 0 {
			t.Fatalf("expected exit code 0, got %d", code)
		}
		if !strings.Contains(out.String(), "No tasks found") {
			t.Errorf("expected empty list message, got: %s", out.String())
		}
	})

	t.Run("Add task success", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"add", "Write unit tests"})
		if code != 0 {
			t.Fatalf("expected exit code 0, got %d. Stderr: %s", code, errOut.String())
		}
		if !strings.Contains(out.String(), "Task 1 created: \"Write unit tests\"") {
			t.Errorf("unexpected output: %s", out.String())
		}
	})

	t.Run("Add task empty description error", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"add", ""})
		if code != 1 {
			t.Fatalf("expected exit code 1 for empty description, got %d", code)
		}
		if !strings.Contains(errOut.String(), "Task description cannot be empty") {
			t.Errorf("expected empty description error message, got: %s", errOut.String())
		}
	})

	t.Run("List tasks shows created task with incomplete status", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"list"})
		if code != 0 {
			t.Fatalf("expected exit code 0, got %d", code)
		}
		output := out.String()
		if !strings.Contains(output, "ID   STATUS       DESCRIPTION") {
			t.Errorf("expected header line, got: %s", output)
		}
		if !strings.Contains(output, "1    [ ]          Write unit tests") {
			t.Errorf("expected task 1 formatted as incomplete, got: %s", output)
		}
	})
}

func TestCLI_Complete(t *testing.T) {
	app, out, errOut := setupTestApp(t)

	// Add task ID 1
	app.Run([]string{"add", "Refactor storage layer"})

	t.Run("Complete existing task", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"complete", "1"})
		if code != 0 {
			t.Fatalf("expected exit code 0, got %d. Stderr: %s", code, errOut.String())
		}
		if !strings.Contains(out.String(), "Task 1 marked as completed") {
			t.Errorf("unexpected output: %s", out.String())
		}
	})

	t.Run("Complete already completed task", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"complete", "1"})
		if code != 0 {
			t.Fatalf("expected exit code 0 for idempotent completion, got %d", code)
		}
		if !strings.Contains(out.String(), "Task 1 is already completed") {
			t.Errorf("unexpected output: %s", out.String())
		}
	})

	t.Run("Complete non-existent task", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"complete", "99"})
		if code != 1 {
			t.Fatalf("expected exit code 1 for non-existent task, got %d", code)
		}
		if !strings.Contains(errOut.String(), "Task with ID 99 was not found") {
			t.Errorf("expected not found error message, got: %s", errOut.String())
		}
	})

	t.Run("Complete with invalid task ID argument", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"complete", "abc"})
		if code != 1 {
			t.Fatalf("expected exit code 1 for invalid ID, got %d", code)
		}
		if !strings.Contains(errOut.String(), "Invalid task ID \"abc\"") {
			t.Errorf("expected invalid ID error message, got: %s", errOut.String())
		}
	})
}

func TestCLI_UnknownCommandAndHelp(t *testing.T) {
	app, out, errOut := setupTestApp(t)

	t.Run("Unknown command", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"unknowncmd"})
		if code != 1 {
			t.Fatalf("expected exit code 1 for unknown command, got %d", code)
		}
		if !strings.Contains(errOut.String(), "Unknown command \"unknowncmd\"") {
			t.Errorf("expected unknown command error message, got: %s", errOut.String())
		}
	})

	t.Run("Help command", func(t *testing.T) {
		out.Reset()
		errOut.Reset()

		code := app.Run([]string{"help"})
		if code != 0 {
			t.Fatalf("expected exit code 0 for help, got %d", code)
		}
		if !strings.Contains(out.String(), "Available commands:") {
			t.Errorf("expected help output, got: %s", out.String())
		}
	})
}

func TestCLI_EndToEndIntegration(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "integration_tasks.json")
	jsonStore := store.NewJSONStore(filePath)
	tm := task.NewTaskManager(jsonStore)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	mockTimer := &MockTimerRunner{}
	app := NewApp(tm, mockTimer, &outBuf, &errBuf)

	// 1. Add Task 1
	if code := app.Run([]string{"add", "Task One"}); code != 0 {
		t.Fatalf("failed to add Task One")
	}

	// 2. Add Task 2
	if code := app.Run([]string{"add", "Task Two"}); code != 0 {
		t.Fatalf("failed to add Task Two")
	}

	// 3. List tasks (1 and 2 incomplete)
	outBuf.Reset()
	if code := app.Run([]string{"list"}); code != 0 {
		t.Fatalf("failed to list tasks")
	}
	output := outBuf.String()
	if !strings.Contains(output, "1    [ ]          Task One") || !strings.Contains(output, "2    [ ]          Task Two") {
		t.Errorf("unexpected task list: %s", output)
	}

	// 4. Complete Task 1
	if code := app.Run([]string{"complete", "1"}); code != 0 {
		t.Fatalf("failed to complete Task 1")
	}

	// 5. Re-create app using same JSON store to simulate process restart
	newTM := task.NewTaskManager(jsonStore)
	newApp := NewApp(newTM, mockTimer, &outBuf, &errBuf)

	outBuf.Reset()
	if code := newApp.Run([]string{"list"}); code != 0 {
		t.Fatalf("failed to list tasks after process restart")
	}
	outputAfterRestart := outBuf.String()

	if !strings.Contains(outputAfterRestart, "1    [x]          Task One") {
		t.Errorf("expected task 1 completed status [x] to persist across process restart, got: %s", outputAfterRestart)
	}
	if !strings.Contains(outputAfterRestart, "2    [ ]          Task Two") {
		t.Errorf("expected task 2 incomplete status [ ] after restart, got: %s", outputAfterRestart)
	}

	// Silence time unused warning
	_ = time.Second
}
