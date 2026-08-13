#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
STAGE_DIR=$(mktemp -d)
trap 'rm -rf "$STAGE_DIR"' EXIT

OPENHAPP_SOURCE_DIR="$ROOT_DIR/daemon" "$ROOT_DIR/packaging/openwrt/source.sh" "$STAGE_DIR"

for path in \
  go.mod \
  cmd \
  internal \
  openhapp.init \
  openhapp.config \
  99-openhapp \
  rpcd-openhapp \
  openhappd-wrapper
 do
  test -e "$STAGE_DIR/$path"
done

printf '%s\n' 'OpenHapp OpenWrt source staging: OK'
