package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/output"
)

func cleanResult() comparator.Result {
	return comparator.Result{}
}

func diffResult() comparator.Result {
	return comparator.Result{
		Missing: []string{"DB_HOST"},
		Extra:   []string{"OLD_KEY"},
		Mismatched: map[string]comparator.Mismatch{
			"LOG_LEVEL": {Base: "info", Comp: "debug"},
		},
	}
}

func TestFormatter_Text_Clean(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(output.FormatText, &buf)
	if err := f.Write(cleanResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected clean message, got: %s", buf.String())
	}
}

func TestFormatter_Text_Diff(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(output.FormatText, &buf)
	if err := f.Write(diffResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MISSING: DB_HOST") {
		t.Errorf("expected MISSING line, got: %s", out)
	}
	if !strings.Contains(out, "EXTRA:   OLD_KEY") {
		t.Errorf("expected EXTRA line, got: %s", out)
	}
	if !strings.Contains(out, "MISMATCH: LOG_LEVEL") {
		t.Errorf("expected MISMATCH line, got: %s", out)
	}
}

func TestFormatter_JSON_Diff(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(output.FormatJSON, &buf)
	if err := f.Write(diffResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in JSON output, got: %s", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}

func TestFormatter_Markdown_Diff(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(output.FormatMarkdown, &buf)
	if err := f.Write(diffResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "## EnvDiff Report") {
		t.Errorf("expected markdown header, got: %s", out)
	}
	if !strings.Contains(out, "### Missing Keys") {
		t.Errorf("expected missing section, got: %s", out)
	}
}

func TestFormatter_Markdown_Clean(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(output.FormatMarkdown, &buf)
	if err := f.Write(cleanResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected clean markdown, got: %s", buf.String())
	}
}
