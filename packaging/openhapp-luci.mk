# OpenHapp LuCI UI package metadata
#
# Companion package for the LuCI web interface.

OPENHAPP_LUCI_PKG_NAME:=openhapp-luci
OPENHAPP_LUCI_PKG_VERSION:=0.1.0
OPENHAPP_LUCI_PKG_RELEASE:=1
OPENHAPP_LUCI_PKG_TITLE:=OpenHapp LuCI Web Interface
OPENHAPP_LUCI_PKG_DESCRIPTION:=LuCI web interface for the OpenHapp OpenWrt client
OPENHAPP_LUCI_PKG_DEPENDS:=+luci-base +rpcd +uhttpd

.PHONY: print
print:
	@echo "$(OPENHAPP_LUCI_PKG_NAME) $(OPENHAPP_LUCI_PKG_VERSION)-$(OPENHAPP_LUCI_PKG_RELEASE)"
