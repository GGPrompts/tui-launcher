package launch

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle = lipgloss.NewStyle().Bold(true)
)

// View renders the Launch tab UI
func (m Model) View() string {
	if m.loading {
		return m.spinner.View() + " Loading configuration...\n"
	}

	if m.err != nil {
		return "Error: " + m.err.Error() + "\n\nPress q to quit.\n"
	}

	var sb strings.Builder

	// Header (3 lines total)
	sb.WriteString("🚀 TUI Launcher - Launch Tab\n")

	// Show current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" && strings.HasPrefix(cwd, homeDir) {
		cwd = "~" + strings.TrimPrefix(cwd, homeDir)
	}
	sb.WriteString("Working Dir: " + cwd)

	if m.insideTmux {
		sb.WriteString(" (tmux)")
	}
	sb.WriteString(" | Mode: ")
	if m.useTmux {
		sb.WriteString("Tmux")
	} else {
		sb.WriteString("Direct")
	}
	sb.WriteString("\n\n")

	// Get layout dimensions
	leftWidth, rightWidth, treeHeight, infoHeight := m.calculateLayout()
	mode := m.getLayoutMode()

	// Define border style
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	switch mode {
	case layoutDesktop:
		// 3-pane layout: Left | Right (top), Info (bottom)
		leftContent := m.viewLeftPane(leftWidth-2, treeHeight)   // -2 for borders
		rightContent := m.viewRightPane(rightWidth-2, treeHeight) // -2 for borders

		leftPane := borderStyle.Width(leftWidth - 2).Render(leftContent)
		rightPane := borderStyle.Width(rightWidth - 2).Render(rightContent)

		// Join left and right panes horizontally
		topRow := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
		sb.WriteString(topRow)
		sb.WriteString("\n")

		// Info pane (full width)
		if infoHeight > 0 {
			infoContent := m.viewInfoPane(m.width-2, infoHeight)
			infoPane := borderStyle.Width(m.width - 2).Render(infoContent)
			sb.WriteString(infoPane)
		}

	case layoutCompact:
		// 2-pane layout: Combined tree (top), Info (bottom)
		treeContent := m.viewCombinedTree(m.width-2, treeHeight)
		treePane := borderStyle.Width(m.width - 2).Render(treeContent)
		sb.WriteString(treePane)
		sb.WriteString("\n")

		// Info pane
		if infoHeight > 0 {
			infoContent := m.viewInfoPane(m.width-2, infoHeight)
			infoPane := borderStyle.Width(m.width - 2).Render(infoContent)
			sb.WriteString(infoPane)
		}

	case layoutMobile:
		// 1-pane layout: Just tree (or info if toggled)
		if m.showingInfo {
			// Show info pane instead of tree
			infoContent := m.viewInfoPane(m.width-2, treeHeight)
			infoPane := borderStyle.Width(m.width - 2).Render(infoContent)
			sb.WriteString(infoPane)
		} else {
			// Show tree
			treeContent := m.viewCombinedTree(m.width-2, treeHeight)
			treePane := borderStyle.Width(m.width - 2).Render(treeContent)
			sb.WriteString(treePane)
		}
	}

	// Status line
	sb.WriteString("\n")
	if len(m.selectedItems) > 0 {
		sb.WriteString(fmt.Sprintf("Selected: %d items | ", len(m.selectedItems)))
	}

	// Footer (adapts to layout mode) - static, no scrolling to prevent flashing
	var footerText string
	switch mode {
	case layoutDesktop:
		footerText = "↑/↓: nav  Tab: panes  Space: expand/select  Enter: launch  e: edit  c: clear  q: quit"
	case layoutCompact:
		footerText = "↑/↓: nav  Tab: switch  Space: select  Enter: launch  e: edit  c: clear  q: quit"
	case layoutMobile:
		footerText = "↑/↓: nav  Tab: switch  i: info  Space: select  Enter: launch  q: quit"
	}

	// Truncate footer if needed (no scrolling)
	if len(footerText) > m.width-2 {
		footerText = footerText[:m.width-5] + "..."
	}
	sb.WriteString(footerText)
	sb.WriteString("\n")

	return sb.String()
}

// viewLeftPane renders the global tools tree (left pane in desktop mode)
func (m Model) viewLeftPane(width, height int) string {
	var lines []string

	// Add title
	title := "Global Tools"
	if m.activePane == paneGlobal {
		title = "> " + title + " <"
	}
	lines = append(lines, title)
	lines = append(lines, "") // Blank line

	// Render tree items
	if len(m.globalTreeItems) == 0 {
		lines = append(lines, "(no global tools)")
	} else {
		for i, ti := range m.globalTreeItems {
			selected := m.selectedItems[ti.item.Path]
			expanded := m.globalExpanded[ti.item.Path]
			line := renderTreeItem(ti, m.globalCursor, i, selected, expanded)

			// Truncate to prevent wrapping
			maxWidth := width - 4 // Account for padding
			if len(line) > maxWidth {
				line = line[:maxWidth-1] + "…"
			}

			// Highlight cursor
			if i == m.globalCursor && m.activePane == paneGlobal {
				line = selectedStyle.Render(line)
			}

			lines = append(lines, line)
		}
	}

	// Fill to exact height
	for len(lines) < height {
		lines = append(lines, "")
	}

	// Truncate if too many lines
	if len(lines) > height {
		lines = lines[:height]
	}

	return strings.Join(lines, "\n")
}

// viewRightPane renders the projects tree (right pane in desktop mode)
func (m Model) viewRightPane(width, height int) string {
	var lines []string

	// Add title
	title := "Projects"
	if m.activePane == paneProject {
		title = "> " + title + " <"
	}
	lines = append(lines, title)
	lines = append(lines, "") // Blank line

	// Render tree items
	if len(m.projectTreeItems) == 0 {
		lines = append(lines, "(no projects)")
	} else {
		for i, ti := range m.projectTreeItems {
			selected := m.selectedItems[ti.item.Path]
			expanded := m.projectExpanded[ti.item.Path]
			line := renderTreeItem(ti, m.projectCursor, i, selected, expanded)

			// Truncate to prevent wrapping
			maxWidth := width - 4 // Account for padding
			if len(line) > maxWidth {
				line = line[:maxWidth-1] + "…"
			}

			// Highlight cursor
			if i == m.projectCursor && m.activePane == paneProject {
				line = selectedStyle.Render(line)
			}

			lines = append(lines, line)
		}
	}

	// Fill to exact height
	for len(lines) < height {
		lines = append(lines, "")
	}

	// Truncate if too many lines
	if len(lines) > height {
		lines = lines[:height]
	}

	return strings.Join(lines, "\n")
}

// viewInfoPane renders the info/help pane (bottom pane)
func (m Model) viewInfoPane(width, height int) string {
	var lines []string

	// Add title
	lines = append(lines, "Info")
	lines = append(lines, "")

	// Show info content or help text
	if m.infoContent != "" {
		// Split content into lines
		contentLines := strings.Split(m.infoContent, "\n")
		for _, line := range contentLines {
			// Truncate to prevent wrapping
			maxWidth := width - 4
			if len(line) > maxWidth {
				line = line[:maxWidth-1] + "…"
			}
			lines = append(lines, line)
		}
	} else {
		// Default help text
		lines = append(lines, "Navigate with arrows or vim keys")
		lines = append(lines, "Space: expand/select  Enter: launch  Tab: switch panes")
		lines = append(lines, "t: toggle mode  c: clear  e: edit config  q: quit")
	}

	// Fill to exact height
	for len(lines) < height {
		lines = append(lines, "")
	}

	// Truncate if too many lines
	if len(lines) > height {
		lines = lines[:height]
	}

	return strings.Join(lines, "\n")
}

// viewCombinedTree renders a combined tree for compact/mobile modes
func (m Model) viewCombinedTree(width, height int) string {
	var lines []string

	// Determine which pane to show in compact mode
	var items []launchTreeItem
	var cursor int
	var expanded map[string]bool
	var title string

	if m.showingProjects {
		items = m.projectTreeItems
		cursor = m.projectCursor
		expanded = m.projectExpanded
		title = "Projects"
	} else {
		items = m.globalTreeItems
		cursor = m.globalCursor
		expanded = m.globalExpanded
		title = "Global Tools"
	}

	// Add title
	lines = append(lines, title+" (Tab to switch)")
	lines = append(lines, "")

	// Render tree items
	if len(items) == 0 {
		lines = append(lines, "(no items)")
	} else {
		for i, ti := range items {
			selected := m.selectedItems[ti.item.Path]
			isExpanded := expanded[ti.item.Path]
			line := renderTreeItem(ti, cursor, i, selected, isExpanded)

			// Truncate to prevent wrapping
			maxWidth := width - 4
			if len(line) > maxWidth {
				line = line[:maxWidth-1] + "…"
			}

			// Highlight cursor
			if i == cursor {
				line = selectedStyle.Render(line)
			}

			lines = append(lines, line)
		}
	}

	// Fill to exact height
	for len(lines) < height {
		lines = append(lines, "")
	}

	// Truncate if too many lines
	if len(lines) > height {
		lines = lines[:height]
	}

	return strings.Join(lines, "\n")
}

// renderTreeItem renders a single tree item with proper indentation
func renderTreeItem(ti launchTreeItem, cursor int, index int, selected bool, expanded bool) string {
	var sb strings.Builder

	// Cursor indicator
	if index == cursor {
		sb.WriteString("> ")
	} else {
		sb.WriteString("  ")
	}

	// Tree lines
	for i, isLast := range ti.parentLasts {
		if i == len(ti.parentLasts)-1 {
			continue
		}
		if isLast {
			sb.WriteString("  ")
		} else {
			sb.WriteString("│ ")
		}
	}

	// Branch character
	if ti.depth > 0 {
		if ti.isLast {
			sb.WriteString("└─")
		} else {
			sb.WriteString("├─")
		}
	}

	// Selection checkbox
	if selected {
		sb.WriteString("☑ ")
	} else if ti.item.ItemType == typeCommand || ti.item.ItemType == typeProfile {
		sb.WriteString("☐ ")
	}

	// Expansion indicator for categories
	if ti.item.ItemType == typeCategory {
		if expanded {
			sb.WriteString("▼ ")
		} else {
			sb.WriteString("▶ ")
		}
	}

	// Icon
	if ti.item.Icon != "" {
		sb.WriteString(ti.item.Icon + " ")
	}

	// Name
	sb.WriteString(ti.item.Name)

	// Type indicator (only show for profiles)
	switch ti.item.ItemType {
	case typeProfile:
		sb.WriteString(fmt.Sprintf(" [%s]", ti.item.LayoutStr))
	}

	return sb.String()
}
