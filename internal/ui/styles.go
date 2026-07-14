package ui

import "github.com/charmbracelet/lipgloss"

// Theme colors: a Tokyo Night-inspired palette.
var (
	colorBg     = lipgloss.Color("#1A1B26")
	colorBgAlt  = lipgloss.Color("#24283B")
	colorAccent = lipgloss.Color("#7AA2F7")
	colorSubtle = lipgloss.Color("#565F89")
	colorText   = lipgloss.Color("#C0CAF5")
	colorGood   = lipgloss.Color("#9ECE6A")
	colorWarn   = lipgloss.Color("#E0AF68")
	colorBad    = lipgloss.Color("#F7768E")
)

var (
	// headerStyle is the full-width app title bar.
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBg).
			Background(colorAccent)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubtle).
			Padding(0, 1)

	panelFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(colorAccent).
				Padding(0, 1)

	panelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText)

	// panelNumberStyle renders a panel's jump-shortcut digit as a small
	// badge ahead of its title.
	panelNumberStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorBg).
				Background(colorSubtle).
				Padding(0, 1)

	panelNumberFocusedStyle = panelNumberStyle.
				Background(colorAccent)

	// statusBarStyle is the full-width help bar background strip.
	statusBarStyle = lipgloss.NewStyle().
			Background(colorBgAlt)

	// keyStyle renders a single keybinding as a keycap-style chip.
	keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBg).
			Background(colorAccent).
			Padding(0, 1)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Background(colorBgAlt)

	errorToastStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBg).
			Background(colorBad).
			Padding(0, 1)
)
