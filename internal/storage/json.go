package storage

import (
	"encoding/json"
	"fmt"
	"github.com/nvbn/termonizer/internal/model"
	"os"
)

type JSON struct {
	path string
}

func NewJSON(path string) *JSON {
	return &JSON{
		path: path,
	}
}

// TODO: it will lead to corrupted files but whatever, do backups
func (j *JSON) Read() ([]model.Goals, error) {
	content := make([]model.Goals, 0)
	raw, err := os.ReadFile(j.path)
	if err != nil {
		return content, nil // file doesn't exists
	}

	if err := json.Unmarshal(raw, &content); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return content, nil
}

func (j *JSON) Write(content []model.Goals) error {
	raw, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	return os.WriteFile(j.path, raw, 0644)
}
