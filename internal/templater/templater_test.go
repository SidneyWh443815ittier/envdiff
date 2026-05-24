package templater_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/templater"
)

func TestGenerate_EmptyMap(t *testing.T) {
	out := templater.Generate(map[string]string{}, templater.DefaultOptions())
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestGenerate_DefaultPlaceholder(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out := templater.Generate(env, templater.DefaultOptions())

	if !strings.Contains(out, "DB_HOST=\n") {
		t.Errorf("expected DB_HOST= with empty value, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PORT=\n") {
		t.Errorf("expected DB_PORT= with empty value, got:\n%s", out)
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"SECRET": "abc123"}
	opts := templater.Options{Placeholder: "CHANGE_ME", AddComments: false}
	out := templater.Generate(env, opts)

	if !strings.Contains(out, "SECRET=CHANGE_ME\n") {
		t.Errorf("expected SECRET=CHANGE_ME, got:\n%s", out)
	}
}

func TestGenerate_WithComments(t *testing.T) {
	env := map[string]string{"API_KEY": "secret"}
	opts := templater.Options{Placeholder: "", AddComments: true}
	out := templater.Generate(env, opts)

	if !strings.Contains(out, "# API_KEY\n") {
		t.Errorf("expected comment line, got:\n%s", out)
	}
	if !strings.Contains(out, "API_KEY=\n") {
		t.Errorf("expected key line, got:\n%s", out)
	}
}

func TestGenerate_SortedKeys(t *testing.T) {
	env := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	out := templater.Generate(env, templater.DefaultOptions())
	lines := strings.Split(strings.TrimSpace(out), "\n")

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line A_KEY, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected last line Z_KEY, got %s", lines[2])
	}
}

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.template")

	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := templater.WriteFile(path, env, templater.DefaultOptions()); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "BAZ=\n") || !strings.Contains(content, "FOO=\n") {
		t.Errorf("unexpected file content:\n%s", content)
	}
}

func TestWriteFile_InvalidPath(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	err := templater.WriteFile("/nonexistent/dir/.env.template", env, templater.DefaultOptions())
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
