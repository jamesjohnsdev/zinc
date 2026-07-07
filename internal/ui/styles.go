package ui

import "github.com/charmbracelet/lipgloss"

// Theme colors, loosely inspired by lazygit's default palette.
var (
	colorAccent    = lipgloss.Color("#7AA2F7")
	colorAccentDim = lipgloss.Color("#3B4261")
	colorSubtle    = lipgloss.Color("#565F89")
	colorText      = lipgloss.Color("#C0CAF5")
	colorGood      = lipgloss.Color("#9ECE6A")
	colorWarn      = lipgloss.Color("#E0AF68")
	colorBad       = lipgloss.Color("#F7768E")
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubtle).
			Padding(0, 1)

	panelFocusedStyle = panelStyle.
				BorderForeground(colorAccent)

	panelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	errorToastStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBad)
)
