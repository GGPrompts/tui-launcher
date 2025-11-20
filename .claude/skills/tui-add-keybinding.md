# TUI Add Keybinding Skill

Add a new keyboard shortcut to your TUI application with proper help documentation.

## Usage

When you want to add a new keyboard shortcut:

```
/tui-add-keybinding "key" "action" "description"
```

Examples:
```
/tui-add-keybinding "r" "refresh" "Refresh data"
/tui-add-keybinding "ctrl+s" "save" "Save current state"
/tui-add-keybinding "/" "search" "Search items"
```

## What I'll Do

I will add a complete keyboard binding by updating the key handler, help text, and any related UI.

### 1. Identify Key Handler Location (`update_keyboard.go` or `update.go`)

Find the main keyboard handler:

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    // ... other cases
    }
}
```

### 2. Add Key Handler (`update_keyboard.go`)

Add the new key binding to the appropriate handler:

```go
func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Global keys (work everywhere)
    switch msg.String() {
    case "q", "ctrl+c":
        return m, tea.Quit

    case "r":  // NEW KEYBINDING
        return m.handleRefresh()

    case "ctrl+s":  // NEW KEYBINDING
        return m.handleSave()

    case "/":  // NEW KEYBINDING
        return m.handleSearch()
    }

    // Context-specific keys
    return m.handleContextKeys(msg)
}
```

**Pattern:** Short single keys ("r", "a") for common actions, Ctrl+ for important actions.

### 3. Implement Action Handler

Create the function that executes the action:

```go
func (m model) handleRefresh() (tea.Model, tea.Cmd) {
    // Perform the refresh action
    m.status = "Refreshing..."

    // Return a command if needed (for async operations)
    return m, func() tea.Msg {
        // Do async work here
        return refreshCompleteMsg{}
    }
}

func (m model) handleSave() (tea.Model, tea.Cmd) {
    m.status = "Saved!"
    // Save logic here
    return m, nil
}

func (m model) handleSearch() (tea.Model, tea.Cmd) {
    m.searchMode = true
    m.searchQuery = ""
    return m, nil
}
```

### 4. Define Help Key (Optional - If Using Key Mapping System)

If your app uses a structured key map (like Bubbletea's `key` package):

```go
import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
    Quit    key.Binding
    Refresh key.Binding  // NEW
    Save    key.Binding  // NEW
    Search  key.Binding  // NEW
}

var keys = keyMap{
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
    Refresh: key.NewBinding(  // NEW
        key.WithKeys("r"),
        key.WithHelp("r", "refresh"),
    ),
    Save: key.NewBinding(  // NEW
        key.WithKeys("ctrl+s"),
        key.WithHelp("ctrl+s", "save"),
    ),
    Search: key.NewBinding(  // NEW
        key.WithKeys("/"),
        key.WithHelp("/", "search"),
    ),
}

// In Update():
case key.Matches(msg, keys.Refresh):  // NEW
    return m.handleRefresh()
```

### 5. Update Help Text (`view.go` or `help.go`)

Add the keybinding to the help display:

#### Simple Help String
```go
func (m model) renderHelp() string {
    help := []string{
        "q: quit",
        "r: refresh",      // NEW
        "ctrl+s: save",    // NEW
        "/: search",       // NEW
        "?: toggle help",
    }
    return dimStyle.Render(strings.Join(help, " | "))
}
```

#### Structured Help (Using Bubbles Help)
```go
import "github.com/charmbracelet/bubbles/help"

func (m model) renderHelp() string {
    h := help.New()

    keys := []key.Binding{
        m.keys.Quit,
        m.keys.Refresh,  // NEW
        m.keys.Save,     // NEW
        m.keys.Search,   // NEW
    }

    return h.View(keys)
}
```

#### Categorized Help
```go
func (m model) renderFullHelp() string {
    help := []string{
        "General:",
        "  q       - Quit",
        "  r       - Refresh",     // NEW
        "  ctrl+s  - Save",         // NEW
        "",
        "Navigation:",
        "  /       - Search",       // NEW
        "  h/l     - Switch panels",
    }

    return lipgloss.JoinVertical(
        lipgloss.Left,
        help...,
    )
}
```

### 6. Add Visual Feedback (Optional)

Show the action result in the status bar:

```go
func (m model) handleRefresh() (tea.Model, tea.Cmd) {
    m.lastAction = "Refreshing..."  // Show in status bar
    m.lastActionTime = time.Now()

    return m, func() tea.Msg {
        // Perform refresh
        time.Sleep(1 * time.Second)
        return refreshCompleteMsg{success: true}
    }
}

// In view.go:
func (m model) renderStatusBar() string {
    status := "Ready"

    if time.Since(m.lastActionTime) < 2*time.Second {
        status = m.lastAction  // Show recent action
    }

    return statusStyle.Render(status)
}
```

### 7. Handle Conflicts

Check for and resolve any key conflicts:

```go
// Good - No conflicts
case "r":
    return m.handleRefresh()

case "R":  // Capital R is different
    return m.handleReload()

// Bad - Conflict!
case "r":
    return m.handleRefresh()
// ...later in code...
case "r":  // ← CONFLICT! This won't execute
    return m.handleReload()
```

**Resolution:** Use different keys or combine with modifiers:
- `r` - Refresh
- `R` (shift+r) - Reload all
- `ctrl+r` - Reset

### 8. Modal/Contextual Keys

Add keys that only work in specific modes:

```go
func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // Search mode keys
    if m.searchMode {
        switch msg.String() {
        case "esc":
            m.searchMode = false
            return m, nil
        case "enter":
            return m.executeSearch()
        default:
            // Add to search query
            m.searchQuery += msg.String()
            return m, nil
        }
    }

    // Normal mode keys
    switch msg.String() {
    case "/":
        m.searchMode = true
        return m, nil
    // ... other keys
    }
}
```

### 9. Add to README/Documentation (Optional)

Update user documentation:

```markdown
## Keyboard Shortcuts

### General
- `q` or `Ctrl+C` - Quit
- `r` - Refresh data
- `Ctrl+S` - Save current state
- `?` - Toggle help

### Search
- `/` - Enter search mode
- `Esc` - Exit search mode
- `Enter` - Execute search
```

## Critical Checks

Before considering the keybinding complete, verify:

### ✅ Key Handler Added
- [ ] Key case added to switch statement
- [ ] Handler function implemented
- [ ] No conflicts with existing keys
- [ ] Uses appropriate modifier (none, ctrl, shift) for importance

### ✅ Action Implementation
- [ ] Handler function exists and is callable
- [ ] Returns proper `(tea.Model, tea.Cmd)` types
- [ ] Updates model state appropriately
- [ ] Returns commands for async operations if needed

### ✅ Help Text Updated
- [ ] Key added to help display
- [ ] Description is clear and concise
- [ ] Follows existing help format
- [ ] Help is visible (not truncated)

### ✅ Visual Feedback
- [ ] Status bar shows action result (if applicable)
- [ ] Error states handled (if applicable)
- [ ] Loading states shown (if async)
- [ ] User knows the action was triggered

### ✅ No Conflicts
- [ ] Key doesn't conflict with existing bindings
- [ ] Key doesn't conflict with terminal shortcuts
- [ ] Key works in intended contexts only

### ✅ Testing
- [ ] Key triggers correct action
- [ ] Works in all intended modes/contexts
- [ ] Doesn't trigger in wrong contexts
- [ ] Help text displays correctly

## Common Key Patterns

### Single Letter Keys
**Use for:** Frequent actions
```
h, j, k, l  - Navigation
r           - Refresh
/           - Search
?           - Help
q           - Quit
```

### Ctrl + Key
**Use for:** Important actions, editor commands
```
Ctrl+C  - Quit/Cancel
Ctrl+S  - Save
Ctrl+R  - Reload
Ctrl+F  - Find
```

### Shift + Key (Capital Letters)
**Use for:** Reverse actions, global versions
```
J  - Jump to bottom (vs j = down one)
H  - Go to first panel (vs h = left one)
R  - Reload all (vs r = refresh current)
```

### Special Keys
**Use for:** Modal operations
```
Enter  - Confirm/Select
Esc    - Cancel/Exit mode
Tab    - Next field/panel
Space  - Toggle/Select
```

### Number Keys
**Use for:** Direct selection
```
1-9  - Select tab/layout
```

## Best Practices

### ✅ DO
- Use intuitive keys (h/j/k/l for navigation)
- Document all keys in help
- Provide visual feedback for actions
- Use modifiers (Ctrl) for destructive actions
- Group related keys logically

### ❌ DON'T
- Override standard terminal shortcuts (Ctrl+C, Ctrl+Z)
- Use obscure key combinations
- Hide important keys from help text
- Create conflicting bindings
- Forget to handle uppercase separately

## Example: Complete Keybinding Addition

**Task:** Add "r" to refresh data

**Files Modified:**

1. `update_keyboard.go`:
```go
case "r":
    return m.handleRefresh()
```

2. `model.go`:
```go
func (m model) handleRefresh() (tea.Model, tea.Cmd) {
    m.status = "Refreshing..."
    return m, refreshDataCmd()
}

func refreshDataCmd() tea.Cmd {
    return func() tea.Msg {
        // Load new data
        data := loadData()
        return dataRefreshedMsg{data: data}
    }
}
```

3. `update.go`:
```go
case dataRefreshedMsg:
    m.data = msg.data
    m.status = "Refreshed!"
    return m, nil
```

4. `view.go`:
```go
func (m model) renderHelp() string {
    return dimStyle.Render("q: quit | r: refresh | ?: help")
}
```

**Result:** Pressing "r" refreshes data, shows "Refreshing..." then "Refreshed!" in status bar.

## Reference

- **Bubbletea Docs:** https://github.com/charmbracelet/bubbletea
- **Bubbles Key Package:** https://github.com/charmbracelet/bubbles/tree/master/key
- **Template Examples:** `examples/*/update_keyboard.go`

## Files Modified

Typical files that will be changed:

- `update_keyboard.go` - Add key case (~3-5 lines)
- `model.go` or `handlers.go` - Add handler function (~10-30 lines)
- `view.go` - Update help text (~1-2 lines)
- `types.go` - Add key map definition if using structured approach (~5-10 lines)
- Optional: `README.md` - Document new keybinding

## Notes

- Single letter keys are easiest to remember
- Ctrl+ combinations feel more "important"
- Always update help text - users won't find hidden keys
- Test in different modes to avoid conflicts
- Consider vim-style navigation (hjkl) for familiarity
