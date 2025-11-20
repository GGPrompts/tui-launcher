# Claude Code Hooks for Tmuxplexer

This directory contains the hook system that enables real-time communication between Claude Code sessions and tmuxplexer.

## Quick Start

1. **Install hooks**:
   ```bash
   ./hooks/install.sh
   ```

2. **Test hooks**:
   ```bash
   ./hooks/test-hooks.sh
   ```

3. **Start Claude in tmux**:
   ```bash
   tmux new -s test-claude
   cd ~/projects/tmuxplexer
   claude
   ```

4. **Monitor state changes** (in another pane):
   ```bash
   watch -n 0.5 'cat /tmp/claude-code-state/*.json | jq .'
   ```

5. **Run tmuxplexer**:
   ```bash
   go build
   ./tmuxplexer
   ```

## Files

- **state-tracker.sh**: Main hook script that tracks Claude's state
- **install.sh**: Installation script for setting up hooks
- **test-hooks.sh**: Test suite for validating hook functionality
- **claude-settings-hooks.json**: Hook configuration to add to ~/.claude/settings.json
- **README.md**: This file

## How It Works

### State Tracking Flow

```
Claude Code Event → Hook Fires → state-tracker.sh → State File → Tmuxplexer Reads
```

### Available Hooks

- **SessionStart**: Claude starts up
- **UserPromptSubmit**: User sends a message
- **PreToolUse**: Before Claude uses a tool (Edit, Read, Bash, etc.)
- **PostToolUse**: After tool execution completes
- **Stop**: Claude finishes responding
- **Notification**: System notifications (including awaiting-input)

### State Transitions

```
idle → processing → tool_use → working → awaiting_input → (cycle repeats)
```

### State File Format

Located at: `/tmp/claude-code-state/<session-id>.json`

```json
{
  "session_id": "abc123",
  "status": "tool_use",
  "current_tool": "Edit",
  "working_dir": "/home/user/project",
  "last_updated": "2025-08-24T14:30:22Z",
  "tmux_pane": "%42",
  "pid": 12345,
  "hook_type": "pre-tool",
  "details": {
    "event": "tool_starting",
    "tool": "Edit",
    "args": {...}
  }
}
```

## Status Indicators

When viewed in tmuxplexer:

- 🟢 **Idle**: Just started, waiting for input
- 🟡 **Processing**: Thinking about user's request
- 🔧 **Tool Use**: Actively executing a tool
- ⚙️ **Working**: Processing tool results
- ⏸️ **Awaiting Input**: Waiting for user response
- ⚪ **Stale**: No updates in >5 seconds (may have crashed)

## Troubleshooting

### Hook not firing

```bash
# Check hooks are configured
cat ~/.claude/settings.json | jq '.hooks'

# Test manually
echo '{}' | ~/.claude/hooks/state-tracker.sh session-start

# Enable debug mode
claude --debug hooks
```

### State file not updating

```bash
# Check state directory
ls -la /tmp/claude-code-state/

# Monitor in real-time
tail -f /tmp/claude-code-state/*.json
```

### Tmuxplexer not showing status

```bash
# Verify Claude session detected
tmux list-sessions -F "#{session_name}|#{pane_current_command}"

# Check state file exists
cat /tmp/claude-code-state/*.json | jq .
```

## Performance

- Hook execution: ~5-10ms overhead per hook
- State file write: ~1ms
- Tmuxplexer read: Only on refresh (2 second interval)
- State files auto-cleanup: After 24 hours

## Security

- State files are world-readable in /tmp (no secrets stored)
- Hooks run with same permissions as Claude Code
- Timeouts prevent hanging Claude
- Failed hooks don't crash Claude

## Advanced Usage

### Custom State Directory

Edit `state-tracker.sh` and change:
```bash
STATE_DIR="/tmp/claude-code-state"
```

### Adjust Stale Threshold

Edit `claude_state.go` and change:
```go
staleThreshold = 5 * time.Second
```

### Add Custom State Events

Extend the switch statement in `state-tracker.sh`:
```bash
case "$HOOK_TYPE" in
    my-custom-event)
        STATUS="custom"
        DETAILS='{"event":"custom_event"}'
        ;;
esac
```

## Documentation

Full documentation: [../docs/claude-hooks-integration.md](../docs/claude-hooks-integration.md)

## Support

Issues? Check:
1. [Troubleshooting section in main docs](../docs/claude-hooks-integration.md#troubleshooting)
2. Run test suite: `./hooks/test-hooks.sh`
3. Check Claude Code logs: `claude --debug hooks`
