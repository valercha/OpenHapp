# OpenHapp OpenWrt package

This directory contains the OpenWrt package recipe for the runtime daemon.

The recipe expects to be used from an OpenWrt build tree where the repository checkout is available at `$(TOPDIR)/../OpenHapp` (or equivalently adjusted by the integrator). The package stages the `daemon` Go module and the runtime integration files, then builds the two daemon binaries with the OpenWrt Go toolchain.
