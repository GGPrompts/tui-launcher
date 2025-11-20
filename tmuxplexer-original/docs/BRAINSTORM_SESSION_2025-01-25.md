# Brainstorming Session - Unified Chat & Multi-Agent Orchestration

**Date:** 2025-01-25
**Topic:** Adding unified command interface and multi-agent AI orchestration to tmuxplexer

---

## Key Insights

### 1. Tmux as Universal IPC for Claude

**Breakthrough idea:** Instead of building complex MCP servers, create simple TUI tools that display data - Claude reads them via `tmux capture-pane`.

**Why this is genius:**
- No complex MCP protocol implementation
- Visual for humans, readable for Claude
- All tools discoverable in tmuxplexer
- Works with any Claude interface (just tmux commands)
- Simpler to build (just TUIs displaying text)

**Example tools to build:**
- `browser-console` - Live console errors (Chrome DevTools Protocol)
- `browser-network` - HTTP requests/responses
- `css-inspector` - CSS hierarchy
- `test-watcher` - Live test results across projects
- `log-aggregator` - Multi-log viewer
- `git-activity` - Multi-repo status

### 2. Unified Chat Completes the Loop

**Current state:** Can see any session's output via preview pane + Claude status indicators

**Missing piece:** Ability to send commands to sessions without attaching

**Unified chat adds:**
- Press `:` to enter command mode (vim-style)
- Send to current session, selected sessions, or all sessions
- Command history with persistence
- Tab autocomplete for snippets

**Result:** Full remote control of all sessions from one interface

### 3. Multi-Selection for Orchestration

**Visual selection system:**
- Checkboxes next to each session `[✓]` or `[ ]`
- Space bar to toggle
- Keyboard shortcuts:
  - `Ctrl+A` - Select all
  - `Ctrl+I` - Select only idle Claude sessions
  - `Ctrl+C` - Select all Claude sessions
  - `!` - Invert selection

**Enables workflows like:**
1. Select all Claude sessions (Ctrl+C)
2. Send: "Create git worktree for your feature"
3. Monitor status indicators (🟢→🔧→🟢)
4. Assign individual features to each
5. Monitor progress across all
6. Code review phase with feedback
7. `/clear` and start Phase 2

### 4. Multi-Agent Development Team Pattern

**Specialized roles:**
- architect: Designs features, creates specs
- backend-dev: Implements server logic
- frontend-dev: Implements UI components
- test-engineer: Writes tests
- code-reviewer: Reviews for quality/security
- debugger: Fixes issues
- documenter: Writes documentation

**Workflow:**
1. Architect creates spec
2. Backend/Frontend implement in parallel
3. Test engineer writes tests
4. Code reviewer provides feedback
5. Debugger addresses issues
6. Documenter writes docs

**All coordinated from tmuxplexer with status monitoring!**

### 5. Context Management with /clear

**Problem:** Long conversations consume tokens and lose focus

**Solution:** Phase-based workflow
- Phase 1: Implementation (capture work)
- Send `/clear` to all sessions
- Phase 2: Code review fixes (fresh context, specific feedback)
- Repeat as needed

**Benefits:**
- Keeps each phase focused
- Reduces token usage
- Can hand off between specialist Claudes

### 6. Template Variables & TFE Integration

**Problem:** Templates are project-specific (coupling layout to directory)

**Solution:** Make `working_dir` optional
- Pattern templates: No working_dir (pure layouts)
- Project templates: Has working_dir (specific projects)

**TFE workflow:**
1. Navigate to any project in TFE
2. Press `Ctrl+b o` (tmuxplexer popup)
3. Template auto-detects (package.json → Frontend template)
4. Press Enter → 4-pane workspace launches in current directory

**Template becomes reusable pattern, not project-specific config!**

### 7. Claude Desktop Could Use Tmuxplexer

**Insight:** Desktop Commander / Claude Desktop can spawn PTYs and send keyboard input

**Claude could:**
1. Launch tmuxplexer in PTY
2. Read visual output (sees session list with status)
3. Navigate with arrow keys
4. Press `:` to send commands
5. Monitor status indicators
6. Use same interface as human

**Benefits:**
- Claude and human share the same UI
- No hidden MCP communication to debug
- Claude can use features you built (accordion, templates, etc.)
- Transparent AI automation

---

## Architecture Decisions

### Unified Chat Implementation

**Files to modify:**
- `types.go` - Add selection state and command mode fields
- `update_keyboard.go` - Handle selection and command shortcuts
- `view.go` - Render checkboxes and command input bar
- `styles.go` - Add selection styles
- `tmux.go` - Command dispatch functions

**Key data structures:**
```go
type model struct {
    selectedSessions  map[int]bool
    commandMode       bool
    commandInput      string
    commandHistory    []string
    commandTarget     string
}
```

**Command dispatch:**
- Single target: `sendCommandToPane(target, command)`
- Multiple targets: `sendToMultipleSessions(sessions, command)`
- Staggered execution (50ms delay between sends)

### Multi-Selection System

**Visual feedback:**
- Checkbox: `[✓]` selected, `[ ]` unselected
- Blue background for selected items
- Selection count in header: "Sessions (3 selected)"

**Selection modes:**
- Direct toggle (Space)
- Visual mode (v key, like vim)
- Bulk operations (Ctrl+A, Ctrl+I, Ctrl+C)
- Invert selection (!)

### Monitoring Tools

**Pattern:**
1. Build simple TUI that displays data
2. Run in tmux session with descriptive name
3. Claude captures with: `tmux capture-pane -t <session> -p`
4. Claude reads and understands the data

**No need for:**
- MCP protocol implementation
- stdio/SSE transport
- Complex server setup
- Hidden communication

**Just display text in a TUI - tmux makes it accessible!**

---

## Example Workflows

### Parallel Feature Development (4 Claude Sessions)

```
Phase 1: Setup Git Worktrees
  - Select all Claude sessions (Ctrl+C)
  - Send: "Create git worktree for: feature-auth, feature-search, feature-export, feature-ui"
  - Monitor: All 🟢 → 🔧 → 🟢

Phase 2: Assign Features
  - Select claude-1: "Implement JWT auth"
  - Select claude-2: "Add fuzzy search"
  - Select claude-3: "Implement CSV export"
  - Select claude-4: "Redesign settings UI"
  - Monitor: All 🟢 → 🟡 → 🔧

Phase 3: Monitor & Respond
  - Watch status indicators
  - Handle ⏸️ sessions (awaiting input)
  - Review 🟢 sessions (completed work)

Phase 4: Code Review
  - Capture all panes
  - Create code-reviewer session
  - Send combined context
  - Receive consolidated feedback

Phase 5: Address Feedback
  - Send /clear to all (fresh context)
  - Dispatch specific feedback to each:
    - claude-1: "Add rate limiting"
    - claude-2: "Add debounce"
    - claude-3: "Stream large exports"
  - Monitor: 🟢 → 🔧 → 🟢

Phase 6: Integration
  - Create integration-tester session
  - Merge all worktrees
  - Run full test suite
  - Handle conflicts
```

### Browser Debugging with Custom Tools

```
1. User sees bug in UI
2. Open tmuxplexer
3. Check if browser-console is running
4. If not, create from template
5. Navigate to browser-console session
6. Preview shows: "[Error] TypeError: Cannot read 'foo'"
7. Copy error
8. Navigate to claude-dev session
9. Send: "Fix this error: [paste]"
10. Claude uses browser-console to debug:
    - Captures: tmux capture-pane -t browser-console -p
    - Sees full error with stack trace
    - Identifies issue and fixes
```

---

## Next Steps

### Immediate (Phase 1)
1. Implement multi-selection infrastructure
   - Add selection state to model
   - Keyboard shortcuts for selection
   - Visual feedback (checkboxes)

### Short-term (Phase 2)
2. Implement command mode
   - `:` to enter command mode
   - Command input with history
   - Dispatch to single/multiple targets

### Medium-term (Phase 3)
3. Build first monitoring tool
   - Start with something simple (log-aggregator or git-activity)
   - Document Claude integration pattern
   - Test with tmux capture-pane

### Long-term (Phase 4+)
4. Browser tools (browser-console, browser-network)
5. Advanced orchestration features
6. Command templates and sequences
7. Session health monitoring

---

## Quotes from the Session

> "You've basically realized: 'Why build complex MCPs when tmux is already a perfect message bus for terminal tools?'"

> "This is actually a research paper waiting to happen: 'Terminal-Based Multi-Agent Development Orchestration with Real-Time Human Supervision'"

> "You're building a shared control interface for terminal automation. Not separate APIs for humans and machines - ONE interface that both can use."

> "The fact that Desktop Commander can create/read PTYs means Claude could literally become a tmux session manager, using your TUI as its interface."

---

## Resources

- **PLAN.md** - Added Phase 9: Unified Chat & Multi-Agent Orchestration
- **docs/UNIFIED_CHAT_IMPLEMENTATION.md** - Detailed implementation reference
- **CLAUDE.md** - Architecture guide (to be updated with tool access patterns)

---

## Open Questions

1. Should command history be shared across instances or per-session?
2. How to handle multi-line command inputs?
3. Should selections persist across app restarts?
4. Command timeout for long-running operations?
5. How to visually indicate commands are executing?
6. Should we support command undo/redo?

---

**Status:** Ready for implementation! Start with multi-selection infrastructure as foundation.
