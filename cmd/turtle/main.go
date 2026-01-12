// ABOUTME: Main entry point for the Turtle terminal tutor
// ABOUTME: Launches the Bubble Tea TUI application

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/2389-research/turtle/internal/tui"
)

func main() {
	// Create the Bubble Tea program
	p := tea.NewProgram(
		tui.NewModel(),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the app
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running turtle: %v\n", err)
		os.Exit(1)
	}
}
