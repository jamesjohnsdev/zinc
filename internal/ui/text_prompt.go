package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TextPrompt is a single-line, titled text input shown as a focused
// overlay panel — the shared base for the commit message, new branch
// name, and issue comment prompts.
type TextPrompt struct {
	title  string
	input  textinput.Model
	active bool
}

// NewTextPrompt constructs a closed TextPrompt with the given title and
// placeholder text.
func NewTextPrompt(title, placeholder string) TextPrompt {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 500
	ti.Prompt = "> "

	return TextPrompt{title: title, input: ti}
}

// Open clears any previous value and focuses the input.
func (p *TextPrompt) Open() {
	p.active = true
	p.input.SetValue("")
	p.input.Focus()
}

// Close blurs and hides the input.
func (p *TextPrompt) Close() {
	p.active = false
	p.input.Blur()
}

// Active reports whether the input is currently open.
func (p TextPrompt) Active() bool {
	return p.active
}

// Value returns the current text.
func (p TextPrompt) Value() string {
	return p.input.Value()
}

// Update forwards msg to the underlying text input.
func (p *TextPrompt) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	return cmd
}

// View renders the prompt as a focused panel.
func (p TextPrompt) View(width int) string {
	title := panelTitleStyle.Render(p.title)
	content := lipgloss.JoinVertical(lipgloss.Left, title, p.input.View())

	return panelFocusedStyle.Width(width).Render(content)
}
