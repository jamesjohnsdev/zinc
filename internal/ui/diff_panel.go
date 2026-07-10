package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// DiffPanel renders a file's diff in a scrollable viewport.
type DiffPanel struct {
	vp    viewport.Model
	title string
}

// NewDiffPanel constructs an empty DiffPanel.
func NewDiffPanel() DiffPanel {
	return DiffPanel{vp: viewport.New(0, 0), title: "Diff"}
}

// SetContent replaces the panel's title and diff text, and scrolls back to
// the top.
func (p *DiffPanel) SetContent(title, diff string) {
	p.title = title
	p.vp.SetContent(styleDiff(diff))
	p.vp.GotoTop()
}

// SetSize sets the inner viewport's dimensions from the panel's outer box
// size, accounting for the border and padding renderPanel-style panels add.
// It must be called (typically on tea.WindowSizeMsg) for PageUp/PageDown to
// scroll by the correct amount.
func (p *DiffPanel) SetSize(width, height int) {
	p.vp.Width = max(width-4, 0)
	p.vp.Height = max(height-3, 0)
}

// PageUp scrolls the diff up by a viewport page.
func (p *DiffPanel) PageUp() {
	p.vp.PageUp()
}

// PageDown scrolls the diff down by a viewport page.
func (p *DiffPanel) PageDown() {
	p.vp.PageDown()
}

// View renders the panel at the given size.
func (p DiffPanel) View(width, height int, focused bool) string {
	style := panelStyle
	if focused {
		style = panelFocusedStyle
	}

	title := panelTitleStyle.Render(p.title)
	content := lipgloss.JoinVertical(lipgloss.Left, title, p.vp.View())

	return style.Width(width).Height(height).Render(content)
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
