# TUI Launcher

**A visual terminal launcher for managing projects, TUI tools, and batch command execution with tmux integration**

Built with Go, [Bubble Tea](https://github.com/charmbracelet/bubbletea), and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Status

🚧 **In Active Development** - Core architecture complete, UI implementation in progress

**Completed:**
- ✅ Type system and architecture
- ✅ Layout system with visual previews
- ✅ Spawn logic (tmux/xterm) using proven tmuxplexer pattern
- ✅ Sample configuration

**Next:**
- ⏳ Config loading
- ⏳ Tree view rendering (porting from TFE)
- ⏳ Keyboard navigation
- ⏳ Multi-select system

## Features

- 🌲 **Tree-based navigation** - Hierarchical organization of projects, tools, and commands
- ☑️  **Multi-select spawning** - Space to select, Enter to launch multiple items
- 📐 **Tmux integration** - Batch launches with configurable layouts (quad split, tiled, etc.)
- 🎯 **Context-aware** - Detects tmux environment and adapts
- 📦 **Project-based** - Set working directories per command
- 🔧 **Saved profiles** - Complex multi-pane setups in one command

## Installation

```bash
# Clone
git clone https://github.com/GGPrompts/tui-launcher.git
cd tui-launcher

# Install (builds, copies to ~/.local/bin, creates 'tl' wrapper)
./install.sh

# Add to PATH (if not already)
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

## Quick Start

```bash
# Launch from anywhere
tl

# Or use full command
tui-launcher
```

## Configuration

Create `~/.config/tui-launcher/config.yaml`:

```yaml
projects:
  - name: TFE
    icon: 🚀
    path: ~/projects/tfe
    commands:
      - name: TFE
        icon: 📂
        command: tfe
        spawn: tmux-split-h
      - name: Dev Server
        icon: 💻
        command: go run .
        spawn: tmux-split-v

tools:
  - category: System Monitoring
    icon: 📊
    items:
      - name: htop
        icon: 💹
        command: htop
        spawn: tmux-split-v
```

## Keyboard Shortcuts

### Navigation
- **↑/↓** or **j/k** - Move cursor (vim keys supported!)
- **→** or **l** - Expand category
- **←** or **h** - Collapse category
- **Tab** - Switch between panes (Global Tools ↔ Projects)
- **Mouse wheel** - Scroll through items

### Selection & Launch
- **Space** - Select/deselect items (multi-select)
- **Enter** - Launch selected item(s)
- **c** - Clear all selections

### Modes
- **d** - Toggle Foreground/Detached mode
  - **Foreground (default)**: Single commands run in terminal, multi-select spawns tmux windows and exits
  - **Detached**: Spawns tmux windows in background, launcher stays open
- **e** - Edit config file
- **i** - Toggle info pane (mobile mode)
- **q** or **Ctrl+C** - Quit

### Multi-Select Workflows

**Foreground Mode (default):**
1. Select multiple items with **Space**
2. Press **Enter** → Each spawns as a tmux window
3. Launcher exits, you're in tmux with multiple windows
4. Use **Ctrl+B w** to switch between windows

**Detached Mode (press 'd'):**
1. Select multiple items with **Space**
2. Press **Enter** → Each spawns as tmux window in background
3. Launcher stays open (spawn more if needed)
4. Press **2** to switch to Sessions tab
5. View all windows with live previews, attach to any

## Development

See [PLAN.md](PLAN.md) for detailed architecture and roadmap.

## License

MIT
