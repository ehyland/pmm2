# pmm (Package Manager Manager) v2

A high-performance, dependency-free version of `pmm` written in Go.

## Features

- **Zero Overhead**: Proxies calls to `npm`, `pnpm`, and `yarn` using `syscall.Exec`.
- **Automatic Multi-version Management**: Reads `packageManager` from `package.json` and installs the correct version automatically.
- **Project Pinning**: easily pin a project to a specific package manager version with `pmm pin`.
- **Native Updates**: Self-updates itself directly from GitHub Releases.
- **Cross-platform**: Works on macOS and Linux (AMD64/ARM64).

## Installation

### Shell Script (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/ehyland/pmm2/main/install.sh | bash
```

After running the script, add the binary directory to your `PATH` in your `.bashrc` or `.zshrc`:

```bash
export PATH="$HOME/.pmm2/bin:$PATH"
```

### Homebrew (macOS)

```bash
brew tap ehyland/tap
brew install pmm
```

## Binary Distribution

`pmm` is distributed as a single multi-call binary. When run as `npm`, `pnpm`, or `yarn` (via symlinks), it behaves as a proxy. When run as `pmm`, it provides management commands.

### Commands

- `pmm update-local`: Updates the `packageManager` in the current project to the latest version.
- `pmm update-default [pm]`: Updates the global default version for a package manager.
- `pmm update-self`: Updates `pmm` itself.
- `pmm pin <pm> <path>`: Pins the project at `<path>` to the latest version of `<pm>`.

## License

MIT
