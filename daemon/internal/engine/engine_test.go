package engine

import (
	"context"
	"testing"
)

func TestNewDefaultsToSingBox(t *testing.T) {
	e := New("")

	if e.Name() != "sing-box" {
		t.Fatalf("unexpected engine name: %q", e.Name())
	}

	if _, ok := e.Backend().(*SingBoxBackend); !ok {
		t.Fatalf("expected Sing-box backend, got %T", e.Backend())
	}
}

func TestStartFailsForUnsupportedEngine(t *testing.T) {
	e := New("xray")

	if err := e.Start(context.Background()); err == nil {
		t.Fatal("expected unsupported engine error")
	}

	if e.Running() {
		t.Fatal("unsupported engine must remain stopped")
	}
}

func TestSetNameSelectsSingBoxBackend(t *testing.T) {
	e := New("xray")
	e.SetName("sing-box")

	if e.Name() != "sing-box" {
		t.Fatalf("unexpected engine name: %q", e.Name())
	}

	if _, ok := e.Backend().(*SingBoxBackend); !ok {
		t.Fatalf("unexpected backend type: %T", e.Backend())
	}
}

func TestSetNameSelectsSingBoxAlias(t *testing.T) {
	e := New("singbox")

	if _, ok := e.Backend().(*SingBoxBackend); !ok {
		t.Fatalf("unexpected backend type: %T", e.Backend())
	}
}

func TestStartChecksSingBoxBackend(t *testing.T) {
	e := New("sing-box")

	// Replace the real backend with a deliberately invalid one so this
	// test remains independent of the host's Sing-box installation.
	e.backend = &SingBoxBackend{
		Binary:  "/nonexistent/sing-box",
		Config:  "/nonexistent/config.json",
		Workdir: t.TempDir(),
	}

	if err := e.Start(context.Background()); err == nil {
		t.Fatal("expected Sing-box check failure")
	}

	if e.Running() {
		t.Fatal("engine must remain stopped after backend check failure")
	}
}
