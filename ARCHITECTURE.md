# pmm2 Architecture Documentation

## Overview

`pmm2` (Package Manager Manager v2) is a high-performance, dependency-free manager for Node.js package managers (`npm`, `pnpm`, `yarn`). rewritten from TypeScript to Go, it focuses on zero-overhead proxying, automatic version switching, and seamless self-updates.

## Design Philosophy

1.  **Zero Overhead**: Proxying to a package manager should not add perceptible latency.
2.  **Stateless Execution**: The binary determines the correct version to run based on the current context (directory tree) without relying on global state where possible.
3.  **Standalone**: Distributed as a single binary with no external dependencies (like Node.js itself) required for the proxying logic.
4.  **Transparency**: The user should forget `pmm` is even there.

---

## Core Components

### 1. Multi-call Binary Logic
The binary behaves differently based on `os.Args[0]` (the name used to invoke it).
- **Management Mode**: If invoked as `pmm`, it provides CLI commands (`pin`, `update-self`, etc.).
- **Shim Mode**: If invoked as `npm`, `npx`, `pnpm`, `pnpx`, or `yarn`, it enters the proxying logic.

### 2. Execution Flow (Shim Mode)
When a shim is called, `pmm2` follows these steps:
1.  **Discovery**: Climbs the directory tree to find the nearest `package.json`.
2.  **Inspection**: Parses the `packageManager` field (e.g., `pnpm@8.6.0`).
3.  **Resolution**:
    - If `packageManager` is found, use that version.
    - If not found, use the global default version stored in `~/.pmm2/defaults.json`.
    - If no default exists, fetch the latest version from the registry and save it as the new default.
4.  **Installation**:
    - Checks `~/.pmm2/drivers/<name>/<version>` for the package manager.
    - If missing, downloads the tarball from the npm registry, extracts it, and creates a small `bin` entry point if necessary.
5.  **Process Replacement**: Uses `syscall.Exec` to replace the `pmm2` process with the target package manager process (usually `node path/to/pm/bin/pm.js`). This ensures that signals, exit codes, and process ownership are handled natively by the OS with zero overhead.

### 3. Registry & Installer
- **Registry**: Interfaces with the npm registry API to fetch version metadata. Supports custom registries via `PMM_NPM_REGISTRY`.
- **Installer**: Handles idempotent installations. It downloads tarballs, verifies contents, and ensures the target directory is atomic (using temporary directories during extraction).

---

## Technical Stack

- **Language**: Go 1.25+
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Self-Update**: [go-selfupdate](https://github.com/creativeprojects/go-selfupdate)
- **Process Management**: `syscall` (specifically `syscall.Exec` for Unix-based systems).
- **Distribution**: [GoReleaser](https://goreleaser.com/) for multi-arch binaries and Homebrew Formulae.

---

## Directory Structure

```text
pmm2/
├── cmd/
│   └── pmm/
│       └── main.go         # Entry point, CLI definitions, and shim routing.
├── internal/
│   ├── config/             # Config structures and environment variable handling.
│   ├── executor/           # The core "Shim" logic and syscall.Exec implementation.
│   ├── installer/          # Logic for downloading and unpacking PM tarballs.
│   ├── inspector/          # package.json and directory climbing logic.
│   ├── registry/           # NPM registry API client.
│   └── defaults/           # Global version fallback management (~/.pmm2/defaults.json).
├── .goreleaser.yaml         # Build and release automation.
├── install.sh              # Bootstrap script for binary installation.
└── README.md
```

---

## Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PMM_DEBUG` | Enables verbose logging to stderr. | `false` |
| `PMM_NPM_REGISTRY` | Custom npm registry URL. | `https://registry.npmjs.org` |
| `PMM_DIR` | Root directory for storage. | `~/.pmm2` |

---

## Self-Update Mechanism

When `pmm update-self` is run:
1.  The `go-selfupdate` library queries the GitHub Releases API for `ehyland/pmm2`.
2.  It compares the latest tag with the hardcoded `version` (set during build via `ldflags`).
3.  If an update is available, it downloads the compressed binary for the current `GOOS` and `GOARCH`.
4.  It replaces the current executable on disk with the new version.
5.  **Atomic Re-execution**: The running process uses `syscall.Exec` to replace itself with the newly downloaded binary, automatically invoking the `setup` command. This ensures the configuration and shim definitions from the *new* version are used to synchronize symlinks.

A manual synchronization can be triggered using `pmm setup`.

---

## Security & Reliability

- **Atomic Writes**: Installations are downloaded to a temporary folder and moved into place only when successful to prevent half-finished installs.
- **Path Isolation**: Package managers are stored in versioned subdirectories in `~/.pmm2/drivers` to avoid conflicts between different project requirements.
- **Signal Passing**: Because of `syscall.Exec`, signals like `SIGINT` (Ctrl+C) are delivered directly to the underlying package manager without an intermediary Go process.
