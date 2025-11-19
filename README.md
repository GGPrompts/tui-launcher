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

## Quick Start

```bash
# Clone
git clone https://github.com/YOUR_USERNAME/tui-launcher.git
cd tui-launcher

# Install dependencies
go mod tidy

# Build
go build

# Run (when complete)
./tui-launcher
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
- **↑/↓** - Move cursor
- **→** - Expand category
- **←** - Collapse category
- **Enter** - Launch item(s)

### Selection
- **Space** - Toggle selection
- **a** - Select all in category
- **c** - Clear selections
- **Esc** - Clear selections / close dialog

### Launching
- **Enter** - Launch (single or batch)
- **Ctrl+Enter** - Quick launch with default layout

## Development

See [PLAN.md](PLAN.md) for detailed architecture and roadmap.

## License

MIT
