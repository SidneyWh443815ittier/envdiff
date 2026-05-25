package patcher_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/patcher"
	"github.com/user/envdiff/internal/parser"
)

// TestPatch_RoundTrip verifies that after patching, the file can be parsed
// and the updated values are returned correctly by the parser.
func TestPatch_RoundTrip(t *testing.T) {
	initial := strings.Join([]string{
		"# environment",
		"APP_ENV=development",
		"LOG_LEVEL=debug",
		"PORT=8080",
		"",
		"# secrets",
		"DB_PASS=secret",
	}, "\n") + "\n"

	p := writeTempEnv(t, initial)

	changes := map[string]string{
		"APP_ENV":   "production",
		"LOG_LEVEL": "warn",
		"NEW_VAR":   "added",
	}

	if err := patcher.Patch(p, changes, patcher.DefaultOptions()); err != nil {
		t.Fatalf("Patch: %v", err)
	}

	env, err := parser.ParseFile(p)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	cases := []struct{ key, want string }{
		{"APP_ENV", "production"},
		{"LOG_LEVEL", "warn"},
		{"PORT", "8080"},
		{"DB_PASS", "secret"},
		{"NEW_VAR", "added"},
	}

	for _, tc := range cases {
		got, ok := env[tc.key]
		if !ok {
			t.Errorf("key %q missing after patch", tc.key)
			continue
		}
		if got != tc.want {
			t.Errorf("key %q: want %q, got %q", tc.key, tc.want, got)
		}
	}
}

// TestPatch_MultiplePatches ensures sequential patches accumulate correctly.
func TestPatch_MultiplePatches(t *testing.T) {
	p := writeTempEnv(t, "APP_ENV=dev\nDEBUG=true\n")

	if err := patcher.Patch(p, map[string]string{"APP_ENV": "staging"}, patcher.DefaultOptions()); err != nil {
		t.Fatalf("first Patch: %v", err)
	}
	if err := patcher.Patch(p, map[string]string{"DEBUG": "false", "EXTRA": "1"}, patcher.DefaultOptions()); err != nil {
		t.Fatalf("second Patch: %v", err)
	}

	env, err := parser.ParseFile(p)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	if env["APP_ENV"] != "staging" {
		t.Errorf("APP_ENV: want staging, got %q", env["APP_ENV"])
	}
	if env["DEBUG"] != "false" {
		t.Errorf("DEBUG: want false, got %q", env["DEBUG"])
	}
	if env["EXTRA"] != "1" {
		t.Errorf("EXTRA: want 1, got %q", env["EXTRA"])
	}
}
