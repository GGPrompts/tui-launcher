# Next Session: Tree Symbol Cleanup & Keyboard Improvements

## Context: What We Just Completed

**Session completed on 2025-01-07**

Successfully implemented expandable session tree with filtering:

1. ✅ **Session Tree Navigation** - Hierarchical view of sessions → windows → panes
2. ✅ **Session Filtering** - Press 'f' to cycle: All → AI → Attached → Detached
3. ✅ **Precise Pane Targeting** - Navigate to specific panes, preview updates, commands go to exact pane
4. ✅ **Universal Command Support** - Send commands to ANY session (not just AI), with safety checks
5. ✅ **Single-pane Session Support** - No need to expand single-pane sessions
6. ✅ **Selection Persistence** - Selection survives auto-refresh (2-second ticker)
7. ✅ **Smart Collapse** - Press ← on pane/window to collapse parent session

**Git Status:**
- Working tree: Modified files not yet committed
- New features ready for testing

---

## 🔍 Issues to Fix This Session

### **Issue 1: Too Many Arrow/Caret Symbols in Tree (Confusing!)**

**Current State:**
The tree view uses multiple arrow/triangle symbols that create visual confusion:

```
▼ ●► terminal-tabs ◆ 🤖               ← TOO MANY ARROWS!
  ├─   0: backend (1 panes) ●
  │  └─   ► Pane 0: bash ●            ← What does ► mean here vs above?
```

**Symbols currently used:**
- `▶/▼` - Expansion indicator (collapsed/expanded session)
- `►` - Selection indicator (currently selected item)
- `◆` - Current session marker (where tmuxplexer is running)
- `●/○` - Attached/detached status
- Tree connectors: `├─`, `└─`, `│`

**Problem:**
1. Too many triangular symbols (`▶`, `▼`, `►`) - hard to distinguish at a glance
2. Selection indicator (`►`) looks too similar to expansion indicator (`▶`)
3. Visual hierarchy is unclear
4. Hard to quickly find what's selected

**Task:**
Clean up the visual hierarchy to reduce confusion:

**Option A - Use highlighting instead of symbols for selection:**
```
▼ ● terminal-tabs ◆ 🤖              ← Not selected
  ├─ 0: backend (1 panes) ●
  │  └─ Pane 0: bash ●              ← Selected (highlighted/bold)
```

**Option B - Use different selection symbol:**
```
▼ ● terminal-tabs ◆ 🤖
  ├─ 0: backend (1 panes) ●
  │  └─ → Pane 0: bash ●            ← Use → instead of ►
```

**Option C - Remove selection symbol, rely on styling only:**
```
▼ ● terminal-tabs ◆ 🤖
  ├─ 0: backend (1 panes) ●
  │  └─ Pane 0: bash ●              ← Bold/inverse/cyan background = selected
```

**Files to modify:**
- `model.go` - `updateSessionsContent()` function (lines ~237-347)
  - Around line 254: Sets prefix `"► "` when selected
  - Around line 297: Session line building with expansion indicator
  - Around lines 310-343: Window and pane rendering
- `view.go` - `renderDynamicPanel()` if styling changes needed (lines ~360-400)
- `styles.go` - Add new style for selected items if needed

**Current selection logic (model.go:253-256):**
```go
prefix := "  "
if selected && m.focusState == FocusSessions {
    prefix = "► "
}
```

**Implementation approach:**
1. Remove the `"► "` prefix from selected items
2. Instead, add a styling tag (like "SELECTED:" already used)
3. In view.go, apply bold/inverse/background color to SELECTED items
4. Test that selection is still obvious when navigating
5. Ensure it works in both focus states

---

### **Issue 2: Keyboard Shortcuts - 'd' for Detach, 'k' for Kill**

**Current State:**
```
'd' or 'D' → Kill session (destructive!)
No detach shortcut
```

**Problem:**
1. **Unintuitive**: 'd' feels like "detach" (common, safe operation)
2. **Too easy to kill**: Accidentally killing sessions is destructive
3. **Against conventions**:
   - Vim uses 'd' for "delete" but requires motion/confirmation
   - 'k' for "kill" follows common Unix conventions (kill command)
   - tmux uses `Ctrl+b d` for detach

**Task:**
Swap keyboard bindings to be more intuitive and safe:

**New bindings:**
- `d` → **Detach** from current session (safe, returns to previous session or shell)
- `k` → **Kill** session (destructive, should require confirmation)

**Implementation:**

**Step 1: Add detach functionality**
- When pressing 'd' on a session, detach from it
- Only works if you're currently attached to that session
- If you're in tmuxplexer's session, show message "Can't detach from current session (press q to quit)"
- If you're viewing a different session, show message "Session is already detached" or "Not currently attached to this session"

**Step 2: Move kill to 'k' with confirmation**
- Change 'd' handler to 'k'
- Add confirmation prompt: "Kill session 'name'? (y/n): "
- Prevent accidental destruction

**Step 3: Add detach tmux command**
File: `tmux.go`
```go
// detachSession detaches from a tmux session
func detachSession(sessionName string) error {
    cmd := exec.Command("tmux", "detach-client", "-s", sessionName)
    return cmd.Run()
}
```

**Files to modify:**
- `update_keyboard.go` - handleKeyPress function
  - Current 'd' handler around lines 464-479
  - Add new 'd' detach handler
  - Change 'd' to 'k' for kill
  - Add confirmation for kill
- `tmux.go` - Add detachSession function
- `types.go` - May need confirmation mode state (or reuse inputMode)

**Current code (update_keyboard.go:464-479):**
```go
case "d", "D":
    // Kill/delete session (Sessions tab) or delete template (Templates tab)
    if true && len(m.sessions) > 0 && m.selectedSession < len(m.sessions) {
        session := m.sessions[m.selectedSession]
        m.statusMsg = "Killing session: " + session.Name + "..."
        return m, m.killSessionCmd(session.Name)
    }
    // ... template deletion code ...
```

**New code structure:**
```go
case "d", "D":
    // Detach from session
    if m.sessionsTab == "sessions" && len(m.sessionTreeItems) > 0 {
        item := m.sessionTreeItems[m.selectedSession]
        if item.Type == "session" && item.Session != nil {
            // Detach logic
            // Check if attached, check if current session, etc.
        }
    }

case "k", "K":
    // Kill session (with confirmation)
    if m.sessionsTab == "sessions" && len(m.sessionTreeItems) > 0 {
        item := m.sessionTreeItems[m.selectedSession]
        if item.Type == "session" && item.Session != nil {
            // Prompt for confirmation
            m.inputMode = "kill_confirm"
            m.inputPrompt = fmt.Sprintf("Kill session '%s'? (y/n): ", item.Session.Name)
            return m, nil
        }
    }
```

**Confirmation handling:**
- Add "kill_confirm" mode to inputMode handling
- On 'y' press, kill the session
- On 'n' or ESC, cancel

---

## 📁 Key Files Reference

**Tree Symbol Rendering:**
- `model.go` - `updateSessionsContent()` (lines 237-347)
  - Line 254: Selection prefix logic
  - Line 261-262: Expansion indicator logic
  - Lines 297-343: Tree rendering with symbols
- `view.go` - `renderDynamicPanel()` (lines 360-400) - Styling tags
- `styles.go` - Style definitions (if adding new selected style)

**Keyboard Handling:**
- `update_keyboard.go` - Main keyboard handler
  - Lines 464-479: Current 'd' kill handler
  - Need to add 'k' handler
  - Input mode handling for confirmation
- `tmux.go` - Tmux commands
  - Add `detachSession()` function
  - Existing `killSession()` function

**Styling tags already in use (for reference):**
- "SELECTED:" - Selected template items
- "CURRENT:" - Current session (cyan, bold)
- "CLAUDE:" - Claude sessions (orange, bold)
- "DETAILS:header:" - Section headers
- "DETAILS:detail:" - Detail lines

---

## 🎯 Success Criteria

### Issue 1: Tree Symbol Cleanup
- [ ] Selection is visually obvious without using `►` symbol
- [ ] Less visual clutter (fewer arrows/triangles)
- [ ] Easy to scan and find selected item
- [ ] Current session (◆) still stands out
- [ ] Expansion state (▶/▼) still clear
- [ ] Tree structure (├─ └─) still readable
- [ ] Works in both normal and popup mode

### Issue 2: Keyboard Shortcuts
- [ ] 'd' detaches from session (safe operation)
- [ ] 'd' shows appropriate message if can't detach
- [ ] 'k' prompts for confirmation before killing
- [ ] 'k' actually kills session after 'y' confirmation
- [ ] ESC or 'n' cancels kill confirmation
- [ ] Help text updated (if you have any)
- [ ] Works on session items in tree
- [ ] Doesn't break template deletion (if 'd' is used there)

---

## 🚀 Getting Started

```bash
# 1. Check current status
cd ~/projects/tmuxplexer
git status

# 2. Find current tree symbol rendering
grep -n "prefix.*►" model.go
grep -n "expansionIndicator" model.go

# 3. Find current 'd' key handler
grep -n 'case "d"' update_keyboard.go

# 4. Build and test current behavior
go build -o tmuxplexer
./tmuxplexer

# 5. Note what symbols you see and which are confusing
# 6. Note what happens when you press 'd' on a session
```

---

## 🧪 Testing Checklist

After making changes, test:

### Tree Symbols:
- [ ] Navigate with ↑/↓ - selection is obvious
- [ ] Expand session with → - expansion indicator clear
- [ ] Collapse with ← - state changes visually
- [ ] Multiple sessions visible - can quickly scan
- [ ] Sessions, windows, and panes all clearly distinguished
- [ ] Current session marker (◆) still visible and obvious
- [ ] AI session indicators (🤖 🔮 ✨) still show
- [ ] Attached/detached (●/○) still show

### Keyboard Shortcuts:
- [ ] Press 'd' on attached session → detaches (if possible)
- [ ] Press 'd' on detached session → shows appropriate message
- [ ] Press 'd' on current session → shows "can't detach" message
- [ ] Press 'k' on session → shows confirmation prompt
- [ ] Type 'y' after confirmation → kills session
- [ ] Type 'n' after confirmation → cancels
- [ ] Press ESC during confirmation → cancels
- [ ] 'd' on template (if applicable) → still works for template deletion

---

## 💡 Design Suggestions

### For Symbol Cleanup:

**Recommended approach:**
1. Remove `►` prefix for selection
2. Use inverse/bold styling on entire line for selected items
3. Keep expansion indicators (▶/▼) as they're functional
4. Keep tree connectors (├─ └─) as they show structure
5. Keep status symbols (●/○ ◆ 🤖) as they're informational

**Example of cleaner design:**
```
▼ ● terminal-tabs ◆ 🤖              ← Expanded, attached, current, AI
  ├─ 0: backend (1 panes) ●
  │  └─ Pane 0: bash ●              ← Inverse/bold = selected (no arrow!)
  └─ 2: logs (1 panes) ●
     └─ Pane 0: tfe ●
```

### For Detach:

**Edge cases to handle:**
- Can't detach from the session tmuxplexer is running in (that's what 'q' does)
- Can't detach if not currently attached (session is already detached)
- In popup mode, detaching might not make sense (maybe disable?)

**Messages:**
- "✓ Detached from session: {name}"
- "⚠️ Session '{name}' is not attached"
- "⚠️ Can't detach from current session (press q to quit)"

### For Kill Confirmation:

**Prompt style:**
```
╭ Confirm Kill Session ──────────╮
│ Kill session 'myproject'?      │
│                                 │
│ [y] Yes, kill it               │
│ [n] No, cancel                 │
│ [ESC] Cancel                   │
╰─────────────────────────────────╯
```

Or simple inline:
```
Kill session 'myproject'? (y/n): _
```

---

## 📝 Notes

- The tree rendering is in `updateSessionsContent()` around line 237-347
- Selection styling uses tags processed in `renderDynamicPanel()`
- Look at how templates use "SELECTED:" tag for reference
- Current session uses "CURRENT:" tag → cyan, bold
- Claude sessions use "CLAUDE:" tag → orange, bold
- Detach command: `tmux detach-client -s <session-name>`
- Kill command: already exists in `tmux.go`
- Test both normal mode and popup mode (`./tmuxplexer --popup`)

**Priority:**
1. Symbol cleanup first (more impactful)
2. Keyboard shortcuts second (once tree is clear)

**Good luck!** 🚀
