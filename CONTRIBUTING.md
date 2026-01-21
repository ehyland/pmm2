# Contributing to pmm2

Thank you for your interest in contributing to `pmm2`!

## Development Environment

- **Go**: 1.25 or later.
- **Node.js**: Required only for testing the package managers that `pmm2` manages.

### Building

```bash
go build ./cmd/pmm2
```

### Running Shims

To test the shim behavior, you can create symlinks in a local `bin` directory:

```bash
mkdir -p ./bin
ln -sf $(pwd)/pmm2 ./bin/npm
ln -sf $(pwd)/pmm2 ./bin/pnpm
ln -sf $(pwd)/pmm2 ./bin/yarn
export PATH="$(pwd)/bin:$PATH"
```

## Release Process

`pmm2` uses [GoReleaser](https://goreleaser.com/) and GitHub Actions for automated versioning and distribution.

### 1. Versioning Strategy

We use [Semantic Versioning](https://semver.org/).

### 2. Creating a Release

Releases are triggered by pushing a git tag.

1.  Ensure all changes are committed and pushed to `main`.
2.  Create a new tag (e.g., `v2.0.0`):
    ```bash
    git tag -a v2.0.0 -m "Release v2.0.0"
    ```
3.  Push the tag to GitHub:
    ```bash
    git push origin v2.0.0
    ```

### 3. Automated Workflow

Once a tag matching `v*` is pushed:

- **GitHub Actions** starts the `release.yml` workflow.
- **GoReleaser** runs to:
  - Compile binaries for all supported platforms (macOS/Linux, AMD64/ARM64).
  - Create a GitHub Release with a generated changelog.
  - Upload the `.tar.gz` archives to the release.
  - Update the Homebrew Formula in the `ehyland/homebrew-tap` repository (requires `HOMEBREW_TAP_GITHUB_TOKEN`).

## Deployment

Since `pmm2` uses a native self-update mechanism, users will be able to upgrade their local installations by running:

```bash
pmm update-self
```

This command polls the GitHub Releases API and replaces the binary with the latest available version.
