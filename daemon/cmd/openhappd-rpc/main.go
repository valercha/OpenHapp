package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
	"github.com/valercha/OpenHapp/daemon/internal/uci"
	"github.com/valercha/OpenHapp/daemon/internal/ubus"
	"github.com/valercha/OpenHapp/daemon/internal/version"
)

func main() {
	payload, err := io.ReadAll(os.Stdin)
	if err != nil {
		fatal(err)
	}

	store := uci.New("/etc/config/openhapp")
	cfg, err := store.Load()
	if err != nil {
		cfg = config.Default()
	}

	st := state.New(version.String())
	st.SetEngine(cfg.Engine)
	st.SetMode(cfg.Mode)
	svc := service.New(cfg, st)
	m := manifest.FromConfig(version.String(), cfg)
	srv := ubus.New(svc, st, cfg, m)
	dispatcher := ubus.NewDispatcher(srv)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Keep start/stop/status operations against one in-process runtime so
	// a single rpcd invocation cannot accidentally leave a child loop running.
	response, dispatchErr := dispatcher.HandleJSON(ctx, payload)
	if dispatchErr != nil && response == nil {
		fatal(dispatchErr)
	}

	if _, err := os.Stdout.Write(response); err != nil {
		fatal(err)
	}
	if err := store.Save(srv.Config()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "openhappd-rpc: persist config: %v\n", err)
	}
	if dispatchErr != nil {
		return
	}
}

func fatal(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "openhappd-rpc: %v\n", err)
	os.Exit(1)
}
