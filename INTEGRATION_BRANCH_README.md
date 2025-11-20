# Integration Work Branch

**Branch:** `feature/tmuxplexer-integration`

## Purpose

This branch contains the work-in-progress integration of tui-launcher and tmuxplexer into a unified workspace manager.

## Current State

### What's Here

1. **tmuxplexer-original/** - Complete copy of tmuxplexer code
   - All Go files (.go)
   - Documentation (README.md, PLAN.md, CLAUDE.md, HOTKEYS.md)
   - Dependencies (go.mod, go.sum)
   - Binary (tmuxplexer)
   - Components, lib, docs, hooks directories

2. **Current tui-launcher code** - Unchanged in root directory
   - model.go, tree.go, spawn.go, types.go, etc.
   - All existing functionality intact

3. **Integration Plan** - INTEGRATION_PLAN.md
   - Detailed 5-phase integration plan
   - Workflow documentation
   - Technical decisions

## Directory Structure

```
tui-launcher/ (root)
├── tmuxplexer-original/     # Tmuxplexer source (reference)
│   ├── model.go
│   ├── tmux.go
│   ├── update_keyboard.go
│   ├── update_mouse.go
│   ├── view.go
│   ├── templates.go
│   ├── claude_state.go
│   ├── components/
│   ├── lib/
│   └── ...
├── model.go                 # TUI Launcher (current)
├── tree.go
├── spawn.go
├── types.go
└── INTEGRATION_PLAN.md      # Integration roadmap
```

## Next Steps

### Phase 1: Code Merge (First Session)
1. Create new directory structure:
   ```
   tabs/
   ├── launch/      # TUI Launcher code
   ├── sessions/    # Tmuxplexer sessions
   └── templates/   # Tmuxplexer templates
   shared/
   ├── tmux.go      # Combined tmux operations
   └── spawn.go     # Combined spawn logic
   ```

2. Extract shared code:
   - Tmux operations (from tmuxplexer/tmux.go + tui-launcher/spawn.go)
   - Configuration loading
   - Common types

3. Create tab routing:
   - Main Update() dispatcher
   - Tab-specific Update() handlers
   - Unified View() composition

### Testing During Integration

**Checkpoint 1:** Launch tab works (preserve existing functionality)
```bash
go build && ./tui-launcher
# Should show current tui-launcher interface
```

**Checkpoint 2:** Sessions tab works (tmuxplexer sessions)
```bash
go build && ./tui-launcher
# Press '2' or Tab to switch to Sessions
# Should list tmux sessions
```

**Checkpoint 3:** Templates tab works (tmuxplexer templates)
```bash
go build && ./tui-launcher
# Press '3' to switch to Templates
# Should show template categories
```

**Checkpoint 4:** All tabs work together
```bash
go build && ./tui-launcher
# Launch commands → auto-switch to Sessions → save as Template
```

## Key Files to Merge

### High Priority (Core Functionality)
- [ ] **tmux.go** - Session/window/pane operations
- [ ] **model.go** - State management (combine both)
- [ ] **types.go** - Type definitions (combine both)
- [ ] **view.go** - Rendering (split by tab)
- [ ] **update_keyboard.go** - Keyboard handling (combine both)

### Medium Priority (Features)
- [ ] **templates.go** - Template loading/saving
- [ ] **claude_state.go** - Claude Code integration
- [ ] **spawn.go** - Already in tui-launcher, enhance with tmuxplexer patterns

### Low Priority (Polish)
- [ ] **update_mouse.go** - Mouse handling (tmuxplexer has good patterns)
- [ ] **styles.go** - Visual consistency
- [ ] **config.go** - Config unification

## Keeping Track

### Preserve from TUI Launcher
- ✅ Hierarchical tree view
- ✅ Multi-select with Space
- ✅ Quick CD feature
- ✅ YAML config
- ✅ Responsive layout (3-pane)
- ✅ Project organization

### Preserve from Tmuxplexer
- ✅ Session list with status indicators
- ✅ Live preview (tmux capture-pane)
- ✅ Template system (COLSxROWS layouts)
- ✅ Save session as template
- ✅ Template wizard
- ✅ Claude Code tracking
- ✅ Popup mode (Ctrl+B O)

### New Combined Features
- 🆕 Launch commands → Sessions tab
- 🆕 Save Launch commands as template
- 🆕 Templates in Launch tab
- 🆕 Cross-tab workflows

## Useful Commands

```bash
# Switch to integration branch
git checkout feature/tmuxplexer-integration

# Build on integration branch
go build

# Compare with original tmuxplexer
diff -u tmuxplexer-original/model.go model.go

# Test tmuxplexer standalone (reference)
cd tmuxplexer-original && ./tmuxplexer

# See what changed from main
git diff main

# Commit progress
git add . && git commit -m "Integration: ..."

# Push to remote (for backup/sharing)
git push -u origin feature/tmuxplexer-integration
```

## Integration Checklist

### Phase 1: Merge (Current)
- [x] Create work branch
- [x] Copy tmuxplexer code
- [ ] Create tab directory structure
- [ ] Extract shared code
- [ ] Create tab routing
- [ ] Test basic tab switching

### Phase 2: Enhanced Launch
- [ ] Add launch feedback
- [ ] Auto-switch to Sessions after launch
- [ ] Visual indicators

### Phase 3: Unified Config
- [ ] Convert templates.json to YAML
- [ ] Merge configs
- [ ] Migration script

### Phase 4: Cross-Tab Features
- [ ] Launch → Sessions integration
- [ ] Sessions → Templates (save)
- [ ] Launch → Templates (convert)

### Phase 5: Polish
- [ ] Popup mode
- [ ] Consistent hotkeys
- [ ] Visual polish
- [ ] Documentation

---

**Status:** Ready for Phase 1 implementation
**Last Updated:** 2025-11-20
**Branch:** feature/tmuxplexer-integration
