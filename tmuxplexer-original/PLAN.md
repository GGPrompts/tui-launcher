# Tmuxplexer - Project Plan

**Professional Tmux Session Manager with 4-Panel Layout + Workspace Templates + AI Orchestration**

---

## 📋 Current Status

**Production-Ready Features:**
- ✅ Phases 1-8: Core TUI, templates, preview, popup mode - **ALL COMPLETE**
- ✅ Phase 9.1: Unified chat/command mode for AI sessions - **COMPLETE**
- ✅ Phase 9.1.1: Template categorization & tree view - **COMPLETE**
- ✅ Phase 9.1.2: Clipboard paste & smart viewport - **COMPLETE**
- ✅ **Phase 10: Unified 3-Panel Adaptive Layout - COMPLETE**
  - Refactored from 4-panel to 3-panel vertical stack
  - Auto-focus behavior (type → command, arrows → sessions, scroll → preview)
  - Click-to-focus mouse support
  - Adaptive panel heights (50/30/20 based on focus)
- ✅ TFE Integration: `--cwd` and `--template` flags - **COMPLETE**

**See:** [docs/CHANGELOG.md](docs/CHANGELOG.md) for complete feature history

---

## 🎯 Next Up: Phase 9.2 - Multi-Selection System

**Goal:** Visual multi-selection for orchestrating commands across multiple AI sessions

**The Gap:** Currently you can send commands to ONE AI session at a time. Phase 9.2 adds:
- Select multiple sessions with checkboxes (Space to toggle)
- Send same command to all selected sessions
- Quick select shortcuts (Ctrl+A for all, Ctrl+C for all Claude, etc.)
- Perfect for multi-agent workflows

### Implementation Checklist

#### 9.2.1: Visual Selection with Checkboxes
- [ ] Add checkboxes to session list: `[✓]` or `[ ]`
- [ ] Space bar to toggle selection
- [ ] Visual feedback (highlight selected sessions)
- [ ] Selection count in header: "Sessions (3 selected)"
- [ ] Persist selection state in model

#### 9.2.2: Selection Modes
- [ ] **Visual mode (v)**: vim-like selection mode
  - Press `v` to enter, Space to select, `v` to exit
  - Move with j/k while in visual mode
- [ ] **Quick select shortcuts:**
  - `Ctrl+A` - Select all sessions
  - `Ctrl+D` - Deselect all
  - `Ctrl+I` - Select only Idle Claude sessions
  - `Ctrl+C` - Select only Claude sessions (any status)
  - `!` - Invert selection (in visual mode)

#### 9.2.3: Multi-Target Command Dispatch
- [ ] Modify command mode to detect selected sessions
- [ ] Send command to all selected sessions (not just current)
- [ ] Staggered execution (50ms delay between sends)
- [ ] Show per-target success/failure
- [ ] Aggregate results in status message
- [ ] Format: "Sent to 3 sessions: claude-1 ✓, claude-2 ✓, backend-api ✓"

#### 9.2.4: Selection Groups (Advanced)
- [ ] Save selection as named group
- [ ] Recall groups quickly (`:group <name>`)
- [ ] Store in `~/.config/tmuxplexer/selection_groups.json`
- [ ] Example groups:
  - "all-features" = feature-auth, feature-search, feature-export
  - "frontend" = webapp-dev, storybook, tests
  - "monitoring" = logs, metrics, alerts

### Use Case Example: Multi-Agent Development

**Scenario:** 3 Claude sessions working on different features

```
# 1. Select all feature sessions
Press 'v' (visual mode)
Navigate and Space to select:
  [✓] claude-auth (feature-auth branch)
  [✓] claude-search (feature-search branch)
  [✓] claude-export (feature-export branch)
Press 'v' to exit visual mode

# 2. Send coordinated commands
Press '1' (command mode)
Type: "Run tests and report status"
Press Enter
→ Command sent to all 3 sessions simultaneously

# 3. Monitor responses in preview
Press '4' and cycle through sessions
→ See each session's test output

# 4. Clear all sessions for next phase
Ensure all 3 still selected
Press '1', type: "/clear"
→ All sessions get fresh context
```

**Benefits:**
- Parallel development across multiple features
- Coordinated testing and deployment
- Bulk operations (git commands, clean builds, etc.)
- Phase-based workflows (research → implement → test → review)

---

## 🚀 Phase 10: Multi-Tool AI Status Detection

**Goal:** Real-time status detection for Aider, Codex, Gemini, and generic AI tools (beyond just Claude)

**Current State:** Claude Code status detection works perfectly via hooks system
**Gap:** No status detection for other AI coding assistants

**Consulted:** OpenAI Codex (GPT-5) - 2025-10-27
**Architecture:** See [docs/ARCHITECTURE_STATUS_DETECTION.md](docs/ARCHITECTURE_STATUS_DETECTION.md)

### Recommended Architecture (from Codex)

**Data Flow:**
```
tmux pane → stream collector → preprocessor → analyzers → aggregator → cache → TFE
```

**Why this approach:**
- ✅ Pipe-based streaming (low overhead, tmux-native)
- ✅ Preprocessing isolates terminal quirks (ANSI, CR, backspace)
- ✅ Analyzer plugins keep tool knowledge isolated
- ✅ Aggregation with TTLs prevents flicker
- ✅ Generic rule engine covers unknown tools

### Implementation Phases

#### 10.1: Foundation (Week 1)
- [ ] Implement core types (Status, Frame, Analyzer interface)
- [ ] Build ANSI stripper with CR/BS handling (critical for spinners)
- [ ] Create aggregator with TTL decay
- [ ] Write status.json output (`~/.cache/tmuxplexer/status.json`)

#### 10.2: Aider Analyzer (Week 1-2)
- [ ] Implement AiderAnalyzer (spinner detection)
- [ ] Build History helper for same-line changes
- [ ] Test with real Aider sessions
- [ ] Validate confidence scoring
- [ ] Status states: idle (prompt visible), active (spinner), editing

#### 10.3: Stream Collector (Week 2)
- [ ] Implement pipe-pane setup/teardown
- [ ] Per-pane log file management (`~/.cache/tmuxplexer/panes/<pane_id>.log`)
- [ ] Tail logs with line buffering (`stdbuf -oL`)
- [ ] Periodic snapshot for recovery
- [ ] Log rotation at 256 KB

#### 10.4: Generic Analyzers (Week 2-3)
- [ ] Implement CodexAnalyzer (heuristic-based)
- [ ] Implement GeminiAnalyzer (if detectable patterns exist)
- [ ] Build generic spinner/prompt detector
- [ ] Create rule-based analyzer (YAML configs)
- [ ] Test with multiple tools simultaneously

#### 10.5: TFE Integration (Week 3)
- [ ] TFE reads status.json
- [ ] Display AI badges next to directories:
  ```
  📁 website         [🤖 refactoring auth]
  📁 backend         [🤖×2 active]
  📁 tmuxplexer      [🤖 idle]
  ```
- [ ] Project path mapping (walk up to find .git)
- [ ] Badge formatting with tool icons

#### 10.6: Polish & Optimization (Week 3-4)
- [ ] Implement tmuxplexerd daemon (10-20 Hz writes)
- [ ] Add health metrics (panes watched, parse errors)
- [ ] Unit tests with canned transcripts
- [ ] Performance tuning (target: <5% CPU for 20 panes)
- [ ] Documentation for analyzer plugin API

### Tool Support Matrix

| Tool | Detection Method | Confidence | Status |
|------|-----------------|------------|--------|
| Claude Code | Hooks system (JSON state files) | 0.99 | ✅ Implemented |
| Aider | Spinner detection + prompt regex | 0.90 | 📋 Planned |
| Codex CLI | Process introspection + heuristics | 0.60 | 📋 Planned |
| Gemini | TBD (process or state file) | TBD | 📋 Research needed |
| Generic | Rule-based YAML configs | 0.50-0.70 | 📋 Planned |

### Extensibility

**Zero-Friction (No Code):**
```yaml
# ~/.config/tmuxplexer/rules.d/my-tool.yaml
process_names: [my-ai, my-assistant]
match_regex: "my-ai>"
active_regex: "Processing|Working"
idle_regex: "my-ai> $"
error_regex: "Error:|Failed:"
confidence: 0.7
ttl_ms: 2000
```

**High-Confidence (Small Analyzer):**
- Implement `Analyzer` interface (few dozen LOC)
- Register at startup
- Full control over detection logic

**Shared Status File Spec:**
```json
{
  "tools": [{
    "tool": "claude",
    "state": "active",
    "description": "reasoning",
    "pane_id": "%3",
    "project": "/home/matt/projects/website",
    "ttl_ms": 2000,
    "at": "2025-10-27T01:54:00Z"
  }]
}
```

Path: `~/.cache/ai-coding-status.json`

Encourage tool authors to write here for consistent status reporting.

---

## 🎨 Phase 11+: Future Vision

### Custom Monitoring Tools (Tmux-based MCPs)

**Insight:** Build simple TUI tools that display data - AI reads them via `tmux capture-pane`

**Built-in Tool Ecosystem:**

**Browser Tools:**
- `browser-console` - Live console errors/warnings (Chrome DevTools Protocol)
- `browser-network` - HTTP requests/responses
- `css-inspector` - CSS hierarchy and computed styles

**Development Tools:**
- `test-watcher` - Live test results across all projects
- `api-inspector` - API request/response viewer
- `log-aggregator` - Tails multiple log files in unified view

**System Tools:**
- `process-monitor` - Process tree with CPU/memory
- `git-activity` - Multi-repo git status dashboard
- `dependency-checker` - Outdated packages across projects

**Usage Pattern:**
```markdown
## Available Monitoring Tools (in CLAUDE.md)

Capture browser console:
```bash
tmux capture-pane -t browser-console -p
```
```

### Interactive Command Sequences

- [ ] Define command sequences in config
- [ ] Example: "deploy" = build → test → push → notify
- [ ] Execute with one keypress
- [ ] Show progress across sessions

### Session Health Monitoring

- [ ] Detect crashed processes (exit code ≠ 0)
- [ ] Detect frozen panes (no output in 1hr)
- [ ] Alert when Claude needs attention (⏸️ Awaiting Input)
- [ ] Visual indicators for session health

### Multi-Session Dashboard (Header Panel Enhancement)

```
Claude Sessions: 5 active
  🟢 2 Idle    🟡 1 Processing
  🔧 1 Working  ⏸️ 1 Needs Input

Next attention: blog-cms (awaiting input)
```

### Natural Language Interface

- [ ] Parse commands like:
  - "split horizontally"
  - "create session named dev"
  - "attach to backend"
  - "send 'git status' to all Claude sessions"
- [ ] Command suggestions (fuzzy matching)
- [ ] Help system with examples

### Remote Session Support

- [ ] Connect to remote tmux servers (SSH)
- [ ] Manage multiple hosts
- [ ] Host profiles

---

## 📊 Success Criteria

### Phase 9.2 Success Metrics
- ✅ Can select 3+ sessions visually with checkboxes
- ✅ Can send one command to all selected sessions
- ✅ Command dispatch completes in <500ms for 10 sessions
- ✅ Clear visual feedback for selection state
- ✅ Works seamlessly in both normal and popup modes

### Phase 10 Success Metrics
- ✅ Aider status detection accuracy >90%
- ✅ Codex status detection accuracy >60%
- ✅ Performance: <5% CPU for 20 panes
- ✅ Latency: Status updates within 500ms
- ✅ TFE displays AI badges correctly
- ✅ Extensible analyzer API documented

---

## 🏗️ Architecture Quick Reference

### Core Components

```
tmuxplexer/
├── main.go                 # Entry point, flag parsing
├── types.go                # Data structures
├── model.go                # Bubbletea model, layout calculations
├── update.go               # Message dispatcher
├── update_keyboard.go      # Keyboard controls
├── update_mouse.go         # Mouse handling
├── view.go                 # Rendering
├── styles.go               # Lipgloss styles
├── config.go               # Configuration
├── templates.go            # Template management
├── tmux.go                 # Tmux integration
├── claude_state.go         # Claude Code state reading
└── hooks/                  # Claude Code hooks
    ├── state-tracker.sh    # Hook script
    └── install.sh          # Hook installer
```

### Key Design Patterns

**Bubble Tea (Elm Architecture):**
- Model: Application state
- Update: Message handling and state transitions
- View: Rendering logic

**Panel System:**
- Weight-based layout calculations
- Dynamic sizing with accordion mode
- Independent panel content management

**Auto-Refresh:**
- 2-second ticker for live updates
- Updates sessions, windows, panes, Claude state

**AI Integration:**
- Hook-based for Claude (highest accuracy)
- Stream-based for other tools (Phase 10)
- Aggregation with confidence scoring

---

## 📝 Questions to Resolve

### Phase 9.2 Questions
- [ ] Should selection persist across panel switches?
- [ ] Should selection be visible in all modes or just when active?
- [ ] Max number of sessions to select simultaneously?
- [ ] Should selection groups sync across popup/normal modes?

### Phase 10 Questions
- [ ] How to handle Aider with custom themes (different spinner glyphs)?
- [ ] Should tmuxplexerd run as systemd service or manual launch?
- [ ] Fallback when stream collection fails?
- [ ] How to detect tool version changes that break analyzers?

---

## 🚀 Getting Started with Next Phase

**To start Phase 9.2:**
1. Read this plan and [docs/UNIFIED_CHAT_IMPLEMENTATION.md](docs/UNIFIED_CHAT_IMPLEMENTATION.md)
2. Review current command mode implementation (update_keyboard.go:1098-1137)
3. Add selection state to model (types.go)
4. Implement checkbox rendering in left panel (view.go)
5. Add Space key handler for toggle (update_keyboard.go)
6. Modify command execution to use selected sessions

**To start Phase 10:**
1. Read [docs/ARCHITECTURE_STATUS_DETECTION.md](docs/ARCHITECTURE_STATUS_DETECTION.md)
2. Study ANSI preprocessing requirements (critical for spinners)
3. Implement core types (Status, Frame, Analyzer interface)
4. Build ANSI stripper with test cases
5. Create AiderAnalyzer as proof of concept
6. Set up pipe-pane collection for one test pane

---

**Last Updated:** 2025-01-29
**Next Milestone:** Phase 9.2 (Multi-Selection) - Target: 1-2 weeks
**After That:** Phase 10 (Multi-Tool AI Detection) - Target: 3-4 weeks
