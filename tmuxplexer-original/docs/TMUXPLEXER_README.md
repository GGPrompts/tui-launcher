# Tmuxplexer Exploration - Complete Documentation Index

This directory contains comprehensive documentation about the Tmuxplexer project for TFE integration.

## Documents

### 1. TMUXPLEXER_SUMMARY.md (Quick Overview)
**Purpose:** Fast reference guide  
**Length:** 9.4 KB  
**Best for:** Getting a quick understanding of what tmuxplexer is and does

**Contains:**
- Quick facts (language, size, status)
- Feature overview
- Architecture diagram
- File breakdown table
- Configuration examples
- Integration strategy for TFE
- Key takeaways

**Start here if:** You have 10 minutes and want the big picture

---

### 2. TMUXPLEXER_ANALYSIS.md (Detailed Technical Analysis)
**Purpose:** Comprehensive technical reference  
**Length:** 30 KB (891 lines)  
**Best for:** Deep understanding, implementation decisions, integration planning

**Contains:**
- Complete project overview (8 phases)
- Full directory structure with file purposes
- Core functionality breakdown
- 4-panel layout system details
- Session state management
- Keyboard & mouse event handling
- Templates system (format, operations, examples)
- Claude Code integration architecture
- Configuration formats (YAML, JSON, CLI flags)
- Dependency analysis
- TFE integration strategies (4 options)
- Key code patterns & examples
- Development workflow
- File size metrics
- Resources & documentation

**Start here if:** You need to implement integration or understand details

---

## Quick Navigation

### For TFE Integration
1. Read: **SUMMARY.md** → Section "TFE Integration Ready"
2. Read: **ANALYSIS.md** → Section "8. TFE INTEGRATION POINTS"
3. Implement: Use the example code in SUMMARY.md

### For Understanding Architecture
1. Read: **SUMMARY.md** → "Architecture Overview"
2. Read: **ANALYSIS.md** → Section "3. CORE FUNCTIONALITY"
3. Read: **ANALYSIS.md** → Section "9. KEY CODE PATTERNS"

### For Understanding Templates
1. Read: **SUMMARY.md** → "Configuration Files"
2. Read: **ANALYSIS.md** → Section "4. TEMPLATES SYSTEM"
3. Look at: `~/.config/tmuxplexer/templates.json` (actual examples)

### For Understanding Claude Integration
1. Read: **SUMMARY.md** → Section "Key Features Implemented" (4. Claude Code Integration)
2. Read: **ANALYSIS.md** → Section "5. CLAUDE CODE INTEGRATION"
3. Check: `/home/matt/projects/tmuxplexer/hooks/README.md`

### For Development
1. Read: **ANALYSIS.md** → Section "10. DEVELOPMENT WORKFLOW"
2. Read: **ANALYSIS.md** → Section "13. QUICK REFERENCE FOR DEVELOPERS"
3. Check: `/home/matt/projects/tmuxplexer/CLAUDE.md` (in the actual repo)

---

## Key Facts

- **Language:** Go 1.24.0
- **Framework:** Bubble Tea (TUI)
- **Total Code:** 4,749 lines
- **Main Files:** 14 Go files
- **Status:** Production-ready
- **Maturity:** All 8 phases complete

---

## Integration Status

### Current State
- ✅ `--cwd` flag implemented
- ✅ `--template` flag implemented
- ✅ Ready for TFE integration

### Required for TFE
- Add context menu item "Launch Dev Workspace"
- Pass current directory as `--cwd`
- Pass template index as `--template`

### Example
```bash
tmuxplexer --cwd ~/projects/myapp --template 0
```

---

## File Locations

- **Tmuxplexer repo:** `/home/matt/projects/tmuxplexer`
- **This analysis:** `/home/matt/projects/TFE/docs/`
  - `TMUXPLEXER_SUMMARY.md` - Quick overview
  - `TMUXPLEXER_ANALYSIS.md` - Detailed analysis
  - `TMUXPLEXER_README.md` - This file

---

## Quick Reference

### Directory Structure
```
tmuxplexer/
├── main.go                    # Entry point
├── types.go                   # Type definitions
├── model.go                   # State & layout
├── update.go                  # Message dispatcher
├── update_keyboard.go         # Keyboard handling (982 lines)
├── update_mouse.go            # Mouse handling
├── view.go                    # View rendering
├── styles.go                  # Styling
├── config.go                  # Configuration (YAML)
├── templates.go               # Templates (JSON)
├── tmux.go                    # Tmux integration (633 lines)
├── claude_state.go            # Claude integration
├── hooks/                     # Claude hooks system
├── components/                # UI components
└── lib/                       # Utility libraries
```

### Configuration Files
- **Tmuxplexer config:** `~/.config/tmuxplexer/config.yaml` (YAML)
- **Templates:** `~/.config/tmuxplexer/templates.json` (JSON)
- **Claude hooks:** `~/.claude/hooks/` (Bash scripts)
- **Claude state:** `/tmp/claude-code-state/` (JSON files)

### Key Commands
```bash
# Build
go build -o tmuxplexer

# Run TUI
./tmuxplexer

# Popup mode (from tmux)
./tmuxplexer --popup

# Create session from template
./tmuxplexer --template 0
./tmuxplexer --cwd /path --template 1

# View templates (no TTY)
./tmuxplexer test_template

# Create session (no TTY)
./tmuxplexer test_create 0
```

---

## Dependencies

### Go Modules
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - UI components
- `github.com/charmbracelet/lipgloss` - Styling
- `gopkg.in/yaml.v3` - YAML parsing

### System Tools
- `tmux` - Required
- `$EDITOR` - Optional (template editing)
- `git` - Optional (branch detection)

---

## Features at a Glance

### 4-Panel Accordion Layout
- Header (stats & quick actions)
- Left (sessions list + details)
- Right (templates list + details)
- Footer (live pane preview with scrollback)
- Focus switching: keys 1, 2, 3, 4
- Accordion mode: `a` key toggles expansion

### Session Management
- List all tmux sessions
- Attach/kill sessions
- View live pane content
- Full scrollback history (PgUp/PgDn)
- Window navigation (arrow keys)
- Auto-refresh every 2 seconds
- Working directory and git branch display

### Templates
- Store multi-pane layouts as JSON
- Create via wizard (`n` key)
- Save running session as template (`s` key)
- Edit in $EDITOR (`e` key)
- Delete with confirmation (`d` key)
- Support for per-pane working directories

### Claude Code Integration
- Real-time status tracking
- Hook system at `~/.claude/hooks/`
- Status indicators: Idle, Processing, Tool Use, Working, Awaiting Input, Stale
- Auto-scroll preview to bottom for Claude sessions
- Orange text styling for Claude sessions

### Popup Mode
- Launch from tmux with `Ctrl+b o`
- 80% width/height floating window
- Session switching without leaving tmux
- Keybinding: `bind-key o run-shell "tmux popup -E -w 80% -h 80% tmuxplexer --popup"`

---

## For Your Exploration

### What to Read
1. Start with **SUMMARY.md** (10 minutes)
2. Deep dive with **ANALYSIS.md** sections relevant to your needs
3. Check the actual code in `/home/matt/projects/tmuxplexer/`

### What to Try
```bash
# Build it
cd ~/projects/tmuxplexer
go build

# See templates
./tmuxplexer test_template

# Create a session
./tmuxplexer --template 0

# Or run in TUI
./tmuxplexer
```

### What to Understand for TFE Integration
- CLI flag handling (section 6.3 in ANALYSIS)
- TFE integration strategies (section 8.4 in ANALYSIS)
- Template system (section 4 in ANALYSIS)
- Message flow pattern (section 9.1 in ANALYSIS)

---

## Document Purposes

| Document | Purpose | Audience | Reading Time |
|----------|---------|----------|--------------|
| **SUMMARY.md** | Quick reference | Everyone | 10 min |
| **ANALYSIS.md** | Technical deep-dive | Developers | 30 min |
| **CLAUDE.md** (in repo) | Development guide | Contributors | 20 min |
| **README.md** (in repo) | User guide | End users | 10 min |
| **PLAN.md** (in repo) | Project roadmap | Project leads | 45 min |

---

## Contact & Resources

### Tmuxplexer Repo
Location: `/home/matt/projects/tmuxplexer`  
Documentation: `CLAUDE.md`, `README.md`, `PLAN.md`, `docs/`

### External Resources
- Bubble Tea: https://github.com/charmbracelet/bubbletea
- Lipgloss: https://github.com/charmbracelet/lipgloss
- Tmux Manual: https://man.openbsd.org/tmux

### Related Projects
- **TFE** (Terminal File Explorer) - Integration target
- **TUITemplate** - Reusable TUI components and patterns

---

## Summary

Tmuxplexer is a production-ready terminal UI for managing tmux sessions with workspace templates. It's well-architected, fully featured, and ready for integration with TFE. The `--cwd` and `--template` CLI flags are already implemented, making integration straightforward.

**For TFE integration:** Add a context menu item that calls `tmuxplexer --cwd $PWD --template <index>` when users want to launch a workspace.

**Next steps:** Review SUMMARY.md, then check ANALYSIS.md section 8 for integration strategies.

