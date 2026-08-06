# OpenHapp runtime daemon package metadata
#
# This file is the first real packaging artifact for the runtime daemon.
# It is intentionally minimal and will be wired into the OpenWrt package
# build system in the next step.

OPENHAPP_PKG_NAME:=openhapp
OPENHAPP_PKG_VERSION:=0.1.0
OPENHAPP_PKG_RELEASE:=1
OPENHAPP_PKG_TITLE:=OpenHapp Runtime Daemon
OPENHAPP_PKG_DESCRIPTION:=Modern OpenWrt VPN client runtime daemon with ubus and LuCI integration
OPENHAPP_PKG_DEPENDS:=+libubox +libubus +procd

.PHONY: print
print:
	@echo "$(OPENHAPP_PKG_NAME) $(OPENHAPP_PKG_VERSION)-$(OPENHAPP_PKG_RELEASE)"
