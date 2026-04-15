# cctl — Claude Code Context Tool

A TUI and CLI tool for managing [Claude Code](https://docs.anthropic.com/en/docs/claude-code) context. Toggle skills, rules, agents, plugins, and MCP servers on and off. Save named profiles to switch between project contexts instantly.

## Why

Claude Code loads everything in `~/.claude/` — every skill, rule, agent, plugin, and MCP server. On a large setup this burns context window tokens on tools you don't need for the current project. A Go backend project doesn't need frontend design skills, TypeScript reviewers, or Playwright MCP.

cctl lets you curate what's active per project:

- **Web project** — enable frontend skills, TypeScript agents, Playwright MCP
- **Go CLI project** — enable Go patterns, Go reviewers, gopls LSP
- **Research project** — enable Obsidian MCP, Zotero, paper-writing skills

Switch between these setups in one command.

## Install

**Prerequisites:** Go 1.21+

```bash
git clone https://github.com/Gauravism2017/cctl.git
cd cctl
make install
```

This builds the binary and copies it to `~/.local/bin/cctl`. Make sure `~/.local/bin` is in your `PATH`.

### From source (without install)

```bash
make build
./bin/cctl
```

## Quick Start

```bash
# Launch the TUI — browse and toggle items interactively
cctl

# Save your current setup as a profile
cctl profile save my-project

# Bind a profile to a project directory
cd ~/my-project
cctl init my-project

# Launch Claude Code with the bound profile
cctl launch
```

## How It Works

cctl uses a **symlink store** to manage items:

```
~/.claude/
├── store/              # Permanent home for ALL items
│   ├── skills/
│   ├── rules/
│   └── agents/
├── skills/             # Symlinks → store/skills/<name>
├── rules/              # Symlinks → store/rules/<name>
├── agents/             # Symlinks → store/agents/<name>
├── settings.json       # Plugin + MCP server state
└── profiles/           # Named profile snapshots (JSON)
```

| Action | What happens |
|--------|-------------|
| **Enable** skill/rule/agent | Create symlink from active dir → store |
| **Disable** skill/rule/agent | Remove symlink (store copy untouched) |
| **Enable** plugin | Add to `enabledPlugins` in `settings.json` |
| **Disable** plugin | Remove from `enabledPlugins` |
| **Enable** MCP server | Move config to `mcpServers` in `settings.json` |
| **Disable** MCP server | Move config to `disabledMcpServers` |

Items installed by Claude Code are automatically adopted into the store on next scan.

## CLI Reference

```
cctl                                # Launch TUI
cctl launch                         # Load bound profile + start Claude Code
cctl launch -p <name>               # Load named profile + start Claude Code
cctl init <profile>                 # Bind profile to current directory
cctl init --remove <profile>        # Unbind profile from current directory
cctl migrate                        # Migrate from vault to symlink store
cctl profile save <name>            # Save current state as profile
cctl profile load <name>            # Load a profile
cctl profile list                   # List all profiles
cctl profile delete <name>          # Delete a profile
```

### Shell completions

```bash
# Zsh
cctl completion zsh > "${fpath[1]}/_cctl"

# Bash
cctl completion bash > /etc/bash_completion.d/cctl

# Fish
cctl completion fish > ~/.config/fish/completions/cctl.fish
```

## TUI

Launch with `cctl` (no arguments).

The TUI has tabs for each item type: **Skills**, **Rules**, **Agents**, **Plugins**, **MCP Servers**, **Profiles**, and a **Summary** of pending changes.

### Keybindings

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Switch tabs |
| `j` / `k` | Move cursor |
| `Space` | Toggle item on/off |
| `/` | Filter items |
| `c` | Clear filter |
| `a` | Enable all visible |
| `n` | Disable all visible |
| `s` | Save changes and exit |
| `q` | Quit without saving |

### Profiles tab

| Key | Action |
|-----|--------|
| `Enter` / `l` | View profile details |
| `Space` | Load profile |
| `e` | Edit profile |
| `d` | Delete profile |
| `n` | Create new profile |

In edit mode: `Enter`/`Space` to expand sections, `Space` on items to toggle, `Tab` to jump sections, `a`/`x` to enable/disable all in section, `s` to save.

## Profiles

Profiles snapshot the enabled/disabled state of all items. Save a profile, switch context, then load it to restore your setup.

### Save and load

```bash
# Snapshot current state
cctl profile save web-dev

# Restore it later
cctl profile load web-dev
```

Loading a profile enables everything in the profile and disables everything else. Items in the profile that no longer exist on disk are skipped.

### Per-project profiles

Bind one or more profiles to a directory:

```bash
cd ~/my-web-project
cctl init web-dev
cctl launch              # loads web-dev profile, starts Claude Code
```

### Multi-profile directories

Monorepos or multi-context directories can bind multiple profiles:

```bash
cd ~/monorepo
cctl init frontend
cctl init backend
cctl init data-pipeline
```

When you run `cctl launch`, an interactive picker appears:

```
Select profile:
  ▸ frontend
    backend
    data-pipeline
```

Use `j`/`k` to navigate, `Enter` to select, `q` to cancel.

To skip the picker, pass the profile directly:

```bash
cctl launch -p backend
```

To remove a binding:

```bash
cctl init --remove data-pipeline
```

### Automatic profile loading

Add a shell wrapper to `~/.zshrc` (or `~/.bashrc`) for automatic profile loading when you run `claude`:

```bash
claude() {
  if [[ -f ".claude/profile" ]]; then
    local profiles
    profiles=($(cat .claude/profile))
    if [[ ${#profiles[@]} -eq 1 ]]; then
      cctl profile load "${profiles[1]}" 2>/dev/null
    fi
  fi
  command claude "$@"
}
```

### Profile storage format

Profiles are stored as JSON in `~/.claude/profiles/`:

```json
{
  "name": "web-dev",
  "created": "2025-04-15T10:30:00Z",
  "skills": ["frontend-design", "api-design"],
  "rules": ["common", "web"],
  "agents": ["code-reviewer.md", "architect.md"],
  "plugins": ["superpowers@claude-plugins-official"],
  "mcpServers": ["context7", "playwright"]
}
```

## Migration

If upgrading from an older vault-based setup (`skills-vault/`, `rules-vault/`, `agents-vault/` directories):

```bash
cctl migrate
```

This moves items from vault directories into `store/` and replaces active items with symlinks. Migration runs automatically on first launch if vault directories are detected.

## Development

### Build

```bash
make build    # Build to bin/cctl
make install  # Build + copy to ~/.local/bin/
make vet      # Run go vet
make clean    # Remove bin/
```

### Project structure

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
  scanner/          Filesystem scanners (skills, rules, agents, plugins, MCP)
  ops/              Enable/disable operations per item type
  tui/              Bubbletea TUI (app, list widget, profile tab, styles)
  profile/          Profile save/load/list/delete logic
  migrate/          Vault → store migration
```

### Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Cobra](https://github.com/spf13/cobra) — CLI framework

## License

MIT
