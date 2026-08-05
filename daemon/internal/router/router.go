package router

import (
	"context"
	"sync"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/engine"
)

// Router represents the future policy-routing layer for OpenHapp.
type Router struct {
	mu      sync.Mutex
	cfg     config.Config
	engine  *engine.Engine
	running bool
}

// New creates a router with the provided runtime configuration.
func New(cfg config.Config, eng *engine.Engine) *Router {
	return &Router{cfg: cfg, engine: eng}
}

// Start prepares the router for runtime use.
func (r *Router) Start(ctx context.Context) error {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return nil
	}

	r.running = true
	if r.engine != nil {
		r.engine.SetName(r.cfg.Engine)
	}
	return nil
}

// Stop stops the router runtime.
func (r *Router) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.running = false
}

// Running reports whether the router runtime is active.
func (r *Router) Running() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.running
}

// UpdateConfig replaces the runtime configuration used by the router.
func (r *Router) UpdateConfig(cfg config.Config) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cfg = cfg
	if r.engine != nil {
		r.engine.SetName(cfg.Engine)
	}
}

// Config returns a copy of the current router configuration.
func (r *Router) Config() config.Config {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.cfg
}

// Mode returns the active routing mode.
func (r *Router) Mode() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cfg.Mode == "" {
		return "proxy"
	}
	return r.cfg.Mode
}
