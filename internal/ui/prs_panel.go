package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/gh"
)

// PRPanel lists open pull requests.
type PRPanel struct {
	listCursor
	prs []gh.PullRequest
}

// NewPRPanel constructs an empty PRPanel.
func NewPRPanel() PRPanel {
	return PRPanel{}
}

// SetPRs replaces the panel's pull request list, clamping the cursor.
func (p *PRPanel) SetPRs(prs []gh.PullRequest) {
	p.prs = prs
	p.setLength(len(prs))
}

// Selected returns the pull request under the cursor, if any.
func (p PRPanel) Selected() (gh.PullRequest, bool) {
	if p.cursor < 0 || p.cursor >= len(p.prs) {
		return gh.PullRequest{}, false
	}
	return p.prs[p.cursor], true
}

// View renders the panel at the given size.
func (p PRPanel) View(width, height int, focused bool) string {
	title := fmt.Sprintf("Pull Requests (%d)", len(p.prs))

	var lines []string
	if len(p.prs) == 0 {
		lines = append(lines, helpDescStyle.Render("no open pull requests"))
	}
	for i, pr := range p.prs {
		cursor := "  "
		if i == p.cursor {
			cursor = "> "
		}
		lines = append(lines, cursor+formatPRLine(pr))
	}

	return renderPanel(width, height, focused, 1, title, lines, p.cursor)
}

func formatPRLine(pr gh.PullRequest) string {
	num := lipgloss.NewStyle().Foreground(colorAccent).Render(fmt.Sprintf("#%d", pr.Number))
	return num + " " + pr.Title + " " + prStateBadge(pr)
}

func prStateBadge(pr gh.PullRequest) string {
	style := lipgloss.NewStyle().Bold(true)
	switch {
	case pr.IsDraft:
		return style.Foreground(colorSubtle).Render("draft")
	case pr.State == "MERGED":
		return style.Foreground(colorAccent).Render("merged")
	case pr.State == "CLOSED":
		return style.Foreground(colorBad).Render("closed")
	default:
		return style.Foreground(colorGood).Render("open")
	}
}
