include $(TOPDIR)/rules.mk

PKG_NAME:=openhapp
PKG_VERSION:=0.1.0
PKG_RELEASE:=1
PKG_BUILD_DIR:=$(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)

include $(INCLUDE_DIR)/package.mk

define Package/openhapp
  SECTION:=net
  CATEGORY:=Network
  TITLE:=OpenHapp runtime daemon
  DEPENDS:=+libubus +libuci +procd +rpcd
endef

define Package/openhapp/description
OpenHapp runtime daemon and OpenWrt rpcd control integration.
endef

define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
	$(CP) $(TOPDIR)/../daemon $(PKG_BUILD_DIR)/daemon
	$(CP) $(TOPDIR)/../daemon/etc/init.d/openhapp $(PKG_BUILD_DIR)/openhapp.init
	$(CP) $(TOPDIR)/../daemon/etc/config/openhapp $(PKG_BUILD_DIR)/openhapp.config
	$(CP) $(TOPDIR)/../daemon/etc/uci-defaults/99-openhapp $(PKG_BUILD_DIR)/99-openhapp
	$(CP) $(TOPDIR)/../daemon/usr/libexec/rpcd/openhapp $(PKG_BUILD_DIR)/rpcd-openhapp
	$(CP) $(TOPDIR)/../daemon/usr/bin/openhappd-wrapper $(PKG_BUILD_DIR)/openhappd-wrapper
endef

define Build/Compile
	cd $(PKG_BUILD_DIR)/daemon && \
		$(GO_BIN) build -trimpath -o $(PKG_BUILD_DIR)/openhappd ./cmd/openhappd && \
		$(GO_BIN) build -trimpath -o $(PKG_BUILD_DIR)/openhappd-rpc ./cmd/openhappd-rpc
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
