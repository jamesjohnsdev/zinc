package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// DiffPanel renders a file's diff in a scrollable viewport, with unified
// diff coloring applied on top of the shared TextViewport base.
type DiffPanel struct {
	TextViewport
}

// NewDiffPanel constructs an empty DiffPanel.
func NewDiffPanel() DiffPanel {
	p := DiffPanel{TextViewport: NewTextViewport()}
	p.title = "Diff"
	return p
}

// SetContent replaces the panel's title and diff text, and scrolls back to
// the top.
func (p *DiffPanel) SetContent(title, diff string) {
	p.TextViewport.SetContent(title, styleDiff(diff))
}

// styleDiff applies simple unified-diff coloring: additions green,
// deletions red, hunk headers accented, and file headers dimmed.
func styleDiff(diff string) string {
	if strings.TrimSpace(diff) == "" {
		return helpDescStyle.Render("no changes")
	}

	lines := strings.Split(diff, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			lines[i] = helpDescStyle.Render(line)
		case strings.HasPrefix(line, "+"):
			lines[i] = lipgloss.NewStyle().Foreground(colorGood).Render(line)
		case strings.HasPrefix(line, "-"):
			lines[i] = lipgloss.NewStyle().Foreground(colorBad).Render(line)
		case strings.HasPrefix(line, "@@"):
			lines[i] = lipgloss.NewStyle().Foreground(colorAccent).Render(line)
		}
	}

	return strings.Join(lines, "\n")
}
