package ubus

import (
	"context"
	"fmt"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

// Server is a minimal ubus-compatible façade for future RPC wiring.
type Server struct {
	svc *service.Service
	st  *state.State
	cfg config.Config
}

// New creates a new ubus server façade.
func New(svc *service.Service, st *state.State, cfg config.Config) *Server {
	return &Server{svc: svc, st: st, cfg: cfg}
}

// Start initializes the ubus façade.
func (s *Server) Start(ctx context.Context) error {
	if s.svc == nil {
		return fmt.Errorf("service is nil")
	}
	return s.svc.Start(ctx)
}

// Stop stops the service.
func (s *Server) Stop() {
	if s.svc == nil {
		return
	}
	s.svc.Stop()
}

// Status returns a snapshot of the current daemon state.
func (s *Server) Status() state.Snapshot {
	if s.st == nil {
		return state.Snapshot{}
	}
	return s.st.Snapshot()
}

// Config returns the runtime configuration snapshot.
func (s *Server) Config() config.Config {
	return s.cfg
}
