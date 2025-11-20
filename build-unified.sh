#!/bin/bash
# build-unified.sh - Build script for unified tab-based tui-launcher
# Phase 1: Launch tab integration

echo "Building tui-launcher with unified tab architecture..."

# Use standard go build (picks up all .go files and packages)
go build -o tui-launcher

if [ $? -eq 0 ]; then
    echo "✓ Build successful: ./tui-launcher"
    echo ""
    echo "Installing to ~/.local/bin/..."
    cp tui-launcher ~/.local/bin/
    echo "✓ Installed to ~/.local/bin/tui-launcher"
    echo ""
    echo "Launch tab is now integrated!"
    echo "Press 't' to toggle between Tmux and Direct modes"
    echo "Press 'e' to edit config"
    echo ""
    echo "Run: tui-launcher (or ./tui-launcher)"
else
    echo "✗ Build failed"
    exit 1
fi
