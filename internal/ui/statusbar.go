package ui

import "strings"

// helpEntry is a single keybinding shown in the status bar.
type helpEntry struct {
	key  string
	desc string
}

var globalHelp = []helpEntry{
	{"tab", "switch panel"},
	{"↑/↓", "navigate"},
	{"r", "refresh"},
	{"q", "quit"},
}

var panelHelp = map[focusIndex][]helpEntry{
	focusFiles: {
		{"space", "stage/unstage"},
		{"a", "stage all"},
		{"c", "commit"},
		{"s", "stash"},
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

var promptHelp = []helpEntry{
	{"enter", "confirm"},
	{"esc", "cancel"},
}

var confirmHelp = []helpEntry{
	{"y", "confirm"},
	{"n/esc", "cancel"},
}

// helpBar renders the status bar's keybinding hints for the model's
// current mode: an open prompt, an open confirm dialog, or the normal
// per-panel keymap layered under the global one.
func (m Model) helpBar() string {
	var entries []helpEntry

	switch {
	case m.commitInput.Active(), m.branchInput.Active():
		entries = promptHelp
	case m.confirm.Active():
		entries = confirmHelp
	default:
		entries = append(entries, panelHelp[m.focus]...)
		entries = append(entries, globalHelp...)
	}

	parts := make([]string, len(entries))
	for i, e := range entries {
		parts[i] = keyStyle.Render(e.key) + " " + helpDescStyle.Render(e.desc)
	}

	return statusBarStyle.Render(strings.Join(parts, "  "))
}
