#!/bin/bash
# Installation script for Claude Code hooks integration with tmuxplexer

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Claude Code Hooks Installation for Tmuxplexer            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check for required tools
echo "Checking dependencies..."

if ! command -v jq &> /dev/null; then
    echo "❌ Error: jq is not installed"
    echo "   Install with: sudo apt install jq"
    exit 1
fi

if ! command -v tmux &> /dev/null; then
    echo "❌ Error: tmux is not installed"
    echo "   Install with: sudo apt install tmux"
    exit 1
fi

if ! command -v claude &> /dev/null; then
    echo "❌ Error: claude is not installed"
    echo "   Install from: https://claude.ai/download"
    exit 1
fi

echo "✓ All dependencies found"
echo ""

# Create Claude hooks directory
HOOKS_DIR="$HOME/.claude/hooks"
echo "Creating hooks directory: $HOOKS_DIR"
mkdir -p "$HOOKS_DIR"

# Copy state tracker script
echo "Installing state-tracker.sh..."
cp "$(dirname "$0")/state-tracker.sh" "$HOOKS_DIR/state-tracker.sh"
chmod +x "$HOOKS_DIR/state-tracker.sh"

# Create state directory
STATE_DIR="/tmp/claude-code-state"
echo "Creating state directory: $STATE_DIR"
mkdir -p "$STATE_DIR"

echo "✓ Files installed"
echo ""

# Check if settings.json exists
SETTINGS_FILE="$HOME/.claude/settings.json"
if [[ ! -f "$SETTINGS_FILE" ]]; then
    echo "Creating new settings.json..."
    cp "$(dirname "$0")/claude-settings-hooks.json" "$SETTINGS_FILE"
    echo "✓ Created $SETTINGS_FILE with hooks configuration"
else
    echo "⚠️  Existing settings.json found at: $SETTINGS_FILE"
    echo ""
    echo "You need to manually merge the hooks configuration:"
    echo "  1. Open: $SETTINGS_FILE"
    echo "  2. Add the 'hooks' section from: $(dirname "$0")/claude-settings-hooks.json"
    echo "  3. Or run: jq -s '.[0] * .[1]' $SETTINGS_FILE $(dirname "$0")/claude-settings-hooks.json > /tmp/merged.json && mv /tmp/merged.json $SETTINGS_FILE"
    echo ""

    read -p "Would you like to automatically merge? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Backup existing settings
        cp "$SETTINGS_FILE" "$SETTINGS_FILE.backup.$(date +%Y%m%d_%H%M%S)"
        echo "✓ Backed up existing settings"

        # Merge settings
        jq -s '.[0] * .[1]' "$SETTINGS_FILE" "$(dirname "$0")/claude-settings-hooks.json" > /tmp/merged.json
        mv /tmp/merged.json "$SETTINGS_FILE"
        echo "✓ Merged hooks configuration into settings.json"
    fi
fi

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Installation Complete!                                    ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "Next steps:"
echo "  1. Test the hooks:"
echo "     ./hooks/test-hooks.sh"
echo ""
echo "  2. Start Claude Code in a tmux session:"
echo "     tmux new -s test-claude"
echo "     cd ~/projects/tmuxplexer"
echo "     claude"
echo ""
echo "  3. In another pane, monitor state:"
echo "     watch -n 0.5 'cat /tmp/claude-code-state/*.json | jq .'"
echo ""
echo "  4. Build and run tmuxplexer:"
echo "     go build"
echo "     ./tmuxplexer"
echo ""
echo "📖 Full documentation: ./docs/claude-hooks-integration.md"
