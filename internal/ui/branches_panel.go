package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/git"
)

// BranchesPanel lists local and remote-tracking branches.
type BranchesPanel struct {
	listCursor
	branches []git.Branch
}

// NewBranchesPanel constructs an empty BranchesPanel.
func NewBranchesPanel() BranchesPanel {
	return BranchesPanel{}
}

// SetBranches replaces the panel's branch list, clamping the cursor.
func (p *BranchesPanel) SetBranches(branches []git.Branch) {
	p.branches = branches
	p.setLength(len(branches))
}

// Selected returns the branch under the cursor, if any.
func (p BranchesPanel) Selected() (git.Branch, bool) {
	if p.cursor < 0 || p.cursor >= len(p.branches) {
		return git.Branch{}, false
	}
	return p.branches[p.cursor], true
}

// Current returns the name of the currently checked-out branch, if known.
func (p BranchesPanel) Current() (string, bool) {
	for _, b := range p.branches {
		if b.Current {
			return b.Name, true
		}
	}
	return "", false
}

// View renders the panel at the given size.
func (p BranchesPanel) View(width, height int, focused bool) string {
	title := fmt.Sprintf("Branches (%d)", len(p.branches))

	var lines []string
	if len(p.branches) == 0 {
		lines = append(lines, helpDescStyle.Render("no branches"))
	}
	for i, b := range p.branches {
		cursor := "  "
		if i == p.cursor {
			cursor = "> "
		}
		lines = append(lines, cursor+formatBranchLine(b))
	}

	return renderPanel(width, height, focused, 2, title, lines, p.cursor)
}

func formatBranchLine(b git.Branch) string {
	marker := "  "
	nameStyle := lipgloss.NewStyle()
	switch {
	case b.Current:
		marker = "* "
		nameStyle = nameStyle.Foreground(colorGood).Bold(true)
	case b.Remote:
		nameStyle = nameStyle.Foreground(colorSubtle)
	}

	line := marker + nameStyle.Render(b.Name)

	var track []string
	if b.Ahead > 0 {
		track = append(track, fmt.Sprintf("↑%d", b.Ahead))
	}
	if b.Behind > 0 {
		track = append(track, fmt.Sprintf("↓%d", b.Behind))
	}
	if len(track) > 0 {
		line += " " + lipgloss.NewStyle().Foreground(colorWarn).Render(strings.Join(track, " "))
	}

	return line
}
