package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"tui-launcher/tabs/launch"
)

// model_unified.go - Unified Model Implementation
// This will eventually replace model.go once all tabs are integrated

// initialUnifiedModel creates the initial unified application state
func initialUnifiedModel(popupMode bool) unifiedModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	// Tab bar takes 3 lines (tabs + hint + separator)
	tabBarHeight := 3
	initialWidth := 80
	initialHeight := 24

	return unifiedModel{
		width:      initialWidth,
		height:     initialHeight,
		currentTab: tabLaunch, // Start with Launch tab
		err:        nil,
		statusMsg:  "Welcome to TUI Launcher - Phase 1 Tab Integration",
		popupMode:  popupMode,
		spinner:    s,
		loading:    false,

		// Initialize launch tab model with adjusted height (subtract tab bar)
		launchModel: launch.New(initialWidth, initialHeight-tabBarHeight),

		// Sessions tab - placeholder
		// sessionsModel: sessionsTabModel{},

		// Templates tab - placeholder
		// templatesModel: templatesTabModel{},
	}
}

// Init initializes the unified model
func (m unifiedModel) Init() tea.Cmd {
	// Type assert to get the actual launch.Model
	launchModel, ok := m.launchModel.(launch.Model)
	if !ok {
		return m.spinner.Tick
	}

	return tea.Batch(
		m.spinner.Tick,
		launchModel.Init(), // Initialize launch tab (loads config)
	)
}

// Update handles messages and routes to appropriate tab
func (m unifiedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		// Global keys (work in all tabs)
		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		// Handle tab switching first
		if handleTabSwitch(&m, key) {
			return m, nil
		}

		// Route to active tab
		cmd := routeUpdateToTab(&m, msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Create adjusted size message for tabs (subtract tab bar height: 3 lines)
		// Tab bar uses: 1 line for tabs, 1 line for hint, 1 line for separator
		tabBarHeight := 3
		adjustedMsg := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - tabBarHeight,
		}

		// Route adjusted resize to active tab
		cmd := routeUpdateToTab(&m, adjustedMsg)
		return m, cmd

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		// Route all other messages to the active tab
		// This includes tab-specific messages like configLoadedMsg
		cmd := routeUpdateToTab(&m, msg)
		return m, cmd
	}

	return m, nil
}

// View renders the unified UI with tab bar
func (m unifiedModel) View() string {
	// Don't check m.loading here - let each tab handle its own loading state
	// The launch tab will show its own loading spinner

	if m.err != nil {
		return "Error: " + m.err.Error() + "\n\nPress q to quit.\n"
	}

	// Render tab bar - this should ALWAYS show
	tabBar := renderTabBar(m.currentTab, m.width)

	// Render active tab content
	content := renderActiveTabContent(m)

	// Combine - tab bar is always shown first
	return tabBar + content
}

// --- Helper Functions ---

// getLayoutMode determines responsive layout (from current model.go)
// Phase 1: Placeholder - will integrate with launch tab
func (m unifiedModel) getLayoutMode() layoutMode {
	if m.height <= 12 {
		return layoutMobile
	}
	if m.width < 80 {
		return layoutCompact
	}
	return layoutDesktop
}

// Note: isInsideTmux() and detectTerminal() are already defined in model.go
// We don't need to redeclare them here
