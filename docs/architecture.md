# OpenHapp Architecture

## Overview

OpenHapp is a modern OpenWrt 25.12+ service for managing proxy profiles, routing, diagnostics, and future transport engines.

## Layers

- **LuCI**: JavaScript user interface for dashboard, servers, subscriptions, diagnostics, and settings.
- **ubus**: RPC boundary between the web UI and the daemon.
- **openhappd**: Go daemon that owns configuration, state, service lifecycle, and engine control.
- **Engine adapters**: Xray today; sing-box later.

## Initial scope

Sprint 1 delivers only the foundation:

- package skeleton
- Go module
- daemon skeleton
- basic docs
- repository hygiene

## Future modules

- subscription manager
- routing manager
- DNS manager
- diagnostics
- failover and auto-select
