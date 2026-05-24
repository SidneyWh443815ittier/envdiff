package masker_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/masker"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_NAME":     "envdiff",
		"DB_PASSWORD":  "s3cr3t",
		"API_KEY":      "abc123",
		"GITHUB_TOKEN": "ghp_xyz",
		"LOG_LEVEL":    "info",
	}
}

func TestMask_DisabledLeavesAllValues(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.Enabled = false

	result := masker.Mask(baseEnv(), opts)
	if result["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected original value, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "abc123" {
		t.Errorf("expected original value, got %q", result["API_KEY"])
	}
}

func TestMask_SensitiveKeysAreMasked(t *testing.T) {
	opts := masker.DefaultOptions()
	result := masker.Mask(baseEnv(), opts)

	sensitive := []string{"DB_PASSWORD", "API_KEY", "GITHUB_TOKEN"}
	for _, k := range sensitive {
		if result[k] != masker.DefaultMask {
			t.Errorf("key %q: expected mask %q, got %q", k, masker.DefaultMask, result[k])
		}
	}
}

func TestMask_NonSensitiveKeysPreserved(t *testing.T) {
	opts := masker.DefaultOptions()
	result := masker.Mask(baseEnv(), opts)

	if result["APP_NAME"] != "envdiff" {
		t.Errorf("expected 'envdiff', got %q", result["APP_NAME"])
	}
	if result["LOG_LEVEL"] != "info" {
		t.Errorf("expected 'info', got %q", result["LOG_LEVEL"])
	}
}

func TestMask_OriginalMapUnmodified(t *testing.T) {
	opts := masker.DefaultOptions()
	env := baseEnv()
	masker.Mask(env, opts)

	if env["DB_PASSWORD"] != "s3cr3t" {
		t.Error("original map was modified")
	}
}

func TestMask_CustomMaskString(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.Mask = "[REDACTED]"
	result := masker.Mask(baseEnv(), opts)

	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected '[REDACTED]', got %q", result["DB_PASSWORD"])
	}
}

func TestIsSensitiveKey(t *testing.T) {
	subs := masker.DefaultOptions().SensitiveSubstrings

	cases := []struct {
		key      string
		expected bool
	}{
		{"DB_PASSWORD", true},
		{"STRIPE_SECRET_KEY", true},
		{"AUTH_TOKEN", true},
		{"APP_NAME", false},
		{"PORT", false},
	}
	for _, tc := range cases {
		got := masker.IsSensitiveKey(tc.key, subs)
		if got != tc.expected {
			t.Errorf("IsSensitiveKey(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}
