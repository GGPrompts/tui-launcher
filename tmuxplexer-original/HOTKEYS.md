# tmuxplexer Hotkeys & Commands

Quick reference for Tmux Session Manager keyboard shortcuts and template workflows.

## 🎯 Panel Navigation

### **3-Panel Unified Layout**
```
┌─ Command Panel (1) ────────────┐  ← Type commands for AI sessions
├─ Sessions/Templates Panel (2) ─┤  ← View sessions OR templates (toggle)
├─ Preview Panel (3) ────────────┤  ← Live pane preview with scrollback
└────────────────────────────────┘
```

### Focus Panels
```
1 - Focus Sessions/Templates panel (top)
    → Press again when focused to toggle Sessions ↔ Templates
2 - Focus Preview panel (middle)
3 - Focus Command panel (bottom)
```

### Panel Cycling
```
Tab       - Cycle through panels (forward)
Shift+Tab - Cycle through panels (backward)
```

**Note:** Each panel must be explicitly focused before its controls work. Use 1/2/3 keys or Tab to switch focus.

## 📋 Sessions Tab (Top Panel)

### Navigate Sessions
```
↑ / k - Move up
↓ / j - Move down
Home  - Jump to top
End   - Jump to bottom
```

### Window Navigation
```
← / h - Previous window
→ / l - Next window
```

### Session Actions
```
Enter   - Attach to selected session
d / D   - Kill selected session
s       - Save session as template (launches save wizard)
r       - Rename session
Ctrl+R  - Refresh sessions & Claude state
```

### Visual Indicators
```
◆       - Current session (cyan + bold)
●       - Attached session
○       - Detached session
🟢 🟡 🔧 - Claude Code status (idle/processing/tool use)
🔮      - Codex session
✨      - Gemini session
📁      - Working directory
branch  - Git branch (if in repo)
```

## 🎨 Templates Tab (Top Panel)

### Navigate Templates
```
↑ / k - Move up (through categories and templates)
↓ / j - Move down
Home  - Jump to top
End   - Jump to bottom
```

### Category Management
```
← / h - Collapse category
→ / l - Expand category
Enter - Toggle category (when on category line)
```

### Template Actions
```
Enter - Create session from selected template
n     - New template (launches creation wizard)
d     - Delete selected template
e     - Edit templates.json in $EDITOR
r     - Reload templates from disk
```

### Template Tree View
```
▼ Projects               ← Expanded category
  ├─ Simple Dev (2x2)    ← Template
  ├─ Frontend Dev (2x2)
  └─ TFE Development (4x2)
▶ Tools                  ← Collapsed category
▶ Agents
```

## 📺 Preview Panel (Middle)

### Scroll Content (Full Scrollback History)
```
PgUp     - Scroll up one page
PgDn     - Scroll down one page
Home / g - Jump to top
End / G  - Jump to bottom
r        - Refresh preview content
```

### Preview Features
```
• Full tmux scrollback history (not just visible area)
• Auto-scroll to bottom for Claude Code sessions
• Scroll position indicator: "Line 123-150 of 500 (45%)"
• Works when preview panel is focused
```

## 🤖 Command Panel (Bottom)

### Command Mode
```
3          - Focus command panel (required)
Type       - Enter command for AI session (only when focused)
Enter      - Send command to selected AI session
Esc        - Unfocus command panel
```

### Command Input Features
```
Ctrl+V   - Paste clipboard (multi-line → single line)
←/→      - Move cursor
Backspace - Delete character before cursor
Delete   - Delete character at cursor
↑/↓      - Browse command history (last 100 commands)
```

### AI Session Filtering
```
• Command mode auto-filters left panel to show ONLY AI sessions
• Prevents accidentally sending commands to production servers
• Shows: 🤖 Claude | 🔮 Codex | ✨ Gemini sessions only
• Filter removed when you exit command mode (Esc)
```

### Command Display
```
When focused:
┌─ Command ● ────────────────────────────────┐
│ Target: 🤖 claude-1 (claude)               │
│ > /clear█                                  │
│ [↑↓] History | [Ctrl+V] Paste | [Enter] Send
└────────────────────────────────────────────┘

When unfocused:
┌─ Command ──────────────────────────────────┐
│ Send commands to AI sessions               │
│ Press '1' or click to focus                │
│ Last: /clear → claude-1 (just now)         │
└────────────────────────────────────────────┘
```

## 🔄 Global Actions

### Refresh & Reload
```
Ctrl+R - Refresh sessions & Claude state
r      - Context-aware:
         • Sessions tab: Rename session
         • Templates tab: Reload templates
         • Preview panel: Refresh preview content
```

### Panel Behavior
```
a / A  - Toggle adaptive mode
         • ON: Panels resize based on focus (default)
         • OFF: Fixed panel heights (40%/40%/20%)
```

### Application Control
```
q       - Quit tmuxplexer
Ctrl+C  - Force quit
Esc     - Exit command mode / close popup
```

## 🤖 Claude Code Integration

### Claude Status Indicators
```
🟢 Idle           - Ready for input
🟡 Processing     - Handling user prompt
🔧 Tool Use       - Executing tool (shows tool name)
⚙️ Working        - Processing results
⏸️ Awaiting Input - Waiting for user
⚪ Stale          - No updates >60s (shows last known status)
```

### AI Tool Detection
```
🤖 Claude Code  - Claude AI assistant
🔮 Codex        - OpenAI Codex (via /codex slash command)
✨ Gemini       - Google Gemini (via /gemini slash command)
```

### Auto-Refresh
```
• Claude status updates every 2 seconds automatically
• State tracked via hooks in /tmp/claude-code-state/
• Shows working directory and git branch
• Manual refresh with Ctrl+R
```

## 🚀 Common Workflows

### Launch Development Workspace
```bash
# 1. Start tmuxplexer
tmuxplexer

# 2. Switch to Templates tab
1      # Focus top panel
1      # Toggle to Templates

# 3. Select template
↓↓     # Navigate to "Frontend Dev (2x2)"

# 4. Create session
Enter  # Session created with 4 panes!
```

### Manage Existing Sessions
```bash
# 1. View sessions
1      # Focus top panel (if not already)
# Sessions tab is default

# 2. Navigate to session
↓↓↓

# 3. Preview windows
←→     # Switch between windows in preview

# 4. Attach or kill
Enter  # Attach to session
d      # Kill session
```

### Save Session as Template
```bash
# 1. Set up a session exactly how you want it
# (manually in tmux: split panes, run commands, etc.)

# 2. Open tmuxplexer
tmuxplexer

# 3. Select the session
↓↓     # Navigate to your session

# 4. Save as template
s      # Launches save wizard

# 5. Fill in wizard
# - Name (pre-filled with session name)
# - Category (Projects/Agents/Tools/Custom)
# - Description (optional)

# 6. Template saved!
# Now appears in Templates tab for reuse
```

### Create Template from Scratch
```bash
# 1. Switch to Templates tab
1      # Focus top panel
1      # Toggle to Templates

# 2. Start wizard
n      # New template

# 3. Step through wizard
# - Template name
# - Description
# - Category
# - Working directory
# - Layout (2x2, 3x3, 4x2, etc.)
# - For each pane: command and title

# 4. Template auto-saved!
```

### Send Commands to AI Sessions
```bash
# 1. Focus command panel
1      # Or just start typing

# 2. Left panel auto-filters to AI sessions only
# Shows: 🤖 Claude | 🔮 Codex | ✨ Gemini

# 3. Select AI session
↓↓     # Navigate to target session

# 4. Type command
/clear # Or any command

# 5. Send
Enter  # Command executed in AI session

# 6. View output
3      # Focus preview panel
# See command output in preview
```

### Browse Template Categories
```bash
# 1. Switch to Templates tab
1      # Focus top panel
1      # Toggle to Templates

# 2. Expand/collapse categories
→      # Expand category
←      # Collapse category
Enter  # Toggle (when on category line)

# 3. Navigate tree
↓      # Move down through categories and templates
↑      # Move up

# 4. Create from template
Enter  # When on template line
```

## 📝 Template Examples

### Simple Dev (2x2)
```json
{
  "name": "Simple Dev (2x2)",
  "description": "Basic development workspace",
  "category": "Projects",
  "working_dir": "~/projects/myapp",
  "layout": "2x2",
  "panes": [
    {"command": "nvim", "title": "Editor"},
    {"command": "bash", "title": "Terminal"},
    {"command": "lazygit", "title": "Git"},
    {"command": "htop", "title": "Monitor"}
  ]
}
```

### Frontend Dev with Claude (2x2)
```json
{
  "name": "Frontend Dev (2x2)",
  "description": "React/Next.js development with AI",
  "category": "Projects",
  "working_dir": "~/projects/frontend",
  "layout": "2x2",
  "panes": [
    {"command": "claude-code .", "title": "Claude AI"},
    {"command": "nvim", "title": "Editor"},
    {"command": "npm run dev", "title": "Dev Server"},
    {"command": "lazygit", "title": "Git"}
  ]
}
```

### Multi-Worktree Dev (2x2)
```json
{
  "name": "Multi-Worktree Dev",
  "description": "Multiple git worktrees for parallel development",
  "category": "Projects",
  "working_dir": "~/projects/myapp",
  "layout": "2x2",
  "panes": [
    {
      "command": "claude-code .",
      "title": "Main",
      "working_dir": "~/projects/myapp"
    },
    {
      "command": "claude-code .",
      "title": "Feature",
      "working_dir": "~/projects/myapp-feature"
    },
    {"command": "lazygit", "title": "Git"},
    {"command": "bash", "title": "Terminal"}
  ]
}
```

### AI Research (3x3)
```json
{
  "name": "AI Research (3x3)",
  "description": "Multiple AI agents + development tools",
  "category": "Agents",
  "working_dir": "~/research",
  "layout": "3x3",
  "panes": [
    {"command": "claude-code .", "title": "Claude"},
    {"command": "opencode", "title": "OpenCode"},
    {"command": "~/projects/tmuxplexer/tmuxplexer", "title": "Sessions"},
    {"command": "~/TFE/TFE", "title": "Files"},
    {"command": "nvim", "title": "Editor"},
    {"command": "~/tkan/tkan", "title": "Tasks"},
    {"command": "~/gh-tui/gh-tui", "title": "GitHub"},
    {"command": "bash", "title": "Terminal"},
    {"command": "htop", "title": "Monitor"}
  ]
}
```

## 🎛️ Advanced Features

### Per-Pane Working Directories
Each pane can override the template's default working directory:
```json
{
  "working_dir": "~/projects",
  "panes": [
    {"command": "nvim", "working_dir": "~/projects/frontend"},
    {"command": "nvim", "working_dir": "~/projects/backend"}
  ]
}
```

### Template Categorization
```json
{
  "category": "Projects"    // Default categories: Projects, Agents, Tools
}
// Or use custom category:
{
  "category": "My Custom Category"
}
```

### CLI Integration (TFE, Scripts)
```bash
# Create session from template without TUI
tmuxplexer --template 0                    # Create from template 0
tmuxplexer --cwd /path/to/dir --template 1 # Override working directory

# Use cases:
# - TFE context menu integration
# - Shell scripts and automation
# - Project launch scripts
```

### Popup Mode (Ctrl+b o)
```bash
# Install tmux keybinding
./install.sh

# Usage:
# From any tmux session, press Ctrl+b o
# → Popup opens (80% width/height)
# → Navigate and select session
# → Press Enter to switch
# → Popup closes automatically
# → Press Esc to close without switching
```

## 🛠️ Configuration

### Templates Location
```bash
~/.config/tmuxplexer/templates.json
```

### Edit Templates
```bash
# From within tmuxplexer (Templates tab)
e    # Opens in $EDITOR

# Or manually
$EDITOR ~/.config/tmuxplexer/templates.json

# Reload after manual changes
r    # Reload templates
```

### Claude Code Hooks Installation
```bash
cd ~/projects/tmuxplexer
./hooks/install.sh  # Installs to ~/.claude/hooks/
```

### Default Editor Fallback
```bash
# Uses $EDITOR environment variable
# Falls back to: micro → nano → vim → vi
export EDITOR=nvim  # Set your preferred editor
```

## 🔧 Troubleshooting

### Session Won't Attach
```bash
# Check if session exists
tmux ls

# Kill stuck session (from tmuxplexer)
d      # Kill selected session

# Or manually
tmux kill-session -t session-name
```

### Template Not Working
```bash
# Validate JSON syntax
cat ~/.config/tmuxplexer/templates.json | jq .

# Check working directory exists
ls ~/projects/myapp

# Test command manually
nvim   # Does the command work?
```

### Claude Status Not Updating
```bash
Ctrl+R    # Force refresh sessions
r         # Refresh preview (when focused on preview)

# Check hooks installation
ls ~/.claude/hooks/state-tracker.sh

# Check state files
ls /tmp/claude-code-state/
```

### Template Tree Not Showing
```bash
# Reload templates
r      # When in Templates tab

# Check category field exists in templates
cat ~/.config/tmuxplexer/templates.json | jq '.[].category'
# Should show category names (auto-migrated to "Uncategorized" if missing)
```

## ⚙️ Keyboard Cheat Sheet

```
NAVIGATION
  1 2 3         Focus panels (top/middle/bottom)
  1 (again)     Toggle Sessions ↔ Templates (when focused)
  Tab           Cycle panels
  ↑↓ / kj       Navigate items
  ←→ / hl       Windows / Categories
  Home/End      Jump to top/bottom
  g / G         Jump to top/bottom (vim-style)

ACTIONS
  Enter         Attach session / Create from template / Toggle category
  s             Save session as template
  n             New template (wizard)
  d / D         Kill session / Delete template
  e             Edit templates.json
  r             Rename / Reload / Refresh (context-aware)

PREVIEW
  PgUp/PgDn     Scroll preview
  Home / g      Top of preview
  End / G       Bottom of preview
  r             Refresh preview

COMMAND MODE
  3             Focus command panel
  Type          Enter command (when focused)
  Enter         Send command to AI session
  Ctrl+V        Paste clipboard
  ↑↓            Browse command history
  Esc           Exit command mode

GLOBAL
  a / A         Toggle adaptive mode (dynamic vs fixed panel heights)
  Ctrl+R        Refresh sessions & Claude state
  X             Fullscreen mode (popup only)
  q             Quit
  Ctrl+C        Force quit
```

---

**Version**: tmuxplexer v2.0
**Last Updated**: 2025-01-03
**Config**: `~/.config/tmuxplexer/templates.json`
**Documentation**: See [CLAUDE.md](CLAUDE.md) for full architecture details
