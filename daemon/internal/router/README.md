# Router

This package provides the policy-routing layer for OpenHapp.

It currently keeps the runtime configuration and exposes a minimal lifecycle façade so the daemon can later bind routing modes, engine selection, and per-profile policy decisions without changing the higher-level service API.

## Responsibilities

- track the active runtime configuration
- expose the current routing mode
- synchronize engine selection with configuration changes
- provide a lifecycle hook for future policy routing initialization

## Future extensions

- per-profile route rules
- domain and CIDR routing tables
- DNS-aware split routing
- OpenWrt firewall integration
- per-subscription routing profiles
