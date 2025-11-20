package launch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

// Import types from main package (these will be used from shared/ later)
type paneType int
type itemType int
type spawnMode int
type tmuxLayout int
type terminalType int
type layoutMode int

// Constants from main package
const (
	paneGlobal paneType = iota
	paneProject
)

const (
	typeCategory itemType = iota
	typeCommand
	typeProfile
)

const (
	spawnXtermWindow spawnMode = iota
	spawnTmuxWindow
	spawnTmuxSplitH
	spawnTmuxSplitV
	spawnTmuxLayout
	spawnCurrentPane
)

const (
	layoutMainVertical tmuxLayout = iota
	layoutMainHorizontal
	layoutTiled
	layoutEvenHorizontal
	layoutEvenVertical
)

const (
	terminalUnknown terminalType = iota
	terminalWindowsTerminal
	terminalWezTerm
	terminalKitty
	terminalITerm2
	terminalXterm
	terminalTermux
)

const (
	layoutDesktop layoutMode = iota
	layoutCompact
	layoutMobile
)

// Imported types
type paneConfig struct {
	Command string `yaml:"command"`
	Cwd     string `yaml:"cwd"`
}

type paneInfo struct {
	description string
	cliFlags    string
	repo        string
	mdPath      string
}

type launchItem struct {
	Name         string        `yaml:"name"`
	Path         string        `yaml:"-"`
	ItemType     itemType      `yaml:"-"`
	Icon         string        `yaml:"icon"`
	Command      string        `yaml:"command"`
	Cwd          string        `yaml:"cwd"`
	DefaultSpawn spawnMode     `yaml:"-"`
	SpawnStr     string        `yaml:"spawn"`
	Children     []launchItem  `yaml:"items"`
	IsProfile    bool          `yaml:"-"`
	Layout       tmuxLayout    `yaml:"-"`
	LayoutStr    string        `yaml:"layout"`
	Panes        []paneConfig  `yaml:"panes"`
}

type launchTreeItem struct {
	item        launchItem
	depth       int
	isLast      bool
	parentLasts []bool
}

type Config struct {
	Projects []ProjectConfig  `yaml:"projects"`
	Tools    []CategoryConfig `yaml:"tools"`
	AI       []CommandConfig  `yaml:"ai"`
	Scripts  []CategoryConfig `yaml:"scripts"`
}

type ProjectConfig struct {
	Name     string          `yaml:"name"`
	Icon     string          `yaml:"icon"`
	Path     string          `yaml:"path"`
	Commands []CommandConfig `yaml:"commands"`
	Profiles []ProfileConfig `yaml:"profiles"`
}

type CategoryConfig struct {
	Category string          `yaml:"category"`
	Icon     string          `yaml:"icon"`
	Items    []CommandConfig `yaml:"items"`
}

type CommandConfig struct {
	Name    string `yaml:"name"`
	Icon    string `yaml:"icon"`
	Command string `yaml:"command"`
	Cwd     string `yaml:"cwd"`
	Spawn   string `yaml:"spawn"`
}

type ProfileConfig struct {
	Name   string       `yaml:"name"`
	Icon   string       `yaml:"icon"`
	Layout string       `yaml:"layout"`
	Panes  []paneConfig `yaml:"panes"`
}

// Messages
type configLoadedMsg struct {
	config Config
	err    error
}

type spawnCompleteMsg struct {
	err error
}

// Model represents the Launch tab state
type Model struct {
	// Display dimensions
	width  int
	height int

	// Multi-pane layout state
	activePane       paneType
	globalItems      []launchItem
	projectItems     []launchItem
	globalTreeItems  []launchTreeItem
	projectTreeItems []launchTreeItem
	globalCursor     int
	projectCursor    int
	globalExpanded   map[string]bool
	projectExpanded  map[string]bool

	// Info pane state
	currentInfo     paneInfo
	infoContent     string
	showingInfo     bool
	showingProjects bool

	// Selection state
	selectedItems map[string]bool

	// Spawn dialog state
	showSpawnDialog bool
	selectedLayout  tmuxLayout
	layoutCursor    int

	// Config
	config Config

	// UI state
	spinner spinner.Model
	loading bool
	err     error

	// Terminal detection
	terminalType terminalType
	insideTmux   bool

	// Spawn mode toggle
	detachedMode bool // When true, spawn in tmux windows (background)
}

// New creates a new Launch tab model
func New(width, height int) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		width:  width,
		height: height,

		// Multi-pane initialization
		activePane:      paneGlobal,
		globalItems:     []launchItem{},
		projectItems:    []launchItem{},
		globalTreeItems: []launchTreeItem{},
		projectTreeItems: []launchTreeItem{},
		globalCursor:     0,
		projectCursor:    0,
		globalExpanded:   make(map[string]bool),
		projectExpanded:  make(map[string]bool),

		// Info pane initialization
		showingInfo:     false,
		showingProjects: false,

		selectedItems:   make(map[string]bool),
		showSpawnDialog: false,
		selectedLayout:  layoutTiled,
		layoutCursor:    0,
		spinner:         s,
		loading:         true,
		terminalType:    detectTerminal(),
		insideTmux:      isInsideTmux(),
		detachedMode:    false, // Default to foreground mode
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadConfig,
	)
}

// loadConfig loads the configuration from disk
func loadConfig() tea.Msg {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Error getting home dir: %v\n", err)
		return configLoadedMsg{err: err}
	}

	configPath := filepath.Join(homeDir, ".config", "tui-launcher", "config.yaml")
	fmt.Fprintf(os.Stderr, "DEBUG: Loading config from: %s\n", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Error reading config: %v\n", err)
		return configLoadedMsg{err: err}
	}

	fmt.Fprintf(os.Stderr, "DEBUG: Config file read, size: %d bytes\n", len(data))

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Error parsing YAML: %v\n", err)
		return configLoadedMsg{err: err}
	}

	fmt.Fprintf(os.Stderr, "DEBUG: Config parsed successfully - Projects: %d, Tools: %d\n", len(config.Projects), len(config.Tools))

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

// getLayoutMode determines which responsive layout to use based on terminal size
func (m Model) getLayoutMode() layoutMode {
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
func (m Model) calculateLayout() (int, int, int, int) {
	// Start with full dimensions
	contentHeight := m.height
	contentWidth := m.width

	// Subtract header (3 lines: title + cwd + blank)
	contentHeight -= 3

	// Subtract footer area (2 lines: status line + footer)
	contentHeight -= 2

	// Account for borders (top + bottom)
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
		treeHeight := (contentHeight * 3) / 4
		infoHeight := contentHeight - treeHeight

		return contentWidth, 0, treeHeight, infoHeight

	case layoutMobile:
		// 1-pane: Just tree, toggle info with 'i' key
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

// updateInfoPane updates the info pane content based on the currently selected item
func (m *Model) updateInfoPane() {
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
		if currentItem.Cwd != "" {
			info.WriteString(fmt.Sprintf("\n📂 Project Directory:\n%s\n", currentItem.Cwd))
			info.WriteString("\n💡 Press Enter to CD into this project\n")
		}

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
