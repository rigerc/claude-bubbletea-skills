// Package state manages persistent loop state for ralphio, using a
// crash-safe write (tmp+rename) strategy to prevent corruption.
package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	stateDir  = ".ralph"
	stateFile = "state.json"
	stateTmp  = "state.json.tmp"

	StatusRunning = "running"
	StatusPaused  = "paused"
	StatusStopped = "stopped"
	StatusError   = "error"
)

// State holds the persistent execution state of the ralphio loop.
type State struct {
	CurrentIteration int       `json:"currentIteration"`
	CurrentTaskID    string    `json:"currentTaskId"`
	LoopStatus       string    `json:"loopStatus"`
	ActiveAdapter    string    `json:"activeAdapter"`
	ActiveModel      string    `json:"activeModel"`
	LastUpdated      time.Time `json:"lastUpdated"`
}

func defaultState() *State {
	return &State{
		LoopStatus:    StatusStopped,
		ActiveAdapter: "claude",
	}
}

func statePath(dir string) string  { return filepath.Join(dir, stateDir, stateFile) }
func tmpPath(dir string) string    { return filepath.Join(dir, stateDir, stateTmp) }
func ralphDir(dir string) string   { return filepath.Join(dir, stateDir) }

// Load reads state from dir/.ralph/state.json. If the file does not exist,
// a default state is returned. If a .tmp file is found without the real file,
// the interrupted write is recovered by renaming .tmp → state.json.
func Load(dir string) (*State, error) {
	real := statePath(dir)
	tmp := tmpPath(dir)

	// Recover from an interrupted previous write.
	if _, err := os.Stat(real); errors.Is(err, os.ErrNotExist) {
		if _, tmpErr := os.Stat(tmp); tmpErr == nil {
			if renameErr := os.Rename(tmp, real); renameErr != nil {
				return nil, fmt.Errorf("recovering state from tmp: %w", renameErr)
			}
		}
	}

	data, err := os.ReadFile(real)
	if errors.Is(err, os.ErrNotExist) {
		return defaultState(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading state file: %w", err)
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing state file: %w", err)
	}
	return &s, nil
}

// Save writes s to dir/.ralph/state.json using a crash-safe tmp+rename
// strategy. The .ralph/ directory is created if it does not exist.
func Save(dir string, s *State) error {
	s.LastUpdated = time.Now()

	if err := os.MkdirAll(ralphDir(dir), 0o755); err != nil {
		return fmt.Errorf("creating state directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding state: %w", err)
	}

	tmp := tmpPath(dir)
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("writing tmp state file: %w", err)
	}

	if err := os.Rename(tmp, statePath(dir)); err != nil {
		return fmt.Errorf("committing state file: %w", err)
	}
	return nil
}
