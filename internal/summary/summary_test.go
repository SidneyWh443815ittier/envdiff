package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/summary"
)

func makeResult(missing, extra []string, mismatched map[string]comparator.ValuePair) comparator.Result {
	return comparator.Result{
		Missing:    missing,
		Extra:      extra,
		Mismatched: mismatched,
	}
}

func TestCompute_Clean(t *testing.T) {
	r := makeResult(nil, nil, nil)
	s := summary.Compute(r)
	if !s.Clean() {
		t.Error("expected clean result")
	}
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
}

func TestCompute_WithIssues(t *testing.T) {
	r := makeResult(
		[]string{"A"},
		[]string{"B", "C"},
		map[string]comparator.ValuePair{"D": {Base: "x", Comp: "y"}},
	)
	s := summary.Compute(r)
	if s.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", s.Missing)
	}
	if s.Extra != 2 {
		t.Errorf("expected Extra=2, got %d", s.Extra)
	}
	if s.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", s.Mismatched)
	}
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
}

func TestWrite_Clean(t *testing.T) {
	var buf bytes.Buffer
	summary.Write(&buf, summary.Stats{})
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestWrite_WithIssues(t *testing.T) {
	var buf bytes.Buffer
	s := summary.Stats{Missing: 2, Extra: 1, Mismatched: 3, Total: 6}
	summary.Write(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "6 issue(s)") {
		t.Errorf("expected issue count in output: %q", out)
	}
	if !strings.Contains(out, "2 missing") {
		t.Errorf("expected missing count in output: %q", out)
	}
	if !strings.Contains(out, "1 extra") {
		t.Errorf("expected extra count in output: %q", out)
	}
	if !strings.Contains(out, "3 mismatched") {
		t.Errorf("expected mismatched count in output: %q", out)
	}
}
