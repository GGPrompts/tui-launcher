# Foreground vs Background Spawn Design

## Current Behavior (v0.2.0)

**All spawns are foreground:**
- User launches command → TUI quits → command runs
- Simple and clean, but TUI can't stay resident

## Proposed: Background Spawn Mode

### Use Cases

1. **Sequential launching**: Launch backend, then frontend, then database monitor without restarting TUI
2. **Persistent workspace**: Keep launcher open as a "mission control" for your dev environment
3. **Mixed workflows**: Launch some things in background, then quit launcher when done

### Implementation Options

#### Option A: Per-Item Spawn Mode
Add `background: true` to config:

```yaml
projects:
  - name: My Web App
    path: ~/projects/webapp
    commands:
      - name: Backend
        command: go run .
        spawn: tmux-split-h
        background: true  # TUI stays open after launch

      - name: Open Editor
        command: nvim .
        spawn: current-pane
        background: false  # TUI quits (default)
```

**Behavior:**
- `background: true` → Spawn command, TUI stays open, can launch more
- `background: false` → Spawn command, TUI quits (current behavior)

#### Option B: Runtime Toggle
Add hotkey to toggle spawn mode:

```
'b' key - Toggle background mode
Status bar shows: "Mode: Background" or "Mode: Foreground"
```

**Behavior:**
- Background mode ON → All launches keep TUI open
- Background mode OFF → All launches quit TUI (current)

#### Option C: Smart Detection
Automatically determine based on spawn mode:

```go
func shouldStayOpen(spawnMode spawnMode) bool {
    switch spawnMode {
    case spawnTmuxSplitH, spawnTmuxSplitV, spawnTmuxWindow:
        return true  // Background - command runs in another pane
    case spawnCurrentPane:
        return false // Foreground - replacing current pane
    case spawnXtermWindow:
        return true  // Background - separate window
    default:
        return false
    }
}
```

**Behavior:**
- Tmux splits/windows → Auto-background (TUI stays)
- Current pane replacement → Auto-foreground (TUI quits)
- Can override with config or hotkey

### Technical Considerations

**If TUI stays open:**
- Need to track spawned processes (optional)
- TUI uses a pane/window permanently
- Can become process manager (show running/stopped)
- More complex lifecycle management

**Challenges:**
- What if user launches in `current-pane` but TUI is in background mode?
  - Solution: Warn or auto-switch to split mode
- How to "finish" and quit launcher?
  - Solution: 'q' still quits, or auto-quit when selections are empty

### Recommended Approach

**Start with Option C (Smart Detection) + Runtime Toggle**

1. **Smart defaults:**
   - Tmux splits/windows → Background (TUI stays)
   - Current pane → Foreground (TUI quits)
   - Xterm windows → Background (TUI stays)

2. **Add 'b' key toggle:**
   - Override smart defaults
   - User control for edge cases
   - Show in status bar

3. **Future: Per-item config (v0.3.0+):**
   - Add `background: true/false` to YAML
   - Most flexible, but more config complexity

### Example Workflow

**With Background Mode:**
```
1. Launch TUI Launcher (tl)
2. Navigate to "Backend" → Press Enter
   → Backend spawns in tmux split
   → TUI stays open (smart: tmux split = background)
3. Navigate to "Frontend" → Press Enter
   → Frontend spawns in another split
   → TUI still open
4. Press 'q' to quit launcher
   → Backend and Frontend keep running
```

**With Foreground Mode (current):**
```
1. Launch TUI Launcher (tl)
2. Select Backend + Frontend with Space
3. Press Enter
   → Both spawn in batch with layout
   → TUI quits immediately
```

### Integration with Multi-Select

Background mode works great with multi-select:
- Select 3 commands
- Press Enter → All spawn
- TUI stays open (background mode)
- Select 2 more commands
- Press Enter → 2 more spawn
- Press 'q' when done

## Process Management (Future)

If TUI stays resident, could add:
- Show running processes in status bar
- 'k' key to kill selected process
- 'r' key to restart failed process
- Color indicators (🟢 running, 🔴 stopped)

This turns launcher into full workspace manager!

## Config Schema Updates

```yaml
# Global default
settings:
  default_background: true  # Keep TUI open by default

projects:
  - name: My Web App
    commands:
      - name: Backend
        command: go run .
        spawn: tmux-split-h
        background: true  # Override global default
```

## Hotkeys

- **'b'** - Toggle background/foreground mode
- **'q'** - Quit launcher (even in background mode)
- **'Q'** - Quit launcher AND kill all spawned processes (future)

---

**Status:** Design proposal for v0.3.0
**Dependencies:** Core 3-pane layout must be stable first
**Complexity:** Medium - requires process tracking
**Value:** High - major UX improvement for dev workflows
