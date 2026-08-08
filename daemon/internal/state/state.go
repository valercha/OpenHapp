package state

import "sync"

// State holds the runtime status of openhappd.
type State struct {
	mu      sync.RWMutex
	running bool
	mode    string
	version string
	engine  string
}

// New creates a new state container.
func New(version string) *State {
	return &State{version: version}
}

// Start marks the daemon as running.
func (s *State) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = true
}

// Stop marks the daemon as stopped.
func (s *State) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}

// SetMode updates the current mode.
func (s *State) SetMode(mode string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mode = mode
}

// SetEngine updates the current engine.
func (s *State) SetEngine(engine string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.engine = engine
}

// Snapshot describes the current runtime state.
type Snapshot struct {
	Running bool   `json:"running"`
	Mode    string `json:"mode"`
	Version string `json:"version"`
	Engine  string `json:"engine"`
}

// Snapshot returns the current state values.
func (s *State) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Snapshot{
		Running: s.running,
		Mode:    s.mode,
		Version: s.version,
		Engine:  s.engine,
	}
}
