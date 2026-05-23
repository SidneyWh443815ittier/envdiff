package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Format defines the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatMarkdown Format = "markdown"
)

// Formatter writes comparison results in a specific format.
type Formatter struct {
	Format Format
	Writer io.Writer
}

// New creates a new Formatter.
func New(format Format, w io.Writer) *Formatter {
	return &Formatter{Format: format, Writer: w}
}

// Write renders the comparison result using the configured format.
func (f *Formatter) Write(result comparator.Result) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(result)
	case FormatMarkdown:
		return f.writeMarkdown(result)
	default:
		return f.writeText(result)
	}
}

func (f *Formatter) writeText(result comparator.Result) error {
	if len(result.Missing) == 0 && len(result.Extra) == 0 && len(result.Mismatched) == 0 {
		fmt.Fprintln(f.Writer, "No differences found.")
		return nil
	}
	for _, k := range result.Missing {
		fmt.Fprintf(f.Writer, "MISSING: %s\n", k)
	}
	for _, k := range result.Extra {
		fmt.Fprintf(f.Writer, "EXTRA:   %s\n", k)
	}
	for k, m := range result.Mismatched {
		fmt.Fprintf(f.Writer, "MISMATCH: %s (base=%q, comp=%q)\n", k, m.Base, m.Comp)
	}
	return nil
}

func (f *Formatter) writeJSON(result comparator.Result) error {
	enc := json.NewEncoder(f.Writer)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

func (f *Formatter) writeMarkdown(result comparator.Result) error {
	sb := &strings.Builder{}
	sb.WriteString("## EnvDiff Report\n\n")
	if len(result.Missing) > 0 {
		sb.WriteString("### Missing Keys\n")
		for _, k := range result.Missing {
			fmt.Fprintf(sb, "- `%s`\n", k)
		}
		sb.WriteString("\n")
	}
	if len(result.Extra) > 0 {
		sb.WriteString("### Extra Keys\n")
		for _, k := range result.Extra {
			fmt.Fprintf(sb, "- `%s`\n", k)
		}
		sb.WriteString("\n")
	}
	if len(result.Mismatched) > 0 {
		sb.WriteString("### Mismatched Values\n")
		for k, m := range result.Mismatched {
			fmt.Fprintf(sb, "- `%s`: base=`%s`, comp=`%s`\n", k, m.Base, m.Comp)
		}
		sb.WriteString("\n")
	}
	if len(result.Missing) == 0 && len(result.Extra) == 0 && len(result.Mismatched) == 0 {
		sb.WriteString("_No differences found._\n")
	}
	_, err := fmt.Fprint(f.Writer, sb.String())
	return err
}
