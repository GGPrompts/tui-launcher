package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check for popup mode flag
	popupMode := false
	for _, arg := range os.Args[1:] {
		if arg == "--popup" {
			popupMode = true
		}
	}

	// Create unified model (Phase 1: Tab-based architecture)
	m := initialUnifiedModel(popupMode)

	// Create program with alt screen and mouse support
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion(), tea.WithMouseAllMotion())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
