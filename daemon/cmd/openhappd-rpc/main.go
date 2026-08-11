package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/valercha/OpenHapp/daemon/internal/ubus"
)

const socketPath = "/var/run/openhapp.sock"

func main() {
	payload, err := io.ReadAll(os.Stdin)
	if err != nil {
		fatal(err)
	}

	var req ubus.Request
	if err := json.Unmarshal(payload, &req); err != nil {
		fatal(fmt.Errorf("decode request: %w", err))
	}

	conn, err := net.DialTimeout("unix", socketPath, 2*1e9)
	if err != nil {
		fatal(fmt.Errorf("connect to OpenHapp control socket: %w", err))
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		fatal(fmt.Errorf("send request: %w", err))
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadBytes('\n')
	if err != nil {
		fatal(fmt.Errorf("read response: %w", err))
	}

	if _, err := os.Stdout.Write(response); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "openhappd-rpc: %v\n", err)
	os.Exit(1)
}

var _ context.Context
