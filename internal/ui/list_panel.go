package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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

	// Content wider than the box would otherwise get word-wrapped by
	// lipgloss's own Render, turning one logical line into two physical
	// rows and silently invalidating the height budget windowLines/
	// splitPanelHeights computed on a one-row-per-line assumption.
	// Truncating instead keeps every line to exactly one row.
	innerWidth := max(width-2, 0) // panelStyle's Padding(0, 1): 1 col each side

	visibleRows := max(height-1, 0) // minus the title line
	visible := windowLines(lines, cursor, visibleRows)

	header := numberStyle.Render(fmt.Sprintf("%d", number)) + " " + panelTitleStyle.Render(title)
	header = ansi.Truncate(header, innerWidth, "…")

	truncated := make([]string, len(visible))
	for i, l := range visible {
		truncated[i] = ansi.Truncate(l, innerWidth, "…")
	}

	parts := []string{header}
	if len(truncated) > 0 {
		parts = append(parts, strings.Join(truncated, "\n"))
	}
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)

	return style.Width(width).Height(height).Render(content)
}

// borderRows is the number of rows a bordered panel box adds beyond its
// content: one for the top border, one for the bottom. Panels use
// Padding(0, 1) so padding contributes no additional rows.
const borderRows = 2

// splitPanelHeights divides total content rows among n stacked bordered
// panels, each of which needs its own borderRows on top of its share of
// the content — so naively dividing total by n would render taller than
// total once borders are added. It returns the content height for each of
// the first n-1 panels and for the last (which absorbs the remainder).
func splitPanelHeights(total, n int) (each, last int) {
	content := max(total-borderRows*n, 0)
	each = content / n
	last = content - each*(n-1)
	return each, last
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
