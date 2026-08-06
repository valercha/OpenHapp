# OpenHapp runtime daemon package metadata
#
# This file now also defines the companion LuCI package metadata so the
# packaging layer can evolve in parallel with the runtime daemon.

OPENHAPP_PKG_NAME:=openhapp
OPENHAPP_PKG_VERSION:=0.1.0
OPENHAPP_PKG_RELEASE:=1
OPENHAPP_PKG_TITLE:=OpenHapp Runtime Daemon
OPENHAPP_PKG_DESCRIPTION:=Modern OpenWrt VPN client runtime daemon with ubus and LuCI integration
OPENHAPP_PKG_DEPENDS:=+libubox +libubus +procd

OPENHAPP_LUCI_PKG_NAME:=openhapp-luci
OPENHAPP_LUCI_PKG_VERSION:=0.1.0
OPENHAPP_LUCI_PKG_RELEASE:=1
OPENHAPP_LUCI_PKG_TITLE:=OpenHapp LuCI UI
OPENHAPP_LUCI_PKG_DESCRIPTION:=LuCI web interface for OpenHapp runtime control and status
OPENHAPP_LUCI_PKG_DEPENDS:=+luci-base +rpcd +uhttpd

.PHONY: print
print:
	@echo "$(OPENHAPP_PKG_NAME) $(OPENHAPP_PKG_VERSION)-$(OPENHAPP_PKG_RELEASE)"
	@echo "$(OPENHAPP_LUCI_PKG_NAME) $(OPENHAPP_LUCI_PKG_VERSION)-$(OPENHAPP_LUCI_PKG_RELEASE)"
