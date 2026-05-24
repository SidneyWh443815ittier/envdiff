package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/grouper"
	"github.com/user/envdiff/internal/sorter"
)

func makeEntries() []sorter.Entry {
	return []sorter.Entry{
		{Key: "DB_HOST", Category: "missing"},
		{Key: "DB_PORT", Category: "extra"},
		{Key: "APP_ENV", Category: "missing"},
		{Key: "PORT", Category: "mismatch", BaseVal: "8080", CompVal: "9090"},
	}
}

func TestGroup_ByCategory(t *testing.T) {
	groups := grouper.Group(makeEntries(), grouper.Options{Strategy: grouper.GroupByCategory})
	catMap := map[string]int{}
	for _, g := range groups {
		catMap[g.Name] = len(g.Entries)
	}
	if catMap["missing"] != 2 {
		t.Errorf("expected 2 missing entries, got %d", catMap["missing"])
	}
	if catMap["extra"] != 1 {
		t.Errorf("expected 1 extra entry, got %d", catMap["extra"])
	}
	if catMap["mismatch"] != 1 {
		t.Errorf("expected 1 mismatch entry, got %d", catMap["mismatch"])
	}
}

func TestGroup_ByPrefix(t *testing.T) {
	groups := grouper.Group(makeEntries(), grouper.Options{Strategy: grouper.GroupByPrefix})
	prefixMap := map[string]int{}
	for _, g := range groups {
		prefixMap[g.Name] = len(g.Entries)
	}
	if prefixMap["DB"] != 2 {
		t.Errorf("expected 2 DB entries, got %d", prefixMap["DB"])
	}
	if prefixMap["APP"] != 1 {
		t.Errorf("expected 1 APP entry, got %d", prefixMap["APP"])
	}
	if prefixMap["PORT"] != 1 {
		t.Errorf("expected 1 PORT entry, got %d", prefixMap["PORT"])
	}
}

func TestGroup_EmptyEntries(t *testing.T) {
	groups := grouper.Group(nil, grouper.DefaultOptions())
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestDefaultOptions_IsCategory(t *testing.T) {
	opts := grouper.DefaultOptions()
	if opts.Strategy != grouper.GroupByCategory {
		t.Errorf("expected GroupByCategory default, got %v", opts.Strategy)
	}
}

func TestGroup_PreservesInsertionOrder(t *testing.T) {
	entries := makeEntries()
	groups := grouper.Group(entries, grouper.Options{Strategy: grouper.GroupByCategory})
	if len(groups) == 0 {
		t.Fatal("expected non-empty groups")
	}
	// First category encountered in entries is "missing"
	if groups[0].Name != "missing" {
		t.Errorf("expected first group to be 'missing', got %q", groups[0].Name)
	}
}
