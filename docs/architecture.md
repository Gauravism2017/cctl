# Architecture

## Symlink Store

cctl uses a symlink store to manage Claude Code items. All items live permanently in `~/.claude/store/`. Enabling an item creates a symlink from the active directory to the store. Disabling removes the symlink — the store copy is never touched.

```
~/.claude/
├── store/              # Permanent home for ALL items
│   ├── skills/         # Skill directories
│   ├── rules/          # Rule directories
│   └── agents/         # Agent files
├── skills/             # Symlinks → store/skills/<name>
├── rules/              # Symlinks → store/rules/<name>
├── agents/             # Symlinks → store/agents/<name>
├── settings.json       # Plugin + MCP server state
└── profiles/           # Named profile snapshots (JSON)
```

## Item Types

| Type | Storage | Toggle Mechanism |
|------|---------|-----------------|
| Skills | `store/skills/` directories | Symlink create/remove |
| Rules | `store/rules/` directories | Symlink create/remove |
| Agents | `store/agents/` files | Symlink create/remove |
| Plugins | `settings.json` | `enabledPlugins` boolean |
| MCP Servers | `settings.json` | Move between `mcpServers` / `disabledMcpServers` |

## Auto-Adoption

Items installed directly by Claude Code (not through cctl) appear as regular files or directories in `~/.claude/skills/`, `~/.claude/rules/`, or `~/.claude/agents/`. On the next scan, cctl detects these non-symlink items, moves them into the store, and replaces them with symlinks. This makes them manageable by cctl without any manual intervention.

## Profile Storage

Profiles are JSON files in `~/.claude/profiles/`:

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

Loading a profile enables everything listed and disables everything else. Items referenced in the profile that no longer exist on disk are silently skipped.

## Per-Project Binding

The `.claude/profile` file in a project directory contains one profile name per line. When `cctl launch` is run:

1. Parse all profile names from the file
2. Filter to profiles that actually exist
3. If one profile remains, auto-select it
4. If multiple remain, show an interactive picker
5. Load the selected profile and exec into Claude Code

## Migration

For users upgrading from the older vault-based system (`skills-vault/`, `rules-vault/`, `agents-vault/`), the `migrate` command moves items from vault directories into `store/` and replaces active items with symlinks. Migration is detected automatically on startup and runs if vault directories are present.
