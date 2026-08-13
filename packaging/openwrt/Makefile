include $(TOPDIR)/rules.mk

PKG_NAME:=openhapp
PKG_VERSION:=0.1.0
PKG_RELEASE:=1

include $(INCLUDE_DIR)/package.mk
include $(TOPDIR)/feeds/packages/lang/golang/golang-package.mk

PKG_BUILD_DEPENDS:=golang/host
PKG_BUILD_PARALLEL:=1
PKG_BUILD_FLAGS:=no-mips16

GO_PKG:=github.com/valercha/OpenHapp/daemon
GO_PKG_BUILD_PKG:=$(GO_PKG)/cmd/openhappd $(GO_PKG)/cmd/openhappd-rpc

PKG_BUILD_DIR:=$(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)
SOURCE_STAGE:=$(PKG_BUILD_DIR)/source
OPENHAPP_SOURCE_DIR?=$(CURDIR)/../../../daemon
OPENHAPP_STAGE_SCRIPT?=$(CURDIR)/source.sh

define Build/Prepare
	mkdir -p $(SOURCE_STAGE)
	OPENHAPP_SOURCE_DIR="$(OPENHAPP_SOURCE_DIR)" $(OPENHAPP_STAGE_SCRIPT) $(SOURCE_STAGE)
endef

define Package/openhapp
  SECTION:=net
  CATEGORY:=Network
  TITLE:=OpenHapp runtime daemon
  DEPENDS:=+libuci +procd +rpcd
endef

define Package/openhapp/description
OpenHapp runtime daemon and OpenWrt rpcd control integration.
endef

define Package/openhapp/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(GO_PKG_BUILD_BIN_DIR)/openhappd $(1)/usr/bin/openhappd
	$(INSTALL_BIN) $(GO_PKG_BUILD_BIN_DIR)/openhappd-rpc $(1)/usr/bin/openhappd-rpc
	$(INSTALL_BIN) $(SOURCE_STAGE)/openhappd-wrapper $(1)/usr/bin/openhappd-wrapper
	$(INSTALL_DIR) $(1)/usr/libexec/rpcd
	$(INSTALL_BIN) $(SOURCE_STAGE)/rpcd-openhapp $(1)/usr/libexec/rpcd/openhapp
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) $(SOURCE_STAGE)/openhapp.init $(1)/etc/init.d/openhapp
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) $(SOURCE_STAGE)/openhapp.config $(1)/etc/config/openhapp
	$(INSTALL_DIR) $(1)/etc/uci-defaults
	$(INSTALL_BIN) $(SOURCE_STAGE)/99-openhapp $(1)/etc/uci-defaults/99-openhapp
endef

$(eval $(call GoBinPackage,openhapp))
$(eval $(call BuildPackage,openhapp))
