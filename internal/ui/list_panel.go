package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// listCursor is the cursor bookkeeping shared by every sidebar list panel.
// Panels embed it to get CursorUp/CursorDown for free; they call setLength
// whenever their backing slice changes so the cursor stays in range.
type listCursor struct {
	cursor int
	length int
}

func (c *listCursor) setLength(n int) {
	c.length = n
	if c.cursor >= n {
		c.cursor = max(n-1, 0)
	}
}

// CursorUp moves the selection up one item.
func (c *listCursor) CursorUp() {
	if c.cursor > 0 {
		c.cursor--
	}
}

// CursorDown moves the selection down one item.
func (c *listCursor) CursorDown() {
	if c.cursor < c.length-1 {
		c.cursor++
	}
}

// renderPanel wraps a titled list of already-formatted lines in the
// standard sidebar panel chrome: a bordered box, dimmed when unfocused and
// accented when focused.
func renderPanel(width, height int, focused bool, title string, lines []string) string {
	style := panelStyle
	if focused {
		style = panelFocusedStyle
	}

	content := lipgloss.JoinVertical(lipgloss.Left, panelTitleStyle.Render(title), strings.Join(lines, "\n"))

	return style.Width(width).Height(height).Render(content)
}
