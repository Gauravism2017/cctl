# CLI Reference

## Commands

### `cctl`

Launch the interactive TUI. No arguments needed.

```bash
cctl
```

### `cctl launch`

Load a profile and start Claude Code.

```bash
cctl launch                    # Use profile from .claude/profile
cctl launch -p <name>          # Use named profile directly
cctl launch --profile <name>   # Same as -p
```

When `.claude/profile` contains multiple profiles and no `-p` flag is given, an interactive picker is shown.

### `cctl init`

Bind a profile to the current directory. Creates the profile from current state if it doesn't exist.

```bash
cctl init <profile-name>              # Bind profile (additive)
cctl init --remove <profile-name>     # Unbind profile
cctl init -r <profile-name>           # Same as --remove
```

Multiple profiles can be bound to the same directory:

```bash
cctl init frontend
cctl init backend
cctl init data-pipeline
```

### `cctl profile`

Manage saved profiles.

```bash
cctl profile save <name>      # Snapshot current state
cctl profile load <name>      # Restore a profile
cctl profile list             # List all profiles
cctl profile delete <name>    # Delete a profile
```

### `cctl migrate`

Migrate from the old vault-based system to the symlink store.

```bash
cctl migrate
```

Runs automatically on first launch if vault directories are detected.

## Shell Completions

Cobra provides built-in shell completion generation:

```bash
# Zsh (add to ~/.zshrc or generate once)
cctl completion zsh > "${fpath[1]}/_cctl"

# Bash
cctl completion bash > /etc/bash_completion.d/cctl

# Fish
cctl completion fish > ~/.config/fish/completions/cctl.fish
```

## Shell Wrapper

For automatic profile loading when running `claude`, add to `~/.zshrc` or `~/.bashrc`:

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
