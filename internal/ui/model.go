// Package ui implements the zinc terminal UI as a Bubble Tea application.
package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/gh"
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

// screenID identifies which top-level screen is shown: the git working
// tree, or the GitHub PRs/issues/workflow runs view.
type screenID int

const (
	screenGit screenID = iota
	screenGitHub
	numScreens
)

// ghFocusIndex identifies which of the GitHub screen's panels has focus.
type ghFocusIndex int

const (
	focusPRs ghFocusIndex = iota
	focusIssues
	focusRuns
	ghNumPanels
)

// Model is the root Bubble Tea model for the application.
type Model struct {
	repo *git.Runner
	gh   *gh.Runner
	// ghRepo targets a specific "[HOST/]OWNER/REPO"; empty resolves from
	// the working directory's git remote (gh's default). Every gh package
	// call already takes this explicitly, so a repo-picker panel can set
	// it later without changing how PRs/issues/runs are loaded.
	ghRepo string

	width  int
	height int

	screen  screenID
	focus   focusIndex
	ghFocus ghFocusIndex

	// ghLoaded defers the first PR/issue/run/repo fetch until the user
	// actually switches to the GitHub screen, so plain git usage never
	// shells out to gh (which needs network + auth) unasked.
	ghLoaded bool

	status   StatusPanel
	branches BranchesPanel
	commits  LogPanel
	stash    StashPanel
	main     DiffPanel

	repoInfo gh.Repo
	prs      PRPanel
	issues   IssuesPanel
	runs     RunsPanel
	ghDetail TextViewport

	commitInput  TextPrompt
	branchInput  TextPrompt
	commentInput TextPrompt
	confirm      ConfirmDialog

	// pendingConfirm runs when the open confirm dialog is accepted with 'y'.
	pendingConfirm tea.Cmd

	err error
}

// NewModel constructs the root application model, operating on the git
// repository containing the current working directory.
func NewModel() Model {
	return Model{
		repo:     git.New("."),
		gh:       gh.New("."),
		status:   NewStatusPanel(),
		branches: NewBranchesPanel(),
		commits:  NewLogPanel(),
		stash:    NewStashPanel(),
		main:     NewDiffPanel(),

		prs:      NewPRPanel(),
		issues:   NewIssuesPanel(),
		runs:     NewRunsPanel(),
		ghDetail: NewTextViewport(),

		commitInput:  NewTextPrompt("Commit message", "commit message"),
		branchInput:  NewTextPrompt("New branch", "new branch name"),
		commentInput: NewTextPrompt("Comment", "comment body"),
		confirm:      NewConfirmDialog(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.loadStatus, m.loadBranches, m.loadLog, m.loadStash)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		l := m.layout()
		m.main.SetSize(l.mainWidth, l.mainHeight)

		return m, nil

	case statusLoadedMsg:
		m.err = msg.err
		if msg.err == nil {
			m.status.SetFiles(msg.files)
		}
		return m, m.loadDiff

	case branchesLoadedMsg:
		m.err = msg.err
		if msg.err == nil {
			m.branches.SetBranches(msg.branches)
		}
		return m, nil

	case diffLoadedMsg:
		content := msg.content
		if msg.err != nil {
			content = "error: " + msg.err.Error()
		}
		m.main.SetContent(msg.title, content)
		return m, nil

	case logLoadedMsg:
		m.err = msg.err
		if msg.err == nil {
			m.commits.SetCommits(msg.commits)
		}
		return m, nil

	case stashLoadedMsg:
		m.err = msg.err
		if msg.err == nil {
			m.stash.SetStashes(msg.stashes)
		}
		return m, nil

	case refreshMsg:
		m.err = msg.err
		return m, tea.Batch(m.loadStatus, m.loadBranches, m.loadLog, m.loadStash, m.loadDiff)

	case repoLoadedMsg:
		m.err = msg.err
		if msg.err == nil {
			m.repoInfo = msg.repo
		}
		return m, nil

	case prsLoadedMsg:
		m.err = msg.err
		if msg.err == nil {
			m.prs.SetPRs(msg.prs)
		}
		return m, m.loadGHDetail

	case issuesLoadedMsg:
		m.err = msg.err
		if msg.err == nil {
			m.issues.SetIssues(msg.issues)
		}
		return m, m.loadGHDetail

	case runsLoadedMsg:
		m.err = msg.err
		if msg.err == nil {
			m.runs.SetRuns(msg.runs)
		}
		return m, m.loadGHDetail

	case ghDetailMsg:
		content := msg.content
		if msg.err != nil {
			content = "error: " + msg.err.Error()
		}
		m.ghDetail.SetContent(msg.title, content)
		return m, nil

	case ghRefreshMsg:
		m.err = msg.err
		return m, m.reloadGH()

	case tea.KeyMsg:
		if m.branchInput.Active() {
			switch msg.String() {
			case "esc":
				m.branchInput.Close()
				return m, nil
			case "enter":
				name := strings.TrimSpace(m.branchInput.Value())
				m.branchInput.Close()
				if name == "" {
					return m, nil
				}
				return m, m.createBranchCmd(name)
			default:
				cmd := m.branchInput.Update(msg)
				return m, cmd
			}
		}

		if m.confirm.Active() {
			switch msg.String() {
			case "y":
				cmd := m.pendingConfirm
				m.confirm.Close()
				m.pendingConfirm = nil
				return m, cmd
			case "n", "esc":
				m.confirm.Close()
				m.pendingConfirm = nil
			}
			return m, nil
		}

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

		if m.commentInput.Active() {
			switch msg.String() {
			case "esc":
				m.commentInput.Close()
				return m, nil
			case "enter":
				body := strings.TrimSpace(m.commentInput.Value())
				m.commentInput.Close()
				issue, ok := m.issues.Selected()
				if body == "" || !ok {
					return m, nil
				}
				return m, m.issueCommentCmd(issue.Number, body)
			default:
				cmd := m.commentInput.Update(msg)
				return m, cmd
			}
		}

		if msg.String() == "g" {
			m.screen = (m.screen + 1) % numScreens
			if m.screen == screenGitHub && !m.ghLoaded {
				m.ghLoaded = true
				return m, tea.Batch(m.loadRepo, m.loadPRs, m.loadIssues, m.loadRuns)
			}
			return m, nil
		}

		if m.screen == screenGitHub {
			return m.updateGitHubScreen(msg)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focus = (m.focus + 1) % numPanels
		case "shift+tab":
			m.focus = (m.focus - 1 + numPanels) % numPanels
		case "1":
			m.focus = focusFiles
		case "2":
			m.focus = focusBranches
		case "3":
			m.focus = focusCommits
		case "4":
			m.focus = focusStash
		case "up", "k":
			switch m.focus {
			case focusFiles:
				m.status.CursorUp()
				return m, m.loadDiff
			case focusBranches:
				m.branches.CursorUp()
			case focusCommits:
				m.commits.CursorUp()
			case focusStash:
				m.stash.CursorUp()
			}
		case "down", "j":
			switch m.focus {
			case focusFiles:
				m.status.CursorDown()
				return m, m.loadDiff
			case focusBranches:
				m.branches.CursorDown()
			case focusCommits:
				m.commits.CursorDown()
			case focusStash:
				m.stash.CursorDown()
			}
		case "ctrl+u":
			m.main.PageUp()
		case "ctrl+d":
			m.main.PageDown()
		case "r":
			return m, tea.Batch(m.loadStatus, m.loadBranches, m.loadLog, m.loadStash)
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
		case "enter":
			if m.focus == focusBranches {
				if b, ok := m.branches.Selected(); ok {
					return m, m.checkoutCmd(b)
				}
			}
		case "n":
			if m.focus == focusBranches {
				m.branchInput.Open()
				return m, textinput.Blink
			}
		case "d":
			switch m.focus {
			case focusBranches:
				if b, ok := m.branches.Selected(); ok && !b.Remote && !b.Current {
					m.confirm.Open("Delete branch '" + b.Name + "'?")
					m.pendingConfirm = m.deleteBranchCmd(b.Name)
				}
			case focusStash:
				if s, ok := m.stash.Selected(); ok {
					m.confirm.Open(fmt.Sprintf("Drop stash@{%d}?", s.Index))
					m.pendingConfirm = m.stashDropCmd(s.Index)
				}
			}
		case "s":
			if m.focus == focusFiles {
				return m, m.stashPushCmd
			}
		case "D":
			if m.focus == focusFiles {
				if f, ok := m.status.Selected(); ok {
					m.confirm.Open("Discard changes to '" + f.Path + "'?")
					m.pendingConfirm = m.discardCmd(f.Path, f.Untracked)
				}
			}
		case "p":
			if m.focus == focusStash {
				if s, ok := m.stash.Selected(); ok {
					return m, m.stashPopCmd(s.Index)
				}
			}
		}
	}

	return m, nil
}

// uiLayout holds the computed dimensions for one frame, plus the
// already-rendered fixed-height chrome (title, optional error toast,
// optional prompt/confirm overlay, help bar) that everything else's size is
// measured against. Both View and the resize handler build this the same
// way so their numbers can never drift apart.
type uiLayout struct {
	sidebarWidth    int
	mainWidth       int
	panelHeight     int
	lastPanelHeight int
	mainHeight      int
	chrome          []string // title, [toast], [overlay], help — in render order
}

func (m Model) layout() uiLayout {
	bodyWidth := max(m.width-2, 0)

	chrome := []string{m.headerBar()}
	if m.err != nil {
		chrome = append(chrome, errorToastStyle.Render("✗ "+m.err.Error()))
	}
	switch {
	case m.commitInput.Active():
		chrome = append(chrome, m.commitInput.View(bodyWidth))
	case m.branchInput.Active():
		chrome = append(chrome, m.branchInput.View(bodyWidth))
	case m.commentInput.Active():
		chrome = append(chrome, m.commentInput.View(bodyWidth))
	case m.confirm.Active():
		chrome = append(chrome, m.confirm.View(bodyWidth))
	}
	chrome = append(chrome, m.helpBar(m.width))

	chromeHeight := 0
	for _, s := range chrome {
		chromeHeight += lipgloss.Height(s)
	}

	bodyHeight := max(m.height-chromeHeight, 0)

	sidebarWidth := bodyWidth / 3
	mainWidth := max(bodyWidth-sidebarWidth, 0)

	panelCount := int(numPanels)
	if m.screen == screenGitHub {
		panelCount = int(ghNumPanels)
	}
	panelHeight, lastPanelHeight := splitPanelHeights(bodyHeight, panelCount)
	mainHeight := max(bodyHeight-borderRows, 0)

	return uiLayout{
		sidebarWidth:    sidebarWidth,
		mainWidth:       mainWidth,
		panelHeight:     panelHeight,
		lastPanelHeight: lastPanelHeight,
		mainHeight:      mainHeight,
		chrome:          chrome,
	}
}

func (m Model) View() string {
	l := m.layout()

	var body string
	if m.screen == screenGitHub {
		body = m.viewGitHubScreen(l)
	} else {
		body = m.viewGitScreen(l)
	}

	help := l.chrome[len(l.chrome)-1]

	var parts []string
	parts = append(parts, l.chrome[:len(l.chrome)-1]...) // title, [toast], [overlay]
	parts = append(parts, body, help)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// viewGitScreen renders the Files/Branches/Commits/Stash sidebar and the
// diff main panel.
func (m Model) viewGitScreen(l uiLayout) string {
	sidebar := lipgloss.JoinVertical(
		lipgloss.Left,
		m.status.View(l.sidebarWidth, l.panelHeight, m.focus == focusFiles),
		m.branches.View(l.sidebarWidth, l.panelHeight, m.focus == focusBranches),
		m.commits.View(l.sidebarWidth, l.panelHeight, m.focus == focusCommits),
		m.stash.View(l.sidebarWidth, l.lastPanelHeight, m.focus == focusStash),
	)

	mainPanel := m.main.View(l.mainWidth, l.mainHeight, false)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainPanel)
}

// statusLoadedMsg reports the result of refreshing the working tree status.
type statusLoadedMsg struct {
	files []git.FileStatus
	err   error
}

// headerBar renders the full-width app title bar: which screen is active,
// plus the current branch (Git screen) or target repo (GitHub screen)
// once known.
func (m Model) headerBar() string {
	if m.screen == screenGitHub {
		header := " zinc · GitHub "
		if m.repoInfo.NameWithOwner != "" {
			header += "· " + m.repoInfo.NameWithOwner + " "
		}
		return headerStyle.Width(m.width).Render(header)
	}

	header := " zinc · Git "
	if name, ok := m.branches.Current(); ok {
		header += "· " + name + " "
	}

	return headerStyle.Width(m.width).Render(header)
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

// logLoadedMsg reports the result of refreshing the commit log.
type logLoadedMsg struct {
	commits []git.Commit
	err     error
}

func (m Model) loadLog() tea.Msg {
	commits, err := m.repo.Log(context.Background(), 200)
	return logLoadedMsg{commits: commits, err: err}
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

// refreshMsg reports that a branch action completed and every panel should
// reload its data.
type refreshMsg struct {
	err error
}

// checkoutCmd checks out a local branch, or creates a local tracking branch
// first if b is a remote-tracking branch.
func (m Model) checkoutCmd(b git.Branch) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		var err error
		if b.Remote {
			err = m.repo.CheckoutRemote(ctx, b.Name)
		} else {
			err = m.repo.Checkout(ctx, b.Name)
		}

		return refreshMsg{err: err}
	}
}

// createBranchCmd creates and checks out a new branch from HEAD.
func (m Model) createBranchCmd(name string) tea.Cmd {
	return func() tea.Msg {
		err := m.repo.CreateBranch(context.Background(), name)
		return refreshMsg{err: err}
	}
}

// deleteBranchCmd deletes a local branch.
func (m Model) deleteBranchCmd(name string) tea.Cmd {
	return func() tea.Msg {
		err := m.repo.DeleteBranch(context.Background(), name, false)
		return refreshMsg{err: err}
	}
}

// stashLoadedMsg reports the result of refreshing the stash list.
type stashLoadedMsg struct {
	stashes []git.Stash
	err     error
}

func (m Model) loadStash() tea.Msg {
	stashes, err := m.repo.StashList(context.Background())
	return stashLoadedMsg{stashes: stashes, err: err}
}

// stashPushCmd stashes all working tree changes, including untracked files.
func (m Model) stashPushCmd() tea.Msg {
	err := m.repo.StashPush(context.Background(), "", true)
	return refreshMsg{err: err}
}

// stashPopCmd applies and removes the stash entry at index.
func (m Model) stashPopCmd(index int) tea.Cmd {
	return func() tea.Msg {
		err := m.repo.StashPop(context.Background(), index)
		return refreshMsg{err: err}
	}
}

// stashDropCmd removes the stash entry at index without applying it.
func (m Model) stashDropCmd(index int) tea.Cmd {
	return func() tea.Msg {
		err := m.repo.StashDrop(context.Background(), index)
		return refreshMsg{err: err}
	}
}

// discardCmd reverts path's working tree changes.
func (m Model) discardCmd(path string, untracked bool) tea.Cmd {
	return func() tea.Msg {
		err := m.repo.Discard(context.Background(), path, untracked)
		return refreshMsg{err: err}
	}
}
