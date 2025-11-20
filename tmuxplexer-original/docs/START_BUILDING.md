# Tmuxplexer Development Kickstart Prompt

## Project Overview

**Tmuxplexer** is a modern, powerful TUI for managing tmux sessions with a unique **4-panel accordion layout**. This is the first TUI that allows ALL four panels (header, left, right, footer) to independently expand when focused.

**Key Innovation**: Unlike traditional dual-pane or tri-panel layouts, tmuxplexer features:
- Header panel (top) - expandable to 50% height
- Left/Right panels (middle row) - expandable to 75% width
- Footer panel (bottom) - expandable to 50% height
- All panels use weight-based dynamic sizing for smooth accordion behavior

## Target Platform

**Primary**: Desktop terminals (120+ cols × 30+ rows recommended)
**Minimum**: 80 cols × 24 rows

## The 4-Panel Layout Architecture

### Panel Content Mapping (from PLAN.md)

**Header Panel** (expandable to 50% height):
- Session statistics (total sessions, windows, panes)
- Quick-launch templates for common session types
- Filter/search controls
- Server connection status

**Left Panel** (expandable to 75% width):
- Session list (scrollable)
- Session metadata (name, windows count, created time, attached status)
- Visual indicators (attached, has activity, etc.)
- Keyboard navigation with vim bindings

**Right Panel** (expandable to 75% width):
- Window/pane details for selected session
- Visual pane layout representation
- Per-window actions (rename, kill, etc.)
- Pane tree view

**Footer Panel** (expandable to 50% height):
- Live pane preview (tmux capture-pane output)
- Natural language command input (future: "create 3 panes in session 'dev'")
- Command history
- Quick action buttons

### Weight-Based Accordion System

The layout uses **weights** instead of fixed pixels:

```go
// Default weights: 1:2:1 vertical (25%:50%:25%), 1:1 horizontal (50%:50%)
headerWeight, middleWeight, footerWeight := 1, 2, 1

// When header focused: 2:1:1 (50%:25%:25%)
if m.accordionMode && m.focusedPanel == "header" {
    headerWeight = 2
    middleWeight = 1
    footerWeight = 1
}

// When left/right focused: 1:4:1 (16.67%:66.67%:16.67%)
// Plus horizontal: left gets 3, right gets 1 (75%:25%)
if m.focusedPanel == "left" {
    headerWeight = 1
    middleWeight = 4
    footerWeight = 1
    leftWidth = (totalWidth * 3) / 4
    rightWidth = (totalWidth * 1) / 4
}

// Calculate actual heights:
totalWeight := headerWeight + middleWeight + footerWeight
headerHeight := (availableHeight * headerWeight) / totalWeight
```

**Why this works**: Instant, proportional resizing with no animations or complex calculations.

## Critical Layout Rules (from CLAUDE.md)

**The 4 Golden Rules** - NEVER violate these:

### Rule 1: Always Account for Borders
```go
contentHeight = totalHeight - titleLines - statusLines - 2 // -2 for panel borders!
```

**Visual breakdown**:
```
Total Terminal Height: 25
- Title Bar:           -3
- Status Bar:          -1
- Panel Borders:       -2  ← CRITICAL: top + bottom borders
─────────────────────────
Content Height:        19 ✓
```

### Rule 2: Never Auto-Wrap in Bordered Panels
```go
// ALWAYS truncate text explicitly
maxTextWidth := panelWidth - 4  // -2 borders, -2 padding
text = truncateString(text, maxTextWidth)

// Helper function (add to your codebase):
func truncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-1] + "…"
}
```

### Rule 3: Match Mouse Detection to Layout
```go
// Calculate EXACT same contentStartY as rendering:
contentStartY := 0
if m.config.UI.ShowTitle {
    contentStartY += 2  // Title bar
}
contentStartY += 1  // Tab/menu bar if present

// Use SAME contentHeight calculation as rendering:
contentWidth, contentHeight := m.calculateLayout()

// Then detect clicks using relative Y:
relY := msg.Y - contentStartY
```

### Rule 4: Use Weights, Not Pixels
```go
// DON'T: leftWidth := 60 (fixed)
// DO: leftWidth := (totalWidth * leftWeight) / totalWeights
```

## Current Project Structure

This project was scaffolded from TUITemplate with:
- Layout: `dual_pane` (will upgrade to 4-panel)
- Components: `panel`, `list`, `status`, `dialog`

**Key files** (following TUITemplate pattern):

```
tmuxplexer/
├── config.yaml           # UI settings (colors, borders, mouse)
├── types.go              # Model struct, ALL state lives here
├── model.go              # Init, layout calculations
├── update.go             # Main update dispatcher
├── update_keyboard.go    # Keyboard handling
├── update_mouse.go       # Mouse click handling
├── view.go               # Main view dispatcher
├── view_components.go    # Individual component renderers
├── tmux.go               # (CREATE THIS) Tmux integration
└── PLAN.md               # Comprehensive project plan
```

## Reference Implementation: TUI Showcase Tab 12

The 4-panel layout is **fully implemented and debugged** in:
`/home/matt/projects/TUITemplate/examples/tui-showcase/`

**Key files to reference**:
- `model.go:calculateFourPanelLayout()` - Weight calculation logic
- `view.go:renderFourPanelTab()` - Panel rendering with borders
- `update_mouse.go:handleFourPanelClick()` - Mouse click detection
- `update_keyboard.go` - Focus switching (keys 1,2,3,4, accordion toggle 'a')

**Testing the reference**:
```bash
cd /home/matt/projects/TUITemplate/examples/tui-showcase
./tui-showcase
# Press Tab until you reach "4-Panel" tab (tab 12)
# Try keys: 1, 2, 3, 4 to focus panels, 'a' to toggle accordion
# Click panels with mouse to focus
```

## Immediate Next Steps

### 1. Copy 4-Panel Layout Functions

Copy these from `/home/matt/projects/TUITemplate/examples/tui-showcase/`:

**From `types.go`**:
```go
// Add to model struct:
headerContent []string // Content for header panel
leftContent   []string // Content for left panel
rightContent  []string // Content for right panel
bottomContent []string // Content for footer panel

focusedPanel  string   // "header", "left", "right", "footer"
accordionMode bool     // Enable/disable accordion expansion
```

**From `model.go`**:
```go
func (m model) calculateFourPanelLayout(availableWidth, availableHeight int) (
    headerHeight, middleHeight, footerHeight int,
    leftWidth, rightWidth int,
) {
    // Copy the entire function from showcase
}
```

**From `view.go`**:
```go
func (m model) renderFourPanelTab(totalWidth, totalHeight int) string {
    // Copy the entire function from showcase
}

func (m model) renderDynamicPanel(panelID string, width, height int, content []string) string {
    // Copy the entire function from showcase
}
```

**From `update_keyboard.go`**:
```go
// Add focus switching:
case "1": m.focusedPanel = "left"
case "2": m.focusedPanel = "right"
case "3": m.focusedPanel = "footer"
case "4": m.focusedPanel = "header"
case "a", "A": m.accordionMode = !m.accordionMode
```

**From `update_mouse.go`**:
```go
func (m model) handleFourPanelClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    // Copy the entire function from showcase
}
```

### 2. Create Tmux Integration Layer

Create `tmux.go` with functions to interact with tmux:

```go
package main

import (
    "os/exec"
    "strings"
)

// listSessions returns all tmux sessions
func listSessions() ([]TmuxSession, error) {
    // Run: tmux list-sessions -F "#{session_name}|#{session_windows}|#{session_attached}"
    // Parse output into []TmuxSession
}

// listWindows returns windows for a session
func listWindows(sessionName string) ([]TmuxWindow, error) {
    // Run: tmux list-windows -t sessionName -F "#{window_index}|#{window_name}|#{window_panes}"
}

// listPanes returns panes for a window
func listPanes(sessionName string, windowIndex int) ([]TmuxPane, error) {
    // Run: tmux list-panes -t sessionName:windowIndex -F "#{pane_id}|#{pane_current_command}"
}

// capturePane returns visible content of a pane
func capturePane(paneID string) (string, error) {
    // Run: tmux capture-pane -p -t paneID
}

// attachSession attaches to a session
func attachSession(sessionName string) error {
    // Run: tmux attach-session -t sessionName
}

// Define types:
type TmuxSession struct {
    Name     string
    Windows  int
    Attached bool
    Created  time.Time
}

type TmuxWindow struct {
    Index int
    Name  string
    Panes int
}

type TmuxPane struct {
    ID      string
    Command string
    Active  bool
}
```

### 3. Wire Up Real Data

Replace placeholder content in `model.go:initialModel()`:

```go
// Instead of:
leftContent: []string{"Placeholder", "session", "list"},

// Use:
sessions, err := listSessions()
if err != nil {
    leftContent = []string{"Error loading sessions: " + err.Error()}
} else {
    leftContent = formatSessionList(sessions)
}
```

### 4. Add Scrolling to Session List

Use offset-based scrolling pattern from `/home/matt/projects/TUITemplate/docs/SCROLLING_AND_RESPONSIVE.md`:

```go
// Add to types.go:
type model struct {
    // ...
    sessions      []TmuxSession
    sessionCursor int // Selected session
    sessionOffset int // Scroll offset
}

// In view_components.go:
func (m model) renderSessionList(width, height int) string {
    visibleCount := height - 4 // Account for borders/title
    start := m.sessionOffset
    end := min(start + visibleCount, len(m.sessions))

    var lines []string
    for i := start; i < end; i++ {
        session := m.sessions[i]
        line := formatSessionLine(session)

        if i == m.sessionCursor {
            line = selectedStyle.Render("▶ " + line)
        } else {
            line = normalStyle.Render("  " + line)
        }
        lines = append(lines, line)
    }

    // Add scroll indicators
    if m.sessionOffset > 0 {
        lines = append([]string{dimStyle.Render("↑ more above")}, lines...)
    }
    if end < len(m.sessions) {
        lines = append(lines, dimStyle.Render("↓ more below"))
    }

    return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
```

### 5. Implement Core Actions

**In `update_keyboard.go`**:
```go
case "enter":
    // Attach to selected session
    if m.sessionCursor < len(m.sessions) {
        session := m.sessions[m.sessionCursor]
        return m, tea.Exec(func() tea.Msg {
            attachSession(session.Name)
            return nil
        })
    }

case "n":
    // Create new session (show dialog)
    m.showingDialog = true
    m.dialogTitle = "New Session"
    return m, nil

case "k":
    // Kill selected session
    if m.sessionCursor < len(m.sessions) {
        // Show confirmation dialog first
    }
```

### 6. Add Live Preview

Update footer panel with live pane preview:

```go
// Add timer for periodic updates:
func (m model) Init() tea.Cmd {
    return tea.Batch(
        tea.EnterAltScreen,
        tickCmd(), // Start periodic updates
    )
}

func tickCmd() tea.Cmd {
    return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

// In update.go:
case tickMsg:
    // Update preview of focused pane
    if m.sessionCursor < len(m.sessions) {
        session := m.sessions[m.sessionCursor]
        // Get active pane, capture content, update bottomContent
    }
    return m, tickCmd() // Schedule next tick
```

## Key Resources

1. **Layout Reference**: `/home/matt/projects/TUITemplate/examples/tui-showcase/` (Tab 12)
2. **Critical Rules**: `/home/matt/projects/TUITemplate/CLAUDE.md`
3. **Scrolling Guide**: `/home/matt/projects/TUITemplate/docs/SCROLLING_AND_RESPONSIVE.md`
4. **Project Plan**: `/home/matt/projects/tmuxplexer/PLAN.md`
5. **4-Panel Demo**: `/home/matt/projects/TUITemplate/examples/tui-showcase/4PANEL_DEMO.md`

## Testing Strategy

1. **Test layout first** with placeholder data (like showcase does)
2. **Test accordion** behavior (all 4 panels expanding/contracting)
3. **Test mouse clicks** on all panels (including expanded states)
4. **Add real tmux data** once layout is solid
5. **Test scrolling** with many sessions (20+)
6. **Test terminal resize** (try 80x24, 120x30, 200x50)

## Success Criteria for MVP

- [ ] 4-panel layout renders correctly
- [ ] All panels expand/contract on focus
- [ ] Mouse clicks work on all panels (including expanded areas)
- [ ] Keyboard focus switching works (1,2,3,4)
- [ ] Session list loads from tmux
- [ ] Session list scrolls with arrow keys
- [ ] Can attach to session with Enter
- [ ] Live preview updates in footer panel
- [ ] No layout bugs (headers visible, no overflow, no flickering)

## Anti-Patterns to Avoid

❌ **DON'T** use `.Height()` on lipgloss styles with borders
❌ **DON'T** let text auto-wrap in bordered panels
❌ **DON'T** calculate mouse coordinates differently than rendering
❌ **DON'T** use fixed pixel widths - always use weights
❌ **DON'T** forget to subtract 2 for borders in height calculations

✅ **DO** truncate content to fit available height
✅ **DO** use offset-based scrolling for lists
✅ **DO** match mouse detection math to rendering math exactly
✅ **DO** test with various terminal sizes
✅ **DO** reference the working showcase implementation

---

**Ready to start?** Begin with step 1 (copy 4-panel layout functions) and build incrementally!
