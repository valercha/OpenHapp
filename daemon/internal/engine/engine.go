package engine

import (
	"context"
	"fmt"
	"sync"
)

// Backend describes a concrete proxy-engine implementation.
type Backend interface {
	Name() string
	Version(context.Context) (string, error)
	Check(context.Context) error
	Running() bool
}

// Info describes the selected engine runtime.
type Info struct {
	Name      string `json:"name"`
	Version   string `json:"version,omitempty"`
	Binary    string `json:"binary,omitempty"`
	Config    string `json:"config,omitempty"`
	Workdir   string `json:"workdir,omitempty"`
	Running   bool   `json:"running"`
	Available bool   `json:"available"`
}

// Engine is the runtime facade for a selected proxy backend.
type Engine struct {
	mu      sync.Mutex
	name    string
	running bool
	backend Backend
}

// New creates a new engine facade.
func New(name string) *Engine {
	if name == "" {
		name = "xray"
	}

	return &Engine{
		name:    name,
		backend: backendFor(name),
	}
}

func backendFor(name string) Backend {
	switch name {
	case "sing-box", "singbox":
		return NewSingBoxBackend()
	default:
		return nil
	}
}

// Name returns the configured engine name.
func (e *Engine) Name() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.name
}

// SetName updates the engine name and selects its backend.
func (e *Engine) SetName(name string) {
	if name == "" {
		name = "xray"
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.name = name
	e.backend = backendFor(name)
}

// Start validates the selected backend and marks the engine active.
//
// For the current Sing-box backend this is intentionally non-invasive:
// it checks the existing installation/configuration but does not start,
// stop, or rewrite the existing system Sing-box process.
func (e *Engine) Start(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return nil
	}

	if e.backend != nil {
		if err := e.backend.Check(ctx); err != nil {
			return fmt.Errorf("%s check failed: %w", e.name, err)
		}
	}

	e.running = true
	return nil
}

// Stop marks the engine as inactive.
func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.running = false
}

// Running reports whether the OpenHapp engine facade is active.
func (e *Engine) Running() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.running
}

// Backend returns the currently selected concrete backend.
func (e *Engine) Backend() Backend {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.backend
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

// Info returns details about the selected engine backend.
func (e *Engine) Info(ctx context.Context) Info {
	e.mu.Lock()
	name := e.name
	running := e.running
	backend := e.backend
	e.mu.Unlock()

	info := Info{
		Name:    name,
		Running: running,
	}

	if backend == nil {
		return info
	}

	info.Available = true

	if version, err := backend.Version(ctx); err == nil {
		info.Version = version
	}

	switch b := backend.(type) {
	case *SingBoxBackend:
		info.Binary = b.Binary
		info.Config = b.Config
		info.Workdir = b.Workdir
		info.Available = b.Available()
		info.Running = b.Running()
	}

	return info
}
