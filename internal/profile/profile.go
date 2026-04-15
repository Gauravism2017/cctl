package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
	"github.com/Gauravism2017/cctl/internal/ops"
	"github.com/Gauravism2017/cctl/internal/scanner"
)

type Profile struct {
	Name       string   `json:"name"`
	Created    string   `json:"created"`
	Skills     []string `json:"skills"`
	Rules      []string `json:"rules"`
	Agents     []string `json:"agents"`
	Plugins    []string `json:"plugins"`
	MCPServers []string `json:"mcpServers,omitempty"`
}

type Entry struct {
	Name    string
	Created string
}

func Save(name string, paths config.Paths) error {
	if err := validateName(name); err != nil {
		return err
	}

	skills := scanner.ScanSkills(paths)
	rules := scanner.ScanRules(paths)
	agents := scanner.ScanAgents(paths)
	plugins := scanner.ScanPlugins(paths)
	mcpServers := scanner.ScanMCPServers(paths)

	p := Profile{
		Name:       name,
		Created:    time.Now().Format(time.RFC3339),
		Skills:     collectEnabled(skills),
		Rules:      collectEnabled(rules),
		Agents:     collectEnabled(agents),
		Plugins:    collectEnabled(plugins),
		MCPServers: collectEnabled(mcpServers),
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profile: %w", err)
	}

	filePath := filepath.Join(paths.Profiles, name+".json")
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("write profile: %w", err)
	}

	return nil
}

func Load(name string, paths config.Paths) (applied int, skipped int, errs []error) {
	p, err := Get(name, paths)
	if err != nil {
		return 0, 0, []error{err}
	}

	skills := scanner.ScanSkills(paths)
	rules := scanner.ScanRules(paths)
	agents := scanner.ScanAgents(paths)
	plugins := scanner.ScanPlugins(paths)
	mcpServers := scanner.ScanMCPServers(paths)

	enabledSet := map[model.ItemType]map[string]bool{
		model.TypeSkill:     toSet(p.Skills),
		model.TypeRule:      toSet(p.Rules),
		model.TypeAgent:     toSet(p.Agents),
		model.TypePlugin:    toSet(p.Plugins),
		model.TypeMCPServer: toSet(p.MCPServers),
	}

	var dirty []model.ConfigItem
	allItems := [][]model.ConfigItem{skills, rules, agents, plugins, mcpServers}
	types := []model.ItemType{model.TypeSkill, model.TypeRule, model.TypeAgent, model.TypePlugin, model.TypeMCPServer}

	for ti, items := range allItems {
		set := enabledSet[types[ti]]
		for i := range items {
			shouldEnable := set[items[i].ID]
			if items[i].Enabled != shouldEnable {
				items[i].Enabled = shouldEnable
				dirty = append(dirty, items[i])
			}
		}
	}

	countSkipped := 0
	for _, set := range enabledSet {
		for id := range set {
			if !existsInItems(id, skills, rules, agents, plugins, mcpServers) {
				countSkipped++
			}
		}
	}

	applyErrs := ops.ApplyChanges(dirty, paths)

	return len(dirty), countSkipped, applyErrs
}

func List(paths config.Paths) ([]Entry, error) {
	entries, err := os.ReadDir(paths.Profiles)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read profiles dir: %w", err)
	}

	var result []Entry
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".json")
		p, err := Get(name, paths)
		if err != nil {
			continue
		}

		result = append(result, Entry{
			Name:    p.Name,
			Created: p.Created,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func Delete(name string, paths config.Paths) error {
	if err := validateName(name); err != nil {
		return err
	}

	filePath := filepath.Join(paths.Profiles, name+".json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("profile %q not found", name)
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}

	return nil
}

func Get(name string, paths config.Paths) (Profile, error) {
	if err := validateName(name); err != nil {
		return Profile{}, err
	}

	filePath := filepath.Join(paths.Profiles, name+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Profile{}, fmt.Errorf("profile %q not found", name)
		}
		return Profile{}, fmt.Errorf("read profile: %w", err)
	}

	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return Profile{}, fmt.Errorf("parse profile: %w", err)
	}

	return p, nil
}

func Update(p Profile, paths config.Paths) error {
	if err := validateName(p.Name); err != nil {
		return err
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profile: %w", err)
	}

	filePath := filepath.Join(paths.Profiles, p.Name+".json")
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("write profile: %w", err)
	}

	return nil
}

func ScanAll(paths config.Paths) (skills, rules, agents, plugins, mcpServers []model.ConfigItem) {
	skills = scanner.ScanSkills(paths)
	rules = scanner.ScanRules(paths)
	agents = scanner.ScanAgents(paths)
	plugins = scanner.ScanPlugins(paths)
	mcpServers = scanner.ScanMCPServers(paths)
	return
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}
	if strings.Contains(name, string(filepath.Separator)) || strings.Contains(name, "..") {
		return fmt.Errorf("invalid profile name: %s", name)
	}
	return nil
}

func collectEnabled(items []model.ConfigItem) []string {
	var ids []string
	for _, item := range items {
		if item.Enabled {
			ids = append(ids, item.ID)
		}
	}
	sort.Strings(ids)
	return ids
}

func toSet(ids []string) map[string]bool {
	m := make(map[string]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}

func existsInItems(id string, itemSets ...[]model.ConfigItem) bool {
	for _, items := range itemSets {
		for _, item := range items {
			if item.ID == id {
				return true
			}
		}
	}
	return false
}
