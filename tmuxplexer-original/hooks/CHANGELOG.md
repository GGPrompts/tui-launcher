# Hooks Changelog

## 2025-11-10 - Session Isolation Fix

### Problem
Multiple Claude Code sessions were incorrectly sharing state updates. When one Claude session was "working", idle Claude sessions in other tmux panes would also show "working" status.

### Root Causes
1. **Missing Matchers**: Hooks configuration lacked matchers, causing hooks to trigger for ALL events instead of being filtered by notification type
2. **Session ID Priority**: State tracker relied on `TMUX_PANE` instead of Claude's official `session_id` from hook data

### Changes

#### 1. Added Matchers to Hooks Configuration (`claude-settings-hooks.json`)
- **All non-notification hooks**: Added `"matcher": "*"` (match all tools/events)
- **Notification hooks**: Added `"matcher": "idle_prompt"` (only trigger on idle/awaiting-input notifications)

This prevents notification hooks from firing for every notification type (permission prompts, auth events, etc.).

#### 2. Improved Session Identification (`state-tracker.sh`)
**New Priority Order**:
1. `session_id` from stdin hook data (most reliable, provided by Claude Code)
2. `CLAUDE_SESSION_ID` environment variable
3. `TMUX_PANE` (fallback for tmux sessions)
4. Process PID (fallback for non-tmux)

**Why This Matters**: Claude Code's `session_id` is guaranteed unique per session, whereas environment variables could potentially leak across processes.

#### 3. Enhanced Notification Handling
Added explicit handling for notification types:
- `idle_prompt` / `awaiting-input`: Sets status to "awaiting_input"
- `permission_prompt`: Preserves current status (doesn't interfere with workflow)
- Other notifications: Preserves current status

### How to Apply

**Option 1: Full Reinstall**
```bash
cd ~/projects/tmuxplexer
./hooks/install.sh
```

**Option 2: Manual Update**
```bash
# 1. Update the state tracker script
cp hooks/state-tracker.sh ~/.claude/hooks/state-tracker.sh
chmod +x ~/.claude/hooks/state-tracker.sh

# 2. Merge hooks configuration into ~/.claude/settings.json
# Add matchers as shown in hooks/claude-settings-hooks.json
```

**Option 3: Quick Merge (with backup)**
```bash
# Backup existing settings
cp ~/.claude/settings.json ~/.claude/settings.json.backup.$(date +%Y%m%d_%H%M%S)

# Merge new hooks config
jq -s '.[0] * .[1]' ~/.claude/settings.json hooks/claude-settings-hooks.json > /tmp/merged.json
mv /tmp/merged.json ~/.claude/settings.json
```

### Verification

1. **Restart all Claude Code sessions** (close and reopen)
2. **Open multiple Claude sessions** in different tmux panes
3. **Send a prompt to one Claude** and verify:
   - Only that Claude's state changes to "processing" → "working"
   - Other idle Claudes remain at "idle" status
4. **Check state files**:
   ```bash
   ls -la /tmp/claude-code-state/
   # Should see one .json file per Claude session

   cat /tmp/claude-code-state/*.json | jq .
   # Each should have unique session_id
   ```

### Expected Behavior After Fix

- ✅ Each Claude Code session has its own isolated state file
- ✅ Hooks only trigger for the specific session that generated the event
- ✅ Notification hooks only fire for `idle_prompt` events (not permission/auth/etc)
- ✅ State files use Claude's official `session_id` for perfect isolation
- ✅ Idle Claude sessions remain "idle" when other sessions are working

### Migration Notes

**No breaking changes** - existing state files will continue to work. The new session identification logic gracefully falls back to `TMUX_PANE` if `session_id` isn't available from hook data.

**State file cleanup**: Old state files with `TMUX_PANE`-based names will naturally expire (24-hour auto-cleanup) and be replaced with `session_id`-based names.
