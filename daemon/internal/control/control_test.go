package control

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/profile"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

func TestServerStatusOverUnixSocket(t *testing.T) {
	socketPath := t.TempDir() + "/openhapp.sock"
	profileStore := profile.NewStore(t.TempDir() + "/openhapp")
	st := state.New("test")
	cfg := config.Default()
	svc := service.New(cfg, st)
	m := manifest.FromConfig("test", cfg)
	server := NewServer(svc, m, socketPath, profileStore)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := server.Start(ctx); err != nil {
		t.Fatalf("start control server: %v", err)
	}

	var conn net.Conn
	var err error

	for i := 0; i < 20; i++ {
		conn, err = net.Dial("unix", socketPath)
		if err == nil {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	if err != nil {
		t.Fatalf("connect control socket: %v", err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(Request{Method: "status"}); err != nil {
		t.Fatalf("encode request: %v", err)
	}

	var response Response
	if err := json.NewDecoder(conn).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Error != "" {
		t.Fatalf("unexpected error: %s", response.Error)
	}

	if err := server.Stop(); err != nil {
		t.Fatalf("stop control server: %v", err)
	}

	if _, err := os.Stat(socketPath); !os.IsNotExist(err) {
		t.Fatalf("control socket still exists after stop: %v", err)
	}
}
