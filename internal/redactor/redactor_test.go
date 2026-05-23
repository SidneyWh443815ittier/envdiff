package redactor_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/redactor"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":      "myapp",
		"DB_PASSWORD":   "s3cr3t",
		"API_KEY":       "abc123",
		"AUTH_TOKEN":    "tok_xyz",
		"LOG_LEVEL":     "info",
		"PRIVATE_KEY":   "-----BEGIN RSA-----",
		"PLAIN_VAR":     "hello",
	}
}

func TestRedact_DisabledLeavesAllValues(t *testing.T) {
	env := baseEnv()
	opts := redactor.Options{Enabled: false, Patterns: redactor.DefaultSensitivePatterns}
	out := redactor.Redact(env, opts)
	for k, want := range env {
		if got := out[k]; got != want {
			t.Errorf("key %s: got %q, want %q", k, got, want)
		}
	}
}

func TestRedact_SensitiveKeysAreRedacted(t *testing.T) {
	env := baseEnv()
	opts := redactor.DefaultOptions()
	out := redactor.Redact(env, opts)

	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "PRIVATE_KEY"}
	for _, k := range sensitive {
		if out[k] != "***REDACTED***" {
			t.Errorf("expected %s to be redacted, got %q", k, out[k])
		}
	}
}

func TestRedact_NonSensitiveKeysPreserved(t *testing.T) {
	env := baseEnv()
	opts := redactor.DefaultOptions()
	out := redactor.Redact(env, opts)

	plain := []string{"APP_NAME", "LOG_LEVEL", "PLAIN_VAR"}
	for _, k := range plain {
		if out[k] != env[k] {
			t.Errorf("key %s should not be redacted, got %q", k, out[k])
		}
	}
}

func TestRedact_OriginalMapUnmodified(t *testing.T) {
	env := baseEnv()
	orig := env["DB_PASSWORD"]
	opts := redactor.DefaultOptions()
	redactor.Redact(env, opts)
	if env["DB_PASSWORD"] != orig {
		t.Errorf("original map was modified")
	}
}

func TestRedact_CustomPatterns(t *testing.T) {
	env := map[string]string{
		"STRIPE_KEY": "sk_live_abc",
		"APP_ENV":    "production",
	}
	opts := redactor.Options{
		Enabled:  true,
		Patterns: []string{"STRIPE"},
	}
	out := redactor.Redact(env, opts)
	if out["STRIPE_KEY"] != "***REDACTED***" {
		t.Errorf("STRIPE_KEY should be redacted")
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be redacted")
	}
}
