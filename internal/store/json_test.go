package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cli-task-manager/internal/task"
)

func TestJSONStore_SaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "tasks.json")
	js := NewJSONStore(filePath)

	t.Run("Missing file loads as empty collection", func(t *testing.T) {
		tasks, nextID, err := js.Load()
		if err != nil {
			t.Fatalf("expected no error loading non-existent file, got %v", err)
		}
		if len(tasks) != 0 {
			t.Errorf("expected empty task slice, got %d items", len(tasks))
		}
		if nextID != 1 {
			t.Errorf("expected default nextID 1, got %d", nextID)
		}
	})

	t.Run("Save and reload tasks", func(t *testing.T) {
		t1, _ := task.NewTask(1, "First task")
		t2, _ := task.NewTask(2, "Second task")
		now := time.Now().UTC()
		t2.Completed = true
		t2.CompletedAt = &now

		saveList := []*task.Task{t1, t2}
		err := js.Save(saveList, 3)
		if err != nil {
			t.Fatalf("failed to save tasks: %v", err)
		}

		loadedTasks, nextID, err := js.Load()
		if err != nil {
			t.Fatalf("failed to reload tasks: %v", err)
		}

		if nextID != 3 {
			t.Errorf("expected nextID 3, got %d", nextID)
		}
		if len(loadedTasks) != 2 {
			t.Fatalf("expected 2 tasks loaded, got %d", len(loadedTasks))
		}
		if loadedTasks[0].Description != "First task" || loadedTasks[0].Completed {
			t.Errorf("unexpected task 1 data: %+v", loadedTasks[0])
		}
		if loadedTasks[1].Description != "Second task" || !loadedTasks[1].Completed {
			t.Errorf("unexpected task 2 data: %+v", loadedTasks[1])
		}
	})

	t.Run("Corrupted JSON returns ErrCorruptedStorage", func(t *testing.T) {
		corruptFile := filepath.Join(tempDir, "corrupt.json")
		if err := os.WriteFile(corruptFile, []byte("{ invalid json content ..."), 0644); err != nil {
			t.Fatalf("failed to write corrupt file: %v", err)
		}

		corruptStore := NewJSONStore(corruptFile)
		_, _, err := corruptStore.Load()
		if !errors.Is(err, task.ErrCorruptedStorage) {
			t.Errorf("expected ErrCorruptedStorage, got %v", err)
		}
	})
}
