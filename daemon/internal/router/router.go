package router

import (
	"sync"

	"github.com/valercha/OpenHapp/daemon/internal/config"
)

// Router represents the future policy-routing layer for OpenHapp.
type Router struct {
	mu  sync.Mutex
	cfg config.Config
}

// New creates a router with the provided runtime configuration.
func New(cfg config.Config) *Router {
	return &Router{cfg: cfg}
}

// UpdateConfig replaces the runtime configuration used by the router.
func (r *Router) UpdateConfig(cfg config.Config) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cfg = cfg
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
