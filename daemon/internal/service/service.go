package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

// Service coordinates the lifecycle of the daemon subsystems.
type Service struct {
	mu     sync.Mutex
	cfg    config.Config
	state  *state.State
	running bool
}

// New creates a new service manager.
func New(cfg config.Config, st *state.State) *Service {
	return &Service{cfg: cfg, state: st}
}

// Start starts the service loop.
func (s *Service) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}

	if s.state == nil {
		return fmt.Errorf("service state is nil")
	}

	s.running = true
	s.state.Start()
	s.state.SetMode(s.cfg.Mode)

	go s.loop(ctx)
	return nil
}

// Stop stops the service loop.
func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	if s.state != nil {
		s.state.Stop()
	}
}

func (s *Service) loop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.Stop()
			return
		case <-ticker.C:
			if s.state != nil {
				_ = s.state.Snapshot()
			}
		}
	}
}
