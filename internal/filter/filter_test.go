package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/filter"
)

func baseResult() comparator.Result {
	return comparator.Result{
		Missing: []string{"APP_PORT", "DB_HOST", "SECRET_KEY"},
		Extra:   []string{"APP_DEBUG", "LEGACY_FLAG"},
		Mismatched: []comparator.Mismatch{
			{Key: "APP_ENV", BaseValue: "production", OtherValue: "staging"},
			{Key: "DB_PASS", BaseValue: "secret", OtherValue: "other"},
		},
	}
}

func TestApply_NoFilter(t *testing.T) {
	opts := filter.NewOptions("", nil)
	result := filter.Apply(baseResult(), opts)

	if len(result.Missing) != 3 {
		t.Errorf("expected 3 missing, got %d", len(result.Missing))
	}
	if len(result.Extra) != 2 {
		t.Errorf("expected 2 extra, got %d", len(result.Extra))
	}
	if len(result.Mismatched) != 2 {
		t.Errorf("expected 2 mismatched, got %d", len(result.Mismatched))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	opts := filter.NewOptions("APP_", nil)
	result := filter.Apply(baseResult(), opts)

	if len(result.Missing) != 1 || result.Missing[0] != "APP_PORT" {
		t.Errorf("expected only APP_PORT in missing, got %v", result.Missing)
	}
	if len(result.Extra) != 1 || result.Extra[0] != "APP_DEBUG" {
		t.Errorf("expected only APP_DEBUG in extra, got %v", result.Extra)
	}
	if len(result.Mismatched) != 1 || result.Mismatched[0].Key != "APP_ENV" {
		t.Errorf("expected only APP_ENV in mismatched, got %v", result.Mismatched)
	}
}

func TestApply_IgnoreKeys(t *testing.T) {
	opts := filter.NewOptions("", []string{"SECRET_KEY", "DB_PASS", "LEGACY_FLAG"})
	result := filter.Apply(baseResult(), opts)

	for _, k := range result.Missing {
		if k == "SECRET_KEY" {
			t.Error("SECRET_KEY should be ignored in missing")
		}
	}
	for _, k := range result.Extra {
		if k == "LEGACY_FLAG" {
			t.Error("LEGACY_FLAG should be ignored in extra")
		}
	}
	for _, m := range result.Mismatched {
		if m.Key == "DB_PASS" {
			t.Error("DB_PASS should be ignored in mismatched")
		}
	}
}

func TestApply_PrefixAndIgnore(t *testing.T) {
	opts := filter.NewOptions("DB_", []string{"DB_HOST"})
	result := filter.Apply(baseResult(), opts)

	if len(result.Missing) != 0 {
		t.Errorf("expected 0 missing after prefix+ignore, got %v", result.Missing)
	}
	if len(result.Mismatched) != 1 || result.Mismatched[0].Key != "DB_PASS" {
		t.Errorf("expected DB_PASS in mismatched, got %v", result.Mismatched)
	}
}
