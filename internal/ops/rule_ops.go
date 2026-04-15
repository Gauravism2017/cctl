package ops

import (
	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

func ApplyRuleChange(item model.ConfigItem, paths config.Paths) error {
	if item.Enabled {
		return enableItem(paths.StoreRules, paths.Rules, item.ID)
	}
	return disableItem(paths.Rules, item.ID)
}
