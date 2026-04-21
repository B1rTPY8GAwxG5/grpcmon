// Package main is the entry point for the grpcmon CLI tool.
// It wires together configuration, capture, TUI, and export subsystems
// to provide a lightweight gRPC traffic monitor for development environments.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/grpcmon/internal/capture"
	"github.com/example/grpcmon/internal/config"
	"github.com/example/grpcmon/internal/export"
	"github.com/example/grpcmon/internal/snapshot"
	"github.com/example/grpcmon/internal/tui"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "grpcmon: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("grpcmon", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		target     = fs.String("target", "", "gRPC server address to monitor (host:port)")
		maxEntries = fs.Int("max-entries", 0, "maximum captured entries to keep in memory (0 = use config default)")
		exportPath = fs.String("export", "", "export captured entries to file on exit (e.g. out.json)")
		saveName   = fs.String("snapshot", "", "save a named snapshot of captured entries on exit")
		cfgFile    = fs.String("config", "", "path to config file (optional)")
		noTUI      = fs.Bool("no-tui", false, "disable interactive TUI; print entries to stdout")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Load configuration, merging file → env → flags.
	cfg := config.Defaults()
	if *cfgFile != "" {
		loaded, err := config.Load(*cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		cfg = loaded
	}
	config.FromEnv(&cfg)

	if *target != "" {
		cfg.Target = *target
	}
	if *maxEntries > 0 {
		cfg.MaxEntries = *maxEntries
	}

	if err := config.Validate(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Set up the capture store.
	store := capture.NewStore(cfg.MaxEntries)

	// Root context cancelled on SIGINT / SIGTERM.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if *noTUI {
		// Headless mode: block until signal, then export if requested.
		fmt.Fprintf(os.Stderr, "grpcmon: listening (target=%s) — press Ctrl+C to stop\n", cfg.Target)
		<-ctx.Done()
	} else {
		// Interactive TUI mode.
		model := tui.New(store)
		if err := model.Run(ctx); err != nil {
			return fmt.Errorf("tui: %w", err)
		}
	}

	// Post-run export / snapshot.
	entries := store.List()

	if *exportPath != "" && len(entries) > 0 {
		f, err := os.Create(*exportPath)
		if err != nil {
			return fmt.Errorf("creating export file: %w", err)
		}
		defer f.Close()
		if err := export.Write(f, cfg.ExportFormat, entries); err != nil {
			return fmt.Errorf("exporting entries: %w", err)
		}
		fmt.Fprintf(os.Stderr, "grpcmon: exported %d entries to %s\n", len(entries), *exportPath)
	}

	if *saveName != "" && len(entries) > 0 {
		if err := snapshot.Save(*saveName, store); err != nil {
			return fmt.Errorf("saving snapshot: %w", err)
		}
		fmt.Fprintf(os.Stderr, "grpcmon: snapshot %q saved (%d entries)\n", *saveName, len(entries))
	}

	return nil
}
