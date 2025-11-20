package main

import (
	"github.com/charmbracelet/bubbles/spinner"
)

// types_unified.go - Unified Type Definitions for Tab-Based Architecture
// This file extends the current types.go with tab routing support

// Tab names for the unified interface
type tabName string

const (
	tabLaunch    tabName = "launch"
	tabSessions  tabName = "sessions"
	tabTemplates tabName = "templates"
)

// String returns the display name for a tab
func (t tabName) String() string {
	switch t {
	case tabLaunch:
		return "Launch"
	case tabSessions:
		return "Sessions"
	case tabTemplates:
		return "Templates"
	default:
		return "Unknown"
	}
}

// --- Unified Model ---
// This will replace the current model struct once tabs are implemented

// unifiedModel represents the top-level application state with tab routing
type unifiedModel struct {
	// Display dimensions (shared across all tabs)
	width  int
	height int

	// Tab routing
	currentTab tabName

	// Error handling (shared)
	err       error
	statusMsg string

	// Popup mode (from tmuxplexer)
	popupMode bool

	// Tab-specific models
	// Note: launch.Model is imported from tabs/launch package
	// The import is added in model_unified.go to avoid circular dependencies
	launchModel interface{} // Will be launch.Model at runtime

	// Sessions tab state (from tmuxplexer) - TBD in implementation
	// sessionsModel sessionsTabModel

	// Templates tab state (from tmuxplexer) - TBD in implementation
	// templatesModel templatesTabModel

	// Shared spinner for loading states
	spinner spinner.Model
	loading bool
}

// --- Launch Tab Model ---
// This mirrors the current model struct for the Launch tab

type launchTabModel struct {
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

	// Terminal detection
	terminalType terminalType
	insideTmux   bool
	useTmux      bool
}

// --- Sessions Tab Model (Placeholder) ---
// Will be populated from tmuxplexer during implementation

type sessionsTabModel struct {
	// Placeholder - will be filled from tmuxplexer-original/types.go
	sessions []TmuxSession
	// ... other fields TBD
}

// --- Templates Tab Model (Placeholder) ---
// Will be populated from tmuxplexer during implementation

type templatesTabModel struct {
	// Placeholder - will be filled from tmuxplexer-original/types.go
	templates []SessionTemplate
	// ... other fields TBD
}

// --- Tmux Types (from tmuxplexer) ---
// These are needed for Sessions and Templates tabs

// TmuxSession represents a tmux session (from tmuxplexer)
type TmuxSession struct {
	Name       string
	Windows    int
	Attached   bool
	Created    string
	LastActive string
	WorkingDir string
	GitBranch  string
	// ClaudeState *ClaudeState // TBD - will add when implementing Sessions tab
	AITool string // "claude", "codex", "gemini", or ""
}

// TmuxWindow represents a window in a session (from tmuxplexer)
type TmuxWindow struct {
	Index  int
	Name   string
	Panes  int
	Active bool
}

// TmuxPane represents a pane in a window (from tmuxplexer)
type TmuxPane struct {
	ID      string
	Index   int
	Width   int
	Height  int
	Active  bool
	Command string
}

// SessionTemplate represents a template for creating tmux sessions (from tmuxplexer)
type SessionTemplate struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    string            `json:"category,omitempty"`
	WorkingDir  string            `json:"working_dir"`
	Layout      string            `json:"layout"` // "2x2", "4x2", "3x3", etc.
	Panes       []PaneTemplate    `json:"panes"`
	Env         map[string]string `json:"env,omitempty"`
}

// PaneTemplate represents a pane configuration in a template (from tmuxplexer)
type PaneTemplate struct {
	Command    string `json:"command"`
	Title      string `json:"title,omitempty"`
	WorkingDir string `json:"working_dir,omitempty"`
}

// --- Helper Functions ---

// nextTab returns the next tab in the cycle
func nextTab(current tabName) tabName {
	switch current {
	case tabLaunch:
		return tabSessions
	case tabSessions:
		return tabTemplates
	case tabTemplates:
		return tabLaunch
	default:
		return tabLaunch
	}
}

// prevTab returns the previous tab in the cycle
func prevTab(current tabName) tabName {
	switch current {
	case tabLaunch:
		return tabTemplates
	case tabSessions:
		return tabLaunch
	case tabTemplates:
		return tabSessions
	default:
		return tabLaunch
	}
}
