# OpenHapp runtime daemon package metadata
#
# OpenWrt 25.12+ runtime package for OpenHapp.
# This file defines the daemon package inputs and keeps the package split
# aligned with the MVP scope: runtime daemon, LuCI UI, and optional docs.

OPENHAPP_PKG_NAME:=openhapp
OPENHAPP_PKG_VERSION:=0.1.0
OPENHAPP_PKG_RELEASE:=1
OPENHAPP_PKG_TITLE:=OpenHapp Runtime Daemon
OPENHAPP_PKG_DESCRIPTION:=Modern OpenWrt VPN client runtime daemon with ubus and LuCI integration
OPENHAPP_PKG_DEPENDS:=+libubox +libubus +procd +uci
OPENHAPP_PKG_FILES:=/usr/bin/openhappd /usr/bin/openhappd-wrapper /etc/init.d/openhapp /etc/config/openhapp /etc/uci-defaults/99-openhapp

OPENHAPP_LUCI_PKG_NAME:=openhapp-luci
OPENHAPP_LUCI_PKG_VERSION:=0.1.0
OPENHAPP_LUCI_PKG_RELEASE:=1
OPENHAPP_LUCI_PKG_TITLE:=OpenHapp LuCI UI
OPENHAPP_LUCI_PKG_DESCRIPTION:=LuCI web interface for OpenHapp runtime control and status
OPENHAPP_LUCI_PKG_DEPENDS:=+luci-base +rpcd +uhttpd +lua
OPENHAPP_LUCI_PKG_FILES:=/www/luci-static/resources/view/openhapp /usr/share/rpcd/acl.d/openhapp.json /etc/config/openhapp /etc/uci-defaults/99-openhapp

.PHONY: print
print:
	@echo "$(OPENHAPP_PKG_NAME) $(OPENHAPP_PKG_VERSION)-$(OPENHAPP_PKG_RELEASE)"
	@echo "$(OPENHAPP_LUCI_PKG_NAME) $(OPENHAPP_LUCI_PKG_VERSION)-$(OPENHAPP_LUCI_PKG_RELEASE)"
