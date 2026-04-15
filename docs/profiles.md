# Profiles

Profiles snapshot the enabled/disabled state of all Claude Code context items. Save a profile, switch projects, then load it to restore your setup.

## Save and Load

```bash
# Snapshot current state
cctl profile save web-dev

# Restore later
cctl profile load web-dev

# See what's available
cctl profile list

# Remove one
cctl profile delete web-dev
```

Loading a profile enables everything listed in it and disables everything else. Items that no longer exist on disk are skipped.

## Per-Project Profiles

Bind a profile to a directory so `cctl launch` picks it up automatically:

```bash
cd ~/my-web-project
cctl init web-dev
cctl launch              # loads web-dev, starts Claude Code
```

If the profile doesn't exist yet, `cctl init` creates it from the current state.

## Multi-Profile Directories

Monorepos or multi-context directories can have multiple profiles:

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

Navigate with `j`/`k`, select with `Enter`, cancel with `q`.

Skip the picker with the `-p` flag:

```bash
cctl launch -p backend
```

Remove a binding:

```bash
cctl init --remove data-pipeline
```

## Example Profiles

### Web Frontend
```bash
cctl profile save web-dev
```
Skills: `frontend-design`, `frontend-patterns`, `tdd-workflow`
Agents: `typescript-reviewer.md`, `a11y-architect.md`, `e2e-runner.md`
Plugins: `typescript-lsp`, `frontend-design`
MCP: `context7`, `playwright`

### Go Backend
```bash
cctl profile save go-backend
```
Skills: `golang-patterns`, `golang-testing`, `tdd-workflow`
Agents: `go-reviewer.md`, `go-build-resolver.md`, `architect.md`
Plugins: `gopls-lsp`, `code-review`
MCP: `context7`, `serena`

### Research / Writing
```bash
cctl profile save research
```
Skills: `deep-research`, `search-first`
Agents: `docs-lookup.md`, `planner.md`
MCP: `context7`, `obsidian`, `zotero`
