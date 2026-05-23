package watcher_test

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/envdiff/internal/watcher"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDefaultOptions(t *testing.T) {
	opts := watcher.DefaultOptions()
	if opts.PollInterval <= 0 {
		t.Error("expected positive poll interval")
	}
	if opts.OnChange == nil {
		t.Error("expected non-nil OnChange")
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	p := writeTempEnv(t, "KEY=value\n")

	var called atomic.Int32
	opts := watcher.Options{
		PollInterval: 50 * time.Millisecond,
		OnChange: func(path string) {
			if path == p {
				called.Add(1)
			}
		},
	}

	w := watcher.New([]string{p}, opts)
	if err := w.Start(); err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer w.Stop()

	// Wait a tick, then modify the file.
	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Allow the watcher to detect the change.
	time.Sleep(150 * time.Millisecond)

	if called.Load() == 0 {
		t.Error("expected OnChange to be called after file modification")
	}
}

func TestWatcher_NoFalsePositive(t *testing.T) {
	p := writeTempEnv(t, "KEY=stable\n")

	var called atomic.Int32
	opts := watcher.Options{
		PollInterval: 50 * time.Millisecond,
		OnChange:     func(_ string) { called.Add(1) },
	}

	w := watcher.New([]string{p}, opts)
	if err := w.Start(); err != nil {
		t.Fatalf("Start() error: %v", err)
	}
	defer w.Stop()

	time.Sleep(200 * time.Millisecond)

	if called.Load() != 0 {
		t.Errorf("expected no change events, got %d", called.Load())
	}
}

func TestWatcher_StartMissingFile(t *testing.T) {
	w := watcher.New([]string{"/nonexistent/.env"}, watcher.DefaultOptions())
	if err := w.Start(); err == nil {
		t.Error("expected error for missing file")
	}
}
