#!/bin/bash
# build-unified.sh - Build script for unified tab-based tui-launcher
# Phase 1: Launch tab integration

echo "Building tui-launcher with unified tab architecture..."

go build -o tui-launcher \
    types.go \
    types_unified.go \
    model_unified.go \
    tab_routing.go \
    main_test_tabs.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful: ./tui-launcher"
    echo ""
    echo "Launch tab is now integrated!"
    echo "Press 1 to see the Launch tab (actual tui-launcher interface)"
    echo "Press 2/3 to see Sessions/Templates placeholders"
    echo ""
    echo "Run: ./tui-launcher"
else
    echo "✗ Build failed"
    exit 1
fi
