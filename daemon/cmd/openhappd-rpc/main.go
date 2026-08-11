package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/valercha/OpenHapp/daemon/internal/ubus"
)

const socketPath = "/var/run/openhapp.sock"

func main() {
	if len(os.Args) < 2 {
		fatal(fmt.Errorf("missing ubus method"))
	}

	params, err := io.ReadAll(os.Stdin)
	if err != nil {
		fatal(err)
	}
	if len(params) == 0 {
		params = []byte("{}")
	}

	if !json.Valid(params) {
		fatal(fmt.Errorf("invalid JSON parameters"))
	}

	req := ubus.Request{Method: os.Args[1], Params: json.RawMessage(params)}

	conn, err := net.DialTimeout("unix", socketPath, 2*time.Second)
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
