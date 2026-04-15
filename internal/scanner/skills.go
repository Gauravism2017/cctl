package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

func ScanSkills(paths config.Paths) []model.ConfigItem {
	adoptNewItems(paths.Skills, paths.StoreSkills)

	entries, err := os.ReadDir(paths.StoreSkills)
	if err != nil {
		return nil
	}

	var items []model.ConfigItem

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		storePath := filepath.Join(paths.StoreSkills, name)
		activePath := filepath.Join(paths.Skills, name)
		enabled := isSymlink(activePath)

		desc := extractSkillDescription(storePath)

		items = append(items, model.ConfigItem{
			Type:            model.TypeSkill,
			ID:              name,
			Name:            name,
			Description:     desc,
			Enabled:         enabled,
			OriginalEnabled: enabled,
			Path:            storePath,
			Source:          "standalone",
		})
	}

	return items
}

func extractSkillDescription(skillDir string) string {
	skillFile := filepath.Join(skillDir, "SKILL.md")
	f, err := os.Open(skillFile)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	passedFrontmatter := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			inFrontmatter = false
			passedFrontmatter = true
			continue
		}

		if inFrontmatter && strings.HasPrefix(line, "description:") {
			desc := strings.TrimPrefix(line, "description:")
			desc = strings.TrimSpace(desc)
			desc = strings.Trim(desc, "\"'")
			if desc != "" {
				return truncate(desc, 80)
			}
		}

		if passedFrontmatter && !inFrontmatter {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}
			if len(trimmed) > 10 {
				return truncate(trimmed, 80)
			}
		}
	}

	return ""
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
