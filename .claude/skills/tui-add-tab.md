# TUI Add Tab Skill

Add a new tab to your TUI application's tabbed layout with full integration.

## Usage

When you want to add a new tab to your tabbed layout:

```
/tui-add-tab TabName "Tab Display Title"
```

Example:
```
/tui-add-tab settings "Settings"
/tui-add-tab logs "System Logs"
```

## What I'll Do

I will add a new tab to your TUI application by updating all necessary files:

### 1. Update Type Definitions (`types.go` or `model.go`)

- Add the new tab constant to your tab enum/constants
- Add tab state tracking fields if needed
- Update any tab-related data structures

Example:
```go
const (
    TabHome = iota
    TabFiles
    TabSettings  // NEW TAB
)
```

### 2. Update View Rendering (`view.go`)

- Add the new tab to the tab bar rendering logic
- Implement the tab's content rendering function
- Apply text truncation to prevent wrapping (CRITICAL!)
- Account for borders in height calculations (CRITICAL!)

Example:
```go
func (m model) renderTabContent() string {
    switch m.currentTab {
    case TabHome:
        return m.renderHomeTab()
    case TabFiles:
        return m.renderFilesTab()
    case TabSettings:  // NEW TAB
        return m.renderSettingsTab()
    default:
        return ""
    }
}

func (m model) renderSettingsTab() string {
    // Calculate available space
    maxTextWidth := width - 4  // -2 for borders, -2 for padding

    // Truncate ALL text before rendering
    title := truncateString("Settings", maxTextWidth)

    // Build content with proper height accounting
    // contentHeight should already account for borders (-2)

    return panel
}
```

### 3. Update Keyboard Handling (`update_keyboard.go`)

- Add keyboard shortcut to switch to the new tab
- Update tab navigation logic (next/prev tab)
- Add any tab-specific keyboard shortcuts

Example:
```go
case "3":
    m.currentTab = TabSettings  // NEW TAB
    return m, nil
```

### 4. Update Mouse Handling (`update_mouse.go`)

- Add click detection for the new tab in the tab bar
- Calculate the tab's position and width
- Handle tab selection via mouse click

Example:
```go
func (m model) handleTabBarClick(x, y int) (tea.Model, tea.Cmd) {
    tabNames := []string{"Home", "Files", "Settings"}  // NEW TAB
    xPos := 2  // Starting position

    for i, name := range tabNames {
        tabWidth := len(name) + 4  // "[ " + name + " ]"

        if x >= xPos && x < xPos+tabWidth {
            m.currentTab = i
            return m, nil
        }

        xPos += tabWidth + 1  // +1 for space between tabs
    }

    return m, nil
}
```

### 5. Update Help Text (if applicable)

- Add the new keyboard shortcut to the help display
- Document the tab's purpose

## Critical Checks

Before considering the task complete, I will verify:

### ✅ Text Truncation
- [ ] All text in the tab (title, subtitle, content) is truncated to prevent wrapping
- [ ] Max text width calculated as: `width - 4` (2 for borders, 2 for padding)
- [ ] Using truncateString() helper or equivalent for all strings

**Why:** Text wrapping in bordered panels causes misalignment (see CLAUDE.md Issue 2)

### ✅ Border Height Accounting
- [ ] Content height calculation subtracts border height: `contentHeight -= 2`
- [ ] No explicit `Height()` set on Lipgloss bordered styles
- [ ] Content fills exact height with empty lines if needed

**Why:** Missing border accounting causes panels to overflow and cover the title bar (see CLAUDE.md Issue 1)

### ✅ Mouse Click Detection
- [ ] Tab click region properly calculated with exact x/y coordinates
- [ ] Tab width includes padding: `len(name) + 4` for "[ name ]" format
- [ ] Click detection uses >= and < (not > and <=)

**Why:** Incorrect boundaries cause clicks to miss or select wrong tab

### ✅ Keyboard Navigation
- [ ] Direct shortcut added (e.g., "3" for third tab)
- [ ] Tab cycling works (next/prev with arrow keys or tab/shift+tab)
- [ ] All tab constants properly defined

### ✅ Responsive Layout
- [ ] Tab works in all terminal sizes
- [ ] Tab bar wraps or scrolls gracefully if too many tabs
- [ ] Content respects minimum terminal size checks

## Files Modified

Typical files that will be changed:

- `types.go` or `model.go` - Tab constants and state
- `view.go` - Tab rendering logic (lines vary by implementation)
- `update_keyboard.go` - Keyboard shortcuts
- `update_mouse.go` - Tab bar click detection
- `styles.go` - Tab styles (if new styles needed)

## Example Output

After running this skill, you'll have:

1. A new tab constant defined
2. A rendering function for the tab's content
3. Keyboard shortcut (e.g., "3") to switch to the tab
4. Mouse click support in the tab bar
5. Proper text truncation and border accounting applied

## Best Practices Applied

- **Text Truncation**: No auto-wrapping, explicit truncation before render
- **Border Accounting**: Height calculations subtract 2 for panel borders
- **Large Click Targets**: Tab buttons are at least 12 chars wide for easy clicking
- **Visual Feedback**: Tab shows different style when active/selected
- **Keyboard First**: Direct number keys for quick access

## Reference

- CLAUDE.md - Critical layout fixes
- MOUSE_SUPPORT_GUIDE.md - Click detection patterns
- Template files: `template/*.go.tmpl`
- Example: `examples/multi-panel/`

## Notes

- This skill enforces the patterns documented in CLAUDE.md
- All text MUST be truncated before rendering to prevent wrapping issues
- Border height MUST be accounted for in all calculations
- Mouse detection MUST match the actual rendered layout
