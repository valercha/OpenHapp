# ubus

This package provides the OpenWrt ubus-facing runtime façade for OpenHapp.

It currently exposes a minimal status/config/manifest interface so the LuCI dashboard can query the daemon for:
- runtime state
- runtime configuration
- manifest snapshot

The runtime snapshots are derived from the active service state and manifest model, so the LuCI UI always sees the current daemon view instead of a stale bootstrap-only copy.

Manifest snapshots are refreshed from the active runtime model, which keeps the LuCI dashboard aligned with the daemon’s current engine and routing mode.

Future work will wire these methods to real ubus object registration and method dispatch so OpenHapp can be controlled entirely through OpenWrt-native RPC.
