# Bug Fix: Config Loading Stuck on "Loading configuration..."

**Issue:** The launcher was stuck showing the loading spinner and never displayed the config.

**Root Cause:** The unified model's `Update()` method was not routing unknown messages (like `configLoadedMsg`) to the active tab.

## The Problem

When the launch tab's `loadConfig()` function returned a `configLoadedMsg`, the message flow was:

1. `loadConfig()` returns `configLoadedMsg`
2. Message goes to `unifiedModel.Update()`
3. `unifiedModel.Update()` only handled `tea.KeyMsg`, `tea.WindowSizeMsg`, and `spinner.TickMsg`
4. **`configLoadedMsg` was NOT handled** - fell through to `return m, nil`
5. Message never reached `launch.Model.Update()` where it would be processed
6. Launch tab's `loading` flag stayed `true` forever
7. Spinner kept showing "Loading configuration..."

## The Fix

### 1. Added default case to route unknown messages (`model_unified.go`)

```go
// Before:
case spinner.TickMsg:
    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg)
    return m, cmd
}

return m, nil

// After:
case spinner.TickMsg:
    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg)
    return m, cmd

default:
    // Route all other messages to the active tab
    // This includes tab-specific messages like configLoadedMsg
    cmd := routeUpdateToTab(&m, msg)
    return m, cmd
}

return m, nil
```

### 2. Removed loading check in unified View (`model_unified.go`)

The unified model was checking its own `m.loading` flag, which was never set to false. Each tab should handle its own loading state.

```go
// Before:
func (m unifiedModel) View() string {
    if m.loading {
        return m.spinner.View() + " Loading configuration...\n"
    }
    // ...

// After:
func (m unifiedModel) View() string {
    // Don't check m.loading here - let each tab handle its own loading state
    // The launch tab will show its own loading spinner
    // ...
```

### 3. Added debug logging (for testing)

Added debug output in `tabs/launch/model.go` to track config loading:

```go
func loadConfig() tea.Msg {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Fprintf(os.Stderr, "DEBUG: Error getting home dir: %v\n", err)
        return configLoadedMsg{err: err}
    }

    configPath := filepath.Join(homeDir, ".config", "tui-launcher", "config.yaml")
    fmt.Fprintf(os.Stderr, "DEBUG: Loading config from: %s\n", configPath)

    data, err := os.ReadFile(configPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "DEBUG: Error reading config: %v\n", err)
        return configLoadedMsg{err: err}
    }

    fmt.Fprintf(os.Stderr, "DEBUG: Config file read, size: %d bytes\n", len(data))
    // ... rest of function
}
```

## Testing

When you run `./tui-launcher` now, you should see debug output like:

```
DEBUG: Loading config from: /home/matt/.config/tui-launcher/config.yaml
DEBUG: Config file read, size: 4897 bytes
DEBUG: Config parsed successfully - Projects: 3, Tools: 4
```

Then the TUI should display the actual Launch tab interface (not the loading spinner).

## What This Teaches Us

**Key Lesson:** When building a message-routing architecture like this, you need a **default case** to handle unknown messages. Otherwise, tab-specific messages get dropped.

**Pattern to follow:**
```go
func (m unifiedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle global keys
    case tea.WindowSizeMsg:
        // Handle resize
    default:
        // IMPORTANT: Route everything else to active tab
        cmd := routeUpdateToTab(&m, msg)
        return m, cmd
    }
}
```

## Files Modified

1. **`model_unified.go`**
   - Added `default:` case to route unknown messages
   - Removed loading check in View()

2. **`tabs/launch/model.go`**
   - Added debug logging to `loadConfig()`

3. **`tabs/launch/update.go`**
   - Added error logging for failed config load

## Next Steps

Once you verify config loads correctly in a real terminal:

1. Remove debug logging (or wrap in `if DEBUG` flag)
2. Test navigation, selection, and spawning
3. Verify Quick CD works
4. Move on to Sessions tab implementation

---

**Status:** ✅ Fixed and ready for testing
**Build:** ✅ Compiles successfully
**Ready for:** Real terminal testing
