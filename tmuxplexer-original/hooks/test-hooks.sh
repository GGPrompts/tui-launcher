#!/bin/bash
# Test script for Claude Code hooks

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Claude Code Hooks Test Suite                             ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

STATE_DIR="/tmp/claude-code-state"
HOOKS_DIR="$HOME/.claude/hooks"

# Check if hook script exists
if [[ ! -f "$HOOKS_DIR/state-tracker.sh" ]]; then
    echo "❌ Error: state-tracker.sh not found at $HOOKS_DIR"
    echo "   Run ./hooks/install.sh first"
    exit 1
fi

# Create test directory
TEST_DIR="$STATE_DIR/test"
mkdir -p "$TEST_DIR"

echo "Testing hook script..."
echo ""

# Test 1: session-start
echo "Test 1: SessionStart hook"
export CLAUDE_SESSION_ID="test-session-001"
echo '{}' | "$HOOKS_DIR/state-tracker.sh" session-start
if [[ -f "$STATE_DIR/test-session-001.json" ]]; then
    echo "✓ State file created"
    STATUS=$(jq -r '.status' "$STATE_DIR/test-session-001.json")
    if [[ "$STATUS" == "idle" ]]; then
        echo "✓ Status is 'idle'"
    else
        echo "❌ Expected status 'idle', got '$STATUS'"
    fi
else
    echo "❌ State file not created"
fi
echo ""

# Test 2: user-prompt
echo "Test 2: UserPromptSubmit hook"
echo '{"prompt":"Test prompt"}' | "$HOOKS_DIR/state-tracker.sh" user-prompt
STATUS=$(jq -r '.status' "$STATE_DIR/test-session-001.json")
if [[ "$STATUS" == "processing" ]]; then
    echo "✓ Status changed to 'processing'"
else
    echo "❌ Expected status 'processing', got '$STATUS'"
fi
echo ""

# Test 3: pre-tool
echo "Test 3: PreToolUse hook"
echo '{"tool_name":"Edit","tool_input":{"file_path":"test.go"}}' | "$HOOKS_DIR/state-tracker.sh" pre-tool
STATUS=$(jq -r '.status' "$STATE_DIR/test-session-001.json")
TOOL=$(jq -r '.current_tool' "$STATE_DIR/test-session-001.json")
if [[ "$STATUS" == "tool_use" ]] && [[ "$TOOL" == "Edit" ]]; then
    echo "✓ Status is 'tool_use', tool is 'Edit'"
else
    echo "❌ Expected status 'tool_use' with tool 'Edit', got status='$STATUS' tool='$TOOL'"
fi
echo ""

# Test 4: post-tool
echo "Test 4: PostToolUse hook"
echo '{"tool_name":"Edit"}' | "$HOOKS_DIR/state-tracker.sh" post-tool
STATUS=$(jq -r '.status' "$STATE_DIR/test-session-001.json")
if [[ "$STATUS" == "working" ]]; then
    echo "✓ Status changed to 'working'"
else
    echo "❌ Expected status 'working', got '$STATUS'"
fi
echo ""

# Test 5: stop
echo "Test 5: Stop hook"
echo '{}' | "$HOOKS_DIR/state-tracker.sh" stop
STATUS=$(jq -r '.status' "$STATE_DIR/test-session-001.json")
if [[ "$STATUS" == "awaiting_input" ]]; then
    echo "✓ Status changed to 'awaiting_input'"
else
    echo "❌ Expected status 'awaiting_input', got '$STATUS'"
fi
echo ""

# Test 6: State file format
echo "Test 6: State file format validation"
STATE_FILE="$STATE_DIR/test-session-001.json"
if jq empty "$STATE_FILE" 2>/dev/null; then
    echo "✓ Valid JSON format"

    # Check required fields
    REQUIRED_FIELDS=("session_id" "status" "working_dir" "last_updated" "tmux_pane" "pid")
    ALL_PRESENT=true
    for field in "${REQUIRED_FIELDS[@]}"; do
        if ! jq -e ".$field" "$STATE_FILE" > /dev/null 2>&1; then
            echo "❌ Missing required field: $field"
            ALL_PRESENT=false
        fi
    done

    if $ALL_PRESENT; then
        echo "✓ All required fields present"
    fi
else
    echo "❌ Invalid JSON format"
fi
echo ""

# Test 7: Timestamp freshness
echo "Test 7: Timestamp validation"
TIMESTAMP=$(jq -r '.last_updated' "$STATE_FILE")
if date -d "$TIMESTAMP" > /dev/null 2>&1; then
    echo "✓ Valid ISO 8601 timestamp: $TIMESTAMP"
else
    echo "❌ Invalid timestamp format: $TIMESTAMP"
fi
echo ""

# Display final state
echo "Final state file contents:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
jq . "$STATE_FILE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Cleanup
echo "Cleaning up test files..."
rm -f "$STATE_DIR/test-session-001.json"
echo "✓ Test complete"
echo ""
echo "Next: Start a real Claude Code session to test integration:"
echo "  tmux new -s test-claude"
echo "  cd ~/projects/tmuxplexer"
echo "  claude"
