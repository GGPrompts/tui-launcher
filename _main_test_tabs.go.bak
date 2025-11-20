package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// main_test_tabs.go - Test Entry Point for Phase 1 Tab Routing
// This is a temporary test file to verify tab routing compiles

func main() {
	// Check for popup mode flag
	popupMode := false
	for _, arg := range os.Args[1:] {
		if arg == "--popup" {
			popupMode = true
		}
	}

	// Create unified model
	m := initialUnifiedModel(popupMode)

	// Create Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
