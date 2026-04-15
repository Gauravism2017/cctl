# TUI Guide

Launch with `cctl` (no arguments). The TUI has tabs for each item type plus profiles and a summary view.

## Tabs

| Tab | Contents |
|-----|----------|
| Skills | Skill directories from the store |
| Rules | Rule directories from the store |
| Agents | Agent files from the store |
| Plugins | Plugins from settings.json |
| MCP Servers | MCP server configs from settings.json |
| Profiles | Saved profiles with view/edit/load/delete |
| Summary | Pending changes before saving |

## Navigation

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Switch tabs |
| `j` / `k` | Move cursor up/down |
| `Space` | Toggle item on/off |
| `/` | Start filtering |
| `c` | Clear filter |
| `a` | Enable all visible items |
| `n` | Disable all visible items |
| `s` | Save changes and exit |
| `q` | Quit without saving |

## Profiles Tab

### List View

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate profiles |
| `Enter` / `l` | View profile details |
| `Space` | Load profile immediately |
| `e` | Edit profile (toggle items on/off) |
| `d` | Delete selected profile |
| `n` | Create new empty profile |

### Detail View

Shows all skills, rules, agents, plugins, and MCP servers in the selected profile.

| Key | Action |
|-----|--------|
| `j` / `k` | Scroll through contents |
| `e` | Edit this profile |
| `Space` | Load this profile |
| `Esc` | Go back to list |

### Edit View

Accordion-style sections for each item type.

| Key | Action |
|-----|--------|
| `Enter` / `Space` on header | Expand/collapse section |
| `Space` on item | Toggle on/off |
| `Tab` | Jump to next section |
| `a` | Enable all in current section |
| `x` | Disable all in current section |
| `s` | Save changes |
| `Esc` | Cancel |
