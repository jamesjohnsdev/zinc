// Package ui implements the zinc terminal UI as a Bubble Tea application.
package ui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
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
	branches BranchesPanel
	commits  placeholderPanel
	stash    placeholderPanel
	main     DiffPanel

	commitInput CommitInput

	err error
}

// NewModel constructs the root application model, operating on the git
// repository containing the current working directory.
func NewModel() Model {
	return Model{
		repo:     git.New("."),
		status:   NewStatusPanel(),
		branches: NewBranchesPanel(),
		commits:  placeholderPanel{title: "Commits", note: "not yet implemented"},
		stash:    placeholderPanel{title: "Stash", note: "not yet implemented"},
		main:     NewDiffPanel(),

		commitInput: NewCommitInput(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadStatus, m.loadBranches)
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
		return m, m.loadDiff

	case branchesLoadedMsg:
		if msg.err == nil {
			m.branches.SetBranches(msg.branches)
		} else {
			m.err = msg.err
		}
		return m, nil

	case diffLoadedMsg:
		content := msg.content
		if msg.err != nil {
			content = "error: " + msg.err.Error()
		}
		m.main.SetContent(msg.title, content)
		return m, nil

	case tea.KeyMsg:
		if m.commitInput.Active() {
			switch msg.String() {
			case "esc":
				m.commitInput.Close()
				return m, nil
			case "enter":
				message := strings.TrimSpace(m.commitInput.Value())
				m.commitInput.Close()
				if message == "" {
					return m, nil
				}
				return m, m.commitCmd(message)
			default:
				cmd := m.commitInput.Update(msg)
				return m, cmd
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focus = (m.focus + 1) % numPanels
		case "shift+tab":
			m.focus = (m.focus - 1 + numPanels) % numPanels
		case "up", "k":
			switch m.focus {
			case focusFiles:
				m.status.CursorUp()
				return m, m.loadDiff
			case focusBranches:
				m.branches.CursorUp()
			}
		case "down", "j":
			switch m.focus {
			case focusFiles:
				m.status.CursorDown()
				return m, m.loadDiff
			case focusBranches:
				m.branches.CursorDown()
			}
		case "ctrl+u":
			m.main.PageUp()
		case "ctrl+d":
			m.main.PageDown()
		case "r":
			return m, tea.Batch(m.loadStatus, m.loadBranches)
		case " ":
			if m.focus == focusFiles {
				return m, m.toggleStaged
			}
		case "a":
			if m.focus == focusFiles {
				return m, m.stageAll
			}
		case "c":
			if m.focus == focusFiles {
				m.commitInput.Open()
				return m, textinput.Blink
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
			keyStyle.Render("c") + " " + helpDescStyle.Render("commit") + "  " +
			keyStyle.Render("^u/^d") + " " + helpDescStyle.Render("scroll diff") + "  " +
			keyStyle.Render("r") + " " + helpDescStyle.Render("refresh") + "  " +
			keyStyle.Render("q") + " " + helpDescStyle.Render("quit"),
	)

	parts := []string{title}
	if m.commitInput.Active() {
		parts = append(parts, m.commitInput.View(bodyWidth))
	}
	parts = append(parts, body, help)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
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

// branchesLoadedMsg reports the result of refreshing the branch list.
type branchesLoadedMsg struct {
	branches []git.Branch
	err      error
}

func (m Model) loadBranches() tea.Msg {
	ctx := context.Background()

	local, err := m.repo.Branches(ctx)
	if err != nil {
		return branchesLoadedMsg{err: err}
	}

	remote, err := m.repo.RemoteBranches(ctx)
	if err != nil {
		return branchesLoadedMsg{err: err}
	}

	return branchesLoadedMsg{branches: append(local, remote...)}
}

// diffLoadedMsg reports the result of loading the diff for the currently
// selected file.
type diffLoadedMsg struct {
	title   string
	content string
	err     error
}

func (m Model) loadDiff() tea.Msg {
	f, ok := m.status.Selected()
	if !ok {
		return diffLoadedMsg{title: "Diff"}
	}

	ctx := context.Background()

	var (
		content string
		err     error
	)
	if f.Untracked {
		content, err = m.repo.DiffUntracked(ctx, f.Path)
	} else {
		staged := f.Staged != '.' && f.Staged != '?'
		content, err = m.repo.Diff(ctx, f.Path, staged)
	}

	return diffLoadedMsg{title: "Diff: " + f.Path, content: content, err: err}
}

// commitCmd commits the currently staged changes with message and refreshes
// the working tree status.
func (m Model) commitCmd(message string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		if err := m.repo.Commit(ctx, message); err != nil {
			return statusLoadedMsg{err: err}
		}

		files, err := m.repo.Status(ctx)
		return statusLoadedMsg{files: files, err: err}
	}
}
