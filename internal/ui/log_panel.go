package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/git"
)

// LogPanel lists the commit history reachable from HEAD.
type LogPanel struct {
	listCursor
	commits []git.Commit
}

// NewLogPanel constructs an empty LogPanel.
func NewLogPanel() LogPanel {
	return LogPanel{}
}

// SetCommits replaces the panel's commit list, clamping the cursor.
func (p *LogPanel) SetCommits(commits []git.Commit) {
	p.commits = commits
	p.setLength(len(commits))
}

// Selected returns the commit under the cursor, if any.
func (p LogPanel) Selected() (git.Commit, bool) {
	if p.cursor < 0 || p.cursor >= len(p.commits) {
		return git.Commit{}, false
	}
	return p.commits[p.cursor], true
}

// View renders the panel at the given size.
func (p LogPanel) View(width, height int, focused bool) string {
	title := fmt.Sprintf("Commits (%d)", len(p.commits))

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

	return renderPanel(width, height, focused, 3, title, lines, p.cursor)
}
