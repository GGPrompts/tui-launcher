# TUI Fix Responsive Layout Skill

Apply all critical responsive layout fixes to prevent panels from overflowing, misaligning, or breaking on different terminal sizes.

## Usage

When you notice any of these issues:
- Panels covering the title bar
- Panels misaligned (one row different)
- Mouse clicks not working after resize
- Layout breaks on small/vertical terminals

Run:
```
/tui-fix-responsive
```

## What I'll Do

I will systematically apply all three critical fixes documented in CLAUDE.md:

### Fix 1: Border Height Accounting (CLAUDE.md:14-54)

**Problem:** Panels overflow and cover the title bar because border height isn't subtracted.

**The Math:**
```
WRONG:
contentHeight = totalHeight - 3 (title) - 1 (status) = totalHeight - 4
Panel renders with borders = contentHeight + 2 (borders)
Actual height used = totalHeight - 4 + 2 = totalHeight - 2 (TOO TALL!)

CORRECT:
contentHeight = totalHeight - 3 (title) - 1 (status) - 2 (borders) = totalHeight - 6
Panel renders with borders = contentHeight + 2
Actual height used = totalHeight - 6 + 2 = totalHeight - 4 ✓
```

**Fix Applied in `model.go`:**
```go
func (m model) calculateLayout() (int, int) {
    contentWidth := m.width
    contentHeight := m.height

    if m.config.UI.ShowTitle {
        contentHeight -= 3 // title bar (3 lines)
    }
    if m.config.UI.ShowStatus {
        contentHeight -= 1 // status bar
    }

    // CRITICAL: Account for panel borders
    contentHeight -= 2 // top + bottom borders

    return contentWidth, contentHeight
}
```

### Fix 2: Text Truncation to Prevent Wrapping (CLAUDE.md:56-98)

**Problem:** When panels resize, text wraps to multiple lines, making panels taller and causing misalignment.

**Example:** `"Weight: 2 | Size: 80x25"` wraps to 2 lines when panel narrows, making it 1 row taller.

**Fix Applied in `view.go`:**
```go
// Calculate max text width to prevent wrapping
maxTextWidth := width - 4 // -2 for borders, -2 for padding

// Truncate ALL text before rendering
title = truncateString(title, maxTextWidth)
subtitle = truncateString(subtitle, maxTextWidth)

// Truncate content lines too
for i := 0; i < availableContentLines && i < len(content); i++ {
    line := truncateString(content[i], maxTextWidth)
    lines = append(lines, line)
}

// Helper function
func truncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-1] + "…"
}
```

**Key Insight:** NEVER let Lipgloss auto-wrap text in bordered panels. Always truncate explicitly.

### Fix 3: Mouse Detection for Vertical/Horizontal Modes (CLAUDE.md:100-136)

**Problem:** When terminal is narrow (<80 cols), panels stack vertically, but mouse clicks still use X coordinates instead of Y.

**Fix Applied in `update_mouse.go`:**
```go
func (m model) handleLeftClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    // ... boundary checks ...

    if m.shouldUseVerticalStack() {
        // Vertical stack mode: use Y coordinates
        topHeight, _ := m.calculateVerticalStackLayout()
        relY := msg.Y - contentStartY

        if relY < topHeight {
            m.focusedPanel = "left"  // Top panel
        } else if relY > topHeight {
            m.focusedPanel = "right" // Bottom panel
        }
    } else {
        // Side-by-side mode: use X coordinates
        leftWidth, _ := m.calculateDualPaneLayout()

        if msg.X < leftWidth {
            m.focusedPanel = "left"
        } else if msg.X > leftWidth {
            m.focusedPanel = "right"
        }
    }

    return m, nil
}
```

**Key Insight:** Mouse detection logic must match the layout orientation (horizontal vs vertical).

### Additional Responsive Fixes

I will also check and apply these related fixes:

#### Vertical Stack Detection
```go
func (m model) shouldUseVerticalStack() bool {
    return m.width < 80  // Switch to vertical when narrow
}
```

#### Mobile/Termux Optimizations (TERMUX_MOBILE_GUIDE.md)

For very small terminals (Termux with keyboard):
```go
func (m model) isMobileCompact() bool {
    return m.height <= 10  // Termux with keyboard open
}

// In calculateLayout():
if m.isMobileCompact() {
    // Remove status bar to save space
    contentHeight -= 0  // No status bar
} else if m.config.UI.ShowStatus {
    contentHeight -= 1  // Normal status bar
}
```

## Critical Checks

After applying fixes, I will verify:

### ✅ Border Height Accounting
- [ ] `calculateLayout()` subtracts 2 for borders: `contentHeight -= 2`
- [ ] No explicit `Height()` calls on bordered Lipgloss styles
- [ ] Panel rendering fills content exactly (padding with empty lines if needed)
- [ ] Total layout math: `titleLines + contentHeight + borderLines + statusLines = totalHeight`

**Test:** On a portrait monitor or small terminal, title bar should NEVER be covered by panels.

### ✅ Text Truncation
- [ ] `truncateString()` helper function exists
- [ ] All titles truncated: `truncateString(title, maxTextWidth)`
- [ ] All subtitles truncated: `truncateString(subtitle, maxTextWidth)`
- [ ] All content lines truncated in loops
- [ ] Max width calculated: `width - 4` (borders + padding)

**Test:** Resize terminal to very narrow width. No text should wrap to multiple lines.

### ✅ Mouse Detection Modes
- [ ] `shouldUseVerticalStack()` function exists and checks width
- [ ] Mouse handler checks layout mode: `if m.shouldUseVerticalStack()`
- [ ] Vertical mode uses Y coordinates for panel detection
- [ ] Horizontal mode uses X coordinates for panel detection
- [ ] Relative coordinates calculated: `relY = msg.Y - contentStartY`

**Test:** Resize to <80 cols. Click on top panel, verify it focuses. Click on bottom panel, verify it focuses.

### ✅ Responsive Breakpoints
- [ ] Vertical stack threshold: width < 80
- [ ] Mobile compact threshold: height <= 10
- [ ] Minimum terminal size check exists
- [ ] Layout adapts smoothly at breakpoints

**Test:** Gradually resize terminal from 120x30 down to 60x10. Layout should adapt without breaking.

### ✅ Border Consistency
- [ ] All panels use same border style
- [ ] No mixing of `Height()` with natural height
- [ ] Content fills exactly to prevent gaps

**Test:** All panel borders should align perfectly, no gaps or overlaps.

## Files Modified

This skill will modify:

- `model.go` (or `types.go`) - Layout calculation functions (lines ~80-98)
- `view.go` - Panel rendering with truncation (lines ~340-352, ~143-180)
- `update_mouse.go` - Mouse click detection (lines ~79-106)
- May add helper functions like `truncateString()`, `shouldUseVerticalStack()`

## Before & After

### Before (Broken Layout)
```
Terminal: 60x25 (portrait)

[Panels covering title - CAN'T SEE APP NAME]
╔═══════════════════════════════╗
║ Panel 1                       ║
║ Weight: 2 | Size:             ║ ← Wrapped text!
║ 80x25 (causing extra height)  ║
╚═══════════════════════════════╝
╔═══════════════════════════════╗
║ Panel 2 (one row lower!)      ║ ← Misaligned!
╚═══════════════════════════════╝
[Status bar pushed off screen]
```

### After (Fixed Layout)
```
Terminal: 60x25 (portrait)

📱 My App | 60×25                    ← Title visible!
╔═══════════════════════════════╗
║ Panel 1                       ║
║ Weight: 2 | Size: 80x2…       ║ ← Truncated!
║                               ║
╚═══════════════════════════════╝
╔═══════════════════════════════╗
║ Panel 2                       ║ ← Aligned!
╚═══════════════════════════════╝
Status: Ready | Mouse: 30,12     ← Status visible!
```

## Debugging Workflow

If issues persist after applying fixes, I will:

1. **Log the math:**
   ```go
   fmt.Printf("Total: %d, Title: %d, Status: %d, Borders: %d, Content: %d\n",
       m.height, titleHeight, statusHeight, 2, contentHeight)
   ```

2. **Verify border rendering:**
   - Check that panels don't set `.Height()` explicitly
   - Confirm content is padded to exact height with empty lines

3. **Test mouse coordinates:**
   - Show mouse position in status bar: `fmt.Sprintf("Mouse: %d,%d", x, y)`
   - Log click events to verify detection logic

4. **Test at critical sizes:**
   - 120x30 (desktop)
   - 80x24 (standard)
   - 60x20 (portrait)
   - 70x10 (Termux with keyboard)

## Best Practices Applied

- **Always subtract borders from height calculations** (CLAUDE.md critical fix #1)
- **Always truncate text before rendering** (CLAUDE.md critical fix #2)
- **Always match mouse logic to layout orientation** (CLAUDE.md critical fix #3)
- **Use relative coordinates for click detection** (Y relative to content start)
- **Provide visual feedback** (status bar shows mouse position)

## Reference Documentation

- **CLAUDE.md** - All three critical fixes documented here
- **TERMUX_MOBILE_GUIDE.md** - Mobile responsive patterns
- **MOUSE_SUPPORT_GUIDE.md** - Mouse detection best practices

## Notes

- These fixes prevent 90% of layout issues in TUI apps
- The fixes are cumulative - all three must be applied
- Always test on portrait/vertical monitors
- Termux testing reveals edge cases desktop terminals hide
- Math is critical: `totalHeight = title + content + borders + status`
