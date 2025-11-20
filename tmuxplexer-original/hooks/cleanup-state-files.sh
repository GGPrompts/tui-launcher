#!/bin/bash
# Cleanup script for Claude Code state files
# Removes stale state files and old debug logs

STATE_DIR="/tmp/claude-code-state"
DEBUG_DIR="$STATE_DIR/debug"

if [[ ! -d "$STATE_DIR" ]]; then
    echo "No state directory found at $STATE_DIR"
    exit 0
fi

# Count files before cleanup
STATE_COUNT=$(find "$STATE_DIR" -maxdepth 1 -name "*.json" -type f 2>/dev/null | wc -l)
DEBUG_COUNT=$(find "$DEBUG_DIR" -name "*.json" -type f 2>/dev/null | wc -l)

echo "Before cleanup:"
echo "  State files: $STATE_COUNT"
echo "  Debug files: $DEBUG_COUNT"

# Clean up state files older than 7 days
echo ""
echo "Cleaning state files older than 7 days..."
REMOVED_STATE=$(find "$STATE_DIR" -maxdepth 1 -name "*.json" -type f -mtime +7 -delete -print 2>/dev/null | wc -l)

# Clean up debug files older than 1 hour (they accumulate VERY quickly - 400+ files/day!)
# Debug files are only useful for active debugging, no point keeping them longer
echo "Cleaning debug files older than 1 hour..."
REMOVED_DEBUG=$(find "$DEBUG_DIR" -name "*.json" -type f -mmin +60 -delete -print 2>/dev/null | wc -l)

# Count files after cleanup
STATE_COUNT_AFTER=$(find "$STATE_DIR" -maxdepth 1 -name "*.json" -type f 2>/dev/null | wc -l)
DEBUG_COUNT_AFTER=$(find "$DEBUG_DIR" -name "*.json" -type f 2>/dev/null | wc -l)

echo ""
echo "Cleanup complete!"
echo "  Removed $REMOVED_STATE state files"
echo "  Removed $REMOVED_DEBUG debug files"
echo ""
echo "After cleanup:"
echo "  State files: $STATE_COUNT_AFTER"
echo "  Debug files: $DEBUG_COUNT_AFTER"

# Show disk usage
DISK_USAGE=$(du -sh "$STATE_DIR" 2>/dev/null | cut -f1)
echo "  Total size: $DISK_USAGE"
