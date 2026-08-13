include $(TOPDIR)/rules.mk

PKG_NAME:=openhapp
PKG_VERSION:=0.1.0
PKG_RELEASE:=1

include $(INCLUDE_DIR)/package.mk
include $(INCLUDE_DIR)/golang.mk

PKG_BUILD_DEPENDS:=golang/host
GO_PKG:=github.com/valercha/OpenHapp/daemon

PKG_BUILD_DIR:=$(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)
SOURCE_STAGE:=$(PKG_BUILD_DIR)/source
OPENHAPP_SOURCE_DIR?=$(TOPDIR)/../OpenHapp/daemon
OPENHAPP_STAGE_SCRIPT?=$(CURDIR)/source.sh

define Build/Prepare
	mkdir -p $(SOURCE_STAGE)
	OPENHAPP_SOURCE_DIR="$(OPENHAPP_SOURCE_DIR)" $(OPENHAPP_STAGE_SCRIPT) $(SOURCE_STAGE)
endef

define Build/Compile
	cd $(SOURCE_STAGE) && $(GO_BIN) build -trimpath -o $(PKG_BUILD_DIR)/openhappd ./cmd/openhappd
	cd $(SOURCE_STAGE) && $(GO_BIN) build -trimpath -o $(PKG_BUILD_DIR)/openhappd-rpc ./cmd/openhappd-rpc
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
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/openhappd $(1)/usr/bin/openhappd
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/openhappd-rpc $(1)/usr/bin/openhappd-rpc
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

$(eval $(call BuildPackage,openhapp))
