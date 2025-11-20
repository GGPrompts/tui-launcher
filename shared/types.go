package shared

// shared/types.go - Shared Type Definitions
// Types used across multiple tabs (Launch, Sessions, Templates)

// ===== SPAWN MODES =====

type SpawnMode int

const (
	SpawnXtermWindow SpawnMode = iota
	SpawnTmuxWindow
	SpawnTmuxSplitH
	SpawnTmuxSplitV
	SpawnTmuxLayout
	SpawnCurrentPane
)

func (s SpawnMode) String() string {
	switch s {
	case SpawnXtermWindow:
		return "XTerm Window"
	case SpawnTmuxWindow:
		return "Tmux Window"
	case SpawnTmuxSplitH:
		return "Tmux Split Horizontal"
	case SpawnTmuxSplitV:
		return "Tmux Split Vertical"
	case SpawnTmuxLayout:
		return "Tmux Layout"
	case SpawnCurrentPane:
		return "Current Pane"
	default:
		return "Unknown"
	}
}

// ===== TMUX LAYOUTS =====

type TmuxLayout int

const (
	LayoutMainVertical TmuxLayout = iota
	LayoutMainHorizontal
	LayoutTiled
	LayoutEvenHorizontal
	LayoutEvenVertical
)

func (l TmuxLayout) String() string {
	switch l {
	case LayoutMainVertical:
		return "main-vertical"
	case LayoutMainHorizontal:
		return "main-horizontal"
	case LayoutTiled:
		return "tiled"
	case LayoutEvenHorizontal:
		return "even-horizontal"
	case LayoutEvenVertical:
		return "even-vertical"
	default:
		return "main-vertical"
	}
}

// ===== LAUNCH ITEM (for spawning) =====

type LaunchItem struct {
	Name    string
	Command string
	Cwd     string
}

// ===== TMUX SESSIONS (from tmuxplexer) =====

type TmuxSession struct {
	Name       string
	Windows    int
	Attached   bool
	Created    string
	LastActive string
	WorkingDir string
	GitBranch  string
	AITool     string // "claude", "codex", "gemini", or ""
}

type TmuxWindow struct {
	Index  int
	Name   string
	Panes  int
	Active bool
}

type TmuxPane struct {
	ID      string
	Index   int
	Width   int
	Height  int
	Active  bool
	Command string
}

// ===== SESSION TEMPLATES (from tmuxplexer) =====

type SessionTemplate struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    string            `json:"category,omitempty"`
	WorkingDir  string            `json:"working_dir"`
	Layout      string            `json:"layout"` // "2x2", "4x2", "3x3", etc.
	Panes       []PaneTemplate    `json:"panes"`
	Env         map[string]string `json:"env,omitempty"`
}

type PaneTemplate struct {
	Command    string `json:"command"`
	Title      string `json:"title,omitempty"`
	WorkingDir string `json:"working_dir,omitempty"`
}

// ===== MESSAGES =====

type SpawnCompleteMsg struct {
	Err error
}

type SessionKilledMsg struct {
	SessionName string
	Err         error
}

type SessionRenamedMsg struct {
	OldName string
	NewName string
	Err     error
}
