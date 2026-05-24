package linter_test

import (
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func TestLint_NoIssues(t *testing.T) {
	lines := []string{
		"APP_NAME=myapp",
		"PORT=8080",
		"# a comment",
		"",
	}
	issues := linter.Lint(lines, linter.DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %+v", len(issues), issues)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	lines := []string{
		"APP_NAME=first",
		"APP_NAME=second",
	}
	issues := linter.Lint(lines, linter.DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "APP_NAME" {
		t.Errorf("expected key APP_NAME, got %s", issues[0].Key)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	lines := []string{"SECRET="}
	issues := linter.Lint(lines, linter.DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %s", issues[0].Key)
	}
}

func TestLint_UnquotedSpaces(t *testing.T) {
	lines := []string{"APP_DESC=hello world"}
	issues := linter.Lint(lines, linter.DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	lines := []string{"appName=value"}
	issues := linter.Lint(lines, linter.DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLint_MalformedLine(t *testing.T) {
	lines := []string{"NODEQUALS"}
	issues := linter.Lint(lines, linter.DefaultOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestLint_DisabledOptions(t *testing.T) {
	lines := []string{
		"lower=value",
		"lower=dup",
	}
	opts := linter.Options{}
	issues := linter.Lint(lines, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues with all rules disabled, got %d", len(issues))
	}
}
