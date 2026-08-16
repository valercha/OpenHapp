package manifest

import (
	"testing"
	"time"

	"github.com/valercha/OpenHapp/daemon/internal/config"
)

func TestDefaultUsesSingBoxWhenEngineIsEmpty(t *testing.T) {
	cfg := config.Default()
	cfg.Engine = ""

	m := Default("0.1.0-test", cfg)

	if m.Engine != "sing-box" {
		t.Fatalf("unexpected engine: %q", m.Engine)
	}
}

func TestDefaultUsesConfiguredEngine(t *testing.T) {
	cfg := config.Default()
	cfg.Engine = "sing-box"
	cfg.Mode = "proxy"

	m := Default("0.1.0-test", cfg)

	if m.Engine != "sing-box" {
		t.Fatalf("unexpected engine: %q", m.Engine)
	}
	if m.Mode != "proxy" {
		t.Fatalf("unexpected mode: %q", m.Mode)
	}
	if m.UpdatedAt.IsZero() {
		t.Fatal("expected manifest timestamp")
	}
	if m.UpdatedAt.After(time.Now().UTC().Add(time.Second)) {
		t.Fatal("manifest timestamp is unexpectedly in the future")
	}
}
