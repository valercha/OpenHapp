package ubus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/valercha/OpenHapp/daemon/internal/profile"
)

// Dispatcher represents the callable surface for the daemon's ubus layer.
type Dispatcher struct {
	srv *Server
}

// NewDispatcher creates a dispatch helper bound to a Server instance.
func NewDispatcher(srv *Server) *Dispatcher {
	return &Dispatcher{srv: srv}
}

// Dispatch routes a method name to the corresponding Server RPC-compatible call.
func (d *Dispatcher) Dispatch(
	ctx context.Context,
	method string,
	params json.RawMessage,
) (any, error) {
	if d == nil || d.srv == nil {
		return nil, fmt.Errorf("ubus dispatcher is nil")
	}

	switch method {
	case "start":
		if err := d.srv.StartRPC(ctx); err != nil {
			return map[string]any{
				"result": "error",
				"error":  err.Error(),
			}, err
		}

		return map[string]any{"result": "ok"}, nil

	case "stop":
		d.srv.StopRPC()
		return map[string]any{"result": "ok"}, nil

	case "status":
		return d.srv.StatusRPC(), nil

	case "engine_info":
		return d.srv.EngineInfo(ctx), nil

	case "config":
		return d.srv.ConfigRPC(), nil

	case "manifest":
		return d.srv.ManifestRPC(), nil

	case "snapshot":
		return d.srv.SnapshotRPC(), nil

	case "profile_list":
		return d.srv.ProfileListRPC()

	case "profile_get":
		var req struct {
			ID string `json:"id"`
		}

		if err := decodeParams(params, &req); err != nil {
			return nil, err
		}

		return d.srv.ProfileGetRPC(req.ID)

	case "profile_add":
		var req profile.Profile

		if err := decodeParams(params, &req); err != nil {
			return nil, err
		}

		if err := d.srv.ProfileAddRPC(req); err != nil {
			return nil, err
		}

		return map[string]any{"result": "ok"}, nil

	case "profile_update":
		var req profile.Profile

		if err := decodeParams(params, &req); err != nil {
			return nil, err
		}

		if err := d.srv.ProfileUpdateRPC(req); err != nil {
			return nil, err
		}

		return map[string]any{"result": "ok"}, nil

	case "profile_delete":
		var req struct {
			ID string `json:"id"`
		}

		if err := decodeParams(params, &req); err != nil {
			return nil, err
		}

		if err := d.srv.ProfileDeleteRPC(req.ID); err != nil {
			return nil, err
		}

		return map[string]any{"result": "ok"}, nil

	default:
		return nil, fmt.Errorf("unknown ubus method: %s", method)
	}
}

func decodeParams(params json.RawMessage, target any) error {
	if len(params) == 0 || string(params) == "null" {
		return nil
	}

	if err := json.Unmarshal(params, target); err != nil {
		return fmt.Errorf("decode ubus params: %w", err)
	}

	return nil
}
