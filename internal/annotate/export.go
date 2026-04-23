package annotate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Save writes all annotations from s to the file at path as JSON.
func Save(s *Store, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("annotate: save: %w", err)
	}
	defer f.Close()
	if err := encode(f, s.List()); err != nil {
		return err
	}
	return f.Close()
}

// LoadInto reads annotations from the file at path and adds them to s.
func LoadInto(s *Store, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("annotate: load: %w", err)
	}
	defer f.Close()
	notes, err := decode(f)
	if err != nil {
		return err
	}
	for id, note := range notes {
		if setErr := s.Set(id, note); setErr != nil {
			return setErr
		}
	}
	return nil
}

func encode(w io.Writer, notes map[string]string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(notes)
}

func decode(r io.Reader) (map[string]string, error) {
	var notes map[string]string
	if err := json.NewDecoder(r).Decode(&notes); err != nil {
		return nil, fmt.Errorf("annotate: decode: %w", err)
	}
	return notes, nil
}
