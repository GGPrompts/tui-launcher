# TUI Launcher - Development Plan

**A visual terminal launcher for managing projects, TUI tools, and batch command execution with tmux integration**

---

## 🚨 Priority Issues for Next Session

### Critical: Command & Hotkey Audit
**Problem:** Several core features appear broken or provide no feedback:
- **'e' key (edit config)** - Not working reliably
- **Launch commands** - Many spawn options don't seem to work
- **No visual feedback** - Commands launch but TUI just disappears (confusing UX)

**Action Items:**
1. **Audit all keyboard shortcuts** - Test each key binding:
   - [ ] 'e' - Edit config (inside/outside tmux)
   - [ ] Enter - Launch commands (all spawn modes)
   - [ ] Space - Selection/expansion
   - [ ] Tab - Pane switching
   - [ ] 'i' - Info toggle (mobile mode)
   - [ ] 't' - Toggle tmux mode
   - [ ] 'c' - Clear selections

2. **Test all spawn modes:**
   - [ ] `tmux-window` - New tmux window
   - [ ] `tmux-split-h` - Horizontal split
   - [ ] `tmux-split-v` - Vertical split
   - [ ] `tmux-layout` - Custom layout
   - [ ] `current-pane` - Replace current pane
   - [ ] `xterm-window` - New xterm window
   - [ ] Direct mode (non-tmux)

3. **Add visual feedback for launches:**
   - [ ] Show "Launching..." message before quit
   - [ ] Display which commands are being spawned
   - [ ] Show spawn mode being used
   - [ ] Add delay or confirmation before exit
   - [ ] Consider toast/notification for successful launches
   - [ ] Error messages if spawn fails

4. **Improve error handling:**
   - [ ] Catch spawn errors and display them
   - [ ] Don't quit if launch fails
   - [ ] Show helpful error messages
   - [ ] Validate commands before spawning

**Expected Outcome:**
- All hotkeys work reliably
- Clear visual feedback when commands launch
- Users understand what's happening (not just "it disappeared")
- Failed launches show helpful error messages

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

**Version:** 0.2.0-dev (in progress)
**Last Updated:** 2025-01-19

### ✅ Completed (v0.1.0 - MVP)
See [CHANGELOG.md](CHANGELOG.md) for full details of completed features:
- Core tree view navigation with multi-select
- Tmux spawn logic with multiple modes
- YAML configuration system
- Profile support for multi-pane setups
- Keyboard/mouse navigation
- Wrapper script (`tl`) for global access

### 🚧 In Progress (v0.2.0)

#### 3-Pane Responsive Layout System
Implementing a responsive layout that adapts to terminal size:

**Desktop Mode** (≥80 width, >12 height):
```
┌──────────────────┬──────────────────┐
│ Global Tools     │ Projects         │
│ ├─ Git           │ ├─ TUI Launcher  │
│ └─ AI            │ └─ TKan          │
├──────────────────┴──────────────────┤
│ Info: lazygit                       │
│ Terminal UI for git commands        │
└─────────────────────────────────────┘
```

**Compact Mode** (<80 width) - Termux landscape:
```
┌────────────────────────────────────┐
│ Global Tools / Projects (Tab)      │
├────────────────────────────────────┤
│ Info pane                          │
└────────────────────────────────────┘
```

**Mobile Mode** (≤12 height) - Termux with keyboard:
```
┌────────────────────────────────────┐
│ Tree only (press 'i' for info)    │
└────────────────────────────────────┘
```

**Status:** Implementation in progress by tt-cc-ofc session
**Files:** See [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) for detailed steps

---

## Roadmap

### v0.2.0 - Responsive 3-Pane Layout (Current)
- [ ] 3-pane layout for desktop (Global | Projects | Info)
- [ ] Responsive breakpoints for Termux compatibility
- [ ] Info pane with markdown file support
- [ ] Tab navigation between panes
- [ ] 'i' key to toggle info in mobile mode
- [ ] Update config schema for item descriptions and info files

### v0.3.0 - Documentation & Polish
- [ ] Update README with installation instructions
- [ ] Add keybindings reference
- [ ] Screenshot/demo GIF
- [ ] Config examples for common workflows
- [ ] TFE integration example
- [ ] Installation script (`install.sh`)

### v0.4.0 - Enhanced Features
- [ ] Favorites system (star items)
- [ ] Recent launches (history)
- [ ] Search/filter (Ctrl+F or /)
- [ ] Command-line args (`--project`, `--tool`)
- [ ] Session management (list, kill, switch)
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
