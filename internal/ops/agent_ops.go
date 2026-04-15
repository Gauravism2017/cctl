package ops

import (
	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
)

func ApplyAgentChange(item model.ConfigItem, paths config.Paths) error {
	if item.Enabled {
		return enableItem(paths.StoreAgents, paths.Agents, item.ID)
	}
	return disableItem(paths.Agents, item.ID)
}
