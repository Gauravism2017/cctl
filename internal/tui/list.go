package tui

import (
	"fmt"
	"strings"

	"github.com/Gauravism2017/cctl/internal/model"
)

type itemList struct {
	items    []model.ConfigItem
	filtered []int
	cursor   int
	filter   string
	offset   int
	height   int
}

func newItemList(items []model.ConfigItem, height int) itemList {
	l := itemList{
		items:  items,
		height: height,
	}
	l.applyFilter("")
	return l
}

func (l *itemList) applyFilter(filter string) {
	l.filter = filter
	l.filtered = nil
	lower := strings.ToLower(filter)

	for i, item := range l.items {
		if filter == "" {
			l.filtered = append(l.filtered, i)
			continue
		}
		if strings.Contains(strings.ToLower(item.Name), lower) ||
			strings.Contains(strings.ToLower(item.Description), lower) ||
			strings.Contains(strings.ToLower(item.ID), lower) {
			l.filtered = append(l.filtered, i)
		}
	}

	if l.cursor >= len(l.filtered) {
		l.cursor = max(0, len(l.filtered)-1)
	}
	l.fixOffset()
}

func (l *itemList) moveUp() {
	if l.cursor > 0 {
		l.cursor--
		l.fixOffset()
	}
}

func (l *itemList) moveDown() {
	if l.cursor < len(l.filtered)-1 {
		l.cursor++
		l.fixOffset()
	}
}

func (l *itemList) fixOffset() {
	if l.cursor < l.offset {
		l.offset = l.cursor
	}
	if l.cursor >= l.offset+l.height {
		l.offset = l.cursor - l.height + 1
	}
}

func (l *itemList) toggleCurrent() {
	if len(l.filtered) == 0 {
		return
	}
	idx := l.filtered[l.cursor]
	l.items[idx].Toggle()
}

func (l *itemList) enableAll() {
	for _, idx := range l.filtered {
		l.items[idx].Enabled = true
	}
}

func (l *itemList) disableAll() {
	for _, idx := range l.filtered {
		l.items[idx].Enabled = false
	}
}

func (l *itemList) currentItem() *model.ConfigItem {
	if len(l.filtered) == 0 {
		return nil
	}
	return &l.items[l.filtered[l.cursor]]
}

func (l *itemList) countEnabled() (enabled, total int) {
	total = len(l.items)
	for _, item := range l.items {
		if item.Enabled {
			enabled++
		}
	}
	return
}

func (l *itemList) countDirty() int {
	count := 0
	for _, item := range l.items {
		if item.Dirty() {
			count++
		}
	}
	return count
}

func (l *itemList) render(width int) string {
	var b strings.Builder

	end := min(l.offset+l.height, len(l.filtered))

	for i := l.offset; i < end; i++ {
		idx := l.filtered[i]
		item := l.items[idx]

		isCursor := i == l.cursor

		checkbox := disabledStyle.Render("[ ]")
		if item.Enabled {
			checkbox = enabledStyle.Render("[x]")
		}

		nameStr := item.Name
		if len(nameStr) > 30 {
			nameStr = nameStr[:27] + "..."
		}

		dirtyMark := "  "
		if item.Dirty() {
			dirtyMark = dirtyStyle.Render("* ")
		}

		descStr := item.Description
		maxDesc := width - 40
		if maxDesc < 10 {
			maxDesc = 10
		}
		if len(descStr) > maxDesc {
			descStr = descStr[:maxDesc-3] + "..."
		}

		prefix := "  "
		if isCursor {
			prefix = "▸ "
		}

		if isCursor {
			line := fmt.Sprintf(" %s%s%s %-30s %s", prefix, dirtyMark, checkbox, nameStr, descStr)
			b.WriteString(selectedStyle.Render(line))
		} else {
			line := fmt.Sprintf(" %s%s%s %-30s %s", prefix, dirtyMark, checkbox, dimStyle.Render(nameStr), descStyle.Render(descStr))
			b.WriteString(line)
		}

		if i < end-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
