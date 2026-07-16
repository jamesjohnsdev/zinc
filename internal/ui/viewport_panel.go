package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// TextViewport is a titled, scrollable block of text — the shared base for
// panels that show a single long document, such as a diff or a pull
// request/issue body.
type TextViewport struct {
	vp    viewport.Model
	title string
}

// NewTextViewport constructs an empty TextViewport.
func NewTextViewport() TextViewport {
	return TextViewport{vp: viewport.New(0, 0)}
}

// SetContent replaces the panel's title and (already-styled) body text, and
// scrolls back to the top.
func (p *TextViewport) SetContent(title, content string) {
	p.title = title
	p.vp.SetContent(content)
	p.vp.GotoTop()
}

// SetSize sets the inner viewport's dimensions from the content width/
// height that will be passed to View. lipgloss's Width/Height already
// include padding (border is applied after, separately), so only the
// panel's own title line needs to be subtracted here. It must be called
// (typically on tea.WindowSizeMsg) for PageUp/PageDown to scroll by the
// correct amount and for content to fill, rather than overflow, its box.
func (p *TextViewport) SetSize(width, height int) {
	p.vp.Width = max(width-2, 0)   // panelStyle's Padding(0, 1): 1 col each side
	p.vp.Height = max(height-1, 0) // the title line
}

// PageUp scrolls the content up by a viewport page.
func (p *TextViewport) PageUp() {
	p.vp.PageUp()
}

// PageDown scrolls the content down by a viewport page.
func (p *TextViewport) PageDown() {
	p.vp.PageDown()
}

// View renders the panel at the given size.
func (p TextViewport) View(width, height int, focused bool) string {
	style := panelStyle
	if focused {
		style = panelFocusedStyle
	}

	title := panelTitleStyle.Render(p.title)
	content := lipgloss.JoinVertical(lipgloss.Left, title, p.vp.View())

	return style.Width(width).Height(height).Render(content)
}
