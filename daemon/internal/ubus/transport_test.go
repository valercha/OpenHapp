package ubus

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/engine"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

func TestHandleJSON(t *testing.T) {
	cfg := config.Default()
	st := state.New("test")
	svc := service.New(cfg, st)
	d := NewDispatcher(New(svc, st, cfg, manifest.FromConfig("test", cfg)))

	payload := []byte(`{"method":"status"}`)
	response, err := d.HandleJSON(context.Background(), payload)
	if err != nil {
		t.Fatalf("handle json: %v", err)
	}

	var decoded Response
	if err := json.Unmarshal(response, &decoded); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if decoded.Error != "" {
		t.Fatalf("unexpected error: %s", decoded.Error)
	}
}

func TestHandleJSONRejectsMalformedRequest(t *testing.T) {
	cfg := config.Default()
	st := state.New("test")
	svc := service.New(cfg, st)
	d := NewDispatcher(New(svc, st, cfg, manifest.FromConfig("test", cfg)))

	if _, err := d.HandleJSON(context.Background(), []byte(`{"method":`)); err == nil {
		t.Fatal("expected malformed JSON error")
	}
}

func TestHandleJSONEngineInfo(t *testing.T) {
	cfg := config.Default()
	cfg.Engine = "sing-box"

	st := state.New("test")
	svc := service.New(cfg, st)

	eng := engine.New("sing-box")
	svc.SetEngine(eng)

	d := NewDispatcher(New(
		svc,
		st,
		cfg,
		manifest.FromConfig("test", cfg),
	))

	payload := []byte(`{"method":"engine_info"}`)
	response, err := d.HandleJSON(context.Background(), payload)
	if err != nil {
		t.Fatalf("handle engine_info: %v", err)
	}

	var decoded Response
	if err := json.Unmarshal(response, &decoded); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if decoded.Error != "" {
		t.Fatalf("unexpected error: %s", decoded.Error)
	}

	raw, err := json.Marshal(decoded.Result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	var info map[string]any
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("decode engine info: %v", err)
	}

	if info["name"] != "sing-box" {
		t.Fatalf("unexpected engine name: %#v", info["name"])
	}

	if info["running"] != false {
		t.Fatalf("unexpected running state: %#v", info["running"])
	}
}
