package ui

import "strings"

// helpEntry is a single keybinding shown in the status bar.
type helpEntry struct {
	key  string
	desc string
}

var gitGlobalHelp = []helpEntry{
	{"g", "GitHub screen"},
	{"tab", "switch panel"},
	{"1-4", "jump to panel"},
	{"↑/↓", "navigate"},
	{"r", "refresh"},
	{"q", "quit"},
}

var gitPanelHelp = map[focusIndex][]helpEntry{
	focusFiles: {
		{"space", "stage/unstage"},
		{"a", "stage all"},
		{"c", "commit"},
		{"s", "stash"},
		{"D", "discard"},
		{"^u/^d", "scroll diff"},
	},
	focusBranches: {
		{"enter", "checkout"},
		{"n", "new branch"},
		{"d", "delete"},
	},
	focusCommits: {},
	focusStash: {
		{"p", "pop"},
		{"d", "drop"},
	},
}

var githubGlobalHelp = []helpEntry{
	{"g", "Git screen"},
	{"tab", "switch panel"},
	{"1-3", "jump to panel"},
	{"↑/↓", "navigate"},
	{"^u/^d", "scroll detail"},
	{"o", "open in browser"},
	{"r", "refresh"},
	{"q", "quit"},
}

var githubPanelHelp = map[ghFocusIndex][]helpEntry{
	focusPRs: {
		{"c", "checkout"},
		{"m", "squash merge"},
		{"x", "close"},
	},
	focusIssues: {
		{"x", "close/reopen"},
		{"C", "comment"},
	},
	focusRuns: {
		{"R", "rerun"},
		{"x", "cancel"},
	},
}

var promptHelp = []helpEntry{
	{"enter", "confirm"},
	{"esc", "cancel"},
}

var confirmHelp = []helpEntry{
	{"y", "confirm"},
	{"n/esc", "cancel"},
}

// helpBar renders the status bar's keybinding hints, as keycap-style
// chips, for the model's current mode: an open prompt, an open confirm
// dialog, or the normal per-panel keymap layered under the global one for
// whichever screen is active. It spans the full given width as a solid
// strip, matching the header bar.
func (m Model) helpBar(width int) string {
	var entries []helpEntry

	switch {
	case m.commitInput.Active(), m.branchInput.Active(), m.commentInput.Active():
		entries = promptHelp
	case m.confirm.Active():
		entries = confirmHelp
	case m.screen == screenGitHub:
		entries = append(entries, githubPanelHelp[m.ghFocus]...)
		entries = append(entries, githubGlobalHelp...)
	default:
		entries = append(entries, gitPanelHelp[m.focus]...)
		entries = append(entries, gitGlobalHelp...)
	}

	parts := make([]string, len(entries))
	for i, e := range entries {
		parts[i] = keyStyle.Render(e.key) + helpDescStyle.Render(" "+e.desc+" ")
	}

	return statusBarStyle.Width(width).Render(strings.Join(parts, ""))
}
