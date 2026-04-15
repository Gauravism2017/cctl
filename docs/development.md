# Development

## Build

```bash
make build    # Build to bin/cctl
make install  # Build + copy to ~/.local/bin/
make vet      # Run go vet
make clean    # Remove bin/
```

## Prerequisites

- Go 1.21+

## Project Structure

```
cmd/cctl/
  main.go           Entry point
  root.go           Root command, TUI launcher, shared setup
  launch.go         Launch command with profile resolution
  init.go           Project-profile binding
  migrate.go        Vault-to-store migration command
  profile.go        Profile CRUD subcommands
  picker.go         Interactive profile selector

internal/
  config/           Path constants, directory management
  model/            ConfigItem struct, ItemType enum
  scanner/          Filesystem scanners per item type
  ops/              Enable/disable operations per item type
  tui/              Bubbletea TUI (app, list widget, profile tab, styles)
  profile/          Profile save/load/list/delete logic
  migrate/          Vault → store migration
```

## Dependencies

| Package | Purpose |
|---------|---------|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | TUI framework (Elm architecture) |
| [Lip Gloss](https://github.com/charmbracelet/lipgloss) | Terminal styling |
| [Cobra](https://github.com/spf13/cobra) | CLI framework |

## How the Scanner Works

Each item type has a `Scan*` function in `internal/scanner/` that:

1. Calls `adoptNewItems()` to detect and adopt non-symlink items into the store
2. Reads the store directory (or `settings.json` for plugins/MCP)
3. Checks whether each item is currently enabled (symlink exists or listed in config)
4. Returns a `[]model.ConfigItem` with both `Enabled` and `OriginalEnabled` set

The `OriginalEnabled` field enables dirty tracking — only items where `Enabled != OriginalEnabled` are written on save.

## How Apply Works

`internal/ops/apply.go` iterates all dirty items and dispatches to type-specific handlers:

- **Skills/Rules/Agents**: Create or remove symlinks via `enableItem`/`disableItem` in `ops/common.go`
- **Plugins**: Read `settings.json`, modify `enabledPlugins`, write back
- **MCP Servers**: Read `settings.json`, move config between `mcpServers` and `disabledMcpServers`, write back

## Contributing

1. Fork the repo
2. Create a feature branch
3. Make changes and run `make vet`
4. Ensure `make build` succeeds
5. Open a PR
