package ubus

import (
	"context"
	"fmt"
	"sync"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

// Server is a minimal ubus-compatible façade for future RPC wiring.
type Server struct {
	mu       sync.RWMutex
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
	s.mu.RLock()
	svc := s.svc
	s.mu.RUnlock()

	if svc == nil {
		return fmt.Errorf("service is nil")
	}

	return svc.Start(ctx)
}

// Stop stops the underlying service.
func (s *Server) Stop() {
	s.mu.RLock()
	svc := s.svc
	s.mu.RUnlock()

	if svc != nil {
		svc.Stop()
	}
}

// Status returns the current runtime snapshot.
func (s *Server) Status() state.Snapshot {
	s.mu.RLock()
	svc := s.svc
	st := s.st
	s.mu.RUnlock()

	if svc != nil {
		return svc.Snapshot()
	}
	if st == nil {
		return state.Snapshot{}
	}
	return st.Snapshot()
}

// Config returns the current runtime configuration snapshot.
func (s *Server) Config() config.Config {
	s.mu.RLock()
	svc := s.svc
	cfg := s.cfg
	s.mu.RUnlock()

	if svc != nil {
		return svc.Config()
	}
	return cfg
}

// Manifest returns the current runtime manifest snapshot.
func (s *Server) Manifest() manifest.Manifest {
	s.mu.RLock()
	m := s.manifest
	cfg := s.cfg
	svc := s.svc
	s.mu.RUnlock()

	if m.Name == "" || m.Version == "" {
		if svc != nil {
			cfg = svc.Config()
		}
		return manifest.FromConfig("0.1.0-dev", cfg).WithTimestamp()
	}
	return m.WithTimestamp()
}

// StartRPC is a compatibility alias for future ubus dispatch wiring.
func (s *Server) StartRPC(ctx context.Context) error { return s.Start(ctx) }

// StopRPC is a compatibility alias for future ubus dispatch wiring.
func (s *Server) StopRPC() { s.Stop() }

// StatusRPC returns the same runtime snapshot used by the future ubus method.
func (s *Server) StatusRPC() state.Snapshot { return s.Status() }

// ConfigRPC returns the same runtime configuration snapshot used by the future ubus method.
func (s *Server) ConfigRPC() config.Config { return s.Config() }

// ManifestRPC returns the same runtime manifest snapshot used by the future ubus method.
func (s *Server) ManifestRPC() manifest.Manifest { return s.Manifest() }
