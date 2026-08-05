# Engine

This package provides the proxy-engine abstraction used by OpenHapp.

It currently exposes a minimal lifecycle façade so the daemon can be wired to Xray first and later extended to sing-box or other engines without changing the higher-level service API.
