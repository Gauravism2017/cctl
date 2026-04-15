package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Gauravism2017/cctl/internal/config"
	"github.com/Gauravism2017/cctl/internal/model"
	"github.com/Gauravism2017/cctl/internal/ops"
	"github.com/Gauravism2017/cctl/internal/scanner"
)

type tabIndex int

const (
	tabSkills   tabIndex = 0
	tabRules    tabIndex = 1
	tabAgents   tabIndex = 2
	tabPlugins  tabIndex = 3
	tabMCP      tabIndex = 4
	tabProfiles tabIndex = 5
	tabSummary  tabIndex = 6
)

var tabNames = []string{"Skills", "Rules", "Agents", "Plugins", "MCP", "Profiles", "Summary"}

type App struct {
	paths       config.Paths
	tabs        []itemList
	profileTab  profileTab
	activeTab   tabIndex
	filtering   bool
	filterInput string
	width       int
	height      int
	message     string
	quitting    bool
}

func NewApp(paths config.Paths, skills, rules, agents, plugins, mcpServers []model.ConfigItem) App {
	listHeight := 20

	return App{
		paths: paths,
		tabs: []itemList{
			newItemList(skills, listHeight),
			newItemList(rules, listHeight),
			newItemList(agents, listHeight),
			newItemList(plugins, listHeight),
			newItemList(mcpServers, listHeight),
		},
		profileTab: newProfileTab(paths, listHeight),
		activeTab:  tabSkills,
		width:      80,
		height:     30,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		listHeight := msg.Height - 8
		if listHeight < 5 {
			listHeight = 5
		}
		for i := range a.tabs {
			a.tabs[i].height = listHeight
		}
		a.profileTab.height = listHeight
		return a, nil

	case profileLoadedMsg:
		a.rescanAll()
		return a, nil

	case tea.KeyMsg:
		if a.filtering {
			return a.handleFilterKey(msg)
		}
		if a.activeTab == tabProfiles {
			return a.handleProfileKey(msg)
		}
		return a.handleNormalKey(msg)
	}

	return a, nil
}

func (a App) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc":
		a.filtering = false
		return a, nil
	case "backspace":
		if len(a.filterInput) > 0 {
			a.filterInput = a.filterInput[:len(a.filterInput)-1]
			a.currentList().applyFilter(a.filterInput)
		}
		return a, nil
	default:
		if len(msg.String()) == 1 {
			a.filterInput += msg.String()
			a.currentList().applyFilter(a.filterInput)
		}
		return a, nil
	}
}

func (a App) handleProfileKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		a.quitting = true
		return a, tea.Quit
	}

	if a.profileTab.mode == profileModeSave || a.profileTab.mode == profileModeEdit {
		cmd, msg2 := a.profileTab.handleKey(msg)
		if msg2 != "" {
			a.message = msg2
		}
		return a, cmd
	}

	switch msg.String() {
	case "q":
		a.quitting = true
		return a, tea.Quit
	case "tab":
		a.activeTab = (a.activeTab + 1) % tabIndex(len(tabNames))
		a.message = ""
		return a, nil
	case "shift+tab":
		a.activeTab = (a.activeTab - 1 + tabIndex(len(tabNames))) % tabIndex(len(tabNames))
		a.message = ""
		return a, nil
	case "s":
		return a.save()
	}

	cmd, msg2 := a.profileTab.handleKey(msg)
	if msg2 != "" {
		a.message = msg2
	}
	return a, cmd
}

func (a *App) rescanAll() {
	listHeight := a.height - 8
	if listHeight < 5 {
		listHeight = 5
	}

	a.tabs = []itemList{
		newItemList(scanner.ScanSkills(a.paths), listHeight),
		newItemList(scanner.ScanRules(a.paths), listHeight),
		newItemList(scanner.ScanAgents(a.paths), listHeight),
		newItemList(scanner.ScanPlugins(a.paths), listHeight),
		newItemList(scanner.ScanMCPServers(a.paths), listHeight),
	}
	a.profileTab.refresh()
}

func (a App) handleNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		a.quitting = true
		return a, tea.Quit

	case "s":
		return a.save()

	case "tab":
		a.activeTab = (a.activeTab + 1) % tabIndex(len(tabNames))
		a.message = ""
		return a, nil

	case "shift+tab":
		a.activeTab = (a.activeTab - 1 + tabIndex(len(tabNames))) % tabIndex(len(tabNames))
		a.message = ""
		return a, nil

	case "j", "down":
		if a.isItemTab() {
			a.currentList().moveDown()
		}
		return a, nil

	case "k", "up":
		if a.isItemTab() {
			a.currentList().moveUp()
		}
		return a, nil

	case " ":
		if a.isItemTab() {
			a.currentList().toggleCurrent()
			a.message = ""
		}
		return a, nil

	case "/":
		if a.isItemTab() {
			a.filtering = true
			a.filterInput = ""
		}
		return a, nil

	case "a":
		if a.isItemTab() {
			a.currentList().enableAll()
			a.message = "Enabled all visible"
		}
		return a, nil

	case "n":
		if a.isItemTab() {
			a.currentList().disableAll()
			a.message = "Disabled all visible"
		}
		return a, nil

	case "c":
		if a.isItemTab() {
			a.filterInput = ""
			a.currentList().applyFilter("")
			a.message = "Filter cleared"
		}
		return a, nil
	}

	return a, nil
}

func (a App) save() (tea.Model, tea.Cmd) {
	var allItems []model.ConfigItem
	for _, tab := range a.tabs[:5] {
		allItems = append(allItems, tab.items...)
	}

	errs := ops.ApplyChanges(allItems, a.paths)
	if len(errs) > 0 {
		var errMsgs []string
		for _, err := range errs {
			errMsgs = append(errMsgs, err.Error())
		}
		a.message = "Errors: " + strings.Join(errMsgs, "; ")
		return a, nil
	}

	dirty := 0
	for _, tab := range a.tabs[:5] {
		dirty += tab.countDirty()
	}

	a.message = fmt.Sprintf("Saved %d changes", dirty)
	a.quitting = true
	return a, tea.Quit
}

func (a App) isItemTab() bool {
	return a.activeTab >= tabSkills && a.activeTab <= tabMCP
}

func (a *App) currentList() *itemList {
	return &a.tabs[a.activeTab]
}

func (a App) View() string {
	if a.quitting {
		if a.message != "" {
			return a.message + "\n"
		}
		return "No changes.\n"
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("cctl — Claude Code Context Manager"))
	b.WriteString("\n")

	var tabs []string
	for i, name := range tabNames {
		if tabIndex(i) == a.activeTab {
			tabs = append(tabs, activeTabStyle.Render(name))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(name))
		}
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	b.WriteString("\n\n")

	if a.activeTab == tabSummary {
		b.WriteString(a.renderSummary())
	} else if a.activeTab == tabProfiles {
		b.WriteString(a.profileTab.render(a.width))
	} else {
		list := a.tabs[a.activeTab]
		enabled, total := list.countEnabled()
		dirty := list.countDirty()

		statusLine := fmt.Sprintf("  Active: %d / %d", enabled, total)
		if dirty > 0 {
			statusLine += dirtyStyle.Render(fmt.Sprintf("  (%d pending changes)", dirty))
		}

		if a.filtering {
			b.WriteString(filterStyle.Render("  Filter: " + a.filterInput + "█"))
			b.WriteString(statusBarStyle.Render(statusLine))
		} else if a.filterInput != "" {
			b.WriteString(filterStyle.Render("  Filter: " + a.filterInput))
			b.WriteString(statusBarStyle.Render(statusLine))
		} else {
			b.WriteString(statusBarStyle.Render(statusLine))
		}
		b.WriteString("\n\n")

		b.WriteString(list.render(a.width))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if a.message != "" {
		b.WriteString(dirtyStyle.Render("  " + a.message))
		b.WriteString("\n")
	}

	if a.filtering {
		b.WriteString(helpStyle.Render("  Enter: apply filter  Esc: cancel  Backspace: delete"))
	} else if a.activeTab == tabProfiles {
		switch a.profileTab.mode {
		case profileModeSave:
			b.WriteString(helpStyle.Render("  Enter: save  Esc: cancel"))
		case profileModeView:
			b.WriteString(helpStyle.Render("  Space: load  e: edit  j/k: scroll  Esc: back"))
		case profileModeEdit:
			b.WriteString(helpStyle.Render("  Enter: expand/collapse  Space: toggle  Tab: next section  a/x: enable/disable section  s: save  Esc: cancel"))
		default:
			b.WriteString(helpStyle.Render("  Enter: view  Space: load  e: edit  d: delete  n: new  Tab: next tab  s: save & exit  q: quit"))
		}
	} else {
		b.WriteString(helpStyle.Render("  Space: toggle  /: filter  c: clear filter  Tab: next tab  a: enable all  n: disable all  s: save & exit  q: quit"))
	}

	return b.String()
}

func (a App) renderSummary() string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("Summary"))
	b.WriteString("\n")

	totalDirty := 0
	for i, name := range tabNames[:5] {
		enabled, total := a.tabs[i].countEnabled()
		dirty := a.tabs[i].countDirty()
		totalDirty += dirty

		line := fmt.Sprintf("  %-10s %d enabled / %d total", name, enabled, total)
		if dirty > 0 {
			line += dirtyStyle.Render(fmt.Sprintf("  (%d changes)", dirty))
		}
		b.WriteString(normalStyle.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(normalStyle.Render(fmt.Sprintf("  %-10s %d saved", "Profiles", len(a.profileTab.profiles))))
	b.WriteString("\n\n")
	if totalDirty > 0 {
		b.WriteString(dirtyStyle.Render(fmt.Sprintf("  %d total pending changes — press 's' to save", totalDirty)))
	} else {
		b.WriteString(normalStyle.Render("  No pending changes"))
	}
	b.WriteString("\n")

	return b.String()
}
