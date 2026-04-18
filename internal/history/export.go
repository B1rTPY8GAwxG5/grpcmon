package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Save writes the current history records as JSON to dir/name.json.
func Save(h *History, dir, name string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("history: mkdir: %w", err)
	}
	records := h.List()
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	dest := filepath.Join(dir, name+".json")
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("history: write: %w", err)
	}
	return nil
}

// LoadInto reads a previously saved history file into h, appending records.
func LoadInto(h *History, dir, name string) error {
	src := filepath.Join(dir, name+".json")
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("history: read: %w", err)
	}
	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("history: unmarshal: %w", err)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, r := range records {
		if len(h.records) >= h.maxSize {
			h.records = h.records[1:]
		}
		h.records = append(h.records, r)
	}
	return nil
}
