# ubus

This package provides the OpenWrt ubus-facing runtime façade for OpenHapp.

It currently exposes a minimal status/config interface so the LuCI dashboard can query the daemon for:
- runtime state
- runtime configuration
- manifest snapshot

Future work will wire these methods to real ubus object registration and method dispatch so OpenHapp can be controlled entirely through OpenWrt-native RPC.
