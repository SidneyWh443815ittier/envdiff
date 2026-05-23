package pipeline_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/config"
	"github.com/user/envdiff/internal/pipeline"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestRun_Clean(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	comp := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	cfg := config.Config{
		BaseFile:  base,
		CompFiles: []string{comp},
		Format:    "text",
	}

	var buf bytes.Buffer
	res, err := pipeline.Run(cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasIssues {
		t.Error("expected no issues")
	}
	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}
}

func TestRun_WithDiff_FailOnDiff(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nMISSING=x\n")
	comp := writeTempEnv(t, "FOO=changed\n")

	cfg := config.Config{
		BaseFile:   base,
		CompFiles:  []string{comp},
		Format:     "text",
		FailOnDiff: true,
	}

	var buf bytes.Buffer
	res, err := pipeline.Run(cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.HasIssues {
		t.Error("expected issues")
	}
	if res.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", res.ExitCode)
	}
}

func TestRun_WithDiff_NoFail(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\n")
	comp := writeTempEnv(t, "FOO=other\n")

	cfg := config.Config{
		BaseFile:   base,
		CompFiles:  []string{comp},
		Format:     "text",
		FailOnDiff: false,
	}

	var buf bytes.Buffer
	res, err := pipeline.Run(cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0 when FailOnDiff=false, got %d", res.ExitCode)
	}
}

func TestRun_OutputContainsDiff(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nSECRET=abc\n")
	comp := writeTempEnv(t, "FOO=bar\n")

	cfg := config.Config{
		BaseFile:  base,
		CompFiles: []string{comp},
		Format:    "text",
	}

	var buf bytes.Buffer
	_, err := pipeline.Run(cfg, &buf)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "SECRET") {
		t.Errorf("expected output to mention SECRET, got: %s", buf.String())
	}
}
