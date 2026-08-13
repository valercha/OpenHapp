# OpenHapp OpenWrt package

This directory contains the OpenWrt package recipe for the runtime daemon.

The recipe stages the `daemon` Go module through `source.sh`, then builds the two daemon binaries with the OpenWrt Go toolchain. `OPENHAPP_SOURCE_DIR` may be overridden by the integrator; by default it points to the repository's `daemon` directory relative to this package directory.

The package installs the daemon, RPC helper, rpcd plugin, procd init script, UCI configuration, and UCI defaults script.