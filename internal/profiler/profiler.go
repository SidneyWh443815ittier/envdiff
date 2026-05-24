// Package profiler collects and reports timing metrics for envdiff pipeline stages.
package profiler

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Stage represents a named timing measurement.
type Stage struct {
	Name     string
	Duration time.Duration
}

// Profiler records elapsed time for named pipeline stages.
type Profiler struct {
	stages []Stage
	start  map[string]time.Time
}

// New returns a new Profiler instance.
func New() *Profiler {
	return &Profiler{
		start: make(map[string]time.Time),
	}
}

// Start begins timing a named stage.
func (p *Profiler) Start(name string) {
	p.start[name] = time.Now()
}

// Stop ends timing for a named stage and records the duration.
// If Start was not called for the name, Stop is a no-op.
func (p *Profiler) Stop(name string) {
	t, ok := p.start[name]
	if !ok {
		return
	}
	p.stages = append(p.stages, Stage{
		Name:     name,
		Duration: time.Since(t),
	})
	delete(p.start, name)
}

// Stages returns all completed stages in the order they were stopped.
func (p *Profiler) Stages() []Stage {
	out := make([]Stage, len(p.stages))
	copy(out, p.stages)
	return out
}

// Total returns the sum of all recorded stage durations.
func (p *Profiler) Total() time.Duration {
	var total time.Duration
	for _, s := range p.stages {
		total += s.Duration
	}
	return total
}

// Write writes a human-readable timing report to w.
func (p *Profiler) Write(w io.Writer) {
	stages := p.Stages()
	sort.Slice(stages, func(i, j int) bool {
		return stages[i].Duration > stages[j].Duration
	})
	fmt.Fprintln(w, "=== Profiler Report ===")
	for _, s := range stages {
		fmt.Fprintf(w, "  %-20s %s\n", s.Name, s.Duration.Round(time.Microsecond))
	}
	fmt.Fprintf(w, "  %-20s %s\n", "TOTAL", p.Total().Round(time.Microsecond))
}
