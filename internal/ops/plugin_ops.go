package ops

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

func ApplyPluginChange(item model.ConfigItem, paths config.Paths) error {
	data, err := os.ReadFile(paths.Settings)
	if err != nil {
		return fmt.Errorf("read settings: %w", err)
	}

	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("parse settings: %w", err)
	}

	plugins, ok := settings["enabledPlugins"].(map[string]any)
	if !ok {
		plugins = make(map[string]any)
		settings["enabledPlugins"] = plugins
	}

	plugins[item.ID] = item.Enabled

	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	tmpFile := paths.Settings + ".tmp"
	if err := os.WriteFile(tmpFile, append(out, '\n'), 0o644); err != nil {
		return fmt.Errorf("write temp: %w", err)
	}

	return os.Rename(tmpFile, paths.Settings)
}
