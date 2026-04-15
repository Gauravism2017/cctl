package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

type mcpSettingsFile struct {
	MCPServers         map[string]json.RawMessage `json:"mcpServers"`
	DisabledMCPServers map[string]json.RawMessage `json:"disabledMcpServers"`
}

func ScanMCPServers(paths config.Paths) []model.ConfigItem {
	var items []model.ConfigItem

	data, err := os.ReadFile(paths.Settings)
	if err != nil {
		return items
	}

	var settings mcpSettingsFile
	if err := json.Unmarshal(data, &settings); err != nil {
		return items
	}

	seen := make(map[string]bool)

	for name, raw := range settings.MCPServers {
		seen[name] = true
		items = append(items, model.ConfigItem{
			Type:            model.TypeMCPServer,
			ID:              name,
			Name:            name,
			Description:     mcpDescription(raw),
			Enabled:         true,
			OriginalEnabled: true,
			Path:            paths.Settings,
		})
	}

	for name, raw := range settings.DisabledMCPServers {
		if seen[name] {
			continue
		}
		items = append(items, model.ConfigItem{
			Type:            model.TypeMCPServer,
			ID:              name,
			Name:            name,
			Description:     mcpDescription(raw),
			Enabled:         false,
			OriginalEnabled: false,
			Path:            paths.Settings,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	return items
}

func mcpDescription(raw json.RawMessage) string {
	var cfg map[string]any
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return ""
	}

	if cmd, ok := cfg["command"].(string); ok {
		return fmt.Sprintf("cmd: %s", cmd)
	}

	if url, ok := cfg["url"].(string); ok {
		serverType, _ := cfg["type"].(string)
		if serverType == "" {
			serverType = "http"
		}
		return fmt.Sprintf("%s: %s", serverType, url)
	}

	return ""
}
