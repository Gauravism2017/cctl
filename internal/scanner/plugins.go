package scanner

import (
	"encoding/json"
	"os"

	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

type settingsFile struct {
	EnabledPlugins map[string]bool `json:"enabledPlugins"`
}

type installedPluginsFile struct {
	Version int                              `json:"version"`
	Plugins map[string][]installedPluginInfo `json:"plugins"`
}

type installedPluginInfo struct {
	Scope       string `json:"scope"`
	ProjectPath string `json:"projectPath,omitempty"`
	InstallPath string `json:"installPath"`
	Version     string `json:"version"`
}

func ScanPlugins(paths config.Paths) []model.ConfigItem {
	var items []model.ConfigItem

	settings := loadSettings(paths.Settings)
	installed := loadInstalled(paths.PluginInstall)

	for pluginKey, enabled := range settings.EnabledPlugins {
		item := model.ConfigItem{
			Type:            model.TypePlugin,
			ID:              pluginKey,
			Name:            pluginKey,
			Enabled:         enabled,
			OriginalEnabled: enabled,
			Path:            paths.Settings,
		}

		if infos, ok := installed.Plugins[pluginKey]; ok && len(infos) > 0 {
			info := infos[0]
			item.Scope = info.Scope
			item.Version = info.Version
			if info.Scope == "local" && info.ProjectPath != "" {
				item.Description = "local: " + info.ProjectPath
			} else {
				item.Description = "scope: " + info.Scope
			}
		}

		items = append(items, item)
	}

	for pluginKey, infos := range installed.Plugins {
		if _, inSettings := settings.EnabledPlugins[pluginKey]; inSettings {
			continue
		}
		if len(infos) == 0 {
			continue
		}
		info := infos[0]
		desc := "scope: " + info.Scope
		if info.Scope == "local" && info.ProjectPath != "" {
			desc = "local: " + info.ProjectPath
		}
		items = append(items, model.ConfigItem{
			Type:            model.TypePlugin,
			ID:              pluginKey,
			Name:            pluginKey,
			Enabled:         false,
			OriginalEnabled: false,
			Path:            paths.Settings,
			Scope:       info.Scope,
			Version:     info.Version,
			Description: desc + " (not in settings)",
		})
	}

	return items
}

func loadSettings(path string) settingsFile {
	s := settingsFile{EnabledPlugins: make(map[string]bool)}
	data, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return s
	}
	if s.EnabledPlugins == nil {
		s.EnabledPlugins = make(map[string]bool)
	}
	return s
}

func loadInstalled(path string) installedPluginsFile {
	f := installedPluginsFile{Plugins: make(map[string][]installedPluginInfo)}
	data, err := os.ReadFile(path)
	if err != nil {
		return f
	}
	if err := json.Unmarshal(data, &f); err != nil {
		return f
	}
	if f.Plugins == nil {
		f.Plugins = make(map[string][]installedPluginInfo)
	}
	return f
}
