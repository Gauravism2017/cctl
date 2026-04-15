package ops

import (
	"fmt"

	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

func ApplyChanges(items []model.ConfigItem, paths config.Paths) []error {
	var errs []error

	if err := paths.EnsureDirs(); err != nil {
		return []error{fmt.Errorf("ensure dirs: %w", err)}
	}

	for _, item := range items {
		if !item.Dirty() {
			continue
		}

		var err error
		switch item.Type {
		case model.TypeSkill:
			err = ApplySkillChange(item, paths)
		case model.TypeRule:
			err = ApplyRuleChange(item, paths)
		case model.TypeAgent:
			err = ApplyAgentChange(item, paths)
		case model.TypePlugin:
			err = ApplyPluginChange(item, paths)
		case model.TypeMCPServer:
			err = ApplyMCPServerChange(item, paths)
		}

		if err != nil {
			errs = append(errs, fmt.Errorf("%s %s: %w", item.Type, item.ID, err))
		}
	}

	return errs
}
