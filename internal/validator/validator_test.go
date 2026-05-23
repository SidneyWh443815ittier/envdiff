package validator_test

import (
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func TestValidate_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"MAX_CONN": "10",
	}
	issues := validator.Validate(env, validator.DefaultOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "",
		"APP_NAME":   "myapp",
	}
	opts := validator.Options{WarnEmptyValues: true, WarnKeyFormat: false}
	issues := validator.Validate(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "SECRET_KEY" {
		t.Errorf("expected issue for SECRET_KEY, got %s", issues[0].Key)
	}
}

func TestValidate_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"db_host": "localhost",
		"APP_ENV": "dev",
	}
	opts := validator.Options{WarnEmptyValues: false, WarnKeyFormat: true}
	issues := validator.Validate(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "db_host" {
		t.Errorf("expected issue for db_host, got %s", issues[0].Key)
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"bad_key": "",
		"GOOD_KEY": "value",
	}
	issues := validator.Validate(env, validator.DefaultOptions())
	// bad_key triggers both empty-value and key-format warnings
	if len(issues) < 2 {
		t.Fatalf("expected at least 2 issues, got %d", len(issues))
	}
}

func TestValidate_DisabledOptions(t *testing.T) {
	env := map[string]string{
		"bad_key": "",
	}
	opts := validator.Options{WarnEmptyValues: false, WarnKeyFormat: false}
	issues := validator.Validate(env, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues when all checks disabled, got %d", len(issues))
	}
}

func TestIssue_String(t *testing.T) {
	i := validator.Issue{Key: "MY_KEY", Message: "value is empty"}
	got := i.String()
	want := "[MY_KEY] value is empty"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
