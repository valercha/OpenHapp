# LuCI

This directory will contain the OpenHapp LuCI web interface assets.

The first UI milestone will expose a Services → OpenHapp section with a dashboard page that shows:
- runtime status
- daemon version
- current engine
- current routing mode
- runtime manifest snapshot

The UI will be wired to the daemon through ubus and later expanded with server management, subscriptions, diagnostics, and routing controls.
