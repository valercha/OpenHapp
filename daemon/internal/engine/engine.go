package engine

import (
	"context"
	"fmt"
	"sync"
)

// Engine models the future proxy engine backend.
type Engine struct {
	mu      sync.Mutex
	name    string
	running bool
}

// New creates a new engine facade.
func New(name string) *Engine {
	if name == "" {
		name = "xray"
	}
	return &Engine{name: name}
}

// Name returns the configured engine name.
func (e *Engine) Name() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.name
}

// Start marks the engine as running.
func (e *Engine) Start(ctx context.Context) error {
	_ = ctx
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.running {
		return nil
	}
	e.running = true
	return nil
}

// Stop marks the engine as stopped.
func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.running = false
}

// Status returns a compact engine status string.
func (e *Engine) Status() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	state := "stopped"
	if e.running {
		state = "running"
	}
	return fmt.Sprintf("engine=%s state=%s", e.name, state)
}
