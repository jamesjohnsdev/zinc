package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/git"
)

// StashPanel lists the repository's stash entries.
type StashPanel struct {
	stashes []git.Stash
	cursor  int
}

// NewStashPanel constructs an empty StashPanel.
func NewStashPanel() StashPanel {
	return StashPanel{}
}

// SetStashes replaces the panel's stash list, clamping the cursor.
func (p *StashPanel) SetStashes(stashes []git.Stash) {
	p.stashes = stashes
	if p.cursor >= len(stashes) {
		p.cursor = max(len(stashes)-1, 0)
	}
}

// Selected returns the stash entry under the cursor, if any.
func (p StashPanel) Selected() (git.Stash, bool) {
	if p.cursor < 0 || p.cursor >= len(p.stashes) {
		return git.Stash{}, false
	}
	return p.stashes[p.cursor], true
}

// CursorUp moves the selection up one stash entry.
func (p *StashPanel) CursorUp() {
	if p.cursor > 0 {
		p.cursor--
	}
}

// CursorDown moves the selection down one stash entry.
func (p *StashPanel) CursorDown() {
	if p.cursor < len(p.stashes)-1 {
		p.cursor++
	}
}

// View renders the panel at the given size.
func (p StashPanel) View(width, height int, focused bool) string {
	style := panelStyle
	if focused {
		style = panelFocusedStyle
	}

	title := panelTitleStyle.Render(fmt.Sprintf("Stash (%d)", len(p.stashes)))

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

	content := lipgloss.JoinVertical(lipgloss.Left, title, strings.Join(lines, "\n"))

	return style.Width(width).Height(height).Render(content)
}
