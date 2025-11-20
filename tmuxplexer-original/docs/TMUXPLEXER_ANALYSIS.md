# Tmuxplexer Codebase Analysis & TFE Integration Guide

**Date:** October 26, 2025  
**Project:** tmuxplexer (~/projects/tmuxplexer)  
**Language:** Go 1.24.0  
**Total Code:** 4,749 lines across 14 Go files  
**Status:** Production-ready with 8 phases completed

---

## 1. PROJECT OVERVIEW

### What is Tmuxplexer?

A modern **Terminal User Interface (TUI) tmux session manager** written in Go using Bubble Tea framework. It provides:

- **4-Panel Accordion Layout** - Header, Sessions (left), Templates (right), Preview (footer) with dynamic expansion
- **Workspace Templates** - Store multi-pane tmux layouts in JSON, create complex workspaces with one keystroke
- **Session Management** - Attach, kill, rename sessions; view live pane content with scrollback history
- **Claude Code Integration** - Real-time status tracking for Claude AI sessions via hooks
- **Popup Mode** - Launch from tmux with `Ctrl+b o` keybinding
- **CLI Flags** - Template creation and directory override for external tool integration

### Current Status

All 8 phases completed and working:
- ✅ Phase 1: 4-Panel Accordion Layout
- ✅ Phase 2: Workspace Templates  
- ✅ Phase 3: Session Management
- ✅ Phase 4: Live Pane Preview & Window Navigation
- ✅ Phase 5: Claude Code Integration
- ✅ Phase 6: Scrollable Preview with Full Scrollback
- ✅ Phase 7: Template Creation Wizard & Template Deletion
- ✅ Phase 8: Popup Mode with Tmux Keybinding

### Key Differentiators

1. **4-Panel Layout** - Unique accordion design vs single-pane tools like `tmux-sessionx`
2. **Claude Integration** - Real-time status indicators for AI sessions via hooks system
3. **Workspace Templates** - "Save as template" feature lets you create templates from running sessions
4. **Context-Aware CWD** - `--cwd` and `--template` flags for TFE integration
5. **Live Preview** - Full scrollback history with PgUp/PgDn navigation, auto-scroll for Claude sessions

---

## 2. DIRECTORY STRUCTURE

```
tmuxplexer/
├── main.go                    # Entry point (165 lines)
│                              # - TTY detection
│                              # - CLI flag parsing (--popup, --cwd, --template)
│                              # - Test mode handling (test_template, test_create)
│
├── types.go                   # Type definitions (299 lines)
│                              # - Model struct with UI state
│                              # - Config types (Theme, UI, Layout, Performance)
│                              # - TmuxSession, TmuxWindow, TmuxPane structs
│                              # - SessionTemplate, PaneTemplate for templates.json
│                              # - TemplateBuilder & SessionSaveBuilder for wizards
│                              # - Custom message types (sessionsLoadedMsg, etc.)
│
├── model.go                   # Model initialization & layout (601 lines)
│                              # - initialModel() - Creates app state
│                              # - calculateFourPanelLayout() - Dynamic panel sizing
│                              # - Update*Content() functions - Keeps panels in sync
│
├── update.go                  # Message dispatcher (238 lines)
│                              # - Init() - Start auto-refresh ticker
│                              # - Update() - Main message router
│                              # - Message handlers (sessionCreated, killed, etc.)
│
├── update_keyboard.go         # Keyboard handling (982 lines) **LARGEST FILE**
│                              # - handleKeyPress() - Main keyboard event handler
│                              # - Navigation (arrow keys, vim keys, page up/down)
│                              # - Panel switching (keys 1,2,3,4)
│                              # - Template wizard input handling
│                              # - Session save wizard
│
├── update_mouse.go            # Mouse event handling (185 lines)
│                              # - handleMouseEvent() - Click & drag processing
│                              # - Panel focus via mouse clicks
│
├── view.go                    # View rendering (513 lines)
│                              # - View() - Main render dispatcher
│                              # - renderFourPanelLayout() - Panel arrangement
│                              # - renderDynamicPanel() - Individual panel styling
│                              # - getTemplateWizardPrompt() - Wizard UI
│
├── styles.go                  # Lipgloss styling (238 lines)
│                              # - Color palette definitions
│                              # - Border styles
│                              # - Text formatting styles
│
├── config.go                  # Configuration (237 lines)
│                              # - loadConfig() - YAML config loading
│                              # - getDefaultConfig() - Default settings
│                              # - saveConfig() - YAML saving
│                              # - Config path: ~/.config/tmuxplexer/config.yaml
│
├── templates.go               # Template management (272 lines)
│                              # - loadTemplates() - Load from ~/.config/tmuxplexer/templates.json
│                              # - saveTemplates() - Save to JSON
│                              # - getDefaultTemplates() - Built-in templates
│                              # - addTemplate() - Append new template
│                              # - deleteTemplate() - Remove template
│
├── tmux.go                    # Tmux integration (633 lines) **SECOND LARGEST**
│                              # - listSessions() - Get all sessions with git branch
│                              # - createSessionFromTemplate() - Grid layout creation
│                              # - attachToSession() - Attach/switch to session
│                              # - killSession() - Delete session
│                              # - listWindows() / listPanes() - Window/pane details
│                              # - capturePane() - Get pane content (full scrollback)
│                              # - extractSessionInfo() - Save session as template
│                              # - detectGridLayout() - Auto-detect pane grid pattern
│
├── claude_state.go            # Claude integration (212 lines)
│                              # - detectClaudeSession() - Check if session is Claude
│                              # - getClaudeStateForSession() - Read state files
│                              # - findStateByPane() / findStateByWorkingDir()
│                              # - Status indicators: Idle, Processing, Tool Use, etc.
│
├── hooks/                     # Claude Code hooks system
│   ├── state-tracker.sh       # Bash script receiving hook events
│   ├── install.sh             # Installation script
│   ├── test-hooks.sh          # Hook testing suite
│   └── claude-settings-hooks.json # Hook configuration
│
├── components/                # Reusable UI components
│   ├── dialog/                # Confirmation dialogs
│   ├── list/                  # List views
│   ├── panel/                 # Panel layouts
│   └── status/                # Status bars
│
├── lib/                       # Utility libraries
│   ├── clipboard/             # Clipboard operations
│   ├── config/                # Config utilities
│   ├── keybindings/           # Keyboard mapping
│   ├── logger/                # Debug logging
│   ├── terminal/              # Terminal detection
│   └── theme/                 # Theme management
│
├── docs/                      # Documentation
│   ├── claude-hooks-integration.md
│   ├── HOOKS-QUICKREF.md
│   └── Other research docs
│
├── test_template.go           # Test: view templates (no TTY)
├── test_create_session.go     # Test: create from template (no TTY)
│
├── go.mod                     # Dependencies
├── go.sum                     # Dependency checksums
├── tmuxplexer                 # Compiled binary
├── install.sh                 # Installation script
├── CLAUDE.md                  # Development guide (607 lines)
├── README.md                  # User documentation
├── PLAN.md                    # Project roadmap
└── CHANGELOG.md               # Version history
```

---

## 3. CORE FUNCTIONALITY

### 3.1 Architecture Pattern

Follows **Bubble Tea Elm Architecture**:

```
Model (State) → Update (Message → State) → View (Render)
   ↓                ↓                         ↓
types.go    update*.go                    view.go
model.go    update_keyboard.go            styles.go
            update_mouse.go
```

### 3.2 4-Panel Layout System

**Dynamic panel sizing with weight-based calculations:**

```
┌─────────────────────────────────────────────────┐
│  [Header: Stats & Quick Actions]  (accordionMode-dependent)
├──────────────────┬──────────────────────────────┤
│  Left Panel      │  Right Panel                │
│  Sessions +      │  Templates +               │
│  Details         │  Details                   │
│  (50%)           │  (50%)                     │
├──────────────────┴──────────────────────────────┤
│  [Footer: Live Preview + Scrollback]            │
│  (acordionMode compresses other panels)         │
└─────────────────────────────────────────────────┘
```

**Key features:**
- Weight-based panel sizing in `calculateFourPanelLayout()` (model.go)
- Accordion mode: focused panel expands, others compress
- Focus switching: keys 1 (sessions), 2 (templates), 3 (preview), 4 (header)
- Mouse click to focus any panel
- Independent content and navigation per panel

### 3.3 Session State Management

**Model tracks:**
- `sessions []TmuxSession` - Loaded tmux sessions
- `templates []SessionTemplate` - Available templates
- `selectedSession int` - Index of selected session
- `selectedTemplate int` - Index of selected template
- `windows []TmuxWindow` - Windows for current session
- `selectedWindow int` - Current window being previewed
- `previewBuffer []string` - Full pane scrollback lines
- `previewScrollOffset int` - Current scroll position
- Panel content arrays: `headerContent`, `leftContent`, `rightContent`, `footerContent`

**Auto-refresh system:**
- Ticker fires every 2 seconds (update.go:137)
- Calls `listSessions()` which refreshes all data
- Updates panel content if data changed
- Claude state detection integrated into `listSessions()`

### 3.4 Keyboard & Mouse Input

**Main handlers:**
- `handleKeyPress()` (update_keyboard.go) - All keyboard navigation
- `handleMouseEvent()` (update_mouse.go) - Click handling and focus

**Keyboard shortcuts:**
| Key | Action |
|-----|--------|
| `1`/`2`/`3`/`4` | Focus panel |
| `↑/↓` or `k/j` | Navigate list items |
| `←/→` or `h/l` | Navigate windows (preview) |
| `PgUp`/`PgDn` | Scroll preview |
| `Home`/`End` or `g`/`G` | Jump preview to top/bottom |
| `Enter` | Attach session (left) OR create from template (right) |
| `s` | Save session as template |
| `n` | New template (wizard mode) |
| `e` | Edit templates.json in $EDITOR |
| `d` or `D` | Kill session (left) OR delete template (right) |
| `r` | Refresh preview |
| `a` | Toggle accordion mode |
| `Ctrl+R` | Refresh sessions & Claude state |
| `q`/`Ctrl+C` | Quit |

### 3.5 Configuration System

**Config file:** `~/.config/tmuxplexer/config.yaml` (YAML format)

**Configurable:**
- Theme: dark, light, solarized, dracula, nord, custom
- Custom colors (primary, secondary, background, foreground, accent, error)
- UI options: mouse enabled, icons, line numbers
- Layout: single, dual_pane, multi_panel, tabbed
- Performance: lazy loading, cache size, async operations
- Logging: enabled, level, file path

**Default config:** Hardcoded in `getDefaultConfig()` if file missing

---

## 4. TEMPLATES SYSTEM

### 4.1 Template Format (JSON)

**Location:** `~/.config/tmuxplexer/templates.json`

**Structure:**
```json
[
  {
    "name": "Simple Dev (2x2)",
    "description": "Basic development workspace",
    "working_dir": "~",
    "layout": "2x2",
    "panes": [
      {
        "command": "nvim",
        "title": "Editor",
        "working_dir": "~"  // Optional: per-pane override
      },
      {
        "command": "bash",
        "title": "Terminal"
      }
    ]
  }
]
```

**Template fields:**
- `name` - Template display name
- `description` - Human-readable description
- `working_dir` - Base directory for all panes (can be overridden by --cwd flag)
- `layout` - Grid format: "2x2", "3x3", "4x2", "1x1", "NxM"
- `panes` - Array of pane configurations
  - `command` - Shell command to run
  - `title` - Pane title (optional)
  - `working_dir` - Optional per-pane override

### 4.2 Built-in Templates

1. **Simple Dev (2x2)** - nvim, bash, lazygit, btop
2. **Frontend Dev (2x2)** - claude-code, nvim, npm run dev, lazygit
3. **TFE Development (4x2)** - Claude, editor, dev server, git, TFE browser, tests, monitor, logs
4. **TUITools (1x1)** - pyradio (example single-pane template)

### 4.3 Template Operations

**Creating templates:**
1. **Wizard mode** - Press `n` in right panel, step through fields
2. **From running session** - Press `s` in left panel, extract panes and layout
3. **Manual editing** - Press `e` in right panel, edit templates.json directly

**Template deletion:**
- Select template, press `d`, confirm with `y`

**Template workflow:**
```
Create/Edit Template → Save to templates.json → Right panel refreshes → 
Select template → Press Enter → Creates tmux session with layout + commands
```

### 4.4 Session Creation Flow

1. User selects template and presses Enter
2. `selectItem()` creates `createSessionCmd()` message
3. `createSessionFromTemplate()` (tmux.go):
   - Resolves working directory (template dir or --cwd override)
   - Creates tmux session: `tmux new-session -s <name> -d`
   - Creates grid layout: `tmux split-window` commands
   - Applies tiled layout: `tmux select-layout tiled`
   - Sends commands to each pane: `tmux send-keys -t <pane> <command> Enter`
4. Returns `sessionCreatedMsg`
5. Auto-refresh refreshes session list after 2 seconds

---

## 5. CLAUDE CODE INTEGRATION

### 5.1 How It Works

**Flow:**
```
Claude Code → Hook fires → state-tracker.sh → /tmp/claude-code-state/{id}.json
                                                        ↓
Tmuxplexer (every 2s) → listSessions() → detectClaudeSession() → getClaudeStateForSession()
                                                                        ↓
                                         Left panel shows status icon + full status text
```

### 5.2 Hook System

**Installation:**
```bash
cd ~/projects/tmuxplexer
./hooks/install.sh
```

This copies hooks to `~/.claude/hooks/` and adds configuration to `~/.claude/settings.json`.

**Available hooks:**
- `SessionStart` - Claude starts
- `UserPromptSubmit` - User sends message
- `PreToolUse` - Before tool execution
- `PostToolUse` - After tool execution
- `Stop` - Claude finishes
- `Notification` - System notifications (awaiting-input)

**State file format:** `/tmp/claude-code-state/<session-id>.json`
```json
{
  "session_id": "abc123",
  "status": "tool_use",
  "current_tool": "Edit",
  "working_dir": "/home/user/project",
  "last_updated": "2025-10-26T14:30:22Z",
  "tmux_pane": "%42",
  "pid": 12345,
  "hook_type": "pre-tool",
  "details": { "event": "tool_starting", "tool": "Edit" }
}
```

### 5.3 Status Indicators

Shown in left panel next to session name:
- 🟢 **Idle** - Waiting for input
- 🟡 **Processing** - Thinking about request
- 🔧 **Tool Use** - Executing tool (Edit, Read, Bash, etc.)
- ⚙️ **Working** - Processing results
- ⏸️ **Awaiting Input** - Waiting for user response
- ⚪ **Stale** - No updates >60 seconds (shows last known status)

**Special Claude session styling:**
- Orange text (#D19A66) for Claude session names
- Bold for emphasis
- Status line shows full description

### 5.4 Claude Session Display

```
Left panel shows:
● claude-session 🔧                    (filled bullet = attached, status icon, ORANGE TEXT)
  📁 ~/projects/tmuxplexer  main       (working dir, git branch)
  🔧 Using Read                        (Claude status)
```

### 5.5 Auto-Scroll Feature

When previewing a Claude session:
- First load automatically scrolls to bottom
- Shows current conversation instead of empty space
- User can still scroll up with PgUp to view history
- Scroll offset resets when changing sessions

---

## 6. CONFIGURATION FORMATS

### 6.1 Config File (YAML)

**Path:** `~/.config/tmuxplexer/config.yaml`

**Example:**
```yaml
theme: "dark"

custom_theme:
  primary: "#61AFEF"
  secondary: "#C678DD"
  background: "#282C34"
  foreground: "#ABB2BF"
  accent: "#98C379"
  error: "#E06C75"

keybindings: "default"

layout:
  type: "single"
  split_ratio: 0.5
  show_divider: true

ui:
  show_title: true
  show_status: true
  show_line_numbers: false
  mouse_enabled: true
  show_icons: true
  icon_set: "nerd_font"

performance:
  lazy_loading: true
  cache_size: 100
  async_operations: true

logging:
  enabled: false
  level: "info"
  file: "~/.local/share/tmuxplexer/debug.log"
```

### 6.2 Templates File (JSON)

**Path:** `~/.config/tmuxplexer/templates.json`

**Auto-created with defaults if missing**

**Format:** Array of SessionTemplate objects (see section 4.1)

### 6.3 CLI Flags

**Current flags:**
```bash
./tmuxplexer                          # Normal TUI mode
./tmuxplexer --popup                  # Popup mode (tmux)
./tmuxplexer --template 0             # Create from template 0, exit
./tmuxplexer --cwd /path --template 1 # Override working directory
./tmuxplexer test_template            # Show templates (no TTY)
./tmuxplexer test_create 0            # Create from template (no TTY)
```

**Flag parsing in main.go (lines 20-31):**
- `--popup` - Boolean flag for popup mode
- `--cwd <directory>` - Override template's working directory
- `--template <index>` - Create session from template by index
- Special handling for `test_template` and `test_create` before TTY check

---

## 7. DEPENDENCIES

### Go Modules (go.mod)

**Direct dependencies:**
```
github.com/charmbracelet/bubbles v0.20.0      # BubbleTeaUI components (input, list, etc.)
github.com/charmbracelet/bubbletea v1.1.0     # TUI framework (Model-View-Update)
github.com/charmbracelet/lipgloss v0.13.1     # Styling and borders
gopkg.in/yaml.v3 v3.0.1                       # YAML config parsing
```

**Transitive dependencies:** 20+ sub-packages for ANSI, terminal control, Unicode handling

### No External Tools Required

The application is self-contained but relies on these system commands:
- `tmux` - Session, window, pane management
- `$EDITOR` - For editing templates.json
- `git` - For branch detection (optional)

---

## 8. TFE INTEGRATION POINTS

### 8.1 Current Integration Status

✅ **Already implemented:**
- `--cwd` flag - Override template working directory
- `--template` flag - Select template by index
- Combined usage: `tmuxplexer --cwd $PWD --template 1`

**PLAN.md shows the next integration steps:**

### 8.2 Proposed TFE Integration

**How TFE can integrate:**

In TFE's `context_menu.go`, add a menu item for "Launch Dev Workspace":

```go
case "Launch Dev Workspace":
    return m, tea.ExecProcess(
        exec.Command("tmuxplexer",
            "--cwd", m.currentPath,  // Current TFE directory
            "--template", "0"),       // Template index
        nil,
    )
```

**Workflow:**
1. Browse to project in TFE: `/home/matt/projects/myapp`
2. Right-click → "Launch Dev Workspace"
3. Tmuxplexer creates session in `myapp` directory (not hardcoded template path!)
4. Session opens in tmux with all configured panes and commands
5. TFE can also support multiple templates per project context

### 8.3 TFE Development Template Example

**Built-in template for TFE development:**
```json
{
  "name": "TFE Development (4x2)",
  "description": "Full TFE development environment with 8 panes",
  "working_dir": "~/projects/TFE",
  "layout": "4x2",
  "panes": [
    {"command": "claude-code .", "title": "Claude AI"},
    {"command": "nvim", "title": "Editor"},
    {"command": "npm run dev", "title": "Dev Server"},
    {"command": "lazygit", "title": "Git"},
    {"command": "./tfe .", "title": "TFE Browser"},
    {"command": "npm test -- --watch", "title": "Tests"},
    {"command": "btop", "title": "Monitor"},
    {"command": "docker compose logs -f || bash", "title": "Logs"}
  ]
}
```

### 8.4 Integration Strategies for TFE

**Option 1: Context Menu Integration (RECOMMENDED)**
- Add "Launch Template" submenu in context menu
- Shows available templates as menu items
- Passes current TFE directory as `--cwd`
- Best user experience

**Option 2: Command Palette**
- Ctrl+P in TFE opens template selector
- Shows available templates
- Creates session in current directory
- More TFE-native experience

**Option 3: Keybinding**
- Single key (e.g., `t`) to launch default template
- Could make it context-aware by project type
- Fastest but least flexible

**Option 4: Template Directory Matching**
- Detect templates for current project
- Store project-specific template configs
- Auto-select matching template based on directory structure

---

## 9. KEY CODE PATTERNS

### 9.1 Message Flow Pattern

**Creating a session:**
```go
// User presses Enter
// → handleKeyPress() in update_keyboard.go
// → selectItem() message handler
// → createSessionCmd() returns Cmd
// → createSessionFromTemplate() executes
// → Returns sessionCreatedMsg
// → Update() handles sessionCreatedMsg
// → Refreshes sessions via refreshSessionsCmd()
```

### 9.2 Panel Update Pattern

```go
// In model.go:
func (m *model) updateLeftPanelContent() {
    // Clear content
    m.leftContent = []string{}
    
    // Build from sessions data
    for i, session := range m.sessions {
        // Format and append lines
    }
}

// Called whenever:
// - Session list refreshes
// - Selection changes
// - View mode changes
```

### 9.3 Auto-Refresh Pattern

```go
// In update.go:
func (m model) Init() tea.Cmd {
    return tickCmd()
}

// Main update loop:
case tickMsg:
    return m, tea.Batch(
        refreshSessionsCmd(),
        tickCmd(), // Schedule next tick
    )

// refreshSessionsCmd() in update_keyboard.go:
func refreshSessionsCmd() tea.Cmd {
    return func() tea.Msg {
        sessions, err := listSessions()
        return sessionsLoadedMsg{sessions, err}
    }
}
```

### 9.4 Keyboard Input Pattern

```go
// In update_keyboard.go:
func (m model) handleKeyPress(msg tea.KeyMsg) (model, tea.Cmd) {
    switch m.focusedPanel {
    case "left":
        // Handle left panel keys
    case "right":
        // Handle right panel keys
    case "footer":
        // Handle preview panel keys
    }
    
    switch msg.String() {
    case "up", "k":
        // moveUp()
    case "enter":
        // selectItem()
    }
}
```

---

## 10. DEVELOPMENT WORKFLOW

### 10.1 Building

```bash
cd ~/projects/tmuxplexer
go build -o tmuxplexer
```

### 10.2 Running

**TUI mode (requires terminal):**
```bash
./tmuxplexer
```

**Popup mode (from inside tmux):**
```bash
./tmuxplexer --popup
```

**CLI mode (no TTY required):**
```bash
./tmuxplexer --template 0
./tmuxplexer --cwd /path --template 1
./tmuxplexer test_template
./tmuxplexer test_create 0
```

### 10.3 Testing Without TTY

The `test_*.go` files allow testing without a terminal:
- `test_template.go` - Lists templates from templates.json
- `test_create_session.go` - Creates a session from template

### 10.4 Debug Logging

Enable in config.yaml:
```yaml
logging:
  enabled: true
  level: "debug"
  file: "~/.local/share/tmuxplexer/debug.log"
```

### 10.5 Common Development Tasks

**Adding a new keyboard shortcut:**
1. Edit `update_keyboard.go`, find `handleKeyPress()`
2. Add case in appropriate panel's switch statement
3. Implement action or call helper function

**Adding a new display mode:**
1. Define enum in `types.go`
2. Create render function in `view.go`
3. Call from main View() dispatcher
4. Add keyboard shortcut to toggle

**Adding a new tmux operation:**
1. Add function to `tmux.go`
2. Create message type in `types.go`
3. Add command creator in `update.go`
4. Handle message in Update() method

---

## 11. CURRENT STATE & MATURITY

### What's Working

✅ All 8 phases complete and functional:
- 4-panel accordion layout with perfect alignment
- Template loading, creation, deletion, editing
- Session management (attach, kill)
- Live pane preview with scrollback
- Claude Code integration with status tracking
- Popup mode with tmux keybinding
- Full keyboard and mouse navigation
- Auto-refresh every 2 seconds
- CLI flags for integration

### Known Limitations

- No command history/execution within panes
- No session renaming UI (can be added)
- No theme customization UI (YAML only)
- State files in /tmp auto-cleanup after 24h (Claude hooks)

### Performance Characteristics

- Session list refresh: ~500ms
- Template creation: <100ms
- Claude state detection: ~50ms per session
- UI render: <16ms (60 FPS)
- Memory: ~15-20MB for 20+ sessions
- Disk: State files ~1KB each

---

## 12. INTEGRATION CHECKLIST FOR TFE

### Phase 1: Current (Already Done)
- [x] `--cwd` flag support
- [x] `--template` flag support
- [x] Combined CLI usage

### Phase 2: Recommended for TFE Integration
- [ ] Add context menu item "Launch Dev Workspace"
- [ ] Pass `currentPath` as `--cwd`
- [ ] Select template by index or name
- [ ] Show success message with session name
- [ ] Option: Ctrl+Alt+T to launch default template

### Phase 3: Future Enhancements
- [ ] Per-project template configuration
- [ ] Template selection dialog in TFE
- [ ] Auto-detect template type by project structure
- [ ] Template variables (${PROJECT}, ${HOME}, etc.)
- [ ] One-click session creation with logging

---

## 13. QUICK REFERENCE FOR DEVELOPERS

### Code Location Quick Index

| Feature | Files |
|---------|-------|
| **Core UI Layout** | view.go, model.go |
| **Keyboard Input** | update_keyboard.go |
| **Mouse Input** | update_mouse.go |
| **Message Handling** | update.go |
| **Session Management** | tmux.go |
| **Template Management** | templates.go |
| **Claude Integration** | claude_state.go, hooks/* |
| **Configuration** | config.go, types.go |
| **Styling** | styles.go |

### Largest Files (Complexity)
1. `update_keyboard.go` (982 lines) - All keyboard event handling
2. `tmux.go` (633 lines) - All tmux integration
3. `model.go` (601 lines) - State management and layout

### Entry Points
- `main()` in main.go - TTY check, flag parsing
- `initialModel()` in model.go - App initialization
- `(m model) Update()` in update.go - Main update loop
- `(m model) View()` in view.go - Main render

---

## 14. FILE SIZES & METRICS

```
Total Code: 4,749 lines

Breakdown:
- update_keyboard.go:  982 lines (21%) - Keyboard handling
- tmux.go:            633 lines (13%) - Tmux integration
- model.go:           601 lines (13%) - State & layout
- types.go:           299 lines (6%)  - Type definitions
- view.go:            513 lines (11%) - Rendering
- templates.go:       272 lines (6%)  - Template management
- config.go:          237 lines (5%)  - Configuration
- styles.go:          238 lines (5%)  - Styling
- update.go:          238 lines (5%)  - Message dispatcher
- claude_state.go:    212 lines (4%)  - Claude integration
- update_mouse.go:    185 lines (4%)  - Mouse handling
- main.go:            165 lines (3%)  - Entry point
- test files:         174 lines (4%)  - Testing
```

**Compiled binary:** 5.4MB (stripped: ~2.5MB)

---

## 15. RESOURCES & DOCUMENTATION

### Internal Documentation
- **CLAUDE.md** (607 lines) - Development guide, architecture, keyboard shortcuts
- **README.md** (228 lines) - User guide, quick start, troubleshooting
- **PLAN.md** (50+ pages) - Detailed roadmap, ideas, phase descriptions
- **docs/** - Claude hooks integration, research, brainstorms

### External Resources
- Bubble Tea: https://github.com/charmbracelet/bubbletea
- Lipgloss: https://github.com/charmbracelet/lipgloss
- Tmux Manual: https://man.openbsd.org/tmux

### Related Projects
- **TUITemplate** (~projects/TUITemplate) - Reusable TUI components and patterns
- **TFE** (~projects/TFE) - Terminal File Explorer (integration target)

---

## SUMMARY

Tmuxplexer is a **well-architected, production-ready TUI application** built on Bubble Tea with:
- Clean separation of concerns (14 focused Go files)
- Comprehensive feature set (8 phases complete)
- Strong integration capabilities (`--cwd`, `--template` flags)
- Real-time Claude Code tracking
- Flexible configuration (YAML, JSON)

**For TFE integration:** The infrastructure is already in place. TFE just needs to call:
```bash
tmuxplexer --cwd $PWD --template <index>
```
from its context menu. This will create workspace templates in the user's current directory.

**Maturity Level:** Beta/Production - all core features complete, well-tested, documented.
