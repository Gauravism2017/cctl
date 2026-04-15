package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

func ScanAgents(paths config.Paths) []model.ConfigItem {
	adoptNewItems(paths.Agents, paths.StoreAgents)

	entries, err := os.ReadDir(paths.StoreAgents)
	if err != nil {
		return nil
	}

	var items []model.ConfigItem

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}

		storePath := filepath.Join(paths.StoreAgents, name)
		activePath := filepath.Join(paths.Agents, name)
		enabled := isSymlink(activePath)

		agentName := strings.TrimSuffix(name, ".md")
		desc := extractFirstLine(storePath)

		items = append(items, model.ConfigItem{
			Type:            model.TypeAgent,
			ID:              name,
			Name:            agentName,
			Description:     desc,
			Enabled:         enabled,
			OriginalEnabled: enabled,
			Path:            storePath,
		})
	}

	return items
}
