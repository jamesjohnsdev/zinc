package ui

import (
	"fmt"
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
// accented when focused, headed by a number badge (matching that panel's
// jump shortcut) and the title. lipgloss's Height only pads content that's
// shorter than the box, it never crops content that's taller — so with
// more lines than fit, renderPanel windows them around cursor itself
// rather than dumping every line and letting the box overflow.
func renderPanel(width, height int, focused bool, number int, title string, lines []string, cursor int) string {
	style := panelStyle
	numberStyle := panelNumberStyle
	if focused {
		style = panelFocusedStyle
		numberStyle = panelNumberFocusedStyle
	}

	visibleRows := max(height-1, 0) // minus the title line
	visible := windowLines(lines, cursor, visibleRows)

	header := numberStyle.Render(fmt.Sprintf("%d", number)) + " " + panelTitleStyle.Render(title)

	parts := []string{header}
	if len(visible) > 0 {
		parts = append(parts, strings.Join(visible, "\n"))
	}
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)

	return style.Width(width).Height(height).Render(content)
}

// windowLines returns at most n consecutive lines from lines, scrolled so
// that index cursor stays visible (centered where possible).
func windowLines(lines []string, cursor, n int) []string {
	if n <= 0 {
		return nil
	}
	if len(lines) <= n {
		return lines
	}

	start := cursor - n/2
	if start < 0 {
		start = 0
	}
	if start+n > len(lines) {
		start = len(lines) - n
	}

	return lines[start : start+n]
}
