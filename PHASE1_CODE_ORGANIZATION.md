# Phase 1: Code Organization Plan

## Overview

This document outlines the code organization for Phase 1 of the TUI Launcher + Tmuxplexer integration. The goal is to create a tab-based architecture that combines both projects into a unified workspace manager.

## Directory Structure

```
tui-launcher/
├── main.go                      # Entry point, tab routing
├── types.go                     # Unified type definitions
├── model.go                     # Unified model initialization
├── update.go                    # Main update dispatcher
├── view.go                      # Main view compositor
├── shared/                      # Shared utilities
│   ├── tmux.go                 # Combined tmux operations
│   ├── config.go               # Configuration loading
│   ├── styles.go               # Shared lipgloss styles
│   └── utils.go                # Common helper functions
├── tabs/
│   ├── launch/                 # Launch tab (current tui-launcher)
│   │   ├── model.go           # Launch tab model
│   │   ├── view.go            # Launch tab rendering
│   │   ├── update.go          # Launch tab update logic
│   │   └── tree.go            # Tree building and navigation
│   ├── sessions/               # Sessions tab (from tmuxplexer)
│   │   ├── model.go           # Sessions tab model
│   │   ├── view.go            # Sessions tab rendering
│   │   ├── update.go          # Sessions tab update logic
│   │   └── preview.go         # Live preview logic
│   └── templates/              # Templates tab (from tmuxplexer)
│       ├── model.go           # Templates tab model
│       ├── view.go            # Templates tab rendering
│       ├── update.go          # Templates tab update logic
│       └── wizard.go          # Template creation wizard
└── tmuxplexer-original/        # Reference (unchanged)
```

## Code Migration Strategy

### 1. Shared Code (`shared/` directory)

#### `shared/tmux.go` - Combined Tmux Operations
Merge tmux functionality from both projects:

**From tui-launcher (`spawn.go`):**
- `insideTmux()` - Detect if inside tmux
- `spawnSingle()` - Single command spawning
- `spawnMultiple()` - Batch command spawning with layouts
- `tmuxSplitHorizontal/Vertical()` - Pane splitting
- `tmuxNewWindow()` - Window creation
- `tmuxSendKeys()` - Send keys to panes
- `generateSessionName()` - Unique session naming

**From tmuxplexer (`tmuxplexer-original/tmux.go`):**
- `listSessions()` - Session enumeration with Claude detection
- `listWindows()` - Window enumeration
- `listPanes()` - Pane enumeration
- `capturePane()` - Pane content capture for previews
- `attachToSession()` - Session attachment
- `killSession()` - Session termination
- `renameSession()` - Session renaming
- `createSessionFromTemplate()` - Template-based session creation
- `createGridLayout()` - Grid layout creation (2x2, 3x3, etc.)
- `getGitBranch()` - Git branch detection

**Integration Notes:**
- Keep ALL functions from both files
- Resolve naming conflicts (both have similar spawn logic)
- Maintain tmuxplexer's Claude detection (detectClaudeSession, getClaudeStateForSession)
- Use tui-launcher's spawn patterns for Launch tab
- Use tmuxplexer's session management for Sessions/Templates tabs

#### `shared/config.go` - Configuration Loading
**Phase 1:** Keep separate configs
- `LoadLaunchConfig()` - Loads `~/.config/tui-launcher/config.yaml`
- `LoadTemplates()` - Loads `~/.config/tmuxplexer/templates.json`

**Phase 3:** Will merge into single YAML config

#### `shared/styles.go` - Shared Lipgloss Styles
Extract common styles from both projects:
- Border styles (from both)
- Color palette (merge both)
- Text styles (bold, dim, etc.)
- Selected/focused styles

#### `shared/utils.go` - Common Utilities
- `expandPath()` - Path expansion (from tui-launcher)
- `visualWidth()` - ANSI-aware width calculation (from tmuxplexer)
- `padToVisualWidth()` - Visual width padding (from tmuxplexer)
- `truncateLine()` - Line truncation (from tmuxplexer)

### 2. Launch Tab (`tabs/launch/`)

#### Code Sources
Migrate from current tui-launcher root:
- `model.go:29-65` → `tabs/launch/model.go` - Launch model struct
- `tree.go` → `tabs/launch/tree.go` - Tree building logic (all)
- `model.go:549-748` → `tabs/launch/view.go` - View functions
- `model.go:76-547` → `tabs/launch/update.go` - Update logic

#### Launch Model (`tabs/launch/model.go`)
```go
package launch

type Model struct {
    // Display dimensions
    width  int
    height int

    // Multi-pane layout state (from current model)
    activePane        paneType
    globalItems       []launchItem
    projectItems      []launchItem
    globalTreeItems   []launchTreeItem
    projectTreeItems  []launchTreeItem
    globalCursor      int
    projectCursor     int
    globalExpanded    map[string]bool
    projectExpanded   map[string]bool

    // Info pane state
    currentInfo       paneInfo
    infoContent       string
    showingInfo       bool
    showingProjects   bool

    // Selection state
    selectedItems     map[string]bool

    // Spawn dialog state
    showSpawnDialog   bool
    selectedLayout    tmuxLayout
    layoutCursor      int

    // Config
    config            Config

    // Terminal detection
    terminalType      terminalType
    insideTmux        bool
    useTmux           bool
}

func New(cfg Config) Model
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

### 3. Sessions Tab (`tabs/sessions/`)

#### Code Sources
Migrate from `tmuxplexer-original/`:
- `model.go:18-125` → `tabs/sessions/model.go` - Sessions model initialization
- `model.go:233-661` → `tabs/sessions/view.go` - Sessions list rendering
- `model.go:877-1133` → `tabs/sessions/preview.go` - Preview panel rendering
- `update_keyboard.go:18-650` → `tabs/sessions/update.go` - Keyboard handling

#### Sessions Model (`tabs/sessions/model.go`)
```go
package sessions

type Model struct {
    // UI State
    width  int
    height int

    // Tmux data
    sessions           []TmuxSession
    selectedSession    int
    currentSessionName string

    // Tree view for sessions
    expandedSessions   map[string]bool
    sessionTreeItems   []SessionTreeItem
    sessionFilter      string // "all", "ai", "attached", "detached"

    // Window navigation
    windows            []TmuxWindow
    selectedWindow     int

    // Preview scrolling
    previewBuffer      []string
    previewScrollOffset int
    previewTotalLines   int

    // Focus state
    focusState         int // FocusSessions or FocusPreview

    // Mouse tracking
    mouseX, mouseY     int
    hoveredItem        string

    // Input mode (rename)
    inputMode          string
    inputBuffer        string
    renameTarget       string
}

func New() Model
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

### 4. Templates Tab (`tabs/templates/`)

#### Code Sources
Migrate from `tmuxplexer-original/`:
- `model.go:18-125` → `tabs/templates/model.go` - Templates model initialization
- `model.go:664-735` → `tabs/templates/view.go` - Templates tree rendering
- `model.go:1136-1300` → `tabs/templates/preview.go` - Template preview
- `update_keyboard.go:651-950` → `tabs/templates/wizard.go` - Template wizard
- `templates.go` → `tabs/templates/io.go` - Template loading/saving

#### Templates Model (`tabs/templates/model.go`)
```go
package templates

type Model struct {
    // UI State
    width  int
    height int

    // Template data
    templates          []SessionTemplate
    selectedTemplate   int

    // Tree view for categorized templates
    expandedCategories map[string]bool
    templateTreeItems  []TemplateTreeItem

    // Preview scrolling
    previewBuffer      []string
    previewScrollOffset int
    previewTotalLines   int

    // Focus state
    focusState         int // FocusTemplates or FocusPreview

    // Template creation wizard
    templateCreationMode bool
    templateBuilder      TemplateBuilder

    // Session save mode
    sessionSaveMode    bool
    sessionBuilder     SessionSaveBuilder

    // Input mode
    inputMode          string
    inputBuffer        string
}

func New() Model
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

## Unified Model Structure

The main model in root will contain tab-specific models:

```go
// types.go
type tabName string

const (
    tabLaunch    tabName = "launch"
    tabSessions  tabName = "sessions"
    tabTemplates tabName = "templates"
)

// model.go
type model struct {
    // Shared state
    width       int
    height      int
    currentTab  tabName
    err         error
    statusMsg   string

    // Tab-specific models
    launchModel    launch.Model
    sessionsModel  sessions.Model
    templatesModel templates.Model

    // Popup mode
    popupMode bool
}
```

## Tab Routing Logic

### Update Dispatcher (`update.go`)
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "1":
            m.currentTab = tabLaunch
            return m, nil
        case "2":
            m.currentTab = tabSessions
            return m, nil
        case "3":
            m.currentTab = tabTemplates
            return m, nil
        case "tab":
            // Cycle forward
            m.currentTab = nextTab(m.currentTab)
            return m, nil
        case "shift+tab":
            // Cycle backward
            m.currentTab = prevTab(m.currentTab)
            return m, nil
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    case tea.WindowSizeMsg:
        // Distribute size to all tabs
        m.width = msg.Width
        m.height = msg.Height
        // Update all tab models
    }

    // Route to active tab
    switch m.currentTab {
    case tabLaunch:
        var cmd tea.Cmd
        m.launchModel, cmd = m.launchModel.Update(msg)
        return m, cmd
    case tabSessions:
        var cmd tea.Cmd
        m.sessionsModel, cmd = m.sessionsModel.Update(msg)
        return m, cmd
    case tabTemplates:
        var cmd tea.Cmd
        m.templatesModel, cmd = m.templatesModel.Update(msg)
        return m, cmd
    }

    return m, nil
}
```

### View Compositor (`view.go`)
```go
func (m model) View() string {
    // Render tab bar
    var tabBar strings.Builder
    tabBar.WriteString(renderTabIndicator("1. Launch", m.currentTab == tabLaunch))
    tabBar.WriteString(renderTabIndicator("2. Sessions", m.currentTab == tabSessions))
    tabBar.WriteString(renderTabIndicator("3. Templates", m.currentTab == tabTemplates))
    tabBar.WriteString("\n\n")

    // Render active tab content
    var content string
    switch m.currentTab {
    case tabLaunch:
        content = m.launchModel.View()
    case tabSessions:
        content = m.sessionsModel.View()
    case tabTemplates:
        content = m.templatesModel.View()
    }

    return tabBar.String() + content
}
```

## Implementation Steps

### Step 1: Create Shared Code
1. Create `shared/tmux.go` - Merge both tmux files
2. Create `shared/config.go` - Separate config loaders
3. Create `shared/styles.go` - Extract common styles
4. Create `shared/utils.go` - Common utilities

### Step 2: Create Launch Tab
1. Extract launch model to `tabs/launch/model.go`
2. Extract tree logic to `tabs/launch/tree.go`
3. Extract view functions to `tabs/launch/view.go`
4. Extract update logic to `tabs/launch/update.go`

### Step 3: Create Sessions Tab
1. Extract sessions model to `tabs/sessions/model.go`
2. Extract sessions view to `tabs/sessions/view.go`
3. Extract preview logic to `tabs/sessions/preview.go`
4. Extract update logic to `tabs/sessions/update.go`

### Step 4: Create Templates Tab
1. Extract templates model to `tabs/templates/model.go`
2. Extract templates view to `tabs/templates/view.go`
3. Extract preview logic to `tabs/templates/preview.go`
4. Extract wizard logic to `tabs/templates/wizard.go`
5. Extract template I/O to `tabs/templates/io.go`

### Step 5: Create Main Entry Point
1. Update `types.go` with unified types
2. Update `model.go` with tab routing model
3. Create `update.go` with dispatcher
4. Create `view.go` with compositor

## Testing Strategy

### Checkpoint 1: Build Verification
```bash
go build
# Should compile without errors
```

### Checkpoint 2: Launch Tab
```bash
./tui-launcher
# Press 1 - should show current tui-launcher interface
# All existing features work (navigation, selection, spawning)
```

### Checkpoint 3: Sessions Tab
```bash
./tui-launcher
# Press 2 - should show tmuxplexer sessions list
# Can navigate, preview, attach to sessions
```

### Checkpoint 4: Templates Tab
```bash
./tui-launcher
# Press 3 - should show tmuxplexer templates
# Can navigate, preview, create from templates
```

### Checkpoint 5: Tab Switching
```bash
./tui-launcher
# Press 1/2/3 - switches tabs
# Press Tab - cycles forward
# Press Shift+Tab - cycles backward
```

## Key Design Decisions

### 1. Keep Separate Configs (Phase 1)
- Launch tab uses `~/.config/tui-launcher/config.yaml`
- Sessions/Templates use `~/.config/tmuxplexer/templates.json`
- Will merge in Phase 3

### 2. Shared Tmux Layer
- All tmux operations in `shared/tmux.go`
- Both spawn patterns coexist (launch vs templates)
- Claude detection remains from tmuxplexer

### 3. Tab Isolation
- Each tab is self-contained (model, view, update)
- Minimal coupling between tabs in Phase 1
- Cross-tab features come in Phase 4

### 4. Preserve All Features
- Launch tab: All current tui-launcher features
- Sessions tab: All tmuxplexer session features
- Templates tab: All tmuxplexer template features
- No feature loss during migration

## Next Phase Integration Points

### Phase 2: Enhanced Launch Tab
- Add post-launch feedback
- Auto-switch to Sessions tab after launching
- Session awareness in Launch tab

### Phase 3: Unified Config
- Merge configs into single YAML
- Migration script for existing configs
- Shared categories across tabs

### Phase 4: Cross-Tab Features
- Launch → Sessions (auto-switch)
- Sessions → Templates (save as template)
- Launch → Templates (convert selections to template)

### Phase 5: Polish
- Popup mode support
- Consistent hotkeys across tabs
- Visual polish and animations
