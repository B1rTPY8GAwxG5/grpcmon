// Package snapshot provides functionality for saving and restoring
// capture store state to/from named snapshots on disk.
package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/grpcmon/grpcmon/internal/capture"
	"github.com/grpcmon/grpcmon/internal/export"
)

// Meta holds metadata about a saved snapshot.
type Meta struct {
	Name      string
	CreatedAt time.Time
	Count     int
}

// Save writes all entries from the store to a named snapshot file under dir.
func Save(dir, name string, store *capture.Store) (Meta, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Meta{}, fmt.Errorf("snapshot: mkdir %s: %w", dir, err)
	}

	entries := store.List()
	path := filepath.Join(dir, name+".json")

	f, err := os.Create(path)
	if err != nil {
		return Meta{}, fmt.Errorf("snapshot: create %s: %w", path, err)
	}
	defer f.Close()

	if err := export.Write(f, entries, export.JSON); err != nil {
		return Meta{}, fmt.Errorf("snapshot: write: %w", err)
	}

	return Meta{Name: name, CreatedAt: time.Now(), Count: len(entries)}, nil
}

// Load reads a named snapshot from dir and returns its entries.
func Load(dir, name string) ([]capture.Entry, error) {
	path := filepath.Join(dir, name+".json")

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open %s: %w", path, err)
	}
	defer f.Close()

	entries, err := export.Read(f, export.JSON)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	return entries, nil
}

// List returns the names of all snapshots stored in dir.
func List(dir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("snapshot: list: %w", err)
	}
	names := make([]string, 0, len(matches))
	for _, m := range matches {
		base := filepath.Base(m)
		names = append(names, base[:len(base)-5])
	}
	return names, nil
}
