#!/bin/sh
set -eu

# Stage the daemon module and OpenWrt integration files next to this package.
# The package recipe consumes OPENHAPP_SOURCE_DIR instead of assuming a
# hard-coded relationship to TOPDIR.

SOURCE_DIR="${OPENHAPP_SOURCE_DIR:-$(CDPATH= cd -- "$(dirname -- "$0")/../../daemon" && pwd)}"
DEST_DIR="${1:?destination directory is required}"

mkdir -p "$DEST_DIR/cmd" "$DEST_DIR/internal" \
    "$DEST_DIR/etc/init.d" "$DEST_DIR/etc/config" "$DEST_DIR/etc/uci-defaults" \
    "$DEST_DIR/usr/libexec/rpcd" "$DEST_DIR/usr/bin"

cp -R "$SOURCE_DIR/cmd/." "$DEST_DIR/cmd/"
cp -R "$SOURCE_DIR/internal/." "$DEST_DIR/internal/"
cp "$SOURCE_DIR/etc/init.d/openhapp" "$DEST_DIR/openhapp.init"
cp "$SOURCE_DIR/etc/config/openhapp" "$DEST_DIR/openhapp.config"
cp "$SOURCE_DIR/etc/uci-defaults/99-openhapp" "$DEST_DIR/99-openhapp"
cp "$SOURCE_DIR/usr/libexec/rpcd/openhapp" "$DEST_DIR/rpcd-openhapp"
cp "$SOURCE_DIR/usr/bin/openhappd-wrapper" "$DEST_DIR/openhappd-wrapper"
