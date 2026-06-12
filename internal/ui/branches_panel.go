package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/git"
)

// BranchesPanel lists local and remote-tracking branches.
type BranchesPanel struct {
	branches []git.Branch
	cursor   int
}

// NewBranchesPanel constructs an empty BranchesPanel.
func NewBranchesPanel() BranchesPanel {
	return BranchesPanel{}
}

// SetBranches replaces the panel's branch list, clamping the cursor.
func (p *BranchesPanel) SetBranches(branches []git.Branch) {
	p.branches = branches
	if p.cursor >= len(branches) {
		p.cursor = max(len(branches)-1, 0)
	}
}

// Selected returns the branch under the cursor, if any.
func (p BranchesPanel) Selected() (git.Branch, bool) {
	if p.cursor < 0 || p.cursor >= len(p.branches) {
		return git.Branch{}, false
	}
	return p.branches[p.cursor], true
}

// CursorUp moves the selection up one branch.
func (p *BranchesPanel) CursorUp() {
	if p.cursor > 0 {
		p.cursor--
	}
}

// CursorDown moves the selection down one branch.
func (p *BranchesPanel) CursorDown() {
	if p.cursor < len(p.branches)-1 {
		p.cursor++
	}
}

// View renders the panel at the given size.
func (p BranchesPanel) View(width, height int, focused bool) string {
	style := panelStyle
	if focused {
		style = panelFocusedStyle
	}

	title := panelTitleStyle.Render(fmt.Sprintf("Branches (%d)", len(p.branches)))

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

	content := lipgloss.JoinVertical(lipgloss.Left, title, strings.Join(lines, "\n"))

	return style.Width(width).Height(height).Render(content)
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
