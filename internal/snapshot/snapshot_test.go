package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/grpcmon/grpcmon/internal/capture"
	"github.com/grpcmon/grpcmon/internal/snapshot"
	"google.golang.org/grpc/codes"
)

func makeStore(t *testing.T, n int) *capture.Store {
	t.Helper()
	s := capture.NewStore(100)
	for i := 0; i < n; i++ {
		s.Add(capture.Entry{
			ID:        capture.NewID(),
			Method:    "/pkg.Svc/Method",
			Status:    codes.OK,
			Timestamp: time.Now(),
			LatencyMs: int64(i * 10),
		})
	}
	return s
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	store := makeStore(t, 3)

	meta, err := snapshot.Save(dir, "test-snap", store)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if meta.Count != 3 {
		t.Errorf("Count = %d, want 3", meta.Count)
	}
	if meta.Name != "test-snap" {
		t.Errorf("Name = %q, want test-snap", meta.Name)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := makeStore(t, 5)

	if _, err := snapshot.Save(dir, "snap", store); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := snapshot.Load(dir, "snap")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 5 {
		t.Errorf("got %d entries, want 5", len(entries))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	dir := t.TempDir()
	_, err := snapshot.Load(dir, "nonexistent")
	if err == nil {
		t.Error("expected error for missing snapshot, got nil")
	}
}

func TestSave_DuplicateName(t *testing.T) {
	dir := t.TempDir()
	store := makeStore(t, 2)

	if _, err := snapshot.Save(dir, "dup", store); err != nil {
		t.Fatalf("first Save: %v", err)
	}
	_, err := snapshot.Save(dir, "dup", store)
	if err == nil {
		t.Error("expected error when saving duplicate snapshot name, got nil")
	}
}

func TestList_ReturnsNames(t *testing.T) {
	dir := t.TempDir()
	store := makeStore(t, 1)

	for _, name := range []string{"alpha", "beta", "gamma"} {
		if _, err := snapshot.Save(dir, name, store); err != nil {
			t.Fatalf("Save %s: %v", name, err)
		}
	}

	names, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("got %d names, want 3", len(names))
	}
}

func TestList_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	names, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestList_IgnoresNonJSON(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(dir+"/ignore.txt", []byte("x"), 0o644)
	store := makeStore(t, 1)
	if _, err := snapshot.Save(dir, "real", store); err != nil {
		t.Fatalf("Save: %v", err)
	}
	names, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 1 || names[0] != "real" {
		t.Errorf("got %v, want [real]", names)
	}
}
