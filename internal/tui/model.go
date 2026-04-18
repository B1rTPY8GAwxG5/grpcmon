// Package tui provides a terminal user interface for browsing captured gRPC entries.
package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/grpcmon/internal/capture"
	"github.com/user/grpcmon/internal/filter"
)

// Model holds the TUI application state.
type Model struct {
	store   *capture.Store
	list    list.Model
	filter  filter.Options
	selected *capture.Entry
	quitting bool
}

// New creates a new TUI Model backed by the given store.
func New(store *capture.Store, opts filter.Options) Model {
	delegate := list.NewDefaultDelegate()
	l := list.New(nil, delegate, 80, 20)
	l.Title = "grpcmon — captured gRPC calls"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)

	m := Model{
		store:  store,
		list:   l,
		filter: opts,
	}
	m.refresh()
	return m
}

func (m *Model) refresh() {
	entries := filter.Apply(m.store.List(), m.filter)
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = entryItem{entry: e}
	}
	m.list.SetItems(items)
}

// Init satisfies tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update satisfies tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// When viewing entry detail, handle navigation keys separately.
		if m.selected != nil {
			switch msg.String() {
			case "backspace", "esc":
				m.selected = nil
				return m, nil
			case "q", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			m.refresh()
			return m, nil
		case "enter":
			if i, ok := m.list.SelectedItem().(entryItem); ok {
				m.selected = &i.entry
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View satisfies tea.Model.
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}
	if m.selected != nil {
		v := fmt.Sprintf("Method:   %s\nStatus:   %s\nDuration: %s\nRequest:  %s\nResponse: %s\n\n[backspace] back  [q] quit\n",
			m.selected.Method,
			m.selected.Status,
			m.selected.Duration,
			m.selected.Request,
			m.selected.Response,
		)
		return v
	}
	return m.list.View() + "\n[r] refresh  [enter] detail  [q] quit\n"
}
