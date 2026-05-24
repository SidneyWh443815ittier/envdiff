package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/merger"
)

func TestMerge_EmptyInput(t *testing.T) {
	result, err := merger.Merge(nil, merger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestMerge_SingleMap(t *testing.T) {
	input := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result, err := merger.Merge([]map[string]string{input}, merger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" || result["BAZ"] != "qux" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := map[string]string{"KEY": "first", "ONLY_A": "a"}
	b := map[string]string{"KEY": "last", "ONLY_B": "b"}
	opts := merger.Options{Strategy: merger.StrategyLast}

	result, err := merger.Merge([]map[string]string{a, b}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "last" {
		t.Errorf("expected 'last', got %q", result["KEY"])
	}
	if result["ONLY_A"] != "a" || result["ONLY_B"] != "b" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	opts := merger.Options{Strategy: merger.StrategyFirst}

	result, err := merger.Merge([]map[string]string{a, b}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", result["KEY"])
	}
}

func TestMerge_StrategyError_NoConflict(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same", "EXTRA": "val"}
	opts := merger.Options{Strategy: merger.StrategyError}

	_, err := merger.Merge([]map[string]string{a, b}, opts)
	if err != nil {
		t.Errorf("expected no error for identical values, got: %v", err)
	}
}

func TestMerge_StrategyError_WithConflict(t *testing.T) {
	a := map[string]string{"KEY": "value1"}
	b := map[string]string{"KEY": "value2"}
	opts := merger.Options{Strategy: merger.StrategyError}

	_, err := merger.Merge([]map[string]string{a, b}, opts)
	if err == nil {
		t.Error("expected error for conflicting values, got nil")
	}
}
