package ui

// placeholderPanel renders a titled panel for a side panel that doesn't
// have a real data source wired up yet. Each use is replaced by a real
// panel implementation as that panel's git data support lands.
type placeholderPanel struct {
	title string
	note  string
}

func (p placeholderPanel) View(width, height int, focused bool) string {
	style := panelStyle
	if focused {
		style = panelFocusedStyle
	}

	content := panelTitleStyle.Render(p.title) + "\n" + helpDescStyle.Render(p.note)

	return style.Width(width).Height(height).Render(content)
}
