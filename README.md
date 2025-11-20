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
- **Mouse wheel** - Scroll through items

### Selection
- **Space** - Context-aware: Expand category OR select command
- **c** - Clear all selections
- **Enter** - Launch selected item(s)

### Modes
- **t** - Toggle tmux mode (tmux spawning vs direct execution)
- **q** or **Ctrl+C** - Quit

### Multi-Select Launch
When multiple items selected:
1. Press **Enter** to open layout dialog
2. Use **↑/↓** to choose layout (quad split, tiled, etc.)
3. Press **Enter** to launch with selected layout

## Development

See [PLAN.md](PLAN.md) for detailed architecture and roadmap.

## License

MIT
