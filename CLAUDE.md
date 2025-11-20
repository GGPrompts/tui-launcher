# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TUI Launcher is a visual terminal launcher for managing projects, TUI tools, and batch command execution with tmux integration, built with Go and the Bubble Tea framework.

## Build and Run Commands

```bash
# Build the project
go build

# Run the launcher
./tui-launcher

# Install dependencies
go mod tidy

# Update the global binary after making changes
go build -o tui-launcher && cp tui-launcher ~/.local/bin/
```

**Note:** The `tl` wrapper at `~/.local/bin/tl` must be **sourced** in your shell config for the Quick CD feature to work:

```bash
# Add to ~/.bashrc or ~/.zshrc
source ~/.local/bin/tl
```

The wrapper function:
- Runs `~/.local/bin/tui-launcher`
- Checks for `~/.tui-launcher_cd_target` file
- Changes to that directory if present (Quick CD feature)

Always update the binary after making code changes to ensure the wrapper uses the latest version.

## Documentation

- **PLAN.md** - Current development roadmap and future features
- **CHANGELOG.md** - Version history and completed features
- **IMPLEMENTATION_PLAN.md** - Detailed 3-pane layout implementation guide
- **README.md** - Project overview and quick start

## Architecture

The codebase follows a Model-View-Update (MVU) pattern using Bubble Tea:

- **model.go**: Core application state and model initialization. Contains config loading logic using gopkg.in/yaml.v3
- **tree.go**: Tree building and rendering logic for hierarchical navigation
- **spawn.go**: Process spawning logic for tmux and terminal execution  
- **layouts.go**: Tmux layout definitions and application
- **types.go**: All type definitions and constants (no emoji variation selectors!)
- **main.go**: Entry point with minimal setup

## Key Implementation Details

### Config System
- Configuration loaded from `~/.config/tui-launcher/config.yaml`
- YAML structure supports projects, tools, AI commands, and scripts sections
- Each item can have commands with working directories and spawn modes
- Profiles support for multi-pane tmux configurations

### Tree Navigation
- Hierarchical tree structure with expandable/collapsible categories
- Uses tree connectors (`├─`, `└─`, `│`) for visual hierarchy
- Context-aware spacebar: expands categories OR selects commands
- Multi-select system with checkbox indicators (`☐`/`☑`)

### Spawn System (Simplified - TFE-inspired)
- **Foreground mode (default)**:
  - Single command: Runs in current terminal, launcher exits
  - Multi-select: Spawns each as tmux window, launcher exits
- **Detached mode (press 'd')**:
  - Single/multi: Spawns tmux windows with `-d` flag (background), launcher stays open
  - Perfect for React apps: spawn multiple named tmux windows, view in Sessions tab
- Uses `tea.ExecProcess` with `tea.ClearScreen` pattern (from TFE)
- Proper terminal state management
- Named tmux windows using item names

### Keyboard Controls
- Arrow keys or vim keys (h/j/k/l) for navigation
- Space for selection (multi-select)
- Enter to launch selected items
- 'd' to toggle between Foreground and Detached modes
- 'c' to clear all selections
- 'e' to edit config file (uses tea.ExecProcess, restarts on exit)
- Tab to switch panes (Global Tools ↔ Projects)
- 'i' to toggle info pane (mobile mode)