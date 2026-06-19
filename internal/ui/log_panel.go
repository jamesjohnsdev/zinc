package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/git"
)

// LogPanel lists the commit history reachable from HEAD.
type LogPanel struct {
	commits []git.Commit
	cursor  int
}

// NewLogPanel constructs an empty LogPanel.
func NewLogPanel() LogPanel {
	return LogPanel{}
}

// SetCommits replaces the panel's commit list, clamping the cursor.
func (p *LogPanel) SetCommits(commits []git.Commit) {
	p.commits = commits
	if p.cursor >= len(commits) {
		p.cursor = max(len(commits)-1, 0)
	}
}

// Selected returns the commit under the cursor, if any.
func (p LogPanel) Selected() (git.Commit, bool) {
	if p.cursor < 0 || p.cursor >= len(p.commits) {
		return git.Commit{}, false
	}
	return p.commits[p.cursor], true
}

// CursorUp moves the selection up one commit.
func (p *LogPanel) CursorUp() {
	if p.cursor > 0 {
		p.cursor--
	}
}

// CursorDown moves the selection down one commit.
func (p *LogPanel) CursorDown() {
	if p.cursor < len(p.commits)-1 {
		p.cursor++
	}
}

// View renders the panel at the given size.
func (p LogPanel) View(width, height int, focused bool) string {
	style := panelStyle
	if focused {
		style = panelFocusedStyle
	}

	title := panelTitleStyle.Render(fmt.Sprintf("Commits (%d)", len(p.commits)))

	var lines []string
	if len(p.commits) == 0 {
		lines = append(lines, helpDescStyle.Render("no commits"))
	}
	for i, c := range p.commits {
		cursor := "  "
		if i == p.cursor {
			cursor = "> "
		}
		hash := lipgloss.NewStyle().Foreground(colorAccent).Render(c.Short)
		lines = append(lines, cursor+hash+" "+c.Subject)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, strings.Join(lines, "\n"))

	return style.Width(width).Height(height).Render(content)
}
