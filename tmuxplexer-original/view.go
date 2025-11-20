package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// view.go - View Rendering
// Purpose: Top-level view rendering and layout
// When to extend: Add new view modes or modify layout logic

// View renders the entire application
func (m model) View() string {
	// Check if terminal size is sufficient
	if !m.isValidSize() {
		return m.renderMinimalView()
	}

	// Handle errors
	if m.err != nil {
		return m.renderErrorView()
	}

	// Render unified 3-panel layout (sessions + preview + command)
	return m.renderUnifiedView()
}

// renderSinglePane renders a single-pane layout
func (m model) renderSinglePane() string {
	var sections []string

	// Title bar
	if m.config.UI.ShowTitle {
		sections = append(sections, m.renderTitleBar())
	}

	// Main content
	sections = append(sections, m.renderMainContent())

	// Status bar
	if m.config.UI.ShowStatus {
		sections = append(sections, m.renderStatusBar())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderDualPane renders a dual-pane layout (side-by-side)
func (m model) renderDualPane() string {
	var sections []string

	// Title bar
	if m.config.UI.ShowTitle {
		sections = append(sections, m.renderTitleBar())
	}

	// Calculate pane dimensions
	leftWidth, rightWidth := m.calculateDualPaneLayout()

	// Left pane
	leftPane := m.renderLeftPane(leftWidth)

	// Divider
	divider := ""
	if m.config.Layout.ShowDivider {
		divider = m.renderDivider()
	}

	// Right pane
	rightPane := m.renderRightPane(rightWidth)

	// Join panes horizontally
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, divider, rightPane)
	sections = append(sections, panes)

	// Status bar
	if m.config.UI.ShowStatus {
		sections = append(sections, m.renderStatusBar())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderMultiPanel renders a multi-panel layout
func (m model) renderMultiPanel() string {
	// Implement multi-panel layout
	// This is a placeholder - customize based on your needs
	return m.renderSinglePane()
}


// Component rendering functions

// renderTitleBar renders the title bar
func (m model) renderTitleBar() string {
	title := titleStyle.Render("Tmuxplexer")
	padding := m.width - lipgloss.Width(title)
	if padding < 0 {
		padding = 0
	}
	return title + strings.Repeat(" ", padding)
}

// renderStatusBar renders the status bar
func (m model) renderStatusBar() string {
	var line1, line2 string

	// Line 1: Current status message or input prompt
	if m.sessionSaveMode {
		// Show session save wizard prompt
		line1 = m.getSessionSavePrompt() + m.inputBuffer + "█"
	} else if m.inputMode == "rename" {
		line1 = m.inputPrompt + m.inputBuffer + "█" // Show cursor
	} else if m.inputMode == "kill_confirm" {
		// Kill confirmation - make it prominent
		line1 = "⚠️  " + m.inputPrompt + " [Press Y to confirm, N to cancel]"
	} else if m.inputMode == "template_delete_confirm" {
		// Template delete confirmation
		line1 = "⚠️  " + m.inputPrompt + " [Press Y to confirm, N to cancel]"
	} else if m.templateCreationMode {
		// Show wizard input prompt
		line1 = m.getTemplateWizardPrompt() + m.inputBuffer + "█"
	} else {
		line1 = m.statusMsg
	}

	// Use scrolling footer (click to activate) or truncate if too long
	maxLen := m.width - 4
	if maxLen > 0 {
		line1 = m.renderScrollingFooter(line1, maxLen)
	}

	// Pad line1 to full width
	padding1 := m.width - visualWidth(line1)
	if padding1 < 0 {
		padding1 = 0
	}
	line1 = line1 + strings.Repeat(" ", padding1)

	// Line 2: Context-aware help text
	line2 = m.getStatusBarHelpText()

	// Use scrolling footer (click to activate) or truncate if too long
	if maxLen > 0 {
		line2 = m.renderScrollingFooter(line2, maxLen)
	}

	// Pad line2 to full width
	padding2 := m.width - visualWidth(line2)
	if padding2 < 0 {
		padding2 = 0
	}
	line2 = line2 + strings.Repeat(" ", padding2)

	// Combine both lines - use warning style for confirmations
	var line1Styled string
	if m.inputMode == "kill_confirm" || m.inputMode == "template_delete_confirm" {
		// Use warning style but with status bar padding
		line1Styled = warningStyle.Copy().Padding(0, 1).Render(line1)
	} else {
		line1Styled = statusStyle.Render(line1)
	}
	return line1Styled + "\n" + statusStyle.Render(line2)
}

// getStatusBarHelpText returns context-aware help text for the status bar
func (m model) getStatusBarHelpText() string {
	// Build help text based on current focus and mode
	var helpParts []string

	// Focus-specific keys
	switch m.focusState {
	case FocusSessions:
		if m.sessionsTab == "sessions" {
			helpParts = append(helpParts, "[↑↓] Navigate", "[PgUp/PgDn] Scroll", "[Enter/o] Attach", "[s] Save", "[d] Detach", "[x] Kill")
		} else {
			helpParts = append(helpParts, "[↑↓] Navigate", "[PgUp/PgDn] Scroll", "[Enter] Create", "[o] Create & Attach", "[n] New", "[e] Edit", "[d] Delete")
		}
	case FocusPreview:
		helpParts = append(helpParts, "[PgUp/PgDn] Scroll", "[r] Refresh")
	case FocusCommand:
		// Command mode: only show keys that work while typing
		helpParts = append(helpParts, "[Enter] Send", "[Esc] Clear", "[↑↓] History", "[Ctrl+V] Paste", "[Ctrl+R] Refresh")
		// Don't show single-letter hotkeys (a, s, 1, 2, 3, q) - they're typed as command input
		return strings.Join(helpParts, " │ ")
	}

	// Panel switching keys (not shown in command mode - they'd be typed as input)
	helpParts = append(helpParts, "[1/2/3] Switch Panels")

	// Adaptive mode toggle (not shown in command mode)
	if m.adaptiveMode {
		helpParts = append(helpParts, "[a] Adaptive: ON")
	} else {
		helpParts = append(helpParts, "[a] Adaptive: OFF")
	}

	// Global keys (not shown in command mode - 'q' would be typed as input)
	helpParts = append(helpParts, "[Ctrl+R] Refresh", "[q] Quit")

	return strings.Join(helpParts, " │ ")
}

// getTemplateWizardPrompt returns the appropriate prompt for the current wizard step
func (m model) getTemplateWizardPrompt() string {
	builder := m.templateBuilder

	switch builder.fieldName {
	case "name":
		return "Step 1/7: Template name: "
	case "description":
		return "Step 2/7: Description (optional): "
	case "category":
		return "Step 3/7: Category (Projects, Agents, Tools, Custom): "
	case "working_dir":
		return "Step 4/7: Working directory: "
	case "layout":
		return "Step 5/7: Layout (e.g., 2x2, 3x3, 4x2): "
	case "pane_command":
		return lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(lipgloss.Left,
				m.getWizardProgressBar(),
				"",
				"Pane "+string(rune('1'+builder.currentPane))+" command: ",
			),
		)
	case "pane_title":
		return lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(lipgloss.Left,
				m.getWizardProgressBar(),
				"",
				"Pane "+string(rune('1'+builder.currentPane))+" title (optional): ",
			),
		)
	case "pane_working_dir":
		return lipgloss.NewStyle().Render(
			lipgloss.JoinVertical(lipgloss.Left,
				m.getWizardProgressBar(),
				"",
				"Pane "+string(rune('1'+builder.currentPane))+" working dir (optional): ",
			),
		)
	default:
		return "Template Wizard: "
	}
}

// getWizardProgressBar returns a progress indicator for the wizard
func (m model) getWizardProgressBar() string {
	builder := m.templateBuilder
	totalSteps := 5 + builder.numPanes*2 // name, desc, category, dir, layout + (command, title) per pane
	currentStep := 5 // Base steps completed

	// Calculate current step based on field
	switch builder.fieldName {
	case "name":
		currentStep = 1
	case "description":
		currentStep = 2
	case "category":
		currentStep = 3
	case "working_dir":
		currentStep = 4
	case "layout":
		currentStep = 5
	case "pane_command":
		currentStep = 5 + builder.currentPane*2 + 1
	case "pane_title":
		currentStep = 5 + builder.currentPane*2 + 2
	}

	return "Creating template: " + builder.name + " | Step " + string(rune('0'+currentStep)) + "/" + string(rune('0'+totalSteps))
}

// renderMainContent renders the main content area
func (m model) renderMainContent() string {
	contentWidth, contentHeight := m.calculateLayout()

	// Implement your main content rendering here
	// Example:
	// return m.renderItemList(contentWidth, contentHeight)

	placeholder := "Main content area\n\n"
	placeholder += "Implement your content rendering in renderMainContent()\n\n"
	placeholder += "Press ? for help\n"
	placeholder += "Press q to quit"

	return contentStyle.Width(contentWidth).Height(contentHeight).Render(placeholder)
}

// renderLeftPane renders the left pane in dual-pane mode
func (m model) renderLeftPane(width int) string {
	_, contentHeight := m.calculateLayout()

	// Implement left pane content
	content := "Left Pane\n\n"
	content += "Width: " + string(rune(width))

	return leftPaneStyle.Width(width).Height(contentHeight).Render(content)
}

// renderRightPane renders the right pane in dual-pane mode
func (m model) renderRightPane(width int) string {
	_, contentHeight := m.calculateLayout()

	// Implement right pane content
	content := "Right Pane (Preview)\n\n"
	content += "Width: " + string(rune(width))

	return rightPaneStyle.Width(width).Height(contentHeight).Render(content)
}

// renderDivider renders the vertical divider between panes
func (m model) renderDivider() string {
	_, contentHeight := m.calculateLayout()
	divider := strings.Repeat("│\n", contentHeight)
	return dividerStyle.Render(divider)
}

// Error and minimal views

// renderErrorView renders an error message
func (m model) renderErrorView() string {
	content := "Error: " + m.err.Error() + "\n\n"
	content += "Press q to quit"
	return errorStyle.Render(content)
}

// renderMinimalView renders a minimal view for small terminals
func (m model) renderMinimalView() string {
	content := "Terminal too small\n"
	content += "Minimum: 60x15\n"
	content += "Press q to quit"
	return errorStyle.Render(content)
}

// renderUnifiedView renders the unified 3-panel adaptive layout
// Layout: Sessions (top, 40-50%) | Preview (middle, 40-30%) | Command (bottom, 20% fixed)
func (m model) renderUnifiedView() string {
	var sections []string

	// Title bar
	if m.config.UI.ShowTitle {
		sections = append(sections, m.renderTitleBar())
	}

	// Calculate available content height
	contentWidth, contentHeight := m.calculateLayout()

	// Add back the 2 lines that calculateLayout() subtracted for borders
	// calculateAdaptivePanelHeights() will handle all 3 panels' borders (6 lines total)
	contentHeight += 2

	// Get adaptive panel heights based on focus state
	sessionsHeight, previewHeight, commandHeight := m.calculateAdaptivePanelHeights(contentHeight)

	// Render each panel with appropriate content (top panel switches between sessions/templates)
	var topPanelContent []string
	var topPanelName string
	if m.sessionsTab == "templates" {
		topPanelContent = m.templatesContent
		topPanelName = "templates"
	} else {
		topPanelContent = m.sessionsContent
		topPanelName = "sessions"
	}
	sessionsPanel := m.renderDynamicPanel(topPanelName, contentWidth, sessionsHeight, topPanelContent)
	previewPanel := m.renderDynamicPanel("preview", contentWidth, previewHeight, m.previewContent)
	commandPanel := m.renderDynamicPanel("command", contentWidth, commandHeight, m.commandContent)

	// Stack panels vertically
	sections = append(sections, sessionsPanel, previewPanel, commandPanel)

	// Status bar
	if m.config.UI.ShowStatus {
		sections = append(sections, m.renderStatusBar())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderDynamicPanel renders a single dynamic panel with border and content
func (m model) renderDynamicPanel(panelName string, width, height int, content []string) string {
	// In unified layout, panels are focused based on focusState
	isFocused := false
	switch m.focusState {
	case FocusSessions:
		isFocused = (panelName == "sessions" || panelName == "templates")
	case FocusPreview:
		isFocused = (panelName == "preview")
	case FocusCommand:
		isFocused = (panelName == "command")
	}

	// Create border style based on focus
	borderColor := lipgloss.Color("240") // Dim gray

	if isFocused {
		borderColor = colorPrimary // Bright blue
	}

	// Panel titles
	titles := map[string]string{
		"sessions":  "Sessions",
		"templates": "Templates",
		"preview":   "Preview",
		"command":   "Command",
	}

	var title string
	if isFocused {
		// Special rendering for Sessions/Templates tabs (top panel)
		if panelName == "sessions" || panelName == "templates" {
			// Show both tabs with active one in blue
			activeTabStyle := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
			inactiveTabStyle := lipgloss.NewStyle().Foreground(colorForeground)

			sessionsText := "Sessions"
			templatesText := "Templates"

			if panelName == "sessions" {
				sessionsText = activeTabStyle.Render(sessionsText)
				templatesText = inactiveTabStyle.Render(templatesText)
			} else {
				sessionsText = inactiveTabStyle.Render(sessionsText)
				templatesText = activeTabStyle.Render(templatesText)
			}

			title = " " + sessionsText + " | " + templatesText + " ● "
		} else {
			// Other panels: simple focused indicator
			title = " " + titles[panelName] + " ● "
		}
	} else {
		// Unfocused: plain title
		if panelName == "sessions" || panelName == "templates" {
			title = " Sessions | Templates "
		} else {
			title = " " + titles[panelName] + " "
		}
	}

	// Calculate max text width
	maxTextWidth := width - 2 // -2 for borders
	if maxTextWidth < 1 {
		maxTextWidth = 1
	}

	// Calculate exact content area height
	// Since title is now in border, we get full inner height for content
	innerHeight := height - 2 // Remove borders
	availableContentLines := innerHeight

	if availableContentLines < 1 {
		availableContentLines = 1
	}

	// Determine scroll offset for sessions/templates panels
	scrollOffset := 0
	if panelName == "sessions" || panelName == "templates" {
		scrollOffset = m.sessionsScrollOffset
	}

	// Build content lines (truncate if too long to prevent wrapping)
	var lines []string
	startIdx := scrollOffset
	for i := 0; i < availableContentLines && (startIdx + i) < len(content); i++ {
		line := content[startIdx + i]

		// Apply styling based on tags
		if line == "DIVIDER" {
			// Full-width divider
			line = strings.Repeat("─", maxTextWidth)
		} else if line == "SESSION_DIVIDER" {
			// Session divider with indent
			dividerWidth := maxTextWidth - 2
			if dividerWidth < 1 {
				dividerWidth = 1
			}
			line = "  " + strings.Repeat("┄", dividerWidth)
		} else if strings.HasPrefix(line, "DETAILS:header:") {
			// Section header style
			text := strings.TrimPrefix(line, "DETAILS:header:")
			line = sectionHeaderStyle.Render(truncateString(text, maxTextWidth))
		} else if strings.HasPrefix(line, "DETAILS:detail:") {
			// Detail text style (dimmed)
			text := strings.TrimPrefix(line, "DETAILS:detail:")
			line = dimmedStyle.Render(truncateString(text, maxTextWidth))
		} else if strings.HasPrefix(line, "HEADER:") {
			// Table header style (bold + primary color)
			text := strings.TrimPrefix(line, "HEADER:")
			line = tableHeaderStyle.Render(truncateString(text, maxTextWidth))
		} else if strings.HasPrefix(line, "CURRENT:") {
			// Current session style (cyan text, bold) - takes precedence over Claude
			text := strings.TrimPrefix(line, "CURRENT:")
			line = currentSessionStyle.Render(truncateString(text, maxTextWidth))
		} else if strings.HasPrefix(line, "CLAUDE:") {
			// Claude session style (orange text)
			text := strings.TrimPrefix(line, "CLAUDE:")
			line = claudeSessionStyle.Render(truncateString(text, maxTextWidth))
		} else if strings.HasPrefix(line, "SELECTED:") {
			// Selected tree item style (bold + underline)
			text := strings.TrimPrefix(line, "SELECTED:")
			line = selectedTreeItemStyle.Render(truncateString(text, maxTextWidth))
		} else {
			// Normal text
			line = truncateString(line, maxTextWidth)
		}

		lines = append(lines, line)
	}

	// Add scroll indicators for sessions/templates panels
	if panelName == "sessions" || panelName == "templates" {
		totalLines := len(content)
		canScrollUp := scrollOffset > 0
		canScrollDown := (startIdx + availableContentLines) < totalLines

		if canScrollUp || canScrollDown {
			// Add scroll position indicator
			scrollInfo := fmt.Sprintf(" Lines %d-%d of %d ",
				scrollOffset+1,
				min(scrollOffset+availableContentLines, totalLines),
				totalLines)

			// If there's room, add the indicator as the last line
			if len(lines) > 0 {
				indicatorStyle := lipgloss.NewStyle().Foreground(colorDimmed).Italic(true)
				lines[len(lines)-1] = indicatorStyle.Render(scrollInfo)
			}
		}
	}

	// Fill remaining space to ensure consistent height
	for len(lines) < innerHeight {
		lines = append(lines, "")
	}

	contentStr := strings.Join(lines, "\n")

	// Create custom border with title in top border (lazygit style)
	border := lipgloss.RoundedBorder()

	// Calculate how much space we have for the title in the top border
	// width - 2 for corner characters
	topBorderSpace := width - 2

	// Create title with padding
	titleLen := lipgloss.Width(title)
	if titleLen > topBorderSpace - 2 {
		// Truncate title if too long
		title = truncateString(title, topBorderSpace - 2)
		titleLen = lipgloss.Width(title)
	}

	// Build top border: ╭ title ─────╮
	// Style border characters to maintain color
	borderStyle := lipgloss.NewStyle().Foreground(borderColor)

	leftBorder := borderStyle.Render(border.TopLeft)
	rightBorder := borderStyle.Render(border.TopRight)
	fillChar := border.Top

	// Calculate fill needed after title
	fillNeeded := topBorderSpace - titleLen
	if fillNeeded < 0 {
		fillNeeded = 0
	}

	// Build top border with styled elements
	fillString := borderStyle.Render(strings.Repeat(fillChar, fillNeeded))
	customTopBorder := leftBorder + title + fillString + rightBorder

	// Build the complete box manually to ensure proper border connection
	var boxLines []string

	// Add top border
	boxLines = append(boxLines, customTopBorder)

	// Add content lines with left/right borders (reuse borderStyle from above)
	contentLines := strings.Split(contentStr, "\n")
	for _, line := range contentLines {
		// Safety: ensure line is never longer than maxTextWidth (re-truncate if needed)
		lineWidth := lipgloss.Width(line)
		if lineWidth > maxTextWidth {
			line = truncateString(line, maxTextWidth)
			lineWidth = lipgloss.Width(line)
		}
		// Ensure line is exactly maxTextWidth (pad if needed)
		if lineWidth < maxTextWidth {
			line = line + strings.Repeat(" ", maxTextWidth-lineWidth)
		}
		// Style the border characters to prevent color bleeding from styled content
		leftBorder := borderStyle.Render(border.Left)
		rightBorder := borderStyle.Render(border.Right)
		boxLines = append(boxLines, leftBorder+line+rightBorder)
	}

	// Add bottom border with styled elements
	bottomLeft := borderStyle.Render(border.BottomLeft)
	bottomFill := borderStyle.Render(strings.Repeat(border.Bottom, maxTextWidth))
	bottomRight := borderStyle.Render(border.BottomRight)
	bottomBorder := bottomLeft + bottomFill + bottomRight
	boxLines = append(boxLines, bottomBorder)

	// Join all lines and apply color
	fullBox := strings.Join(boxLines, "\n")

	return lipgloss.NewStyle().
		Foreground(borderColor).
		Render(fullBox)
}

// renderVerticalDivider renders a vertical divider between panels
func renderVerticalDivider(height int) string {
	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	var lines []string
	for i := 0; i < height; i++ {
		lines = append(lines, "│")
	}

	return dividerStyle.Render(strings.Join(lines, "\n"))
}

// Helper functions

// truncateString truncates a string to fit within maxWidth
func truncateString(s string, maxWidth int) string {
	// Use lipgloss.Width to properly measure visual width (ignoring ANSI codes)
	currentWidth := lipgloss.Width(s)
	if currentWidth <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		// Very narrow, just show first few runes
		runes := []rune(s)
		if len(runes) > maxWidth {
			return string(runes[:maxWidth])
		}
		return s
	}

	// Truncate by removing runes from the end until we fit
	runes := []rune(s)
	targetWidth := maxWidth - 3 // Reserve space for "..."

	for len(runes) > 0 && lipgloss.Width(string(runes)) > targetWidth {
		runes = runes[:len(runes)-1]
	}

	return string(runes) + "..."
}

// padRight pads a string with spaces to reach the desired width
func padRight(s string, width int) string {
	currentWidth := lipgloss.Width(s)
	if currentWidth >= width {
		return s
	}
	return s + strings.Repeat(" ", width-currentWidth)
}

// centerString centers a string within the given width
func centerString(s string, width int) string {
	strWidth := lipgloss.Width(s)
	if strWidth >= width {
		return s
	}
	leftPad := (width - strWidth) / 2
	rightPad := width - strWidth - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

// renderScrollingFooter renders footer text with horizontal scrolling if enabled
// If text fits within width, returns as-is. If scrolling is enabled and text is too long,
// creates a looping marquee effect. Otherwise truncates with "..."
func (m model) renderScrollingFooter(text string, availableWidth int) string {
	textLen := visualWidth(text)

	// If text fits, no modification needed
	if textLen <= availableWidth {
		return text
	}

	// If scrolling is active, create looping marquee
	if m.footerScrolling {
		// Add visual indicator and separator for smooth loop
		indicator := "⏵ " // Indicates scrolling is active
		paddedText := indicator + text + "   •   " + indicator + text

		// Convert to runes to handle multi-byte unicode characters (↑, ↓, •, etc.)
		runes := []rune(paddedText)
		runeCount := len(runes)

		// Calculate scroll position with wrapping
		scrollPos := m.footerOffset % runeCount

		// Extract visible portion (by rune, not byte)
		var result strings.Builder
		for i := 0; i < availableWidth && i < runeCount; i++ {
			charPos := (scrollPos + i) % runeCount
			result.WriteRune(runes[charPos])
		}

		return result.String()
	}

	// Not scrolling - truncate with "..."
	return truncateString(text, availableWidth)
}
