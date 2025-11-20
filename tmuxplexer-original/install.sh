#!/bin/bash
# Installation script for tmuxplexer tmux keybinding
# This adds Ctrl+b o to launch tmuxplexer in a popup

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Tmuxplexer - Tmux Keybinding Installation                ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check for required tools
echo "Checking dependencies..."

if ! command -v tmux &> /dev/null; then
    echo "❌ Error: tmux is not installed"
    echo "   Install with: sudo apt install tmux"
    exit 1
fi

echo "✓ tmux found"
echo ""

# Get absolute path to tmuxplexer binary
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TMUXPLEXER_BIN="$SCRIPT_DIR/tmuxplexer"

# Check if binary exists
if [[ ! -f "$TMUXPLEXER_BIN" ]]; then
    echo "⚠️  tmuxplexer binary not found at: $TMUXPLEXER_BIN"
    echo ""
    read -p "Would you like to build it now? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Building tmuxplexer..."
        cd "$SCRIPT_DIR"
        go build -o tmuxplexer
        echo "✓ Build complete"
    else
        echo "Please build tmuxplexer first:"
        echo "  cd $SCRIPT_DIR"
        echo "  go build -o tmuxplexer"
        exit 1
    fi
fi

echo "✓ tmuxplexer binary found"
echo ""

# Tmux config file
TMUX_CONF="$HOME/.tmux.conf"

# The keybinding to add
KEYBINDING="bind-key o run-shell \"tmux popup -E -w 80% -h 80% -d '#{pane_current_path}' $TMUXPLEXER_BIN --popup\""

# Check if tmux.conf exists
if [[ ! -f "$TMUX_CONF" ]]; then
    echo "Creating new ~/.tmux.conf..."
    echo "# Tmux configuration" > "$TMUX_CONF"
    echo "" >> "$TMUX_CONF"
fi

# Check if keybinding already exists
if grep -q "bind-key o.*tmuxplexer.*--popup" "$TMUX_CONF"; then
    echo "⚠️  Tmuxplexer keybinding already exists in ~/.tmux.conf"
    echo ""
    read -p "Would you like to update it? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Remove old binding
        sed -i.backup "/bind-key o.*tmuxplexer.*--popup/d" "$TMUX_CONF"
        echo "✓ Removed old keybinding"
    else
        echo "Skipping keybinding update"
        exit 0
    fi
fi

# Add keybinding
echo "" >> "$TMUX_CONF"
echo "# Tmuxplexer popup mode - Press Ctrl+b o" >> "$TMUX_CONF"
echo "$KEYBINDING" >> "$TMUX_CONF"
echo "" >> "$TMUX_CONF"

echo "✓ Added keybinding to ~/.tmux.conf"
echo ""

# Reload tmux config if inside tmux
if [[ -n "$TMUX" ]]; then
    echo "Reloading tmux configuration..."
    tmux source-file "$TMUX_CONF"
    echo "✓ Tmux config reloaded"
else
    echo "⚠️  Not inside tmux session"
    echo "   The keybinding will be available when you start/restart tmux"
fi

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Installation Complete!                                    ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "Usage:"
echo "  1. From any tmux session, press: Ctrl+b o"
echo "  2. Tmuxplexer will open in a popup (80% width/height)"
echo "  3. Select a session and press Enter to switch"
echo "  4. Or press 'q' to close the popup"
echo ""
echo "Keybinding added:"
echo "  $KEYBINDING"
echo ""
echo "To test:"
echo "  tmux"
echo "  Ctrl+b o"
echo ""
