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

OpenHapp is being structured as a small set of installable units so the runtime daemon, the LuCI interface, and any optional assets can evolve independently.

The first release will focus on a minimal but functional runtime package that can:
- install cleanly on OpenWrt 25.12+
- register a procd service
- expose a stable runtime configuration model
- provide a LuCI entry point for status and control

## OpenWrt package plan

The packaging layer is split into a small number of artifacts that are meant to be installed together or separately:

- `openhapp`: runtime daemon, init script, UCI config, and runtime helpers
- `openhapp-luci`: LuCI dashboard and actions pages
- `openhapp-docs`: optional documentation files

The runtime package will own the daemon process, its procd service, default UCI config, and the ubus-facing runtime API surface used by the UI.

The LuCI package will own the menu entries, ACL rules, and views that call into the daemon's runtime API.

## Runtime package responsibilities

The `openhapp` package will include:
- `/usr/bin/openhappd`
- `/etc/init.d/openhapp`
- `/etc/config/openhapp`
- runtime files needed by the daemon

## LuCI package responsibilities

The `openhapp-luci` package will include:
- `menu.d` wiring
- ACL rules
- dashboard and actions views
- JavaScript resources for status/control pages

## Next artifacts

The next files in this directory will define the actual package metadata and build inputs for the daemon and LuCI layers.
