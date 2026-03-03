### `plugin-system` repository

#### Project Overview

This repository contains the `plugin-system`, a plugin for the Slidebolt system that provides system-level information as sensors. It's designed to be loaded and managed by the Slidebolt application.

#### Architecture

The `plugin-system` is a Go application that implements the `runner.Plugin` interface from the `slidebolt/sdk-runner`. It integrates with the Slidebolt ecosystem to provide the following features:

-   **Virtual Device**: Creates a virtual device named `system-device`.
-   **System Sensors**: Exposes three sensors as entities of the `system-device`:
    -   `system-time`: The current time.
    -   `system-date`: The current date.
    -   `system-cpu`: The current CPU usage percentage.

The plugin runs a loop that periodically emits events with the latest values for these sensors.

#### Key Files

| File | Description |
| :--- | :--- |
| `go.mod` | Defines the Go module and its dependencies on the `slidebolt/sdk-runner` and `slidebolt/sdk-types`. |
| `main.go` | The main entry point that initializes the `SystemPlugin` and starts it using the `sdk-runner`. |
| `plugin.go` | The core logic of the plugin. It implements the `runner.Plugin` interface, manages the virtual device and its entities, and emits system events. |

#### Available Commands

This plugin is not intended to be run directly by the user. It is a component that is automatically loaded and managed by the Slidebolt system.

#### Standalone Discovery Mode

This plugin supports a standalone discovery mode for rapid testing and diagnostics without requiring the full Slidebolt stack (NATS, Gateway, etc.).

To run discovery and output the results to JSON:
```bash
./plugin-system -discover
```

**Note**: Ensure any required environment variables (e.g., API keys, URLs) are set before running.
