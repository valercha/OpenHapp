package ubus

import (
	"context"
	"encoding/json"
	"fmt"
)

// Request is the transport-neutral representation of a ubus method call.
type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// Response is the transport-neutral representation returned to the caller.
type Response struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// HandleJSON dispatches a JSON-encoded request through the existing dispatcher.
func (d *Dispatcher) HandleJSON(ctx context.Context, payload []byte) ([]byte, error) {
	var req Request
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, fmt.Errorf("decode ubus request: %w", err)
	}

	result, err := d.Dispatch(ctx, req.Method, req.Params)
	resp := Response{Result: result}
	if err != nil {
		resp.Error = err.Error()
	}

	encoded, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		return nil, fmt.Errorf("encode ubus response: %w", marshalErr)
	}
	return encoded, err
}
