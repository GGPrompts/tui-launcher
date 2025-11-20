# Tmuxplexer Phase 5: Claude Code Dashboard & Hooks Integration

## Project Context

Building **tmuxplexer** - a modern TUI tmux session manager with workspace templates in Go (Bubble Tea framework).

**Location:** `~/projects/tmuxplexer`
**Repository:** https://github.com/GGPrompts/Tmuxplexer

## Current Status: Ready for Phase 5! 🚀

### ✅ Phases 1-4 Complete (All Working!)
- **Phase 1:** 4-Panel Accordion Layout
- **Phase 2:** Workspace Templates
- **Phase 3:** Session Management (attach, kill, auto-refresh)
- **Phase 4:** Live Preview & Window Navigation
- **Latest:** Context-aware status bar + attach outside tmux fixed

### 🎉 NEW: Claude Hooks Integration System (Just Built!)

We just created a complete **real-time Claude state tracking system** using hooks! This is the foundation for Phase 5.

## What Was Just Created

### 7 New Files Ready to Use:

1. **`hooks/state-tracker.sh`** - Bash script that captures all Claude hook events and writes state to JSON files
2. **`claude_state.go`** - Go code ready to integrate into tmuxplexer (read Claude state, parse JSON, format display)
3. **`hooks/install.sh`** - One-command installation script
4. **`hooks/test-hooks.sh`** - Test suite to verify hooks work
5. **`docs/claude-hooks-integration.md`** - Complete technical guide (75+ lines)
6. **`docs/HOOKS-QUICKREF.md`** - Quick reference cheat sheet
7. **`hooks/claude-settings-hooks.json`** - Config template for ~/.claude/settings.json

### How the Hooks System Works

```
Claude Event (user prompt, tool use, completion)
       ↓
Hook fires (7 different hook types available)
       ↓
state-tracker.sh executes
       ↓
Writes JSON to /tmp/claude-code-state/<session-id>.json
       ↓
Tmuxplexer reads state file
       ↓
Displays real-time status with icons
```

### Available Hooks (All Configured!)

| Hook Type | When It Fires | Status Written | Use Case |
|-----------|---------------|----------------|----------|
| SessionStart | Claude starts | `idle` | Initialize tracking |
| UserPromptSubmit | User sends message | `processing` | User activity |
| PreToolUse | Before tool runs | `tool_use` | Claude working |
| PostToolUse | After tool completes | `working` | Processing results |
| Stop | Response complete | `awaiting_input` | Ready for input |
| SubagentStop | Subagent finishes | varies | Multi-agent |
| Notification | System notify | varies | Awaiting-input bell |

### State File Format

Each Claude session writes to `/tmp/claude-code-state/<session-id>.json`:

```json
{
  "session_id": "abc123",
  "status": "tool_use",
  "current_tool": "Edit",
  "working_dir": "/home/matt/projects/tmuxplexer",
  "last_updated": "2025-08-24T14:30:22Z",
  "tmux_pane": "%42",
  "pid": 12345,
  "hook_type": "pre-tool",
  "details": {
    "event": "tool_starting",
    "tool": "Edit"
  }
}
```

### Status Indicators Ready to Display

- 🟢 **Idle** - Just started
- 🟡 **Processing** - Thinking about prompt
- 🔧 **Tool Use** - Executing Edit/Read/Bash
- ⚙️ **Working** - Processing tool results
- ⏸️ **Awaiting Input** - Ready for next prompt
- ⚪ **Stale** - No updates >5 sec (hung/crashed)

## Phase 5 Goals

Transform tmuxplexer into the **ultimate Claude Code dashboard** with:

### Priority 1: Claude State Integration (THIS SESSION!)
Display real-time Claude status for each session using the hooks system.

**Target UI:**
```
┌─ Sessions ────────────────────────────────┐
│ TEMPLATES                                  │
│   Simple Dev (2x2)                         │
│   Frontend Dev (2x2)                       │
│                                            │
│ ACTIVE SESSIONS                            │
│ ► ● project-1-claude                       │
│     🔧 Tool Use (Edit) │ ~/project-1      │
│   ○ project-2-claude                       │
│     ⏸️ Awaiting Input │ ~/project-2       │
│   ○ docs-claude                            │
│     🟡 Processing │ ~/docs                │
├────────────────────────────────────────────┤
│ 🔧 TOOL USE: Edit                          │
│ File: model.go                             │
│ Working: ~/project-1                       │
│                                            │
│ Status updates every 2 seconds            │
└────────────────────────────────────────────┘
```

### Priority 2: Scrollable Preview (Quick Win!)
Add PgUp/PgDn scrolling to footer preview panel.

### Priority 3: Send Commands to Pane
Interactive command input to send text to any Claude session.

## Step-by-Step Implementation Plan

### Step 1: Install & Test Hooks (5-10 minutes)

```bash
cd ~/projects/tmuxplexer

# 1. Install the hooks system
./hooks/install.sh

# 2. Verify installation
cat ~/.claude/settings.json | grep -A 20 hooks

# 3. Test hooks work
./hooks/test-hooks.sh

# 4. Start a Claude session and verify state files appear
tmux new -s test-claude
claude
# In another terminal:
watch -n 1 'ls -lh /tmp/claude-code-state/'
# Send a prompt in Claude, watch state file update!
```

**Success criteria:**
- ✅ Hooks installed in ~/.claude/settings.json
- ✅ Test script shows all hooks working
- ✅ State files appear in /tmp/claude-code-state/
- ✅ State updates when you interact with Claude

### Step 2: Integrate claude_state.go (30-45 minutes)

The file `claude_state.go` is ready - just needs integration:

```bash
# File is already created at:
# ~/projects/tmuxplexer/claude_state.go
```

**Tasks:**

1. **Add ClaudeState to TmuxSession type** (types.go):
```go
type TmuxSession struct {
    Name       string
    Windows    int
    Attached   bool
    Created    string
    LastActive string
    // NEW: Add Claude state
    ClaudeState *ClaudeState  // nil if not a Claude session
}
```

2. **Call getClaudeStateForSession() in listSessions()** (tmux.go):
```go
func listSessions() ([]TmuxSession, error) {
    // ... existing code to get sessions ...

    for i := range sessions {
        // Try to get Claude state for this session
        claudeState := getClaudeStateForSession(sessions[i].Name)
        sessions[i].ClaudeState = claudeState
    }

    return sessions, nil
}
```

3. **Display Claude state in left panel** (model.go, updateLeftPanelContent()):
```go
for i, item := range m.listItems {
    if item.Type == "session" {
        prefix := "  "
        if i == m.selectedItem {
            prefix = "► "
        }

        icon := "○"
        if item.Session.Attached {
            icon = "●"
        }

        // Add Claude status indicator
        statusIcon := ""
        if item.Session.ClaudeState != nil {
            statusIcon = getClaudeStatusIcon(item.Session.ClaudeState.Status)
        }

        lines = append(lines, fmt.Sprintf("%s%s %s%s",
            prefix, icon, statusIcon, item.Name))

        // Optional: Show status text on second line
        if item.Session.ClaudeState != nil && i == m.selectedItem {
            statusText := item.Session.ClaudeState.FormatStatus()
            lines = append(lines, fmt.Sprintf("    %s", statusText))
        }
    }
}
```

4. **Display detailed Claude info in right panel** (model.go, updateRightPanelContent()):
```go
if selectedItem.Type == "session" && selectedItem.Session.ClaudeState != nil {
    lines = append(lines, "")
    lines = append(lines, "🤖 CLAUDE CODE")
    lines = append(lines, "")

    cs := selectedItem.Session.ClaudeState
    lines = append(lines, fmt.Sprintf("Status: %s", cs.FormatStatus()))

    if cs.CurrentTool != "" {
        lines = append(lines, fmt.Sprintf("Tool: %s", cs.CurrentTool))
    }

    if cs.WorkingDir != "" {
        lines = append(lines, fmt.Sprintf("Directory: %s", cs.WorkingDir))
    }

    lines = append(lines, fmt.Sprintf("Updated: %s", cs.LastUpdated.Format("15:04:05")))
    lines = append(lines, "")
}
```

5. **Build and test**:
```bash
go build -o tmuxplexer
./tmuxplexer
```

**Success criteria:**
- ✅ Tmuxplexer compiles without errors
- ✅ Sessions with Claude show status icons
- ✅ Right panel shows Claude state details
- ✅ Status updates every 2 seconds (auto-refresh)

### Step 3: Enhanced Status Bar (15 minutes)

Update context-aware status to show Claude info:

```go
// In update_keyboard.go getContextualStatusMessage()
if selectedItem.Type == "session" {
    // Check for Claude state
    if selectedItem.Session.ClaudeState != nil {
        cs := selectedItem.Session.ClaudeState
        baseStatus := fmt.Sprintf("Claude: %s", cs.FormatStatus())

        switch m.focusedPanel {
        case "left":
            return fmt.Sprintf("%s | [Enter] Attach | [d/K] Kill | [q] Quit", baseStatus)
        // ... other cases
        }
    }

    // ... existing session status logic
}
```

### Step 4: Scrollable Preview (BONUS - if time allows)

Add scroll state and PgUp/PgDn handling:

```go
// In types.go
type model struct {
    // ... existing fields
    previewScroll int       // Current scroll position
    previewLines  []string  // Full captured content
}

// In update_keyboard.go
case "pgup":
    if m.focusedPanel == "footer" {
        m.previewScroll -= 10
        if m.previewScroll < 0 {
            m.previewScroll = 0
        }
        m.updateFooterContent()
        m.statusMsg = fmt.Sprintf("Preview scroll: line %d", m.previewScroll)
    }

case "pgdown":
    if m.focusedPanel == "footer" {
        m.previewScroll += 10
        maxScroll := len(m.previewLines) - m.footerHeight
        if m.previewScroll > maxScroll {
            m.previewScroll = maxScroll
        }
        m.updateFooterContent()
        m.statusMsg = fmt.Sprintf("Preview scroll: line %d", m.previewScroll)
    }
```

## Testing Checklist

### Basic Functionality
```bash
# 1. Create test Claude sessions
tmux new -s web-app -d "cd ~/projects/web && claude"
tmux new -s api -d "cd ~/projects/api && claude"

# 2. Launch tmuxplexer
./tmuxplexer

# 3. Verify display
# ✓ Sessions show with Claude indicators
# ✓ Status icons appear (🟢🟡🔧⚙️⏸️)
# ✓ Right panel shows Claude details

# 4. Test state updates
# In web-app session: Send a prompt
# Watch tmuxplexer: Status should change to 🟡 Processing
# Then 🔧 Tool Use when Claude uses tools
# Then ⏸️ Awaiting Input when done

# 5. Test auto-refresh
# Leave tmuxplexer running
# Interact with Claude in tmux
# Verify status updates every 2 seconds
```

### Edge Cases
- [ ] Non-Claude sessions show normally (no state)
- [ ] Stale state shows ⚪ after 5 seconds
- [ ] Multiple Claude sessions tracked independently
- [ ] State persists across tmuxplexer restarts
- [ ] Handles missing state files gracefully

## Common Issues & Solutions

### Issue: Hooks not firing
```bash
# Check hooks are installed
cat ~/.claude/settings.json | grep hooks

# Check script is executable
ls -l ~/projects/tmuxplexer/hooks/state-tracker.sh
chmod +x ~/projects/tmuxplexer/hooks/state-tracker.sh

# Test manually
echo '{"session_id":"test"}' | ~/projects/tmuxplexer/hooks/state-tracker.sh user-prompt-submit
ls /tmp/claude-code-state/
```

### Issue: State files not appearing
```bash
# Check directory exists and is writable
ls -ld /tmp/claude-code-state/
mkdir -p /tmp/claude-code-state/
chmod 755 /tmp/claude-code-state/

# Check for errors in hook execution
claude --debug hooks
```

### Issue: Status not updating in tmuxplexer
```bash
# Verify auto-refresh is working
# Should update every 2 seconds (tickCmd in update.go)

# Check state file is recent
ls -lh /tmp/claude-code-state/*.json
cat /tmp/claude-code-state/*.json | jq .
```

## File Locations Reference

### New Files (Just Created)
```
~/projects/tmuxplexer/
├── claude_state.go                      # Go integration (READY TO USE)
├── hooks/
│   ├── state-tracker.sh                 # Hook script
│   ├── install.sh                       # Installation
│   ├── test-hooks.sh                    # Test suite
│   ├── README.md                        # Hooks docs
│   └── claude-settings-hooks.json       # Config template
└── docs/
    ├── claude-hooks-integration.md      # Complete guide
    └── HOOKS-QUICKREF.md                # Quick reference
```

### Existing Files to Modify
```
types.go              # Add ClaudeState field to TmuxSession
tmux.go               # Call getClaudeStateForSession()
model.go              # Display Claude state in panels
update_keyboard.go    # Enhanced status messages
```

### External Files
```
~/.claude/settings.json                  # Hooks configuration
/tmp/claude-code-state/<session-id>.json # State files (auto-created)
```

## Success Criteria for Phase 5

By end of this session:
- [x] Hooks installed and tested
- [ ] claude_state.go integrated into tmuxplexer
- [ ] Claude sessions show status icons in left panel
- [ ] Right panel shows detailed Claude state
- [ ] Status updates automatically every 2 seconds
- [ ] Context-aware status bar shows Claude info
- [ ] (Bonus) Scrollable preview with PgUp/PgDn
- [ ] All tests pass
- [ ] Documentation updated
- [ ] Committed and pushed to GitHub

## Documentation to Update

After implementation:
- [ ] README.md - Add Claude integration section
- [ ] CLAUDE.md - Document new claude_state.go patterns
- [ ] Add screenshots showing Claude status indicators
- [ ] Update keyboard shortcuts if adding new keys

## Commit Message Template

```
feat: Phase 5 - Claude Code Dashboard Integration

Integrated real-time Claude status tracking using hooks system.

Features:
- Displays Claude session status with icons (🟢🟡🔧⚙️⏸️)
- Shows current tool, working directory, last update
- Auto-refreshes Claude state every 2 seconds
- Detects stale/hung sessions (>5 seconds)
- Status updates via hooks system (7 hook types)

Implementation:
- Added claude_state.go for state file reading
- Extended TmuxSession with ClaudeState field
- Enhanced left/right panels with Claude info
- Updated context-aware status bar
- Installed hooks via state-tracker.sh script

State file location: /tmp/claude-code-state/<session-id>.json

Testing:
✓ Multiple Claude sessions tracked independently
✓ Status icons update in real-time
✓ Stale detection works (shows ⚪ after 5s)
✓ Graceful handling of non-Claude sessions
✓ All hooks tested and working

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Future Enhancements (Phase 6+)

Once Claude integration is working:
1. **Send commands to pane** - Type prompts from tmuxplexer
2. **Command history** - Recent prompts per session
3. **Multi-line input** - Longer prompts with editor
4. **Tool execution log** - Show recent tool uses
5. **Session grouping** - Group related Claude sessions
6. **Custom indicators** - User-defined status icons
7. **Sound/visual alerts** - Notify when Claude finishes

## Questions for This Session

1. Should we show Claude status in left panel inline or on hover?
2. Scrollable preview - essential or can wait?
3. Want to add sound when Claude awaiting input (system bell)?
4. Should we cache state reads or query on every render?
5. Display format for "Tool Use" - just icon or full details?

## Quick Start for Next Session

```bash
# 1. Install hooks
cd ~/projects/tmuxplexer
./hooks/install.sh

# 2. Test hooks work
./hooks/test-hooks.sh

# 3. Review claude_state.go
cat claude_state.go

# 4. Start integration
# Modify types.go, tmux.go, model.go per steps above

# 5. Build and test
go build -o tmuxplexer
./tmuxplexer
```

---

**Current State:** Hooks system complete! Ready to integrate into tmuxplexer UI.

**What to do first:** Install hooks with `./hooks/install.sh`, test with `./hooks/test-hooks.sh`, then start integrating claude_state.go!

🚀 Let's build the ultimate Claude Code dashboard! 🎉
