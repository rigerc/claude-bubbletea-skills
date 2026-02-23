// Package plan manages the task list for ralphio, persisted as tasks.json
// using the same crash-safe tmp+rename write strategy as the state package.
package plan

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	tasksFile = "tasks.json"
	tasksTmp  = "tasks.json.tmp"

	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusSkipped    = "skipped"
)

// Task describes a single unit of work to be executed by ralphio.
type Task struct {
	ID                string `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Priority          int    `json:"priority"`
	Status            string `json:"status"`
	RetryCount        int    `json:"retryCount"`
	MaxRetries        int    `json:"maxRetries"`
	ValidationCommand string `json:"validationCommand"`
}

// Manager reads and writes the task list in a given directory.
type Manager struct {
	dir string
}

// NewManager returns a Manager for the specified directory.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) tasksPath() string { return filepath.Join(m.dir, tasksFile) }
func (m *Manager) tmpPath() string   { return filepath.Join(m.dir, tasksTmp) }

// LoadTasks reads tasks from tasks.json in the manager's directory.
// Returns an empty slice if the file does not exist.
func (m *Manager) LoadTasks() ([]Task, error) {
	data, err := os.ReadFile(m.tasksPath())
	if errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading tasks file: %w", err)
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("parsing tasks file: %w", err)
	}
	return tasks, nil
}

// SaveTasks persists tasks to tasks.json using a crash-safe tmp+rename write.
func (m *Manager) SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding tasks: %w", err)
	}

	if err := os.WriteFile(m.tmpPath(), data, 0o644); err != nil {
		return fmt.Errorf("writing tmp tasks file: %w", err)
	}

	if err := os.Rename(m.tmpPath(), m.tasksPath()); err != nil {
		return fmt.Errorf("committing tasks file: %w", err)
	}
	return nil
}

// NextTask returns a pointer to the first pending task with the lowest
// Priority number. Returns nil if no pending task exists.
func NextTask(tasks []Task) *Task {
	pending := make([]Task, 0, len(tasks))
	for i := range tasks {
		if tasks[i].Status == StatusPending {
			pending = append(pending, tasks[i])
		}
	}
	if len(pending) == 0 {
		return nil
	}

	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Priority < pending[j].Priority
	})

	// Return pointer into the original slice so the caller can mutate it.
	for i := range tasks {
		if tasks[i].ID == pending[0].ID {
			return &tasks[i]
		}
	}
	return nil
}

// UpdateTask mutates the task with the given id in-place, setting its Status
// and RetryCount. Returns false if no task with that id is found.
func UpdateTask(tasks []Task, id, status string, retryCount int) bool {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].RetryCount = retryCount
			return true
		}
	}
	return false
}
