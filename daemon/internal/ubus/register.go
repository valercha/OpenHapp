package ubus

import (
	"context"
	"fmt"
)

// Dispatcher represents the minimal callable surface for the daemon's ubus layer.
// It is intentionally small so the real OpenWrt ubus wiring can be added later
// without changing the higher-level service/UI contracts.
type Dispatcher struct {
	srv *Server
}

// NewDispatcher creates a dispatch helper bound to a Server instance.
func NewDispatcher(srv *Server) *Dispatcher {
	return &Dispatcher{srv: srv}
}

// Dispatch routes a method name to the corresponding Server RPC-compatible call.
func (d *Dispatcher) Dispatch(ctx context.Context, method string) (any, error) {
	if d == nil || d.srv == nil {
		return nil, fmt.Errorf("ubus dispatcher is nil")
	}

	switch method {
	case "start":
		if err := d.srv.StartRPC(ctx); err != nil {
			return map[string]any{"result": "error", "error": err.Error()}, err
		}
		return map[string]any{"result": "ok"}, nil
	case "stop":
		d.srv.StopRPC()
		return map[string]any{"result": "ok"}, nil
	case "status":
		return d.srv.StatusRPC(), nil
	case "config":
		return d.srv.ConfigRPC(), nil
	case "manifest":
		return d.srv.ManifestRPC(), nil
	case "snapshot":
		return d.srv.SnapshotRPC(), nil
	default:
		return nil, fmt.Errorf("unknown ubus method: %s", method)
	}
}
