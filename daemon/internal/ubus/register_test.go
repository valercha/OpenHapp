package ubus

import (
	"context"
	"testing"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

func TestDispatcherMethods(t *testing.T) {
	cfg := config.Default()
	st := state.New("test")
	svc := service.New(cfg, st)
	srv := New(svc, st, cfg, manifest.FromConfig("test", cfg))
	d := NewDispatcher(srv)
	ctx := context.Background()

	methods := []string{"status", "config", "manifest", "snapshot", "start", "stop"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			if _, err := d.Dispatch(ctx, method); err != nil {
				t.Fatalf("dispatch %s: %v", method, err)
			}
		})
	}
}

func TestDispatcherRejectsUnknownMethod(t *testing.T) {
	cfg := config.Default()
	st := state.New("test")
	svc := service.New(cfg, st)
	d := NewDispatcher(New(svc, st, cfg, manifest.FromConfig("test", cfg)))

	if _, err := d.Dispatch(context.Background(), "unknown"); err == nil {
		t.Fatal("expected unknown method error")
	}
}
