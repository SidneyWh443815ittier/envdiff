// Package watcher monitors .env files for changes and triggers a diff callback.
package watcher

import (
	"fmt"
	"os"
	"time"
)

// Options configures the file watcher.
type Options struct {
	// PollInterval is how often to check for file changes.
	PollInterval time.Duration
	// OnChange is called when any watched file changes.
	OnChange func(path string)
}

// DefaultOptions returns sensible watcher defaults.
func DefaultOptions() Options {
	return Options{
		PollInterval: 2 * time.Second,
		OnChange:     func(path string) {},
	}
}

// Watcher polls a set of file paths for modification time changes.
type Watcher struct {
	paths   []string
	opts    Options
	mtimes  map[string]time.Time
	stopCh  chan struct{}
}

// New creates a new Watcher for the given paths.
func New(paths []string, opts Options) *Watcher {
	return &Watcher{
		paths:  paths,
		opts:   opts,
		mtimes: make(map[string]time.Time),
		stopCh: make(chan struct{}),
	}
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() error {
	if err := w.snapshot(); err != nil {
		return fmt.Errorf("watcher: initial snapshot failed: %w", err)
	}
	go w.loop()
	return nil
}

// Stop signals the watcher to halt.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

func (w *Watcher) loop() {
	ticker := time.NewTicker(w.opts.PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.check()
		case <-w.stopCh:
			return
		}
	}
}

func (w *Watcher) snapshot() error {
	for _, p := range w.paths {
		info, err := os.Stat(p)
		if err != nil {
			return err
		}
		w.mtimes[p] = info.ModTime()
	}
	return nil
}

func (w *Watcher) check() {
	for _, p := range w.paths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if prev, ok := w.mtimes[p]; ok && info.ModTime().After(prev) {
			w.mtimes[p] = info.ModTime()
			w.opts.OnChange(p)
		}
	}
}
