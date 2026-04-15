package scanner

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

func ScanRules(paths config.Paths) []model.ConfigItem {
	adoptNewItems(paths.Rules, paths.StoreRules)

	entries, err := os.ReadDir(paths.StoreRules)
	if err != nil {
		return nil
	}

	var items []model.ConfigItem

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		storePath := filepath.Join(paths.StoreRules, name)
		activePath := filepath.Join(paths.Rules, name)
		enabled := isSymlink(activePath)

		if entry.IsDir() {
			fileCount := countFiles(storePath)
			desc := name + " rules (" + strconv.Itoa(fileCount) + " files)"

			items = append(items, model.ConfigItem{
				Type:            model.TypeRule,
				ID:              name,
				Name:            name + "/",
				Description:     desc,
				Enabled:         enabled,
				OriginalEnabled: enabled,
				Path:            storePath,
				Category:        name,
			})
		} else if strings.HasSuffix(name, ".md") {
			items = append(items, model.ConfigItem{
				Type:            model.TypeRule,
				ID:              name,
				Name:            name,
				Description:     extractFirstLine(storePath),
				Enabled:         enabled,
				OriginalEnabled: enabled,
				Path:            storePath,
				Category:        "root",
			})
		}
	}

	return items
}

func countFiles(dir string) int {
	count := 0
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			count++
		}
	}
	return count
}

func extractFirstLine(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed == "---" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			trimmed = strings.TrimLeft(trimmed, "# ")
		}
		return truncate(trimmed, 80)
	}

	return ""
}
