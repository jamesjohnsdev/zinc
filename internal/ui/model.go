// Package ui implements the zinc terminal UI as a Bubble Tea application.
package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is the root Bubble Tea model for the application.
type Model struct {
	width  int
	height int
}

// NewModel constructs the root application model.
func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	title := titleStyle.Render("zinc")
	body := panelFocusedStyle.
		Width(max(m.width-2, 0)).
		Height(max(m.height-6, 0)).
		Render("git TUI — press q to quit")
	help := statusBarStyle.Render(keyStyle.Render("q") + " " + helpDescStyle.Render("quit"))

	return lipgloss.JoinVertical(lipgloss.Left, title, body, help)
}
