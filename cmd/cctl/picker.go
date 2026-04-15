package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	pickerCursor  = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).SetString("▸ ")
	pickerNormal  = lipgloss.NewStyle().PaddingLeft(2)
	pickerTitle   = lipgloss.NewStyle().Bold(true).MarginBottom(1)
)

type pickerModel struct {
	choices  []string
	cursor   int
	selected string
	aborted  bool
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.aborted = true
			return m, tea.Quit
		case "enter":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m pickerModel) View() string {
	s := pickerTitle.Render("Select profile:") + "\n"
	for i, choice := range m.choices {
		if i == m.cursor {
			s += pickerCursor.Render(choice) + "\n"
		} else {
			s += pickerNormal.Render(choice) + "\n"
		}
	}
	s += "\n" + lipgloss.NewStyle().Faint(true).Render("j/k to move · enter to select · q to cancel")
	return s
}

func runPicker(choices []string) (string, error) {
	m := pickerModel{choices: choices}
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("picker error: %w", err)
	}

	final := result.(pickerModel)
	if final.aborted {
		return "", fmt.Errorf("cancelled")
	}
	return final.selected, nil
}
