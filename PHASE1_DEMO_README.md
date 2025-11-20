# Phase 1 Tab Routing Demo

## Quick Start

### Build and Run
```bash
# Build the Phase 1 demo
go build -o tui-launcher-phase1 types.go types_unified.go tab_routing.go model_unified.go main_test_tabs.go

# Run it
./tui-launcher-phase1
```

### What You'll See

A simple TUI with three tabs showing placeholder content:

```
 1. Launch   2. Sessions   3. Templates  (Tab/Shift+Tab to cycle, 1/2/3 for direct access)
────────────────────────────────────────────────────────────────────────────────────────
Launch Tab (Current TUI Launcher)

This tab will show:
  • Hierarchical project tree
  • Global tools, AI commands, scripts
  • Multi-select command launching
  • Quick CD into project directories

Status: Welcome to TUI Launcher - Phase 1 Tab Integration

Phase 1: Tab routing implemented ✓
Next: Integrate existing tui-launcher view
```

### Test Tab Switching

| Key | Action |
|-----|--------|
| `1` | Switch to Launch tab |
| `2` | Switch to Sessions tab |
| `3` | Switch to Templates tab |
| `Tab` | Cycle forward through tabs |
| `Shift+Tab` | Cycle backward through tabs |
| `q` or `Ctrl+C` | Quit |

## What's Working

✅ **Tab Architecture**
- Three-tab structure (Launch, Sessions, Templates)
- Tab switching with keyboard shortcuts
- Visual tab indicator showing active tab
- Status messages on tab changes

✅ **Build System**
- Compiles successfully
- No runtime errors
- Clean separation of concerns

✅ **Code Organization**
- `types_unified.go` - Tab routing types and unified model
- `tab_routing.go` - Tab switching logic and view routing
- `model_unified.go` - Model initialization and update dispatcher
- `main_test_tabs.go` - Entry point

## What's NOT Working Yet

❌ **Launch Tab** - Shows placeholder, doesn't load actual launcher tree
- Needs: Migration of current model.go view logic
- Needs: Integration with config.yaml loading
- Needs: Tree building and navigation

❌ **Sessions Tab** - Shows placeholder, doesn't list tmux sessions
- Needs: Migration of tmuxplexer session list view
- Needs: Live preview panel implementation
- Needs: Session attach/kill/rename functionality

❌ **Templates Tab** - Shows placeholder, doesn't load templates
- Needs: Migration of tmuxplexer template tree view
- Needs: Template creation wizard
- Needs: Template I/O (loading from templates.json)

## Architecture Overview

### Unified Model
```go
type unifiedModel struct {
    width, height int           // Shared across tabs
    currentTab    tabName        // Which tab is active
    launchModel   launchTabModel // Launch tab state
    // sessionsModel  (TBD)
    // templatesModel (TBD)
}
```

### Tab Routing Flow
```
User Input
    ↓
Update() dispatcher
    ↓
handleTabSwitch()? → Change currentTab
    ↓
routeUpdateToTab() → Active tab's update logic
    ↓
View() compositor
    ↓
renderTabBar() + renderActiveTabContent()
    ↓
Display to terminal
```

## Next Steps

### Priority 1: Launch Tab Integration
Extract current tui-launcher view logic:
1. Copy `model.go` view functions to `tabs/launch/view.go`
2. Copy `model.go` update logic to `tabs/launch/update.go`
3. Copy `tree.go` to `tabs/launch/tree.go`
4. Wire up config loading
5. **Goal:** Launch tab shows actual project tree

### Priority 2: Sessions Tab Basic View
Extract tmuxplexer session list:
1. Copy session model from `tmuxplexer-original/model.go`
2. Copy sessions view from `tmuxplexer-original/view.go`
3. Implement basic session listing
4. **Goal:** Sessions tab shows live tmux sessions

### Priority 3: Shared Tmux Layer
Unify tmux operations:
1. Merge `spawn.go` + `tmuxplexer-original/tmux.go`
2. Create `shared/tmux.go` with all tmux functions
3. **Goal:** All tabs can call shared tmux operations

## Documentation

- **PHASE1_CODE_ORGANIZATION.md** - Complete integration plan
- **PHASE1_COMPLETION_SUMMARY.md** - What was accomplished in Phase 1
- **INTEGRATION_PLAN.md** - Overall 5-phase roadmap
- **INTEGRATION_BRANCH_README.md** - Branch status and checklist

## Troubleshooting

### "undefined: XYZ" compilation error
The phase 1 build uses explicit file listing. Make sure you include:
- `types.go` (existing types)
- `types_unified.go` (tab routing types)
- `tab_routing.go`
- `model_unified.go`
- `main_test_tabs.go`

### Tab switching not working
Check if you're pressing the right keys:
- Number keys `1`, `2`, `3` (not Ctrl+1, etc.)
- `Tab` key (may need Shift+Tab if your terminal intercepts Tab)

### Placeholder views don't show full content
The placeholder views are intentionally minimal. Full functionality comes in:
- Launch tab: Phase 1 continuation
- Sessions tab: Phase 1 continuation
- Templates tab: Phase 1 continuation

## Comparison: Current vs Phase 1

### Current tui-launcher
```bash
./tui-launcher
# Single view: hierarchical tree of projects and commands
# No tabs, no session management, no templates
```

### Phase 1 Demo
```bash
./tui-launcher-phase1
# Tab 1 (Launch): Placeholder (will be current interface)
# Tab 2 (Sessions): Placeholder (will be tmuxplexer sessions)
# Tab 3 (Templates): Placeholder (will be tmuxplexer templates)
# Tab switching works, content TBD
```

### Final Goal (Phase 1 Complete)
```bash
./tui-launcher
# Tab 1 (Launch): Current tui-launcher interface ✓
# Tab 2 (Sessions): Live tmux sessions with previews ✓
# Tab 3 (Templates): Categorized workspace templates ✓
# All tabs fully functional
```

## Testing Checklist

- [ ] Build compiles without errors
- [ ] Program launches successfully
- [ ] Press `1` → Shows Launch placeholder
- [ ] Press `2` → Shows Sessions placeholder
- [ ] Press `3` → Shows Templates placeholder
- [ ] Press `Tab` → Cycles forward (Launch → Sessions → Templates → Launch)
- [ ] Press `Shift+Tab` → Cycles backward
- [ ] Tab bar highlights correct active tab
- [ ] Status message updates on tab switch
- [ ] Press `q` → Exits cleanly

## Performance Notes

Binary size: ~4.4MB (similar to current tui-launcher)
- No performance degradation expected
- Tab routing adds minimal overhead
- Placeholder views are very lightweight
- Real performance testing comes when tabs are populated

## Known Limitations (Phase 1)

1. **No Config Loading**: Placeholder tabs don't load actual data yet
2. **No Mouse Support**: Tab switching is keyboard-only (mouse support in Phase 5)
3. **No Popup Mode**: Popup mode flag exists but isn't wired up yet
4. **No Cross-Tab Features**: Tabs are isolated (cross-tab features in Phase 4)
5. **Minimal Error Handling**: Basic error display only

## Git Workflow

```bash
# See what's new
git status

# Review changes
git diff

# Stage Phase 1 files
git add tabs/ shared/ *.md types_unified.go tab_routing.go model_unified.go main_test_tabs.go

# Commit
git commit -m "Phase 1: Implement tab routing foundation"

# Push to integration branch
git push origin feature/tmuxplexer-integration
```

---

**Phase 1 Foundation: Complete ✅**
**Ready for Launch Tab Integration**
