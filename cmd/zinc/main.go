// Command zinc is a terminal UI for managing git repositories and, where
// available, their associated GitHub pull requests and Actions workflows.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jamesjohnsdev/zinc/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "zinc: %v\n", err)
		os.Exit(1)
	}
}
