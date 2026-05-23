package comparator_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/comparator"
)

func TestCompare_NoDiff(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := comparator.Compare(base, target)

	if result.HasDiff() {
		t.Errorf("expected no diff, got %+v", result)
	}
}

func TestCompare_MissingKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar", "MISSING": "value"}
	target := map[string]string{"FOO": "bar"}

	result := comparator.Compare(base, target)

	if len(result.Missing) != 1 || result.Missing[0] != "MISSING" {
		t.Errorf("expected MISSING in Missing, got %v", result.Missing)
	}
	if len(result.Extra) != 0 {
		t.Errorf("expected no extra keys, got %v", result.Extra)
	}
}

func TestCompare_ExtraKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "EXTRA": "value"}

	result := comparator.Compare(base, target)

	if len(result.Extra) != 1 || result.Extra[0] != "EXTRA" {
		t.Errorf("expected EXTRA in Extra, got %v", result.Extra)
	}
	if len(result.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", result.Missing)
	}
}

func TestCompare_ValueMismatch(t *testing.T) {
	base := map[string]string{"FOO": "original"}
	target := map[string]string{"FOO": "changed"}

	result := comparator.Compare(base, target)

	diff, ok := result.Mismatch["FOO"]
	if !ok {
		t.Fatal("expected FOO in Mismatch")
	}
	if diff.Base != "original" || diff.Target != "changed" {
		t.Errorf("unexpected diff values: %+v", diff)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := comparator.Compare(map[string]string{}, map[string]string{})
	if result.HasDiff() {
		t.Errorf("expected no diff for empty maps, got %+v", result)
	}
}

func TestHasDiff_True(t *testing.T) {
	result := comparator.Result{
		Missing:  []string{"KEY"},
		Mismatch: make(map[string]comparator.Diff),
	}
	if !result.HasDiff() {
		t.Error("expected HasDiff to return true")
	}
}
