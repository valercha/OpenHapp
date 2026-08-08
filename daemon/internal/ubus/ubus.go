package ubus

import (
	"context"
	"fmt"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

// Server is a minimal ubus-compatible façade for future RPC wiring.
type Server struct {
	svc      *service.Service
	st       *state.State
	cfg      config.Config
	manifest manifest.Manifest
}

// New creates a new ubus server façade.
func New(svc *service.Service, st *state.State, cfg config.Config, m manifest.Manifest) *Server {
	return &Server{svc: svc, st: st, cfg: cfg, manifest: m}
}

// Start initializes the ubus façade.
func (s *Server) Start(ctx context.Context) error {
	if s.svc == nil {
		return fmt.Errorf("service is nil")
	}

	return s.svc.Start(ctx)
}

// Stop stops the underlying service.
func (s *Server) Stop() {
	if s.svc != nil {
		s.svc.Stop()
	}
}

// Status returns the current runtime snapshot.
func (s *Server) Status() state.Snapshot {
	if s.svc != nil {
		return s.svc.Snapshot()
	}
	if s.st == nil {
		return state.Snapshot{}
	}
	return s.st.Snapshot()
}

// Config returns the current runtime configuration snapshot.
func (s *Server) Config() config.Config {
	if s.svc != nil {
		return s.svc.Config()
	}
	return s.cfg
}

// Manifest returns the current runtime manifest snapshot.
func (s *Server) Manifest() manifest.Manifest {
	if s.manifest.Name == "" || s.manifest.Version == "" {
		cfg := s.Config()
		return manifest.FromConfig("0.1.0-dev", cfg).WithTimestamp()
	}
	return s.manifest.WithTimestamp()
}

// StartRPC is a compatibility alias for future ubus dispatch wiring.
func (s *Server) StartRPC(ctx context.Context) error {
	return s.Start(ctx)
}

// StopRPC is a compatibility alias for future ubus dispatch wiring.
func (s *Server) StopRPC() {
	s.Stop()
}

// StatusRPC returns the same runtime snapshot used by the future ubus method.
func (s *Server) StatusRPC() state.Snapshot {
	return s.Status()
}

// ConfigRPC returns the same runtime configuration snapshot used by the future ubus method.
func (s *Server) ConfigRPC() config.Config {
	return s.Config()
}

// ManifestRPC returns the same runtime manifest snapshot used by the future ubus method.
func (s *Server) ManifestRPC() manifest.Manifest {
	return s.Manifest()
}
