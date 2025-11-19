package main

import (
	"github.com/charmbracelet/bubbles/spinner"
)

// Version is the current version of tui-launcher
const Version = "1.0.0"

// itemType represents different types of launch items
type itemType int

const (
	typeCategory itemType = iota // Expandable folder/category
	typeCommand                   // Single executable command
	typeProfile                   // Multi-launch configuration
)

func (t itemType) String() string {
	switch t {
	case typeCategory:
		return "Category"
	case typeCommand:
		return "Command"
	case typeProfile:
		return "Profile"
	default:
		return "Unknown"
	}
}

// spawnMode represents different ways to spawn commands
type spawnMode int

const (
	spawnXtermWindow spawnMode = iota // New xterm window
	spawnTmuxWindow                   // New tmux window
	spawnTmuxSplitH                   // Tmux split horizontal
	spawnTmuxSplitV                   // Tmux split vertical
	spawnTmuxLayout                   // Tmux with custom layout
	spawnCurrentPane                  // Replace current pane
)

func (s spawnMode) String() string {
	switch s {
	case spawnXtermWindow:
		return "XTerm Window"
	case spawnTmuxWindow:
		return "Tmux Window"
	case spawnTmuxSplitH:
		return "Tmux Split Horizontal"
	case spawnTmuxSplitV:
		return "Tmux Split Vertical"
	case spawnTmuxLayout:
		return "Tmux Layout"
	case spawnCurrentPane:
		return "Current Pane"
	default:
		return "Unknown"
	}
}

// tmuxLayout represents tmux layout types
type tmuxLayout int

const (
	layoutMainVertical tmuxLayout = iota // Main pane left, others stacked right
	layoutMainHorizontal                 // Main pane top, others stacked below
	layoutTiled                          // Grid layout
	layoutEvenHorizontal                 // Equal width columns
	layoutEvenVertical                   // Equal height rows
)

func (l tmuxLayout) String() string {
	switch l {
	case layoutMainVertical:
		return "main-vertical"
	case layoutMainHorizontal:
		return "main-horizontal"
	case layoutTiled:
		return "tiled"
	case layoutEvenHorizontal:
		return "even-horizontal"
	case layoutEvenVertical:
		return "even-vertical"
	default:
		return "main-vertical"
	}
}

// terminalType represents different terminal emulators with varying emoji rendering
type terminalType int

const (
	terminalUnknown terminalType = iota
	terminalWindowsTerminal
	terminalWezTerm
	terminalKitty
	terminalITerm2
	terminalXterm
	terminalTermux
)

func (t terminalType) String() string {
	switch t {
	case terminalWindowsTerminal:
		return "Windows Terminal"
	case terminalWezTerm:
		return "WezTerm"
	case terminalKitty:
		return "Kitty"
	case terminalITerm2:
		return "iTerm2"
	case terminalXterm:
		return "xterm"
	case terminalTermux:
		return "Termux"
	default:
		return "Unknown"
	}
}

// Emoji constants (NO variation selectors - U+FE0F causes width bugs!)
const (
	// Spawn mode emojis
	emojiXtermWindow  = "🪟" // U+1FA9F
	emojiTmuxWindow   = "🔲" // U+1F532
	emojiTmuxSplitH   = "⬛" // U+2B1B
	emojiTmuxSplitV   = "⬜" // U+2B1C
	emojiTmuxLayout   = "📐" // U+1F4D0
	emojiCurrentPane  = "🎯" // U+1F3AF
	emojiBatchSpawn   = "🎛" // U+1F39B (NO U+FE0F!)

	// Item type icons
	emojiProject      = "📦" // U+1F4E6
	emojiFolder       = "📁" // U+1F4C1
	emojiCommand      = "⚡" // U+26A1
	emojiProfile      = "🔧" // U+1F527
	emojiTUITool      = "🛠" // U+1F6E0 (NO U+FE0F!)
	emojiAI           = "🤖" // U+1F916
	emojiScript       = "📜" // U+1F4DC
	emojiFavorite     = "⭐" // U+2B50
	emojiMonitoring   = "📊" // U+1F4CA
	emojiGit          = "🗄" // U+1F5C4 (NO U+FE0F!)
	emojiDatabase     = "💾" // U+1F4BE
	emojiServer       = "🖥" // U+1F5A5 (NO U+FE0F!)

	// Status indicators
	emojiExpanded     = "▼"  // U+25BC
	emojiCollapsed    = "▶"  // U+25B6
	emojiSelected     = "☑"  // U+2611
	emojiUnselected   = "☐"  // U+2610
	emojiRunning      = "🟢" // U+1F7E2
	emojiStopped      = "🔴" // U+1F534
)

// paneConfig represents a single pane in a profile
type paneConfig struct {
	Command string `yaml:"command"`
	Cwd     string `yaml:"cwd"`
}

// launchItem represents a command, category, or profile in the launcher
type launchItem struct {
	Name         string        `yaml:"name"`
	Path         string        `yaml:"-"` // Unique identifier (computed)
	ItemType     itemType      `yaml:"-"` // Determined from config structure
	Icon         string        `yaml:"icon"`
	Command      string        `yaml:"command"`
	Cwd          string        `yaml:"cwd"`
	DefaultSpawn spawnMode     `yaml:"-"` // Parsed from spawn string
	SpawnStr     string        `yaml:"spawn"` // String from config
	Children     []launchItem  `yaml:"items"`

	// For profiles
	IsProfile    bool          `yaml:"-"`
	Layout       tmuxLayout    `yaml:"-"` // Parsed from layout string
	LayoutStr    string        `yaml:"layout"` // String from config
	Panes        []paneConfig  `yaml:"panes"`
}

// launchTreeItem represents an item in the flattened tree view
// (Similar to TFE's treeItem)
type launchTreeItem struct {
	item        launchItem
	depth       int
	isLast      bool
	parentLasts []bool // Track which parent levels are last items
}

// model represents the application state
type model struct {
	// Display dimensions
	width  int
	height int

	// Tree navigation
	cursor        int
	treeItems     []launchTreeItem
	expandedItems map[string]bool

	// Selection state
	selectedItems map[string]bool

	// Spawn dialog state
	showSpawnDialog bool
	selectedLayout  tmuxLayout
	layoutCursor    int // For layout picker in dialog

	// Config
	config        Config

	// UI state
	spinner       spinner.Model
	loading       bool
	err           error

	// Terminal detection
	terminalType  terminalType
	insideTmux    bool
}

// Config represents the configuration file structure
type Config struct {
	Projects []ProjectConfig  `yaml:"projects"`
	Tools    []CategoryConfig `yaml:"tools"`
	AI       []CommandConfig  `yaml:"ai"`
	Scripts  []CategoryConfig `yaml:"scripts"`
}

// ProjectConfig represents a project with commands and profiles
type ProjectConfig struct {
	Name     string          `yaml:"name"`
	Icon     string          `yaml:"icon"`
	Path     string          `yaml:"path"`
	Commands []CommandConfig `yaml:"commands"`
	Profiles []ProfileConfig `yaml:"profiles"`
}

// CategoryConfig represents a category of commands
type CategoryConfig struct {
	Category string          `yaml:"category"`
	Icon     string          `yaml:"icon"`
	Items    []CommandConfig `yaml:"items"`
}

// CommandConfig represents a single command
type CommandConfig struct {
	Name    string `yaml:"name"`
	Icon    string `yaml:"icon"`
	Command string `yaml:"command"`
	Cwd     string `yaml:"cwd"`
	Spawn   string `yaml:"spawn"`
}

// ProfileConfig represents a multi-pane launch configuration
type ProfileConfig struct {
	Name   string       `yaml:"name"`
	Icon   string       `yaml:"icon"`
	Layout string       `yaml:"layout"`
	Panes  []paneConfig `yaml:"panes"`
}

// layoutOption represents a layout choice in the spawn dialog
type layoutOption struct {
	layout      tmuxLayout
	name        string
	description string
	preview     string // ASCII art preview
}

// spawnCompleteMsg is sent when spawning completes
type spawnCompleteMsg struct {
	err error
}

// configLoadedMsg is sent when config loads
type configLoadedMsg struct {
	config Config
	err    error
}
