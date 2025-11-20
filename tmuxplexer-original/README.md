# Tmuxplexer

A modern TUI tmux session manager with workspace templates. Create complex tmux layouts instantly!

## Features

✨ **Workspace Templates** - Save and launch multi-pane layouts with a single keypress
🎯 **Dual-Panel Layout** - Sessions (left) and Templates (right) with independent navigation
⚡ **Instant Session Creation** - Templates → Live tmux sessions in seconds
🎨 **Beautiful TUI** - Bubble Tea framework with lipgloss styling
📺 **Scrollable Pane Preview** - Full scrollback history with PgUp/PgDn navigation
🤖 **Claude Code Integration** - Real-time status tracking for Claude sessions
🔄 **Auto-Refresh** - Session list, Claude status, and preview update every 2 seconds
⌨️ **Full Keyboard Control** - Navigate sessions, windows, and panels
🎛️ **Session Management** - Attach, kill, and monitor sessions easily

## Quick Start

### ⚠️ Run in a Real Terminal

Tmuxplexer is a **TUI application** and requires a proper terminal:

```bash
# In Windows Terminal, iTerm2, or any terminal emulator:
cd ~/projects/tmuxplexer
./tmuxplexer
```

**Not working?** Make sure you're in a real terminal, not a background process or IDE terminal.

### Navigation

**Panel Selection:**
- **1** - Focus left panel (Sessions)
- **2** - Focus right panel (Templates)
- **3** - Focus footer panel (Preview)
- **4** - Focus header panel (Stats)

**Within Panels:**
- **↑/↓** or **k/j** - Navigate items in focused panel
- **←/→** or **h/l** - Navigate windows in preview panel
- **PgUp/PgDn** - Scroll preview content (footer panel only)
- **Home/End** or **g/G** - Jump to top/bottom of preview (footer panel only)
- **Enter** - Attach to session (left panel) OR create from template (right panel)
- **d** or **K** - Kill selected session (left panel only)
- **e** - Edit templates in your editor (right panel only)

**Global:**
- **a** - Toggle accordion mode
- **Ctrl+R** - Refresh session list and Claude status
- **q** - Quit

### Test Commands (No TTY Required)

```bash
# View available templates
./tmuxplexer test_template

# Create a session from template 0 (Simple Dev 2x2)
./tmuxplexer test_create 0

# Create from template 1 (Frontend Dev 2x2)
./tmuxplexer test_create 1

# List created sessions
tmux ls
```

## Workspace Templates

Templates are stored in `~/.config/tmuxplexer/templates.json`

### Example Template

```json
{
  "name": "Frontend Dev (2x2)",
  "description": "Frontend workspace with Claude, editor, dev server, and git",
  "working_dir": "~/projects/my-app",
  "layout": "2x2",
  "panes": [
    {"command": "claude-code .", "title": "Claude AI"},
    {"command": "nvim", "title": "Editor"},
    {"command": "npm run dev", "title": "Dev Server"},
    {"command": "lazygit", "title": "Git"}
  ]
}
```

### Supported Layouts

- **2x2** - 4 panes in a 2×2 grid
- **4x2** - 8 panes in a 4×2 grid
- **3x3** - 9 panes in a 3×3 grid
- Any **COLSxROWS** format

### Built-in Templates

1. **Simple Dev (2x2)** - Editor, terminal, git, monitor
2. **Frontend Dev (2x2)** - Claude, editor, dev server, git
3. **TFE Development (4x2)** - Full 8-pane workspace
4. **Monitoring Wall (4x2)** - System monitoring dashboard

## How It Works

1. **Load Templates** - Reads `~/.config/tmuxplexer/templates.json`
2. **Dual-Panel UI** - Left panel shows sessions, right panel shows templates
3. **Claude Integration** - Detects Claude Code sessions and shows real-time status
4. **Select & Create** - Press Enter on template to create session
5. **Create Grid Layout** - Splits tmux window into NxM panes
6. **Run Commands** - Sends configured commands to each pane
7. **Manage Sessions** - View details, attach, kill sessions in left panel
8. **Live Preview** - See real-time pane content in footer panel
9. **Auto-Refresh** - Sessions, Claude status, and preview update every 2 seconds

## Troubleshooting

### "No TTY detected" Error

Run in a real terminal application (Windows Terminal, iTerm2, etc.), not in background scripts or IDE terminals.

### Session Creation Fails

1. Make sure tmux is installed: `tmux -V`
2. Check template syntax in `~/.config/tmuxplexer/templates.json`
3. Verify working directories exist

## Development Status

### ✅ Phase 1: Layout System (Complete)
- 4-panel accordion layout with perfect alignment
- Dynamic panel sizing with weight-based calculations
- Panel focus switching (keys 1,2,3,4)
- Accordion mode toggle ('a')

### ✅ Phase 2: Workspace Templates (Complete)
- Template loading from ~/.config/tmuxplexer/templates.json
- Session creation from templates
- Grid layouts (2x2, 4x2, 3x3, any COLSxROWS)
- 4 built-in templates

### ✅ Phase 3: Session Management (Complete)
- Attach to sessions (Enter) - detects inside/outside tmux
- Kill sessions (d/K)
- Session details in right panel
- Auto-refresh every 2 seconds
- Window list with active indicators

### ✅ Phase 4: Live Preview & Window Navigation (Complete)
- Live pane content preview in footer panel
- Window navigation with ←/→ arrows
- Auto-refresh of preview content
- Window position indicator (1/3, 2/3, etc.)

### ✅ Phase 5: Claude Code Integration & Layout Reorganization (Complete)
- Claude Code hooks integration for real-time status tracking
- Status indicators: 🟢 Idle, 🟡 Processing, 🔧 Tool Use, ⚙️ Working, ⏸️ Awaiting Input, ⚪ Stale
- Dual-panel layout: Sessions (left) + Templates (right)
- Independent navigation in each panel
- Session details inline in left panel
- Template details inline in right panel
- Auto-detection of Claude sessions via working directory

### ✅ Phase 6: Scrollable Preview (Complete)
- Full scrollback history capture from tmux panes
- PgUp/PgDn scrolling through preview content
- Home/End (g/G) to jump to top/bottom
- Scroll position indicator showing current line and percentage
- Auto-reset scroll position when changing sessions or windows

### 🚀 Future Enhancements (Phase 7+)

Potential features for future development:

- **Send Commands** - Type and send commands to selected pane
- **Command History** - View and rerun previous commands
- **Multi-Claude Dashboard** - Manage multiple Claude sessions from one view
- **Save Session as Template** - Export current session layout to templates.json
- **Session Rename** - Rename sessions on the fly
- **Custom Themes** - User-configurable color schemes
- **Search/Filter** - Quick search for sessions and templates

## Development

### Project Structure

```
tmuxplexer/
├── main.go              # Entry point
├── types.go             # Type definitions
├── model.go             # Model initialization
├── update.go            # Message dispatcher
├── update_keyboard.go   # Keyboard handling
├── update_mouse.go      # Mouse handling
├── view.go              # View rendering
├── styles.go            # Lipgloss styles
├── config.go            # Configuration
├── templates.go         # Template management
├── tmux.go              # Tmux integration
├── claude_state.go      # Claude Code integration
├── hooks/               # Claude Code hooks for state tracking
│   ├── state-tracker.sh # Hook script for state updates
│   ├── install.sh       # Hooks installation script
│   └── test-hooks.sh    # Test suite for hooks
├── docs/                # Documentation
└── test_*.go            # Test commands
```

### Building

```bash
go build -o tmuxplexer
```

## Why Tmuxplexer?

**The Problem:** Managing 5+ terminal windows on an ultrawide monitor. Manually recreating layouts every day.

**The Solution:** Workspace templates that launch instantly. Save your perfect layouts once, reuse forever.

## License

MIT

## Author

Built with ❤️ using [Bubble Tea](https://github.com/charmbracelet/bubbletea)
