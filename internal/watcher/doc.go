// Package watcher provides a lightweight polling-based file watcher
// that notifies callers when monitored .env files are modified on disk.
// It is intended for use in watch mode where envdiff re-runs automatically
// on file changes without requiring inotify or OS-level FS events.
package watcher
