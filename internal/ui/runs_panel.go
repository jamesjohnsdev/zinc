package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/jamesjohnsdev/zinc/internal/gh"
)

// RunsPanel lists recent Actions workflow runs.
type RunsPanel struct {
	listCursor
	runs []gh.WorkflowRun
}

// NewRunsPanel constructs an empty RunsPanel.
func NewRunsPanel() RunsPanel {
	return RunsPanel{}
}

// SetRuns replaces the panel's workflow run list, clamping the cursor.
func (p *RunsPanel) SetRuns(runs []gh.WorkflowRun) {
	p.runs = runs
	p.setLength(len(runs))
}

// Selected returns the workflow run under the cursor, if any.
func (p RunsPanel) Selected() (gh.WorkflowRun, bool) {
	if p.cursor < 0 || p.cursor >= len(p.runs) {
		return gh.WorkflowRun{}, false
	}
	return p.runs[p.cursor], true
}

// View renders the panel at the given size.
func (p RunsPanel) View(width, height int, focused bool) string {
	title := fmt.Sprintf("Workflow Runs (%d)", len(p.runs))

	var lines []string
	if len(p.runs) == 0 {
		lines = append(lines, helpDescStyle.Render("no workflow runs"))
	}
	for i, run := range p.runs {
		cursor := "  "
		if i == p.cursor {
			cursor = "> "
		}
		lines = append(lines, cursor+formatRunLine(run))
	}

	return renderPanel(width, height, focused, 3, title, lines, p.cursor)
}

func formatRunLine(run gh.WorkflowRun) string {
	icon, style := runStatusIcon(run)
	branch := lipgloss.NewStyle().Foreground(colorSubtle).Render("(" + run.HeadBranch + ")")
	return style.Render(icon) + " " + run.WorkflowName + " " + branch
}

func runStatusIcon(run gh.WorkflowRun) (string, lipgloss.Style) {
	style := lipgloss.NewStyle().Bold(true)

	if run.Status != "completed" {
		if run.Status == "in_progress" {
			return "●", style.Foreground(colorWarn)
		}
		return "○", style.Foreground(colorSubtle) // queued, waiting, ...
	}

	switch run.Conclusion {
	case "success":
		return "✓", style.Foreground(colorGood)
	case "cancelled", "skipped":
		return "○", style.Foreground(colorSubtle)
	default: // failure, timed_out, action_required, ...
		return "✗", style.Foreground(colorBad)
	}
}
