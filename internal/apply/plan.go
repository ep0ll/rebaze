package apply

import (
	"encoding/json"
	"fmt"
	"os"
)

// Plan is the executable mutation document used by apply / preview.
// It is a practical subset of the MutationBundle in schema.json.
type Plan struct {
	SchemaVersion string         `json:"schemaVersion,omitempty"`
	Kind          string         `json:"kind,omitempty"`
	Image         string         `json:"image"`
	Tag           string         `json:"tag,omitempty"`
	Config        *ConfigMut     `json:"config,omitempty"`
	Layers        *LayerMut      `json:"layers,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
}

type ConfigMut struct {
	SetEnv        map[string]string `json:"setEnv,omitempty"`
	UnsetEnv      []string          `json:"unsetEnv,omitempty"`
	SetLabel      map[string]string `json:"setLabel,omitempty"`
	UnsetLabel    []string          `json:"unsetLabel,omitempty"`
	SetEntrypoint []string          `json:"setEntrypoint,omitempty"`
	SetCmd        []string          `json:"setCmd,omitempty"`
	SetUser       string            `json:"setUser,omitempty"`
	SetWorkdir    string            `json:"setWorkdir,omitempty"`
}

type LayerMut struct {
	Delete []int         `json:"delete,omitempty"`
	Insert []LayerInsert `json:"insert,omitempty"`
	Append []string      `json:"append,omitempty"`
}

type LayerInsert struct {
	Index int    `json:"index"`
	From  string `json:"from"`
}

func LoadPlan(path string) (*Plan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read plan: %w", err)
	}
	var p Plan
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse plan: %w", err)
	}
	if p.Image == "" {
		return nil, fmt.Errorf("plan.image is required")
	}
	return &p, nil
}
