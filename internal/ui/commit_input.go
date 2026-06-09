package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CommitInput is a single-line text prompt for entering a commit message.
type CommitInput struct {
	input  textinput.Model
	active bool
}

// NewCommitInput constructs a closed CommitInput.
func NewCommitInput() CommitInput {
	ti := textinput.New()
	ti.Placeholder = "commit message"
	ti.CharLimit = 200
	ti.Prompt = "> "

	return CommitInput{input: ti}
}

// Open clears any previous value and focuses the input.
func (c *CommitInput) Open() {
	c.active = true
	c.input.SetValue("")
	c.input.Focus()
}

// Close blurs and hides the input.
func (c *CommitInput) Close() {
	c.active = false
	c.input.Blur()
}

// Active reports whether the input is currently open.
func (c CommitInput) Active() bool {
	return c.active
}

// Value returns the current text.
func (c CommitInput) Value() string {
	return c.input.Value()
}

// Update forwards msg to the underlying text input.
func (c *CommitInput) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)
	return cmd
}

// View renders the input as a focused panel.
func (c CommitInput) View(width int) string {
	title := panelTitleStyle.Render("Commit message")
	content := lipgloss.JoinVertical(lipgloss.Left, title, c.input.View())

	return panelFocusedStyle.Width(width).Render(content)
}
