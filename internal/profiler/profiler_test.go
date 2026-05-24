package profiler_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/profiler"
)

func TestNew_EmptyProfiler(t *testing.T) {
	p := profiler.New()
	if len(p.Stages()) != 0 {
		t.Fatal("expected no stages on new profiler")
	}
	if p.Total() != 0 {
		t.Fatal("expected zero total on new profiler")
	}
}

func TestStartStop_RecordsStage(t *testing.T) {
	p := profiler.New()
	p.Start("parse")
	time.Sleep(5 * time.Millisecond)
	p.Stop("parse")

	stages := p.Stages()
	if len(stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(stages))
	}
	if stages[0].Name != "parse" {
		t.Errorf("expected stage name 'parse', got %q", stages[0].Name)
	}
	if stages[0].Duration < 5*time.Millisecond {
		t.Errorf("expected duration >= 5ms, got %s", stages[0].Duration)
	}
}

func TestStop_WithoutStart_IsNoop(t *testing.T) {
	p := profiler.New()
	p.Stop("unknown")
	if len(p.Stages()) != 0 {
		t.Fatal("expected no stages after Stop without Start")
	}
}

func TestTotal_SumsAllStages(t *testing.T) {
	p := profiler.New()
	p.Start("a")
	time.Sleep(2 * time.Millisecond)
	p.Stop("a")
	p.Start("b")
	time.Sleep(2 * time.Millisecond)
	p.Stop("b")

	if p.Total() < 4*time.Millisecond {
		t.Errorf("expected total >= 4ms, got %s", p.Total())
	}
	if len(p.Stages()) != 2 {
		t.Fatalf("expected 2 stages, got %d", len(p.Stages()))
	}
}

func TestWrite_ContainsStageNames(t *testing.T) {
	p := profiler.New()
	p.Start("compare")
	time.Sleep(1 * time.Millisecond)
	p.Stop("compare")
	p.Start("filter")
	time.Sleep(1 * time.Millisecond)
	p.Stop("filter")

	var buf bytes.Buffer
	p.Write(&buf)
	out := buf.String()

	for _, want := range []string{"compare", "filter", "TOTAL", "Profiler Report"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestStages_ReturnsCopy(t *testing.T) {
	p := profiler.New()
	p.Start("load")
	p.Stop("load")

	s1 := p.Stages()
	s1[0].Name = "mutated"
	s2 := p.Stages()

	if s2[0].Name == "mutated" {
		t.Error("Stages() should return a copy, not a reference")
	}
}
