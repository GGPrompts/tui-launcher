package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"tui-launcher/tabs/launch"
)

// tab_routing.go - Tab Routing Implementation
// Handles tab switching and routing events to active tab

// --- Tab Routing Update Logic ---

// handleTabSwitch processes tab switching keys (1/2/3, Tab, Shift+Tab)
// Returns true if a tab switch occurred, false otherwise
func handleTabSwitch(m *unifiedModel, key string) bool {
	switch key {
	case "1":
		if m.currentTab != tabLaunch {
			m.currentTab = tabLaunch
			m.statusMsg = "Switched to Launch tab"
			return true
		}
	case "2":
		if m.currentTab != tabSessions {
			m.currentTab = tabSessions
			m.statusMsg = "Switched to Sessions tab"
			return true
		}
	case "3":
		if m.currentTab != tabTemplates {
			m.currentTab = tabTemplates
			m.statusMsg = "Switched to Templates tab"
			return true
		}
	case "tab":
		m.currentTab = nextTab(m.currentTab)
		m.statusMsg = fmt.Sprintf("Switched to %s tab", m.currentTab.String())
		return true
	case "shift+tab":
		m.currentTab = prevTab(m.currentTab)
		m.statusMsg = fmt.Sprintf("Switched to %s tab", m.currentTab.String())
		return true
	}
	return false
}

// routeUpdateToTab routes a message to the active tab's update function
// Phase 1: Only Launch tab is implemented
func routeUpdateToTab(m *unifiedModel, msg tea.Msg) tea.Cmd {
	switch m.currentTab {
	case tabLaunch:
		// Route to launch tab update (will be implemented)
		return routeToLaunchTab(m, msg)

	case tabSessions:
		// Phase 1: Placeholder - just show message
		m.statusMsg = "Sessions tab - coming soon (Phase 1)"
		return nil

	case tabTemplates:
		// Phase 1: Placeholder - just show message
		m.statusMsg = "Templates tab - coming soon (Phase 1)"
		return nil

	default:
		return nil
	}
}

// routeToLaunchTab routes messages to the launch tab
func routeToLaunchTab(m *unifiedModel, msg tea.Msg) tea.Cmd {
	// Type assert to get the actual launch.Model
	launchModel, ok := m.launchModel.(launch.Model)
	if !ok {
		m.statusMsg = "Error: Launch model type mismatch"
		return nil
	}

	// Call the launch tab's Update method
	updatedModel, cmd := launchModel.Update(msg)

	// Store the updated model back
	m.launchModel = updatedModel

	return cmd
}

// --- Tab Rendering ---

// renderTabBar renders the tab indicator bar at the top
func renderTabBar(currentTab tabName, width int) string {
	var result strings.Builder

	// Render tabs with simple text markers
	// Using > < to indicate active tab
	var launchTab, sessionsTab, templatesTab string

	if currentTab == tabLaunch {
		launchTab = "> 1. Launch <"
		sessionsTab = "  2. Sessions  "
		templatesTab = "  3. Templates  "
	} else if currentTab == tabSessions {
		launchTab = "  1. Launch  "
		sessionsTab = "> 2. Sessions <"
		templatesTab = "  3. Templates  "
	} else if currentTab == tabTemplates {
		launchTab = "  1. Launch  "
		sessionsTab = "  2. Sessions  "
		templatesTab = "> 3. Templates <"
	} else {
		// Default to launch
		launchTab = "> 1. Launch <"
		sessionsTab = "  2. Sessions  "
		templatesTab = "  3. Templates  "
	}

	// Combine tabs
	result.WriteString(launchTab)
	result.WriteString(sessionsTab)
	result.WriteString(templatesTab)

	// Add navigation hint
	result.WriteString("  (Tab/Shift+Tab to cycle, 1/2/3 for direct access)")
	result.WriteString("\n")

	// Add separator line
	separator := strings.Repeat("─", width)
	result.WriteString(separator)
	result.WriteString("\n")

	return result.String()
}

// renderActiveTabContent routes rendering to the active tab
func renderActiveTabContent(m unifiedModel) string {
	switch m.currentTab {
	case tabLaunch:
		// Type assert and render launch tab
		launchModel, ok := m.launchModel.(launch.Model)
		if !ok {
			return "Error: Launch model type mismatch"
		}
		return launchModel.View()

	case tabSessions:
		return renderSessionsTabPlaceholder(m)

	case tabTemplates:
		return renderTemplatesTabPlaceholder(m)

	default:
		return "Unknown tab"
	}
}

// --- Placeholder Tab Renderers (Phase 1) ---

func renderLaunchTabPlaceholder(m unifiedModel) string {
	var content strings.Builder

	content.WriteString("Launch Tab (Current TUI Launcher)\n\n")
	content.WriteString("This tab will show:\n")
	content.WriteString("  • Hierarchical project tree\n")
	content.WriteString("  • Global tools, AI commands, scripts\n")
	content.WriteString("  • Multi-select command launching\n")
	content.WriteString("  • Quick CD into project directories\n\n")

	if m.statusMsg != "" {
		content.WriteString(fmt.Sprintf("Status: %s\n", m.statusMsg))
	}

	content.WriteString("\nPhase 1: Tab routing implemented ✓\n")
	content.WriteString("Next: Integrate existing tui-launcher view\n")

	return content.String()
}

func renderSessionsTabPlaceholder(m unifiedModel) string {
	var content strings.Builder

	content.WriteString("Sessions Tab (From Tmuxplexer)\n\n")
	content.WriteString("This tab will show:\n")
	content.WriteString("  • Active tmux sessions\n")
	content.WriteString("  • Live session previews\n")
	content.WriteString("  • Session management (attach, kill, rename)\n")
	content.WriteString("  • Claude Code status tracking\n\n")

	if m.statusMsg != "" {
		content.WriteString(fmt.Sprintf("Status: %s\n", m.statusMsg))
	}

	content.WriteString("\nPhase 1: Coming soon\n")

	return content.String()
}

func renderTemplatesTabPlaceholder(m unifiedModel) string {
	var content strings.Builder

	content.WriteString("Templates Tab (From Tmuxplexer)\n\n")
	content.WriteString("This tab will show:\n")
	content.WriteString("  • Categorized workspace templates\n")
	content.WriteString("  • Template creation wizard\n")
	content.WriteString("  • Quick session creation from layouts\n")
	content.WriteString("  • Template preview and editing\n\n")

	if m.statusMsg != "" {
		content.WriteString(fmt.Sprintf("Status: %s\n", m.statusMsg))
	}

	content.WriteString("\nPhase 1: Coming soon\n")

	return content.String()
}
