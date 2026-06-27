package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/git"
)

// StashPanel lists the repository's stash entries.
type StashPanel struct {
	listCursor
	stashes []git.Stash
}

// NewStashPanel constructs an empty StashPanel.
func NewStashPanel() StashPanel {
	return StashPanel{}
}

// SetStashes replaces the panel's stash list, clamping the cursor.
func (p *StashPanel) SetStashes(stashes []git.Stash) {
	p.stashes = stashes
	p.setLength(len(stashes))
}

// Selected returns the stash entry under the cursor, if any.
func (p StashPanel) Selected() (git.Stash, bool) {
	if p.cursor < 0 || p.cursor >= len(p.stashes) {
		return git.Stash{}, false
	}
	return p.stashes[p.cursor], true
}

// View renders the panel at the given size.
func (p StashPanel) View(width, height int, focused bool) string {
	title := fmt.Sprintf("Stash (%d)", len(p.stashes))

	var lines []string
	if len(p.stashes) == 0 {
		lines = append(lines, helpDescStyle.Render("no stashes"))
	}
	for i, s := range p.stashes {
		cursor := "  "
		if i == p.cursor {
			cursor = "> "
		}
		idx := lipgloss.NewStyle().Foreground(colorAccent).Render(fmt.Sprintf("%d:", s.Index))
		lines = append(lines, cursor+idx+" "+s.Message)
	}

	return renderPanel(width, height, focused, title, lines)
}
