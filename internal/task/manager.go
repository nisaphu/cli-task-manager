package task

import (
	"fmt"
	"time"
)

// StoreInterface defines storage capability for TaskManager.
type StoreInterface interface {
	Load() ([]*Task, int, error)
	Save(tasks []*Task, nextID int) error
}

// MemoryStore provides an in-memory Store implementation for testing without disk I/O.
type MemoryStore struct {
	tasks  []*Task
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:  []*Task{},
		nextID: 1,
	}
}

func (m *MemoryStore) Load() ([]*Task, int, error) {
	return m.tasks, m.nextID, nil
}

func (m *MemoryStore) Save(tasks []*Task, nextID int) error {
	m.tasks = tasks
	m.nextID = nextID
	return nil
}

// TaskManager coordinates task business operations and persistence.
type TaskManager struct {
	store StoreInterface
}

// NewTaskManager creates a new TaskManager wrapping a Store.
func NewTaskManager(store StoreInterface) *TaskManager {
	return &TaskManager{store: store}
}

// AddTask creates a new task with a non-empty description, assigns a unique ID, and saves it.
func (m *TaskManager) AddTask(description string) (*Task, error) {
	tasks, nextID, err := m.store.Load()
	if err != nil {
		return nil, err
	}

	newTask, err := NewTask(nextID, description)
	if err != nil {
		return nil, err
	}

	tasks = append(tasks, newTask)
	nextID++

	if err := m.store.Save(tasks, nextID); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return newTask, nil
}

// ListTasks returns all tasks currently managed by the system.
func (m *TaskManager) ListTasks() ([]*Task, error) {
	tasks, _, err := m.store.Load()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// CompleteTask marks a task as completed using its unique identifier.
func (m *TaskManager) CompleteTask(id int) (*Task, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskID
	}

	tasks, nextID, err := m.store.Load()
	if err != nil {
		return nil, err
	}

	var targetTask *Task
	for _, t := range tasks {
		if t.ID == id {
			targetTask = t
			break
		}
	}

	if targetTask == nil {
		return nil, fmt.Errorf("%w: ID %d", ErrTaskNotFound, id)
	}

	if targetTask.Completed {
		return targetTask, ErrTaskAlreadyCompleted
	}

	now := time.Now().UTC()
	targetTask.Completed = true
	targetTask.CompletedAt = &now

	if err := m.store.Save(tasks, nextID); err != nil {
		return nil, fmt.Errorf("failed to save task completion: %w", err)
	}

	return targetTask, nil
}
