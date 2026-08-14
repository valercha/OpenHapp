package service

import (
	"context"
	"testing"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/engine"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

func TestStartFailsWhenEngineCheckFails(t *testing.T) {
	cfg := config.Default()
	cfg.Engine = "sing-box"

	st := state.New("test")
	svc := New(cfg, st)

	eng := engine.New("sing-box")
	engBackend := eng.Backend().(*engine.SingBoxBackend)
	engBackend.Binary = "/nonexistent/sing-box"
	engBackend.Config = "/nonexistent/config.json"
	svc.SetEngine(eng)

	if err := svc.Start(context.Background()); err == nil {
		t.Fatal("expected engine check error")
	}

	if svc.Running() {
		t.Fatal("service must remain stopped when engine check fails")
	}

	if st.Snapshot().Running {
		t.Fatal("state must remain stopped when engine check fails")
	}
}
