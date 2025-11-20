# TUI Launcher - Development Plan

**A visual terminal launcher for managing projects, TUI tools, and batch command execution with tmux integration**

---

## ✅ Recently Completed (2025-11-20)

### Launch Tab Improvements
**Status:** All priority issues resolved! ✅

**Completed:**
- ✅ 'e' key (edit config) - Fixed using `tea.ExecProcess` pattern from TFE
- ✅ Launch commands - Simplified to foreground-only (like TFE) for single commands
- ✅ Added 'd' key - Toggle between Foreground and Detached modes
- ✅ Multi-select spawning - Works in both modes now
- ✅ Visual feedback - Mode indicator in header shows "Foreground" or "Detached"

**Working Keyboard Shortcuts:**
- ✅ 'e' - Edit config (opens micro/nano/vim, restarts launcher on exit)
- ✅ Enter - Launch commands (foreground or detached based on mode)
- ✅ Space - Selection/expansion
- ✅ Tab - Pane switching
- ✅ 'i' - Info toggle (mobile mode)
- ✅ 'd' - Toggle Foreground/Detached mode
- ✅ 'c' - Clear selections

**Spawn Modes (Simplified):**
- ✅ **Foreground mode (default)** - Single command runs in current terminal, exits launcher
- ✅ **Foreground + multi-select** - Spawns each as tmux window, exits launcher (user in tmux with multiple windows)
- ✅ **Detached mode** - Single/multi spawns tmux windows with `-d` flag, stays in launcher
- ✅ **Profiles** - Use tmuxplexer templates for complex multi-pane layouts

---

## Vision

Create a tree-based TUI launcher that allows:
- **Visual organization** of projects, tools, scripts, and AI commands
- **Multi-select spawning** (Space to select, Enter to launch)
- **Tmux integration** for batch launches with configurable layouts
- **Context-aware spawning** (inside tmux vs standalone)
- **Project-based working directories** for proper context
- **Saved profiles** for complex multi-pane setups
- **Responsive layouts** adapting from desktop to mobile (Termux)

### Why Not Just mcfly?

**mcfly:** Great for "I ran this before, what was it?" (command history search)
**tui-launcher:** "I want to start X in Y context with Z layout" (workspace orchestrator)

They complement each other - mcfly for ad-hoc commands, launcher for organized workflows.

---

## Current Status

**Version:** 0.3.0-dev (Phase 1: Tmuxplexer Integration)
**Last Updated:** 2025-11-20 (Evening)
**Branch:** feature/tmuxplexer-integration

### ✅ Completed (v0.2.0 - 3-Pane Layout)
See [CHANGELOG.md](CHANGELOG.md) for full details:
- ✅ 3-pane responsive layout (Desktop/Compact/Mobile modes)
- ✅ Global Tools | Projects | Info pane structure
- ✅ Tab navigation between panes
- ✅ Info pane with item details
- ✅ Responsive breakpoints for Termux
- ✅ Core tree view navigation with multi-select
- ✅ YAML configuration system
- ✅ Profile support for multi-pane setups
- ✅ Keyboard/mouse navigation
- ✅ Wrapper script (`tl`) for global access

### ✅ Completed (v0.3.0 - Launch Tab Refinement)
- ✅ Simplified spawn logic (TFE-inspired foreground mode)
- ✅ Foreground/Detached mode toggle ('d' key)
- ✅ Multi-select tmux spawning (foreground and detached)
- ✅ Fixed 'e' key (edit config) using `tea.ExecProcess`
- ✅ Mode indicator in header
- ✅ Debug logging for troubleshooting

### 🚧 In Progress (v0.3.0 - Tmuxplexer Integration)

**Goal:** Integrate tmuxplexer's session management and template features into tui-launcher as tabs.

#### Phase 1: Tab Architecture & Launch Tab ✅ COMPLETE

**Completed (2025-11-20):**
- ✅ Created unified tab-based architecture
- ✅ Implemented tab routing (1/2/3 keys, Tab/Shift+Tab cycling)
- ✅ Created `tabs/launch/` package with full Launch tab functionality
- ✅ Created `shared/` layer merging spawn.go + tmuxplexer tmux operations
- ✅ Migrated all tui-launcher features to Launch tab
- ✅ Build succeeds, compiles cleanly
- ✅ Fixed config loading bug (message routing)
- ✅ Fixed tab bar display bug (now shows on initial launch)

**Ready for Testing:**
- 🔍 Real terminal testing needed (config loading, spawning, Quick CD)
- 🔍 Verify tab bar displays correctly on launch
- 🔍 Test all keyboard shortcuts and spawn modes

**Architecture:**
```
tabs/launch/          # Launch tab (existing tui-launcher features)
  ├── model.go        # Model, Init(), config loading
  ├── view.go         # Multi-pane rendering
  ├── update.go       # Keyboard/mouse handling
  └── tree.go         # Tree building from config

shared/               # Unified tmux operations
  ├── tmux.go         # Spawn + session management
  └── types.go        # Shared type definitions

model_unified.go      # Tab routing coordinator
tab_routing.go        # Message routing to active tab
types_unified.go      # Tab types (tabName, unifiedModel)
```

**Tab System:**
```
┌─ 1. Launch ──┬─ 2. Sessions ──┬─ 3. Templates ──┐
│ [Active Tab Content Below]                      │
└──────────────────────────────────────────────────┘
```

- **Tab 1 (Launch):** ✅ Full tui-launcher functionality
- **Tab 2 (Sessions):** 🔜 Tmux session management (from tmuxplexer)
- **Tab 3 (Templates):** 🔜 Workspace templates (from tmuxplexer)

**Files:** See [PHASE1_LAUNCH_TAB_COMPLETE.md](PHASE1_LAUNCH_TAB_COMPLETE.md) for detailed implementation

#### Next Steps (Phase 1 Continuation)

**Priority 1: Real Terminal Testing**
- [ ] Verify config loads from ~/.config/tui-launcher/config.yaml
- [ ] Test navigation (arrows, vim keys, mouse)
- [ ] Test multi-select and spawning
- [ ] Test Quick CD functionality
- [ ] Test config editing (e key)
- [ ] Verify all spawn modes work

**Priority 2: Sessions Tab (Partial)**
- [ ] Create `tabs/sessions/` package
- [ ] Copy session model from tmuxplexer
- [ ] Implement basic sessions list view
- [ ] Wire into tab routing
- [ ] **Checkpoint:** Press 2 → See tmux sessions list

**Priority 3: Templates Tab (Partial)**
- [ ] Create `tabs/templates/` package
- [ ] Copy templates model from tmuxplexer
- [ ] Implement template tree view
- [ ] Wire into tab routing
- [ ] **Checkpoint:** Press 3 → See template list

---

## Roadmap

### v0.2.0 - Responsive 3-Pane Layout ✅ COMPLETE
- ✅ 3-pane layout for desktop (Global | Projects | Info)
- ✅ Responsive breakpoints for Termux compatibility
- ✅ Info pane with item details
- ✅ Tab navigation between panes (Tab key)
- ✅ 'i' key to toggle info in mobile mode
- ✅ Config schema supports project paths and profiles

### v0.3.0 - Tmuxplexer Integration (CURRENT)

**Phase 1: Tab Architecture & Launch Tab** ✅ COMPLETE
- ✅ Unified tab-based architecture (1/2/3 keys)
- ✅ Launch tab with all tui-launcher features
- ✅ Shared tmux operations layer
- ✅ Config loading and tree building
- ✅ Tab bar displays on initial launch
- 🔍 Real terminal testing (config, spawning, Quick CD)

**Phase 2: Sessions Tab** 🔜 NEXT
- [ ] Tmux sessions list (from tmuxplexer)
- [ ] Session management (attach, kill, rename)
- [ ] Live session preview
- [ ] Claude Code status tracking
- [ ] Window navigation
- [ ] Auto-refresh (2-second interval)

**Phase 3: Templates Tab** 🔜 UPCOMING
- [ ] Workspace templates tree view
- [ ] Categorized templates (Projects, Agents, Tools)
- [ ] Template creation wizard
- [ ] Save session as template
- [ ] Template preview and editing
- [ ] Grid layout support (2x2, 3x3, etc.)

**Phase 4: Unified Config** 🔜 FUTURE
- [ ] Merge launcher.yaml + templates.json → single config
- [ ] Migration script for existing configs
- [ ] Unified template and command definitions

**Phase 5: Cross-Tab Features** 🔜 FUTURE
- [ ] Launch → Sessions (auto-switch after spawn)
- [ ] Sessions → Templates (save as template)
- [ ] Templates → Launch (show in tree)
- [ ] Popup mode integration (Ctrl+B O)

### v0.4.0 - Documentation & Polish
- [ ] Update README with tab-based interface
- [ ] Add keybindings reference (all tabs)
- [ ] Screenshot/demo GIF
- [ ] Config examples for templates
- [ ] Installation script (`install.sh`)

### v0.5.0 - Enhanced Features
- [ ] Favorites system (star items)
- [ ] Recent launches (history)
- [ ] Search/filter (Ctrl+F or /)
- [ ] Command-line args (`--project`, `--tool`, `--template`)
- [ ] Error handling improvements

---

## Future Ideas (Post v1.0)

### Advanced Features
- **Remote spawning** - SSH into servers and spawn there
- **Docker integration** - Launch containers with tmux inside
- **Command templates** - Variables in commands ({{project}}, {{branch}})
- **Conditional commands** - Only show if file/dir exists
- **Status indicators** - Show if process is running (🟢/🔴)
- **Web UI** - Launch via browser (for remote access)

### AI Integration
- **AI config generation** - Ask Claude to generate launch configs
- **Smart suggestions** - Recommend commands based on context
- **Natural language** - "Launch dev environment for TKan"

### Collaboration
- **Export/import** - Share configs with team
- **Team profiles** - Shared workspace configurations
- **Config sync** - Sync across machines

---

## Architecture Overview

### File Structure
```
tui-launcher/
├── main.go               # Entry point
├── types.go              # Type definitions
├── model.go              # Model initialization
├── tree.go               # Tree building/rendering
├── spawn.go              # Spawn logic (tmux/xterm)
├── layouts.go            # Tmux layout definitions
├── go.mod
├── PLAN.md              # This file
├── CHANGELOG.md         # Version history
├── IMPLEMENTATION_PLAN.md  # 3-pane layout details
├── CLAUDE.md            # Instructions for Claude Code
└── README.md
```

### Core Principles
- **MVU Pattern**: Model-View-Update (Bubble Tea)
- **Responsive Design**: Adapt to terminal size (Golden Rules)
- **Separation of Concerns**: Global tools vs project commands
- **Proven Patterns**: Reuse from TFE (tree view) and tmuxplexer (spawn)
- **Mobile First**: Optimize for Termux alongside desktop

---

## Configuration System

YAML-based config at `~/.config/tui-launcher/config.yaml`

### Current Schema (v0.1.0)
```yaml
projects:
  - name: Project Name
    icon: 🚀
    path: ~/projects/project-name
    commands:
      - name: Command Name
        icon: 📂
        command: command to run
        cwd: ~/specific/directory  # Optional
        spawn: tmux-split-h        # Optional
    profiles:
      - name: Profile Name
        icon: 🔧
        layout: main-vertical
        panes:
          - command: command1
          - command: command2

tools:
  - category: Category Name
    icon: 🔧
    items:
      - name: Tool Name
        icon: 🎯
        command: tool-command
        spawn: tmux-window
```

### Planned Schema (v0.2.0)
Add info/documentation support:
```yaml
tools:
  - category: Git
    icon: 🔧
    items:
      - name: lazygit
        icon: 🎯
        command: lazygit
        description: "Terminal UI for git commands"
        info_file: ~/.config/tui-launcher/docs/lazygit.md
        repo: "https://github.com/jesseduffield/lazygit"
```

---

## Development Workflow

### Building & Testing
```bash
# Build
go build

# Run locally
./tui-launcher

# Install globally
go build -o tui-launcher && cp tui-launcher ~/.local/bin/

# Update dependencies
go mod tidy
```

### Using with Wrapper
The `tl` wrapper at `~/.local/bin/tl` calls the binary. Always rebuild and copy after changes:
```bash
go build -o tui-launcher && cp tui-launcher ~/.local/bin/
```

### Skills & Documentation
- `.claude/skills/bubbletea/` - TUI development patterns
- `IMPLEMENTATION_PLAN.md` - Detailed implementation steps
- `CHANGELOG.md` - Version history

---

## Success Metrics

**Must Have (v0.2.0):**
- ✅ 3-pane layout works on desktop
- ✅ Responsive layout works in Termux
- ✅ Info pane shows helpful content
- ✅ No visual glitches (borders align, no overflow)
- ✅ All existing functionality preserved

**Should Have (v0.3.0):**
- Documentation is clear and complete
- Installation is easy (one script)
- Examples cover common use cases

**Nice to Have (v0.4.0+):**
- Search/filter for quick access
- Favorites for most-used commands
- Command history tracking
- TFE integration seamless

---

**Status:** Active Development
**Contributors:** matt, Claude Code
**License:** MIT
