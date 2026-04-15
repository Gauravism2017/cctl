package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Paths struct {
	ClaudeHome    string
	Skills        string
	Rules         string
	Agents        string
	StoreSkills   string
	StoreRules    string
	StoreAgents   string
	Settings      string
	PluginInstall string
	Profiles      string
}

func NewPaths() (Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Paths{}, fmt.Errorf("get home directory: %w", err)
	}
	claudeHome := filepath.Join(home, ".claude")
	store := filepath.Join(claudeHome, "store")

	return Paths{
		ClaudeHome:    claudeHome,
		Skills:        filepath.Join(claudeHome, "skills"),
		Rules:         filepath.Join(claudeHome, "rules"),
		Agents:        filepath.Join(claudeHome, "agents"),
		StoreSkills:   filepath.Join(store, "skills"),
		StoreRules:    filepath.Join(store, "rules"),
		StoreAgents:   filepath.Join(store, "agents"),
		Settings:      filepath.Join(claudeHome, "settings.json"),
		PluginInstall: filepath.Join(claudeHome, "plugins", "installed_plugins.json"),
		Profiles:      filepath.Join(claudeHome, "profiles"),
	}, nil
}

func (p Paths) EnsureDirs() error {
	dirs := []string{
		p.StoreSkills,
		p.StoreRules,
		p.StoreAgents,
		p.Skills,
		p.Rules,
		p.Agents,
		p.Profiles,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}
	return nil
}

func (p Paths) NeedsMigration() bool {
	_, err := os.Stat(filepath.Join(p.ClaudeHome, "store"))
	return os.IsNotExist(err)
}

func (p Paths) VaultPath(itemType string) string {
	return filepath.Join(p.ClaudeHome, itemType+"-vault")
}
