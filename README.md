# OpenHapp

OpenHapp is a modern, open-source VPN management platform for OpenWrt.

## Current MVP

The repository now contains the core runtime path needed for the first working MVP:

- Go daemon (`openhappd`) with UCI-backed configuration persistence
- long-running Unix control socket for daemon control
- `ubus`-compatible dispatcher and JSON transport
- OpenWrt `rpcd` plugin and RPC helper
- `procd` service integration
- LuCI dashboard and start/stop controls
- LuCI ACL for the implemented `ubus` methods
- OpenWrt package recipe for the daemon and runtime integration files
- automated Go/JSON CI checks

## Architecture

LuCI → rpcd → `openhappd-rpc` → `/var/run/openhapp.sock` → `openhappd` → dispatcher/service/state

UCI remains the persistent configuration source of truth.

## Packaging

The OpenWrt package recipe is under `packaging/openwrt/`. It stages the daemon module, builds `openhappd` and `openhappd-rpc` with the OpenWrt Go toolchain, and installs the runtime integration files.

An actual OpenWrt SDK/build-tree integration test is still required before claiming release-ready status.

## Status

MVP implementation is in progress. Production/release readiness depends on successful OpenWrt SDK build and router integration testing.
