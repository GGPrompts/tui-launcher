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
		footerTick(),
	)
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.treeItems)-1 {
				m.cursor++
			}

		case "right", "l":
			// Expand current category
			if m.cursor < len(m.treeItems) {
				currentItem := m.treeItems[m.cursor].item
				if currentItem.ItemType == typeCategory {
					m.expandedItems[currentItem.Path] = true
					m.treeItems = flattenTree(m.rootItems, m.expandedItems)
				}
			}

		case "left", "h":
			// Collapse current category
			if m.cursor < len(m.treeItems) {
				currentItem := m.treeItems[m.cursor].item
				if currentItem.ItemType == typeCategory {
					m.expandedItems[currentItem.Path] = false
					m.treeItems = flattenTree(m.rootItems, m.expandedItems)
				}
			}

		case " ":
			// Context-aware action: expand/collapse categories, or toggle selection for commands
			if m.cursor < len(m.treeItems) {
				currentItem := m.treeItems[m.cursor].item

				if currentItem.ItemType == typeCategory {
					// Toggle expansion for categories
					if m.expandedItems[currentItem.Path] {
						m.expandedItems[currentItem.Path] = false
					} else {
						m.expandedItems[currentItem.Path] = true
					}
					m.treeItems = flattenTree(m.rootItems, m.expandedItems)

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

		case "enter":
			// Launch selected items or current item
			if m.cursor < len(m.treeItems) {
				currentItem := m.treeItems[m.cursor].item

				if len(m.selectedItems) > 0 {
					// Launch all selected items in batch
					var itemsToLaunch []launchItem

					// Collect all selected items
					for _, ti := range m.treeItems {
						if m.selectedItems[ti.item.Path] {
							itemsToLaunch = append(itemsToLaunch, ti.item)
						}
					}

					if len(itemsToLaunch) > 0 {
						// Use default layout for batch launch
						return m, spawnMultiple(itemsToLaunch, m.selectedLayout)
					}

				} else if currentItem.ItemType == typeCommand {
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

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			// Scroll up
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.MouseWheelDown:
			// Scroll down
			if m.cursor < len(m.treeItems)-1 {
				m.cursor++
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case configLoadedMsg:
		m.loading = false
		m.config = msg.config
		m.err = msg.err
		if msg.err == nil {
			// Build tree from config
			m.rootItems = buildTreeFromConfig(msg.config)
			m.treeItems = flattenTree(m.rootItems, m.expandedItems)
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

	case footerTickMsg:
		// Increment footer scroll offset
		m.footerOffset++
		return m, footerTick()
	}

	return m, nil
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

	// Header
	sb.WriteString("🚀 TUI Launcher v" + Version + "\n")

	// Show current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "unknown"
	}
	// Shorten home directory to ~
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

	// Render tree
	if len(m.treeItems) == 0 {
		sb.WriteString("No items configured.\n")
	} else {
		for i, ti := range m.treeItems {
			selected := m.selectedItems[ti.item.Path]
			expanded := m.expandedItems[ti.item.Path]
			line := renderTreeItem(ti, m.cursor, i, selected, expanded)

			// Make the cursor line bold
			if i == m.cursor {
				line = selectedStyle.Render(line)
			}

			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}

	// Status line
	sb.WriteString("\n")
	if len(m.selectedItems) > 0 {
		sb.WriteString(fmt.Sprintf("Selected: %d items | ", len(m.selectedItems)))
	}

	// Footer
	footerText := "↑/↓: navigate  Space: expand/select  c: clear  t: toggle mode  Enter: launch  q: quit"
	sb.WriteString(renderScrollingFooter(footerText, m.width, m.footerOffset))
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
