package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// model.go - Model Management
// Purpose: Model initialization and layout calculations
// When to extend: Add new initialization logic or layout calculation functions here

// initialModel creates the initial application state
func initialModel(cfg Config, popupMode bool) model {
	// Load templates
	templates, err := loadTemplates()
	if err != nil {
		templates = []SessionTemplate{} // Use empty list on error
	}

	// Load sessions
	sessions, err := listSessions()
	if err != nil {
		sessions = []TmuxSession{}
	}

	// Get current session name for auto-selection
	currentSessionName := getCurrentSessionName()

	// Auto-select current session (if we're in tmux)
	selectedSession := 0
	if currentSessionName != "" {
		for i, session := range sessions {
			if session.Name == currentSessionName {
				selectedSession = i
				break
			}
		}
	}

	m := model{
		config:             cfg,
		width:              0,
		height:             0,
		focusedComponent:   "main",
		statusMsg:          "",
		popupMode:          popupMode,
		currentSessionName: currentSessionName,

		// Focus state initialization (default to sessions list)
		focusState: FocusSessions,
		lastUpperPanelFocus: FocusSessions, // Track last upper panel for adaptive sizing

		// Adaptive mode initialization (default to enabled)
		adaptiveMode: true,

		// Tab state initialization (default to sessions tab)
		sessionsTab: "sessions",

		// Template and session data
		templates:        templates,
		sessions:         sessions,
		selectedSession:  selectedSession,
		selectedTemplate: 0,
		expandedCategories: map[string]bool{
			"Projects": true, // Auto-expand Projects category by default
		},
		templateTreeItems: []TemplateTreeItem{},

		// Session tree initialization
		expandedSessions: map[string]bool{},
		sessionTreeItems: []SessionTreeItem{},
		sessionFilter:    FilterAll, // Default to showing all sessions

		// Command mode initialization
		commandInput:   "",
		commandCursor:  0,
		commandHistory: []string{},
		historyIndex:   -1,

		// Placeholder content (will be replaced with real tmux data)
		sessionsContent: []string{
			"No sessions loaded yet",
			"",
			"Press 'n' to create a new session",
			"",
			"Navigation:",
			"  ↑/k - Move up",
			"  ↓/j - Move down",
			"  Enter - Attach to session",
		},
		previewContent: []string{
			"Preview will appear here",
			"",
			"This panel shows the live content",
			"of the active pane in the selected",
			"session.",
			"",
			"Updates automatically when you",
			"select a different session.",
		},
		commandContent: []string{
			"> _",
			"Type a command to send to the selected session",
		},
		templatesContent: []string{
			"Templates (for future use)",
		},
	}

	// Update all panel content
	m.updateSessionsContent()
	m.updatePreviewContent()
	m.updateCommandContent()
	// Templates content updated on-demand when Templates tab is shown

	// Set initial context-aware status message
	m.statusMsg = m.getContextualStatusMessage()

	return m
}

// setSize updates the model dimensions and recalculates layouts
func (m *model) setSize(width, height int) {
	m.width = width
	m.height = height

	// Recalculate all panel content with new dimensions
	// This is critical for popup mode where initial size is 0x0
	m.updateSessionsContent()
	m.updatePreviewContent()
	m.updateCommandContent()
	m.updateTemplatesContent()
}

// calculateLayout computes layout dimensions based on config
func (m model) calculateLayout() (int, int) {
	contentWidth := m.width
	contentHeight := m.height

	// Adjust for UI elements
	if m.config.UI.ShowTitle {
		contentHeight -= 1 // title bar height (1 line)
	}
	if m.config.UI.ShowStatus {
		contentHeight -= 2 // status bar height (2 lines: status + help)
	}

	// CRITICAL: Account for panel borders (fixes overflow issue)
	contentHeight -= 2 // top + bottom borders

	return contentWidth, contentHeight
}

// calculateDualPaneLayout computes left and right pane widths
func (m model) calculateDualPaneLayout() (int, int) {
	contentWidth, _ := m.calculateLayout()

	dividerWidth := 0
	if m.config.Layout.ShowDivider {
		dividerWidth = 1
	}

	leftWidth := int(float64(contentWidth-dividerWidth) * m.config.Layout.SplitRatio)
	rightWidth := contentWidth - leftWidth - dividerWidth

	return leftWidth, rightWidth
}

// Helper functions for common operations

// visualWidth calculates the visual width of a string, accounting for ANSI codes and emojis
func visualWidth(s string) int {
	// Strip ANSI codes first
	stripped := ""
	inAnsi := false

	for _, ch := range s {
		// Detect start of ANSI escape sequence
		if ch == '\033' {
			inAnsi = true
			continue
		}

		// Skip characters inside ANSI sequences
		if inAnsi {
			// ANSI sequences end with a letter (A-Z, a-z)
			if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') {
				inAnsi = false
			}
			continue
		}

		stripped += string(ch)
	}

	// Strip variation selectors to work around go-runewidth bug
	// VS incorrectly reports width=1 instead of width=0
	stripped = strings.ReplaceAll(stripped, "\uFE0F", "") // VS-16 (emoji presentation)
	stripped = strings.ReplaceAll(stripped, "\uFE0E", "") // VS-15 (text presentation)

	// Now use StringWidth on the whole stripped string
	return runewidth.StringWidth(stripped)
}

// padToVisualWidth pads a string to a target visual width
func padToVisualWidth(s string, targetWidth int) string {
	currentWidth := visualWidth(s)

	if currentWidth >= targetWidth {
		return s
	}

	padding := targetWidth - currentWidth
	return s + strings.Repeat(" ", padding)
}

// getContentArea returns the available content area dimensions
func (m model) getContentArea() (width, height int) {
	return m.calculateLayout()
}

// isValidSize checks if the terminal size is sufficient
// Lowered minimum for tabbed layout (was 80x24, now 60x15)
func (m model) isValidSize() bool {
	return m.width >= 60 && m.height >= 15
}

// updateSessionsContent updates the left panel with templates and sessions
func (m *model) updateSessionsContent() {
	var lines []string

	// Add session statistics at the top with filter indicator
	attachedCount := 0
	detachedCount := 0
	aiToolCounts := make(map[string]int)

	for _, session := range m.sessions {
		if session.Attached {
			attachedCount++
		} else {
			detachedCount++
		}
		if session.AITool != "" {
			aiToolCounts[session.AITool]++
		}
	}

	// Stats line with filter indicator
	filterText := ""
	switch m.sessionFilter {
	case FilterAll:
		filterText = "All"
	case FilterAI:
		filterText = "AI Only"
	case FilterAttached:
		filterText = "Attached"
	case FilterDetached:
		filterText = "Detached"
	}

	statsLine := fmt.Sprintf("📊 %s | 🎯 %d showing", filterText, len(m.sessions))
	lines = append(lines, statsLine)
	lines = append(lines, fmt.Sprintf("[f] Filter | ● %d attached | ○ %d detached | 🤖 %d AI",
		attachedCount, detachedCount, aiToolCounts["claude"]+aiToolCounts["codex"]+aiToolCounts["gemini"]))
	lines = append(lines, "DIVIDER")
	lines = append(lines, "")

	if len(m.sessions) == 0 {
		lines = append(lines, "  (no sessions match filter)")
		lines = append(lines, "")
		lines = append(lines, "Press 'f' to change filter")
		m.sessionsContent = lines
		return
	}

	// Use table view for AI filter, tree view for everything else
	if m.sessionFilter == FilterAI {
		// Render compact table view for AI sessions
		tableLines := m.renderSessionTableView()
		lines = append(lines, tableLines...)
		m.sessionsContent = lines
		return
	}

	// Update tree items with current filter
	m.updateSessionTreeItems()

	// Render tree view
	for i, item := range m.sessionTreeItems {
		var line string

		// Build indentation with tree characters
		indent := "  " // Base padding

		// Draw vertical lines for parent levels
		for j := 0; j < item.Depth; j++ {
			if j < len(item.ParentLasts) && !item.ParentLasts[j] {
				indent += "│  "
			} else {
				indent += "   "
			}
		}

		// Selection indicator
		selected := i == m.selectedSession

		if item.Type == "session" {
			// Check if session is "simple" (1 window, 1 pane) - don't show expansion indicator
			isSimpleSession := false
			if item.Session.Windows == 1 {
				// Check if that window has only 1 pane
				windows, err := listWindows(item.Session.Name)
				if err == nil && len(windows) == 1 && windows[0].Panes == 1 {
					isSimpleSession = true
				}
			}

			// Session: show expansion indicator only if not a simple session
			expansionIndicator := ""
			if !isSimpleSession {
				expansionIndicator = "▶ "
				if m.expandedSessions[item.Session.Name] {
					expansionIndicator = "▼ "
				}
			}

			icon := "○"
			if item.Session.Attached {
				icon = "●"
			}

			// Mark current session
			currentMarker := ""
			if item.Session.Name == m.currentSessionName {
				currentMarker = " ◆"
			}

			// AI tool indicator
			toolIcon := ""
			if item.Session.AITool != "" {
				switch item.Session.AITool {
				case "claude":
					toolIcon = " 🤖"
				case "codex":
					toolIcon = " 🔮"
				case "gemini":
					toolIcon = " ✨"
				}
			}

			// Claude status
			statusSuffix := ""
			if item.Session.ClaudeState != nil {
				statusIcon := getClaudeStatusIcon(item.Session.ClaudeState)
				statusSuffix = fmt.Sprintf(" %s", statusIcon)
			}

			line = fmt.Sprintf("%s%s%s %s%s%s", indent, expansionIndicator, icon, item.Session.Name, currentMarker, toolIcon+statusSuffix)

			// Apply styling tags
			if item.Session.Name == m.currentSessionName {
				line = "CURRENT:" + line
			} else if item.Session.ClaudeState != nil {
				line = "CLAUDE:" + line
			}

			// Add selection tag
			if selected {
				line = "SELECTED:" + line
			}

			// Add the session name line
			lines = append(lines, line)

			// Show directory and git branch on a second line (with same indentation + 2 spaces)
			if item.Session.WorkingDir != "" {
				// Shorten home directory to ~
				displayDir := item.Session.WorkingDir
				if home := os.Getenv("HOME"); home != "" {
					displayDir = strings.Replace(displayDir, home, "~", 1)
				}

				detailIndent := indent + "  " // Match base indent + 2 spaces
				if item.Session.GitBranch != "" {
					lines = append(lines, fmt.Sprintf("%s📁 %s  %s", detailIndent, displayDir, item.Session.GitBranch))
				} else {
					lines = append(lines, fmt.Sprintf("%s📁 %s", detailIndent, displayDir))
				}
			}

			// Show Claude status on a third line for Claude sessions
			if item.Session.ClaudeState != nil {
				statusText := formatClaudeStatus(item.Session.ClaudeState)
				detailIndent := indent + "  " // Match base indent + 2 spaces
				lines = append(lines, fmt.Sprintf("%s%s", detailIndent, statusText))
			}

			// Skip adding line here since we already added it above
			continue
		} else if item.Type == "window" {
			// Window: show tree connector
			connector := "├─ "
			if item.IsLast {
				connector = "└─ "
			}

			activeMarker := ""
			if item.Window.Active {
				activeMarker = " ●"
			}

			line = fmt.Sprintf("%s%s%s%s", indent, connector, item.Name, activeMarker)

			if selected {
				line = "SELECTED:" + line
			}
		} else if item.Type == "pane" {
			// Pane: show tree connector
			connector := "├─ "
			if item.IsLast {
				connector = "└─ "
			}

			activeMarker := ""
			if item.Pane.Active {
				activeMarker = " ●"
			}

			line = fmt.Sprintf("%s%s%s%s", indent, connector, item.Name, activeMarker)

			if selected {
				line = "SELECTED:" + line
			}
		}

		lines = append(lines, line)
	}

	m.sessionsContent = lines
}

// renderSessionTableView renders AI sessions in a compact table format
func (m *model) renderSessionTableView() []string {
	var lines []string

	// Calculate available width (content width minus padding and borders)
	contentWidth, _ := m.calculateLayout()
	availableWidth := contentWidth - 6 // Account for borders and padding

	// Apply session filter to get only filtered sessions
	filteredSessions := []TmuxSession{}
	for i, session := range m.sessions {
		include := false
		switch m.sessionFilter {
		case FilterAll:
			include = true
		case FilterAI:
			include = (session.AITool != "")
		case FilterAttached:
			include = session.Attached
		case FilterDetached:
			include = !session.Attached
		}
		if include {
			filteredSessions = append(filteredSessions, session)
			_ = i // filteredIndices not needed
		}
	}

	// Distribute column widths dynamically (following TFE pattern)
	// Account for spacing between columns (6 columns × 2 spaces = 12 chars) and padding
	usableWidth := availableWidth - 14

	// Fixed-width columns
	sessWidth := 3
	windowsWidth := 7

	// Calculate space used by fixed columns and their spacing
	// St (3) + Name (var) + Win (7) + Dir (var) + Branch (var) + Status (var)
	// Spacing: 5 gaps × 2 spaces = 10 chars
	fixedColumnsWidth := sessWidth + windowsWidth
	spacingWidth := 10

	// Remaining width for dynamic columns
	remainingWidth := usableWidth - fixedColumnsWidth - spacingWidth
	if remainingWidth < 30 {
		remainingWidth = 30 // Minimum safety
	}

	// Distribute dynamic columns proportionally
	// Name: 30%, Directory: 35%, Branch: 15%
	// Status: gets whatever is left (instead of fixed 20%)
	nameWidth := (remainingWidth * 30) / 100
	dirWidth := (remainingWidth * 35) / 100
	branchWidth := (remainingWidth * 15) / 100

	// Status column gets the remaining space
	statusWidth := remainingWidth - nameWidth - dirWidth - branchWidth

	// Ensure minimum widths
	if nameWidth < 10 {
		nameWidth = 10
	}
	if dirWidth < 15 {
		dirWidth = 15
	}
	if branchWidth < 8 {
		branchWidth = 8
	}
	if statusWidth < 10 {
		statusWidth = 10
	}

	// Build header with styling
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("87")). // Bright blue
		PaddingLeft(2)

	// Build header string (use regular spacing, not visual padding)
	headerStr := fmt.Sprintf("%-*s  %-*s  %-*s  %-*s  %-*s  %-*s",
		sessWidth, "St",
		nameWidth, "Name",
		windowsWidth, "Win",
		dirWidth, "Directory",
		branchWidth, "Branch",
		statusWidth, "Status")

	// Apply styling to header
	headerLine := headerStyle.Render(headerStr)
	lines = append(lines, headerLine+"\033[0m") // Reset ANSI codes

	// Render each filtered session as a table row
	for i, session := range filteredSessions {
		// Check if this is the selected item
		selected := false
		if m.selectedSession < len(m.sessionTreeItems) {
			selectedItem := m.sessionTreeItems[m.selectedSession]
			if selectedItem.Type == "session" && selectedItem.Session.Name == session.Name {
				selected = true
			}
		}

		// Column 1: Session state (attached/current)
		sessIcon := "○"
		if session.Attached {
			sessIcon = "●"
		}
		currentMarker := " "
		if session.Name == m.currentSessionName {
			currentMarker = "◆"
		}
		sessCol := fmt.Sprintf("%s%s", sessIcon, currentMarker)

		// Column 2: Name with AI tool icon
		toolIcon := ""
		switch session.AITool {
		case "claude":
			toolIcon = "🤖"
		case "codex":
			toolIcon = "🔮"
		case "gemini":
			toolIcon = "✨"
		}
		displayName := session.Name
		if toolIcon != "" {
			displayName = toolIcon + " " + displayName
		}

		// Truncate name if needed using visual width
		if visualWidth(displayName) > nameWidth {
			// Truncate based on byte length (approximate)
			for visualWidth(displayName) > nameWidth-2 {
				displayName = displayName[:len(displayName)-1]
			}
			displayName = displayName + ".."
		}

		// Column 3: Window count
		windowCount := fmt.Sprintf("%d", session.Windows)

		// Column 4: Directory (shortened)
		displayDir := session.WorkingDir
		if home := os.Getenv("HOME"); home != "" {
			displayDir = strings.Replace(displayDir, home, "~", 1)
		}
		if len(displayDir) > dirWidth {
			displayDir = "..." + displayDir[len(displayDir)-(dirWidth-3):]
		}

		// Column 5: Git branch
		displayBranch := session.GitBranch
		if displayBranch == "" {
			displayBranch = "-"
		}
		if len(displayBranch) > branchWidth {
			displayBranch = displayBranch[:branchWidth-2] + ".."
		}

		// Column 6: Status
		statusText := "-"
		if session.ClaudeState != nil {
			statusText = formatClaudeStatus(session.ClaudeState)
		}
		if visualWidth(statusText) > statusWidth {
			// Truncate status text if needed
			for visualWidth(statusText) > statusWidth-2 && len(statusText) > 0 {
				statusText = statusText[:len(statusText)-1]
			}
			statusText = statusText + ".."
		}

		// Build row using padToVisualWidth for name column (contains emojis)
		paddedName := padToVisualWidth(displayName, nameWidth)
		row := fmt.Sprintf("%-*s  %s  %-*s  %-*s  %-*s  %-*s",
			sessWidth, sessCol,
			paddedName,
			windowsWidth, windowCount,
			dirWidth, displayDir,
			branchWidth, displayBranch,
			statusWidth, statusText)

		// Apply row styling
		var styledRow string
		if selected {
			// Selected row: blue background (overrides everything)
			styledRow = selectedStyle.Render(row)
		} else {
			// Non-selected rows: apply foreground color + optional alternating background
			var rowStyle lipgloss.Style

			// Determine foreground color
			if session.Name == m.currentSessionName {
				// Current session: cyan + bold
				rowStyle = currentSessionStyle.Copy()
			} else if session.ClaudeState != nil {
				// Claude session: orange + bold
				rowStyle = claudeSessionStyle.Copy()
			} else {
				// Regular session: default style
				rowStyle = lipgloss.NewStyle()
			}

			// Add alternating background for easier reading
			if i%2 == 0 {
				rowStyle = rowStyle.Background(lipgloss.AdaptiveColor{Light: "#eeeeee", Dark: "#333333"})
			}

			styledRow = rowStyle.Render(row)
		}

		// Add padding and ANSI reset
		lines = append(lines, "  "+styledRow+"\033[0m")
	}

	return lines
}

// updateTemplatesContent updates the right panel with templates in tree view
func (m *model) updateTemplatesContent() {
	var lines []string

	lines = append(lines, "TEMPLATES")
	lines = append(lines, "")

	if len(m.templates) == 0 {
		lines = append(lines, "  (no templates)")
		lines = append(lines, "")
		lines = append(lines, "Add templates to ~/.config/tmuxplexer/templates.json")
		m.templatesContent = lines
		return
	}

	// Update tree items
	m.updateTemplateTreeItems()

	// Render tree view
	for i, item := range m.templateTreeItems {
		var line string

		// Build indentation with tree characters
		indent := "  " // Base padding

		// Draw vertical lines for parent levels
		for j := 0; j < item.Depth; j++ {
			if j < len(item.ParentLasts) && !item.ParentLasts[j] {
				indent += "│  "
			} else {
				indent += "   "
			}
		}

		// Selection indicator and tree connector
		selected := i == m.selectedTemplate
		prefix := "  "
		if selected {
			prefix = "► "
		}

		if item.Type == "category" {
			// Category: show expansion indicator
			expansionIndicator := "▶ "
			if m.expandedCategories[item.Name] {
				expansionIndicator = "▼ "
			}
			line = fmt.Sprintf("%s%s%s", indent, expansionIndicator, item.Name)

			// Add SELECTED: tag for visual styling
			if selected {
				line = "SELECTED:" + line
			}
		} else {
			// Template: show tree connector
			connector := "├─ "
			if item.IsLast {
				connector = "└─ "
			}
			line = fmt.Sprintf("%s%s%s%s (%s)", indent, connector, prefix, item.Name, item.Template.Layout)

			// Add SELECTED: tag for visual styling
			if selected {
				line = "SELECTED:" + line
			}
		}

		lines = append(lines, line)
	}

	m.templatesContent = lines
}

// updateCommandContent updates the header panel with command input UI
func (m *model) updateCommandContent() {
	var lines []string

	// Determine target from selected tree item
	targetText := "no pane selected"
	if m.selectedSession >= 0 && m.selectedSession < len(m.sessionTreeItems) {
		item := m.sessionTreeItems[m.selectedSession]

		if item.Session != nil {
			// Add AI tool icon if applicable
			toolIcon := ""
			if item.Session.AITool != "" {
				switch item.Session.AITool {
				case "claude":
					toolIcon = "🤖 "
				case "codex":
					toolIcon = "🔮 "
				case "gemini":
					toolIcon = "✨ "
				}
			}

			// Build target description based on item type
			if item.Type == "pane" && item.Pane != nil {
				targetText = fmt.Sprintf("%s%s > Window %d > Pane %d ✓",
					toolIcon, item.Session.Name, item.WindowIndex, item.PaneIndex)
			} else if item.Type == "window" && item.Window != nil {
				targetText = fmt.Sprintf("%s%s > Window %d (select a pane)",
					toolIcon, item.Session.Name, item.WindowIndex)
			} else if item.Type == "session" && item.Session.Windows == 1 {
				// Check if single-pane session
				targetText = fmt.Sprintf("%s%s (single pane) ✓", toolIcon, item.Session.Name)
			} else {
				targetText = fmt.Sprintf("%s%s (expand and select a pane)",
					toolIcon, item.Session.Name)
			}
		}
	}

	// When on Chat tab, show command input
	if m.focusState == FocusCommand {
		lines = append(lines, "Target: "+targetText)

		// Calculate command panel height using the same adaptive calculation as rendering
		contentWidth, contentHeight := m.calculateLayout()

		// Add back the 2 lines that calculateLayout() subtracted for borders
		contentHeight += 2

		// Use the same adaptive height calculation as rendering and mouse handling
		_, _, commandHeight := m.calculateAdaptivePanelHeights(contentHeight)

		// Calculate available lines for command input
		// commandHeight - 2 (borders) - 1 (target line) - 1 (help text)
		availableCommandLines := commandHeight - 4
		if availableCommandLines < 1 {
			availableCommandLines = 1
		}

		// Wrap command input to multiple lines
		maxTextWidth := contentWidth - 2
		if maxTextWidth < 10 {
			maxTextWidth = 10
		}

		// Get wrapped command lines with cursor position
		allCommandLines, cursorLineIdx := wrapCommandInput(m.commandInput, m.commandCursor, maxTextWidth)
		totalCommandLines := len(allCommandLines)

		// Create viewport around cursor to keep it visible
		var visibleCommandLines []string
		showTopIndicator := false
		showBottomIndicator := false

		if totalCommandLines <= availableCommandLines {
			// All lines fit, show them all
			visibleCommandLines = allCommandLines
		} else {
			// Need viewport - reserve space for indicators
			// If we show both indicators, we need 2 lines for them
			// So we have (availableCommandLines - 2) for actual command lines
			maxIndicatorLines := 2 // Top + bottom indicators
			viewportSize := availableCommandLines - maxIndicatorLines
			if viewportSize < 1 {
				viewportSize = 1
			}

			// Center viewport on cursor line
			viewportStart := cursorLineIdx - (viewportSize / 2)
			if viewportStart < 0 {
				viewportStart = 0
			}
			viewportEnd := viewportStart + viewportSize
			if viewportEnd > totalCommandLines {
				viewportEnd = totalCommandLines
				viewportStart = viewportEnd - viewportSize
				if viewportStart < 0 {
					viewportStart = 0
				}
			}

			visibleCommandLines = allCommandLines[viewportStart:viewportEnd]
			showTopIndicator = viewportStart > 0
			showBottomIndicator = viewportEnd < totalCommandLines
		}

		// Add scroll indicators
		if showTopIndicator {
			lines = append(lines, "  ↑ more above...")
		}
		lines = append(lines, visibleCommandLines...)
		if showBottomIndicator {
			lines = append(lines, "  ↓ more below...")
		}

		// Show help text with character count for long commands
		helpText := "[↑↓] History | [Ctrl+V] Paste | [Enter] Send | [Esc] Clear"
		if len([]rune(m.commandInput)) > 100 {
			helpText = fmt.Sprintf("[↑↓] History | [Ctrl+V] Paste | %d chars | [Enter] Send | [Esc] Clear", len([]rune(m.commandInput)))
		}
		lines = append(lines, helpText)
	} else {
		// When not focused, show hint text
		lines = append(lines, "Send commands to any pane")
		lines = append(lines, "Press '3' to focus | Select a pane with ↑↓ first")

		// Show last command if available
		if m.lastCommand != "" {
			lastCmdLine := fmt.Sprintf("Last: %s → %s", m.lastCommand, m.lastCommandTarget)
			if m.lastCommandTime != "" {
				lastCmdLine += " (" + m.lastCommandTime + ")"
			}
			lines = append(lines, lastCmdLine)
		}
	}

	m.commandContent = lines
}

// updatePreviewContent updates the footer panel with live pane preview
func (m *model) updatePreviewContent() {
	var lines []string

	// Calculate max text width for footer panel (account for borders)
	maxTextWidth := m.width - 2
	if maxTextWidth < 1 {
		maxTextWidth = 1
	}

	// If we're on the templates tab, show template details in preview pane
	if m.sessionsTab == "templates" {
		m.updateTemplatePreview()
		return
	}

	// Check if we have a selected tree item
	if len(m.sessionTreeItems) == 0 || m.selectedSession >= len(m.sessionTreeItems) {
		lines = []string{
			"Preview Panel",
			"",
			"Select a session (left panel) to see live pane content",
			"",
			"Features:",
			"  • Live content from panes",
			"  • Navigate tree with ↑/↓",
			"  • Expand sessions with Enter",
		}
		m.previewContent = lines
		return
	}

	// Get the selected tree item
	item := m.sessionTreeItems[m.selectedSession]

	// Get the session from the tree item
	if item.Session == nil {
		m.previewContent = []string{"No session selected"}
		return
	}
	session := *item.Session

	// If this is the collapsed current session, show a helpful message
	// But allow previewing individual windows/panes within the current session
	if item.Type == "session" && session.Name == m.currentSessionName {
		lines = []string{
			fmt.Sprintf("Current Session: %s", session.Name),
			"",
			"You are currently in this session.",
			"",
			"Quick Actions:",
			"  • Press → to expand and view windows/panes",
			"  • Press 's' to save this session as a template",
			"  • Press 'r' to rename this session",
			"  • Press ↑/↓ to select another session",
		}
		m.previewContent = lines
		m.previewBuffer = []string{}
		m.previewTotalLines = 0
		return
	}

	// Determine which pane to preview based on tree item type
	var paneToPreview *TmuxPane
	var windowToPreview *TmuxWindow

	if item.Type == "pane" && item.Pane != nil {
		// Direct pane selection - show this specific pane
		paneToPreview = item.Pane
		windowToPreview = item.Window
	} else if item.Type == "window" && item.Window != nil {
		// Window selection - show active pane in this window
		windowToPreview = item.Window

		// Get panes for this window
		panes, err := listPanes(session.Name, windowToPreview.Index)
		if err == nil && len(panes) > 0 {
			// Find active pane
			for _, pane := range panes {
				if pane.Active {
					paneToPreview = &pane
					break
				}
			}
			// Fallback to first pane if no active pane found
			if paneToPreview == nil {
				paneToPreview = &panes[0]
			}
		}
	} else {
		// Session selection - show active pane in active window
		windows, err := listWindows(session.Name)
		if err != nil {
			lines = []string{
				fmt.Sprintf("Preview: %s", session.Name),
				"",
				"Error getting windows: " + err.Error(),
			}
			m.previewContent = lines
			return
		}
		if len(windows) == 0 {
			lines = []string{
				fmt.Sprintf("Preview: %s", session.Name),
				"",
				"Session has no windows",
			}
			m.previewContent = lines
			return
		}

		// Find active window
		for _, window := range windows {
			if window.Active {
				windowToPreview = &window
				break
			}
		}

		if windowToPreview == nil {
			windowToPreview = &windows[0]
		}

		// Get panes for active window
		panes, err := listPanes(session.Name, windowToPreview.Index)
		if err == nil && len(panes) > 0 {
			for _, pane := range panes {
				if pane.Active {
					paneToPreview = &pane
					break
				}
			}
			if paneToPreview == nil {
				paneToPreview = &panes[0]
			}
		}
	}

	if paneToPreview == nil || windowToPreview == nil {
		lines = []string{
			fmt.Sprintf("Preview: %s", session.Name),
			"",
		}
		if windowToPreview == nil {
			lines = append(lines, "No window found (this shouldn't happen)")
		} else {
			lines = append(lines, fmt.Sprintf("No pane found in window %d", windowToPreview.Index))
		}
		m.previewContent = lines
		return
	}

	// Capture pane content (full scrollback history)
	content, err := capturePane(paneToPreview.ID)
	if err != nil {
		lines = []string{
			fmt.Sprintf("Preview: %s - Window %d: %s - Pane %d (%s)",
				session.Name, windowToPreview.Index, windowToPreview.Name,
				paneToPreview.Index, paneToPreview.ID),
			"",
			"Failed to capture pane content:",
			err.Error(),
			"",
			"Note: tmux can capture any pane regardless of visibility",
			"If you're seeing this, there might be a permissions issue",
		}
		m.previewContent = lines
		m.previewBuffer = []string{}
		m.previewTotalLines = 0
		return
	}

	// Split content into lines and truncate each line to fit panel width
	contentLines := splitLines(content)

	// Truncate each line to fit within the panel borders
	for i, line := range contentLines {
		contentLines[i] = truncateLine(line, maxTextWidth)
	}

	m.previewBuffer = contentLines
	m.previewTotalLines = len(contentLines)

	// Build header for preview (also truncate header lines)
	helpLine := ""
	if m.focusState == FocusPreview {
		// When preview focused, show scrolling controls
		helpLine = fmt.Sprintf("↑↓/PgUp/PgDn/Wheel Scroll | [r] Refresh | Total Lines: %d", m.previewTotalLines)
	} else {
		// When not focused, show navigation hint
		helpLine = fmt.Sprintf("Navigate tree with ↑/↓ | Press [2] to scroll | Total Lines: %d", m.previewTotalLines)
	}

	headerLines := []string{
		truncateLine(fmt.Sprintf("Preview: %s - Window %d: %s - Pane %d: %s",
			session.Name, windowToPreview.Index, windowToPreview.Name,
			paneToPreview.Index, paneToPreview.Command), maxTextWidth),
		truncateLine(helpLine, maxTextWidth),
		"", // separator
	}

	// Calculate how many content lines can fit in preview panel
	_, totalContentHeight := m.calculateLayout()

	// Add back the 2 lines that calculateLayout() subtracted for borders
	totalContentHeight += 2

	// Use the same adaptive height calculation as rendering and mouse handling
	_, previewHeight, _ := m.calculateAdaptivePanelHeights(totalContentHeight)

	// Calculate actual content height (preview height minus header and borders)
	contentHeight := previewHeight - len(headerLines) - 2 // -2 for borders

	// Auto-scroll to bottom for Claude sessions ONCE per session (to see current chat, not empty space at top)
	// Only do this if we haven't already auto-scrolled for this session
	if session.ClaudeState != nil && m.previewScrollOffset == 0 && m.autoScrolledSession != session.Name {
		// Set scroll offset to show the bottom of the buffer
		maxOffset := m.previewTotalLines - contentHeight
		if maxOffset > 0 {
			m.previewScrollOffset = maxOffset
			m.autoScrolledSession = session.Name // Mark this session as auto-scrolled
		}
	}

	// Get scrollable window of content
	startLine := m.previewScrollOffset
	endLine := startLine + contentHeight

	if startLine < 0 {
		startLine = 0
	}
	if endLine > m.previewTotalLines {
		endLine = m.previewTotalLines
	}
	if startLine >= m.previewTotalLines {
		startLine = 0
		m.previewScrollOffset = 0
	}

	// Build footer content with header + visible window of content
	lines = headerLines
	if startLine < endLine {
		lines = append(lines, contentLines[startLine:endLine]...)
	}

	// Add scroll position indicator if scrollable
	if m.previewTotalLines > contentHeight {
		scrollPercent := 0
		if m.previewTotalLines > 0 {
			scrollPercent = (m.previewScrollOffset * 100) / m.previewTotalLines
		}
		lines = append(lines, "")
		lines = append(lines, truncateLine(fmt.Sprintf("── Scroll: %d%% (Line %d-%d of %d) ──",
			scrollPercent, startLine+1, endLine, m.previewTotalLines), maxTextWidth))
	}

	m.previewContent = lines
}

// updateTemplatePreview updates the preview panel with template details when on the templates tab
func (m *model) updateTemplatePreview() {
	var allLines []string

	// Calculate max text width
	maxTextWidth := m.width - 2
	if maxTextWidth < 1 {
		maxTextWidth = 1
	}

	// Check if we have a selected template
	if len(m.templateTreeItems) == 0 || m.selectedTemplate >= len(m.templateTreeItems) {
		allLines = []string{
			"Template Preview",
			"",
			"Select a template to view its details",
			"",
			"Features:",
			"  • View template configuration",
			"  • See pane layout and commands",
			"  • Navigate with ↑/↓",
			"  • Press [2] to focus and scroll this preview",
		}
		m.previewContent = allLines
		m.previewBuffer = allLines
		m.previewTotalLines = len(allLines)
		return
	}

	selectedItem := m.templateTreeItems[m.selectedTemplate]

	// If a category is selected, show category info
	if selectedItem.Type == "category" {
		// Count templates in this category
		templateCount := 0
		for _, item := range m.templateTreeItems {
			if item.Type == "template" && item.Template != nil && item.Template.Category == selectedItem.Name {
				templateCount++
			}
		}

		allLines = []string{
			fmt.Sprintf("Category: %s", selectedItem.Name),
			"",
			fmt.Sprintf("Templates: %d", templateCount),
			"",
			"Actions:",
			"  • Press Enter or → to expand/collapse",
			"  • Press ↑/↓ to navigate templates",
			"  • Press [2] to focus preview",
			"  • Press 'n' to create new template",
		}
		m.previewContent = allLines
		m.previewBuffer = allLines
		m.previewTotalLines = len(allLines)
		return
	}

	// Template is selected - show details
	if selectedItem.Type == "template" && selectedItem.Template != nil {
		template := selectedItem.Template

		allLines = append(allLines, "● TEMPLATE DETAILS")
		allLines = append(allLines, "")
		allLines = append(allLines, "Name: "+template.Name)
		if template.Category != "" {
			allLines = append(allLines, "Category: "+template.Category)
		}
		if template.Description != "" {
			allLines = append(allLines, "Description: "+template.Description)
		}
		allLines = append(allLines, "")
		allLines = append(allLines, "Layout: "+template.Layout)
		allLines = append(allLines, "Working Dir: "+template.WorkingDir)
		allLines = append(allLines, "")
		allLines = append(allLines, fmt.Sprintf("Panes: %d", len(template.Panes)))
		allLines = append(allLines, "")

		// Show pane details
		for i, pane := range template.Panes {
			allLines = append(allLines, fmt.Sprintf("Pane %d:", i+1))
			if pane.Title != "" {
				allLines = append(allLines, "  Title: "+pane.Title)
			}
			if pane.Command != "" {
				allLines = append(allLines, "  Command: "+pane.Command)
			}
			if pane.WorkingDir != "" {
				allLines = append(allLines, "  Dir: "+pane.WorkingDir)
			} else {
				allLines = append(allLines, "  Dir: "+template.WorkingDir+" (default)")
			}
			allLines = append(allLines, "")
		}

		allLines = append(allLines, "")
		allLines = append(allLines, "Actions:")
		allLines = append(allLines, "  • Press Enter to create session from this template")
		allLines = append(allLines, "  • Press 'o' to create and attach immediately")
		allLines = append(allLines, "  • Press 'e' to edit templates in your editor")
		allLines = append(allLines, "  • Press 'd' to delete this template")
		allLines = append(allLines, "  • Press [2] to focus and scroll this preview")
	}

	// Store full content in buffer
	m.previewBuffer = allLines
	m.previewTotalLines = len(allLines)

	// Calculate how many lines can fit in preview panel
	_, totalContentHeight := m.calculateLayout()
	totalContentHeight += 2
	_, previewHeight, _ := m.calculateAdaptivePanelHeights(totalContentHeight)

	// Calculate visible content height (subtract borders and header)
	headerLines := 1 // We'll add a simple header
	contentHeight := previewHeight - headerLines - 2 // -2 for borders

	if contentHeight < 1 {
		contentHeight = 1
	}

	// Calculate scroll window
	startLine := m.previewScrollOffset
	endLine := startLine + contentHeight

	if startLine < 0 {
		startLine = 0
		m.previewScrollOffset = 0
	}
	if endLine > m.previewTotalLines {
		endLine = m.previewTotalLines
	}
	if startLine >= m.previewTotalLines {
		startLine = 0
		m.previewScrollOffset = 0
		endLine = contentHeight
		if endLine > m.previewTotalLines {
			endLine = m.previewTotalLines
		}
	}

	// Build visible content with header
	var displayLines []string

	// Add header with scroll info if scrollable
	if m.previewTotalLines > contentHeight {
		scrollPercent := 0
		if m.previewTotalLines > 0 {
			scrollPercent = (m.previewScrollOffset * 100) / m.previewTotalLines
		}
		header := fmt.Sprintf("Template Details (Scroll: %d%% - Line %d-%d of %d)",
			scrollPercent, startLine+1, endLine, m.previewTotalLines)
		displayLines = append(displayLines, header)
	} else {
		displayLines = append(displayLines, "Template Details")
	}
	displayLines = append(displayLines, "") // separator

	// Add visible window of content
	if startLine < endLine && endLine <= len(allLines) {
		displayLines = append(displayLines, allLines[startLine:endLine]...)
	}

	m.previewContent = displayLines
}

// splitLines splits a string into lines, handling different line ending styles
func splitLines(content string) []string {
	// Replace \r\n with \n, then split on \n
	content = strings.Replace(content, "\r\n", "\n", -1)
	return strings.Split(content, "\n")
}

// truncateLine truncates a line to fit within maxWidth (visual width, accounting for ANSI codes)
func truncateLine(s string, maxWidth int) string {
	// Use lipgloss.Width to properly measure visual width (ignoring ANSI codes)
	currentWidth := lipgloss.Width(s)
	if currentWidth <= maxWidth {
		return s
	}

	// Need to truncate - remove runes from the end until we fit
	runes := []rune(s)
	for len(runes) > 0 && lipgloss.Width(string(runes)) > maxWidth {
		runes = runes[:len(runes)-1]
	}

	return string(runes)
}

// calculateAdaptivePanelHeights computes panel heights for the unified 3-panel layout
// When adaptiveMode is true, panels adapt based on focus state:
// - Sessions focused: 50% / 30% / 20%
// - Preview focused: 30% / 50% / 20%
// - Command focused or default: 40% / 40% / 20%
// When adaptiveMode is false, panels use fixed proportions (40% / 40% / 20%)
// Command panel always stays at 20% (fixed height for typing)
func (m model) calculateAdaptivePanelHeights(availableHeight int) (sessionsHeight, previewHeight, commandHeight int) {
	// Golden Rule #1: Account for borders BEFORE calculating heights
	// Each panel has top+bottom border (2 lines), so subtract 6 total (3 panels × 2)
	innerHeight := availableHeight - 6 // -2 per panel for borders

	if innerHeight < 15 {
		// Minimum viable layout for tiny terminals
		return 5, 5, 5
	}

	// Command panel is always fixed at 20% (minimum 5 lines for input + help text)
	commandHeight = innerHeight * 20 / 100
	if commandHeight < 5 {
		commandHeight = 5
	}

	// Remaining height for sessions and preview
	remaining := innerHeight - commandHeight

	// Check if adaptive mode is enabled
	if !m.adaptiveMode {
		// Fixed mode: always use balanced 40/40 split
		sessionsHeight = remaining / 2
		previewHeight = remaining - sessionsHeight
		return sessionsHeight, previewHeight, commandHeight
	}

	// Adaptive distribution based on focus
	// When command panel is focused, use lastUpperPanelFocus to maintain sizing
	focusForSizing := m.focusState
	if m.focusState == FocusCommand {
		focusForSizing = m.lastUpperPanelFocus
	}

	switch focusForSizing {
	case FocusSessions:
		// Sessions expanded: 50%, Preview compressed: 30%
		sessionsHeight = remaining * 50 / 80 // 50/(50+30)
		previewHeight = remaining - sessionsHeight

	case FocusPreview:
		// Preview expanded: 50%, Sessions compressed: 30%
		previewHeight = remaining * 50 / 80 // 50/(50+30)
		sessionsHeight = remaining - previewHeight

	default:
		// Balanced: 40/40 split (fallback for startup or unknown state)
		sessionsHeight = remaining / 2
		previewHeight = remaining - sessionsHeight
	}

	return sessionsHeight, previewHeight, commandHeight
}

// buildTemplateTreeItems builds a flattened tree structure for categorized templates
func (m *model) buildTemplateTreeItems() []TemplateTreeItem {
	items := []TemplateTreeItem{}
	
	// Group templates by category
	categories := make(map[string][]int) // category -> template indices
	for i, template := range m.templates {
		category := template.Category
		if category == "" {
			category = "Uncategorized"
		}
		categories[category] = append(categories[category], i)
	}
	
	// Sort category names for consistent display
	categoryNames := make([]string, 0, len(categories))
	for name := range categories {
		categoryNames = append(categoryNames, name)
	}
	// Sort alphabetically for consistent ordering
	sort.Strings(categoryNames)

	// Iterate over sorted category names
	for categoryIndex, category := range categoryNames {
		templateIndices := categories[category]
		isLastCategory := categoryIndex == len(categories)-1
		
		// Add category item
		categoryItem := TemplateTreeItem{
			Type:          "category",
			Name:          category,
			Category:      "",
			Template:      nil,
			Depth:         0,
			IsLast:        isLastCategory,
			ParentLasts:   []bool{},
			TemplateIndex: -1,
		}
		items = append(items, categoryItem)
		
		// If category is expanded, add its templates
		if m.expandedCategories[category] {
			for j, templateIdx := range templateIndices {
				isLastTemplate := j == len(templateIndices)-1
				template := m.templates[templateIdx]
				
				templateItem := TemplateTreeItem{
					Type:          "template",
					Name:          template.Name,
					Category:      category,
					Template:      &template,
					Depth:         1,
					IsLast:        isLastTemplate,
					ParentLasts:   []bool{isLastCategory},
					TemplateIndex: templateIdx,
				}
				items = append(items, templateItem)
			}
		}
	}

	return items
}

// updateTemplateTreeItems rebuilds the tree items cache
func (m *model) updateTemplateTreeItems() {
	m.templateTreeItems = m.buildTemplateTreeItems()
}

// getDefaultCategories returns a list of default category names
func getDefaultCategories() []string {
	return []string{
		"Projects",
		"Agents",
		"Tools",
		"Custom",
	}
}

// buildSessionTreeItems builds a flattened tree from sessions with filters applied
func (m *model) buildSessionTreeItems() []SessionTreeItem {
	items := []SessionTreeItem{}

	// Apply session filter
	filteredSessions := []int{} // indices of sessions that pass the filter
	for i, session := range m.sessions {
		include := false
		switch m.sessionFilter {
		case FilterAll:
			include = true
		case FilterAI:
			include = (session.AITool != "")
		case FilterAttached:
			include = session.Attached
		case FilterDetached:
			include = !session.Attached
		}
		if include {
			filteredSessions = append(filteredSessions, i)
		}
	}

	// Build tree from filtered sessions
	for sessionIdx, sessionIndex := range filteredSessions {
		session := m.sessions[sessionIndex]
		isLastSession := sessionIdx == len(filteredSessions)-1

		// Add session item
		sessionItem := SessionTreeItem{
			Type:         "session",
			Name:         session.Name,
			Session:      &session,
			Window:       nil,
			Pane:         nil,
			Depth:        0,
			IsLast:       isLastSession,
			ParentLasts:  []bool{},
			SessionIndex: sessionIndex,
			WindowIndex:  -1,
			PaneIndex:    -1,
		}
		items = append(items, sessionItem)

		// If session is expanded, add its windows and panes
		if m.expandedSessions[session.Name] {
			// Get windows for this session
			windows, err := listWindows(session.Name)
			if err != nil {
				windows = []TmuxWindow{}
			}

			// Skip tree structure for sessions with only 1 window and 1 pane
			if len(windows) == 1 && windows[0].Panes == 1 {
				continue
			}

			for winIdx, window := range windows {
				isLastWindow := winIdx == len(windows)-1

				// Add window item
				windowItem := SessionTreeItem{
					Type:         "window",
					Name:         fmt.Sprintf("%d: %s (%d panes)", window.Index, window.Name, window.Panes),
					Session:      &session,
					Window:       &window,
					Pane:         nil,
					Depth:        1,
					IsLast:       isLastWindow,
					ParentLasts:  []bool{isLastSession},
					SessionIndex: sessionIndex,
					WindowIndex:  winIdx,
					PaneIndex:    -1,
				}
				items = append(items, windowItem)

				// Get panes for this window
				panes, err := listPanes(session.Name, window.Index)
				if err != nil {
					panes = []TmuxPane{}
				}

				for paneIdx, pane := range panes {
					isLastPane := paneIdx == len(panes)-1

					// Add pane item
					paneItem := SessionTreeItem{
						Type:         "pane",
						Name:         fmt.Sprintf("Pane %d: %s", pane.Index, pane.Command),
						Session:      &session,
						Window:       &window,
						Pane:         &pane,
						Depth:        2,
						IsLast:       isLastPane,
						ParentLasts:  []bool{isLastSession, isLastWindow},
						SessionIndex: sessionIndex,
						WindowIndex:  winIdx,
						PaneIndex:    paneIdx,
					}
					items = append(items, paneItem)
				}
			}
		}
	}

	return items
}

// updateSessionTreeItems rebuilds the session tree items cache
func (m *model) updateSessionTreeItems() {
	m.sessionTreeItems = m.buildSessionTreeItems()
}

// wrapCommandInput wraps the command input to multiple lines with cursor
// Returns array of lines that fit within maxWidth, and the line index where cursor is located
func wrapCommandInput(input string, cursorPos int, maxWidth int) ([]string, int) {
	if maxWidth < 10 {
		maxWidth = 10 // Minimum reasonable width
	}

	// Reserve 1 char for cursor (since cursor at end adds a char)
	// First line has "> " prefix (2 chars) + cursor space (1 char)
	firstLineWidth := maxWidth - 3
	otherLineWidth := maxWidth - 1

	var lines []string
	cursorLineIdx := 0
	runes := []rune(input)
	totalRunes := len(runes)

	// Track position in input
	pos := 0

	// First line
	if pos < totalRunes {
		endPos := pos + firstLineWidth
		if endPos > totalRunes {
			endPos = totalRunes
		}

		line := "> " + string(runes[pos:endPos])

		// Add cursor if it's on this line
		if cursorPos >= pos && cursorPos <= endPos {
			// Insert cursor at correct position
			lineRunes := []rune(line)
			cursorOffset := 2 + (cursorPos - pos) // 2 for "> " prefix
			if cursorOffset == len(lineRunes) {
				line = string(lineRunes) + "█"
			} else {
				line = string(lineRunes[:cursorOffset]) + "█" + string(lineRunes[cursorOffset:])
			}
			cursorLineIdx = len(lines)
		}

		lines = append(lines, line)
		pos = endPos
	}

	// Subsequent lines
	for pos < totalRunes {
		endPos := pos + otherLineWidth
		if endPos > totalRunes {
			endPos = totalRunes
		}

		line := string(runes[pos:endPos])

		// Add cursor if it's on this line
		if cursorPos >= pos && cursorPos < endPos {
			lineRunes := []rune(line)
			cursorOffset := cursorPos - pos
			if cursorOffset == len(lineRunes) {
				line = string(lineRunes) + "█"
			} else {
				line = string(lineRunes[:cursorOffset]) + "█" + string(lineRunes[cursorOffset:])
			}
			cursorLineIdx = len(lines)
		} else if cursorPos == totalRunes && endPos == totalRunes && pos < totalRunes {
			// Cursor at end of input on this line
			line = line + "█"
			cursorLineIdx = len(lines)
		}

		lines = append(lines, line)
		pos = endPos
	}

	// If input is empty, show cursor on first line
	if totalRunes == 0 {
		lines = append(lines, "> █")
		cursorLineIdx = 0
	}

	return lines, cursorLineIdx
}
