package model

type ItemType string

const (
	TypeSkill  ItemType = "skill"
	TypeRule   ItemType = "rule"
	TypeAgent  ItemType = "agent"
	TypePlugin    ItemType = "plugin"
	TypeMCPServer ItemType = "mcp-server"
)

type ConfigItem struct {
	Type            ItemType
	ID              string
	Name            string
	Description     string
	Enabled         bool
	OriginalEnabled bool
	Path            string
	Source          string
	Scope           string
	Category        string
	Version         string
}

func (c *ConfigItem) Toggle() {
	c.Enabled = !c.Enabled
}

func (c ConfigItem) Dirty() bool {
	return c.Enabled != c.OriginalEnabled
}

func (c ConfigItem) StatusIcon() string {
	if c.Enabled {
		return "[x]"
	}
	return "[ ]"
}
