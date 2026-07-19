package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/gh"
)

// IssuesPanel lists open issues.
type IssuesPanel struct {
	listCursor
	issues []gh.Issue
}

// NewIssuesPanel constructs an empty IssuesPanel.
func NewIssuesPanel() IssuesPanel {
	return IssuesPanel{}
}

// SetIssues replaces the panel's issue list, clamping the cursor.
func (p *IssuesPanel) SetIssues(issues []gh.Issue) {
	p.issues = issues
	p.setLength(len(issues))
}

// Selected returns the issue under the cursor, if any.
func (p IssuesPanel) Selected() (gh.Issue, bool) {
	if p.cursor < 0 || p.cursor >= len(p.issues) {
		return gh.Issue{}, false
	}
	return p.issues[p.cursor], true
}

// View renders the panel at the given size.
func (p IssuesPanel) View(width, height int, focused bool) string {
	title := fmt.Sprintf("Issues (%d)", len(p.issues))

	var lines []string
	if len(p.issues) == 0 {
		lines = append(lines, helpDescStyle.Render("no open issues"))
	}
	for i, issue := range p.issues {
		cursor := "  "
		if i == p.cursor {
			cursor = "> "
		}
		lines = append(lines, cursor+formatIssueLine(issue))
	}

	return renderPanel(width, height, focused, 2, title, lines, p.cursor)
}

func formatIssueLine(issue gh.Issue) string {
	num := lipgloss.NewStyle().Foreground(colorAccent).Render(fmt.Sprintf("#%d", issue.Number))

	badge := lipgloss.NewStyle().Bold(true).Foreground(colorGood).Render("open")
	if issue.State == "CLOSED" {
		badge = lipgloss.NewStyle().Bold(true).Foreground(colorBad).Render("closed")
	}

	return num + " " + issue.Title + " " + badge
}
