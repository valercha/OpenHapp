package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/control"
	"github.com/valercha/OpenHapp/daemon/internal/engine"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/profile"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
	"github.com/valercha/OpenHapp/daemon/internal/ubus"
	"github.com/valercha/OpenHapp/daemon/internal/uci"
	"github.com/valercha/OpenHapp/daemon/internal/version"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("openhappd %s starting", version.String())

	store := uci.New("/etc/config/openhapp")
	profileStore := profile.NewStore("/etc/config/openhapp")
	cfg, err := store.Load()
	if err != nil {
		log.Printf("failed to load UCI config, falling back to defaults: %v", err)
		cfg = config.Default()
	}

	st := state.New(version.String())
	st.SetEngine(cfg.Engine)
	st.SetMode(cfg.Mode)
	m := manifest.FromConfig(version.String(), cfg).WithTimestamp()
	eng := engine.New(cfg.Engine)
	svc := service.New(cfg, st)
	svc.SetEngine(eng)
	bus := ubus.New(svc, st, cfg, m, profileStore)
	controlServer := control.NewServer(svc, m, "", profileStore)

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
	if err := controlServer.Start(ctx); err != nil {
		bus.Stop()
		log.Fatalf("failed to start control socket: %v", err)
	}

	<-ctx.Done()
	_ = controlServer.Stop()
	bus.Stop()
	log.Printf("openhappd stopped: %+v", bus.Status())
}
