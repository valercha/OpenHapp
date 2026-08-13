include $(TOPDIR)/rules.mk

PKG_NAME:=openhapp
PKG_VERSION:=0.1.0
PKG_RELEASE:=1

include $(INCLUDE_DIR)/package.mk
include $(INCLUDE_DIR)/golang.mk

PKG_BUILD_DEPENDS:=golang/host
GO_PKG:=github.com/valercha/OpenHapp/daemon

PKG_BUILD_DIR:=$(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)
SRC_DIR:=$(TOPDIR)/../daemon

# The package is built from the repository checkout. Only the daemon module
# and runtime integration files are staged into the OpenWrt build directory.
define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
	$(CP) $(SRC_DIR)/go.mod $(PKG_BUILD_DIR)/go.mod
	$(CP) $(SRC_DIR)/cmd $(PKG_BUILD_DIR)/cmd
	$(CP) $(SRC_DIR)/internal $(PKG_BUILD_DIR)/internal
	$(CP) $(SRC_DIR)/etc/init.d/openhapp $(PKG_BUILD_DIR)/openhapp.init
	$(CP) $(SRC_DIR)/etc/config/openhapp $(PKG_BUILD_DIR)/openhapp.config
	$(CP) $(SRC_DIR)/etc/uci-defaults/99-openhapp $(PKG_BUILD_DIR)/99-openhapp
	$(CP) $(SRC_DIR)/usr/libexec/rpcd/openhapp $(PKG_BUILD_DIR)/rpcd-openhapp
	$(CP) $(SRC_DIR)/usr/bin/openhappd-wrapper $(PKG_BUILD_DIR)/openhappd-wrapper
endef

define Build/Compile
	cd $(PKG_BUILD_DIR) && $(GO_BIN) build -trimpath -o $(PKG_BUILD_DIR)/openhappd ./cmd/openhappd
	cd $(PKG_BUILD_DIR) && $(GO_BIN) build -trimpath -o $(PKG_BUILD_DIR)/openhappd-rpc ./cmd/openhappd-rpc
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
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/openhappd-wrapper $(1)/usr/bin/openhappd-wrapper
	$(INSTALL_DIR) $(1)/usr/libexec/rpcd
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/rpcd-openhapp $(1)/usr/libexec/rpcd/openhapp
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/openhapp.init $(1)/etc/init.d/openhapp
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) $(PKG_BUILD_DIR)/openhapp.config $(1)/etc/config/openhapp
	$(INSTALL_DIR) $(1)/etc/uci-defaults
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/99-openhapp $(1)/etc/uci-defaults/99-openhapp
endef

$(eval $(call BuildPackage,openhapp))
