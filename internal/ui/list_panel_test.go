package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestRenderPanelTruncatesLongLines guards against a line wider than the
// panel silently costing an extra physical row via lipgloss's own word
// wrap, which would invalidate the one-row-per-line assumption
// windowLines/splitPanelHeights rely on.
func TestRenderPanelTruncatesLongLines(t *testing.T) {
	longLine := "  " + strings.Repeat("a very long title that will not fit ", 5)

	const height = 6

	out := renderPanel(30, height, false, 1, "Title", []string{longLine}, 0)
	if want := height + borderRows; lipgloss.Height(out) != want {
		t.Errorf("rendered height = %d, want exactly %d (content height + border, not inflated by word-wrap)", lipgloss.Height(out), want)
	}
}

func TestWindowLinesScrollsAroundCursor(t *testing.T) {
	lines := []string{"a", "b", "c", "d", "e"}

	got := windowLines(lines, 4, 2)
	want := []string{"d", "e"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("windowLines(cursor=4, n=2) = %v, want %v", got, want)
	}

	if got := windowLines(lines, 0, 10); len(got) != len(lines) {
		t.Errorf("windowLines with n >= len(lines) should return all lines, got %v", got)
	}
}
