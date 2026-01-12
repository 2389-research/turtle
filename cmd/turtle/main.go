// ABOUTME: Main entry point for the Turtle terminal tutor
// ABOUTME: Launches the Bubble Tea TUI application

package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/2389-research/turtle/internal/tui"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("turtle %s\n", version)
		os.Exit(0)
	}

	// Create the Bubble Tea program
	p := tea.NewProgram(
		tui.NewMissionTUI(),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the app
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running turtle: %v\n", err)
		os.Exit(1)
	}
}
