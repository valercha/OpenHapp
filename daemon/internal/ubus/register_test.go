package ubus

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/engine"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/profile"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

func TestDispatcherMethods(t *testing.T) {
	cfg := config.Default()
	st := state.New("test")
	svc := service.New(cfg, st)
	profiles := profile.NewStore(t.TempDir() + "/openhapp")
	srv := New(svc, st, cfg, manifest.FromConfig("test", cfg), profiles)
	d := NewDispatcher(srv)
	ctx := context.Background()

	methods := []string{"status", "config", "manifest", "snapshot", "engine_info", "start", "stop"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			if _, err := d.Dispatch(ctx, method, nil); err != nil {
				t.Fatalf("dispatch %s: %v", method, err)
			}
		})
	}
}

func TestDispatcherEngineInfo(t *testing.T) {
	cfg := config.Default()
	cfg.Engine = "sing-box"

	st := state.New("test")
	svc := service.New(cfg, st)

	eng := engine.New("sing-box")
	svc.SetEngine(eng)

	profiles := profile.NewStore(t.TempDir() + "/openhapp")

	srv := New(
		svc,
		st,
		cfg,
		manifest.FromConfig("test", cfg),
		profiles,
	)

	d := NewDispatcher(srv)

	result, err := d.Dispatch(
		context.Background(),
		"engine_info",
		nil,
	)
	if err != nil {
		t.Fatalf("dispatch engine_info: %v", err)
	}

	info, ok := result.(engine.Info)
	if !ok {
		t.Fatalf("unexpected result type: %T", result)
	}

	if info.Name != "sing-box" {
		t.Fatalf("unexpected engine name: %q", info.Name)
	}

	if info.Running {
		t.Fatal("engine should not be running")
	}
}

func TestDispatcherRejectsUnknownMethod(t *testing.T) {
	cfg := config.Default()
	st := state.New("test")
	svc := service.New(cfg, st)
	profiles := profile.NewStore(t.TempDir() + "/openhapp")
	d := NewDispatcher(New(
		svc,
		st,
		cfg,
		manifest.FromConfig("test", cfg),
		profiles,
	))

	if _, err := d.Dispatch(context.Background(), "unknown", nil); err == nil {
		t.Fatal("expected unknown method error")
	}
}

func TestDispatcherProfileCRUD(t *testing.T) {
	cfg := config.Default()
	st := state.New("test")
	svc := service.New(cfg, st)

	profiles := profile.NewStore(t.TempDir() + "/openhapp")
	srv := New(svc, st, cfg, manifest.FromConfig("test", cfg), profiles)
	d := NewDispatcher(srv)
	ctx := context.Background()

	payload := json.RawMessage(`{
		"id":"de-01",
		"name":"Germany 01",
		"type":"vless",
		"server":"example.com",
		"port":443,
		"enabled":true,
		"properties":{
			"uuid":"test-uuid"
		}
	}`)

	result, err := d.Dispatch(ctx, "profile_add", payload)
	if err != nil {
		t.Fatalf("profile_add: %v", err)
	}

	response, ok := result.(map[string]any)
	if !ok || response["result"] != "ok" {
		t.Fatalf("unexpected profile_add response: %#v", result)
	}

	result, err = d.Dispatch(ctx, "profile_list", nil)
	if err != nil {
		t.Fatalf("profile_list: %v", err)
	}

	list, ok := result.([]profile.Profile)
	if !ok {
		t.Fatalf("unexpected profile_list type: %T", result)
	}

	if len(list) != 1 || list[0].ID != "de-01" {
		t.Fatalf("unexpected profile list: %+v", list)
	}

	result, err = d.Dispatch(
		ctx,
		"profile_get",
		json.RawMessage(`{"id":"de-01"}`),
	)
	if err != nil {
		t.Fatalf("profile_get: %v", err)
	}

	got, ok := result.(profile.Profile)
	if !ok {
		t.Fatalf("unexpected profile_get type: %T", result)
	}

	if got.Name != "Germany 01" {
		t.Fatalf("unexpected profile: %+v", got)
	}

	result, err = d.Dispatch(
		ctx,
		"profile_update",
		json.RawMessage(`{
			"id":"de-01",
			"name":"Germany Updated",
			"type":"vless",
			"server":"example.com",
			"port":8443,
			"enabled":true
		}`),
	)
	if err != nil {
		t.Fatalf("profile_update: %v", err)
	}

	response, ok = result.(map[string]any)
	if !ok || response["result"] != "ok" {
		t.Fatalf("unexpected profile_update response: %#v", result)
	}

	result, err = d.Dispatch(
		ctx,
		"profile_delete",
		json.RawMessage(`{"id":"de-01"}`),
	)
	if err != nil {
		t.Fatalf("profile_delete: %v", err)
	}

	response, ok = result.(map[string]any)
	if !ok || response["result"] != "ok" {
		t.Fatalf("unexpected profile_delete response: %#v", result)
	}

	result, err = d.Dispatch(ctx, "profile_list", nil)
	if err != nil {
		t.Fatalf("profile_list after delete: %v", err)
	}

	list, ok = result.([]profile.Profile)
	if !ok {
		t.Fatalf("unexpected profile list type after delete: %T", result)
	}

	if len(list) != 0 {
		t.Fatalf("expected empty profile list, got: %+v", list)
	}
}
