package resolver

import (
	"testing"
)

func TestResolve_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	out, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" || out["PORT"] != "8080" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestResolve_SimpleReference(t *testing.T) {
	env := map[string]string{
		"SCHEME":   "https",
		"BASE_URL": "${SCHEME}://example.com",
	}
	out, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["BASE_URL"]; got != "https://example.com" {
		t.Errorf("expected https://example.com, got %q", got)
	}
}

func TestResolve_ChainedReferences(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "${A}_world",
		"C": "${B}!",
	}
	out, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["C"]; got != "hello_world!" {
		t.Errorf("expected hello_world!, got %q", got)
	}
}

func TestResolve_MissingRef_NoFail(t *testing.T) {
	env := map[string]string{
		"URL": "${MISSING}://host",
	}
	opts := DefaultOptions()
	opts.FailOnMissing = false
	out, err := Resolve(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// reference left as-is
	if got := out["URL"]; got != "${MISSING}://host" {
		t.Errorf("expected reference preserved, got %q", got)
	}
}

func TestResolve_MissingRef_FailOnMissing(t *testing.T) {
	env := map[string]string{
		"URL": "${MISSING}://host",
	}
	opts := DefaultOptions()
	opts.FailOnMissing = true
	_, err := Resolve(env, opts)
	if err == nil {
		t.Fatal("expected error for missing reference, got nil")
	}
}

func TestResolve_MaxDepthExceeded(t *testing.T) {
	// A -> B -> A creates a cycle; depth limit should trigger.
	env := map[string]string{
		"A": "${B}",
		"B": "${A}",
	}
	opts := DefaultOptions()
	opts.MaxDepth = 3
	_, err := Resolve(env, opts)
	if err == nil {
		t.Fatal("expected depth error for cyclic reference, got nil")
	}
}

func TestResolve_OriginalMapUnmodified(t *testing.T) {
	env := map[string]string{
		"SCHEME":   "http",
		"BASE_URL": "${SCHEME}://localhost",
	}
	original := env["BASE_URL"]
	_, _ = Resolve(env, DefaultOptions())
	if env["BASE_URL"] != original {
		t.Errorf("original map was mutated")
	}
}
