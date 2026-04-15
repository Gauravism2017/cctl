package ops

import (
	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

func ApplySkillChange(item model.ConfigItem, paths config.Paths) error {
	if item.Enabled {
		return enableItem(paths.StoreSkills, paths.Skills, item.ID)
	}
	return disableItem(paths.Skills, item.ID)
}
