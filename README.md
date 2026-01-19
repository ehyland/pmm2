# 📦 PMM - _NodeJS_ Package Manager Manager

Just like corepack, plus

🔥 &nbsp; Better [fnm](https://github.com/Schniz/fnm/issues/566) support.

🔥 &nbsp; `pmm update-local` update `packageManager` for your local project

🔥 &nbsp; `pmm update-default <package-manager> [version]` update the global fallback version

🔥 &nbsp; `pmm pin <package-manager> <path>` add `packageManager` field to your package

🔥 &nbsp; Installs package managers from your configured npm registry (or set `PMM_NPM_REGISTRY` to use an alternative)

> ⏲️ &nbsp; No Windows support at this point (happy to accept a pull request)

## Install

```shell
curl -o- https://raw.githubusercontent.com/ehyland/pmm/main/install.sh | bash
```

## Usage

Add `packageManager` field to your projects `package.json`.

e.g.

```json
{
  "packageManager": "pnpm@7.5.0"
}
```

Then use your package manager as you usually would. Behind the scenes, `pmm` will automatically install and run the package manager version in your `package.json`.

The first time you run `npm` or `pnpm` outside of a configured project / in a global context, pmm will get the latest version of your package manager and set it as the global default. The default can then be updated with `pmm update-default <package-manager> [version]`.

## Uninstall

Simply remove the `~/.pmm` dir and the enabling script in your `~/.bashrc`

```shell
export PMM_DIR="$HOME/.pmm"
[ -s "$PMM_DIR/package/enable.sh" ] && \. "$PMM_DIR/package/enable.sh"  # This loads pmm shims
```

## Contributing to `pmm`

See [CONTRIBUTING.md](./CONTRIBUTING.md)
