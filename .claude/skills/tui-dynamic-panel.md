# TUI Dynamic Panel Skill

Create LazyGit-style resizable panels that expand when focused using weight-based layouts.

## Usage

When you want panels to dynamically resize based on focus (accordion mode):

```
/tui-dynamic-panel
```

This will implement the weight-based layout system inspired by LazyGit.

## What I'll Do

I will implement dynamic panel resizing using a weight-based system, where focused panels get more space.

### 1. Add Accordion Mode Config (`types.go` or `config.go`)

Add configuration for accordion/expand-on-focus mode:

```go
type UIConfig struct {
    ShowTitle         bool
    ShowStatus        bool
    MouseEnabled      bool
    AccordionMode     bool  // NEW: Expand focused panel
    FocusedWeight     int   // NEW: Weight for focused panel (default: 2)
    UnfocusedWeight   int   // NEW: Weight for unfocused panel (default: 1)
}

// Default config
func defaultConfig() Config {
    return Config{
        UI: UIConfig{
            ShowTitle:       true,
            ShowStatus:      true,
            MouseEnabled:    true,
            AccordionMode:   true,   // Enable by default
            FocusedWeight:   2,      // Focused panel gets 2x space
            UnfocusedWeight: 1,      // Unfocused panels get 1x space
        },
    }
}
```

### 2. Update Model to Track Focus (`types.go` or `model.go`)

Ensure the model tracks which panel is focused:

```go
type model struct {
    width          int
    height         int
    currentLayout  LayoutMode
    focusedPanel   string       // "left", "right", "top", etc.
    config         Config
    // ... other fields
}
```

### 3. Implement Weight-Based Layout Calculation (`model.go`)

**Key Concept:** Use weights instead of fixed sizes!

```go
// calculateDualPaneLayout returns widths based on focus and accordion mode
func (m model) calculateDualPaneLayout() (leftWidth, rightWidth int) {
    totalWidth := m.width

    if !m.config.UI.AccordionMode {
        // Accordion disabled: equal split (50/50)
        leftWidth = totalWidth / 2
        rightWidth = totalWidth - leftWidth
        return leftWidth, rightWidth
    }

    // Accordion enabled: use weight-based calculation
    leftWeight := m.config.UI.UnfocusedWeight
    rightWeight := m.config.UI.UnfocusedWeight

    if m.focusedPanel == "left" {
        leftWeight = m.config.UI.FocusedWeight  // Focused gets more weight
    } else if m.focusedPanel == "right" {
        rightWeight = m.config.UI.FocusedWeight
    }

    // Calculate widths from weights
    totalWeight := leftWeight + rightWeight
    leftWidth = (totalWidth * leftWeight) / totalWeight
    rightWidth = totalWidth - leftWidth

    return leftWidth, rightWidth
}
```

**The Math:**
- Equal weights (1:1) → 50/50 split
- Focused weight (2:1) → 66/33 split
- Focused weight (3:1) → 75/25 split

**Why weights work:**
```
Example: 120 cols, left focused (weight 2), right unfocused (weight 1)
Total weight = 2 + 1 = 3
Left width = (120 * 2) / 3 = 80 cols (66%)
Right width = (120 * 1) / 3 = 40 cols (33%)
```

### 4. Multi-Panel Accordion (Vertical Stack)

For multiple panels stacked vertically:

```go
// calculateMultiPanelHeights returns heights based on accordion mode
func (m model) calculateMultiPanelHeights(panelNames []string, totalHeight int) []int {
    numPanels := len(panelNames)

    if !m.config.UI.AccordionMode {
        // Equal distribution
        baseHeight := totalHeight / numPanels
        heights := make([]int, numPanels)

        for i := 0; i < numPanels; i++ {
            heights[i] = baseHeight
        }

        // Give remainder to first panel
        heights[0] += totalHeight % numPanels

        return heights
    }

    // Accordion mode: focused panel gets most space
    weights := make([]int, numPanels)
    totalWeight := 0

    for i, name := range panelNames {
        if name == m.focusedPanel {
            weights[i] = m.config.UI.FocusedWeight  // Focused: 2x
        } else {
            weights[i] = m.config.UI.UnfocusedWeight  // Unfocused: 1x
        }
        totalWeight += weights[i]
    }

    // Calculate heights from weights
    heights := make([]int, numPanels)
    usedHeight := 0

    for i := 0; i < numPanels-1; i++ {
        heights[i] = (totalHeight * weights[i]) / totalWeight
        usedHeight += heights[i]
    }

    // Last panel gets remainder (prevents rounding errors)
    heights[numPanels-1] = totalHeight - usedHeight

    return heights
}
```

**Example:** 30 rows, 4 panels, "panel2" focused (weight 2), others weight 1
```
Total weight = 1 + 2 + 1 + 1 = 5
Panel 1: (30 * 1) / 5 = 6 rows
Panel 2: (30 * 2) / 5 = 12 rows ← FOCUSED
Panel 3: (30 * 1) / 5 = 6 rows
Panel 4: (30 * 1) / 5 = 6 rows
```

### 5. Update Rendering to Use Dynamic Widths (`view.go`)

Update rendering to use the weight-based calculations:

```go
func (m model) renderDualLayout() string {
    // Get dynamic widths based on focus
    leftWidth, rightWidth := m.calculateDualPaneLayout()

    // Height same for both
    _, contentHeight := m.calculateLayout()

    // Render panels with calculated widths
    leftPanel := m.renderLeftPanel(leftWidth, contentHeight)
    rightPanel := m.renderRightPanel(rightWidth, contentHeight)

    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        leftPanel,
        rightPanel,
    )
}
```

### 6. Focus Switching with Keyboard (`update_keyboard.go`)

Switching focus triggers instant resize:

```go
func (m model) handlePanelNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "h", "left":
        m.focusedPanel = "left"
        // Layout automatically recalculates on next render!
        return m, nil

    case "l", "right":
        m.focusedPanel = "right"
        return m, nil
    }

    return m, nil
}
```

**Key Insight:** No animations needed! Bubbletea rerenders, weights recalculate, panels resize smoothly.

### 7. Focus Switching with Mouse (`update_mouse.go`)

Clicking a panel focuses it and triggers resize:

```go
func (m model) handleLeftClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    // Get CURRENT widths (before focus change)
    leftWidth, _ := m.calculateDualPaneLayout()

    // Detect which panel was clicked
    if msg.X < leftWidth {
        m.focusedPanel = "left"
    } else {
        m.focusedPanel = "right"
    }

    // Next render will use NEW widths based on new focus
    return m, nil
}
```

**Important:** Use current widths for detection, new widths apply on next render.

### 8. Toggle Accordion Mode (`update_keyboard.go`)

Allow users to toggle accordion mode on/off:

```go
case "a":
    // Toggle accordion mode
    m.config.UI.AccordionMode = !m.config.UI.AccordionMode
    return m, nil
```

### 9. Visual Feedback (`view.go`)

Show accordion status in title or status bar:

```go
func (m model) renderTitleBar() string {
    mode := "Equal"
    if m.config.UI.AccordionMode {
        mode = "Accordion"
    }

    title := fmt.Sprintf("App | Layout: %s | Mode: %s",
        m.currentLayout, mode)

    return titleStyle.Render(title)
}
```

## Critical Checks

After implementing dynamic panels, verify:

### ✅ Weight Calculation
- [ ] Weight-based calculation function exists
- [ ] Handles equal weights (accordion off): all panels same size
- [ ] Handles focused weight: focused panel larger
- [ ] Total weight calculated: `sum(all weights)`
- [ ] Width/height formula: `(total * weight) / totalWeight`
- [ ] Last panel gets remainder to prevent rounding gaps

### ✅ Focus Management
- [ ] Model tracks `focusedPanel` string
- [ ] Keyboard shortcuts change focus
- [ ] Mouse clicks change focus
- [ ] Focus change triggers automatic recalculation

### ✅ Layout Recalculation
- [ ] Layout functions called on every render
- [ ] Weights updated based on current focus
- [ ] No caching of old dimensions
- [ ] Smooth resize (no flicker)

### ✅ Accordion Toggle
- [ ] Config has `AccordionMode` boolean
- [ ] Keyboard shortcut toggles mode (e.g., "a")
- [ ] When off: equal split
- [ ] When on: weight-based split
- [ ] Visual indicator shows current mode

### ✅ Mouse Detection Alignment
- [ ] Click detection uses CURRENT widths (before focus change)
- [ ] Boundaries recalculated on each click
- [ ] No stale width values
- [ ] Clicks work immediately after resize

### ✅ Multi-Panel Support
- [ ] Vertical stacking uses height weights
- [ ] Horizontal splitting uses width weights
- [ ] Supports 2+ panels
- [ ] All panels accounted for in weight calculations

### ✅ Edge Cases
- [ ] Handles single panel (no resizing)
- [ ] Handles very small terminals
- [ ] Prevents panels from becoming too small (optional minimum)
- [ ] Handles focus on non-existent panel (fallback to first panel)

### ✅ Visual Consistency
- [ ] Focused panel visually distinct (highlighted border)
- [ ] Resize is instant (no animation lag)
- [ ] No gaps between panels
- [ ] Borders align perfectly

## Before & After

### Before (Static Layout)
```
Terminal: 120 cols

╔═══════════════════════════════╗ ╔═══════════════════════════════╗
║ Left Panel                    ║ ║ Right Panel                   ║
║ (always 60 cols)              ║ ║ (always 60 cols)              ║
║                               ║ ║                               ║
╚═══════════════════════════════╝ ╚═══════════════════════════════╝
```

### After (Dynamic Layout - Left Focused)
```
Terminal: 120 cols

╔══════════════════════════════════════════╗ ╔═══════════════╗
║ Left Panel (FOCUSED)                     ║ ║ Right Panel   ║
║ (80 cols = 66%)                          ║ ║ (40 cols=33%) ║
║                                          ║ ║               ║
╚══════════════════════════════════════════╝ ╚═══════════════╝
```

### After (Dynamic Layout - Right Focused)
```
Terminal: 120 cols

╔═══════════════╗ ╔══════════════════════════════════════════╗
║ Left Panel    ║ ║ Right Panel (FOCUSED)                    ║
║ (40 cols=33%) ║ ║ (80 cols = 66%)                          ║
║               ║ ║                                          ║
╚═══════════════╝ ╚══════════════════════════════════════════╝
```

## Configuration Options

Users can customize weights:

```yaml
# config.yaml
ui:
  accordion_mode: true
  focused_weight: 3      # Focused panel gets 3x space
  unfocused_weight: 1    # Unfocused panels get 1x space
```

**Result:** Focused panel gets 75% (3/(3+1)), unfocused gets 25%.

## Advanced: Minimum Panel Sizes

Prevent panels from becoming unusably small:

```go
func (m model) calculateDualPaneLayout() (leftWidth, rightWidth int) {
    const minPanelWidth = 20  // Minimum 20 cols per panel

    totalWidth := m.width

    // Calculate with weights
    leftWeight, rightWeight := 1, 1
    if m.config.UI.AccordionMode && m.focusedPanel == "left" {
        leftWeight = m.config.UI.FocusedWeight
    } else if m.config.UI.AccordionMode && m.focusedPanel == "right" {
        rightWeight = m.config.UI.FocusedWeight
    }

    totalWeight := leftWeight + rightWeight
    leftWidth = (totalWidth * leftWeight) / totalWeight
    rightWidth = totalWidth - leftWidth

    // Enforce minimums
    if leftWidth < minPanelWidth {
        leftWidth = minPanelWidth
        rightWidth = totalWidth - leftWidth
    }
    if rightWidth < minPanelWidth {
        rightWidth = minPanelWidth
        leftWidth = totalWidth - rightWidth
    }

    return leftWidth, rightWidth
}
```

## Best Practices Applied

- **Proportional Sizing:** Use weights, not fixed pixels
- **Immediate Resize:** No animations, instant recalculation
- **Simple Math:** Total weight, then divide proportionally
- **Configurable:** Users can adjust weight ratios
- **Smooth UX:** Focus change automatically triggers resize
- **No Caching:** Always recalculate on render (cheap operation)

## LazyGit Reference

This pattern is inspired by LazyGit's implementation:

- **Box Layout System:** Weights determine proportions
- **Focus-Based Resizing:** Focused items get more space
- **Accordion Mode:** Configurable expand-on-focus
- **Smooth Transitions:** Instant, no animation needed

See `LAZYGIT_ANALYSIS.md` for detailed breakdown.

## Files Modified

- `types.go` or `model.go` - Add `focusedPanel`, config options (~10 lines)
- `config.go` - Add accordion config (~5 lines)
- `model.go` - Weight-based layout calculation functions (~30-50 lines)
- `view.go` - Use dynamic widths/heights (~10 lines)
- `update_keyboard.go` - Focus switching, accordion toggle (~10 lines)
- `update_mouse.go` - Update click detection (~5 lines)

## Testing Checklist

1. **Focus left panel:** Press `h`, verify left expands
2. **Focus right panel:** Press `l`, verify right expands
3. **Click panels:** Click left, verify it expands; click right, verify it expands
4. **Toggle accordion:** Press `a`, verify equal split; press `a` again, verify weighted split
5. **Resize terminal:** Drag to different sizes, verify proportions maintained
6. **Multi-panel:** Test with 3+ panels, verify weights work
7. **Edge cases:** Very small terminal, single panel

## Reference

- **LAZYGIT_ANALYSIS.md** - Complete breakdown of LazyGit's weight system
- **CLAUDE.md** - Layout calculation best practices
- **Examples:** Check LazyGit source for inspiration

## Notes

- Weight-based layouts are more flexible than fixed sizes
- The math is simple: `width = (total * weight) / totalWeight`
- No animations needed - Bubbletea handles smooth redraws
- Users love the "snap to focus" behavior
- Perfect for file browsers, code editors, dashboard apps
