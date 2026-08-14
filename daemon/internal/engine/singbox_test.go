package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSingBoxBackendDefaults(t *testing.T) {
	b := NewSingBoxBackend()

	if b.Binary != defaultSingBoxBinary {
		t.Fatalf("unexpected binary: %q", b.Binary)
	}

	if b.Config != defaultSingBoxConfig {
		t.Fatalf("unexpected config: %q", b.Config)
	}

	if b.Workdir != defaultSingBoxWorkdir {
		t.Fatalf("unexpected workdir: %q", b.Workdir)
	}
}

func TestSingBoxCheckMissingBinary(t *testing.T) {
	b := &SingBoxBackend{
		Binary:  "/nonexistent/sing-box",
		Config:  "/nonexistent/config.json",
		Workdir: t.TempDir(),
	}

	if err := b.Check(context.Background()); err == nil {
		t.Fatal("expected error for missing binary")
	}
}

func TestSingBoxCheckMissingConfig(t *testing.T) {
	b := &SingBoxBackend{
		Binary:  os.Args[0],
		Config:  filepath.Join(t.TempDir(), "missing.json"),
		Workdir: t.TempDir(),
	}

	if err := b.Check(context.Background()); err == nil {
		t.Fatal("expected error for missing config")
	}
}
