# Tmuxplexer Changelog

All completed phases and features documented here.

---

## Phase 1: 4-Panel Accordion Layout ✅ COMPLETED

**Status:** Fully implemented and working perfectly

### Core Layout
- ✅ Four-panel design (Header, Left, Right, Footer)
- ✅ Weight-based dynamic panel sizing with focus expansion
- ✅ Panel focus switching (keys: 1, 2, 3, 4)
- ✅ Accordion toggle mode (key: a)
- ✅ Inline headers in borders (lazygit style) - saves 4 lines
- ✅ Perfect vertical border alignment
- ✅ Click any panel to focus (mouse support)

### Panel Functions
- **Header Panel**: Session statistics, filters, quick actions
- **Left Panel**: Active sessions list with details inline
- **Right Panel**: Templates list with details inline
- **Footer Panel**: Live preview of selected session's active pane

**Key Achievement:** Clean, space-efficient layout with no alignment issues

---

## Phase 2: Workspace Templates ✅ COMPLETED

**Status:** Full template system implemented

### Template Management
- ✅ Templates stored in `~/.config/tmuxplexer/templates.json`
- ✅ Visual distinction between templates (offline) and sessions (online)
- ✅ Create session from template (Enter key on right panel)
- ✅ Support for multiple layouts: 2x2, 3x3, 4x2, custom grids
- ✅ Per-pane commands and titles
- ✅ Per-pane working directories (multi-worktree support)
- ✅ Template metadata (name, description)

### Template Operations
- ✅ **Create from template**: Press Enter → instant workspace
- ✅ **Save session as template**: Press 's' on left panel → preserve current layout
- ✅ **Edit templates**: Press 'e' on right panel → opens in $EDITOR
- ✅ **Delete template**: Press 'd' on right panel with confirmation
- ✅ **Template wizard**: Interactive creation flow (press 'n' on right panel)

### Template Wizard Features
- ✅ Step-by-step field input (name, description, working dir, layout)
- ✅ Auto-calculates pane count from layout string
- ✅ Per-pane command and title configuration
- ✅ Progress indicator shows current step
- ✅ ESC to cancel at any time

### Save Session as Template
- ✅ Extract current session layout and pane info
- ✅ Detect grid layout automatically (2x2, 3x3, etc.)
- ✅ Capture working directory per pane
- ✅ Capture running commands in each pane
- ✅ Capture pane titles (if set)
- ✅ Save wizard with pre-filled session name

**Key Achievement:** "Configure by Example" - set up session manually, then save as template

---

## Phase 3: Session Management ✅ COMPLETED

**Status:** Core session operations working

### Basic Operations
- ✅ List all tmux sessions with status
- ✅ Attach to session (Enter key or click)
- ✅ Kill session (key: d/D with confirmation)
- ✅ Rename session (key: r, inline editing)
- ✅ Create new session (basic name input)
- ✅ Auto-refresh every 2 seconds (ticker-based)

### Session Display
- ✅ Session name with attached/detached indicator (●/○)
- ✅ Window count and pane count
- ✅ Working directory with tilde expansion (~)
- ✅ Git branch detection and display
- ✅ Claude Code status integration (icons and text)
- ✅ Visual indicators for session state

### Status Indicators
- ✅ Attached sessions marked with filled bullet (●)
- ✅ Detached sessions marked with unfilled bullet (○)
- ✅ Current session marked with diamond (◆) in popup mode

**Key Achievement:** Real-time session monitoring with auto-refresh

---

## Phase 4: Live Pane Preview & Window Navigation ✅ COMPLETED

**Status:** Full preview system implemented

### Preview Features
- ✅ Live pane content capture in footer panel
- ✅ Window navigation (←/→ or h/l keys)
- ✅ Preview updates automatically with session selection
- ✅ Preview refreshes with auto-refresh ticker (2 seconds)
- ✅ Manual refresh (key: r when footer focused)

### Display
- ✅ Shows active pane content for selected window
- ✅ Window indicator shows current/total windows
- ✅ Clean rendering with borders

**Key Achievement:** See session output without attaching

---

## Phase 5: Claude Code Integration ✅ COMPLETED

**Status:** Real-time Claude status detection working

### Hooks Integration
- ✅ Bash hooks system (`hooks/state-tracker.sh`)
- ✅ State files written to `/tmp/claude-code-state/*.json`
- ✅ Hook events captured: SessionStart, UserPromptSubmit, PreToolUse, PostToolUse, Stop, Notification
- ✅ State file reading and parsing
- ✅ Session detection via pane introspection

### Status Indicators
- ✅ Real-time status icons in session list:
  - 🟢 Idle (ready for input)
  - 🟡 Processing (handling user prompt)
  - 🔧 Tool Use (executing tool)
  - ⚙️ Working (processing results)
  - ⏸️ Awaiting Input (waiting for user)
  - ⚪ Stale (>60s old, shows last known status)

### Display Integration
- ✅ Claude session names highlighted in **orange** and **bold**
- ✅ Status icon next to session name
- ✅ Full status text on third line
- ✅ Stale state handling (state files older than 60s marked as stale)

### Installation
- ✅ Hook installation script (`hooks/install.sh`)
- ✅ Documentation in `docs/claude-hooks-integration.md`

**Key Achievement:** At-a-glance Claude status monitoring without selecting session

---

## Phase 6: Scrollable Preview with Full Scrollback ✅ COMPLETED

**Status:** Full history scrolling implemented

### Scrolling Features
- ✅ Captures full pane history (`tmux capture-pane -S -`)
- ✅ Stores complete pane content as array of lines
- ✅ PgUp/PgDn scrolling (by viewport height)
- ✅ Home/End and g/G keys (jump to top/bottom)
- ✅ Scroll position indicator: "Scroll: 45% (Line 123-150 of 500)"
- ✅ Auto-reset scroll position when changing sessions

### Claude-Specific Features
- ✅ Auto-scroll to bottom for Claude sessions on initial load
- ✅ Shows current conversation instead of empty top of terminal
- ✅ User can still scroll up to view history

### UX
- ✅ Scrolling only works when footer panel focused (prevents accidents)
- ✅ Scroll indicator only shows when content exceeds viewport
- ✅ Status message updates with current line position

**Key Achievement:** Full Claude Code conversation history accessible

---

## Phase 7: Template Creation Wizard & Deletion ✅ COMPLETED

**Status:** Complete template management workflow

### Creation Wizard (key: n on right panel)
- ✅ Interactive step-by-step creation
- ✅ Step 1: Template name (required)
- ✅ Step 2: Description (optional, supports spaces)
- ✅ Step 3: Working directory (defaults to ~)
- ✅ Step 4: Layout (e.g., 2x2, 3x3, 4x2)
- ✅ Step 5+: Per-pane configuration (command, title)
- ✅ Auto-save to templates.json
- ✅ Auto-refresh template list

### Template Deletion (key: d on right panel)
- ✅ Confirmation prompt (y/n)
- ✅ Removes template from templates.json
- ✅ Auto-refresh template list
- ✅ Adjusts selection after deletion

### Template Editing (key: e on right panel)
- ✅ Opens templates.json in $EDITOR
- ✅ Editor detection: micro → nano → vim → vi
- ✅ Blocks until editor exits
- ✅ Auto-reload templates after save

**Key Achievement:** Complete template lifecycle management in-app

---

## Phase 8: Popup Mode ✅ COMPLETED

**Status:** Fully functional popup mode

### Popup Features
- ✅ Launch as floating popup within tmux (`--popup` flag)
- ✅ Keybinding: `Ctrl+b o` (via install.sh)
- ✅ 80% width/height by default
- ✅ Opens in current pane's working directory
- ✅ ESC or D to close popup without switching
- ✅ Enter to switch to selected session

### Behavior
- ✅ Uses `tmux switch-client` (instant switch, stays in tmux)
- ✅ Normal mode uses `tmux attach-session` (replaces process)
- ✅ Popup closes automatically after session selection

### Installation
- ✅ `install.sh` script adds keybinding to `~/.tmux.conf`
- ✅ Works alongside tmux-sessionx (different keybinding)

**Key Achievement:** Quick session switching without leaving workflow

---

## Phase 9.1: Unified Chat/Command Mode ✅ COMPLETED

**Status:** AI session command interface fully implemented

### Command Interface
- ✅ Focus header panel (key: 1) to enter command mode
- ✅ Type commands with cursor navigation (←/→)
- ✅ Command history (↑/↓ arrows, last 100 commands)
- ✅ Enter to send, Esc to exit
- ✅ Backspace/Delete for editing

### AI Session Filtering
- ✅ **Automatic filtering**: Shows only AI sessions (Claude, Codex, Gemini) when in command mode
- ✅ AI tool detection: `detectClaudeSession()`, `detectCodexSession()`, `detectGeminiSession()`
- ✅ Tool icons: 🤖 Claude | 🔮 Codex | ✨ Gemini
- ✅ `AITool` field in session struct

### Command Execution
- ✅ Send commands to selected AI session's active pane
- ✅ Uses `tmux send-keys` wrapper
- ✅ Works without attaching to session
- ✅ Preview output in footer panel (press 4)

### Header Panel States
- ✅ **Focused**: Shows target, command input, help text
- ✅ **Unfocused**: Shows hint and last command sent

### Session Stats
- ✅ Shows AI session breakdown in left panel when in command mode
- ✅ Format: "🤖 2 AI sessions | Claude:1 | Codex:1"

### Safety
- ✅ **AI-only filter**: Prevents accidentally sending commands to production servers
- ✅ Only Claude/Codex/Gemini sessions shown in command mode

**Key Achievement:** Send commands to AI sessions without attaching (perfect for `/clear`, git commands, etc.)

---

## Phase 9.1.1: Template Categorization & Tree View ✅ COMPLETED

**Status:** Hierarchical template organization fully implemented

### Tree View Display
- ✅ Category-based tree view in right panel
- ✅ Expand/collapse categories with Enter key
- ✅ Tree connectors: `▶/▼` for categories, `├─/└─` for templates
- ✅ Indentation and visual hierarchy
- ✅ Category names show at top level (depth 0)
- ✅ Templates nested under categories (depth 1)

### Category Management
- ✅ `Category` field added to `SessionTemplate` struct
- ✅ Default categories: "Projects", "Agents", "Tools", "Uncategorized"
- ✅ Custom category support (user can type any category name)
- ✅ "Projects" category auto-expands on startup
- ✅ Expansion state tracked in `expandedCategories` map

### Template Creation Wizard Updates
- ✅ Category selection added as **Step 3** (after description)
- ✅ Wizard flow: Name → Description → **Category** → Working Dir → Layout → Panes
- ✅ Step numbers updated from `/6` to `/7` throughout
- ✅ Category prompt suggests: "Projects, Agents, Tools, Custom, or type your own"
- ✅ Defaults to "Uncategorized" if user skips

### Session Save Wizard Updates
- ✅ Category selection added as **Step 2** (after name)
- ✅ Save flow: Name → **Category** → Description
- ✅ Step numbers updated from `2/2` to `3/3`
- ✅ Consistent category prompt and defaults

### Template Migration
- ✅ `migrateTemplates()` function auto-assigns "Uncategorized" to old templates
- ✅ Migration runs automatically on `loadTemplates()`
- ✅ Migrated templates auto-saved back to disk
- ✅ No user intervention required for existing templates

### Default Template Categories
- ✅ "Simple Dev", "Frontend Dev", "TFE Development" → **Projects**
- ✅ "Monitoring Wall" → **Tools**

### Implementation Files
- ✅ `types.go`: Added `Category` field to `SessionTemplate` and `TemplateBuilder`
- ✅ `model.go`: Tree view rendering with `buildTemplateTreeItems()`, `updateTemplateTreeItems()`
- ✅ `update_keyboard.go`: Wizard updates for both create and save flows
- ✅ `view.go`: Wizard prompt updates and progress bar calculations
- ✅ `templates.go`: Migration function and categorized default templates

**Key Achievement:** Organized template management with collapsible categories, making it easy to organize and find templates as the template library grows

---

## Phase 9.1.2: Clipboard Paste & Smart Viewport ✅ COMPLETED

**Status:** Enhanced command input with paste support and intelligent scrolling

### Clipboard Paste Support
- ✅ **Ctrl+V Paste**: Paste clipboard content at cursor position
- ✅ **Multi-line Conversion**: Newlines automatically converted to spaces
- ✅ **Unicode Support**: Proper handling via `[]rune()` conversions
- ✅ **Large Paste Support**: No size limits, handles 5-10KB prompts
- ✅ **Visual Feedback**: Status message shows paste confirmation with character count
- ✅ **Error Handling**: Graceful failure if clipboard unavailable

### Multi-line Command Wrapping
- ✅ **Adaptive Wrapping**: Text wraps to fill available panel width
- ✅ **Panel Expansion**: More lines visible when header panel expands (accordion mode)
- ✅ **Cursor Tracking**: Cursor position maintained correctly across wrapped lines
- ✅ **Border Safety**: All lines guaranteed to fit within borders (no overflow)

### Smart Viewport System
- ✅ **Cursor Always Visible**: Viewport automatically centers on cursor position
- ✅ **Scroll Indicators**: Clear "↑ more above..." and "↓ more below..." indicators
- ✅ **Help Text Protected**: Help text always visible at bottom (space reserved)
- ✅ **Dynamic Calculation**: Adapts to panel height and accordion mode
- ✅ **Smart Space Allocation**: Reserves space for indicators when needed

### Implementation Details
- ✅ Dependency: `github.com/atotto/clipboard`
- ✅ `wrapCommandInput()`: Multi-line wrapping with cursor tracking (model.go:890-970)
- ✅ `updateHeaderContent()`: Viewport logic and scroll indicators (model.go:494-576)
- ✅ Safety truncation in panel rendering (view.go:440-445)

### Use Cases
- ✅ Paste large prompt templates from TFE (Terminal File Explorer)
- ✅ Paste multi-line git commit messages
- ✅ Paste Claude slash commands from documentation
- ✅ Paste complex scripts (auto-converted to single line)

**Key Achievement:** Seamless paste support for large prompt templates with intelligent viewport that keeps cursor visible and provides clear scroll feedback

---

## TFE Integration: CLI Flags ✅ COMPLETED

**Status:** Context-aware working directory support

### CLI Flags
- ✅ `--cwd <directory>`: Override template's working directory
- ✅ `--template <index>`: Create session from template and exit (no TUI)
- ✅ Combined usage: `tmuxplexer --cwd $PWD --template 0`

### Use Cases
- ✅ Launch templates in current directory context
- ✅ TFE context menu integration (launch from file browser)
- ✅ Automated session creation from scripts

### Implementation
- ✅ Flag parsing in main.go
- ✅ Override logic in session creation
- ✅ Backward compatible (no flag = use template's dir)

**Key Achievement:** Templates now context-aware, perfect for TFE integration

---

## Phase 10: Unified 3-Panel Adaptive Layout ✅ COMPLETED

**Status:** Complete refactor from 4-panel accordion to unified 3-panel vertical stack

### Layout Architecture
- ✅ **3-panel vertical stack**: Sessions (top) | Preview (middle) | Command (bottom)
- ✅ **Adaptive height distribution**: Panels resize based on focus state
  - Sessions focused: 50% / 30% / 20%
  - Preview focused: 30% / 50% / 20%
  - Command focused: Maintains previous upper panel sizing (no resize)
- ✅ Command panel always 20% (fixed for typing comfort)
- ✅ Smooth visual transitions when focus changes
- ✅ Upper panels don't resize when focusing command panel (prevents disorientation)

### Focus Management
- ✅ **Manual focus switching**:
  - Key `1`: Focus command panel
  - Key `2`: Focus sessions panel
  - Tab/Shift+Tab: Cycle through panels
- ✅ **Auto-focus behavior** (natural workflow):
  - Typing any character → auto-focus command panel
  - Arrow keys (↑↓) → auto-focus sessions panel
  - Scroll keys (PgUp/PgDn/Home/End/g/G) → auto-focus preview panel

### Mouse Interactions
- ✅ **Click-to-focus**: Click any panel to focus it
  - Clicking sessions → expands to 50%
  - Clicking preview → expands to 50%
  - Clicking command → focuses for typing
- ✅ **Mouse wheel scrolling**:
  - Scroll preview content when preview/command focused
  - Scroll sessions list when sessions focused
- ✅ **Y-coordinate detection**: Vertical stack uses Golden Rule #3

### Command Input Polish
- ✅ Multi-line command wrapping with cursor (█)
- ✅ Smart viewport with scroll indicators (↑ more above... / ↓ more below...)
- ✅ Cursor always visible in viewport
- ✅ Help text shows target session and controls
- ✅ Character count for long commands (>100 chars)
- ✅ Last command display when unfocused
- ✅ Clipboard paste support (Ctrl+V)

### Code Quality
- ✅ Removed all legacy 4-panel/tab layout code
- ✅ Deleted `view_tabs.go.bak`
- ✅ Removed commented legacy functions:
  - `calculateFourPanelLayout()` (model.go)
  - `renderFourPanelLayout()` (view.go)
  - `getPanelAtPosition()` (update_mouse.go)
  - `handleFourPanelClick()` (update_mouse.go)
- ✅ Clean codebase with no legacy references

### Implementation Files
- ✅ `types.go`: Focus state constants (FocusSessions, FocusPreview, FocusCommand)
- ✅ `model.go`: `calculateAdaptivePanelHeights()` - adaptive 40/40/20 → 50/30/20 logic
- ✅ `view.go`: `renderUnifiedView()` - 3-panel vertical stack rendering
- ✅ `update_keyboard.go`: Auto-focus behavior, focus cycling
- ✅ `update_mouse.go`: Click detection for adaptive panels, wheel scrolling

### User Experience
- ✅ Natural workflow: type → command, arrows → sessions, scroll → preview
- ✅ Visual feedback on focus changes (border color, panel expansion)
- ✅ No flicker or layout breaks
- ✅ Works in small terminals (60×15 minimum)
- ✅ Works in popup mode (`./tmuxplexer --popup`)

**Key Achievement:** Unified layout with intelligent auto-focus - user workflow drives panel focus

---

## Phase 10.1: Template Preview & Focus-Based Scrolling ✅ COMPLETED

**Status:** Enhanced template workflow and preview panel UX

### Template Preview in Middle Panel
- ✅ **Template details shown in preview panel** (not at bottom of list)
  - Shows: layout, category, description, pane configurations
  - Shows: working directories, commands, titles per pane
  - Shows: action hints (Enter to create, 'o' to attach, 'e' to edit, 'd' to delete)
- ✅ **Category preview**: Shows template count when category selected
- ✅ **Clean Templates tab**: No cramped details at bottom, just tree view
- ✅ **Better readability**: Full width preview with proper spacing

### Focus-Based Preview Scrolling
- ✅ **Preview scrolling respects focus**: Arrow keys only scroll preview when focused
  - Sessions/Templates focused (press `1`): Arrow keys navigate list
  - Preview focused (press `2`): Arrow keys scroll preview content
- ✅ **Mouse wheel scrolling**: Works when preview is focused
- ✅ **Scroll indicators**: Shows position and total lines when scrolling
- ✅ **Page navigation**: PgUp/PgDn, Home/End, g/G all work
- ✅ **Works for both**: Session previews AND template details

### Adaptive Sizing Enhancement
- ✅ **Command panel doesn't resize upper panels**: Prevents disorienting jumps
  - Before: Click command → panels 1 & 2 resize to 40/40 (jarring!)
  - After: Click command → panels 1 & 2 maintain previous sizing (smooth!)
- ✅ **lastUpperPanelFocus tracking**: Remembers last focus state (Sessions or Preview)
- ✅ **Resize only when switching 1↔2**: Panels only adapt when actually switching between Sessions and Preview

### Template Editing
- ✅ **Enabled 'e' key**: Edit templates.json in default editor
- ✅ **Editor detection**: $EDITOR → micro → nano → vim → vi
- ✅ **Auto-reload**: Templates refresh when editor closes
- ✅ **Help text updated**: Shows [e] Edit in Templates tab status bar

### Implementation
- ✅ `updateTemplatePreview()` (model.go): Renders template details in preview
- ✅ `updatePreviewContent()` (model.go): Routes to template preview when on Templates tab
- ✅ `moveUp()`/`moveDown()` (update_keyboard.go): Check focus before scrolling vs navigating
- ✅ `calculateAdaptivePanelHeights()` (model.go): Uses lastUpperPanelFocus for command panel
- ✅ Focus tracking in all focus-change operations (keyboard, mouse, tab cycling)

**Key Achievement:** Smooth, focus-aware UI that respects user intention - no surprise resizes or navigation conflicts

---

## Additional Features

### Directory and Git Branch Display
- ✅ Working directory shown for each session (📁)
- ✅ Home directory shortened to `~`
- ✅ Git branch detection via `git rev-parse --abbrev-ref HEAD`
- ✅ Non-git directories show directory only
- ✅ Updates with auto-refresh (2-second interval)

### Mouse Support
- ✅ Click any panel to focus
- ✅ Panel expansion on click (accordion mode)
- ✅ Mouse wheel scrolling (preview panel)

### Keyboard Shortcuts (Complete List)
| Key | Action |
|-----|--------|
| `1` | Focus command panel (bottom) |
| `2` | Focus sessions panel (top) |
| `3` | Reserved for future use |
| `Tab/Shift+Tab` | Cycle through panels |
| `↑/↓` or `k/j` | Navigate sessions (auto-focus sessions panel) |
| `←/→` or `h/l` | Navigate windows |
| `PgUp/PgDn` | Scroll preview (auto-focus preview panel) |
| `Home/End` or `g/G` | Jump to top/bottom of preview (auto-focus preview panel) |
| `Enter` | Attach to session / send command |
| `s` | Save session as template |
| `r` | Refresh preview or rename session |
| `d/D` | Kill session |
| `Ctrl+R` | Refresh sessions & Claude state |
| `Ctrl+V` | Paste clipboard (command mode only) |
| `Typing` | Auto-focus command panel and insert character |
| `ESC` | Exit command mode / close popup |
| `q/Ctrl+C` | Quit |

---

## Documentation

### Created Documentation Files
- ✅ `CLAUDE.md`: Project overview, architecture, usage guide
- ✅ `docs/claude-hooks-integration.md`: Claude Code integration details
- ✅ `docs/HOOKS-QUICKREF.md`: Quick reference for hooks system
- ✅ `docs/ARCHITECTURE_STATUS_DETECTION.md`: AI status detection architecture
- ✅ `docs/UNIFIED_CHAT_IMPLEMENTATION.md`: Command mode implementation details
- ✅ `README.md`: Basic project information

---

## Testing & Quality

### Test Commands
- ✅ `test_template`: Non-TTY template validation
- ✅ `test_create N`: Non-TTY session creation test
- ✅ Both commands bypass TUI for development testing

---

## Performance & Stability

### Auto-Refresh System
- ✅ 2-second ticker for live updates
- ✅ Updates session list, window list, pane preview
- ✅ Updates Claude Code state
- ✅ Minimal performance impact

### Error Handling
- ✅ Graceful handling of missing tmux
- ✅ Graceful handling of dead sessions
- ✅ Graceful handling of missing state files
- ✅ Stale state detection (>60s old)

---

## Current Production Status

**All core features working and stable:**
- 4-panel accordion layout with perfect alignment
- Complete template system (create, edit, delete, save from session)
- Real-time Claude Code integration
- Full pane preview with scrolling
- Popup mode with keybinding
- AI session command interface with filtering
- Context-aware template launching (--cwd flag)

**Ready for daily use!**
