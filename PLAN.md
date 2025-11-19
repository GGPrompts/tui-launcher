# TUI Launcher - Development Plan

**A visual terminal launcher for managing projects, TUI tools, and batch command execution with tmux integration**

---

## Vision

Create a tree-based TUI launcher that allows:
- **Visual organization** of projects, tools, scripts, and AI commands
- **Multi-select spawning** (Space to select, Enter to launch)
- **Tmux integration** for batch launches with configurable layouts
- **Context-aware spawning** (inside tmux vs standalone)
- **Project-based working directories** for proper context
- **Saved profiles** for complex multi-pane setups

### Why Not Just mcfly?

**mcfly:** Great for "I ran this before, what was it?" (command history search)
**tui-launcher:** "I want to start X in Y context with Z layout" (workspace orchestrator)

They complement each other - mcfly for ad-hoc commands, launcher for organized workflows.

---

## Core Features

### 1. Tree View Navigation
- ✅ Hierarchical tree structure (borrowed from TFE)
- ✅ Expandable categories with `▶`/`▼` indicators
- ✅ Tree connectors: `├─`, `└─`, `│`
- ✅ Emoji icons for visual identification
- ✅ Multi-level nesting (Projects → Commands, Categories → Tools)
- ✅ Smooth scrolling in narrow terminals (Termux-optimized)

### 2. Multi-Select System
- **Space:** Toggle selection (☐ → ☑)
- **Enter:** Launch selected items (or single item at cursor)
- **Esc:** Clear selections
- **a:** Select all in current category
- **c:** Clear all selections
- Visual indicators: ☑ checkbox for selected items
- Status bar shows selection count

### 3. Spawn Modes

#### Single Item Launch
- Quick launch with default spawn mode
- Configurable per-item defaults

#### Multi-Item Launch
- Batch spawn dialog appears when 2+ items selected
- Choose tmux layout:
  - 📐 `main-vertical` - Main pane left, others stacked right
  - 📏 `main-horizontal` - Main pane top, others stacked below
  - 🔲 `tiled` - Grid layout
  - ⚡ `even-horizontal` - Equal width columns
  - ⚡ `even-vertical` - Equal height rows

#### Spawn Options
- 🪟 **New xterm window** - Separate window (non-tmux)
- 🔲 **Tmux new window** - New window in current session
- ⬛ **Tmux split horizontal** - Split current pane horizontally
- ⬜ **Tmux split vertical** - Split current pane vertically
- 📐 **Tmux layout** - Multi-pane with layout chooser
- 🎯 **Current pane** - Replace current pane (detach on exit)

### 4. Context Awareness

**Inside tmux:**
- Default to tmux splits/windows
- Offer session switching
- Detect current session/window

**Outside tmux:**
- Create new tmux session for multi-launch
- Fall back to xterm windows
- Auto-attach to created sessions

### 5. Configuration System

YAML-based config at `~/.config/tui-launcher/config.yaml`

```yaml
projects:
  - name: TFE
    icon: 🚀
    path: ~/projects/tfe
    commands:
      - name: TFE
        icon: 📂
        command: tfe
        spawn: tmux-split-h
      - name: Dev Server
        icon: 💻
        command: go run .
        spawn: tmux-split-v

    profiles:
      - name: Dev Environment
        icon: 🔧
        layout: main-vertical
        panes:
          - command: tfe
          - command: go run .
          - command: tail -f logs/debug.log

tools:
  - category: System Monitoring
    icon: 📊
    items:
      - name: htop
        icon: 💹
        command: htop
        spawn: tmux-split-v

ai:
  - name: Claude Code
    icon: 💬
    command: claude
    spawn: tmux-split-h
```

### 6. TFE Integration

Add to TFE's context menu (`~/.config/tfe/tools.yaml`):
```yaml
tools:
  - name: "Launch Environment"
    command: "tui-launcher --project {{file}}"
    icon: "🚀"
    showFor: "directories"
```

Right-click project folder in TFE → Launch Environment → Opens launcher focused on that project!

---

## Architecture

### File Structure
```
tui-launcher/
├── main.go               # Entry point (minimal)
├── types.go              # Type definitions
├── model.go              # Model initialization
├── update.go             # Update dispatcher
├── update_keyboard.go    # Keyboard handling
├── update_mouse.go       # Mouse handling
├── view.go               # View rendering
├── styles.go             # Lipgloss styles
├── config.go             # Config loading/parsing
├── tree.go               # Tree building/rendering
├── spawn.go              # Spawn logic (tmux/xterm)
├── dialog.go             # Spawn dialog UI
├── helpers.go            # Utility functions
├── go.mod
├── PLAN.md              # This file
└── README.md
```

### Key Types

```go
type launchItem struct {
    name         string
    path         string      // Unique identifier
    itemType     itemType    // category, command, profile
    icon         string
    command      string
    cwd          string      // Working directory
    defaultSpawn spawnMode
    children     []launchItem // For categories

    // For profiles
    isProfile    bool
    layout       tmuxLayout
    panes        []paneConfig
}

type itemType int
const (
    typeCategory itemType = iota  // Expandable folder
    typeCommand                    // Single executable
    typeProfile                    // Multi-launch config
)

type spawnMode int
const (
    spawnXtermWindow spawnMode = iota
    spawnTmuxWindow
    spawnTmuxSplitH
    spawnTmuxSplitV
    spawnTmuxLayout
    spawnCurrentPane
)

type tmuxLayout int
const (
    layoutMainVertical tmuxLayout = iota
    layoutMainHorizontal
    layoutTiled
    layoutEvenHorizontal
    layoutEvenVertical
)

type launchTreeItem struct {
    item        launchItem
    depth       int
    isLast      bool
    parentLasts []bool
}

type model struct {
    // Display
    width, height     int
    cursor            int

    // Tree state
    items             []launchItem
    treeItems         []launchTreeItem
    expandedItems     map[string]bool

    // Selection
    selectedItems     map[string]bool

    // Spawn dialog
    showSpawnDialog   bool
    selectedLayout    tmuxLayout

    // Config
    config            Config
}
```

---

## Implementation Phases

### Phase 1: Core Tree View ✅
**Goal:** Basic tree navigation with TFE patterns

- [x] Project structure
- [ ] Port `treeItem` structure from TFE
- [ ] Port `buildTreeItems()` logic
- [ ] Port `renderTreeView()` with emoji width handling
- [ ] Basic keyboard navigation (arrows, enter, esc)
- [ ] Expand/collapse categories
- [ ] Config loading (YAML)

**Files:** `types.go`, `tree.go`, `config.go`, `update_keyboard.go`, `view.go`

### Phase 2: Multi-Select System
**Goal:** Space to select, visual indicators

- [ ] Selection state tracking (`selectedItems` map)
- [ ] Space key to toggle selection
- [ ] Checkbox rendering (☐/☑)
- [ ] Visual feedback (highlight selected items)
- [ ] Selection count in status bar
- [ ] Clear selections (Esc)
- [ ] Select all in category (a)

**Files:** `update_keyboard.go`, `tree.go`, `view.go`

### Phase 3: Single Spawn Logic
**Goal:** Launch individual commands

- [ ] Detect if inside tmux (`$TMUX` env var)
- [ ] Spawn in tmux split horizontal
- [ ] Spawn in tmux split vertical
- [ ] Spawn in tmux new window
- [ ] Spawn in xterm window
- [ ] Set working directory per command
- [ ] Default spawn mode per item

**Files:** `spawn.go`, `helpers.go`

### Phase 4: Multi-Spawn Dialog
**Goal:** Batch launches with layout selection

- [ ] Spawn dialog UI
- [ ] Layout picker (arrow keys)
- [ ] Show selected items list
- [ ] Preview layout visually (ASCII art)
- [ ] Confirm/cancel
- [ ] Multi-spawn execution
- [ ] Apply tmux layout after spawning

**Files:** `dialog.go`, `spawn.go`, `view.go`

### Phase 5: Profile Support
**Goal:** Saved multi-pane configurations

- [ ] Profile type in config
- [ ] Profile rendering in tree
- [ ] Launch profile (create session + panes + layout)
- [ ] Session naming
- [ ] Auto-attach to created session

**Files:** `config.go`, `spawn.go`, `tree.go`

### Phase 6: Polish & Features
**Goal:** Production-ready launcher

- [ ] Mouse support (click to select/expand)
- [ ] Favorites system (star items)
- [ ] Recent launches (history)
- [ ] Search/filter (Ctrl+F)
- [ ] Command-line args (`--project`, `--tool`)
- [ ] TFE integration example
- [ ] Error handling (command not found, etc.)
- [ ] Session management (list, kill, switch)

**Files:** `update_mouse.go`, `favorites.go`, `search.go`, `main.go`

### Phase 7: Documentation
**Goal:** Usable by others

- [ ] README with screenshots
- [ ] Config examples
- [ ] Keyboard shortcuts reference
- [ ] Integration guide (TFE, tmux, etc.)
- [ ] Installation instructions
- [ ] Example configs for common workflows

**Files:** `README.md`, `HOTKEYS.md`, `examples/`

---

## Emoji System

**Critical:** NO variation selectors (U+FE0F)!
(Learned from TFE - causes width calculation bugs in go-runewidth)

### Icons Used
```go
// Spawn modes
emojiXtermWindow    = "🪟"  // U+1FA9F
emojiTmuxWindow     = "🔲"  // U+1F532
emojiTmuxSplitH     = "⬛"  // U+2B1B
emojiTmuxSplitV     = "⬜"  // U+2B1C
emojiTmuxLayout     = "📐"  // U+1F4D0
emojiCurrentPane    = "🎯"  // U+1F3AF
emojiBatchSpawn     = "🎛"  // U+1F39B (NO U+FE0F!)

// Item types
emojiProject        = "📦"  // U+1F4E6
emojiFolder         = "📁"  // U+1F4C1
emojiCommand        = "⚡"  // U+26A1
emojiProfile        = "🔧"  // U+1F527
emojiTUITool        = "🛠"  // U+1F6E0 (NO U+FE0F!)
emojiAI             = "🤖"  // U+1F916
emojiScript         = "📜"  // U+1F4DC
emojiFavorite       = "⭐"  // U+2B50

// Status
emojiExpanded       = "▼"   // U+25BC
emojiCollapsed      = "▶"   // U+25B6
emojiSelected       = "☑"   // U+2611
emojiUnselected     = "☐"   // U+2610
```

**Validation:**
```bash
echo -n "🛠" | xxd  # Should NOT see efb88f (U+FE0F)
```

---

## Key Interactions

### Navigation
- **↑/↓**: Move cursor
- **→**: Expand category
- **←**: Collapse category
- **Home**: Jump to top
- **End**: Jump to bottom
- **PgUp/PgDn**: Page up/down

### Selection
- **Space**: Toggle selection at cursor
- **a**: Select all in current category
- **c**: Clear all selections
- **Esc**: Clear selections (or close dialog)

### Launching
- **Enter**: Launch (single or show dialog for multi)
- **Ctrl+Enter**: Quick multi-launch (skip dialog, use default layout)
- **s**: Show spawn options for current item
- **p**: Launch as profile (if item is profile)

### Dialog Navigation
- **↑/↓**: Select layout option
- **Enter**: Confirm launch
- **Esc**: Cancel

---

## Technical Patterns from TFE

### 1. Emoji Width Handling
```go
// Port from TFE's render_file_list.go
func (m model) runeWidth(r rune) int {
    if r >= 0xFE00 && r <= 0xFE0F { // Variation selectors
        if m.terminalType == terminalWindowsTerminal {
            return 1
        }
        return 0
    }
    return runewidth.RuneWidth(r)
}
```

### 2. Tree Building (Recursive)
```go
// Port from TFE's buildTreeItems
func (m model) buildLaunchTree(items []launchItem, depth int, parentLasts []bool) []launchTreeItem {
    treeItems := []launchTreeItem{}

    for i, item := range items {
        isLast := i == len(items)-1

        // Add current item
        treeItems = append(treeItems, launchTreeItem{
            item:        item,
            depth:       depth,
            isLast:      isLast,
            parentLasts: append([]bool{}, parentLasts...),
        })

        // Recursively add children if expanded
        if item.itemType == typeCategory && m.expandedItems[item.path] {
            newParentLasts := append(parentLasts, isLast)
            children := m.buildLaunchTree(item.children, depth+1, newParentLasts)
            treeItems = append(treeItems, children...)
        }
    }

    return treeItems
}
```

### 3. Config System
```go
// Similar to TFE's config.go
type Config struct {
    Projects []ProjectConfig `yaml:"projects"`
    Tools    []CategoryConfig `yaml:"tools"`
    AI       []CommandConfig `yaml:"ai"`
    Scripts  []CategoryConfig `yaml:"scripts"`
}

func loadConfig() (Config, error) {
    configPath := filepath.Join(os.Getenv("HOME"), ".config/tui-launcher/config.yaml")
    data, err := os.ReadFile(configPath)
    // ... parse YAML
}
```

---

## Example Usage Flows

### Flow 1: Quick Single Launch
1. Open launcher: `tui-launcher`
2. Navigate to "TFE" with arrow keys
3. Press Enter → TFE launches in tmux split

### Flow 2: Multi-Select Batch Launch
1. Navigate to "MyApp" project
2. Space on "TFE" → ☑
3. Space on "npm run dev" → ☑
4. Space on "tail -f logs/app.log" → ☑
5. Press Enter → Spawn dialog appears
6. Select layout: `main-vertical`
7. Press Enter → Launches all 3 in tmux layout

Result:
```
┌──────────┬───────────────┐
│          │ npm run dev   │
│   TFE    ├───────────────┤
│          │ tail -f logs  │
└──────────┴───────────────┘
```

### Flow 3: Launch Profile
1. Navigate to "Dev Environment" (profile icon 🔧)
2. Press Enter → Entire dev environment spawns with saved layout
3. Auto-switches to new tmux session

### Flow 4: From TFE Integration
1. In TFE, navigate to project folder
2. Right-click → "Launch Environment"
3. Launcher opens pre-focused on that project
4. Quick select tools → batch launch

---

## Success Metrics

**Must Have:**
- ✅ Tree navigation smooth in Termux
- ✅ Multi-select works intuitively
- ✅ Batch tmux launches work correctly
- ✅ Config loads from YAML
- ✅ Working directory set properly per command

**Should Have:**
- Profiles save time vs manual setup
- Mouse support for convenience
- TFE integration is seamless
- Error messages are helpful

**Nice to Have:**
- Session management (list/kill/switch)
- Command history/favorites
- Fuzzy search for items
- Custom keybindings

---

## Current Status

**Phase:** 1 - Core Tree View (In Progress) + Phase 3 Spawn Logic (Completed!)

**Completed:**
- ✅ Project structure and types (`types.go`)
- ✅ Layout system with visual previews (`layouts.go`)
- ✅ Spawn logic using tmuxplexer pattern (`spawn.go`)
- ✅ Sample config with real projects (`~/.config/tui-launcher/config.yaml`)
- ✅ Comprehensive planning (`PLAN.md`, `LAYOUTS_DEMO.md`)

**Next Steps:**
1. Implement config loading (`config.go`)
2. Port tree building logic from TFE (`tree.go`)
3. Set up basic rendering (`view.go`, `styles.go`)
4. Add keyboard navigation (`update_keyboard.go`)
5. Implement model initialization (`model.go`)

**Blockers:** None

**Notes:**
- Leveraging TUITemplate architecture
- Reusing proven patterns from TFE (tree view) and tmuxplexer (spawn logic)
- Emoji width handling already solved in TFE
- Tmux spawn uses "create all panes, then apply layout" pattern (reliable!)

---

## Future Ideas (Post v1.0)

- **Remote spawning** - SSH into servers and spawn there
- **Docker integration** - Launch containers with tmux inside
- **Command templates** - Variables in commands ({{project}}, {{branch}})
- **Conditional commands** - Only show if file/dir exists
- **Status indicators** - Show if process is running (🟢/🔴)
- **Web UI** - Launch via browser (for remote access)
- **AI integration** - Ask Claude to generate launch configs
- **Export/import** - Share configs with team

---

**Last Updated:** 2025-01-19
**Status:** Planning Complete, Ready to Build! 🚀
