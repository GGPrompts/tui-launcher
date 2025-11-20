# Tmuxplexer Exploration - Complete Analysis

## Quick Facts

**Project:** Terminal User Interface tmux session manager  
**Language:** Go 1.24.0  
**Size:** 4,749 lines of code across 14 files  
**Status:** Production-ready (8 phases complete)  
**Architecture:** Bubble Tea (Model-View-Update pattern)  
**Configuration:** YAML (config) + JSON (templates)  

---

## What Does It Do?

Tmuxplexer is a modern terminal UI for managing tmux sessions with:

1. **4-Panel Accordion Layout** - Header (stats), Sessions (left), Templates (right), Preview (footer)
2. **Workspace Templates** - Save multi-pane layouts to JSON, recreate with one keystroke
3. **Session Management** - Attach, kill, view live pane content with full scrollback
4. **Claude Code Integration** - Real-time status tracking (Idle, Processing, Tool Use, etc.)
5. **Popup Mode** - Launch from tmux with `Ctrl+b o` keybinding
6. **CLI Integration** - Flags for context-aware workspace creation

---

## Architecture Overview

```
┌─────────────────────────────────────────┐
│         Bubble Tea Framework             │
│  (Model-View-Update pattern)             │
├─────────────────────────────────────────┤
│                                          │
│  Model (state)         → types.go        │
│  ├── Sessions          → model.go        │
│  ├── Templates                           │
│  ├── Panels            → templates.go    │
│  └── Input modes                         │
│                                          │
│  Update (events)       → update*.go      │
│  ├── Keyboard          → update_keyboard.go (982 lines)
│  ├── Mouse            → update_mouse.go  │
│  └── Messages         → update.go        │
│                                          │
│  View (rendering)      → view.go         │
│  ├── 4-panel layout    → styles.go       │
│  ├── Dynamic sizing                      │
│  └── Content formatting                  │
│                                          │
│  Integration           → tmux.go         │
│  ├── Session ops       → claude_state.go │
│  ├── Templates         → config.go       │
│  └── Claude hooks      → hooks/*         │
└─────────────────────────────────────────┘
```

---

## File Breakdown

| File | Lines | Purpose |
|------|-------|---------|
| **update_keyboard.go** | 982 | Keyboard event handling (largest) |
| **tmux.go** | 633 | Tmux integration & session ops |
| **model.go** | 601 | State management & layout |
| **view.go** | 513 | View rendering & panels |
| **templates.go** | 272 | Template loading/saving |
| **types.go** | 299 | Type definitions & structs |
| **config.go** | 237 | Configuration (YAML) |
| **styles.go** | 238 | Lipgloss styling |
| **update.go** | 238 | Message dispatcher |
| **claude_state.go** | 212 | Claude integration |
| **update_mouse.go** | 185 | Mouse event handling |
| **main.go** | 165 | Entry point + TTY check |
| **Other** | 174 | Test commands, misc |
| **Total** | 4,749 | All Go code |

---

## Key Features Implemented

### 1. 4-Panel Layout System
- Weight-based dynamic sizing
- Accordion mode: focused panel expands
- Focus switching: keys 1, 2, 3, 4
- Mouse click to focus panels

### 2. Session Management
- List all tmux sessions
- Attach/kill sessions
- View working directory and git branch
- Live pane preview with full scrollback (PgUp/PgDn)
- Window navigation (arrow keys)
- Auto-refresh every 2 seconds

### 3. Template System
- Load from `~/.config/tmuxplexer/templates.json`
- Create via wizard (`n` key)
- Save running session as template (`s` key)
- Edit manually (`e` key)
- Delete templates (`d` key)
- Per-pane working directory override

### 4. Claude Code Integration
- Hook system at `~/.claude/hooks/`
- State files at `/tmp/claude-code-state/`
- Real-time status: Idle, Processing, Tool Use, Working, Awaiting Input, Stale
- Auto-scroll preview to bottom for Claude sessions
- Orange text styling for Claude sessions

### 5. Popup Mode
- Launch from tmux with `Ctrl+b o`
- 80% width/height floating window
- Session switching without leaving tmux
- Keybinding in `~/.tmux.conf`

### 6. CLI Integration
- `--popup` - Popup mode
- `--template <index>` - Create from template
- `--cwd <directory>` - Override working directory
- `test_template` - List templates (no TTY)
- `test_create <index>` - Create session (no TTY)

---

## Configuration Files

### Templates (`~/.config/tmuxplexer/templates.json`)

```json
[
  {
    "name": "Simple Dev (2x2)",
    "description": "Basic dev workspace",
    "working_dir": "~",
    "layout": "2x2",
    "panes": [
      {"command": "nvim", "title": "Editor"},
      {"command": "bash", "title": "Terminal"},
      {"command": "lazygit", "title": "Git"},
      {"command": "btop", "title": "Monitor"}
    ]
  }
]
```

### Config (`~/.config/tmuxplexer/config.yaml`)

```yaml
theme: "dark"
ui:
  mouse_enabled: true
  show_icons: true
layout:
  type: "single"
  split_ratio: 0.5
performance:
  lazy_loading: true
  cache_size: 100
```

---

## TFE Integration Ready

### Already Implemented
- ✅ `--cwd` flag for directory override
- ✅ `--template` flag for template selection
- ✅ Combined: `tmuxplexer --cwd $PWD --template 1`

### Recommended Integration
In TFE's `context_menu.go`:
```go
case "Launch Dev Workspace":
    exec.Command("tmuxplexer",
        "--cwd", m.currentPath,
        "--template", "0")
```

### Example Workflow
1. Browse to `/home/matt/projects/myapp` in TFE
2. Right-click → "Launch Dev Workspace"
3. Tmuxplexer creates session in `myapp` directory
4. Session has all configured panes and commands running

---

## Dependencies

### Go Modules
```
github.com/charmbracelet/bubbletea v1.1.0  # TUI framework
github.com/charmbracelet/bubbles v0.20.0   # UI components
github.com/charmbracelet/lipgloss v0.13.1  # Styling
gopkg.in/yaml.v3 v3.0.1                    # YAML parsing
```

### System Requirements
- `tmux` - For session management
- `$EDITOR` - For template editing (optional)
- `git` - For branch detection (optional)

### No External UI Libraries
- Pure Bubble Tea + Lipgloss
- No heavy GUI frameworks
- ~5.4MB compiled binary

---

## Code Patterns

### Message Flow (Elm Architecture)
```
User Input → Keyboard/Mouse Event → Handler → Message
           → Update() routes message → State change
           → View() re-renders → Display update
```

### Panel Update Pattern
```go
// When sessions refresh:
listSessions() → sessionsLoadedMsg
              → Update() handles
              → updateLeftPanelContent()
              → m.View() re-renders
```

### Auto-Refresh
- Ticker every 2 seconds
- Refreshes sessions and Claude state
- Updates panel content if changed
- No polling overhead when idle

---

## Documentation

### In-Repository
- **CLAUDE.md** (607 lines) - Development guide & architecture
- **README.md** (228 lines) - User guide & quick start
- **PLAN.md** (50+ pages) - Detailed roadmap & ideas
- **docs/** - Claude hooks integration, research

### Code Examples
- `test_template.go` - List templates without TTY
- `test_create_session.go` - Create session without TTY
- Built-in templates - TFE Development (4x2 grid)

---

## Current Maturity

### Production Ready
✅ All 8 development phases complete:
1. 4-Panel Accordion Layout
2. Workspace Templates
3. Session Management
4. Live Pane Preview & Window Navigation
5. Claude Code Integration
6. Scrollable Preview with Scrollback
7. Template Creation Wizard & Deletion
8. Popup Mode with Tmux Keybinding

### Performance
- Session refresh: ~500ms
- Template creation: <100ms
- Claude state detection: ~50ms per session
- UI render: <16ms (60 FPS)
- Memory: 15-20MB for 20+ sessions

### Known Limitations
- No command execution within panes
- No session renaming UI
- No theme customization UI (YAML only)
- State files auto-cleanup after 24h

---

## Next Steps for TFE Integration

### Phase 1 (Quick Win)
- [x] Understand tmuxplexer architecture
- [x] Review integration points
- [ ] Add context menu "Launch Workspace"
- [ ] Test with existing templates

### Phase 2 (Enhanced)
- [ ] Template selection dialog in TFE
- [ ] Per-project template config
- [ ] Auto-detect project type

### Phase 3 (Future)
- [ ] Template variables (${PROJECT}, ${HOME})
- [ ] Per-pane command variables
- [ ] Logging of workspace launches

---

## Key Takeaways

1. **Well-Architected** - Clean separation of concerns, 14 focused files
2. **Integration-Ready** - CLI flags already in place for TFE
3. **Production-Ready** - All phases complete, well-tested
4. **Claude-Aware** - Real-time status tracking via hooks
5. **Extensible** - Clear patterns for adding features
6. **User-Friendly** - Both keyboard and mouse support
7. **Lightweight** - No heavy dependencies, ~5MB binary

---

## Files

Main analysis document saved to:  
`/home/matt/projects/TFE/docs/TMUXPLEXER_ANALYSIS.md` (891 lines)

Covers:
- Project overview
- Directory structure
- Core functionality
- Templates system
- Claude integration
- Configuration formats
- Dependencies
- TFE integration strategies
- Code patterns
- Development workflow
- Quick references

