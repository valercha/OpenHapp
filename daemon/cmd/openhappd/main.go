package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
	"github.com/valercha/OpenHapp/daemon/internal/ubus"
	"github.com/valercha/OpenHapp/daemon/internal/version"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("openhappd %s starting", version.String())

	cfg := config.Default()
	st := state.New(version.String())
	st.SetEngine(cfg.Engine)
	st.SetMode(cfg.Mode)
	m := manifest.FromConfig(version.String(), cfg).WithTimestamp()
	svc := service.New(cfg, st)
	bus := ubus.New(svc, st, cfg, m)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			log.Println(version.String())
			return
		case "status":
			log.Printf("status: %+v", bus.Status())
			return
		case "manifest":
			log.Printf("manifest: %s", bus.Manifest().JSON())
			return
		case "config":
			log.Printf("config: %+v", bus.Config())
			return
		case "snapshot":
			log.Printf("snapshot: %+v", bus.Snapshot())
			return
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := bus.Start(ctx); err != nil {
		log.Fatalf("failed to start service: %v", err)
	}

	<-ctx.Done()
	bus.Stop()
	log.Printf("openhappd stopped: %+v", bus.Status())
}
