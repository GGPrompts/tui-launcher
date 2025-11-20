# TUI New Layout Mode Skill

Create a completely new layout mode for your TUI application with full plumbing.

## Usage

When you want to add a new layout mode (e.g., "triple-split", "dashboard", "grid"):

```
/tui-new-layout LayoutName "Description"
```

Examples:
```
/tui-new-layout triple "Three vertical panels"
/tui-new-layout grid "2x2 grid layout"
/tui-new-layout dashboard "Dashboard with widgets"
```

## What I'll Do

I will create a complete new layout mode by following the established pattern from existing layouts (single/dual/multi/tabbed).

### 1. Define Layout Constant (`types.go` or `model.go`)

Add the new layout to your layout enum:

```go
type LayoutMode int

const (
    LayoutSingle LayoutMode = iota
    LayoutDual
    LayoutMulti
    LayoutTabbed
    LayoutTriple  // NEW LAYOUT
)

// Add display name
func (l LayoutMode) String() string {
    switch l {
    case LayoutSingle:
        return "Single"
    case LayoutDual:
        return "Dual"
    case LayoutMulti:
        return "Multi"
    case LayoutTabbed:
        return "Tabbed"
    case LayoutTriple:  // NEW LAYOUT
        return "Triple"
    default:
        return "Unknown"
    }
}
```

### 2. Add Layout Calculation Function (`model.go`)

Create a function to calculate panel dimensions for the new layout:

```go
// calculateTripleLayout returns dimensions for three vertical panels
func (m model) calculateTripleLayout() (leftWidth, middleWidth, rightWidth, height int) {
    totalWidth := m.width
    contentHeight := m.height

    // Subtract UI elements (CRITICAL: includes borders!)
    if m.config.UI.ShowTitle {
        contentHeight -= 3
    }
    if m.config.UI.ShowStatus {
        contentHeight -= 1
    }
    contentHeight -= 2  // CRITICAL: Panel borders

    // Divide width into three panels
    // Left: 25%, Middle: 50%, Right: 25%
    leftWidth = totalWidth / 4
    rightWidth = totalWidth / 4
    middleWidth = totalWidth - leftWidth - rightWidth

    return leftWidth, middleWidth, rightWidth, contentHeight
}
```

**Key Pattern:** All layout calculation functions must:
- Start with total dimensions
- Subtract title bar height (if shown)
- Subtract status bar height (if shown)
- **CRITICAL:** Subtract border height (2 lines)
- Return individual panel dimensions

### 3. Create Rendering Function (`view.go`)

Add rendering logic for the new layout:

```go
func (m model) renderTripleLayout() string {
    leftWidth, middleWidth, rightWidth, height := m.calculateTripleLayout()

    // Calculate max text width for truncation (CRITICAL!)
    maxTextLeft := leftWidth - 4
    maxTextMiddle := middleWidth - 4
    maxTextRight := rightWidth - 4

    // Render each panel with truncation
    leftPanel := m.renderPanel(
        truncateString("Left Panel", maxTextLeft),
        truncateString("Info", maxTextLeft),
        m.leftContent,
        leftWidth,
        height,
        maxTextLeft,
        m.focusedPanel == "left",
    )

    middlePanel := m.renderPanel(
        truncateString("Middle Panel", maxTextMiddle),
        truncateString("Main Content", maxTextMiddle),
        m.middleContent,
        middleWidth,
        height,
        maxTextMiddle,
        m.focusedPanel == "middle",
    )

    rightPanel := m.renderPanel(
        truncateString("Right Panel", maxTextRight),
        truncateString("Details", maxTextRight),
        m.rightContent,
        rightWidth,
        height,
        maxTextRight,
        m.focusedPanel == "right",
    )

    // Join panels horizontally
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        leftPanel,
        middlePanel,
        rightPanel,
    )
}

// Generic panel renderer (CRITICAL: proper height handling)
func (m model) renderPanel(title, subtitle string, content []string, width, height, maxTextWidth int, focused bool) string {
    // Height calculation
    innerHeight := height - 4  // -2 for title/subtitle, -2 for borders (handled by content)

    // Build content lines with truncation
    lines := []string{}
    for i := 0; i < innerHeight && i < len(content); i++ {
        line := truncateString(content[i], maxTextWidth)
        lines.append(lines, line)
    }

    // Fill remaining height with empty lines
    for len(lines) < innerHeight {
        lines = append(lines, "")
    }

    // Render with border (borders add 2 to height)
    style := panelStyle
    if focused {
        style = focusedPanelStyle
    }

    panel := lipgloss.JoinVertical(
        lipgloss.Left,
        titleStyle.Render(title),
        subtitleStyle.Render(subtitle),
        strings.Join(lines, "\n"),
    )

    return style.Render(panel)
}
```

### 4. Update Main View Router (`view.go`)

Add the new layout to the main rendering switch:

```go
func (m model) renderMainContent() string {
    switch m.currentLayout {
    case LayoutSingle:
        return m.renderSingleLayout()
    case LayoutDual:
        return m.renderDualLayout()
    case LayoutMulti:
        return m.renderMultiLayout()
    case LayoutTabbed:
        return m.renderTabbedLayout()
    case LayoutTriple:  // NEW LAYOUT
        return m.renderTripleLayout()
    default:
        return "Unknown layout"
    }
}
```

### 5. Add Keyboard Shortcuts (`update_keyboard.go`)

Add shortcut to switch to the new layout:

```go
func (m model) handleLayoutSwitch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "1":
        m.currentLayout = LayoutSingle
    case "2":
        m.currentLayout = LayoutDual
    case "3":
        m.currentLayout = LayoutMulti
    case "4":
        m.currentLayout = LayoutTabbed
    case "5":  // NEW LAYOUT
        m.currentLayout = LayoutTriple
        m.focusedPanel = "middle"  // Default focus
    }
    return m, nil
}
```

Add focus switching for the new panels:

```go
func (m model) handlePanelFocus(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    if m.currentLayout == LayoutTriple {
        switch msg.String() {
        case "h", "left":
            m.focusedPanel = "left"
        case "l", "right":
            m.focusedPanel = "right"
        case "k", "up":
            m.focusedPanel = "middle"
        }
    }
    return m, nil
}
```

### 6. Add Mouse Support (`update_mouse.go`)

Add click detection for the new layout:

```go
func (m model) handleLeftClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    if m.currentLayout == LayoutTriple {
        return m.handleTripleLayoutClick(msg)
    }
    // ... other layouts
}

func (m model) handleTripleLayoutClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    // Calculate panel boundaries
    leftWidth, middleWidth, rightWidth, _ := m.calculateTripleLayout()

    // Determine content start Y
    contentStartY := 0
    if m.config.UI.ShowTitle {
        contentStartY += 3
    }

    // Check if click is in content area
    if msg.Y < contentStartY || msg.Y >= m.height-1 {
        return m, nil  // Click in title/status area
    }

    // Detect which panel was clicked (X coordinate)
    if msg.X < leftWidth {
        m.focusedPanel = "left"
    } else if msg.X < leftWidth + middleWidth {
        m.focusedPanel = "middle"
    } else {
        m.focusedPanel = "right"
    }

    return m, nil
}
```

**Key Pattern:** Mouse detection must:
- Calculate exact panel boundaries
- Account for title bar offset
- Use correct coordinate (X for horizontal, Y for vertical)
- Check boundaries with >= and < (not > and <=)

### 7. Update Help Text (`view.go` or `help.go`)

Document the new layout and its shortcuts:

```go
func (m model) renderHelp() string {
    help := []string{
        "Layouts: 1=Single 2=Dual 3=Multi 4=Tabbed 5=Triple",
        "Navigate: h/l (left/right) k (middle)",
        "Quit: q or Ctrl+C",
    }
    return strings.Join(help, " | ")
}
```

### 8. Add Initial State (`model.go`)

Initialize any layout-specific state:

```go
func initialModel() model {
    return model{
        currentLayout: LayoutTriple,  // Start with new layout
        focusedPanel:  "middle",      // Default focus
        leftContent:   []string{"Left item 1", "Left item 2"},
        middleContent: []string{"Main content line 1", "Main content line 2"},
        rightContent:  []string{"Detail 1", "Detail 2"},
        // ... other fields
    }
}
```

## Critical Checks

Before considering the layout complete, I will verify:

### ✅ Layout Calculation
- [ ] Function exists: `calculateXXXLayout()`
- [ ] Subtracts title height (3 lines if shown)
- [ ] Subtracts status height (1 line if shown)
- [ ] **CRITICAL:** Subtracts border height (2 lines)
- [ ] Returns all panel dimensions
- [ ] Math adds up: `sum(panel widths/heights) = total - UI elements`

### ✅ Panel Rendering
- [ ] All text truncated before rendering
- [ ] Max text width: `panelWidth - 4`
- [ ] Content fills exact height (padding with empty lines)
- [ ] No explicit `Height()` on bordered styles
- [ ] Borders rendered naturally by Lipgloss

### ✅ Layout Routing
- [ ] New layout added to main switch statement
- [ ] Renders correctly via `renderXXXLayout()`
- [ ] Default case handles unknown layouts

### ✅ Keyboard Support
- [ ] Number key assigned to switch to layout
- [ ] Focus switching works for all panels
- [ ] Help text updated with new shortcuts
- [ ] No conflicts with existing shortcuts

### ✅ Mouse Support
- [ ] Click handler exists for new layout
- [ ] Panel boundaries calculated correctly
- [ ] Uses correct coordinate system (X/Y)
- [ ] Boundary checks use >= and <
- [ ] Content offset accounts for title bar

### ✅ Responsive Design
- [ ] Layout works at 120x30 (desktop)
- [ ] Layout works at 80x24 (standard)
- [ ] Layout adapts at 60x20 (portrait)
- [ ] Layout handles <80 cols gracefully
- [ ] Optional: Switches to vertical stack when narrow

### ✅ Visual Consistency
- [ ] Uses existing panel styles
- [ ] Focused panel visually distinct
- [ ] Borders align properly
- [ ] No gaps or overlaps

## Example Layouts

### Triple Vertical Split (25% | 50% | 25%)
```
╔════════╦══════════════════╦════════╗
║ Left   ║ Middle (Focused) ║ Right  ║
║        ║                  ║        ║
║ Items  ║ Main Content     ║ Details║
╚════════╩══════════════════╩════════╝
```

### 2x2 Grid
```
╔═══════════════╦═══════════════╗
║ Top Left      ║ Top Right     ║
║               ║               ║
╠═══════════════╬═══════════════╣
║ Bottom Left   ║ Bottom Right  ║
║               ║               ║
╚═══════════════╩═══════════════╝
```

### Dashboard (Top Bar + 3 Bottom Panels)
```
╔═══════════════════════════════╗
║ Header / Summary              ║
╠═════════╦═════════╦═══════════╣
║ Panel 1 ║ Panel 2 ║ Panel 3   ║
║         ║         ║           ║
╚═════════╩═════════╩═══════════╝
```

## Files Modified

Typical files that will be changed:

- `types.go` or `model.go` - Layout constant and calculation (~20-40 lines)
- `view.go` - Rendering function (~40-80 lines)
- `update_keyboard.go` - Keyboard shortcuts (~10-20 lines)
- `update_mouse.go` - Mouse detection (~20-40 lines)
- Optional: `styles.go` - Layout-specific styles

## Best Practices Applied

- **Border Accounting:** Always subtract 2 from height for borders
- **Text Truncation:** Truncate ALL strings before rendering
- **Exact Dimensions:** Content fills exact height, no more, no less
- **Consistent Patterns:** Follow existing layout structure
- **Mouse Alignment:** Detection math matches rendering math
- **Keyboard First:** Direct number keys for layout switching
- **Visual Feedback:** Focused panel clearly distinguished

## Testing Checklist

After creating the layout, test:

1. **Switch to layout:** Press assigned number key
2. **Focus switching:** Use h/j/k/l or arrow keys
3. **Mouse clicks:** Click each panel, verify focus changes
4. **Resize terminal:** Test at multiple sizes
5. **Portrait mode:** Test at <80 cols width
6. **Mobile mode:** Test at 70x10 (Termux simulation)
7. **Border alignment:** All borders line up perfectly
8. **Text wrapping:** No text wraps at any width

## Reference

- **CLAUDE.md** - Critical layout fixes and patterns
- **LAZYGIT_ANALYSIS.md** - Layout calculation patterns
- **MOUSE_SUPPORT_GUIDE.md** - Click detection
- **Existing layouts** - Follow the same pattern

## Notes

- All layouts follow the same plumbing pattern
- The template enforces best practices automatically
- Layout calculation is the most critical part
- Mouse detection must exactly match rendered boundaries
- Test on portrait/vertical monitors to catch issues early
