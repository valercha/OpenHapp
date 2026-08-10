# Packaging

This directory contains the OpenWrt packaging assets for OpenHapp.

Target platform:
- OpenWrt 25.12+
- apk-based package management
- procd service integration
- LuCI JS UI integration

Planned package layout:
- `openhapp` runtime daemon package
- `openhapp-luci` LuCI web interface package
- `openhapp-docs` optional documentation assets

## Packaging goals

OpenHapp is structured as a small set of installable units so the runtime daemon, the LuCI interface, and any optional assets can evolve independently.

The first release focuses on a minimal but functional runtime package that can:
- install cleanly on OpenWrt 25.12+
- register a procd service
- expose a stable runtime configuration model
- provide a LuCI entry point for status and control
- persist settings through UCI
- expose runtime state through ubus
- ship a wrapper binary path for installation consistency

## Next artifacts

The next files in this directory will define the actual package metadata and build inputs for the daemon and LuCI layers.
