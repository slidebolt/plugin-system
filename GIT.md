# Git Workflow for plugin-system

This repository contains the Slidebolt System Plugin, which provides internal system entities and monitoring. It produces a standalone binary.

## Dependencies
- **Internal:**
  - `sb-contract`: Core interfaces and shared structures.
  - `sb-domain`: Shared domain models.
  - `sb-messenger-sdk`: Shared messaging interfaces.
  - `sb-runtime`: Core execution environment.
  - `sb-storage-sdk`: Shared storage interfaces.
  - `sb-testkit`: Testing utilities (for BDD and unit tests).
- **External:** 
  - Standard Go library and NATS.
  - `github.com/cucumber/godog`: BDD testing framework.

## Build Process
- **Type:** Go Application (Plugin).
- **Consumption:** Run as a background plugin service.
- **Artifacts:** Produces a binary named `plugin-system`.
- **Command:** `go build -o plugin-system ./cmd/plugin-system`
- **Validation:** 
  - Validated through unit tests: `go test -v ./...`
  - Validated through BDD tests: `go test -v ./cmd/plugin-system`
  - Validated by successful compilation of the binary.

## Pre-requisites & Publishing
As a core plugin, `plugin-system` must be updated whenever the core domain, messaging, storage, or testkit SDKs are changed.

**Before publishing:**
1. Determine current tag: `git tag | sort -V | tail -n 1`
2. Ensure all local tests pass: `go test -v ./...`
3. Ensure the binary builds: `go build -o plugin-system ./cmd/plugin-system`

**Publishing Order:**
1. Ensure all internal dependencies are tagged and pushed.
2. Update `plugin-system/go.mod` to reference the latest tags.
3. Determine next semantic version for `plugin-system` (e.g., `v1.0.4`).
4. Commit and push the changes to `main`.
5. Tag the repository: `git tag v1.0.4`.
6. Push the tag: `git push origin main v1.0.4`.

## Update Workflow & Verification
1. **Modify:** Update system plugin logic in `internal/` or `app/`.
2. **Verify Local:**
   - Run `go mod tidy`.
   - Run `go test ./...`.
   - Run `go test ./cmd/plugin-system` (BDD features).
   - Run `go build -o plugin-system ./cmd/plugin-system`.
3. **Commit:** Ensure the commit message clearly describes the system plugin change.
4. **Tag & Push:** (Follow the Publishing Order above).
