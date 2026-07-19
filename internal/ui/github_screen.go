package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/gh"
)

// updateGitHubScreen handles key presses for the GitHub screen once the
// shared overlay checks (prompts, confirm dialog) in Update have already
// been ruled out.
func (m Model) updateGitHubScreen(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "tab":
		m.ghFocus = (m.ghFocus + 1) % ghNumPanels
		return m, m.loadGHDetail
	case "shift+tab":
		m.ghFocus = (m.ghFocus - 1 + ghNumPanels) % ghNumPanels
		return m, m.loadGHDetail
	case "1":
		m.ghFocus = focusPRs
		return m, m.loadGHDetail
	case "2":
		m.ghFocus = focusIssues
		return m, m.loadGHDetail
	case "3":
		m.ghFocus = focusRuns
		return m, m.loadGHDetail

	case "up", "k":
		switch m.ghFocus {
		case focusPRs:
			m.prs.CursorUp()
		case focusIssues:
			m.issues.CursorUp()
		case focusRuns:
			m.runs.CursorUp()
		}
		return m, m.loadGHDetail
	case "down", "j":
		switch m.ghFocus {
		case focusPRs:
			m.prs.CursorDown()
		case focusIssues:
			m.issues.CursorDown()
		case focusRuns:
			m.runs.CursorDown()
		}
		return m, m.loadGHDetail

	case "ctrl+u":
		m.ghDetail.PageUp()
	case "ctrl+d":
		m.ghDetail.PageDown()

	case "r":
		return m, tea.Batch(m.loadRepo, m.loadPRs, m.loadIssues, m.loadRuns)

	case "c": // checkout PR
		if m.ghFocus == focusPRs {
			if pr, ok := m.prs.Selected(); ok {
				return m, m.prCheckoutCmd(pr.Number)
			}
		}

	case "m": // squash-merge PR
		if m.ghFocus == focusPRs {
			if pr, ok := m.prs.Selected(); ok && !pr.IsDraft {
				m.confirm.Open(fmt.Sprintf("Squash-merge PR #%d?", pr.Number))
				m.pendingConfirm = m.prMergeCmd(pr.Number, gh.MergeSquash)
			}
		}

	case "x": // close PR / close-or-reopen issue / cancel run
		switch m.ghFocus {
		case focusPRs:
			if pr, ok := m.prs.Selected(); ok && pr.State == "OPEN" {
				m.confirm.Open(fmt.Sprintf("Close PR #%d?", pr.Number))
				m.pendingConfirm = m.prCloseCmd(pr.Number)
			}
		case focusIssues:
			if issue, ok := m.issues.Selected(); ok {
				if issue.State == "CLOSED" {
					m.confirm.Open(fmt.Sprintf("Reopen issue #%d?", issue.Number))
					m.pendingConfirm = m.issueReopenCmd(issue.Number)
				} else {
					m.confirm.Open(fmt.Sprintf("Close issue #%d?", issue.Number))
					m.pendingConfirm = m.issueCloseCmd(issue.Number)
				}
			}
		case focusRuns:
			if run, ok := m.runs.Selected(); ok && run.Status != "completed" {
				m.confirm.Open("Cancel run '" + run.WorkflowName + "'?")
				m.pendingConfirm = m.runCancelCmd(run.DatabaseID)
			}
		}

	case "C": // comment on issue
		if m.ghFocus == focusIssues {
			if _, ok := m.issues.Selected(); ok {
				m.commentInput.Open()
				return m, textinput.Blink
			}
		}

	case "R": // rerun workflow
		if m.ghFocus == focusRuns {
			if run, ok := m.runs.Selected(); ok {
				return m, m.runRerunCmd(run.DatabaseID)
			}
		}

	case "o": // open in browser
		switch m.ghFocus {
		case focusPRs:
			if pr, ok := m.prs.Selected(); ok {
				return m, m.prWebCmd(pr.Number)
			}
		case focusIssues:
			if issue, ok := m.issues.Selected(); ok {
				return m, m.issueWebCmd(issue.Number)
			}
		case focusRuns:
			if run, ok := m.runs.Selected(); ok {
				return m, m.runWebCmd(run.DatabaseID)
			}
		}
	}

	return m, nil
}

// viewGitHubScreen renders the PRs/Issues/Runs sidebar and the detail main
// panel.
func (m Model) viewGitHubScreen(l uiLayout) string {
	sidebar := lipgloss.JoinVertical(
		lipgloss.Left,
		m.prs.View(l.sidebarWidth, l.panelHeight, m.ghFocus == focusPRs),
		m.issues.View(l.sidebarWidth, l.panelHeight, m.ghFocus == focusIssues),
		m.runs.View(l.sidebarWidth, l.lastPanelHeight, m.ghFocus == focusRuns),
	)

	detail := m.ghDetail.View(l.mainWidth, l.mainHeight, false)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, detail)
}

// repoLoadedMsg reports the result of loading the target repo's summary.
type repoLoadedMsg struct {
	repo gh.Repo
	err  error
}

func (m Model) loadRepo() tea.Msg {
	repo, err := m.gh.ViewRepo(context.Background(), m.ghRepo)
	return repoLoadedMsg{repo: repo, err: err}
}

// prsLoadedMsg reports the result of refreshing the pull request list.
type prsLoadedMsg struct {
	prs []gh.PullRequest
	err error
}

func (m Model) loadPRs() tea.Msg {
	prs, err := m.gh.PullRequests(context.Background(), m.ghRepo)
	return prsLoadedMsg{prs: prs, err: err}
}

// issuesLoadedMsg reports the result of refreshing the issue list.
type issuesLoadedMsg struct {
	issues []gh.Issue
	err    error
}

func (m Model) loadIssues() tea.Msg {
	issues, err := m.gh.Issues(context.Background(), m.ghRepo)
	return issuesLoadedMsg{issues: issues, err: err}
}

// runsLoadedMsg reports the result of refreshing the workflow run list.
type runsLoadedMsg struct {
	runs []gh.WorkflowRun
	err  error
}

func (m Model) loadRuns() tea.Msg {
	runs, err := m.gh.WorkflowRuns(context.Background(), m.ghRepo)
	return runsLoadedMsg{runs: runs, err: err}
}

// ghDetailMsg reports the result of loading detail for whichever GitHub
// screen panel currently has focus.
type ghDetailMsg struct {
	title   string
	content string
	err     error
}

// loadGHDetail fetches detail for the item selected in the currently
// focused panel. PR/issue detail (including body) requires a gh call;
// run detail is built entirely from data the list already has.
func (m Model) loadGHDetail() tea.Msg {
	switch m.ghFocus {
	case focusPRs:
		pr, ok := m.prs.Selected()
		if !ok {
			return ghDetailMsg{title: "Detail"}
		}
		full, err := m.gh.ViewPR(context.Background(), m.ghRepo, pr.Number)
		if err != nil {
			return ghDetailMsg{title: fmt.Sprintf("PR #%d", pr.Number), err: err}
		}
		return ghDetailMsg{
			title:   fmt.Sprintf("PR #%d: %s", full.Number, full.Title),
			content: nonEmpty(full.Body, "(no description)"),
		}

	case focusIssues:
		issue, ok := m.issues.Selected()
		if !ok {
			return ghDetailMsg{title: "Detail"}
		}
		full, err := m.gh.ViewIssue(context.Background(), m.ghRepo, issue.Number)
		if err != nil {
			return ghDetailMsg{title: fmt.Sprintf("Issue #%d", issue.Number), err: err}
		}
		return ghDetailMsg{
			title:   fmt.Sprintf("Issue #%d: %s", full.Number, full.Title),
			content: nonEmpty(full.Body, "(no description)"),
		}

	case focusRuns:
		run, ok := m.runs.Selected()
		if !ok {
			return ghDetailMsg{title: "Detail"}
		}
		return ghDetailMsg{title: "Run: " + run.WorkflowName, content: formatRunDetail(run)}
	}

	return ghDetailMsg{title: "Detail"}
}

func nonEmpty(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func formatRunDetail(run gh.WorkflowRun) string {
	lines := []string{
		"Workflow:   " + run.WorkflowName,
		"Branch:     " + run.HeadBranch,
		"Event:      " + run.Event,
		"Status:     " + run.Status,
		"Conclusion: " + nonEmpty(run.Conclusion, "(pending)"),
		"Created:    " + run.CreatedAt.Local().Format(time.RFC1123),
		"URL:        " + run.URL,
	}
	return strings.Join(lines, "\n")
}

// ghRefreshMsg reports that a PR/issue/run action completed and every
// GitHub panel should reload its data.
type ghRefreshMsg struct {
	err error
}

func (m Model) reloadGH() tea.Cmd {
	return tea.Batch(m.loadPRs, m.loadIssues, m.loadRuns)
}

func (m Model) prCheckoutCmd(number int) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.PRCheckout(context.Background(), m.ghRepo, number)
		return refreshMsg{err: err} // checkout changes the working tree; refresh the Git screen
	}
}

func (m Model) prMergeCmd(number int, method gh.MergeMethod) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.PRMerge(context.Background(), m.ghRepo, number, method)
		return ghRefreshMsg{err: err}
	}
}

func (m Model) prCloseCmd(number int) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.PRClose(context.Background(), m.ghRepo, number)
		return ghRefreshMsg{err: err}
	}
}

func (m Model) issueCloseCmd(number int) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.IssueClose(context.Background(), m.ghRepo, number)
		return ghRefreshMsg{err: err}
	}
}

func (m Model) issueReopenCmd(number int) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.IssueReopen(context.Background(), m.ghRepo, number)
		return ghRefreshMsg{err: err}
	}
}

func (m Model) issueCommentCmd(number int, body string) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.IssueComment(context.Background(), m.ghRepo, number, body)
		return ghRefreshMsg{err: err}
	}
}

func (m Model) runRerunCmd(id int64) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.RunRerun(context.Background(), m.ghRepo, id)
		return ghRefreshMsg{err: err}
	}
}

func (m Model) runCancelCmd(id int64) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.RunCancel(context.Background(), m.ghRepo, id)
		return ghRefreshMsg{err: err}
	}
}

func (m Model) prWebCmd(number int) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.PRWeb(context.Background(), m.ghRepo, number)
		return ghRefreshMsg{err: err}
	}
}

func (m Model) issueWebCmd(number int) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.IssueWeb(context.Background(), m.ghRepo, number)
		return ghRefreshMsg{err: err}
	}
}

func (m Model) runWebCmd(id int64) tea.Cmd {
	return func() tea.Msg {
		err := m.gh.RunWeb(context.Background(), m.ghRepo, id)
		return ghRefreshMsg{err: err}
	}
}
