package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
	"github.com/Gauravism2017/cctl/internal/profile"
)

type profileMode int

const (
	profileModeList profileMode = 0
	profileModeSave profileMode = 1
	profileModeView profileMode = 2
	profileModeEdit profileMode = 3
)

type profileLoadedMsg struct{}

type editRow struct {
	id         string
	name       string
	itemType   model.ItemType
	enabled    bool
	isHeader   bool
	sectionIdx int
}

type profileTab struct {
	profiles    []profile.Entry
	cursor      int
	offset      int
	height      int
	mode        profileMode
	nameInput   string
	message     string
	paths       config.Paths
	viewProfile profile.Profile
	viewScroll  int
	editProfile     profile.Profile
	editRows        []editRow
	editCursor      int
	editOffset      int
	expandedSection int
}

func newProfileTab(paths config.Paths, height int) profileTab {
	entries, _ := profile.List(paths)
	return profileTab{
		profiles: entries,
		height:   height,
		paths:    paths,
	}
}

func (t *profileTab) refresh() {
	entries, _ := profile.List(t.paths)
	t.profiles = entries
	if t.cursor >= len(t.profiles) {
		t.cursor = max(0, len(t.profiles)-1)
	}
	t.fixOffset()
}

func (t *profileTab) handleKey(msg tea.KeyMsg) (tea.Cmd, string) {
	switch t.mode {
	case profileModeSave:
		return t.handleSaveKey(msg)
	case profileModeView:
		return t.handleViewKey(msg)
	case profileModeEdit:
		return t.handleEditKey(msg)
	default:
		return t.handleListKey(msg)
	}
}

func (t *profileTab) handleListKey(msg tea.KeyMsg) (tea.Cmd, string) {
	switch msg.String() {
	case "j", "down":
		if t.cursor < len(t.profiles)-1 {
			t.cursor++
			t.fixOffset()
		}
		return nil, ""

	case "k", "up":
		if t.cursor > 0 {
			t.cursor--
			t.fixOffset()
		}
		return nil, ""

	case "enter", "l":
		if len(t.profiles) == 0 {
			return nil, ""
		}
		entry := t.profiles[t.cursor]
		p, err := profile.Get(entry.Name, t.paths)
		if err != nil {
			return nil, "Error: " + err.Error()
		}
		t.viewProfile = p
		t.viewScroll = 0
		t.mode = profileModeView
		return nil, ""

	case " ":
		if len(t.profiles) == 0 {
			return nil, ""
		}
		entry := t.profiles[t.cursor]
		applied, skipped, errs := profile.Load(entry.Name, t.paths)
		if len(errs) > 0 {
			var msgs []string
			for _, err := range errs {
				msgs = append(msgs, err.Error())
			}
			return nil, "Load errors: " + strings.Join(msgs, "; ")
		}
		result := fmt.Sprintf("Loaded %q: %d changes", entry.Name, applied)
		if skipped > 0 {
			result += fmt.Sprintf(", %d skipped", skipped)
		}
		return func() tea.Msg { return profileLoadedMsg{} }, result

	case "e":
		if len(t.profiles) == 0 {
			return nil, ""
		}
		entry := t.profiles[t.cursor]
		p, err := profile.Get(entry.Name, t.paths)
		if err != nil {
			return nil, "Error: " + err.Error()
		}
		t.enterEditMode(p)
		return nil, ""

	case "d":
		if len(t.profiles) == 0 {
			return nil, ""
		}
		entry := t.profiles[t.cursor]
		if err := profile.Delete(entry.Name, t.paths); err != nil {
			return nil, "Delete error: " + err.Error()
		}
		t.refresh()
		return nil, fmt.Sprintf("Deleted %q", entry.Name)

	case "n":
		t.mode = profileModeSave
		t.nameInput = ""
		return nil, ""
	}

	return nil, ""
}

func (t *profileTab) handleViewKey(msg tea.KeyMsg) (tea.Cmd, string) {
	switch msg.String() {
	case "esc", "h", "backspace":
		t.mode = profileModeList
		return nil, ""

	case "j", "down":
		t.viewScroll++
		return nil, ""

	case "k", "up":
		if t.viewScroll > 0 {
			t.viewScroll--
		}
		return nil, ""

	case "e":
		t.enterEditMode(t.viewProfile)
		return nil, ""

	case " ":
		applied, skipped, errs := profile.Load(t.viewProfile.Name, t.paths)
		if len(errs) > 0 {
			var msgs []string
			for _, err := range errs {
				msgs = append(msgs, err.Error())
			}
			return nil, "Load errors: " + strings.Join(msgs, "; ")
		}
		t.mode = profileModeList
		result := fmt.Sprintf("Loaded %q: %d changes", t.viewProfile.Name, applied)
		if skipped > 0 {
			result += fmt.Sprintf(", %d skipped", skipped)
		}
		return func() tea.Msg { return profileLoadedMsg{} }, result
	}

	return nil, ""
}

func (t *profileTab) handleEditKey(msg tea.KeyMsg) (tea.Cmd, string) {
	switch msg.String() {
	case "esc":
		t.mode = profileModeList
		return nil, ""

	case "j", "down":
		t.editMoveDown()
		return nil, ""

	case "k", "up":
		t.editMoveUp()
		return nil, ""

	case "enter", " ":
		if t.editCursor >= 0 && t.editCursor < len(t.editRows) {
			row := &t.editRows[t.editCursor]
			if row.isHeader {
				t.toggleSection(row.sectionIdx)
			} else {
				row.enabled = !row.enabled
			}
		}
		return nil, ""

	case "tab":
		t.editNextSection()
		return nil, ""

	case "a":
		t.editSetSection(true)
		return nil, ""

	case "x":
		t.editSetSection(false)
		return nil, ""

	case "s":
		t.saveEdit()
		t.refresh()
		return nil, fmt.Sprintf("Updated profile %q", t.editProfile.Name)
	}

	return nil, ""
}

func (t *profileTab) handleSaveKey(msg tea.KeyMsg) (tea.Cmd, string) {
	switch msg.String() {
	case "esc":
		t.mode = profileModeList
		t.nameInput = ""
		return nil, ""

	case "enter":
		name := strings.TrimSpace(t.nameInput)
		if name == "" {
			return nil, "Profile name cannot be empty"
		}
		emptyProfile := profile.Profile{
			Name:    name,
			Created: time.Now().Format(time.RFC3339),
		}
		if err := profile.Update(emptyProfile, t.paths); err != nil {
			return nil, "Save error: " + err.Error()
		}
		t.nameInput = ""
		t.enterEditMode(emptyProfile)
		return nil, fmt.Sprintf("Created %q — now add items", name)

	case "backspace":
		if len(t.nameInput) > 0 {
			t.nameInput = t.nameInput[:len(t.nameInput)-1]
		}
		return nil, ""

	default:
		if len(msg.String()) == 1 {
			t.nameInput += msg.String()
		}
		return nil, ""
	}
}

func (t *profileTab) enterEditMode(p profile.Profile) {
	t.editProfile = p
	t.editRows = buildEditRows(p, t.paths)
	t.editCursor = 0
	t.editOffset = 0
	t.expandedSection = -1
	t.mode = profileModeEdit
}

func buildEditRows(p profile.Profile, paths config.Paths) []editRow {
	skills, rules, agents, plugins, mcpServers := profile.ScanAll(paths)

	enabledSkills := toStringSet(p.Skills)
	enabledRules := toStringSet(p.Rules)
	enabledAgents := toStringSet(p.Agents)
	enabledPlugins := toStringSet(p.Plugins)
	enabledMCP := toStringSet(p.MCPServers)

	var rows []editRow

	rows = append(rows, editRow{isHeader: true, name: "Skills", sectionIdx: 0})
	rows = append(rows, itemsToRows(skills, model.TypeSkill, enabledSkills, 0)...)

	rows = append(rows, editRow{isHeader: true, name: "Rules", sectionIdx: 1})
	rows = append(rows, itemsToRows(rules, model.TypeRule, enabledRules, 1)...)

	rows = append(rows, editRow{isHeader: true, name: "Agents", sectionIdx: 2})
	rows = append(rows, itemsToRows(agents, model.TypeAgent, enabledAgents, 2)...)

	rows = append(rows, editRow{isHeader: true, name: "Plugins", sectionIdx: 3})
	rows = append(rows, itemsToRows(plugins, model.TypePlugin, enabledPlugins, 3)...)

	rows = append(rows, editRow{isHeader: true, name: "MCP Servers", sectionIdx: 4})
	rows = append(rows, itemsToRows(mcpServers, model.TypeMCPServer, enabledMCP, 4)...)

	return rows
}

func itemsToRows(items []model.ConfigItem, itemType model.ItemType, enabled map[string]bool, sectionIdx int) []editRow {
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	var rows []editRow
	for _, item := range items {
		rows = append(rows, editRow{
			id:         item.ID,
			name:       item.Name,
			itemType:   itemType,
			enabled:    enabled[item.ID],
			sectionIdx: sectionIdx,
		})
	}
	return rows
}

func (t *profileTab) isEditRowVisible(i int) bool {
	row := t.editRows[i]
	if row.isHeader {
		return true
	}
	return row.sectionIdx == t.expandedSection
}

func (t *profileTab) editMoveDown() {
	for i := t.editCursor + 1; i < len(t.editRows); i++ {
		if t.isEditRowVisible(i) {
			t.editCursor = i
			t.editFixOffset()
			return
		}
	}
}

func (t *profileTab) editMoveUp() {
	for i := t.editCursor - 1; i >= 0; i-- {
		if t.isEditRowVisible(i) {
			t.editCursor = i
			t.editFixOffset()
			return
		}
	}
}

func (t *profileTab) editNextSection() {
	currentSection := t.editRows[t.editCursor].sectionIdx
	for i := t.editCursor + 1; i < len(t.editRows); i++ {
		if t.editRows[i].isHeader && t.editRows[i].sectionIdx != currentSection {
			t.expandedSection = t.editRows[i].sectionIdx
			t.editCursor = i
			t.editFixOffset()
			return
		}
	}
	for i := 0; i < t.editCursor; i++ {
		if t.editRows[i].isHeader {
			t.expandedSection = t.editRows[i].sectionIdx
			t.editCursor = i
			t.editFixOffset()
			return
		}
	}
}

func (t *profileTab) toggleSection(sectionIdx int) {
	if t.expandedSection == sectionIdx {
		t.expandedSection = -1
	} else {
		t.expandedSection = sectionIdx
	}
	t.editFixOffset()
}

func (t *profileTab) editSetSection(enabled bool) {
	sectionIdx := t.editRows[t.editCursor].sectionIdx
	for i, row := range t.editRows {
		if !row.isHeader && row.sectionIdx == sectionIdx {
			t.editRows[i].enabled = enabled
		}
	}
}

func (t *profileTab) editFixOffset() {
	visiblePos := 0
	for i := 0; i < t.editCursor; i++ {
		if t.isEditRowVisible(i) {
			visiblePos++
		}
	}
	if visiblePos < t.editOffset {
		t.editOffset = visiblePos
	}
	if visiblePos >= t.editOffset+t.height {
		t.editOffset = visiblePos - t.height + 1
	}
}

func (t *profileTab) saveEdit() {
	var skills, rules, agents, plugins, mcpServers []string

	for _, row := range t.editRows {
		if row.isHeader || !row.enabled {
			continue
		}
		switch row.itemType {
		case model.TypeSkill:
			skills = append(skills, row.id)
		case model.TypeRule:
			rules = append(rules, row.id)
		case model.TypeAgent:
			agents = append(agents, row.id)
		case model.TypePlugin:
			plugins = append(plugins, row.id)
		case model.TypeMCPServer:
			mcpServers = append(mcpServers, row.id)
		}
	}

	sort.Strings(skills)
	sort.Strings(rules)
	sort.Strings(agents)
	sort.Strings(plugins)
	sort.Strings(mcpServers)

	t.editProfile.Skills = skills
	t.editProfile.Rules = rules
	t.editProfile.Agents = agents
	t.editProfile.Plugins = plugins
	t.editProfile.MCPServers = mcpServers

	_ = profile.Update(t.editProfile, t.paths)
	t.mode = profileModeList
}

func (t *profileTab) fixOffset() {
	if t.cursor < t.offset {
		t.offset = t.cursor
	}
	if t.cursor >= t.offset+t.height {
		t.offset = t.cursor - t.height + 1
	}
}

func (t *profileTab) render(width int) string {
	switch t.mode {
	case profileModeSave:
		return t.renderSave()
	case profileModeView:
		return t.renderView(width)
	case profileModeEdit:
		return t.renderEdit(width)
	default:
		return t.renderList(width)
	}
}

func (t *profileTab) renderSave() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render("Save Profile"))
	b.WriteString("\n\n")
	b.WriteString("  Name: ")
	b.WriteString(profileInputStyle.Render(t.nameInput + "█"))
	return b.String()
}

func (t *profileTab) renderList(width int) string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("Profiles"))
	b.WriteString("\n")

	if len(t.profiles) == 0 {
		b.WriteString(dimStyle.Render("  No profiles saved yet. Press 'n' to create one."))
		return b.String()
	}

	end := min(t.offset+t.height, len(t.profiles))

	for i := t.offset; i < end; i++ {
		entry := t.profiles[i]
		isCursor := i == t.cursor

		prefix := "  "
		if isCursor {
			prefix = "▸ "
		}

		created := formatDate(entry.Created)

		if isCursor {
			line := fmt.Sprintf(" %s%-20s  %s", prefix, entry.Name, created)
			b.WriteString(selectedStyle.Render(line))
		} else {
			line := fmt.Sprintf(" %s%s  %s", prefix, profileNameStyle.Render(fmt.Sprintf("%-20s", entry.Name)), profileDateStyle.Render(created))
			b.WriteString(line)
		}

		if i < end-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (t *profileTab) renderView(width int) string {
	var lines []string

	p := t.viewProfile
	lines = append(lines, headerStyle.Render(p.Name))
	lines = append(lines, profileDateStyle.Render("  Created: "+formatDate(p.Created)))
	lines = append(lines, "")

	lines = append(lines, renderSection("Skills", p.Skills)...)
	lines = append(lines, renderSection("Rules", p.Rules)...)
	lines = append(lines, renderSection("Agents", p.Agents)...)
	lines = append(lines, renderSection("Plugins", p.Plugins)...)
	lines = append(lines, renderSection("MCP Servers", p.MCPServers)...)

	viewHeight := t.height
	if viewHeight < 5 {
		viewHeight = 5
	}

	maxScroll := len(lines) - viewHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	if t.viewScroll > maxScroll {
		t.viewScroll = maxScroll
	}

	end := min(t.viewScroll+viewHeight, len(lines))
	visible := lines[t.viewScroll:end]

	var b strings.Builder
	b.WriteString(strings.Join(visible, "\n"))

	if maxScroll > 0 {
		b.WriteString("\n")
		b.WriteString(dimStyle.Render(fmt.Sprintf("  (%d/%d)", t.viewScroll+viewHeight, len(lines))))
	}

	return b.String()
}

func (t *profileTab) renderEdit(width int) string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("Edit: " + t.editProfile.Name))
	b.WriteString("\n")

	visIdx := 0
	rendered := 0

	for i, row := range t.editRows {
		if !t.isEditRowVisible(i) {
			continue
		}

		if visIdx < t.editOffset {
			visIdx++
			continue
		}

		if rendered >= t.height {
			break
		}

		if row.isHeader {
			expanded := row.sectionIdx == t.expandedSection
			count := t.countSectionEnabled(i)
			total := t.countSectionTotal(i)

			indicator := "▸"
			if expanded {
				indicator = "▾"
			}

			isCursor := i == t.editCursor
			headerText := fmt.Sprintf("  %s %s (%d/%d)", indicator, row.name, count, total)
			if isCursor {
				b.WriteString("\n")
				b.WriteString(selectedStyle.Render(headerText))
			} else {
				b.WriteString("\n")
				b.WriteString(enabledStyle.Render(headerText))
			}
			b.WriteString("\n")
		} else {
			isCursor := i == t.editCursor

			prefix := "    "
			if isCursor {
				prefix = "  ▸ "
			}

			checkbox := disabledStyle.Render("[ ]")
			if row.enabled {
				checkbox = enabledStyle.Render("[x]")
			}

			nameStr := row.name
			if len(nameStr) > 40 {
				nameStr = nameStr[:37] + "..."
			}

			if isCursor {
				line := fmt.Sprintf("%s%s %s", prefix, checkbox, nameStr)
				b.WriteString(selectedStyle.Render(line))
			} else {
				line := fmt.Sprintf("%s%s %s", prefix, checkbox, dimStyle.Render(nameStr))
				b.WriteString(line)
			}
			b.WriteString("\n")
		}

		visIdx++
		rendered++
	}

	return b.String()
}

func (t *profileTab) countSectionEnabled(headerIdx int) int {
	count := 0
	for i := headerIdx + 1; i < len(t.editRows) && !t.editRows[i].isHeader; i++ {
		if t.editRows[i].enabled {
			count++
		}
	}
	return count
}

func (t *profileTab) countSectionTotal(headerIdx int) int {
	count := 0
	for i := headerIdx + 1; i < len(t.editRows) && !t.editRows[i].isHeader; i++ {
		count++
	}
	return count
}

func renderSection(title string, items []string) []string {
	var lines []string
	count := len(items)
	lines = append(lines, enabledStyle.Render(fmt.Sprintf("  %s (%d)", title, count)))

	if count == 0 {
		lines = append(lines, dimStyle.Render("    (none)"))
	} else {
		for _, item := range items {
			lines = append(lines, normalStyle.Render("    "+item))
		}
	}
	lines = append(lines, "")
	return lines
}

func toStringSet(ids []string) map[string]bool {
	m := make(map[string]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}

func formatDate(rfc3339 string) string {
	t, err := time.Parse(time.RFC3339, rfc3339)
	if err != nil {
		return rfc3339
	}
	return t.Format("2006-01-02 15:04")
}
