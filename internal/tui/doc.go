// Package tui implements a lightweight terminal user interface for grpcmon.
//
// It renders a scrollable list of captured gRPC calls sourced from a
// capture.Store, supports basic filtering via filter.Options, and allows the
// user to inspect individual entry details interactively.
//
// # Usage
//
//	store := capture.NewStore(500)
//	model := tui.New(store, filter.Options{Method: "/svc.Greeter/SayHello"})
//	if err := tea.NewProgram(model).Start(); err != nil {
//	    log.Fatal(err)
//	}
//
// Key bindings:
//
//	[r]      refresh the list from the store
//	[enter]  view full detail for the selected entry
//	[q]      quit the application
package tui
