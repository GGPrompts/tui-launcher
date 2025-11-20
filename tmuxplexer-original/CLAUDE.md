# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Tmuxplexer is a modern TUI (Terminal User Interface) tmux session manager built with Go and the Bubble Tea framework. It provides a 4-panel accordion layout for managing tmux sessions and creating complex multi-pane layouts from workspace templates.

**Current Status**: All core features implemented and working:
- ✅ Phase 1: 4-Panel Accordion Layout (refactored to 3-panel in Phase 10)
- ✅ Phase 2: Workspace Templates (create, edit, delete)
- ✅ Phase 3: Session Management (attach, kill, rename, auto-refresh)
- ✅ Phase 4: Live Pane Preview & Window Navigation
- ✅ Phase 5: Claude Code Integration & Layout Reorganization
- ✅ Phase 6: Scrollable Preview with Full Scrollback History
- ✅ Phase 7: Template Creation Wizard & Template Deletion
- ✅ Phase 8: Popup Mode with Tmux Keybinding (Ctrl+b o)
- ✅ Phase 9.1: Unified Chat/Command Mode (AI session commands)
- ✅ Phase 9.1.1: Template Categorization & Tree View
- ✅ Phase 10: Unified 3-Panel Adaptive Layout (explicit focus, click-to-focus)

**📋 Documentation:**
- **[PLAN.md](PLAN.md)**: Current roadmap and next phases (9.2, 10, 11+)
- **[docs/CHANGELOG.md](docs/CHANGELOG.md)**: Complete history of implemented features

## Build and Run Commands

```bash
# Build the application
go build -o tmuxplexer

# Run in TUI mode (requires a real terminal with TTY)
./tmuxplexer

# Run in popup mode (from within tmux)
./tmuxplexer --popup

# Create session from template via CLI (no TTY required)
./tmuxplexer --template 0                    # Create from template 0
./tmuxplexer --cwd /path/to/dir --template 1 # Override working directory

# Install tmux keybinding (Ctrl+b o)
./install.sh

# Test commands (work without TTY)
./tmuxplexer test_template           # View available templates
./tmuxplexer test_create 0           # Create session from template 0
```

**Important**: The TUI requires a proper terminal (Windows Terminal, iTerm2, etc.). It won't work in background processes or without TTY.

### CLI Flags for TFE Integration

**`--cwd <directory>`**: Override the template's working directory
- Use case: Launch templates in the current directory context
- Example: `tmuxplexer --cwd $PWD --template 0`

**`--template <index>`**: Create a session from template and exit (no TUI)
- Use case: Automated session creation from scripts or TFE context menu
- Example: `tmuxplexer --template 1`

**Combined usage for TFE integration:**
```bash
# From TFE context menu at /home/matt/projects/myapp
tmuxplexer --cwd /home/matt/projects/myapp --template 0

# This creates a session using template 0 but with myapp as the working directory
# instead of the hardcoded directory in the template
```

## Popup Mode

Tmuxplexer can be launched as a floating popup within any tmux session:

**Installation:**
```bash
./install.sh
```

This adds the keybinding `Ctrl+b o` to your `~/.tmux.conf`:
```bash
bind-key o run-shell "tmux popup -E -w 80% -h 80% -d '#{pane_current_path}' tmuxplexer --popup"
```

**Usage:**
1. From any tmux session, press `Ctrl+b o`
2. Tmuxplexer opens in a floating popup (80% width/height)
3. Navigate and select a session (or create from template)
4. Press `Enter` to switch to that session
5. The popup closes automatically
6. Press `ESC` or `q` to close popup without switching

**Key Differences from Normal Mode:**
- **Popup Mode**: Uses `tmux switch-client` (instant switch, stays in tmux)
- **Normal Mode**: Uses `tmux attach-session` (replaces current process)

**Why Popup Mode?**
- Quick session switching without leaving your current workflow
- Visual session browser overlaid on your work
- Stays within tmux (no process replacement)
- Perfect for managing multiple sessions on-the-fly

## Core Architecture

### Bubble Tea Pattern (Elm Architecture)

The codebase follows the Bubble Tea framework's Model-View-Update pattern:

- **Model** (`types.go`, `model.go`): Application state and initialization
- **Update** (`update.go`, `update_keyboard.go`, `update_mouse.go`): Message handling and state transitions
- **View** (`view.go`): Rendering logic

### File Organization by Responsibility

- `main.go` - Entry point ONLY. No business logic allowed here.
- `types.go` - All type definitions, structs, enums, and constants
- `model.go` - Model initialization and layout calculations
- `update.go` - Main update dispatcher for non-input events
- `update_keyboard.go` - Keyboard input handling
- `update_mouse.go` - Mouse input handling
- `view.go` - View rendering and layout
- `styles.go` - Lipgloss styling definitions
- `config.go` - Configuration management
- `templates.go` - Workspace template loading/saving
- `tmux.go` - Tmux integration (session/window/pane management)
- `claude_state.go` - Claude Code integration (state reading and detection)
- `hooks/` - Claude Code hooks for real-time state tracking
- `docs/` - Documentation and integration guides
- `test_*.go` - Test mode commands for non-TTY environments

### Unified 3-Panel Adaptive Layout

The UI uses a unified 3-panel vertical stack:

1. **Sessions/Templates Panel** (top, 40-50%, press `1`):
   - **Sessions Tab**: Active sessions list + details inline (shows Claude status)
   - **Templates Tab**: Template tree view with categories
   - Press `1` when focused to toggle between Sessions/Templates tabs
2. **Preview Panel** (middle, 30-50%, press `2`):
   - **Sessions Tab**: Live preview of selected session's active pane
   - **Templates Tab**: Template details (layout, panes, commands, working dirs)
   - Scrollable with arrow keys/PgUp/PgDn when focused
3. **Command Panel** (bottom, 20% fixed, press `3`): Command input for AI sessions

**Adaptive Height Distribution:**
- Sessions focused: 50% / 30% / 20% (sessions expanded)
- Preview focused: 30% / 50% / 20% (preview expanded)
- Command focused: Maintains previous upper panel sizing (no resize when focusing command)

**Focus Behavior:**
- Press `3` to focus command panel (required for typing commands)
- Press `1` to focus sessions/templates panel (use arrow keys to navigate)
- Press `2` to focus preview panel (use scroll keys to scroll)
- Each panel must be explicitly focused before its controls work

Layout calculations use `calculateAdaptivePanelHeights()` (model.go:770-814). Panels dynamically resize based on `focusState` for optimal workflow.

### Tmux Integration Layer

All tmux operations are isolated in `tmux.go`:

- **Session Management**: `listSessions()`, `createSession()`, `killSession()`, `attachToSession()`
- **Window/Pane Operations**: `listWindows()`, `listPanes()`, `capturePane()`
- **Template-Based Creation**: `createSessionFromTemplate()`, `createGridLayout()`

Key insight: The code handles both running inside tmux (uses `switch-client`) and outside tmux (uses `attach-session`) via `isInsideTmux()`.

### Template System

Templates are stored at `~/.config/tmuxplexer/templates.json`. Each template defines:

- `layout`: Grid dimensions (e.g., "2x2", "4x2", "3x3")
- `panes`: Array of pane configurations with commands and titles
- `working_dir`: Base directory for the session (can be overridden per-pane)
- `description`: Human-readable description of the template

**Per-Pane Working Directories:**
Each pane can optionally specify its own `working_dir`, enabling multi-worktree workflows:
```json
{
  "name": "Multi-Worktree Dev",
  "working_dir": "~/projects/myapp",
  "panes": [
    {"command": "claude-code .", "working_dir": "~/projects/myapp"},
    {"command": "claude-code .", "working_dir": "~/projects/myapp-feature"}
  ]
}
```

The `createGridLayout()` function creates all panes first, then applies tmux's "tiled" layout for even distribution, rather than manually tracking shifting pane indices.

### State Management

The model maintains separate state for each content type:

- `templates`: Loaded workspace templates (Templates tab)
- `sessions`: Current tmux sessions (Sessions tab)
- `selectedTemplate`: Index of selected template in Templates tab
- `selectedSession`: Index of selected session in Sessions tab
- `windows`: Windows for selected session
- `previewBuffer`: Full scrollback history lines for preview
- `previewScrollOffset`: Current scroll position in preview
- `previewTotalLines`: Total lines in preview buffer
- Panel content arrays: `sessionsContent`, `templatesContent`, `previewContent`, `commandContent`

Update functions (`updateSessionsContent()`, `updateTemplatesContent()`, `updatePreviewContent()`, `updateCommandContent()`) keep panel content synchronized with model state. Each panel independently manages its own selection and display logic.

### Auto-Refresh System

The application uses a ticker (update.go:137) that fires every 2 seconds to refresh session data automatically. This provides live updates of session status, window counts, pane content, and Claude Code state.

### Claude Code Integration

Real-time Claude Code status tracking via hooks system:

**Architecture:**
- `claude_state.go`: State file reading, session detection, status formatting
- `hooks/state-tracker.sh`: Bash script that captures hook events and writes JSON state files
- `/tmp/claude-code-state/*.json`: State files written by hooks (auto-cleaned after 24h)

**Detection Flow:**
1. Hook fires in Claude Code (SessionStart, UserPromptSubmit, PreToolUse, PostToolUse, Stop, Notification)
2. `state-tracker.sh` receives hook data via stdin
3. Script writes state to `/tmp/claude-code-state/{session-id}.json`
4. `listSessions()` calls `detectClaudeSession()` to check if session is running Claude
5. `getClaudeStateForSession()` reads state file and populates `ClaudeState` field
6. Sessions tab displays status icon next to session name

**Status Indicators (Enhanced with Tool Details):**
- 🟢 **Idle**: Ready for input
- 🟡 **Processing**: Handling user prompt
- 🔧 **Tool Use**: Shows tool name and details
  - `🔧 Read: model.go` - Reading a file
  - `🔧 Edit: types.go` - Editing a file
  - `🔧 Bash: npm test` - Running a command
  - `🔧 Grep: searchPattern` - Searching code
  - `🔧 Task: description` - Sub-agent task
- ⚙️ **Working**: Processing tool results (shows tool name)
- ⏸️ **Awaiting Input**: Waiting for user
- ⚪ **Stale**: No updates >60 seconds (shows last known status)

**Stale State Behavior:**
- State files older than 60 seconds are marked as "stale" but still displayed
- This ensures Claude sessions are visible even when idle for extended periods
- The stale indicator shows the last known status from the state file
- Once a hook fires (any Claude activity), the status updates to fresh
- On initial launch, idle Claude sessions will show with ⚪ icon until activity occurs

**Installation:**
```bash
cd ~/projects/tmuxplexer
./hooks/install.sh  # Installs hooks to ~/.claude/hooks/
```

**Cleanup & Maintenance:**

The hooks create state files and debug logs in `/tmp/claude-code-state/`. To prevent accumulation:

```bash
# Manual cleanup (removes files >7 days old)
~/projects/tmuxplexer/hooks/cleanup-state-files.sh

# Or add to crontab for automatic daily cleanup
crontab -e
# Add: 0 2 * * * ~/projects/tmuxplexer/hooks/cleanup-state-files.sh
```

Cleanup removes:
- State files older than 7 days
- Debug files older than 1 hour (they accumulate VERY quickly - 400+ files/day!)

**Optional: Tmux Status Bar Integration**

You can also display Claude activity in your tmux status bar for even more visibility:

```bash
# 1. Copy the status script to Claude hooks directory
cp ~/projects/tmuxplexer/hooks/tmux-status-claude.sh ~/.claude/hooks/
chmod +x ~/.claude/hooks/tmux-status-claude.sh

# 2. Add to your ~/.tmux.conf (replace or append to status-right)
set -g status-right '#(~/.claude/hooks/tmux-status-claude.sh) | %H:%M %d-%b'

# 3. Reload tmux configuration
tmux source-file ~/.tmux.conf
```

This shows real-time Claude status in the tmux status bar:
- `🟢 Ready` - Claude is idle
- `🟡 Processing` - Processing your request
- `🔧 Read: model.go` - Reading a specific file
- `🔧 Bash: npm test...` - Running a command
- `⚙️ Edit` - Processing edit results

See `docs/claude-hooks-integration.md` for complete technical documentation.

### Template Categorization & Tree View

Hierarchical organization of templates with collapsible categories:

**Tree View Display:**
- Templates organized by category in Templates tab (top panel)
- Categories shown with expand/collapse indicators: `▶` (collapsed) / `▼` (expanded)
- Templates nested under categories with tree connectors: `├─` / `└─`
- Visual indentation shows hierarchy
- Press `Enter` or `→` on category to expand, `←` to collapse
- Press `Enter` on template to create session
- Template details shown in preview panel (middle) when selected

**Category Management:**
- Default categories: "Projects", "Agents", "Tools", "Uncategorized"
- Custom categories supported (type any name in wizard)
- "Projects" category auto-expands on startup
- Category field stored in `SessionTemplate` struct
- Migration auto-assigns "Uncategorized" to old templates

**Implementation:**
- `Category` field in `SessionTemplate` (types.go:247)
- `expandedCategories` map tracks expansion state (model.go:47-49)
- `TemplateTreeItem` struct for flattened tree (types.go:270-280)
- `buildTemplateTreeItems()` (model.go:794-857): Groups templates by category
- `updateTemplateTreeItems()` (model.go:859-862): Rebuilds tree cache
- `updateRightPanelContent()` (model.go:354-451): Renders tree view
- `migrateTemplates()` (templates.go:70-79): Auto-migration for old templates

**Visual Structure:**
```
▼ Projects
  ├─ Simple Dev (2x2)
  ├─ Frontend Dev (2x2)
  └─ TFE Development (4x2)
▶ Tools
▶ Agents
```

### Template Creation Wizard

Interactive template creation via the 'n' key (Templates tab):

**User Flow:**
1. Switch to Templates tab (press `1` twice if needed)
2. Press `n` to start wizard
3. Step through fields sequentially:
   - **Step 1**: Template name (required)
   - **Step 2**: Description (optional, supports spaces)
   - **Step 3**: Category (defaults to "Uncategorized", suggests: Projects, Agents, Tools, Custom)
   - **Step 4**: Working directory (defaults to `~`)
   - **Step 5**: Layout (e.g., `2x2`, `3x3`, `4x2`)
   - **Step 6+**: For each pane:
     - Pane command (e.g., `nvim`, `bash`, `claude-code .`)
     - Pane title (optional)
4. Template auto-saves to `~/.config/tmuxplexer/templates.json`
5. Template list refreshes automatically with new category

**Implementation:**
- `TemplateBuilder` (types.go): Tracks wizard state across steps
- `startTemplateCreation()` (update_keyboard.go): Initializes wizard
- `handleTemplateCreationInput()` (update_keyboard.go): Handles keyboard input per step
- `advanceTemplateWizard()` (update_keyboard.go): Validates and advances to next field
- `calculatePaneCount()` (update_keyboard.go): Parses layout string (e.g., "2x2" → 4 panes)
- `addTemplate()` (templates.go): Appends template to templates.json
- `getTemplateWizardPrompt()` (view.go): Renders step-specific prompt with progress

**Features:**
- Press ESC at any time to cancel
- Spaces allowed in template names/descriptions
- Progress indicator shows current step
- Auto-calculates pane count from layout

### Save Session as Template

Save existing tmux sessions as templates via the 's' key (Sessions tab):

**User Flow:**
1. Ensure you're on Sessions tab (press `1` if needed)
2. Select a running session (↑/↓)
3. Press `s` to save session as template
4. Application extracts session info automatically:
   - Detects grid layout (2x2, 3x3, etc.)
   - Captures working directory for each pane
   - Captures running command in each pane
   - Captures pane titles (if set)
5. Step through wizard:
   - **Step 1**: Template name (pre-filled with session name)
   - **Step 2**: Category (defaults to "Uncategorized", suggests: Projects, Agents, Tools, Custom)
   - **Step 3**: Description (optional)
6. Template auto-saves to `~/.config/tmuxplexer/templates.json`
7. Template list refreshes automatically under the selected category

**Key Features:**
- **"Configure by Example"**: Set up a session exactly how you want it, then save it
- **Per-Pane Working Directories**: Each pane's working directory is preserved
- **Smart Layout Detection**: Automatically detects grid layouts
- **Common Working Dir**: Uses the most common directory as the template default
- **Instant Replication**: Saved templates appear in the Templates tab for reuse

**Implementation:**
- `extractSessionInfo()` (tmux.go): Extracts pane info (working dirs, commands, dimensions)
- `detectGridLayout()` (tmux.go): Detects grid pattern from pane positions
- `SessionSaveBuilder` (types.go): Tracks save wizard state
- `extractSessionCmd()` (update_keyboard.go): Command to extract session
- `handleSessionSaveInput()` (update_keyboard.go): Handles wizard input
- `sessionExtractedMsg` (update.go): Message handler for extracted session data

**Example Workflow:**
```bash
# 1. Create a session manually with tmux
tmux new-session -s mydev -c ~/projects/myapp
# Split into 2x2 grid, run commands...

# 2. Open tmuxplexer
./tmuxplexer

# 3. Press '1' to focus sessions, select 'mydev', press 's'
# 4. Enter name: "My Dev Setup"
# 5. Template saved!

# 6. Later: Press '2', select "My Dev Setup", press Enter
# New session created with identical layout and working directories
```

**Advantages Over Manual Creation:**
- No need to remember layout configurations
- Pane working directories are preserved
- Commands are captured for reference
- Faster than using the manual wizard for complex setups

### Template Deletion

Delete templates via the 'd' key (Templates tab):

**User Flow:**
1. Switch to Templates tab (press `1` twice if needed)
2. Select template to delete (↑/↓)
3. Press `d` to delete
4. Confirm deletion (y/n)
5. Template removed from `~/.config/tmuxplexer/templates.json`

**Implementation:**
- Confirmation prompt prevents accidental deletion
- `deleteTemplate()` (templates.go): Removes template at index
- Auto-refreshes template list and adjusts selection

### Template Editing

In-app template editing via the 'e' key (Templates tab):

**User Flow:**
1. Switch to Templates tab (press `1` twice if needed)
2. Press `e` to open templates.json in editor
3. Edit/save in editor (micro, nano, vim, etc.)
4. Editor closes → templates automatically reload
5. Templates tab updates with new templates

**Editor Detection:**
Uses `$EDITOR` environment variable, falling back to: micro → nano → vim → vi

**Implementation:**
- `getUserEditor()` (templates.go): Detects user's preferred editor
- `openTemplatesInEditor()` (templates.go): Opens templates.json, blocks until editor exits
- `editTemplatesCmd()` (update_keyboard.go): Command that opens editor then reloads templates
- `templatesReloadedMsg` (update.go): Message handler that updates templates in model

### Scrollable Preview

Full scrollback history viewing in the footer panel (Phase 6):

**Architecture:**
- `capturePane()` (tmux.go): Modified to use `-S -` flag for full history capture (not just visible area)
- `previewBuffer` (types.go): Stores complete pane content as array of lines
- `previewScrollOffset` (types.go): Current scroll position (line number)
- `previewTotalLines` (types.go): Total lines in buffer

**User Flow:**
1. Select a session or template to populate preview
2. Focus preview panel (press `2`)
3. Use arrow keys/PgUp/PgDn to scroll through content
4. Use Home/End or g/G to jump to top/bottom
5. Scroll position indicator shows: "Scroll: 45% (Line 123-150 of 500)"

**Implementation Details:**
- `updateFooterContent()` (model.go): Captures full pane content, stores in buffer, calculates visible window based on scroll offset
- `pageUp()`/`pageDown()` (update_keyboard.go): Scrolls by viewport height (footer panel height - header - borders)
- `moveToTop()`/`moveToBottom()` (update_keyboard.go): Jumps to start/end of buffer
- Auto-reset: Scroll offset resets to 0 when changing sessions or windows

**Key Design Decisions:**
- Scrolling only works when footer panel is focused (prevents accidental scrolling)
- Full history capture means Claude Code output history is fully accessible
- Scroll position indicator only shows when content exceeds viewport (avoids clutter for short content)
- Status message updates with current line position during scrolling

### Claude Code Preview Features

Special preview enhancements for Claude Code sessions:

**Auto-Scroll to Bottom:**
When a session is running Claude Code (detected via `ClaudeState`), the preview automatically scrolls to the bottom of the terminal on initial load. This shows the current conversation instead of empty space at the top.

**Implementation:**
- In `updateFooterContent()` (model.go:455-462): After capturing pane content, checks if `session.ClaudeState != nil`
- If Claude session and scroll offset is 0, sets offset to bottom: `maxOffset = totalLines - visibleHeight`
- User can still scroll up to view history using PgUp or scroll down to latest with PgDn/End

**Manual Refresh:**
- Press `r` when footer panel is focused to manually refresh preview content
- Useful for forcing an update without waiting for auto-refresh (2-second interval)
- Updates help text in header shows `[r] Refresh` hint

**Status Display:**
Claude sessions show full status under session names in the Sessions tab, with **orange text** to make them visually distinct:
```
○ myproject 🟢                       (regular session - normal text)
  📁 ~/projects/myproject  main
  🟢 Idle
● claude-session 🔧                   (Claude session - ORANGE TEXT, bold)
  📁 ~/projects/tmuxplexer  feature-branch
  🔧 Using Read
```

**Visual Features:**
- **Current session** (where you are now) is rendered in **cyan (#56B6C2)** and **bold** with `◆` marker
- **Claude sessions** are rendered in **orange (#D19A66)** and **bold**
- The filled/unfilled bullet (●/○) indicates attached/detached status
- Status icon appears on the same line as the session name
- Full status text appears on the third line

This allows viewing all sessions at a glance with instant visual identification of your current location and AI sessions.

### Directory and Git Branch Display

Each session now displays its working directory and git branch (if in a git repository):

**Display Format:**
- Line 1: Session name with status icon
- Line 2: 📁 Working directory  git-branch (if in git repo)
- Line 3: Claude status (if Claude session)

**Features:**
- Home directory shortened to `~` for readability
- Git branch detected automatically via `git rev-parse --abbrev-ref HEAD`
- Non-git directories show directory only
- Updates with auto-refresh (2-second interval)

**Implementation:**
- `getGitBranch()` (tmux.go:24-31): Detects git branch for a directory
- `WorkingDir` and `GitBranch` fields added to `TmuxSession` (types.go:162-163)
- Working directory retrieved via tmux's `#{pane_current_path}` variable (tmux.go:84-90)
- Display formatting in `updateLeftPanelContent()` (model.go:199-218)

**Claude Session Styling:**
- `colorOrange` (#D19A66) added to color palette (styles.go:21)
- `claudeSessionStyle` style defined with orange, bold text (styles.go:90-92)
- "CLAUDE:" prefix tag added to Claude session lines (model.go:200-204)
- Tag processed in `renderDynamicPanel()` to apply orange style (view.go:366-369)

## Common Development Patterns

### Adding a New Tmux Operation

1. Add the function to `tmux.go` (keep all tmux commands isolated here)
2. Create a message type in `types.go` (e.g., `sessionRenamedMsg`)
3. Add command creator in `update.go` (e.g., `renameSessionCmd()`)
4. Handle the message in `Update()` method (update.go:22)
5. Update relevant panel content in model

### Modifying Layout Behavior

Layout calculations are centralized in `model.go`:
- `calculateFourPanelLayout()`: Panel dimension calculations
- `renderFourPanelLayout()`: Panel rendering (view.go:190)
- `renderDynamicPanel()`: Individual panel styling with borders

### Adding Keyboard Shortcuts

Add handling in `update_keyboard.go`. The file contains focused logic for keyboard navigation, panel switching, and session operations.

### Command Mode (Phase 9.1: Unified Chat) ✅ COMPLETED

Send commands to AI sessions without attaching via the header panel command interface:

**Key Feature:** **AI Sessions Only** - Command mode automatically filters to show only Claude Code, Codex, and Gemini sessions. This prevents accidentally sending commands to production servers or non-AI sessions.

**User Flow:**
1. Press `3` (or click command panel) to focus command panel
2. **Sessions tab auto-filters** to show only AI sessions (Claude, Codex, Gemini)
3. Navigate with ↑/↓ to select target AI session
4. Type command (e.g., `/clear` for Claude, `git status`, `ls -la`)
5. Press `Enter` to send command to selected AI session
6. Command executes in session's active pane
7. Press `Esc` to return to sessions tab (filter removed)

**Features:**
- **AI Session Filtering**: Automatically shows only Claude/Codex/Gemini when in command mode
- **Tool Icons**: 🤖 Claude | 🔮 Codex | ✨ Gemini
- **Command History**: Use ↑/↓ arrows to browse history (last 100 commands)
- **Cursor Navigation**: ←/→ arrows move cursor, Backspace/Delete edit
- **Clipboard Paste**: Ctrl+V pastes clipboard at cursor (multi-line converted to single line)
- **Multi-line Wrapping**: Long commands wrap to fill the header panel (adapts when panel expands)
- **Smart Viewport**: Cursor always visible; scroll indicators show hidden content (↑ more above... / ↓ more below...)
- **Long Command Handling**: Commands >100 chars show character count
- **Last Command Display**: When unfocused, shows last sent command and target
- **Session Stats**: Shows AI session breakdown (e.g., "🤖 2 AI sessions | Claude:1 | Codex:1")

**Implementation:**
- `AITool` field (types.go:174): Marks sessions as "claude", "codex", "gemini", or ""
- `detectCodexSession()`, `detectGeminiSession()` (claude_state.go:39-58): AI tool detection
- AI session filtering (model.go:169-181, update_keyboard.go:1100-1106): Filter logic
- `commandInput` (types.go:73-79): Command mode state
- `updateHeaderContent()` (model.go:379-405): Renders command UI with AI session target
- `executeCommand()` (update_keyboard.go:1098-1137): Sends command to AI sessions only
- `sendKeysToSession()` (tmux.go:635-642): Tmux send-keys wrapper

**Header Panel States:**

*Focused with Very Long Command (Viewport with Scroll Indicators):*
```
┌─ Command ● ────────────────────────────────┐
│ Target: 🤖 claude-1 (claude)               │
│   ↑ more above...                          │
│ oject:** matt **Problem Description:** {{ │
│ PROBLEM_DESCRIPTION}} **Error Message (i █│
│   ↓ more below...                          │
│ [↑↓] History | [Ctrl+V] Paste | 487 chars │
└────────────────────────────────────────────┘
```

*Focused with Long Command (All Lines Fit):*
```
┌─ Command ● ────────────────────────────────┐
│ Target: 🤖 claude-1 (claude)               │
│ > I'm debugging an issue in: **File:** /h │
│ ome/matt/.prompts/debug-help.prompty **Pr │
│ oject:** matt **Problem Description:** {{ │
│ PROBLEM_DESCRIPTION}} **Error Message █   │
│ [↑↓] History | [Ctrl+V] Paste | 487 chars │
└────────────────────────────────────────────┘
```

*Focused with Short Command:*
```
┌─ Command ● ────────────────────────────────┐
│ Target: 🤖 claude-1 (claude)               │
│ > /clear█                                  │
│ [↑↓] History | [Ctrl+V] Paste | [Enter] Send | [Esc] Unfocus │
└────────────────────────────────────────────┘
```

*Unfocused (Hint):*
```
┌─ Command ──────────────────────────────────┐
│ Send commands to AI sessions               │
│ Press '1' or click to focus                │
│ Last: /clear → claude-1 (just now)         │
└────────────────────────────────────────────┘
```

**Sessions Tab (When in Command Mode):**
```
┌─ Sessions ─────────────────────────────────┐
│ 🤖 2 AI sessions | Claude:1 | Codex:1      │
│ ────────────────────────────────────────── │
│ ► claude-1 🟢                              │
│   📁 ~/projects/feature-auth  main         │
│   🟢 Idle                                  │
│                                            │
│   codex-debug 🔮                           │
│   📁 ~/projects/bugfix  hotfix-123         │
└────────────────────────────────────────────┘
```

**Use Cases:**
- **Send `/clear` to Claude** for fresh context without attaching
- **Query Codex** for root cause analysis via slash commands
- **Ask Gemini** for alternative approaches to problems
- Send git commands to AI sessions for version control tasks
- Review command output in preview panel (press `4`)
- **Safe by design**: Can't accidentally send to production servers (only AI tools shown)
- **Multi-agent orchestration**: Perfect for coordinating multiple Claude/Codex/Gemini sessions (Phase 9.2+)

**Future Enhancements (Phase 9.2):**
- Multi-selection with checkboxes (Space to toggle)
- Send to multiple sessions simultaneously
- Command snippets/autocomplete
- Broadcast modes (all Claude sessions, all attached, etc.)

## Configuration

Templates are loaded from `~/.config/tmuxplexer/templates.json`. If the file doesn't exist, `loadTemplates()` creates it with default templates automatically.

## Keyboard Shortcuts

Current implementation (see `update_keyboard.go`):

| Key | Action | Implementation |
|-----|--------|---------------|
| `1` | Focus Sessions/Templates panel (top) / Toggle tab when focused | Sets `focusState = FocusSessions`, or toggles `sessionsTab` if already focused |
| `2` | Focus Preview panel (middle) | Sets `focusState = FocusPreview` |
| `3` | Focus Command panel (bottom) | Sets `focusState = FocusCommand` - required to type commands |
| `Tab/Shift+Tab` | Cycle through panels | Cycles `focusState` forward/backward |
| `↑/↓` or `k/j` | Navigate sessions/templates | `moveUp()` / `moveDown()` - navigates within focused panel |
| `←/→` or `h/l` | Navigate windows / expand/collapse categories | `moveLeft()` / `moveRight()` - windows in sessions tab, categories in templates tab |
| `PgUp/PgDn` | Scroll preview | `pageUp()` / `pageDown()` - scrolls within focused preview panel |
| `Home` or `g` | Jump to top of preview | `moveToTop()` - jumps to top of focused preview panel |
| `End` or `G` | Jump to bottom of preview | `moveToBottom()` - jumps to bottom of focused preview panel |
| `Enter` | Attach to session / create from template / toggle category | `selectItem()` - context-aware: attaches (sessions tab), creates session (templates tab), toggles category expansion |
| `s` | Save session as template (Sessions tab only) | `extractSessionCmd()` - extracts current session layout and opens save wizard |
| `n` | Create new template from scratch (Templates tab only) | `startTemplateCreation()` - opens interactive template creation wizard |
| `r` | Refresh preview or rename session | `updatePreviewContent()` / rename mode |
| `d` or `D` | Kill session | `killSessionCmd()` |
| `Ctrl+R` | Refresh sessions & Claude state | `refreshSessionsCmd()` |
| `Ctrl+V` | Paste clipboard (command mode only) | Inserts clipboard at cursor, converts multi-line to single line |
| `ESC` | Exit command mode / close popup | Returns to sessions panel or quits in popup mode |
| `q` or `Ctrl+C` | Quit | `tea.Quit` |
| `X` (popup mode only) | Detach and launch fullscreen | Exits popup, launches tmuxplexer in full terminal |

**Focus Behavior:**
- Use number keys (1/2/3) or Tab/Shift+Tab to switch focus between panels
- Each panel must be focused to use its controls
- Click any panel to focus it manually

## Tmux Usage Tips

### Detaching from Sessions

When you're inside a tmux session (e.g., running Pyradio or other long-running processes), you can detach and leave them running in the background:

**How to Detach:**
- Press `Ctrl+b` then `d` (tmux's default detach keybinding)
- This leaves the session running in the background
- All processes continue running (music keeps playing, builds keep running, etc.)

**Re-attach via Tmuxplexer:**
1. Launch tmuxplexer: `./tmuxplexer` or `Ctrl+b o` (popup)
2. Navigate to the session in Sessions tab (↑/↓)
3. Your current session is marked with `◆` in the session list
4. Press `Enter` to switch to/attach to any session

**In Popup Mode:**
- Press `ESC` or `D` to quickly close the popup without switching sessions
- Your current session stays active

**Example Workflow:**
```bash
# Create a session with Pyradio
./tmuxplexer  # Press Enter on "Music" template
# Pyradio starts playing music

# Detach: Ctrl+b, d
# Music keeps playing in background, you're back to shell

# Later: re-attach
./tmuxplexer  # Select session, press Enter
# Back to Pyradio, music still playing
```

This is the core power of tmux: persistent sessions that survive detachment and terminal closure.

## Testing Without TTY

Use test commands for development without a full terminal:
- `test_template`: Validates template loading
- `test_create N`: Tests session creation from template N

These commands bypass the TUI and provide direct command-line output.

## Reference: TUITemplate Project

When adding new features or components to Tmuxplexer, reference the **TUITemplate** project at `~/projects/TUITemplate`.

### What is TUITemplate?

TUITemplate is a production-ready template for building TUI applications with Go, Bubbletea, and Lipgloss. It contains reusable components, patterns, and comprehensive documentation extracted from real-world TUI development.

### Available Components

**UI Components** (`~/projects/TUITemplate/components/`):
- `dialog/` - Confirm dialogs, input dialogs, progress dialogs, modals
- `input/` - Text input, multiline input, forms, autocomplete
- `list/` - Simple lists, filtered lists, tree views
- `menu/` - Context menus, command palettes, menu bars
- `panel/` - Single-pane, dual-pane, multi-panel, tabbed layouts
- `preview/` - Text preview, markdown, syntax highlighting, images, hex viewer
- `status/` - Status bars, title bars, breadcrumbs
- `table/` - Simple and interactive tables

**Utility Libraries** (`~/projects/TUITemplate/lib/`):
- `clipboard/` - System clipboard integration
- `config/` - YAML configuration with hot-reload and validation
- `keybindings/` - Customizable keyboard shortcuts
- `logger/` - Debug logging to file
- `terminal/` - Terminal capability detection
- `theme/` - Theme system with color schemes

### How to Use as Reference

1. **Adding a New Component**: Copy the component from TUITemplate and adapt it
   ```bash
   cp -r ~/projects/TUITemplate/components/dialog ./components/
   # Adapt the component to Tmuxplexer's needs
   ```

2. **Implementing New Features**: Check TUITemplate's examples and documentation
   - `~/projects/TUITemplate/examples/` - Working example applications
   - `~/projects/TUITemplate/docs/research/` - Research on 80+ TUI tools and 20+ libraries

3. **Best Practices**: Reference TUITemplate's architecture patterns
   - Component isolation and reusability
   - Error handling patterns
   - Performance optimizations (lazy loading, virtual scrolling)
   - Testing approaches

### Example Use Cases

- **Adding a confirmation dialog** → Use `components/dialog/confirm.go` as reference
- **Implementing a command palette** → See `components/menu/command_palette.go`
- **Adding syntax highlighting** → Check `components/preview/syntax.go`
- **Improving input handling** → Reference `lib/keybindings/`
- **Adding themes** → See `lib/theme/` for theme system implementation

### Documentation

Comprehensive research available at:
- `~/projects/TUITemplate/docs/research/ECOSYSTEM_QUICK_REFERENCE.md` - Bubbletea ecosystem
- `~/projects/TUITemplate/docs/research/TUI_APPLICATIONS.md` - 80+ TUI tools analysis
- `~/projects/TUITemplate/ARCHITECTURE.md` - Architecture patterns and best practices
