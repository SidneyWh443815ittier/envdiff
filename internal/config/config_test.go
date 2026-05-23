package config

import (
	"testing"
)

func TestParse_BasicArgs(t *testing.T) {
	cfg, err := Parse([]string{".env.base", ".env.prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseFile != ".env.base" {
		t.Errorf("expected BaseFile '.env.base', got %q", cfg.BaseFile)
	}
	if len(cfg.CompFiles) != 1 || cfg.CompFiles[0] != ".env.prod" {
		t.Errorf("expected CompFiles ['.env.prod'], got %v", cfg.CompFiles)
	}
}

func TestParse_MultipleCompFiles(t *testing.T) {
	cfg, err := Parse([]string{".env", ".env.staging", ".env.prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.CompFiles) != 2 {
		t.Errorf("expected 2 comp files, got %d", len(cfg.CompFiles))
	}
}

func TestParse_WithFlags(t *testing.T) {
	cfg, err := Parse([]string{"-prefix", "APP_", "-fail", "-quiet", ".env", ".env.prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Prefix != "APP_" {
		t.Errorf("expected prefix 'APP_', got %q", cfg.Prefix)
	}
	if !cfg.FailOnDiff {
		t.Error("expected FailOnDiff to be true")
	}
	if !cfg.Quiet {
		t.Error("expected Quiet to be true")
	}
}

func TestParse_IgnoreKeys(t *testing.T) {
	cfg, err := Parse([]string{"-ignore", "SECRET,TOKEN, DEBUG", ".env", ".env.prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.IgnoreKeys) != 3 {
		t.Fatalf("expected 3 ignore keys, got %d", len(cfg.IgnoreKeys))
	}
	if cfg.IgnoreKeys[2] != "DEBUG" {
		t.Errorf("expected trimmed key 'DEBUG', got %q", cfg.IgnoreKeys[2])
	}
}

func TestParse_MissingArgs(t *testing.T) {
	_, err := Parse([]string{".env.base"})
	if err == nil {
		t.Error("expected error for missing compare file, got nil")
	}
}

func TestParse_NoArgs(t *testing.T) {
	_, err := Parse([]string{})
	if err == nil {
		t.Error("expected error for no arguments, got nil")
	}
}

func TestParse_JSONFlag(t *testing.T) {
	cfg, err := Parse([]string{"-json", ".env", ".env.prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.JSONOutput {
		t.Error("expected JSONOutput to be true")
	}
}
