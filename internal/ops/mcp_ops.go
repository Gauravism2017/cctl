package ops

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

func ApplyMCPServerChange(item model.ConfigItem, paths config.Paths) error {
	data, err := os.ReadFile(paths.Settings)
	if err != nil {
		return fmt.Errorf("read settings: %w", err)
	}

	var settings map[string]json.RawMessage
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("parse settings: %w", err)
	}

	active := extractRawMap(settings, "mcpServers")
	disabled := extractRawMap(settings, "disabledMcpServers")

	if item.Enabled {
		if cfg, ok := disabled[item.ID]; ok {
			active[item.ID] = cfg
			delete(disabled, item.ID)
		}
	} else {
		if cfg, ok := active[item.ID]; ok {
			disabled[item.ID] = cfg
			delete(active, item.ID)
		}
	}

	settings["mcpServers"], _ = json.Marshal(active)
	if len(disabled) > 0 {
		settings["disabledMcpServers"], _ = json.Marshal(disabled)
	} else {
		delete(settings, "disabledMcpServers")
	}

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

func extractRawMap(settings map[string]json.RawMessage, key string) map[string]json.RawMessage {
	raw, ok := settings[key]
	if !ok {
		return make(map[string]json.RawMessage)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return make(map[string]json.RawMessage)
	}
	return m
}
