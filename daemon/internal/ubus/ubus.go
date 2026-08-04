package ubus

import (
	"context"
	"fmt"

	"github.com/valercha/OpenHapp/daemon/internal/service"
)

// Server is a minimal ubus-compatible façade for future RPC wiring.
type Server struct {
	svc *service.Service
}

// New creates a new ubus server façade.
func New(svc *service.Service) *Server {
	return &Server{svc: svc}
}

// Start initializes the ubus façade.
func (s *Server) Start(ctx context.Context) error {
	if s.svc == nil {
		return fmt.Errorf("service is nil")
	}

	return s.svc.Start(ctx)
}
