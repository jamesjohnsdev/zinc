package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TestViewFitsTerminalHeight guards against the layout regressing into
// rendering more rows than the terminal has, which previously left the top
// of the UI scrolled off screen: sidebar panels are bordered boxes, and the
// height budget must account for that border overhead. Sizes below ~80x24
// are below this layout's structural minimum (four stacked bordered
// panels each need at least a title row plus two border rows) and aren't
// covered here, same as most multi-pane TUIs.
func TestViewFitsTerminalHeight(t *testing.T) {
	sizes := []struct{ w, h int }{
		{80, 24},
		{100, 30},
		{120, 40},
		{200, 60},
	}

	for _, sz := range sizes {
		m := NewModel()
		updated, _ := m.Update(tea.WindowSizeMsg{Width: sz.w, Height: sz.h})
		m = updated.(Model)

		view := m.View()
		if got := lipgloss.Height(view); got > sz.h {
			t.Errorf("terminal %dx%d: rendered height %d exceeds terminal height", sz.w, sz.h, got)
		}
	}
}
