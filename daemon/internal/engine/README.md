# Engine

This package provides the proxy-engine abstraction used by OpenHapp.

The engine layer is responsible for starting, stopping, naming, and reporting the runtime state of the selected proxy backend.

Current behavior:
- defaults to `xray`
- exposes lifecycle primitives for the daemon
- can later be backed by Xray, sing-box, or another engine without changing the higher-level service API

Planned extensions:
- process supervision
- config generation
- health checks
- restart and failover hooks
