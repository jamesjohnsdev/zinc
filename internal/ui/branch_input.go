package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BranchInput is a single-line prompt for naming a new branch.
type BranchInput struct {
	input  textinput.Model
	active bool
}

// NewBranchInput constructs a closed BranchInput.
func NewBranchInput() BranchInput {
	ti := textinput.New()
	ti.Placeholder = "new branch name"
	ti.CharLimit = 200
	ti.Prompt = "> "

	return BranchInput{input: ti}
}

// Open clears any previous value and focuses the input.
func (b *BranchInput) Open() {
	b.active = true
	b.input.SetValue("")
	b.input.Focus()
}

// Close blurs and hides the input.
func (b *BranchInput) Close() {
	b.active = false
	b.input.Blur()
}

// Active reports whether the input is currently open.
func (b BranchInput) Active() bool {
	return b.active
}

// Value returns the current text.
func (b BranchInput) Value() string {
	return b.input.Value()
}

// Update forwards msg to the underlying text input.
func (b *BranchInput) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	b.input, cmd = b.input.Update(msg)
	return cmd
}

// View renders the input as a focused panel.
func (b BranchInput) View(width int) string {
	title := panelTitleStyle.Render("New branch")
	content := lipgloss.JoinVertical(lipgloss.Left, title, b.input.View())

	return panelFocusedStyle.Width(width).Render(content)
}
