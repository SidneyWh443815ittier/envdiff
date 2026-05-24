package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/sorter"
)

func makeResult() comparator.Result {
	return comparator.Result{
		Missing: []string{"ZEBRA", "ALPHA"},
		Extra:   []string{"MANGO"},
		Mismatched: map[string][2]string{
			"PORT": {"8080", "9090"},
		},
	}
}

func TestFlatten_ProducesAllCategories(t *testing.T) {
	result := makeResult()
	entries := sorter.Flatten(result)
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
	cats := map[string]int{}
	for _, e := range entries {
		cats[e.Category]++
	}
	if cats["missing"] != 2 || cats["extra"] != 1 || cats["mismatch"] != 1 {
		t.Errorf("unexpected category counts: %v", cats)
	}
}

func TestSort_ByKey_Ascending(t *testing.T) {
	entries := sorter.Flatten(makeResult())
	sorted := sorter.Sort(entries, sorter.Options{SortBy: sorter.SortByKey})
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Key < sorted[i-1].Key {
			t.Errorf("out of order: %s before %s", sorted[i-1].Key, sorted[i].Key)
		}
	}
}

func TestSort_ByKey_Descending(t *testing.T) {
	entries := sorter.Flatten(makeResult())
	sorted := sorter.Sort(entries, sorter.Options{SortBy: sorter.SortByKey, Descending: true})
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Key > sorted[i-1].Key {
			t.Errorf("out of order descending: %s before %s", sorted[i-1].Key, sorted[i].Key)
		}
	}
}

func TestSort_BySeverity_MissingFirst(t *testing.T) {
	entries := sorter.Flatten(makeResult())
	sorted := sorter.Sort(entries, sorter.Options{SortBy: sorter.SortBySeverity})
	if sorted[0].Category != "missing" {
		t.Errorf("expected first entry to be 'missing', got %q", sorted[0].Category)
	}
}

func TestSort_ByCategory_Alphabetical(t *testing.T) {
	entries := sorter.Flatten(makeResult())
	sorted := sorter.Sort(entries, sorter.Options{SortBy: sorter.SortByCategory})
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Category < sorted[i-1].Category {
			t.Errorf("category out of order: %s before %s", sorted[i-1].Category, sorted[i].Category)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := sorter.DefaultOptions()
	if opts.SortBy != sorter.SortByKey {
		t.Errorf("expected SortByKey default, got %v", opts.SortBy)
	}
	if opts.Descending {
		t.Error("expected Descending to be false by default")
	}
}
