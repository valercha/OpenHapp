# OpenWrt SDK check

The repository contains the package recipe and local staging checks, but a real OpenWrt SDK build requires an OpenWrt build environment/SDK.

Required validation on an OpenWrt SDK matching the target release:

```sh
make defconfig
make package/openhapp/compile V=s
make package/openhapp/install V=s
```

The SDK check is intentionally kept separate from host CI because the host runner does not provide the OpenWrt package build system or target toolchain.

Release readiness requires a successful SDK package build and installation test on the target OpenWrt release.
