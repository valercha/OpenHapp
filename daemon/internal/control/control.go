package control

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/ubus"
)

const defaultSocketPath = "/var/run/openhapp.sock"

type Server struct {
	mu         sync.Mutex
	svc        *service.Service
	socketPath string
	listener   net.Listener
}

type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func NewServer(svc *service.Service, socketPath string) *Server {
	if socketPath == "" {
		socketPath = defaultSocketPath
	}
	return &Server{svc: svc, socketPath: socketPath}
}

func (s *Server) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener != nil {
		return nil
	}
	if err := os.Remove(s.socketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove control socket: %w", err)
	}
	ln, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("listen on control socket: %w", err)
	}
	if err := os.Chmod(s.socketPath, 0o660); err != nil {
		_ = ln.Close()
		return fmt.Errorf("chmod control socket: %w", err)
	}
	s.listener = ln
	go s.acceptLoop(ctx, ln)
	return nil
}

func (s *Server) Stop() error {
	s.mu.Lock()
	ln := s.listener
	s.listener = nil
	s.mu.Unlock()
	if ln == nil {
		return nil
	}
	return ln.Close()
}

func (s *Server) acceptLoop(ctx context.Context, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
		go s.handleConn(ctx, conn)
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	var req Request
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		_ = json.NewEncoder(conn).Encode(Response{Error: err.Error()})
		return
	}

	if s.svc == nil {
		_ = json.NewEncoder(conn).Encode(Response{Error: "service is nil"})
		return
	}

	srv := ubus.New(s.svc, nil, s.svc.Config(), ubus.DefaultManifest(s.svc.Config()))
	dispatcher := ubus.NewDispatcher(srv)
	result, err := dispatcher.Dispatch(ctx, req.Method)
	resp := Response{Result: result}
	if err != nil {
		resp.Error = err.Error()
	}
	_ = json.NewEncoder(conn).Encode(resp)
}
