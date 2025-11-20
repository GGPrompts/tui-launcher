# TUI Launcher + Tmuxplexer Integration Plan

## Vision: Unified Workspace Manager

Combine the best of both tools into a single, powerful TUI for managing:
- **Projects** (tui-launcher) - Organized command launching with multi-select
- **Sessions** (tmuxplexer) - Live tmux session management with previews
- **Templates** (tmuxplexer) - Workspace layouts for instant multi-pane setups

---

## Current State Analysis

### TUI Launcher Strengths
✅ **Hierarchical organization** - Projects, tools, AI, scripts in a tree
✅ **Multi-select launching** - Space to select, Enter to batch launch
✅ **Context-aware spawning** - Different modes (split-h, split-v, window)
✅ **Quick CD** - Press Enter on projects to CD into them
✅ **Responsive layout** - Adapts to terminal size (desktop/compact/mobile)
✅ **YAML config** - Easy to edit and share

### Tmuxplexer Strengths
✅ **Session management** - List, attach, kill tmux sessions
✅ **Live previews** - tmux capture-pane shows actual session content
✅ **Workspace templates** - COLSxROWS layouts (2x2, 3x3, 4x2, etc.)
✅ **Auto-refresh** - Updates every 2 seconds
✅ **Claude Code integration** - Shows Claude status in sessions
✅ **Scrollable previews** - Full scrollback history (PgUp/PgDn)
✅ **Template wizard** - Interactive creation mode
✅ **Popup mode** - Ctrl+B O spawns TUI in 90% screen overlay (perfect for quick access!)

### Current Pain Points

**TUI Launcher:**
- ❌ No visibility into what's running
- ❌ Can't see session previews
- ❌ No way to attach to existing sessions
- ❌ Commands launch and disappear (confusing UX)

**Tmuxplexer:**
- ❌ Only handles tmux sessions (not general commands)
- ❌ No hierarchical organization
- ❌ No project-based context
- ❌ JSON config is harder to edit than YAML

---

## Integration Architecture

### Proposed Tab Structure

```
┌────────────────────────────────────────────────────────────┐
│ [Launch] [Sessions] [Templates]                            │
├─────────────────────┬──────────────────────────────────────┤
│                     │                                      │
│  Tab-specific       │  Info/Preview Pane                   │
│  content            │                                      │
│                     │                                      │
├─────────────────────┴──────────────────────────────────────┤
│ Footer: Hotkeys for current tab                           │
└────────────────────────────────────────────────────────────┘
```

### Tab 1: Launch (Current TUI Launcher)
**Left Pane:** Global Tools | Projects tree
**Right Pane:** Info pane with item details
**Actions:**
- Space: Select commands
- Enter: Launch selected OR CD into project
- Shows project commands when in project directory

### Tab 2: Sessions (From Tmuxplexer)
**Left Pane:** Active tmux sessions tree
- Group by: All, AI (Claude), Attached, Detached
- Visual indicators: 🟢 attached, 🔴 detached
- Show window count per session

**Right Pane:** Live preview of selected session
- tmux capture-pane output
- Scrollable with PgUp/PgDn
- Auto-refresh every 2s
- Switch windows with ←/→

**Actions:**
- Enter: Attach to session
- d/K: Kill session
- r: Rename session
- s: Save as template
- Ctrl+R: Refresh

### Tab 3: Templates (From Tmuxplexer)
**Left Pane:** Template categories tree
- Categorize by: Work, Personal, AI, Dev, etc.
- Show layout type (2x2, 3x3, etc.)

**Right Pane:** Template preview/editor
- Show pane layout diagram
- List commands for each pane
- Working directory

**Actions:**
- Enter: Create session from template
- n: New template (wizard)
- e: Edit in $EDITOR
- Delete: Remove template

---

## Implementation Phases

### Phase 1: Code Merge & Refactor
**Goal:** Combine codebases without breaking existing functionality

**Tasks:**
1. **Create new repo structure:**
   ```
   workspace-manager/
   ├── main.go                  # Entry point with tab routing
   ├── types.go                 # Combined type definitions
   ├── model.go                 # Combined model state
   ├── tabs/
   │   ├── launch/              # TUI Launcher code
   │   │   ├── view.go
   │   │   ├── update.go
   │   │   └── tree.go
   │   ├── sessions/            # Tmuxplexer sessions
   │   │   ├── view.go
   │   │   ├── update.go
   │   │   └── preview.go
   │   └── templates/           # Tmuxplexer templates
   │       ├── view.go
   │       ├── update.go
   │       └── wizard.go
   ├── shared/
   │   ├── tmux.go              # Tmux operations
   │   ├── spawn.go             # Spawn logic
   │   └── config.go            # Config loading
   └── config/
       ├── launcher.yaml        # TUI Launcher config
       └── templates.json       # Tmuxplexer templates
   ```

2. **Unify model struct:**
   - Combine tui-launcher's model + tmuxplexer's model
   - Add `currentTab` field ("launch", "sessions", "templates")
   - Shared fields: width, height, err, statusMsg
   - Tab-specific fields in sub-structs

3. **Create tab routing:**
   - Numbers 1/2/3 or Tab key to switch tabs
   - Each tab handles its own Update() and View()
   - Share common components (borders, footers)

### Phase 2: Enhanced Launch Tab
**Goal:** Add feedback and session awareness to Launch tab

**Tasks:**
1. **Add launch feedback:**
   - Show "Launching X commands..." before quit
   - Display spawn mode being used
   - Add 500ms delay to see message
   - Show errors if spawn fails

2. **Session awareness:**
   - After launching, switch to Sessions tab
   - Show newly created session in preview
   - Option to attach immediately

3. **Visual feedback improvements:**
   - Progress indicators for batch launches
   - Success/error messages
   - Command validation before launch

### Phase 3: Unified Config System
**Goal:** Single YAML config for everything

**Tasks:**
1. **Convert templates.json to YAML:**
   ```yaml
   templates:
     - name: "Frontend Dev (2x2)"
       category: "Web Development"
       description: "Full frontend workspace"
       working_dir: ~/projects/my-app
       layout: "2x2"
       panes:
         - command: "claude-code ."
           title: "Claude AI"
         - command: "nvim"
           title: "Editor"
         - command: "npm run dev"
           title: "Dev Server"
         - command: "lazygit"
           title: "Git"
   ```

2. **Integrate with launcher config:**
   - Templates can reference projects
   - Projects can have default templates
   - Shared categories across both

3. **Migration script:**
   - Convert existing templates.json to YAML
   - Merge with launcher config
   - Validate and backup

### Phase 4: Cross-Tab Features
**Goal:** Make tabs work together seamlessly

**Tasks:**
1. **Launch → Sessions:**
   - After launching commands, auto-switch to Sessions tab
   - Highlight newly created session
   - Option to attach immediately

2. **Sessions → Templates:**
   - "Save as template" creates new template from session
   - Preserves pane layout and commands
   - Adds to templates.yaml

3. **Templates → Launch:**
   - Templates can be added to launcher tree
   - Show templates in Launch tab under "Workspaces"
   - Quick access without switching tabs

### Phase 5: Polish & UX
**Goal:** Smooth, intuitive experience

**Tasks:**
1. **Popup mode integration:**
   - Detect when running in tmux popup (check `$TMUX_POPUP`)
   - Switch sessions instead of attach (avoid nested tmux)
   - Bind to Ctrl+B O for global access
   - Use 90% screen size (or configurable)
   - Quick escape on second 'q' press

2. **Consistent hotkeys:**
   - Tab/1/2/3: Switch tabs
   - Enter: Primary action (attach/launch/create)
   - Space: Select (Launch tab only)
   - d/K: Delete/kill
   - e: Edit config/template
   - q: Quit (or hide popup if in popup mode)

3. **Visual consistency:**
   - Same border styles across tabs
   - Unified color scheme
   - Consistent info pane layout

4. **Smart defaults:**
   - Remember last tab
   - Auto-refresh sessions
   - Preserve scroll positions

---

## Technical Decisions

### 1. Config Format
**Decision:** Use YAML for everything
**Rationale:**
- More readable than JSON
- Supports comments
- Easier to edit manually
- Consistent with launcher config

### 2. Code Organization
**Decision:** Tab-based modules in `tabs/` directory
**Rationale:**
- Clear separation of concerns
- Each tab is self-contained
- Easy to test independently
- Share common code via `shared/`

### 3. Tab Switching
**Decision:** Numbers 1/2/3 + Tab key
**Rationale:**
- Fast direct access with numbers
- Tab for cycling through
- Consistent with tmuxplexer's panel focus

### 4. Naming
**Decision:** Rename to "workspace-manager" or keep "tui-launcher"?
**Options:**
- `workspace-manager` - More descriptive, broader scope
- `tui-launcher` - Keep existing name, backwards compatible
- `tmux-workspace` - Focus on tmux integration

**Recommendation:** `workspace-manager` (better reflects unified purpose)

---

## Migration Path

### For Existing TUI Launcher Users
1. Config stays in `~/.config/tui-launcher/config.yaml`
2. Binary stays at `~/.local/bin/tui-launcher`
3. Wrapper script `tl` still works
4. Launch tab looks identical
5. Add Sessions/Templates tabs as bonus features

### For Existing Tmuxplexer Users
1. Convert `templates.json` to YAML automatically
2. Sessions tab works the same
3. Gain organized command launching
4. Keep all existing hotkeys in Sessions tab

### New Users
1. Simple install script
2. Example configs for both tabs
3. Guided tour on first launch
4. Templates for common workflows

---

## Success Metrics

**Must Have:**
- ✅ All existing tui-launcher features work
- ✅ All existing tmuxplexer features work
- ✅ Tab switching is smooth
- ✅ Config migration is automatic
- ✅ No breaking changes to existing workflows

**Should Have:**
- Launch → Sessions auto-switch
- Session → Template save
- Unified YAML config
- Visual feedback for launches

**Nice to Have:**
- Templates in Launch tab
- Cross-tab search
- Session filtering/grouping
- Template categories

---

## Timeline Estimate

**Phase 1 (Code Merge):** 1-2 sessions
- Combine codebases
- Tab routing
- Basic tab switching

**Phase 2 (Enhanced Launch):** 1 session
- Launch feedback
- Session awareness

**Phase 3 (Unified Config):** 1 session
- YAML conversion
- Migration script

**Phase 4 (Cross-Tab):** 1 session
- Auto-switch
- Save as template

**Phase 5 (Polish):** 1 session
- Hotkey consistency
- Visual polish

**Total:** 5-6 sessions (or ~1-2 weeks at current pace)

---

## Next Steps

1. **Discuss approach** - Do you want to:
   - Start fresh repo (`workspace-manager`)?
   - Merge into `tui-launcher`?
   - Keep separate and share code?

2. **Choose starting phase** - Which to tackle first:
   - Phase 1 (merge now)?
   - Phase 2 (fix launch feedback first)?
   - Phase 3 (config unification)?

3. **Test tmuxplexer features** - Which do you use most:
   - Session previews?
   - Template wizard?
   - Claude status tracking?

4. **Config preferences** - Keep both configs or merge immediately?

---

**Status:** Planning
**Decision Needed:** Integration approach
**Next Action:** Await user direction
