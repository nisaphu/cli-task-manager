package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"cli-task-manager/internal/task"
)

// TaskData represents the JSON root document layout.
type TaskData struct {
	Version int          `json:"version"`
	NextID  int          `json:"next_id"`
	Tasks   []*task.Task `json:"tasks"`
}

// Store defines the interface for persisting and loading task collections.
type Store interface {
	Load() ([]*task.Task, int, error)
	Save(tasks []*task.Task, nextID int) error
}

// JSONStore implements Store interface using a local JSON file per ADR 0001.
type JSONStore struct {
	filePath string
}

// NewJSONStore creates a new JSONStore instance for the target file path.
func NewJSONStore(filePath string) *JSONStore {
	return &JSONStore{filePath: filePath}
}

// FilePath returns the configured JSON file path.
func (s *JSONStore) FilePath() string {
	return s.filePath
}

// Load reads tasks from the local JSON file. If the file does not exist,
// it returns an empty task collection with nextID = 1.
func (s *JSONStore) Load() ([]*task.Task, int, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*task.Task{}, 1, nil
		}
		return nil, 0, fmt.Errorf("%w: unable to read file: %v", task.ErrCorruptedStorage, err)
	}

	if len(data) == 0 {
		return []*task.Task{}, 1, nil
	}

	var td TaskData
	if err := json.Unmarshal(data, &td); err != nil {
		return nil, 0, fmt.Errorf("%w: invalid JSON format: %v", task.ErrCorruptedStorage, err)
	}

	if td.Tasks == nil {
		td.Tasks = []*task.Task{}
	}
	if td.NextID < 1 {
		td.NextID = 1
	}

	return td.Tasks, td.NextID, nil
}

// Save writes tasks to the local JSON file using atomic file replace.
func (s *JSONStore) Save(tasks []*task.Task, nextID int) error {
	if tasks == nil {
		tasks = []*task.Task{}
	}

	td := TaskData{
		Version: 1,
		NextID:  nextID,
		Tasks:   tasks,
	}

	data, err := json.MarshalIndent(td, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize task data: %w", err)
	}

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Atomic write: create temp file in same directory, write data, then rename.
	tmpFile, err := os.CreateTemp(dir, "tasks-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpName)
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tmpName, s.filePath); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("failed to replace storage file: %w", err)
	}

	return nil
}
