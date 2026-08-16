# Engine

This package provides the proxy-engine abstraction used by OpenHapp.

The engine layer is responsible for starting, stopping, naming, and reporting the runtime state of the selected proxy backend.

Current behavior:
- defaults to `sing-box`
- exposes lifecycle primitives for the daemon
- provides a non-invasive Sing-box backend that validates the existing installation and reports runtime state
- does not currently own or manage the external Sing-box process
- can later support managed engine process supervision without changing the higher-level service API

Planned extensions:
- explicit engine ownership model
- managed process supervision
- config generation
- health checks
- restart and failover hooks
