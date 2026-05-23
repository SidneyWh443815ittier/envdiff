package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/envdiff/internal/comparator"
	"github.com/envdiff/internal/reporter"
)

func TestReport_Clean(t *testing.T) {
	result := comparator.Result{}
	var buf bytes.Buffer
	reporter.Report(&buf, result, ".env", ".env.production")
	output := buf.String()
	if !strings.Contains(output, "No differences") {
		t.Errorf("expected clean message, got: %s", output)
	}
}

func TestReport_MissingKeys(t *testing.T) {
	result := comparator.Result{
		Missing: map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"},
	}
	var buf bytes.Buffer
	reporter.Report(&buf, result, ".env", ".env.production")
	output := buf.String()
	if !strings.Contains(output, "Missing keys") {
		t.Errorf("expected missing keys section, got: %s", output)
	}
	if !strings.Contains(output, "DB_HOST") || !strings.Contains(output, "API_KEY") {
		t.Errorf("expected missing key names in output, got: %s", output)
	}
}

func TestReport_ExtraKeys(t *testing.T) {
	result := comparator.Result{
		Extra: map[string]string{"NEW_FEATURE": "true"},
	}
	var buf bytes.Buffer
	reporter.Report(&buf, result, ".env", ".env.production")
	output := buf.String()
	if !strings.Contains(output, "Extra keys") {
		t.Errorf("expected extra keys section, got: %s", output)
	}
	if !strings.Contains(output, "NEW_FEATURE") {
		t.Errorf("expected extra key name in output, got: %s", output)
	}
}

func TestReport_MismatchedValues(t *testing.T) {
	result := comparator.Result{
		Mismatched: map[string]comparator.Mismatch{
			"LOG_LEVEL": {Base: "debug", Target: "info"},
		},
	}
	var buf bytes.Buffer
	reporter.Report(&buf, result, ".env", ".env.production")
	output := buf.String()
	if !strings.Contains(output, "Mismatched values") {
		t.Errorf("expected mismatched section, got: %s", output)
	}
	if !strings.Contains(output, "LOG_LEVEL") {
		t.Errorf("expected key name in mismatch output, got: %s", output)
	}
	if !strings.Contains(output, "debug") || !strings.Contains(output, "info") {
		t.Errorf("expected base and target values in output, got: %s", output)
	}
}

func TestExitCode_Clean(t *testing.T) {
	result := comparator.Result{}
	if code := reporter.ExitCode(result); code != 0 {
		t.Errorf("expected exit code 0 for clean result, got %d", code)
	}
}

func TestExitCode_WithDiff(t *testing.T) {
	result := comparator.Result{
		Missing: map[string]string{"KEY": "val"},
	}
	if code := reporter.ExitCode(result); code != 1 {
		t.Errorf("expected exit code 1 for diff result, got %d", code)
	}
}
