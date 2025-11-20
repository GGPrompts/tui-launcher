package launch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"

	"tui-launcher/shared"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.MouseMsg:
		return m.handleMouseEvent(msg)

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

			// Update info pane for initial selection
			m.updateInfoPane()
		} else {
			// Debug: show error if config failed to load
			fmt.Fprintf(os.Stderr, "Config load error: %v\n", msg.err)
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

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
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
		return m.navigateUp(), nil

	case "down", "j":
		return m.navigateDown(), nil

	case "right", "l":
		return m.expandCategory(), nil

	case "left", "h":
		return m.collapseCategory(), nil

	case " ":
		return m.handleSpaceKey(), nil

	case "c":
		// Clear all selections
		m.selectedItems = make(map[string]bool)

	case "d":
		// Toggle detached mode (spawn in tmux vs foreground)
		m.detachedMode = !m.detachedMode

	case "e":
		// Edit config file
		if m.insideTmux {
			// If inside tmux, spawn editor in a new split
			return m, editConfigInTmux()
		} else {
			// Use ExecProcess to properly handle the editor
			return m, editConfigExec()
		}

	case "enter":
		return m.handleEnterKey()
	}

	return m, nil
}

// handleMouseEvent handles mouse input
func (m Model) handleMouseEvent(msg tea.MouseMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.MouseLeft:
		// Click to switch panes in desktop mode
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
		return m.navigateUp(), nil

	case tea.MouseWheelDown:
		return m.navigateDown(), nil
	}

	return m, nil
}

// navigateUp moves cursor up in active pane
func (m Model) navigateUp() Model {
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
	m.updateInfoPane()
	return m
}

// navigateDown moves cursor down in active pane
func (m Model) navigateDown() Model {
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
	m.updateInfoPane()
	return m
}

// expandCategory expands the current category
func (m Model) expandCategory() Model {
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

	return m
}

// collapseCategory collapses the current category
func (m Model) collapseCategory() Model {
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

	return m
}

// handleSpaceKey handles space key (expand/collapse or select)
func (m Model) handleSpaceKey() Model {
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

	return m
}

// handleEnterKey handles enter key (launch or CD)
func (m Model) handleEnterKey() (Model, tea.Cmd) {
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
		// Special handling for project categories - CD into them
		if currentItem.ItemType == typeCategory && currentItem.Cwd != "" {
			// This is a project category with a directory - CD into it
			if err := writeCDTarget(currentItem.Cwd); err != nil {
				m.err = fmt.Errorf("failed to write CD target: %w", err)
				return m, nil
			}
			return m, tea.Quit
		}

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
				// Multi-select
				if m.detachedMode {
					// Detached: Spawn each in background, stay in launcher
					return m, spawnMultipleDetached(itemsToLaunch)
				} else {
					// Foreground: Spawn each as tmux window, exit launcher
					// User lands in tmux with multiple windows to switch between
					return m, spawnMultipleForeground(itemsToLaunch)
				}
			}

		} else {
			// No selection - launch current item if it's a command or profile
			if currentItem.ItemType == typeCommand {
				// Check if detached mode is enabled
				if m.detachedMode {
					// Spawn in tmux window (background/detached)
					return m, spawnInTmuxWindow(currentItem)
				} else {
					// Launch in foreground (like TFE)
					return m, runCommandDirectly(currentItem)
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

	return m, nil
}

// writeCDTarget writes the target directory to a file so the shell wrapper can cd after exit
func writeCDTarget(path string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	targetFile := filepath.Join(homeDir, ".tui-launcher_cd_target")
	return os.WriteFile(targetFile, []byte(path), 0644)
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

// editConfigExec returns a command to edit the config file using tea.ExecProcess
func editConfigExec() tea.Cmd {
	return func() tea.Msg {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil
		}

		configPath := filepath.Join(homeDir, ".config", "tui-launcher", "config.yaml")

		// Get editor from environment, fallback to sensible defaults
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = os.Getenv("VISUAL")
		}
		if editor == "" {
			// Check which editors are available (micro is first since it's user-friendly)
			for _, e := range []string{"micro", "nano", "vim", "vi"} {
				if _, err := exec.LookPath(e); err == nil {
					editor = e
					break
				}
			}
		}
		if editor == "" {
			// No editor found
			return nil
		}

		c := exec.Command(editor, configPath)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		// Use TFE's proven pattern: ClearScreen + ExecProcess + Sequence()
		return tea.Sequence(
			tea.ClearScreen,
			tea.ExecProcess(c, func(err error) tea.Msg {
				// After editor exits, restart the launcher
				if err != nil {
					return tea.Quit()
				}
				// Restart the launcher
				restart := exec.Command(os.Args[0])
				restart.Stdin = os.Stdin
				restart.Stdout = os.Stdout
				restart.Stderr = os.Stderr
				restart.Run()
				os.Exit(0)
				return nil
			}),
		)()
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

// runCommandDirectly runs a command directly in the current terminal using tea.ExecProcess
func runCommandDirectly(item launchItem) tea.Cmd {
	return func() tea.Msg {
		// Build the command with working directory handling
		cmdStr := item.Command
		if item.Cwd != "" {
			cmdStr = fmt.Sprintf("cd %s && %s", shellescape(item.Cwd), item.Command)
		}

		c := exec.Command("sh", "-c", cmdStr)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		// Use TFE's pattern: ClearScreen + ExecProcess
		return tea.Sequence(
			tea.ClearScreen,
			tea.ExecProcess(c, func(err error) tea.Msg {
				return tea.Quit()
			}),
		)()
	}
}

// shellescape escapes a string for safe use in shell commands
func shellescape(s string) string {
	// Simple quote escaping - wrap in single quotes and escape any single quotes
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// spawnInTmuxWindow spawns a command in a new tmux window (detached/background)
func spawnInTmuxWindow(item launchItem) tea.Cmd {
	return func() tea.Msg {
		// Build command with working directory
		cmdStr := item.Command
		if item.Cwd != "" {
			cmdStr = fmt.Sprintf("cd %s && %s", shellescape(item.Cwd), item.Command)
		}

		// Create tmux window with item name as window name
		// -d flag = don't switch to the new window (stay in launcher)
		windowName := item.Name
		cmd := exec.Command("tmux", "new-window", "-d", "-n", windowName, "sh", "-c", cmdStr)

		if err := cmd.Run(); err != nil {
			return spawnCompleteMsg{err: err}
		}

		// Don't quit - stay in launcher so user can spawn more
		return nil
	}
}

// spawnMultipleForeground spawns multiple commands as tmux windows and exits launcher
func spawnMultipleForeground(items []launchItem) tea.Cmd {
	return func() tea.Msg {
		fmt.Fprintf(os.Stderr, "DEBUG: spawnMultipleForeground called with %d items\n", len(items))

		for i, item := range items {
			fmt.Fprintf(os.Stderr, "DEBUG: Spawning item %d: %s\n", i+1, item.Name)

			// Build command with working directory
			cmdStr := item.Command
			if item.Cwd != "" {
				cmdStr = fmt.Sprintf("cd %s && %s", shellescape(item.Cwd), item.Command)
			}

			// Create tmux window with item name
			// No -d flag = switches to each new window
			windowName := item.Name
			cmd := exec.Command("tmux", "new-window", "-n", windowName, "sh", "-c", cmdStr)

			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "DEBUG: Failed to spawn %s: %v\n", item.Name, err)
				continue
			}

			fmt.Fprintf(os.Stderr, "DEBUG: Successfully spawned %s\n", item.Name)
			time.Sleep(10 * time.Millisecond)
		}

		fmt.Fprintf(os.Stderr, "DEBUG: Finished spawning, now quitting launcher\n")

		// Quit launcher - user is now in tmux with multiple windows
		return tea.Quit()
	}
}

// spawnMultipleDetached spawns multiple commands in separate tmux windows
func spawnMultipleDetached(items []launchItem) tea.Cmd {
	return func() tea.Msg {
		fmt.Fprintf(os.Stderr, "DEBUG: spawnMultipleDetached called with %d items\n", len(items))

		for i, item := range items {
			fmt.Fprintf(os.Stderr, "DEBUG: Spawning item %d: %s\n", i+1, item.Name)

			// Build command with working directory
			cmdStr := item.Command
			if item.Cwd != "" {
				cmdStr = fmt.Sprintf("cd %s && %s", shellescape(item.Cwd), item.Command)
			}

			// Create tmux window with item name
			// -d flag = don't switch to the new window (stay in launcher)
			windowName := item.Name
			cmd := exec.Command("tmux", "new-window", "-d", "-n", windowName, "sh", "-c", cmdStr)

			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "DEBUG: Failed to spawn %s: %v\n", item.Name, err)
				// Continue spawning others even if one fails
				continue
			}

			fmt.Fprintf(os.Stderr, "DEBUG: Successfully spawned %s\n", item.Name)

			// Small delay between spawns
			time.Sleep(10 * time.Millisecond)
		}

		fmt.Fprintf(os.Stderr, "DEBUG: Finished spawning all items\n")

		// Don't quit - stay in launcher
		return nil
	}
}

// spawnSingle wraps shared.SpawnSingle to convert types
func spawnSingle(item launchItem, mode spawnMode) tea.Cmd {
	// Convert to shared types
	sharedItem := shared.LaunchItem{
		Name:    item.Name,
		Command: item.Command,
		Cwd:     item.Cwd,
	}

	// Convert spawn mode
	var sharedMode shared.SpawnMode
	switch mode {
	case spawnXtermWindow:
		sharedMode = shared.SpawnXtermWindow
	case spawnTmuxWindow:
		sharedMode = shared.SpawnTmuxWindow
	case spawnTmuxSplitH:
		sharedMode = shared.SpawnTmuxSplitH
	case spawnTmuxSplitV:
		sharedMode = shared.SpawnTmuxSplitV
	case spawnTmuxLayout:
		sharedMode = shared.SpawnTmuxLayout
	case spawnCurrentPane:
		sharedMode = shared.SpawnCurrentPane
	}

	return shared.SpawnSingle(sharedItem, sharedMode)
}

func spawnMultiple(items []launchItem, layout tmuxLayout) tea.Cmd {
	// Convert to shared types
	sharedItems := make([]shared.LaunchItem, len(items))
	for i, item := range items {
		sharedItems[i] = shared.LaunchItem{
			Name:    item.Name,
			Command: item.Command,
			Cwd:     item.Cwd,
		}
	}

	// Convert layout
	var sharedLayout shared.TmuxLayout
	switch layout {
	case layoutMainVertical:
		sharedLayout = shared.LayoutMainVertical
	case layoutMainHorizontal:
		sharedLayout = shared.LayoutMainHorizontal
	case layoutTiled:
		sharedLayout = shared.LayoutTiled
	case layoutEvenHorizontal:
		sharedLayout = shared.LayoutEvenHorizontal
	case layoutEvenVertical:
		sharedLayout = shared.LayoutEvenVertical
	}

	return shared.SpawnMultiple(sharedItems, sharedLayout)
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return strings.Replace(path, "~", homeDir, 1)
	}
	return path
}
