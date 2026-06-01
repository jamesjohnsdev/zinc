// Package ui implements the zinc terminal UI as a Bubble Tea application.
package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/git"
)

// Model is the root Bubble Tea model for the application.
type Model struct {
	repo *git.Runner

	width  int
	height int

	status StatusPanel

	err error
}

// NewModel constructs the root application model, operating on the git
// repository containing the current working directory.
func NewModel() Model {
	return Model{
		repo:   git.New("."),
		status: NewStatusPanel(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadStatus
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case statusLoadedMsg:
		m.err = msg.err
		if msg.err == nil {
			m.status.SetFiles(msg.files)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.status.CursorUp()
		case "down", "j":
			m.status.CursorDown()
		case "r":
			return m, m.loadStatus
		}
	}

	return m, nil
}

func (m Model) View() string {
	m.status.SetFocused(true)

	title := titleStyle.Render("zinc")

	bodyWidth := max(m.width-2, 0)
	bodyHeight := max(m.height-6, 0)

	var body string
	if m.err != nil {
		body = panelFocusedStyle.
			Width(bodyWidth).
			Height(bodyHeight).
			Render("error: " + m.err.Error())
	} else {
		body = m.status.View(bodyWidth, bodyHeight)
	}

	help := statusBarStyle.Render(
		keyStyle.Render("↑/↓") + " " + helpDescStyle.Render("navigate") + "  " +
			keyStyle.Render("r") + " " + helpDescStyle.Render("refresh") + "  " +
			keyStyle.Render("q") + " " + helpDescStyle.Render("quit"),
	)

	return lipgloss.JoinVertical(lipgloss.Left, title, body, help)
}

// statusLoadedMsg reports the result of refreshing the working tree status.
type statusLoadedMsg struct {
	files []git.FileStatus
	err   error
}

func (m Model) loadStatus() tea.Msg {
	files, err := m.repo.Status(context.Background())
	return statusLoadedMsg{files: files, err: err}
}
