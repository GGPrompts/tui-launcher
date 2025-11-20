# Claude Code Hooks - Quick Reference

## Installation (One-time)

```bash
cd ~/projects/tmuxplexer
./hooks/install.sh
```

## Testing

```bash
# Run test suite
./hooks/test-hooks.sh

# Watch state changes live
watch -n 0.5 'cat /tmp/claude-code-state/*.json | jq .'
```

## Available Hooks

| Hook | When It Fires | Status Change | Use Case |
|------|---------------|---------------|----------|
| `SessionStart` | Claude starts | → `idle` | Initialize tracking |
| `UserPromptSubmit` | User sends message | → `processing` | User is active |
| `PreToolUse` | Before tool execution | → `tool_use` | Claude working |
| `PostToolUse` | After tool completes | → `working` | Processing results |
| `Stop` | Claude finishes | → `awaiting_input` | Waiting for user |
| `Notification` | System notification | varies | Bell/alerts |

## Status Indicators

| Icon | Status | Meaning |
|------|--------|---------|
| 🟢 | idle | Just started |
| 🟡 | processing | Thinking |
| 🔧 | tool_use | Using Edit/Read/Bash |
| ⚙️ | working | Processing results |
| ⏸️ | awaiting_input | Waiting for you |
| ⚪ | stale | No updates (>5s) |

## State File Location

```
/tmp/claude-code-state/<session-id>.json
```

## Hook Script Location

```
~/.claude/hooks/state-tracker.sh
```

## Configuration Location

```
~/.claude/settings.json
```

## Quick Commands

```bash
# List all state files
ls -lh /tmp/claude-code-state/

# View state for specific session
cat /tmp/claude-code-state/abc123.json | jq .

# Monitor all sessions
watch -n 1 'for f in /tmp/claude-code-state/*.json; do echo "=== $f ==="; jq . $f; done'

# Count active Claude sessions
ls /tmp/claude-code-state/*.json 2>/dev/null | wc -l

# Cleanup stale state files
find /tmp/claude-code-state -name "*.json" -mtime +1 -delete

# Test hook manually
echo '{"test":"data"}' | ~/.claude/hooks/state-tracker.sh session-start

# Check Claude settings for hooks
cat ~/.claude/settings.json | jq '.hooks'

# Debug hooks in Claude
claude --debug hooks
```

## State Lifecycle

```
SessionStart
     ↓
   [idle]
     ↓
UserPromptSubmit
     ↓
 [processing]
     ↓
 PreToolUse
     ↓
 [tool_use]
     ↓
PostToolUse
     ↓
 [working]
     ↓
    Stop
     ↓
[awaiting_input]
     ↓
(cycle repeats)
```

## Hook Data Format

### Input (stdin to hook script)

```json
{
  "tool_name": "Edit",
  "tool_input": {
    "file_path": "/path/to/file.go"
  },
  "prompt": "User's prompt text",
  "notification_type": "awaiting-input"
}
```

### Output (state file)

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
    "tool": "Edit"
  }
}
```

## Tmuxplexer Integration

### Go Types

```go
type ClaudeState struct {
    SessionID   string
    Status      string  // idle, processing, tool_use, working, awaiting_input
    CurrentTool string
    WorkingDir  string
    LastUpdated string
    TmuxPane    string
    PID         int
    Details     map[string]interface{}
}
```

### Detection

```go
// Check if session runs Claude
detectClaudeSession(sessionName) bool

// Get state for session
getClaudeStateForSession(sessionName, paneID) (*ClaudeState, error)

// Format for display
formatClaudeStatus(state) string       // "🟢 Idle"
getClaudeStatusIcon(state) string      // "🟢"
```

## Environment Variables

| Variable | Purpose | Example |
|----------|---------|---------|
| `CLAUDE_SESSION_ID` | Override session ID | `export CLAUDE_SESSION_ID=my-session` |
| `TMUX_PANE` | Auto-detected pane ID | `%42` |
| `PWD` | Working directory | `/home/user/project` |

## Performance

- Hook execution: ~5-10ms
- State file write: ~1ms
- No network calls
- Auto-cleanup after 24h
- Stale after 5s

## Troubleshooting Checklist

- [ ] Is jq installed? (`which jq`)
- [ ] Are hooks configured? (`cat ~/.claude/settings.json | jq .hooks`)
- [ ] Is script executable? (`ls -l ~/.claude/hooks/state-tracker.sh`)
- [ ] Does state dir exist? (`ls -la /tmp/claude-code-state/`)
- [ ] Are state files being created? (`watch ls /tmp/claude-code-state/`)
- [ ] Is timestamp fresh? (`cat /tmp/claude-code-state/*.json | jq .last_updated`)
- [ ] Is Claude detecting session? (Check tmuxplexer display)

## Common Issues

### "jq: command not found"
```bash
sudo apt install jq
```

### "Permission denied: state-tracker.sh"
```bash
chmod +x ~/.claude/hooks/state-tracker.sh
```

### "No state file created"
```bash
# Check hooks are configured
cat ~/.claude/settings.json | jq .hooks

# Test manually
echo '{}' | ~/.claude/hooks/state-tracker.sh session-start
ls /tmp/claude-code-state/
```

### "Status always shows 'Stale'"
```bash
# Check timestamp is recent
cat /tmp/claude-code-state/*.json | jq .last_updated

# Verify hooks are firing
claude --debug hooks
```

## Resources

- Full docs: [docs/claude-hooks-integration.md](./claude-hooks-integration.md)
- Hook README: [hooks/README.md](../hooks/README.md)
- Installation: [hooks/install.sh](../hooks/install.sh)
- Tests: [hooks/test-hooks.sh](../hooks/test-hooks.sh)
