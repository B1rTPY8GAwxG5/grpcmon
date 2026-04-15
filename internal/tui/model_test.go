package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/grpcmon/internal/capture"
	"github.com/user/grpcmon/internal/filter"
)

func makeStore(t *testing.T, entries ...capture.Entry) *capture.Store {
	t.Helper()
	s := capture.NewStore(100)
	for _, e := range entries {
		s.Add(e)
	}
	return s
}

func sampleEntry(method, status string) capture.Entry {
	return capture.Entry{
		ID:        "test-id",
		Method:    method,
		Status:    status,
		Timestamp: time.Now(),
		Duration:  10 * time.Millisecond,
		Request:   `{"key":"value"}`,
		Response:  `{"ok":true}`,
	}
}

func TestNew_InitialisesListWithEntries(t *testing.T) {
	s := makeStore(t,
		sampleEntry("/svc.Foo/Bar", "OK"),
		sampleEntry("/svc.Foo/Baz", "NOT_FOUND"),
	)
	m := New(s, filter.Options{})
	if got := len(m.list.Items()); got != 2 {
		t.Fatalf("expected 2 items, got %d", got)
	}
}

func TestUpdate_QuitKey(t *testing.T) {
	s := makeStore(t)
	m := New(s, filter.Options{})
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	m2 := updated.(Model)
	if !m2.quitting {
		t.Fatal("expected quitting to be true after 'q' key")
	}
}

func TestUpdate_RefreshKey(t *testing.T) {
	s := makeStore(t)
	m := New(s, filter.Options{})
	s.Add(sampleEntry("/svc.New/Call", "OK"))
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	m2 := updated.(Model)
	if got := len(m2.list.Items()); got != 1 {
		t.Fatalf("expected 1 item after refresh, got %d", got)
	}
}

func TestEntryItem_Title(t *testing.T) {
	e := sampleEntry("/svc.Foo/Bar", "OK")
	item := entryItem{entry: e}
	title := item.Title()
	if title == "" {
		t.Fatal("expected non-empty title")
	}
}

func TestEntryItem_FilterValue(t *testing.T) {
	e := sampleEntry("/svc.Foo/Bar", "OK")
	item := entryItem{entry: e}
	if item.FilterValue() != "/svc.Foo/Bar" {
		t.Fatalf("unexpected filter value: %s", item.FilterValue())
	}
}
