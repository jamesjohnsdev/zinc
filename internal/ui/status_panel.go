package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/git"
)

// StatusPanel renders the working tree file list: staged, unstaged,
// unmerged, and untracked paths.
type StatusPanel struct {
	listCursor
	files []git.FileStatus
}

// NewStatusPanel constructs an empty StatusPanel.
func NewStatusPanel() StatusPanel {
	return StatusPanel{}
}

// SetFiles replaces the panel's file list, clamping the cursor if needed.
func (p *StatusPanel) SetFiles(files []git.FileStatus) {
	p.files = files
	p.setLength(len(files))
}

// Selected returns the file under the cursor, if any.
func (p StatusPanel) Selected() (git.FileStatus, bool) {
	if p.cursor < 0 || p.cursor >= len(p.files) {
		return git.FileStatus{}, false
	}
	return p.files[p.cursor], true
}

// View renders the panel at the given size.
func (p StatusPanel) View(width, height int, focused bool) string {
	title := fmt.Sprintf("Files (%d)", len(p.files))

	var lines []string
	if len(p.files) == 0 {
		lines = append(lines, helpDescStyle.Render("clean working tree"))
	}
	for i, f := range p.files {
		cursor := "  "
		if i == p.cursor {
			cursor = "> "
		}
		lines = append(lines, cursor+formatFileLine(f))
	}

	return renderPanel(width, height, focused, 1, title, lines, p.cursor)
}

func formatFileLine(f git.FileStatus) string {
	code, style := statusIndicator(f)
	return style.Render(code) + " " + f.Path
}

func statusIndicator(f git.FileStatus) (string, lipgloss.Style) {
	bold := lipgloss.NewStyle().Bold(true)
	switch {
	case f.Unmerged:
		return "UU", bold.Foreground(colorBad)
	case f.Untracked:
		return "??", bold.Foreground(colorWarn)
	case f.Staged != '.' && f.Unstaged != '.':
		return string([]byte{f.Staged, f.Unstaged}), bold.Foreground(colorWarn)
	case f.Staged != '.':
		return string([]byte{f.Staged, ' '}), bold.Foreground(colorGood)
	default:
		return string([]byte{' ', f.Unstaged}), bold.Foreground(colorBad)
	}
}
