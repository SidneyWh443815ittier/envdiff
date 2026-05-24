package interpolator

import (
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("values should be unchanged, got %v", out)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://example.com",
		"API":      "${BASE_URL}/api",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "https://example.com/api"; out["API"] != want {
		t.Errorf("API: got %q, want %q", out["API"], want)
	}
}

func TestInterpolate_BareStyle(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"ADDR": "$HOST:8080",
	}
	out, err := Interpolate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "localhost:8080"; out["ADDR"] != want {
		t.Errorf("ADDR: got %q, want %q", out["ADDR"], want)
	}
}

func TestInterpolate_MissingRef_NoFail(t *testing.T) {
	env := map[string]string{"URL": "${MISSING}/path"}
	opts := DefaultOptions()
	opts.Placeholder = "UNKNOWN"
	out, err := Interpolate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "UNKNOWN/path"; out["URL"] != want {
		t.Errorf("URL: got %q, want %q", out["URL"], want)
	}
}

func TestInterpolate_MissingRef_FailOnMissing(t *testing.T) {
	env := map[string]string{"URL": "${MISSING}/path"}
	opts := DefaultOptions()
	opts.FailOnMissing = true
	_, err := Interpolate(env, opts)
	if err == nil {
		t.Fatal("expected error for missing reference, got nil")
	}
}

func TestInterpolate_OriginalUnmodified(t *testing.T) {
	env := map[string]string{
		"BASE": "http://base",
		"FULL": "${BASE}/v1",
	}
	original := map[string]string{
		"BASE": "http://base",
		"FULL": "${BASE}/v1",
	}
	_, _ = Interpolate(env, DefaultOptions())
	for k, v := range original {
		if env[k] != v {
			t.Errorf("original mutated: key %q changed to %q", k, env[k])
		}
	}
}
