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

// Snapshot holds the daemon runtime data in one payload.
type Snapshot struct {
	Status   state.Snapshot    `json:"status"`
	Config   config.Config     `json:"config"`
	Manifest manifest.Manifest `json:"manifest"`
}

// Server is the OpenHapp ubus-facing runtime façade.
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

// Start starts the service runtime.
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
		snap := svc.Snapshot()
		stateSnap := state.Snapshot{}
		if st != nil {
			stateSnap = st.Snapshot()
		}
		stateSnap.Running = snap.Running
		stateSnap.Mode = snap.Config.Mode
		stateSnap.Engine = snap.Config.Engine
		return stateSnap
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

// Snapshot returns the full runtime payload used by UI and ubus dispatch.
func (s *Server) Snapshot() Snapshot {
	return Snapshot{Status: s.Status(), Config: s.Config(), Manifest: s.Manifest()}
}

// StartRPC is the ubus start method.
func (s *Server) StartRPC(ctx context.Context) error { return s.Start(ctx) }

// StopRPC is the ubus stop method.
func (s *Server) StopRPC() { s.Stop() }

// StatusRPC is the ubus status method.
func (s *Server) StatusRPC() state.Snapshot { return s.Status() }

// ConfigRPC is the ubus config method.
func (s *Server) ConfigRPC() config.Config { return s.Config() }

// ManifestRPC is the ubus manifest method.
func (s *Server) ManifestRPC() manifest.Manifest { return s.Manifest() }

// SnapshotRPC is the ubus snapshot method.
func (s *Server) SnapshotRPC() Snapshot { return s.Snapshot() }
