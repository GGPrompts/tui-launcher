# Unified Chat Implementation Reference

**Quick reference for implementing the unified chat and multi-selection features**

Created: 2025-01-25
Status: Design phase - ready for implementation

---

## Overview

Transform tmuxplexer into a command center by adding:
1. **Unified command interface** - Send commands to any session without attaching
2. **Multi-selection system** - Select multiple sessions with checkboxes
3. **Multi-agent orchestration** - Coordinate multiple Claude sessions in parallel

---

## Keyboard Shortcuts

### Selection Mode

```
Space       Toggle selection of current session
v           Enter/exit visual selection mode
Ctrl+A      Select all sessions
Ctrl+D      Deselect all
Ctrl+I      Select only idle Claude sessions
Ctrl+C      Select only Claude sessions (any status)
!           Invert selection
```

### Command Mode

```
:           Enter command mode
Esc         Cancel command mode
Enter       Send command to target(s)
↑/↓         Navigate command history
Tab         Autocomplete from snippets
```

### Navigation (with selections)

```
j/k         Move cursor (preserves selections)
Space       Toggle selection while navigating
```

---

## Implementation Checklist

### Phase 1: Multi-Selection Infrastructure

**Files to modify:**
- `types.go` - Add selection state fields
- `update_keyboard.go` - Add selection shortcuts
- `view.go` - Render checkboxes and selection styling
- `styles.go` - Add selection styles

**Key additions:**

```go
// In types.go - model struct
type model struct {
    // ... existing fields

    // Multi-selection
    selectedSessions  map[int]bool
    selectionMode     bool
    lastSelectedIndex int

    // Command mode
    commandMode       bool
    commandInput      string
    commandHistory    []string
    commandTarget     string
    commandTargetMode string // "current", "selected", "all"
}
```

**Helper methods needed:**
- `toggleSelection(index int)`
- `isSelected(index int) bool`
- `getSelectedCount() int`
- `clearSelections()`
- `selectAll()`
- `selectOnlyIdle()`
- `selectOnlyClaude()`
- `getSelectedSessions() []TmuxSession`

### Phase 2: Command Mode

**Files to modify:**
- `update_keyboard.go` - Handle `:` key and command input
- `view.go` - Render command input bar
- `tmux.go` - Add command dispatch functions

**Key functions:**

```go
// In tmux.go
func sendCommandToPane(target string, command string) tea.Cmd
func (m model) sendToMultipleSessions(sessions []TmuxSession, command string) tea.Cmd
func (m model) dispatchCommand() tea.Cmd
```

**Message types needed:**

```go
type commandSentMsg struct {
    target  string
    command string
    output  string
    err     error
}

type multiCommandSentMsg struct {
    results map[string]commandResult
    command string
}

type commandResult struct {
    target  string
    command string
    err     error
}
```

### Phase 3: Visual Feedback

**Files to modify:**
- `styles.go` - Add selection styles
- `view.go` - Render session lines with checkboxes

**Styles needed:**

```go
selectedStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("62")).  // Blue
    Foreground(lipgloss.Color("230")). // White
    Bold(true)

focusedStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("236")). // Dark gray
    Foreground(lipgloss.Color("255")). // Bright white
    Bold(true)

selectedAndFocusedStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("63")).  // Brighter blue
    Foreground(lipgloss.Color("255")).
    Bold(true)
```

**Rendering pattern:**

```go
func (m model) renderSessionLine(index int, session TmuxSession) string {
    checkbox := "[ ]"
    if m.isSelected(index) {
        checkbox = "[✓]"
    }

    // Build line with checkbox, icon, name, status
    line := fmt.Sprintf("%s %s %s %s %s",
        checkbox, icon, name, statusIcon, statusText)

    // Apply appropriate style
    if index == m.selectedSession && m.isSelected(index) {
        return m.styles.selectedAndFocused.Render(line)
    } else if index == m.selectedSession {
        return m.styles.focused.Render(line)
    } else if m.isSelected(index) {
        return m.styles.selected.Render(line)
    }

    return line
}
```

---

## Command History Persistence

**Location:** `~/.config/tmuxplexer/command_history.json`

**Format:**
```json
{
  "history": [
    "git status",
    "npm test",
    "go build -o tmuxplexer",
    "Create git worktree for feature-auth"
  ],
  "max_entries": 100
}
```

**Functions needed:**

```go
func loadCommandHistory() ([]string, error)
func saveCommandHistory(history []string) error
func addToHistory(command string) []string
func getPreviousCommand(currentIndex int) string
func getNextCommand(currentIndex int) string
```

---

## Command Snippets

**Location:** `~/.config/tmuxplexer/snippets.json`

**Format:**
```json
{
  "snippets": {
    "gs": "git status",
    "gc": "git commit -m \"",
    "gp": "git push",
    "npr": "npm run",
    "test": "go test ./...",
    "build": "go build -o"
  }
}
```

**Implementation:**

```go
// In update_keyboard.go - when in command mode
case "tab":
    if snippet, ok := m.snippets[m.commandInput]; ok {
        m.commandInput = snippet
    }
```

---

## Multi-Agent Orchestration Examples

### Example 1: Parallel Feature Development

```go
// User workflow in tmuxplexer:

// 1. Select all Claude sessions
// Press: Ctrl+C (selects all Claude sessions)

// 2. Send setup command
// Press: ':'
// Type: "Create git worktree: feature-auth, feature-search, feature-export"
// Press: Enter

// 3. Wait for completion (watch status indicators)
// All sessions: 🟢 → 🔧 → 🟢

// 4. Assign individual features
// Navigate to claude-1, deselect others
// Press: ':'
// Type: "Implement JWT authentication with refresh tokens"
// Press: Enter
// (Repeat for claude-2, claude-3 with different features)

// 5. Monitor progress
// Watch status indicators change: 🟢 → 🟡 → 🔧 → ⏸️ or 🟢

// 6. Handle awaiting input
// Navigate to ⏸️ session
// Read preview
// Press: ':'
// Type: response
// Press: Enter
```

### Example 2: Code Review Phase

```bash
# After features complete, capture all work:

# For each Claude session:
tmux capture-pane -t claude-1 -S -3000 -p > claude-1-work.txt
tmux capture-pane -t claude-2 -S -3000 -p > claude-2-work.txt
tmux capture-pane -t claude-3 -S -3000 -p > claude-3-work.txt

# Create code reviewer session
# Send combined context for review
# Receive feedback

# Clear all Claude sessions for Phase 2
# Select all Claude sessions (Ctrl+C)
# Press: ':'
# Type: "/clear"
# Press: Enter

# Send targeted feedback to each
# (navigate and send individually)
```

---

## tmux Commands Used

### Send Keys
```bash
tmux send-keys -t <session>:<window>.<pane> "command" Enter
```

### Capture Pane
```bash
# Last 100 lines
tmux capture-pane -t <session> -p

# Last 3000 lines (full history)
tmux capture-pane -t <session> -S -3000 -p

# Specific pane
tmux capture-pane -t session:window.pane -p
```

### List Sessions with Details
```bash
tmux list-sessions -F "#{session_id}|#{session_name}|#{session_attached}|#{session_created}"
```

---

## Testing Plan

### Unit Tests

1. **Selection logic**
   - Toggle selection
   - Select all/deselect all
   - Select by criteria (idle, Claude-only)
   - Invert selection

2. **Command dispatch**
   - Single target
   - Multiple targets
   - Error handling
   - Command history

3. **Target detection**
   - Current session
   - Selected sessions
   - All sessions
   - Filtered sessions

### Integration Tests

1. **End-to-end workflows**
   - Select → Send → Verify
   - Multi-select → Broadcast → Verify all
   - Command history → Recall → Send

2. **Error scenarios**
   - Session doesn't exist
   - Command fails
   - Network timeout (for remote sessions)

### Manual Testing

1. **Visual feedback**
   - Checkboxes toggle correctly
   - Selection count updates
   - Styles apply properly
   - Command input bar appears

2. **Keyboard shortcuts**
   - All shortcuts work
   - No conflicts
   - Intuitive behavior

3. **Multi-agent orchestration**
   - Create 3-4 test sessions
   - Select and send commands
   - Monitor status changes
   - Verify commands executed

---

## Performance Considerations

### Optimization Strategies

1. **Staggered command dispatch**
   - 50ms delay between sends
   - Prevents tmux server overload
   - Maintains command order

2. **Async operations**
   - Command dispatch returns tea.Cmd
   - UI stays responsive
   - Results update incrementally

3. **Efficient rendering**
   - Only re-render changed sessions
   - Debounce rapid selection changes
   - Cache rendered lines

### Benchmarks to Track

- Time to select all sessions (should be < 10ms)
- Time to send command to 10 sessions (should be < 1s)
- UI responsiveness during command dispatch
- Memory usage with 100+ command history

---

## Future Enhancements

### Command Templates

```json
{
  "templates": {
    "deploy": [
      "go build -o app",
      "go test ./...",
      "git push origin main",
      "notify-send 'Deploy complete'"
    ],
    "setup-feature": [
      "git worktree add ../{{feature-name}}",
      "cd ../{{feature-name}}",
      "git checkout -b {{feature-name}}"
    ]
  }
}
```

### Session Groups

```json
{
  "groups": {
    "all-features": ["claude-1", "claude-2", "claude-3"],
    "monitoring": ["browser-console", "test-watcher", "log-aggregator"],
    "dev-team": ["architect", "backend-dev", "frontend-dev"]
  }
}
```

### Smart Targeting

```
Send to @claude-sessions
Send to @idle-sessions
Send to @group:all-features
Send to @pattern:feature-*
```

---

## Questions to Resolve

- [ ] Should command history be shared across all instances or per-session?
- [ ] How to handle very long command inputs (multi-line)?
- [ ] Should selection persist across app restarts?
- [ ] Command timeout for long-running operations?
- [ ] Visual indicator while commands are executing?
- [ ] Undo/redo for selections?

---

## References

- Main plan: `PLAN.md` Phase 9
- Architecture: `CLAUDE.md`
- Keyboard shortcuts: `HOTKEYS.md` (to be updated)
- Templates: `templates.go`
- Command dispatch: `tmux.go`

---

**Ready to implement!** Start with Phase 1 (multi-selection infrastructure) as the foundation.
