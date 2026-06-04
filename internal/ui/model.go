// Package ui implements the zinc terminal UI as a Bubble Tea application.
package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/git"
)

// focusIndex identifies which of the sidebar panels currently has focus.
type focusIndex int

const (
	focusFiles focusIndex = iota
	focusBranches
	focusCommits
	focusStash
	numPanels
)

// Model is the root Bubble Tea model for the application.
type Model struct {
	repo *git.Runner

	width  int
	height int

	focus focusIndex

	status   StatusPanel
	branches placeholderPanel
	commits  placeholderPanel
	stash    placeholderPanel
	main     placeholderPanel

	err error
}

// NewModel constructs the root application model, operating on the git
// repository containing the current working directory.
func NewModel() Model {
	return Model{
		repo:     git.New("."),
		status:   NewStatusPanel(),
		branches: placeholderPanel{title: "Branches", note: "not yet implemented"},
		commits:  placeholderPanel{title: "Commits", note: "not yet implemented"},
		stash:    placeholderPanel{title: "Stash", note: "not yet implemented"},
		main:     placeholderPanel{title: "Diff", note: "select a file to see its diff"},
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
		case "tab":
			m.focus = (m.focus + 1) % numPanels
		case "shift+tab":
			m.focus = (m.focus - 1 + numPanels) % numPanels
		case "up", "k":
			if m.focus == focusFiles {
				m.status.CursorUp()
			}
		case "down", "j":
			if m.focus == focusFiles {
				m.status.CursorDown()
			}
		case "r":
			return m, m.loadStatus
		case " ":
			if m.focus == focusFiles {
				return m, m.toggleStaged
			}
		case "a":
			if m.focus == focusFiles {
				return m, m.stageAll
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	title := titleStyle.Render("zinc")

	bodyWidth := max(m.width-2, 0)
	bodyHeight := max(m.height-6, 0)

	sidebarWidth := bodyWidth / 3
	mainWidth := max(bodyWidth-sidebarWidth, 0)

	panelHeight := bodyHeight / int(numPanels)
	lastPanelHeight := bodyHeight - panelHeight*(int(numPanels)-1)

	var filesView string
	if m.err != nil {
		filesView = panelFocusedStyle.
			Width(sidebarWidth).
			Height(panelHeight).
			Render("error: " + m.err.Error())
	} else {
		filesView = m.status.View(sidebarWidth, panelHeight, m.focus == focusFiles)
	}

	sidebar := lipgloss.JoinVertical(
		lipgloss.Left,
		filesView,
		m.branches.View(sidebarWidth, panelHeight, m.focus == focusBranches),
		m.commits.View(sidebarWidth, panelHeight, m.focus == focusCommits),
		m.stash.View(sidebarWidth, lastPanelHeight, m.focus == focusStash),
	)

	mainPanel := m.main.View(mainWidth, bodyHeight, false)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainPanel)

	help := statusBarStyle.Render(
		keyStyle.Render("tab") + " " + helpDescStyle.Render("switch panel") + "  " +
			keyStyle.Render("↑/↓") + " " + helpDescStyle.Render("navigate") + "  " +
			keyStyle.Render("space") + " " + helpDescStyle.Render("stage/unstage") + "  " +
			keyStyle.Render("a") + " " + helpDescStyle.Render("stage all") + "  " +
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

// toggleStaged stages the selected unstaged/untracked file, or unstages it
// if it already has staged content.
func (m Model) toggleStaged() tea.Msg {
	f, ok := m.status.Selected()
	if !ok {
		return statusLoadedMsg{}
	}

	ctx := context.Background()

	var err error
	if f.Staged != '.' && f.Staged != '?' {
		err = m.repo.Unstage(ctx, f.Path)
	} else {
		err = m.repo.Stage(ctx, f.Path)
	}
	if err != nil {
		return statusLoadedMsg{err: err}
	}

	files, err := m.repo.Status(ctx)
	return statusLoadedMsg{files: files, err: err}
}

func (m Model) stageAll() tea.Msg {
	ctx := context.Background()
	if err := m.repo.StageAll(ctx); err != nil {
		return statusLoadedMsg{err: err}
	}

	files, err := m.repo.Status(ctx)
	return statusLoadedMsg{files: files, err: err}
}
