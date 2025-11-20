package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

var (
	selectedStyle = lipgloss.NewStyle().Bold(true)
)

// footerTick sends periodic messages to animate footer scrolling
func footerTick() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return footerTickMsg{}
	})
}

// initialModel creates the initial application model
func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		width:          80,
		height:         24,
		cursor:         0,
		treeItems:      []launchTreeItem{},
		expandedItems:  make(map[string]bool),

		// Multi-pane initialization
		activePane:       paneGlobal,
		globalItems:      []launchItem{},
		projectItems:     []launchItem{},
		globalTreeItems:  []launchTreeItem{},
		projectTreeItems: []launchTreeItem{},
		globalCursor:     0,
		projectCursor:    0,
		globalExpanded:   make(map[string]bool),
		projectExpanded:  make(map[string]bool),

		// Info pane initialization
		showingInfo:     false,
		showingProjects: false,

		selectedItems:  make(map[string]bool),
		showSpawnDialog: false,
		selectedLayout: layoutTiled,
		layoutCursor:   0,
		spinner:        s,
		loading:        true,
		terminalType:   detectTerminal(),
		insideTmux:     isInsideTmux(),
		useTmux:        true, // Default to tmux mode
	}
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadConfig,
	)
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab":
			// Handle Tab key based on layout mode
			mode := m.getLayoutMode()
			switch mode {
			case layoutDesktop:
				// Toggle between left and right panes
				if m.activePane == paneGlobal {
					m.activePane = paneProject
				} else {
					m.activePane = paneGlobal
				}
			case layoutCompact:
				// Toggle between global and projects
				m.showingProjects = !m.showingProjects
			case layoutMobile:
				// Toggle between global and projects in mobile mode
				m.showingProjects = !m.showingProjects
			}
			// Update info pane after switching panes
			m.updateInfoPane()

		case "i":
			// Toggle info pane in mobile mode
			if m.getLayoutMode() == layoutMobile {
				m.showingInfo = !m.showingInfo
			}

		case "up", "k":
			// Navigate up in the active pane
			mode := m.getLayoutMode()
			if mode == layoutDesktop {
				if m.activePane == paneGlobal {
					if m.globalCursor > 0 {
						m.globalCursor--
					}
				} else {
					if m.projectCursor > 0 {
						m.projectCursor--
					}
				}
			} else {
				// Compact/mobile mode - navigate in current view
				if m.showingProjects {
					if m.projectCursor > 0 {
						m.projectCursor--
					}
				} else {
					if m.globalCursor > 0 {
						m.globalCursor--
					}
				}
			}
			// Update info pane after cursor movement
			m.updateInfoPane()

		case "down", "j":
			// Navigate down in the active pane
			mode := m.getLayoutMode()
			if mode == layoutDesktop {
				if m.activePane == paneGlobal {
					if m.globalCursor < len(m.globalTreeItems)-1 {
						m.globalCursor++
					}
				} else {
					if m.projectCursor < len(m.projectTreeItems)-1 {
						m.projectCursor++
					}
				}
			} else {
				// Compact/mobile mode - navigate in current view
				if m.showingProjects {
					if m.projectCursor < len(m.projectTreeItems)-1 {
						m.projectCursor++
					}
				} else {
					if m.globalCursor < len(m.globalTreeItems)-1 {
						m.globalCursor++
					}
				}
			}
			// Update info pane after cursor movement
			m.updateInfoPane()

		case "right", "l":
			// Expand current category in active pane
			mode := m.getLayoutMode()
			if mode == layoutDesktop {
				if m.activePane == paneGlobal {
					if m.globalCursor < len(m.globalTreeItems) {
						currentItem := m.globalTreeItems[m.globalCursor].item
						if currentItem.ItemType == typeCategory {
							m.globalExpanded[currentItem.Path] = true
							m.globalTreeItems = flattenTree(m.globalItems, m.globalExpanded)
						}
					}
				} else {
					if m.projectCursor < len(m.projectTreeItems) {
						currentItem := m.projectTreeItems[m.projectCursor].item
						if currentItem.ItemType == typeCategory {
							m.projectExpanded[currentItem.Path] = true
							m.projectTreeItems = flattenTree(m.projectItems, m.projectExpanded)
						}
					}
				}
			} else {
				// Compact/mobile mode
				if m.showingProjects {
					if m.projectCursor < len(m.projectTreeItems) {
						currentItem := m.projectTreeItems[m.projectCursor].item
						if currentItem.ItemType == typeCategory {
							m.projectExpanded[currentItem.Path] = true
							m.projectTreeItems = flattenTree(m.projectItems, m.projectExpanded)
						}
					}
				} else {
					if m.globalCursor < len(m.globalTreeItems) {
						currentItem := m.globalTreeItems[m.globalCursor].item
						if currentItem.ItemType == typeCategory {
							m.globalExpanded[currentItem.Path] = true
							m.globalTreeItems = flattenTree(m.globalItems, m.globalExpanded)
						}
					}
				}
			}

		case "left", "h":
			// Collapse current category in active pane
			mode := m.getLayoutMode()
			if mode == layoutDesktop {
				if m.activePane == paneGlobal {
					if m.globalCursor < len(m.globalTreeItems) {
						currentItem := m.globalTreeItems[m.globalCursor].item
						if currentItem.ItemType == typeCategory {
							m.globalExpanded[currentItem.Path] = false
							m.globalTreeItems = flattenTree(m.globalItems, m.globalExpanded)
						}
					}
				} else {
					if m.projectCursor < len(m.projectTreeItems) {
						currentItem := m.projectTreeItems[m.projectCursor].item
						if currentItem.ItemType == typeCategory {
							m.projectExpanded[currentItem.Path] = false
							m.projectTreeItems = flattenTree(m.projectItems, m.projectExpanded)
						}
					}
				}
			} else {
				// Compact/mobile mode
				if m.showingProjects {
					if m.projectCursor < len(m.projectTreeItems) {
						currentItem := m.projectTreeItems[m.projectCursor].item
						if currentItem.ItemType == typeCategory {
							m.projectExpanded[currentItem.Path] = false
							m.projectTreeItems = flattenTree(m.projectItems, m.projectExpanded)
						}
					}
				} else {
					if m.globalCursor < len(m.globalTreeItems) {
						currentItem := m.globalTreeItems[m.globalCursor].item
						if currentItem.ItemType == typeCategory {
							m.globalExpanded[currentItem.Path] = false
							m.globalTreeItems = flattenTree(m.globalItems, m.globalExpanded)
						}
					}
				}
			}

		case " ":
			// Context-aware action: expand/collapse categories, or toggle selection for commands
			mode := m.getLayoutMode()
			var currentItem launchItem
			var isValid bool

			// Get current item based on active pane and mode
			if mode == layoutDesktop {
				if m.activePane == paneGlobal {
					if m.globalCursor < len(m.globalTreeItems) {
						currentItem = m.globalTreeItems[m.globalCursor].item
						isValid = true
					}
				} else {
					if m.projectCursor < len(m.projectTreeItems) {
						currentItem = m.projectTreeItems[m.projectCursor].item
						isValid = true
					}
				}
			} else {
				// Compact/mobile mode
				if m.showingProjects {
					if m.projectCursor < len(m.projectTreeItems) {
						currentItem = m.projectTreeItems[m.projectCursor].item
						isValid = true
					}
				} else {
					if m.globalCursor < len(m.globalTreeItems) {
						currentItem = m.globalTreeItems[m.globalCursor].item
						isValid = true
					}
				}
			}

			if isValid {
				if currentItem.ItemType == typeCategory {
					// Toggle expansion for categories
					if mode == layoutDesktop {
						if m.activePane == paneGlobal {
							m.globalExpanded[currentItem.Path] = !m.globalExpanded[currentItem.Path]
							m.globalTreeItems = flattenTree(m.globalItems, m.globalExpanded)
						} else {
							m.projectExpanded[currentItem.Path] = !m.projectExpanded[currentItem.Path]
							m.projectTreeItems = flattenTree(m.projectItems, m.projectExpanded)
						}
					} else {
						if m.showingProjects {
							m.projectExpanded[currentItem.Path] = !m.projectExpanded[currentItem.Path]
							m.projectTreeItems = flattenTree(m.projectItems, m.projectExpanded)
						} else {
							m.globalExpanded[currentItem.Path] = !m.globalExpanded[currentItem.Path]
							m.globalTreeItems = flattenTree(m.globalItems, m.globalExpanded)
						}
					}
				} else if currentItem.ItemType == typeCommand || currentItem.ItemType == typeProfile {
					// Toggle selection for commands/profiles
					if m.selectedItems[currentItem.Path] {
						delete(m.selectedItems, currentItem.Path)
					} else {
						m.selectedItems[currentItem.Path] = true
					}
				}
			}

		case "c":
			// Clear all selections
			m.selectedItems = make(map[string]bool)

		case "t":
			// Toggle tmux/xterm mode
			m.useTmux = !m.useTmux

		case "e":
			// Edit config file
			if m.insideTmux {
				// If inside tmux, spawn editor in a new split
				return m, editConfigInTmux()
			} else {
				// Otherwise quit and open editor
				return m, tea.Sequence(
					tea.Quit,
					editConfig,
				)
			}

		case "enter":
			// Launch selected items or current item
			mode := m.getLayoutMode()
			var currentItem launchItem
			var isValid bool

			// Get current item from active pane
			if mode == layoutDesktop {
				if m.activePane == paneGlobal {
					if m.globalCursor < len(m.globalTreeItems) {
						currentItem = m.globalTreeItems[m.globalCursor].item
						isValid = true
					}
				} else {
					if m.projectCursor < len(m.projectTreeItems) {
						currentItem = m.projectTreeItems[m.projectCursor].item
						isValid = true
					}
				}
			} else {
				// Compact/mobile mode
				if m.showingProjects {
					if m.projectCursor < len(m.projectTreeItems) {
						currentItem = m.projectTreeItems[m.projectCursor].item
						isValid = true
					}
				} else {
					if m.globalCursor < len(m.globalTreeItems) {
						currentItem = m.globalTreeItems[m.globalCursor].item
						isValid = true
					}
				}
			}

			if isValid {
				// If items are selected, launch them
				if len(m.selectedItems) > 0 {
					// Launch all selected items in batch
					var itemsToLaunch []launchItem

					// Collect all selected items from both panes
					for _, ti := range m.globalTreeItems {
						if m.selectedItems[ti.item.Path] {
							itemsToLaunch = append(itemsToLaunch, ti.item)
						}
					}
					for _, ti := range m.projectTreeItems {
						if m.selectedItems[ti.item.Path] {
							itemsToLaunch = append(itemsToLaunch, ti.item)
						}
					}

					if len(itemsToLaunch) > 0 {
						// Use default layout for batch launch
						return m, spawnMultiple(itemsToLaunch, m.selectedLayout)
					}

				} else {
					// No selection - launch current item if it's a command or profile
					if currentItem.ItemType == typeCommand {
						// Launch single command
						if !m.useTmux {
							// Non-tmux mode: run command directly in current terminal
							return m, tea.Sequence(
								tea.Quit,
								runCommandDirectly(currentItem),
							)
						} else {
							// Tmux mode: use configured spawn mode
							return m, spawnSingle(currentItem, currentItem.DefaultSpawn)
						}

					} else if currentItem.ItemType == typeProfile {
						// Launch profile (convert panes to launch items)
						var itemsToLaunch []launchItem
						for i, pane := range currentItem.Panes {
							item := launchItem{
								Name:    fmt.Sprintf("%s-pane-%d", currentItem.Name, i),
								Command: pane.Command,
								Cwd:     expandPath(pane.Cwd),
							}
							itemsToLaunch = append(itemsToLaunch, item)
						}
						return m, spawnMultiple(itemsToLaunch, currentItem.Layout)
					}
				}
			}
		}

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			// Click to switch panes in desktop mode
			// GOLDEN RULE #3: Match mouse detection to layout (X for horizontal split)
			mode := m.getLayoutMode()
			if mode == layoutDesktop {
				leftWidth, _, _, _ := m.calculateLayout()
				// Header is 3 lines tall
				if msg.Y > 3 {
					if msg.X < leftWidth {
						m.activePane = paneGlobal
					} else {
						m.activePane = paneProject
					}
					m.updateInfoPane()
				}
			}

		case tea.MouseWheelUp:
			// Scroll up in active pane
			mode := m.getLayoutMode()
			if mode == layoutDesktop {
				if m.activePane == paneGlobal {
					if m.globalCursor > 0 {
						m.globalCursor--
					}
				} else {
					if m.projectCursor > 0 {
						m.projectCursor--
					}
				}
			} else {
				// Compact/mobile mode
				if m.showingProjects {
					if m.projectCursor > 0 {
						m.projectCursor--
					}
				} else {
					if m.globalCursor > 0 {
						m.globalCursor--
					}
				}
			}
			// Update info pane after scrolling
			m.updateInfoPane()

		case tea.MouseWheelDown:
			// Scroll down in active pane
			mode := m.getLayoutMode()
			if mode == layoutDesktop {
				if m.activePane == paneGlobal {
					if m.globalCursor < len(m.globalTreeItems)-1 {
						m.globalCursor++
					}
				} else {
					if m.projectCursor < len(m.projectTreeItems)-1 {
						m.projectCursor++
					}
				}
			} else {
				// Compact/mobile mode
				if m.showingProjects {
					if m.projectCursor < len(m.projectTreeItems)-1 {
						m.projectCursor++
					}
				} else {
					if m.globalCursor < len(m.globalTreeItems)-1 {
						m.globalCursor++
					}
				}
			}
			// Update info pane after scrolling
			m.updateInfoPane()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case configLoadedMsg:
		m.loading = false
		m.config = msg.config
		m.err = msg.err
		if msg.err == nil {
			// Build trees from config (split into global and project panes)
			m.globalItems, m.projectItems = buildTreeFromConfig(msg.config)
			m.globalTreeItems = flattenTree(m.globalItems, m.globalExpanded)
			m.projectTreeItems = flattenTree(m.projectItems, m.projectExpanded)

			// Keep legacy single-pane view for backwards compatibility
			// Combine both for the old view
			m.rootItems = append(m.globalItems, m.projectItems...)
			m.treeItems = flattenTree(m.rootItems, m.expandedItems)

			// Update info pane for initial selection
			m.updateInfoPane()
		}

	case spawnCompleteMsg:
		m.err = msg.err
		// Clear selections after launch
		if msg.err == nil {
			m.selectedItems = make(map[string]bool)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// viewLeftPane renders the global tools tree (left pane in desktop mode)
func (m model) viewLeftPane(width, height int) string {
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

			// GOLDEN RULE #2: Truncate to prevent wrapping
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

	// Fill to exact height (GOLDEN RULE #1: height already accounts for borders)
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
func (m model) viewRightPane(width, height int) string {
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

			// GOLDEN RULE #2: Truncate to prevent wrapping
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
func (m model) viewInfoPane(width, height int) string {
	var lines []string

	// Add title
	lines = append(lines, "Info")
	lines = append(lines, "")

	// Show info content or help text
	if m.infoContent != "" {
		// Split content into lines
		contentLines := strings.Split(m.infoContent, "\n")
		for _, line := range contentLines {
			// GOLDEN RULE #2: Truncate to prevent wrapping
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
func (m model) viewCombinedTree(width, height int) string {
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
	lines = append(lines, title + " (Tab to switch)")
	lines = append(lines, "")

	// Render tree items
	if len(items) == 0 {
		lines = append(lines, "(no items)")
	} else {
		for i, ti := range items {
			selected := m.selectedItems[ti.item.Path]
			isExpanded := expanded[ti.item.Path]
			line := renderTreeItem(ti, cursor, i, selected, isExpanded)

			// GOLDEN RULE #2: Truncate to prevent wrapping
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

// View renders the UI
func (m model) View() string {
	if m.loading {
		return m.spinner.View() + " Loading configuration...\n"
	}

	if m.err != nil {
		return "Error: " + m.err.Error() + "\n\nPress q to quit.\n"
	}

	var sb strings.Builder

	// Header (3 lines total)
	sb.WriteString("🚀 TUI Launcher v" + Version)

	// Show layout mode for debugging
	mode := m.getLayoutMode()
	sb.WriteString(fmt.Sprintf(" [%s]", mode.String()))
	sb.WriteString("\n")

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
	mode = m.getLayoutMode()

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

// renderScrollingFooter renders footer text with horizontal scrolling if needed
func renderScrollingFooter(text string, width int, offset int) string {
	// Reserve some space for padding and borders
	availableWidth := width - 4
	if availableWidth < 1 {
		availableWidth = 40 // Fallback minimum
	}

	// Convert to runes for proper unicode handling
	textRunes := []rune(text)
	textLen := len(textRunes)

	// If text fits, no scrolling needed
	if textLen <= availableWidth {
		return text
	}

	// Add padding to create smooth loop (as runes)
	separatorRunes := []rune("   •   ")
	paddedRunes := append(textRunes, separatorRunes...)
	paddedRunes = append(paddedRunes, textRunes...)

	// Calculate scroll position with wrapping
	scrollPos := offset % len(paddedRunes)

	// Extract visible portion
	var result []rune
	for i := 0; i < availableWidth; i++ {
		charPos := (scrollPos + i) % len(paddedRunes)
		result = append(result, paddedRunes[charPos])
	}

	return string(result)
}

// loadConfig loads the configuration from disk
func loadConfig() tea.Msg {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return configLoadedMsg{err: err}
	}

	configPath := filepath.Join(homeDir, ".config", "tui-launcher", "config.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return configLoadedMsg{err: err}
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return configLoadedMsg{err: err}
	}

	return configLoadedMsg{
		config: config,
		err:    nil,
	}
}

// detectTerminal detects the terminal emulator type
func detectTerminal() terminalType {
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")

	switch {
	case termProgram == "WezTerm":
		return terminalWezTerm
	case termProgram == "iTerm.app":
		return terminalITerm2
	case term == "xterm-kitty":
		return terminalKitty
	case term == "xterm-256color":
		return terminalXterm
	case os.Getenv("PREFIX") == "/data/data/com.termux/files/usr":
		return terminalTermux
	case os.Getenv("WT_SESSION") != "":
		return terminalWindowsTerminal
	default:
		return terminalUnknown
	}
}

// isInsideTmux checks if we're running inside a tmux session
func isInsideTmux() bool {
	return os.Getenv("TMUX") != ""
}

// runCommandDirectly runs a command directly in the current terminal (non-tmux mode)
func runCommandDirectly(item launchItem) tea.Cmd {
	return func() tea.Msg {
		// Change to working directory if specified
		if item.Cwd != "" {
			if err := os.Chdir(item.Cwd); err != nil {
				fmt.Printf("Error changing directory to %s: %v\n", item.Cwd, err)
				os.Exit(1)
			}
		}

		// Print what we're running
		fmt.Printf("Running: %s\n", item.Command)

		// Execute command using shell
		cmd := exec.Command("sh", "-c", item.Command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Command failed: %v\n", err)
			os.Exit(1)
		}

		os.Exit(0)
		return nil
	}
}

// editConfig opens the config file in the user's editor
func editConfig() tea.Msg {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Press Enter to continue...\n")
		fmt.Scanln()
		os.Exit(1)
	}

	configPath := filepath.Join(homeDir, ".config", "tui-launcher", "config.yaml")

	// Get editor from environment, fallback to sensible defaults
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Check which editors are available
		for _, e := range []string{"nano", "vim", "vi", "micro"} {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
	}
	if editor == "" {
		fmt.Printf("No editor found. Set $EDITOR or install nano/vim/vi\n")
		fmt.Printf("Press Enter to continue...\n")
		fmt.Scanln()
		os.Exit(1)
	}

	// Small delay to ensure terminal is restored
	time.Sleep(100 * time.Millisecond)

	// Open editor
	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Editor failed: %v\n", err)
		fmt.Printf("Press Enter to continue...\n")
		fmt.Scanln()
		os.Exit(1)
	}

	fmt.Printf("\nConfig saved. Press Enter to restart launcher...\n")
	fmt.Scanln()

	// Restart the launcher
	cmd = exec.Command(os.Args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	os.Exit(0)
	return nil
}

// updateInfoPane updates the info pane content based on the currently selected item
func (m *model) updateInfoPane() {
	mode := m.getLayoutMode()
	var currentItem launchItem
	var isValid bool

	// Get current item based on active pane and mode
	if mode == layoutDesktop {
		if m.activePane == paneGlobal {
			if m.globalCursor < len(m.globalTreeItems) {
				currentItem = m.globalTreeItems[m.globalCursor].item
				isValid = true
			}
		} else {
			if m.projectCursor < len(m.projectTreeItems) {
				currentItem = m.projectTreeItems[m.projectCursor].item
				isValid = true
			}
		}
	} else {
		// Compact/mobile mode
		if m.showingProjects {
			if m.projectCursor < len(m.projectTreeItems) {
				currentItem = m.projectTreeItems[m.projectCursor].item
				isValid = true
			}
		} else {
			if m.globalCursor < len(m.globalTreeItems) {
				currentItem = m.globalTreeItems[m.globalCursor].item
				isValid = true
			}
		}
	}

	if !isValid {
		m.infoContent = ""
		return
	}

	// Build info content based on item type
	var info strings.Builder

	// Name and icon
	if currentItem.Icon != "" {
		info.WriteString(currentItem.Icon + " ")
	}
	info.WriteString(currentItem.Name + "\n")
	info.WriteString(strings.Repeat("─", len(currentItem.Name)+2) + "\n\n")

	// Type-specific info
	switch currentItem.ItemType {
	case typeCategory:
		info.WriteString(fmt.Sprintf("Type: Category\n"))
		info.WriteString(fmt.Sprintf("Children: %d items\n", len(currentItem.Children)))

	case typeCommand:
		info.WriteString(fmt.Sprintf("Type: Command\n"))
		if currentItem.Command != "" {
			info.WriteString(fmt.Sprintf("Command: %s\n", currentItem.Command))
		}
		if currentItem.Cwd != "" {
			info.WriteString(fmt.Sprintf("Working Dir: %s\n", currentItem.Cwd))
		}
		if currentItem.SpawnStr != "" {
			info.WriteString(fmt.Sprintf("Spawn Mode: %s\n", currentItem.SpawnStr))
		}

	case typeProfile:
		info.WriteString(fmt.Sprintf("Type: Profile\n"))
		info.WriteString(fmt.Sprintf("Layout: %s\n", currentItem.LayoutStr))
		info.WriteString(fmt.Sprintf("Panes: %d\n", len(currentItem.Panes)))
		if len(currentItem.Panes) > 0 {
			info.WriteString("\nPane Commands:\n")
			for i, pane := range currentItem.Panes {
				info.WriteString(fmt.Sprintf("  %d. %s\n", i+1, pane.Command))
			}
		}
	}

	m.infoContent = info.String()
}

// getLayoutMode determines which responsive layout to use based on terminal size
func (m model) getLayoutMode() layoutMode {
	// Mobile mode: Very small height (Termux with keyboard open)
	if m.height <= 12 {
		return layoutMobile
	}

	// Compact mode: Narrow terminal (portrait or phone)
	if m.width < 80 {
		return layoutCompact
	}

	// Desktop mode: Normal terminal
	return layoutDesktop
}

// calculateLayout returns pane dimensions based on current layout mode
// Returns (leftWidth, rightWidth, treeHeight, infoHeight)
// Follows Golden Rule #1: Always subtract borders from height calculations
// Follows Golden Rule #4: Use proportional sizing (weights), not pixels
func (m model) calculateLayout() (int, int, int, int) {
	// Start with full dimensions
	contentHeight := m.height
	contentWidth := m.width

	// Subtract header (3 lines: title + cwd + blank)
	contentHeight -= 3

	// Subtract footer area (2 lines: status line + footer)
	contentHeight -= 2

	// GOLDEN RULE #1: Account for borders BEFORE rendering
	// Subtract 2 for panel borders (top + bottom)
	contentHeight -= 2

	// Get current layout mode
	mode := m.getLayoutMode()

	switch mode {
	case layoutDesktop:
		// 3-pane: Left | Right | Bottom
		// Split width proportionally (50/50)
		leftWidth := contentWidth / 2
		rightWidth := contentWidth - leftWidth

		// Split height (2/3 for trees, 1/3 for info)
		treeHeight := (contentHeight * 2) / 3
		infoHeight := contentHeight - treeHeight

		return leftWidth, rightWidth, treeHeight, infoHeight

	case layoutCompact:
		// 2-pane: Combined tree on top | Info on bottom
		// Use full width for tree
		treeHeight := (contentHeight * 3) / 4
		infoHeight := contentHeight - treeHeight

		return contentWidth, 0, treeHeight, infoHeight

	case layoutMobile:
		// 1-pane: Just tree, toggle info with 'i' key
		// Use all available space for tree
		return contentWidth, 0, contentHeight, 0

	default:
		// Fallback to desktop
		leftWidth := contentWidth / 2
		rightWidth := contentWidth - leftWidth
		treeHeight := (contentHeight * 2) / 3
		infoHeight := contentHeight - treeHeight
		return leftWidth, rightWidth, treeHeight, infoHeight
	}
}

// editConfigInTmux opens the config file in a tmux split
func editConfigInTmux() tea.Cmd {
	return func() tea.Msg {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return spawnCompleteMsg{err: err}
		}

		configPath := filepath.Join(homeDir, ".config", "tui-launcher", "config.yaml")

		// Get editor from environment, fallback to sensible defaults
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = os.Getenv("VISUAL")
		}
		if editor == "" {
			// Check which editors are available
			for _, e := range []string{"nano", "vim", "vi", "micro"} {
				if _, err := exec.LookPath(e); err == nil {
					editor = e
					break
				}
			}
		}
		if editor == "" {
			return spawnCompleteMsg{err: fmt.Errorf("no editor found. Set $EDITOR or install nano/vim/vi")}
		}

		// Open in tmux split
		cmd := exec.Command("tmux", "split-window", "-h", editor, configPath)
		if err := cmd.Run(); err != nil {
			return spawnCompleteMsg{err: fmt.Errorf("failed to open editor: %w", err)}
		}

		return spawnCompleteMsg{err: nil}
	}
}
