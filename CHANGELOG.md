# Changelog

All notable changes to TUI Launcher will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

## [0.3.0] - 2025-11-20 - Launch Tab Refinement

### Added
- **Foreground/Detached mode toggle** - Press 'd' to switch between modes
  - **Foreground mode (default)**: Commands run in current terminal, launcher exits
  - **Detached mode**: Commands spawn as background tmux windows, launcher stays open
- **Mode indicator in header** - Shows current mode (Foreground/Detached)
- **Multi-select tmux spawning**:
  - Foreground mode: Spawns each item as tmux window, exits launcher
  - Detached mode: Spawns each with `-d` flag (background), stays in launcher
- **TFE-inspired command execution** - Uses `tea.ExecProcess` with `tea.ClearScreen` pattern
- **Debug logging** - Stderr output for troubleshooting spawn operations

### Changed
- **Simplified spawn logic** - Removed complex tmux/direct mode toggle ('t' key)
- **Single command foreground launch** - Always runs in current terminal by default
- **Edit config ('e' key)** - Now uses `tea.ExecProcess` for reliable editor launching
- **Multi-select behavior** - Creates tmux windows instead of using layout dialog
- **Config spawn field** - Now ignored for single commands (use 'd' toggle instead)

### Fixed
- **Edit config not working** - Fixed by using proper `tea.ExecProcess` pattern
- **Commands spawning in foreground when expecting detached** - Added `-d` flag to tmux commands
- **Multi-select only launching first item** - Now spawns all selected items properly
- **Terminal state after external programs** - Proper cleanup with `tea.ClearScreen`

### Technical Details
- Launch tab fully migrated to `tabs/launch/` package
- Unified model architecture with tab routing
- Removed `useTmux` field, replaced with `detachedMode`
- Error handling continues through all items (doesn't stop on first error)

## [0.2.0] - 2025-11-19 - Responsive Layout & Quick CD

### Added
- **3-pane responsive layout system** with adaptive UI modes:
  - **Desktop mode** (≥80 width, >12 height): Global Tools (left) | Projects (right) | Info (bottom)
  - **Compact mode** (<80 width): Combined tree + info for narrow terminals
  - **Mobile mode** (≤12 height): Single pane with 'i' key to toggle info
  - Layout mode indicator in header for debugging
- **Multi-pane navigation**:
  - Tab key switches between panes (desktop) or toggles Global/Projects (compact/mobile)
  - Independent cursor tracking for each pane
  - Arrow keys navigate within active pane
  - Visual indicators show which pane is focused
- **Dynamic info pane** that shows:
  - Item name, type, and icon
  - Command details and working directory
  - Profile layout and pane information
  - Category children count
  - Project directory path with CD hint
- **Quick CD feature** (TFE-style):
  - Press Enter on a project to CD into that directory
  - Wrapper script handles directory change after exit
  - Info pane shows project directory and "Press Enter to CD" hint
  - Multi-select still works: select commands from multiple projects, then launch
- **Mouse support**:
  - Mouse wheel scrolls within active pane
  - Click left/right pane to switch focus in desktop mode
  - Follows Bubbletea Golden Rule #3 (X-coordinate for horizontal splits)
- **Bubbletea skill integration** for TUI best practices:
  - Golden Rule #1: Always account for borders before rendering
  - Golden Rule #2: Explicitly truncate text to prevent wrapping
  - Golden Rule #3: Match mouse detection to layout orientation
  - Golden Rule #4: Use proportional sizing (weights), not hardcoded pixels

### Changed
- **Spawn behavior now smarter (foreground by default)**:
  - Commands without explicit `spawn:` mode now use current pane/terminal (foreground)
  - Outside tmux: Replaces current terminal instead of spawning xterm window
  - Inside tmux: Uses current pane instead of creating split
  - Explicit spawn modes (tmux-window, tmux-split-h, etc.) still work as configured
  - You can always relaunch with `tl` - no need for TUI to stay resident
- Footer is now static (no auto-scrolling) to prevent flashing
- Footer text shortened to fit better on small screens
- Tree building refactored to split items into global/project panes
- Projects now store directory path for CD functionality

### Fixed
- Command execution now uses proper shell invocation (`sh -c`)
- Edit config (E key) now works correctly:
  - Opens in tmux split when inside tmux
  - Opens in terminal with proper TTY restoration when outside tmux
  - Prioritizes nano/vim for better compatibility
- Enter key now launches items without requiring Space selection first
- Working directory paths now properly quoted to handle spaces
- Top panels no longer extend off screen (proper height calculation)
- Footer flashing eliminated by removing auto-scroll
- Mouse wheel now works correctly with multi-pane layout

### Technical Details
- Height calculation properly accounts for: header (3 lines) + footer (2 lines) + borders (2 lines)
- Text truncation prevents wrapping and overflow
- Proportional pane sizing adapts to terminal width
- Wrapper script pattern: `~/.tui-launcher_cd_target` for CD communication

## [0.1.0] - 2025-01-19 - MVP Release

### Added

#### Phase 1: Core Tree View
- Hierarchical tree structure with expandable/collapsible categories
- Tree connectors (`├─`, `└─`, `│`) for visual hierarchy
- Emoji icons for visual identification (no variation selectors)
- Multi-level nesting (Projects → Commands, Categories → Tools)
- Keyboard navigation (arrow keys, vim keys h/j/k/l)
- Mouse wheel scrolling support
- Working directory display in header
- Bold text styling for cursor selection using Lipgloss
- YAML config loading with gopkg.in/yaml.v3

#### Phase 2: Multi-Select System
- Space key to toggle selection (context-aware)
- Checkbox rendering (☐/☑) for selected items
- Selection count displayed in status bar
- Clear all selections with 'c' key
- Visual feedback for selected items

#### Phase 3: Spawn Logic
- Single command launch with Enter key
- Batch multi-select launch with configurable tmux layouts:
  - main-vertical, main-horizontal, tiled, even-horizontal, even-vertical
- Tmux integration (splits, windows, sessions)
- Profile support for multi-pane configurations
- Working directory handling per command with `cwd` field
- Toggle mode between Tmux and Direct execution ('t' key)
- Auto-clear selections after successful launch
- 10ms delay between pane creation for stability (tmuxplexer pattern)
- Auto-detection of tmux environment

#### Configuration System
- YAML-based config at `~/.config/tui-launcher/config.yaml`
- Support for projects, tools, AI commands, and scripts sections
- Per-item spawn mode configuration
- Working directory (`cwd`) and path (`path`) support
- Profile configurations for complex multi-pane setups
- Full config with 30+ real tools organized by category

#### UI Polish
- Animated footer scrolling (unicode-safe)
- Context-aware spacebar: expands categories OR selects commands
- Clean spawn mode display (hidden from tree view)
- Proper unicode handling for smooth scrolling
- Responsive layout adjustments

#### Developer Experience
- Bash wrapper script at `~/.local/bin/tl`
- Update command: `go build -o tui-launcher && cp tui-launcher ~/.local/bin/`
- Clear project structure following MVU pattern:
  - `main.go` - Entry point
  - `types.go` - Type definitions
  - `model.go` - State and initialization
  - `tree.go` - Tree building and rendering
  - `spawn.go` - Process spawning logic
  - `layouts.go` - Tmux layout definitions

### Architecture
- Model-View-Update (MVU) pattern using Bubble Tea
- Separation of global tools and project-specific commands
- Proven patterns from TFE (tree view) and tmuxplexer (spawn logic)
- Terminal type detection (Windows Terminal vs others)
- Termux optimization for mobile terminals

### Keyboard Shortcuts
- **Navigation**: ↑/↓ or j/k - Move cursor
- **Navigation**: → or l - Expand category
- **Navigation**: ← or h - Collapse category
- **Selection**: Space - Toggle selection (or expand category)
- **Selection**: c - Clear all selections
- **Launch**: Enter - Launch selected items or current item
- **Modes**: t - Toggle between Tmux and Direct execution
- **Config**: e - Edit config file
- **Exit**: q or Ctrl+C - Quit application

### Notes
- Based on TUITemplate architecture
- Reuses proven patterns from TFE (tree view)
- Reuses tmuxplexer spawn patterns
- Unicode handling optimized for smooth scrolling
- Ready for daily use
