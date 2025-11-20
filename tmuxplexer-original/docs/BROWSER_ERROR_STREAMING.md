# Browser Error Streaming to Terminal

**Status**: Planning / Proof of Concept
**Purpose**: Enable automatic browser error/console streaming to tmux panes for seamless Claude-assisted debugging

## Problem Statement

When developing frontend applications with Claude Code assistance, there's a constant manual workflow:

1. Code change is made
2. Browser shows error
3. Developer manually copies error from browser console
4. Developer pastes error to Claude
5. Claude provides fix
6. Repeat

This breaks flow state and adds significant friction to the debugging loop.

## Proposed Solution

Create a pipeline that automatically streams browser errors, console logs, and network issues directly to a tmux pane, which Claude can read via `tmux capture-pane`.

### Benefits

- **Zero Copy/Paste**: Errors automatically visible to Claude
- **Contextual**: Errors appear with timestamps, source locations, stack traces
- **Real-Time**: Claude sees errors as they happen, can track patterns
- **Full History**: Scrollback buffer preserves entire debugging session
- **Multi-Modal**: Captures console logs, network errors, CSS issues, performance metrics
- **Tmuxplexer Integration**: Natural fit with existing tmux workflow

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Browser Context                          │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │         Browser Extension (Manifest V3)                   │  │
│  │  - Intercepts console.error/warn/log                      │  │
│  │  - Captures uncaught exceptions                           │  │
│  │  - Monitors network requests (fetch/XHR)                  │  │
│  │  - Tracks React/Vue errors (via window.onerror)           │  │
│  │  - Observes CSS/layout issues (optional)                  │  │
│  └────────────────┬──────────────────────────────────────────┘  │
│                   │ chrome.runtime.sendNativeMessage()          │
└───────────────────┼──────────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Native Bridge Layer                         │
│                                                                   │
│  Option A: WebSocket Server (localhost:8765)                    │
│  - Browser extension connects via WebSocket                     │
│  - Real-time bidirectional communication                        │
│  - Can send commands back (e.g., "clear errors")                │
│                                                                   │
│  Option B: Native Messaging Host                                │
│  - Chrome native messaging protocol                             │
│  - More secure, no network port needed                          │
│  - Requires native app manifest registration                    │
│                                                                   │
│  Option C: File Watch (Simplest)                                │
│  - Extension writes to /tmp/browser-errors.jsonl                │
│  - Terminal watches with tail -f                                │
│  - No server needed, minimal setup                              │
└────────────────┬─────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Terminal (Tmux Pane)                        │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Pane: Browser Error Monitor                            │   │
│  │                                                          │   │
│  │  [14:32:15] ❌ TypeError: Cannot read 'user' of null    │   │
│  │              at LoginForm.tsx:45                         │   │
│  │              Stack: LoginForm > App > root              │   │
│  │                                                          │   │
│  │  [14:32:16] 🌐 POST /api/login → 401 Unauthorized       │   │
│  │              Response: {"error": "Invalid token"}       │   │
│  │                                                          │   │
│  │  [14:32:20] ⚠️  Warning: useEffect missing dependency  │   │
│  │              at Dashboard.tsx:89                         │   │
│  └─────────────────────────────────────────────────────────┘   │
└────────────────┬─────────────────────────────────────────────────┘
                 │
                 ▼ tmux capture-pane -p -t {session}:{window}.{pane}
┌─────────────────────────────────────────────────────────────────┐
│                      Claude Code Context                         │
│                                                                   │
│  Claude can now:                                                 │
│  1. Capture tmuxplexer pane to see layout                       │
│  2. Identify error monitor pane ID                              │
│  3. Capture error monitor pane for latest errors                │
│  4. Analyze errors in full context                              │
│  5. Provide fixes without manual copy/paste                     │
└─────────────────────────────────────────────────────────────────┘
```

## Data Format

### Error Event Schema

```json
{
  "timestamp": 1704148335123,
  "type": "error" | "warning" | "log" | "network" | "performance",
  "level": "error" | "warn" | "info",
  "message": "TypeError: Cannot read property 'user' of null",
  "source": {
    "file": "LoginForm.tsx",
    "line": 45,
    "column": 12,
    "url": "http://localhost:3000/src/LoginForm.tsx"
  },
  "stack": [
    "LoginForm.handleSubmit (LoginForm.tsx:45)",
    "onClick (Button.tsx:12)",
    "callCallback (react-dom.js:4164)"
  ],
  "context": {
    "url": "http://localhost:3000/login",
    "userAgent": "Chrome/120.0.0.0",
    "viewport": "1920x1080",
    "component": "LoginForm"
  },
  "metadata": {
    "reactVersion": "18.2.0",
    "buildId": "dev-20250124-143215"
  }
}
```

### Network Error Schema

```json
{
  "timestamp": 1704148336789,
  "type": "network",
  "method": "POST",
  "url": "/api/login",
  "status": 401,
  "statusText": "Unauthorized",
  "duration": 245,
  "request": {
    "headers": {"Content-Type": "application/json"},
    "body": "{\"email\":\"user@example.com\"}"
  },
  "response": {
    "headers": {"Content-Type": "application/json"},
    "body": "{\"error\":\"Invalid token\"}",
    "size": 27
  }
}
```

## Implementation Phases

### Phase 0: Proof of Concept (30 minutes)

**Goal**: Validate the concept with minimal code

```bash
# Terminal window 1: Mock error stream
while true; do
  echo "[$(date +%H:%M:%S)] Error: Mock error at line $RANDOM"
  sleep 5
done

# Terminal window 2: Capture with Claude
tmux capture-pane -p -t {session}:{window}.{pane}
```

**Success Criteria**: Can see error stream in tmux, Claude can capture it

### Phase 1: File-Based Bridge (2 hours)

**Goal**: Real browser errors to terminal via file watch

**Components**:

1. **Dev Server Endpoint** (`src/lib/dev-error-endpoint.js`):
```javascript
// Add to Next.js/Vite/Express
let errorLog = [];

app.get('/_dev/errors', (req, res) => {
  res.json(errorLog);
});

app.post('/_dev/errors', (req, res) => {
  errorLog.push({
    timestamp: Date.now(),
    ...req.body
  });
  // Keep last 100 errors
  if (errorLog.length > 100) errorLog.shift();
  res.json({ success: true });
});

app.delete('/_dev/errors', (req, res) => {
  errorLog = [];
  res.json({ success: true });
});
```

2. **Client-Side Error Capture** (add to app entry point):
```javascript
// src/lib/dev-error-reporter.js
if (process.env.NODE_ENV === 'development') {
  const reportError = (error) => {
    fetch('/_dev/errors', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message: error.message,
        stack: error.stack,
        url: window.location.href,
        timestamp: Date.now()
      })
    }).catch(() => {}); // Fail silently
  };

  window.addEventListener('error', (e) => reportError(e.error));
  window.addEventListener('unhandledrejection', (e) =>
    reportError(new Error(e.reason))
  );

  const originalError = console.error;
  console.error = (...args) => {
    originalError(...args);
    reportError(new Error(args.join(' ')));
  };
}
```

3. **Terminal Monitor Script** (`scripts/browser-error-monitor.sh`):
```bash
#!/bin/bash
# Monitor browser errors from dev server

DEV_SERVER=${1:-http://localhost:3000}
POLL_INTERVAL=${2:-1}

echo "Monitoring browser errors from $DEV_SERVER..."
echo "Press Ctrl+C to stop"
echo ""

LAST_COUNT=0

while true; do
  RESPONSE=$(curl -s "$DEV_SERVER/_dev/errors" | jq -r '.[] |
    "[\(.timestamp | strftime("%H:%M:%S"))] ❌ \(.message)\n   at \(.url)"')

  CURRENT_COUNT=$(echo "$RESPONSE" | wc -l)

  if [ "$CURRENT_COUNT" -ne "$LAST_COUNT" ]; then
    clear
    echo "=== Browser Errors (Last 100) ==="
    echo ""
    echo "$RESPONSE"
    LAST_COUNT=$CURRENT_COUNT
  fi

  sleep "$POLL_INTERVAL"
done
```

4. **Tmuxplexer Template**:
```json
{
  "name": "Frontend Debug",
  "layout": "2x2",
  "working_dir": "~/projects/myapp",
  "panes": [
    {
      "command": "claude-code .",
      "title": "Claude AI"
    },
    {
      "command": "npm run dev",
      "title": "Dev Server"
    },
    {
      "command": "bash ~/projects/tmuxplexer/scripts/browser-error-monitor.sh",
      "title": "Browser Errors"
    },
    {
      "command": "nvim",
      "title": "Editor"
    }
  ]
}
```

**Success Criteria**:
- Browser errors appear in tmux pane within 1 second
- Claude can capture and read errors
- Errors persist in scrollback history

### Phase 2: Browser Extension (1-2 days)

**Goal**: Capture ALL browser activity without modifying app code

**Components**:

1. **Chrome Extension** (`browser-extension/manifest.json`):
```json
{
  "manifest_version": 3,
  "name": "Tmux DevTools Bridge",
  "version": "1.0.0",
  "description": "Stream browser errors to tmux",
  "permissions": [
    "debugger",
    "webRequest",
    "storage"
  ],
  "host_permissions": ["<all_urls>"],
  "background": {
    "service_worker": "background.js"
  },
  "content_scripts": [{
    "matches": ["<all_urls>"],
    "js": ["content.js"],
    "run_at": "document_start"
  }]
}
```

2. **Content Script** (`browser-extension/content.js`):
```javascript
// Inject into page context to capture console
(function() {
  const originalError = console.error;
  const originalWarn = console.warn;
  const originalLog = console.log;

  function sendToBackground(level, args) {
    const error = new Error();
    chrome.runtime.sendMessage({
      type: 'console',
      level,
      message: args.map(a => String(a)).join(' '),
      stack: error.stack,
      url: window.location.href,
      timestamp: Date.now()
    });
  }

  console.error = function(...args) {
    originalError.apply(console, args);
    sendToBackground('error', args);
  };

  console.warn = function(...args) {
    originalWarn.apply(console, args);
    sendToBackground('warn', args);
  };

  // Capture uncaught errors
  window.addEventListener('error', (e) => {
    sendToBackground('error', [e.message]);
  });

  window.addEventListener('unhandledrejection', (e) => {
    sendToBackground('error', ['Unhandled Promise:', e.reason]);
  });
})();
```

3. **Background Worker** (`browser-extension/background.js`):
```javascript
let ws = null;

function connectWebSocket() {
  ws = new WebSocket('ws://localhost:8765');

  ws.onopen = () => {
    console.log('Connected to terminal bridge');
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
    setTimeout(connectWebSocket, 5000);
  };

  ws.onclose = () => {
    console.log('Disconnected from terminal bridge');
    setTimeout(connectWebSocket, 5000);
  };
}

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      ...message,
      tabId: sender.tab?.id,
      tabUrl: sender.tab?.url
    }));
  }
});

// Start connection
connectWebSocket();

// Monitor network requests
chrome.webRequest.onCompleted.addListener(
  (details) => {
    if (details.statusCode >= 400) {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
          type: 'network',
          method: details.method,
          url: details.url,
          status: details.statusCode,
          timestamp: Date.now()
        }));
      }
    }
  },
  { urls: ["<all_urls>"] }
);
```

4. **WebSocket Bridge Server** (`scripts/browser-bridge-server.js`):
```javascript
const WebSocket = require('ws');
const chalk = require('chalk');

const wss = new WebSocket.Server({ port: 8765 });

console.log(chalk.green('🔌 Browser bridge listening on ws://localhost:8765'));
console.log(chalk.dim('Waiting for browser connection...\n'));

wss.on('connection', (ws) => {
  console.log(chalk.cyan('✓ Browser connected'));

  ws.on('message', (data) => {
    const msg = JSON.parse(data);

    const timestamp = new Date(msg.timestamp).toLocaleTimeString();

    switch (msg.type) {
      case 'console':
        const icon = msg.level === 'error' ? '❌' :
                     msg.level === 'warn' ? '⚠️' : 'ℹ️';
        const color = msg.level === 'error' ? chalk.red :
                      msg.level === 'warn' ? chalk.yellow : chalk.blue;

        console.log(
          chalk.dim(`[${timestamp}]`),
          icon,
          color(msg.message)
        );

        if (msg.stack) {
          const stackLines = msg.stack.split('\n').slice(1, 4);
          stackLines.forEach(line => {
            console.log(chalk.dim('   ' + line.trim()));
          });
        }
        break;

      case 'network':
        if (msg.status >= 400) {
          console.log(
            chalk.dim(`[${timestamp}]`),
            '🌐',
            chalk.red(`${msg.method} ${msg.url} → ${msg.status}`)
          );
        }
        break;
    }

    console.log(); // Blank line for readability
  });

  ws.on('close', () => {
    console.log(chalk.yellow('✗ Browser disconnected'));
  });
});
```

**Success Criteria**:
- Extension captures all console output
- Network errors (4xx, 5xx) appear in terminal
- Works without modifying application code
- Reconnects automatically if bridge server restarts

### Phase 3: Advanced Features (Optional)

**Features to Consider**:

1. **React DevTools Integration**: Capture component errors and warnings
2. **Performance Monitoring**: Web Vitals, long tasks, memory usage
3. **CSS Issue Detection**: Missing styles, layout shifts
4. **Filtering**: Only show errors matching patterns (file, component, severity)
5. **Source Maps**: Resolve minified code to original source
6. **Recording**: Save error sessions for replay
7. **Notifications**: Alert when critical errors occur
8. **Two-Way Commands**: Send commands from terminal to clear console, trigger actions

## Integration with Tmuxplexer

### Template Updates

Add error monitor to default templates:

```json
{
  "name": "Full Stack Dev (4x2)",
  "layout": "4x2",
  "panes": [
    {"command": "claude-code .", "title": "Claude AI"},
    {"command": "nvim", "title": "Editor"},
    {"command": "npm run dev", "title": "Frontend"},
    {"command": "npm run server", "title": "Backend"},
    {"command": "bash ~/tmuxplexer/scripts/browser-bridge.sh", "title": "Browser Errors"},
    {"command": "lazygit", "title": "Git"},
    {"command": "npm test -- --watch", "title": "Tests"},
    {"command": "btop", "title": "Monitor"}
  ]
}
```

### Claude Instructions

Add to project `.claude/` or `CLAUDE.md`:

```markdown
## Browser Error Debugging

This project has automatic browser error streaming to tmux.

To see current browser errors:
1. Capture tmuxplexer to identify pane IDs: `tmux capture-pane -p -t {session}:0.0`
2. Find "Browser Errors" pane (usually window 0, pane 4)
3. Capture errors: `tmux capture-pane -p -t {session}:0.4 -S -50`

The error pane shows:
- JavaScript errors with stack traces
- Network failures (4xx/5xx responses)
- Console warnings
- React component errors

All errors include timestamps and source locations.
```

### Pane Detection

Enhance tmuxplexer to detect error monitor panes:

```go
// tmux.go - Add to TmuxPane struct
type TmuxPane struct {
  ID     string
  Active bool
  Title  string
  // New field:
  IsErrorMonitor bool // Detected by title or command
}

// Detect error monitor panes
func detectErrorMonitorPane(session, window string) *TmuxPane {
  panes, _ := listPanes(session, window)
  for _, pane := range panes {
    if strings.Contains(pane.Title, "Browser") &&
       strings.Contains(pane.Title, "Error") {
      return &pane
    }
  }
  return nil
}
```

## Testing Strategy

### Unit Tests

1. **Error Serialization**: Test JSON schema validation
2. **Message Filtering**: Test error level filtering
3. **Stack Trace Parsing**: Test source map resolution

### Integration Tests

1. **Extension → Bridge**: Test WebSocket connection
2. **Bridge → Terminal**: Test formatted output
3. **Terminal → Claude**: Test capture-pane readability

### End-to-End Tests

1. Create session with error monitor template
2. Trigger browser error (intentional bug)
3. Verify error appears in tmux pane within 1 second
4. Capture pane and verify format is parseable
5. Verify scrollback preserves history

## Performance Considerations

### Browser Impact

- Extension should be lightweight (<1MB memory)
- Debounce rapid errors (max 10/second)
- Use content script injection only when devtools open

### Terminal Impact

- Buffer errors (flush every 100ms)
- Limit scrollback to 10,000 lines
- Use efficient formatting (no heavy JSON parsing per line)

### Network Impact

- WebSocket reconnection with exponential backoff
- Compress large payloads (stack traces)
- Local-only (localhost), no external requests

## Security Considerations

1. **Development Only**: Only activate in NODE_ENV=development
2. **Localhost Binding**: WebSocket server binds to 127.0.0.1 only
3. **No Sensitive Data**: Don't log auth tokens, passwords, PII
4. **Extension Permissions**: Minimize required permissions
5. **Input Validation**: Sanitize all messages before display

## Alternative Approaches Considered

### 1. Chrome DevTools Protocol (CDP)

**Pros**: Full access to all browser events
**Cons**: Requires running Chrome in debug mode, more complex

### 2. Proxy Server (mitmproxy)

**Pros**: No browser extension needed
**Cons**: Only captures network, misses console logs

### 3. Browser DevTools API

**Pros**: Official API, well-documented
**Cons**: Requires extension, similar to our approach

### 4. Remote Debugging Protocol

**Pros**: Can debug remotely
**Cons**: Complex setup, not suitable for local dev

## Success Metrics

1. **Time Saved**: Measure reduction in copy/paste operations
2. **Error Detection Speed**: Time from error to Claude awareness
3. **Adoption**: Number of developers using the workflow
4. **Reliability**: Uptime of error streaming (target: 99%)

## References

- [Chrome Extension Manifest V3](https://developer.chrome.com/docs/extensions/mv3/)
- [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/)
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
- [Tmux Capture Pane](https://man7.org/linux/man-pages/man1/tmux.1.html#BUFFERS)
- [Source Maps Specification](https://sourcemaps.info/spec.html)

## Future Enhancements

1. **Multi-Browser Support**: Firefox, Safari, Edge
2. **Mobile Debugging**: Remote debugging for mobile devices
3. **Log Aggregation**: Combine frontend + backend + database logs
4. **Error Analytics**: Track error frequency, patterns over time
5. **AI Error Suggestions**: Claude proactively suggests fixes for common errors
6. **Error Replay**: Time-travel debugging with full state snapshots

## Open Questions

1. Should we support production error monitoring (with opt-in)?
2. How to handle source maps in Docker/remote dev environments?
3. Should we integrate with existing tools (Sentry, LogRocket)?
4. How to handle multiple browser tabs/windows?
5. Should we build a TUI for the error monitor pane (interactive filtering)?

---

# Multi-AI Collaboration via Tmux Orchestration

**Status**: Planning / Proof of Concept
**Purpose**: Enable multi-turn AI-to-AI conversations through tmux for complex debugging and consensus building
**Related**: See `~/ObsidianVault/Projects/TFE/CODEX_COLLAB_DESIGN.md` for original concept

## Problem Statement

Slash commands (`/codex`, `/gemini`, etc.) work great for single-query AI consultation, but complex debugging often requires:

1. **Multi-turn conversations**: Building context over several exchanges
2. **Multiple AI perspectives**: Getting consensus from different models
3. **Stateful sessions**: AI maintains conversation history
4. **Visual monitoring**: Seeing all AI interactions in parallel
5. **Tool integration**: Combining AI responses with logs, metrics, browser errors

**Current workflow limitation**:
```
User → Claude → /codex → wait → copy/paste response → /gemini → wait → compare
```

This breaks flow state and requires manual orchestration.

## Proposed Solution

Transform tmuxplexer into an **AI Orchestration Platform** where Claude can:

1. **Send questions to AI panes** via `tmux send-keys`
2. **Capture responses** via `tmux capture-pane` (already implemented)
3. **Orchestrate multi-turn conversations** by generating follow-ups
4. **Synthesize consensus** from multiple AI responses
5. **Monitor progress** visually in real-time

### Key Insight

Tmuxplexer already has 50% of this working:
- ✅ Visual layout of all tools/AIs
- ✅ Capture pane output (browser errors, logs, metrics)
- ⏳ **Missing**: Send input to panes (the last piece!)

Once we add "send to pane", tmuxplexer becomes the orchestration layer.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│             Tmuxplexer (AI Orchestration Layer)                  │
│                                                                   │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Claude Code (Orchestrator)                               │  │
│  │  - Sees layout of all AIs and tools                       │  │
│  │  - Sends questions to specific panes                      │  │
│  │  - Captures responses and generates follow-ups            │  │
│  │  - Synthesizes multi-AI consensus                         │  │
│  │  - Correlates with logs/errors/metrics                    │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                   │
│  Sends questions via:      tmux send-keys -t {pane} "..." C-m   │
│  Reads responses via:      tmux capture-pane -p -t {pane}       │
└───────────────┬───────────────┬───────────────┬──────────────────┘
                │               │               │
                ▼               ▼               ▼
    ┌───────────────────┐ ┌───────────────┐ ┌──────────────────┐
    │  Pane 0: Codex    │ │ Pane 1: Gemini│ │ Pane 2: Cursor   │
    │  (GPT-5 Deep      │ │ (Gemini 2.0   │ │ (Claude Opus)    │
    │   Reasoning)      │ │  Flash Exp)   │ │                  │
    │                   │ │               │ │                  │
    │ codex>            │ │ gemini>       │ │ cursor>          │
    │ [Waiting...]      │ │ [Thinking...] │ │ [Ready]          │
    └───────────────────┘ └───────────────┘ └──────────────────┘
                │               │               │
                └───────────────┴───────────────┘
                              │
                              ▼
                    Synthesized Consensus
```

### Additional Tool Integration

Not just AI-to-AI, but also AI-to-Tool:

```
┌──────────────────────────────────────────────────────────┐
│  Multi-Modal Debugging Dashboard                          │
├──────────────────────────────────────────────────────────┤
│  Pane 0: Codex         │  Pane 1: Browser Errors         │
│  "Check auth flow"     │  [14:32] 401 Unauthorized       │
│  Response: "JWT..."    │  POST /api/login                │
├──────────────────────────────────────────────────────────┤
│  Pane 2: lnav (Logs)   │  Pane 3: bottom (Metrics)       │
│  ERROR: Invalid token  │  CPU: 45% (node spike)          │
│  auth.ts:89            │  Memory: 2.1GB                  │
└──────────────────────────────────────────────────────────┘

Claude orchestrates:
1. See browser error (401)
2. Ask Codex: "Why 401 on /api/login?"
3. Capture logs showing "Invalid token"
4. See CPU spike in metrics
5. Synthesize: "Auth middleware rejecting valid tokens, causing retries"
```

## Use Cases

### 1. Multi-Turn Debugging Dialogue

**User**: "Debug this authentication issue"

**Tmuxplexer orchestrates**:

```
Turn 1: Claude → Codex (Pane 0)
  Q: "What's wrong with this auth flow?"
  A: "The JWT validation is checking the wrong claim. Should be 'sub' not 'userId'"

Turn 2: Claude → Cursor (Pane 1)
  Q: "Codex says JWT claim is wrong. Where's the validation code?"
  A: "Found in middleware/auth.ts:45, checks 'userId' but should check 'sub'"

Turn 3: Claude → Gemini (Pane 2)
  Q: "Both AIs agree on the issue. Are there other places with this pattern?"
  A: "Yes, found 3 more: api/users.ts:12, api/posts.ts:34, lib/verify.ts:89"

Claude synthesizes:
"All three AIs agree. The JWT validation uses 'userId' claim instead of
standard 'sub' claim. This appears in 4 files. Here's the fix..."

[Claude implements fix across all 4 files]
```

### 2. Multi-AI Consensus Building

**User**: "Which architecture is better for this feature?"

**Tmuxplexer orchestrates**:

```
Parallel Query (all at once):
├─ Codex:   "Microservices - better scalability"
├─ Gemini:  "Monolith - simpler for current team size"
└─ Cursor:  "Modular monolith - best of both worlds"

Claude analyzes responses:
"2 of 3 AIs recommend monolithic approach. Codex suggests microservices
but your team is 3 people. Cursor's 'modular monolith' balances both.
Recommendation: Start modular monolith, extract services later if needed."
```

### 3. Specialized AI Routing

**Claude's orchestration logic**:

```go
func (m model) routeQuestion(question string) []int {
  var panes []int

  // Route by question type
  if isPerformanceIssue(question) {
    panes = append(panes, m.findPane("Cursor"))  // Best at optimization
    panes = append(panes, m.findPane("Gemini"))  // Fast, can run tests
  } else if isArchitectureQuestion(question) {
    panes = append(panes, m.findPane("Codex"))   // Deep reasoning
    panes = append(panes, m.findPane("Cursor"))  // Practical experience
  } else if isUIBug(question) {
    panes = append(panes, m.findPane("Cursor"))  // Frontend expert
    panes = append(panes, m.findPane("Browser")) // Check console errors
  }

  return panes
}
```

### 4. Correlation Across Sources

**User**: "Why is the app slow?"

**Tmuxplexer correlates**:

```
1. Capture browser errors pane:
   "[14:30:15] Warning: Slow render (3200ms)"

2. Ask Codex (focus on architecture):
   Q: "Why would rendering take 3+ seconds?"
   A: "Check for unnecessary re-renders, large lists without virtualization"

3. Capture bottom metrics pane:
   "Memory: 4.2GB (2GB increase in last 5min)"

4. Ask Gemini (focus on memory):
   Q: "Memory grew 2GB during slow render. What's leaking?"
   A: "Likely large data structures in component state"

5. Capture lnav logs:
   "API returned 50,000 records to /dashboard"

Claude synthesizes:
"Found the issue: /dashboard API returns 50K records (2GB), no pagination.
Frontend renders all at once without virtualization. Fix: Add pagination
to API + react-window for virtual scrolling."
```

## Implementation Phases

### Phase 1: Send Input to Panes (Foundation)

**Goal**: Enable sending text/commands to any pane

**Components**:

1. **Input Mode** (add to `update_keyboard.go`):
```go
case "i": // Insert mode - type to pane
  if m.focusedPanel == "footer" && len(m.windows) > 0 {
    m.mode = "input"
    m.inputPrompt = "Send to pane: "
    m.inputBuffer = ""
    return m, nil
  }

case "enter": // Send input
  if m.mode == "input" {
    return m, m.sendToPaneCmd(m.inputBuffer)
  }
```

2. **Pane Command Sender** (add to `tmux.go`):
```go
// sendToPane sends text to a specific pane and presses Enter
func sendToPane(sessionName string, windowIndex int, paneID string, text string) error {
  target := fmt.Sprintf("%s:%d.%s", sessionName, windowIndex, paneID)

  cmd := exec.Command("tmux", "send-keys", "-t", target, text, "C-m")
  return cmd.Run()
}

// sendToPaneRaw sends text without pressing Enter (for interactive prompts)
func sendToPaneRaw(sessionName string, windowIndex int, paneID string, text string) error {
  target := fmt.Sprintf("%s:%d.%s", sessionName, windowIndex, paneID)

  cmd := exec.Command("tmux", "send-keys", "-t", target, "-l", text)
  return cmd.Run()
}
```

3. **UI Feedback** (add to `view.go`):
```go
// Show input mode in status bar
if m.mode == "input" {
  status := fmt.Sprintf("INPUT MODE | Pane %s | Type command (ESC to cancel): %s",
    m.getSelectedPaneID(), m.inputBuffer)
  return statusStyle.Render(status)
}
```

**Success Criteria**:
- Press `i` on focused pane
- Type command
- Press Enter
- Command appears in target pane

### Phase 2: AI Response Detection (Smart Capture)

**Goal**: Automatically detect when AI finishes responding

**Components**:

1. **Response Parser** (new file: `ai_detection.go`):
```go
// AIProvider identifies which AI is running in a pane
type AIProvider int

const (
  AIUnknown AIProvider = iota
  AICodex        // OpenAI Codex CLI
  AIGemini       // Google Gemini CLI
  AICursor       // Cursor CLI
  AIWindsurf     // Windsurf CLI
  AIClaude       // Claude Code
)

// DetectAIProvider identifies AI by pane title or command
func DetectAIProvider(pane TmuxPane) AIProvider {
  title := strings.ToLower(pane.Title)

  if strings.Contains(title, "codex") {
    return AICodex
  } else if strings.Contains(title, "gemini") {
    return AIGemini
  } else if strings.Contains(title, "cursor") {
    return AICursor
  } else if strings.Contains(title, "windsurf") {
    return AIWindsurf
  } else if strings.Contains(title, "claude") {
    return AIClaude
  }

  return AIUnknown
}

// WaitForAIResponse polls pane until response is complete
func WaitForAIResponse(sessionName string, windowIndex int, paneID string,
                       provider AIProvider, timeout time.Duration) (string, error) {

  startTime := time.Now()
  lastOutput := ""

  for time.Since(startTime) < timeout {
    // Capture current pane content
    output, err := capturePane(paneID)
    if err != nil {
      return "", err
    }

    // Check if response is complete based on AI provider
    if isResponseComplete(output, lastOutput, provider) {
      return extractResponse(output, provider), nil
    }

    lastOutput = output
    time.Sleep(2 * time.Second) // Poll every 2 seconds
  }

  return "", fmt.Errorf("timeout waiting for AI response")
}

// isResponseComplete detects when AI finishes based on output patterns
func isResponseComplete(current, previous string, provider AIProvider) bool {
  switch provider {
  case AICodex:
    // Codex shows "tokens used: 1234" when done
    return strings.Contains(current, "tokens used:")

  case AIGemini:
    // Gemini shows prompt "gemini>" when ready
    lines := strings.Split(current, "\n")
    if len(lines) > 0 {
      lastLine := strings.TrimSpace(lines[len(lines)-1])
      return lastLine == "gemini>"
    }

  case AICursor:
    // Cursor shows cursor blinking indicator
    return strings.Contains(current, "cursor>") &&
           len(current) > len(previous) && // Content stopped growing
           strings.Count(current, "\n") == strings.Count(previous, "\n")

  case AIWindsurf:
    // Windsurf shows "Ready" status
    return strings.Contains(current, "[Ready]")
  }

  return false
}

// extractResponse removes prompts and metadata, returns just AI's answer
func extractResponse(output string, provider AIProvider) string {
  lines := strings.Split(output, "\n")

  var responseLines []string
  inResponse := false

  for _, line := range lines {
    // Skip empty lines at start
    if !inResponse && strings.TrimSpace(line) == "" {
      continue
    }

    // Start of response (after user's question)
    if !inResponse && !isPromptLine(line, provider) {
      inResponse = true
    }

    // End of response (metadata/token count)
    if inResponse && isMetadataLine(line, provider) {
      break
    }

    if inResponse && !isPromptLine(line, provider) {
      responseLines = append(responseLines, line)
    }
  }

  return strings.TrimSpace(strings.Join(responseLines, "\n"))
}
```

2. **Polling Strategy** (configurable):
```go
type ResponsePollingConfig struct {
  InitialDelay  time.Duration // Wait before first check (e.g. 1s)
  PollInterval  time.Duration // Check every N seconds (e.g. 2s)
  MaxTimeout    time.Duration // Give up after N seconds (e.g. 60s)
  StableChecks  int           // Require N stable checks to confirm done
}

var DefaultPollingConfig = ResponsePollingConfig{
  InitialDelay: 1 * time.Second,
  PollInterval: 2 * time.Second,
  MaxTimeout:   60 * time.Second,
  StableChecks: 2, // Output unchanged for 2 checks = done
}
```

**Success Criteria**:
- Send question to Codex pane
- Function waits automatically
- Returns when "tokens used:" appears
- Extracts clean response without prompts/metadata

### Phase 3: Multi-Turn Orchestrator (The Magic)

**Goal**: Enable Claude to conduct multi-turn conversations with AIs

**Components**:

1. **Collaboration Session** (new file: `collaboration.go`):
```go
// CollaborationTurn represents one Q&A exchange
type CollaborationTurn struct {
  TurnNumber int
  PaneID     string
  AIProvider AIProvider
  Question   string
  Response   string
  Timestamp  time.Time
  Duration   time.Duration
}

// CollaborationSession manages multi-turn AI dialogue
type CollaborationSession struct {
  SessionID    string
  InitialQuery string
  Panes        []string // Which AI panes to involve
  MaxTurns     int
  CurrentTurn  int
  History      []CollaborationTurn
  Consensus    string // Synthesized result
  Status       string // "running", "complete", "error"
}

// StartCollaboration begins a multi-turn AI consultation
func (m model) startCollaborationCmd(query string, panes []string, turns int) tea.Cmd {
  return func() tea.Msg {
    session := &CollaborationSession{
      SessionID:    generateSessionID(),
      InitialQuery: query,
      Panes:        panes,
      MaxTurns:     turns,
      Status:       "running",
    }

    // Execute turns
    for turn := 1; turn <= turns; turn++ {
      session.CurrentTurn = turn

      for _, paneID := range panes {
        // Get AI provider type
        pane := m.findPaneByID(paneID)
        provider := DetectAIProvider(pane)

        // Generate question for this turn
        question := generateQuestionForTurn(session, turn, provider)

        // Send to pane
        startTime := time.Now()
        err := sendToPane(m.sessionName, pane.WindowIndex, paneID, question)
        if err != nil {
          continue
        }

        // Wait for response
        response, err := WaitForAIResponse(
          m.sessionName,
          pane.WindowIndex,
          paneID,
          provider,
          60*time.Second,
        )

        // Record turn
        session.History = append(session.History, CollaborationTurn{
          TurnNumber: turn,
          PaneID:     paneID,
          AIProvider: provider,
          Question:   question,
          Response:   response,
          Timestamp:  startTime,
          Duration:   time.Since(startTime),
        })
      }

      // Between-turn analysis
      if turn < turns {
        // Let Claude analyze responses and plan next turn
        // This could be a brief pause or call to another AI
        time.Sleep(1 * time.Second)
      }
    }

    // Synthesize consensus from all turns
    session.Consensus = synthesizeConsensus(session.History)
    session.Status = "complete"

    return collaborationCompleteMsg{session}
  }
}

// generateQuestionForTurn creates contextual question based on previous turns
func generateQuestionForTurn(session *CollaborationSession, turn int, provider AIProvider) string {
  if turn == 1 {
    // First turn: ask initial question
    return session.InitialQuery
  }

  // Subsequent turns: build context from previous responses
  context := fmt.Sprintf("Previous question: %s\n\n", session.InitialQuery)

  // Add relevant previous responses
  for _, prevTurn := range session.History {
    if prevTurn.TurnNumber == turn - 1 {
      context += fmt.Sprintf("%s said: %s\n\n",
        aiProviderName(prevTurn.AIProvider),
        truncate(prevTurn.Response, 200))
    }
  }

  // Generate follow-up question (this would be Claude's logic)
  // For now, a simple template:
  followUp := fmt.Sprintf("Based on the previous analysis, %s",
    getFollowUpPrompt(turn, provider))

  return context + followUp
}

// synthesizeConsensus combines all AI responses into final recommendation
func synthesizeConsensus(history []CollaborationTurn) string {
  var consensus strings.Builder

  consensus.WriteString("## Multi-AI Analysis Results\n\n")

  // Group by turn
  turnMap := make(map[int][]CollaborationTurn)
  for _, turn := range history {
    turnMap[turn.TurnNumber] = append(turnMap[turn.TurnNumber], turn)
  }

  // Show progression
  for i := 1; i <= len(turnMap); i++ {
    consensus.WriteString(fmt.Sprintf("### Turn %d\n\n", i))

    for _, turn := range turnMap[i] {
      consensus.WriteString(fmt.Sprintf("**%s**: %s\n\n",
        aiProviderName(turn.AIProvider),
        summarize(turn.Response)))
    }
  }

  // Find common themes
  consensus.WriteString("### Consensus\n\n")
  agreements := findAgreements(history)
  for _, agreement := range agreements {
    consensus.WriteString(fmt.Sprintf("- ✓ All AIs agree: %s\n", agreement))
  }

  return consensus.String()
}
```

2. **Message Types** (add to `types.go`):
```go
// collaborationStartedMsg signals start of multi-turn session
type collaborationStartedMsg struct {
  sessionID string
}

// collaborationTurnMsg reports progress of each turn
type collaborationTurnMsg struct {
  sessionID string
  turn      int
  paneID    string
  question  string
  response  string
}

// collaborationCompleteMsg delivers final results
type collaborationCompleteMsg struct {
  session *CollaborationSession
}
```

3. **UI Updates** (add to `update.go`):
```go
case collaborationTurnMsg:
  msg := msg.(collaborationTurnMsg)
  m.statusMsg = fmt.Sprintf("Collaboration Turn %d: %s responded",
    msg.turn, msg.paneID)
  return m, nil

case collaborationCompleteMsg:
  msg := msg.(collaborationCompleteMsg)
  m.statusMsg = "Collaboration complete! View results in footer."
  m.collaborationResults = msg.session
  m.updateFooterContent() // Show consensus
  return m, nil
```

**Success Criteria**:
- User triggers collaboration with 3 AIs, 3 turns
- Claude sends initial question to all 3 panes
- Claude waits for all responses
- Claude generates follow-ups based on responses
- Process repeats for 3 turns
- Claude synthesizes consensus
- User sees complete conversation history

### Phase 4: Templates and Presets

**Goal**: Pre-configured AI collaboration setups

**Components**:

1. **AI Arena Template** (add to `templates.json`):
```json
{
  "name": "AI Debug Arena (3-Way)",
  "description": "Multi-AI debugging with Codex, Gemini, Cursor",
  "layout": "2x2",
  "working_dir": "~",
  "panes": [
    {
      "command": "codex",
      "title": "Codex (GPT-5)"
    },
    {
      "command": "gemini -m gemini-2.0-flash-exp",
      "title": "Gemini 2.0"
    },
    {
      "command": "cursor",
      "title": "Cursor (Claude)"
    },
    {
      "command": "bash",
      "title": "Orchestrator"
    }
  ],
  "collaboration": {
    "enabled": true,
    "default_turns": 3,
    "ai_panes": [0, 1, 2]
  }
}
```

2. **Full Debug Dashboard Template**:
```json
{
  "name": "Full Stack Debug Dashboard",
  "description": "AI collaboration + monitoring tools",
  "layout": "4x2",
  "working_dir": "~/projects/myapp",
  "panes": [
    {"command": "codex", "title": "Codex"},
    {"command": "gemini -m gemini-2.0-flash-exp", "title": "Gemini"},
    {"command": "npm run dev", "title": "Dev Server"},
    {"command": "node ~/tmuxplexer/scripts/browser-bridge.js", "title": "Browser Errors"},
    {"command": "lnav ~/projects/myapp/logs/*.log", "title": "Application Logs"},
    {"command": "bottom", "title": "System Metrics"},
    {"command": "lazygit", "title": "Git"},
    {"command": "bash", "title": "Terminal"}
  ],
  "collaboration": {
    "enabled": true,
    "default_turns": 2,
    "ai_panes": [0, 1],
    "tool_panes": [3, 4, 5]
  }
}
```

**Success Criteria**:
- Create session from "AI Debug Arena" template
- All AI panes initialize correctly
- Tmuxplexer detects AI panes automatically
- User can start collaboration with one command

### Phase 5: UI Enhancements

**Goal**: Visual feedback for collaboration sessions

**Components**:

1. **Collaboration Status Panel** (add to header when active):
```go
func (m model) renderCollaborationStatus() []string {
  if m.activeCollaboration == nil {
    return []string{}
  }

  var lines []string

  session := m.activeCollaboration
  lines = append(lines, fmt.Sprintf("🤝 Collaboration: Turn %d/%d",
    session.CurrentTurn, session.MaxTurns))
  lines = append(lines, "")

  // Show progress for each AI
  for _, paneID := range session.Panes {
    pane := m.findPaneByID(paneID)
    provider := DetectAIProvider(pane)

    status := "⏳ Waiting"
    if hasTurnCompleted(session, session.CurrentTurn, paneID) {
      status = "✓ Complete"
    } else if isTurnInProgress(session, session.CurrentTurn, paneID) {
      status = "💭 Thinking..."
    }

    lines = append(lines, fmt.Sprintf("  %s: %s",
      aiProviderName(provider), status))
  }

  return lines
}
```

2. **Collaboration History View** (footer panel when collaboration completes):
```go
func (m model) renderCollaborationResults() []string {
  if m.collaborationResults == nil {
    return []string{}
  }

  var lines []string
  session := m.collaborationResults

  lines = append(lines, "═══ COLLABORATION RESULTS ═══")
  lines = append(lines, "")
  lines = append(lines, fmt.Sprintf("Query: %s", session.InitialQuery))
  lines = append(lines, fmt.Sprintf("Turns: %d | Duration: %s",
    len(session.History), totalDuration(session)))
  lines = append(lines, "")

  // Show turn-by-turn
  for i := 1; i <= session.MaxTurns; i++ {
    lines = append(lines, fmt.Sprintf("─── Turn %d ───", i))

    for _, turn := range session.History {
      if turn.TurnNumber == i {
        lines = append(lines, fmt.Sprintf("🤖 %s (%s):",
          aiProviderName(turn.AIProvider),
          turn.Duration.Round(time.Second)))

        // Show truncated response
        responseLines := strings.Split(turn.Response, "\n")
        for j, line := range responseLines {
          if j < 5 { // Show first 5 lines
            lines = append(lines, "  "+line)
          }
        }
        if len(responseLines) > 5 {
          lines = append(lines, fmt.Sprintf("  ... (%d more lines)",
            len(responseLines)-5))
        }
        lines = append(lines, "")
      }
    }
  }

  // Show consensus
  lines = append(lines, "═══ CONSENSUS ═══")
  lines = append(lines, "")
  consensusLines := strings.Split(session.Consensus, "\n")
  lines = append(lines, consensusLines...)

  return lines
}
```

3. **Keyboard Shortcuts** (add to `update_keyboard.go`):
```go
case "c": // Start collaboration
  if m.focusedPanel == "left" && len(m.sessions) > 0 {
    // Get AI panes from current session
    aiPanes := m.detectAIPanes()
    if len(aiPanes) > 0 {
      m.statusMsg = "Starting collaboration..."
      return m, m.startCollaborationCmd(
        "Analyze the current codebase for potential issues",
        aiPanes,
        3, // 3 turns
      )
    }
  }
  return m, nil

case "C": // Configure collaboration
  // Open dialog to set question, turns, panes
  m.mode = "collab_config"
  return m, nil
```

**Success Criteria**:
- Press `c` to start collaboration
- Header shows progress ("Turn 2/3", "Gemini: Thinking...")
- Footer updates with real-time responses
- When complete, full conversation history is scrollable
- Consensus appears at bottom

## Integration with Browser Error Streaming

**Combining both features creates a powerful debugging loop**:

```
┌────────────────────────────────────────────────────────────┐
│  Full Stack Debug Dashboard (4x2 layout)                   │
├────────────────────────────────────────────────────────────┤
│  Pane 0: Codex            │  Pane 1: Gemini                │
│  [Ready for questions]    │  [Ready for questions]         │
├────────────────────────────────────────────────────────────┤
│  Pane 2: Dev Server       │  Pane 3: Browser Errors        │
│  ✓ Compiled successfully  │  [14:32] ❌ TypeError: ...     │
│  http://localhost:3000    │  [14:33] 🌐 POST → 401         │
├────────────────────────────────────────────────────────────┤
│  Pane 4: Application Logs │  Pane 5: System Metrics        │
│  ERROR: Token invalid     │  CPU: 45% Memory: 2.1GB        │
│  auth.ts:89               │  node: High CPU spike          │
├────────────────────────────────────────────────────────────┤
│  Pane 6: Git Status       │  Pane 7: Orchestrator          │
│  main ✓ No changes        │  $ tmuxplexer                  │
└────────────────────────────────────────────────────────────┘
```

**Automated debugging workflow**:

1. Browser error appears in Pane 3
2. Claude captures error: `tmux capture-pane -p -t 0:0.3`
3. Claude captures logs: `tmux capture-pane -p -t 0:0.4`
4. Claude starts collaboration:
   - Sends error + logs to Codex (Pane 0)
   - Sends same to Gemini (Pane 1)
5. Both AIs analyze simultaneously
6. Claude synthesizes responses
7. Claude implements fix
8. Browser error disappears automatically

**This is the holy grail: Zero-friction AI-assisted debugging!**

## Keyboard Shortcuts

```
New shortcuts for collaboration:
  i         - Insert mode (send to focused pane)
  c         - Start collaboration (auto-detect AI panes)
  C         - Configure collaboration (custom question/turns)
  ESC       - Cancel input mode or stop collaboration
  v         - View collaboration history (toggle footer view)
```

## Claude Instructions

**Add to project `CLAUDE.md` or `.claude/instructions.md`:**

```markdown
## Multi-AI Collaboration

This project uses tmuxplexer for AI orchestration.

### Available AIs

Check which AIs are running:
1. Capture tmuxplexer: `tmux capture-pane -p -t {session}:0.0`
2. Look for panes with "Codex", "Gemini", "Cursor", etc.

### Starting a Collaboration

To get multiple AI opinions on a complex issue:

1. Identify AI pane IDs from tmuxplexer layout
2. Send question to each AI:
   ```bash
   tmux send-keys -t {session}:0.{pane_id} "Your question here" C-m
   ```
3. Wait 10-30 seconds for responses
4. Capture each response:
   ```bash
   tmux capture-pane -p -t {session}:0.{pane_id} -S -50
   ```
5. Extract AI's answer (skip prompts and metadata)
6. Generate follow-up questions based on responses
7. Repeat for 2-3 turns
8. Synthesize consensus from all AIs

### Example: 3-Way Debugging

```bash
# Turn 1: Ask all AIs
tmux send-keys -t main:0.0 "Why is auth failing?" C-m  # Codex
tmux send-keys -t main:0.1 "Why is auth failing?" C-m  # Gemini
tmux send-keys -t main:0.2 "Why is auth failing?" C-m  # Cursor

sleep 15

# Capture responses
codex_resp=$(tmux capture-pane -p -t main:0.0 -S -30)
gemini_resp=$(tmux capture-pane -p -t main:0.1 -S -30)
cursor_resp=$(tmux capture-pane -p -t main:0.2 -S -30)

# Turn 2: Follow-ups based on Turn 1
# (Continue pattern...)
```

### Automated Collaboration

When tmuxplexer Phase 3 is complete, use:
- `c` key to start collaboration
- Specify turns and panes
- View results in footer panel

### Combining with Browser Errors

For frontend debugging:
1. Browser errors appear in Pane 3 automatically
2. Capture errors: `tmux capture-pane -p -t main:0.3`
3. Send to AIs for analysis
4. Correlate with application logs (Pane 4)
5. Check metrics (Pane 5) for performance impact
```

## Testing Strategy

### Unit Tests

1. **sendToPane**: Test command sending
2. **DetectAIProvider**: Test AI detection logic
3. **isResponseComplete**: Test completion detection
4. **extractResponse**: Test response parsing

### Integration Tests

1. **Single AI Query**: Send question, wait for response
2. **Multi-AI Parallel**: Send to 3 AIs simultaneously
3. **Multi-Turn Conversation**: 3 turns with 2 AIs
4. **Error Recovery**: Handle timeout, connection loss

### End-to-End Tests

1. Create "AI Arena" session
2. Start collaboration with 3 AIs
3. Verify all responses captured
4. Verify consensus generated
5. Verify UI updates correctly

## Performance Considerations

### Response Waiting

- Use configurable polling intervals (default: 2s)
- Implement timeout limits (default: 60s)
- Parallel queries (don't wait sequentially)
- Cancel in-flight requests if user interrupts

### Resource Usage

- Each AI session consumes memory (CLI instances)
- Limit concurrent collaborations to 1 per tmuxplexer session
- Clean up old collaboration histories (keep last 10)
- Option to export collaboration results before clearing

### Network/API Costs

- Multi-turn conversations use more tokens
- Each AI charges separately
- Provide token usage summary after collaboration
- Option to set max token limits per turn

## Safety & Usage Controls

**Critical safeguards to prevent hitting your weekly/monthly subscription limits**:

### Modern Subscription Reality

Users don't pay per token anymore - they have subscription tiers with usage caps:

- **Claude Max**: Resets weekly (e.g., 30% remaining today, 100% tomorrow at midnight)
- **Codex**: Monthly usage limit based on tier
- **Gemini**: Daily/monthly caps depending on plan
- **Cursor**: Usage-based limits per billing cycle

**The real concern**: Burning through your weekly cap on Monday and being throttled until Sunday!

**Critical safeguards to prevent runaway usage**:

### Hard Limits (Non-Negotiable)

```go
// collaboration.go - Safety limits
const (
  MaxTurnsAllowed      = 5     // Absolute max turns per collaboration
  MaxAIsPerSession     = 4     // Max concurrent AIs
  MaxTokensPerTurn     = 4000  // Abort if single response exceeds this
  MaxTotalTokens       = 20000 // Abort entire collaboration if exceeded
  DefaultTurns         = 2     // Conservative default
  MaxConcurrentCollab  = 1     // Only one collaboration at a time
)

// CollaborationLimits enforced before starting
type CollaborationLimits struct {
  MaxTurns       int
  MaxAIs         int
  MaxTokensTotal int
  RequireApproval bool // Require user approval for each turn
}

var DefaultLimits = CollaborationLimits{
  MaxTurns:        2,
  MaxAIs:          3,
  MaxTokensTotal:  10000,
  RequireApproval: false, // Can enable for extra safety
}
```

### Usage Estimation (Before Starting)

```go
// Subscription tier configuration
type SubscriptionTier struct {
  Provider     AIProvider
  LimitType    string  // "weekly", "daily", "monthly"
  UsageLimit   int     // Total requests or tokens allowed
  CurrentUsage int     // How much used so far
  ResetTime    time.Time // When usage resets
  ResetDay     string  // e.g., "Sunday midnight", "1st of month"
}

var UserSubscriptions = map[AIProvider]SubscriptionTier{
  AIClaude: {
    Provider:     AIClaude,
    LimitType:    "weekly",
    UsageLimit:   100,  // 100 requests per week (Claude Max example)
    CurrentUsage: 70,   // 70 used, 30 remaining (30% left)
    ResetTime:    time.Parse("2025-01-26 00:00:00"), // Tonight at midnight
    ResetDay:     "Sunday midnight",
  },
  AICodex: {
    Provider:     AICodex,
    LimitType:    "monthly",
    UsageLimit:   500,  // 500 requests per month
    CurrentUsage: 120,  // 24% used
    ResetTime:    time.Parse("2025-02-01 00:00:00"),
    ResetDay:     "1st of month",
  },
  AIGemini: {
    Provider:     AIGemini,
    LimitType:    "daily",
    UsageLimit:   1000, // 1000 requests per day
    CurrentUsage: 234,  // 23% used
    ResetTime:    time.Now().Add(24 * time.Hour),
    ResetDay:     "Daily at midnight",
  },
}

// EstimateCollaborationUsage calculates usage impact
func EstimateCollaborationUsage(query string, providers []AIProvider, turns int) UsageEstimate {
  estimate := UsageEstimate{
    NumTurns: turns,
    PerAI:    make(map[AIProvider]AIUsageImpact),
  }

  for _, provider := range providers {
    tier := UserSubscriptions[provider]
    requestsNeeded := turns // One request per turn

    remaining := tier.UsageLimit - tier.CurrentUsage
    percentOfRemaining := (float64(requestsNeeded) / float64(remaining)) * 100
    percentOfTotal := (float64(tier.CurrentUsage + requestsNeeded) / float64(tier.UsageLimit)) * 100

    estimate.PerAI[provider] = AIUsageImpact{
      Provider:           provider,
      RequestsNeeded:     requestsNeeded,
      CurrentUsage:       tier.CurrentUsage,
      RemainingBefore:    remaining,
      RemainingAfter:     remaining - requestsNeeded,
      PercentOfRemaining: percentOfRemaining,
      PercentOfTotal:     percentOfTotal,
      ResetTime:          tier.ResetTime,
    }
  }

  return estimate
}

type UsageEstimate struct {
  NumTurns int
  PerAI    map[AIProvider]AIUsageImpact
}

type AIUsageImpact struct {
  Provider           AIProvider
  RequestsNeeded     int
  CurrentUsage       int
  RemainingBefore    int
  RemainingAfter     int
  PercentOfRemaining float64 // What % of remaining quota will be used
  PercentOfTotal     float64 // Total usage % after collaboration
  ResetTime          time.Time
}

// Show estimate to user BEFORE starting
func (m model) showCollaborationEstimate(estimate UsageEstimate) string {
  var output strings.Builder

  output.WriteString("═══ COLLABORATION USAGE ESTIMATE ═══\n\n")
  output.WriteString(fmt.Sprintf("Turns: %d\n\n", estimate.NumTurns))

  for provider, impact := range estimate.PerAI {
    name := aiProviderName(provider)
    icon := getUsageIcon(impact.PercentOfTotal)

    output.WriteString(fmt.Sprintf("%s %s:\n", icon, name))
    output.WriteString(fmt.Sprintf("  Currently: %d/%d used (%.0f%%)\n",
      impact.CurrentUsage, impact.CurrentUsage + impact.RemainingBefore,
      (float64(impact.CurrentUsage) / float64(impact.CurrentUsage + impact.RemainingBefore)) * 100))
    output.WriteString(fmt.Sprintf("  This will use: %d requests (%.0f%% of remaining)\n",
      impact.RequestsNeeded, impact.PercentOfRemaining))
    output.WriteString(fmt.Sprintf("  After: %d remaining (%.0f%% used)\n",
      impact.RemainingAfter, impact.PercentOfTotal))

    // Time until reset
    timeUntil := time.Until(impact.ResetTime)
    if timeUntil < 24*time.Hour {
      output.WriteString(fmt.Sprintf("  ⏰ Resets in %s\n", formatDuration(timeUntil)))
    } else {
      output.WriteString(fmt.Sprintf("  Resets: %s\n", impact.ResetTime.Format("Jan 2, 3:04pm")))
    }
    output.WriteString("\n")
  }

  output.WriteString("Continue? [y/N]: ")
  return output.String()
}

func getUsageIcon(percentUsed float64) string {
  if percentUsed < 50 {
    return "🟢" // Green - plenty left
  } else if percentUsed < 75 {
    return "🟡" // Yellow - getting high
  } else if percentUsed < 90 {
    return "🟠" // Orange - careful!
  } else {
    return "🔴" // Red - almost out!
  }
}
```

### Turn-by-Turn Approval (Optional Safety Mode)

```go
// ApprovalMode: User must approve each turn before sending
func (m model) startCollaborationWithApproval(query string, panes []string, turns int) tea.Cmd {
  return func() tea.Msg {
    session := &CollaborationSession{
      RequireApproval: true,
      // ... rest of setup
    }

    for turn := 1; turn <= turns; turn++ {
      if session.RequireApproval {
        // Show preview of what will be sent
        preview := generateQuestionForTurn(session, turn, provider)

        // Prompt user
        fmt.Printf("\n=== Turn %d/%d ===\n", turn, turns)
        fmt.Printf("Will send to %d AIs:\n%s\n\n", len(panes), preview)
        fmt.Printf("Continue? [y/N/stop]: ")

        var response string
        fmt.Scanln(&response)

        if response != "y" && response != "Y" {
          session.Status = "cancelled_by_user"
          return collaborationCancelledMsg{session, turn}
        }
      }

      // Send to AIs...
    }
  }
}
```

### Emergency Stop (Cancel Mid-Collaboration)

```go
// Keyboard shortcut to abort running collaboration
case "x", "X": // Emergency stop
  if m.activeCollaboration != nil {
    m.statusMsg = "🛑 Collaboration CANCELLED (tokens saved!)"

    // Log what was saved
    saved := estimateRemainingCost(m.activeCollaboration)
    m.statusMsg += fmt.Sprintf(" | Saved ~$%.2f", saved)

    m.activeCollaboration.Status = "cancelled"
    m.activeCollaboration = nil
    return m, nil
  }
```

### Real-Time Usage Tracking

```go
// Track requests during collaboration
type CollaborationUsage struct {
  TurnNumber      int
  RequestsUsed    int
  PerAI           map[AIProvider]int // Requests per AI
  QuotaRemaining  map[AIProvider]int // How much left per AI
}

func (session *CollaborationSession) trackUsage(turn CollaborationTurn) {
  // Increment request count for this AI
  session.RequestsUsed++
  session.PerAI[turn.AIProvider]++

  // Update remaining quota
  tier := UserSubscriptions[turn.AIProvider]
  newUsage := tier.CurrentUsage + session.PerAI[turn.AIProvider]
  session.QuotaRemaining[turn.AIProvider] = tier.UsageLimit - newUsage

  // Check if any AI is running low
  for provider, remaining := range session.QuotaRemaining {
    tier := UserSubscriptions[provider]
    percentUsed := (float64(tier.UsageLimit - remaining) / float64(tier.UsageLimit)) * 100

    if percentUsed > 90 {
      session.Status = "quota_warning"
      session.StatusMsg = fmt.Sprintf("⚠️  %s is at %.0f%% (resets %s)",
        aiProviderName(provider), percentUsed, formatResetTime(tier.ResetTime))
    }
  }
}

// Show in header during collaboration
func (m model) renderUsageWarning() string {
  if m.activeCollaboration == nil {
    return ""
  }

  var parts []string

  // Show usage for each AI being used
  for provider, used := range m.activeCollaboration.PerAI {
    tier := UserSubscriptions[provider]
    total := tier.UsageLimit
    current := tier.CurrentUsage + used
    percent := (float64(current) / float64(total)) * 100

    icon := getUsageIcon(percent)
    name := aiProviderName(provider)

    // Show remaining count and reset time if < 24h
    remaining := total - current
    timeUntil := time.Until(tier.ResetTime)

    if timeUntil < 24*time.Hour {
      parts = append(parts, fmt.Sprintf("%s %s: %d left (resets in %s)",
        icon, name, remaining, formatDuration(timeUntil)))
    } else {
      parts = append(parts, fmt.Sprintf("%s %s: %d/%d (%.0f%%)",
        icon, name, current, total, percent))
    }
  }

  return strings.Join(parts, " | ")
}
```

### Depth vs Breadth: Understanding Turn Limits

**Important distinction**:
- ✅ **Depth per turn** (reasoning effort): UNLIMITED - Let AIs think deeply!
- ⚠️ **Breadth** (number of turns): LIMITED - Prevent infinite loops

```go
// Codex can still use full reasoning effort per turn
const CodexDefaultFlags = "-m gpt-5 -c model_reasoning_effort=\"high\""

// The turn limit prevents:
// Turn 1: "What's wrong?"
// Turn 2: "Can you check file X?"
// Turn 3: "What about file Y?"
// Turn 4: "And file Z?"
// Turn 5: "Also check file A?"  ← This is the problem we're preventing
// ... (could loop forever asking follow-ups)

// Each individual turn can be as deep/expensive as needed:
// Codex with high reasoning: ~4000 tokens/turn ✅ ALLOWED
// GPT-5 extended thinking: ~8000 tokens/turn ✅ ALLOWED
// The limit is on NUMBER of back-and-forth exchanges, not thinking depth
```

### Conservative Presets

```go
// Safe defaults for different scenarios
var CollaborationPresets = map[string]CollaborationConfig{
  "quick": {
    Name:    "Quick Check (1 turn, 2 AIs)",
    Turns:   1,
    MaxAIs:  2,
    Timeout: 30 * time.Second,
    CodexFlags: "-m gpt-5 -c model_reasoning_effort=\"medium\"", // Still deep!
  },
  "standard": {
    Name:    "Standard Debug (2 turns, 3 AIs)",
    Turns:   2,
    MaxAIs:  3,
    Timeout: 60 * time.Second,
    CodexFlags: "-m gpt-5 -c model_reasoning_effort=\"high\"", // FULL reasoning ✅
  },
  "deep": {
    Name:    "Deep Analysis (3 turns, 3 AIs)",
    Turns:   3,
    MaxAIs:  3,
    Timeout: 120 * time.Second, // Longer timeout for deep thinking
    CodexFlags: "-m gpt-5 -c model_reasoning_effort=\"high\"", // FULL reasoning ✅
  },
  // NO "unlimited" option - always capped at 5 turns
}
```

### Per-AI Configuration

```go
// Each AI can have custom flags while respecting turn limits
type AIConfiguration struct {
  Provider AIProvider
  Command  string
  Flags    []string
  Timeout  time.Duration
}

var DefaultAIConfigs = map[AIProvider]AIConfiguration{
  AICodex: {
    Provider: AICodex,
    Command:  "codex",
    Flags: []string{
      "-m", "gpt-5",
      "-c", "model_reasoning_effort=\"high\"",  // ✅ Full deep reasoning
      "--sandbox", "read-only",                 // Safety
    },
    Timeout: 120 * time.Second, // Long timeout for deep thinking
  },
  AIGemini: {
    Provider: AIGemini,
    Command:  "gemini",
    Flags: []string{
      "-m", "gemini-2.0-flash-thinking-exp",   // ✅ Thinking mode
    },
    Timeout: 60 * time.Second,
  },
  AICursor: {
    Provider: AICursor,
    Command:  "cursor",
    Flags:    []string{}, // Cursor handles its own config
    Timeout:  90 * time.Second,
  },
}

// When sending to Codex, use full command:
// codex -m gpt-5 -c model_reasoning_effort="high" --sandbox read-only "Your question here"
//
// The turn limit (e.g. 2) means:
// - Turn 1: Codex thinks DEEPLY about initial question (uses full reasoning)
// - Turn 2: Codex thinks DEEPLY about follow-up (uses full reasoning)
// - Turn 3: STOPPED (prevents infinite loop)
//
// Each turn is expensive but thorough. We just limit how many times we go back and forth.
```

### Configuration File

```json
// ~/.config/tmuxplexer/collaboration.json
{
  "safety": {
    "max_turns": 3,
    "require_confirmation": true,
    "show_usage_estimates": true,
    "warn_at_percent": 75,      // Warn when any AI hits 75% usage
    "block_at_percent": 95      // Block if any AI would exceed 95%
  },
  "defaults": {
    "turns": 2,
    "timeout_seconds": 60,
    "preset": "standard"
  },
  "subscriptions": {
    "claude": {
      "tier": "max",
      "limit_type": "weekly",
      "limit": 100,               // 100 requests per week
      "reset_day": "sunday",
      "reset_time": "00:00"
    },
    "codex": {
      "tier": "pro",
      "limit_type": "monthly",
      "limit": 500,               // 500 requests per month
      "reset_day": "1st",
      "reset_time": "00:00"
    },
    "gemini": {
      "tier": "advanced",
      "limit_type": "daily",
      "limit": 1000,              // 1000 requests per day
      "reset_day": "daily",
      "reset_time": "00:00"
    }
  },
  "smart_usage": {
    "go_crazy_mode": {
      "enabled": true,
      "description": "When usage is low (<30%) and reset is soon (<24h), allow more aggressive usage",
      "threshold_percent": 30,
      "hours_until_reset": 24,
      "max_turns_boost": 5        // Allow up to 5 turns instead of default 2
    },
    "conservation_mode": {
      "enabled": true,
      "description": "When usage is high (>75%), be more conservative",
      "threshold_percent": 75,
      "max_turns_limit": 1,       // Only allow 1 turn
      "require_per_turn_approval": true
    }
  }
}
```

### UI Warnings

```go
// Before starting collaboration, show clear warning with usage impact
func (m model) renderCollaborationWarning(estimate UsageEstimate) string {
  var output strings.Builder

  output.WriteString("╔══════════════════════════════════════════════════════╗\n")
  output.WriteString("║      MULTI-AI COLLABORATION USAGE WARNING            ║\n")
  output.WriteString("╠══════════════════════════════════════════════════════╣\n")
  output.WriteString(fmt.Sprintf("║  Configuration: 3 AIs × %d turns = %d requests       ║\n",
    estimate.NumTurns, estimate.NumTurns * 3))
  output.WriteString("║                                                      ║\n")

  // Show impact per AI
  for provider, impact := range estimate.PerAI {
    name := aiProviderName(provider)
    icon := getUsageIcon(impact.PercentOfTotal)

    output.WriteString(fmt.Sprintf("║  %s %s:%-42s║\n", icon, name, ""))
    output.WriteString(fmt.Sprintf("║    Current: %d/%d (%.0f%% used)%-20s║\n",
      impact.CurrentUsage,
      impact.CurrentUsage + impact.RemainingBefore,
      (float64(impact.CurrentUsage) / float64(impact.CurrentUsage + impact.RemainingBefore)) * 100,
      ""))
    output.WriteString(fmt.Sprintf("║    After: %.0f%% used, %d requests left%-14s║\n",
      impact.PercentOfTotal,
      impact.RemainingAfter,
      ""))

    // Show reset time if < 24h
    timeUntil := time.Until(impact.ResetTime)
    if timeUntil < 24*time.Hour {
      output.WriteString(fmt.Sprintf("║    ⏰ Resets in %s%-32s║\n",
        formatDuration(timeUntil), ""))
    }
    output.WriteString("║                                                      ║\n")
  }

  output.WriteString("║  Safety limits:                                      ║\n")
  output.WriteString(fmt.Sprintf("║    • Max turns: %d (hard limit)%-24s║\n", MaxTurnsAllowed, ""))
  output.WriteString("║    • Emergency stop: Press 'x' anytime               ║\n")
  output.WriteString("║                                                      ║\n")
  output.WriteString("║  Press 'Y' to confirm, any other key to cancel       ║\n")
  output.WriteString("╚══════════════════════════════════════════════════════╝\n")

  return output.String()
}

// Example output:
//
// ╔══════════════════════════════════════════════════════╗
// ║      MULTI-AI COLLABORATION USAGE WARNING            ║
// ╠══════════════════════════════════════════════════════╣
// ║  Configuration: 3 AIs × 2 turns = 6 requests         ║
// ║                                                      ║
// ║  🟡 Claude Max:                                      ║
// ║    Current: 70/100 (70% used)                        ║
// ║    After: 76% used, 24 requests left                 ║
// ║    ⏰ Resets in 8 hours                              ║
// ║                                                      ║
// ║  🟢 Codex:                                           ║
// ║    Current: 120/500 (24% used)                       ║
// ║    After: 28% used, 358 requests left                ║
// ║                                                      ║
// ║  🟢 Gemini:                                          ║
// ║    Current: 234/1000 (23% used)                      ║
// ║    After: 26% used, 764 requests left                ║
// ║    ⏰ Resets in 14 hours                             ║
// ║                                                      ║
// ║  Safety limits:                                      ║
// ║    • Max turns: 5 (hard limit)                       ║
// ║    • Emergency stop: Press 'x' anytime               ║
// ║                                                      ║
// ║  Press 'Y' to confirm, any other key to cancel       ║
// ╚══════════════════════════════════════════════════════╝
```

### Token Counting by Provider

```go
// Extract actual token usage from each AI's response
func extractTokenCount(response string, provider AIProvider) int {
  switch provider {
  case AICodex:
    // Codex shows: "tokens used: 1234"
    if match := regexp.MustCompile(`tokens used: (\d+)`).FindStringSubmatch(response); len(match) > 1 {
      tokens, _ := strconv.Atoi(match[1])
      return tokens
    }

  case AIGemini:
    // Gemini shows: "Token count: 1234"
    if match := regexp.MustCompile(`Token count: (\d+)`).FindStringSubmatch(response); len(match) > 1 {
      tokens, _ := strconv.Atoi(match[1])
      return tokens
    }

  case AICursor:
    // Cursor may not show - estimate based on response length
    return estimateTokens(response)
  }

  // Fallback: estimate
  return estimateTokens(response)
}

func estimateTokens(text string) int {
  // Rough estimate: ~1.3 tokens per word
  words := strings.Fields(text)
  return int(float64(len(words)) * 1.3)
}

// Calculate cost per provider
func calculateProviderCost(tokens int, provider AIProvider) float64 {
  costPer1k := map[AIProvider]float64{
    AICodex:   0.020,  // $0.02 per 1k tokens (GPT-5)
    AIGemini:  0.001,  // $0.001 per 1k tokens (Flash)
    AICursor:  0.010,  // $0.01 per 1k tokens (Claude)
  }

  rate := costPer1k[provider]
  return (float64(tokens) / 1000.0) * rate
}
```

### Smart Usage Modes (Adaptive Behavior)

```go
// Automatically adjust limits based on usage situation
func (m model) getSmartTurnLimit(providers []AIProvider) int {
  config := loadCollaborationConfig()

  for _, provider := range providers {
    tier := UserSubscriptions[provider]
    percentUsed := (float64(tier.CurrentUsage) / float64(tier.UsageLimit)) * 100
    hoursUntilReset := time.Until(tier.ResetTime).Hours()

    // GO CRAZY MODE: Low usage + about to reset = max it out!
    if config.SmartUsage.GoCrazyMode.Enabled {
      if percentUsed < float64(config.SmartUsage.GoCrazyMode.ThresholdPercent) &&
         hoursUntilReset < float64(config.SmartUsage.GoCrazyMode.HoursUntilReset) {

        return config.SmartUsage.GoCrazyMode.MaxTurnsBoost
      }
    }

    // CONSERVATION MODE: High usage = be careful
    if config.SmartUsage.ConservationMode.Enabled {
      if percentUsed > float64(config.SmartUsage.ConservationMode.ThresholdPercent) {
        return config.SmartUsage.ConservationMode.MaxTurnsLimit
      }
    }
  }

  // Default
  return config.Defaults.Turns
}

// Example: User's scenario (30% used, 8 hours until Sunday midnight reset)
//
// Claude Max:
//   Usage: 70/100 (70% used)
//   Reset: Sunday 00:00 (8 hours from now)
//
// Smart mode kicks in:
//   ✅ Usage < 30%? NO (70% used)
//   ❌ Conservation mode triggered (>75%? YES at 70%)
//
// Actually wait... user said they have 30% REMAINING, so 70% used
// Let me recalculate:
//
// If 30% remaining = 30/100 requests left
// Then 70/100 used (70% usage)
//
// Hmm, but user is going crazy with 5-6 terminals, so maybe they're at 70% used
// and feeling safe to burn the last 30% before midnight reset?
//
// Let's document both scenarios:
```

### Real-World Usage Scenarios

**Scenario 1: "Go Crazy Mode" (User's current situation)**
```
Claude Max: 70/100 used (30 remaining), resets in 8 hours

Thinking: "Reset is tonight, I still have 30% left, let's use it!"

Smart mode:
  - 70% used = Close to threshold but not conservation mode yet
  - 8 hours until reset = Short window
  - Decision: Allow 3-5 turns (more than default 2)
  - Rationale: Better to use it than lose it!

User behavior: Opens 5-6 tmux terminals with different AIs,
               runs multiple collaborations in parallel
```

**Scenario 2: "Conservation Mode" (Monday morning after reset)**
```
Claude Max: 15/100 used (85 remaining), resets in 6 days

Thinking: "Just reset, I have the whole week, be conservative"

Smart mode:
  - 15% used = Plenty remaining
  - 6 days until reset = Long window
  - Decision: Default 2 turns, no special treatment
  - Rationale: Pace yourself for the week
```

**Scenario 3: "Conservation Mode" (Friday afternoon)**
```
Claude Max: 92/100 used (8 remaining), resets in 2 days

Thinking: "I'm almost out and it's only Friday!"

Smart mode:
  - 92% used = HIGH USAGE, conservation mode triggered
  - 2 days until reset = Still time
  - Decision: Max 1 turn, require approval for each
  - Warning: "⚠️  Only 8 requests left for 2 days!"
  - Rationale: Don't get throttled before the weekend
```

**Scenario 4: "Last Hour Yolo Mode"**
```
Claude Max: 85/100 used (15 remaining), resets in 1 hour

Thinking: "Reset in 1 hour, burn what's left!"

Smart mode:
  - 85% used but doesn't matter
  - <1 hour until reset = YOLO TIME
  - Decision: Remove limits entirely, max out at 5 turns
  - Message: "🔥 Go crazy! Resets in 47 minutes"
  - Rationale: Use it or lose it!
```

## Security Considerations

1. **API Keys**: Don't log API keys from AI responses
2. **Sensitive Data**: Warning if sending code with secrets
3. **Local Only**: All AIs run locally (Codex, Gemini CLIs)
4. **Sandboxing**: Consider read-only mode for AI tools
5. **User Approval**: Optionally require approval before sending each turn
6. **Usage Limits**: Respect subscription tier limits and prevent quota exhaustion
7. **Emergency Stop**: Always allow user to cancel mid-collaboration

## Success Metrics

1. **Time to Debug**: Measure reduction in debugging time
2. **Consensus Accuracy**: How often multi-AI agrees vs single AI
3. **User Satisfaction**: Survey developers using the feature
4. **Adoption Rate**: % of tmuxplexer users enabling collaboration
5. **Cost Efficiency**: Compare token usage multi-turn vs multiple single queries

## Alternative Approaches Considered

### 1. Slash Commands Only

**Pros**: Simple, works today
**Cons**: Manual orchestration, no state persistence, no visual monitoring

### 2. External Orchestration Tool

**Pros**: Language-agnostic, more flexible
**Cons**: Adds complexity, breaks tmux workflow

### 3. Custom AI Wrapper

**Pros**: More control over AI interactions
**Cons**: Reinventing wheel, maintenance burden

### 4. VS Code Extension

**Pros**: Rich UI, existing ecosystem
**Cons**: Not terminal-native, doesn't fit tmux workflow

**Verdict**: Tmuxplexer is the perfect fit because:
- Already captures pane output
- Visual layout shows all AIs
- Terminal-native workflow
- Integrates with other tools (logs, metrics, errors)

## Open Questions

1. How to handle API rate limits across multiple AIs?
2. Should we cache AI responses to avoid redundant queries?
3. How to version control collaboration sessions (save to git)?
4. Should we support non-CLI AIs (web APIs, language models)?
5. How to visualize disagreement between AIs?
6. Should Claude automatically trigger collaboration on repeated failures?

## Related Projects

- **OpenAI Codex**: Deep reasoning AI
- **Google Gemini**: Fast, multimodal AI
- **Cursor**: VS Code with Claude integration
- **Aider**: AI pair programming in terminal
- **Smol Developer**: Multi-agent coding assistant

---

**Status**: Ready for Phase 0 (Proof of Concept)
**Next Step**: Test sending input to panes manually, validate detection works
**Estimated Time to MVP**:
- Phase 1 (Send Input): 2-3 hours
- Phase 2 (Response Detection): 4-6 hours
- Phase 3 (Orchestrator): 8-12 hours
- **Total**: ~16-20 hours to working multi-AI collaboration
