# TASK: plugin-system

## Status: Functional

## Issues

### 1. Hard-Coded Cross-Plugin Dependency on plugin-wiz (Fixed)
The direct RPC calls to `plugin-wiz` for cycling RGB colors have been removed. This logic is no longer present in the system sensors plugin.

- [x] Remove `pushWizColor` and its call from `tickLoop`

### 2. Independent NATS Connection (Fixed)
The independent NATS connection and registry watcher have been removed. The plugin no longer maintains its own NATS connection, adhering to the SDK contract where the runner manages the transport.

- [x] Evaluate whether the registry watcher is necessary; if so, find an SDK-compliant way to observe plugin registration (Removed as it was only used for the cross-plugin dependency)

### 3. OnCommand Returns Error for All Commands (By Design - Documented)
`OnCommand` now includes a comment clarifying that system sensors are read-only by design.

- [x] Add a comment clarifying this is by design
