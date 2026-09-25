package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Event struct {
	Timestamp   time.Time `json:"timestamp"`
	Action      string    `json:"action"`
	Source      string    `json:"source"`
	Destination string    `json:"destination,omitempty"`
	Digest      string    `json:"digest,omitempty"`
	Previous    string    `json:"previous,omitempty"`
	Notes       []string  `json:"notes,omitempty"`
}

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	return filepath.Join(dir, "rebaze", "history.json"), nil
}

func Load(path string) ([]Event, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ev []Event
	if err := json.Unmarshal(data, &ev); err != nil {
		return nil, fmt.Errorf("parse history: %w", err)
	}
	return ev, nil
}

func Append(path string, ev Event) error {
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now().UTC()
	}
	existing, err := Load(path)
	if err != nil {
		return err
	}
	existing = append(existing, ev)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func LastFor(path, dest string) (*Event, error) {
	ev, err := Load(path)
	if err != nil {
		return nil, err
	}
	for i := len(ev) - 1; i >= 0; i-- {
		if ev[i].Destination == dest || ev[i].Source == dest {
			return &ev[i], nil
		}
	}
	return nil, fmt.Errorf("no history entry for %s", dest)
}
