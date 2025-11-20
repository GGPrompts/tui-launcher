# 3-Pane Layout Implementation Plan

## Overview
Transform the current single-pane tree view into a 3-pane adaptive layout:
```
┌─────────────────┬─────────────────┐
│  Global Tools   │    Projects     │
│  (Left Pane)    │  (Right Pane)   │
├─────────────────┴─────────────────┤
│        Info/Help (Bottom)         │
└───────────────────────────────────┘
```

## Step 1: Update Model Structure

### Add to types.go:
```go
type paneType int
const (
    paneGlobal paneType = iota  // Left pane
    paneProject                  // Right pane
)

type paneInfo struct {
    description string
    cliFlags    string
    repo        string
    mdPath      string  // Path to .md file for detailed info
}
```

### Update model struct:
```go
type model struct {
    // Existing fields...
    
    // New pane management
    activePane      paneType        // Which pane has focus
    globalItems     []launchItem    // Items for left pane
    projectItems    []launchItem    // Items for right pane
    globalCursor    int             // Cursor position in global pane
    projectCursor   int             // Cursor position in project pane
    
    // Info pane
    currentInfo     paneInfo        // Info for selected item
    infoContent     string          // Rendered markdown content
    
    // Layout
    leftPaneWidth   int             // Width of left pane
    rightPaneWidth  int             // Width of right pane
    infoPaneHeight  int             // Height of info pane
}
```

## Step 2: Refactor Tree Building

### In tree.go, split the tree building:
```go
func buildTreeFromConfig(config Config) ([]launchItem, []launchItem) {
    var globalItems []launchItem
    var projectItems []launchItem
    
    // Projects go to right pane
    for _, proj := range config.Projects {
        // Build project tree...
        projectItems = append(projectItems, projectItem)
    }
    
    // Tools, AI, Scripts go to left pane
    for _, cat := range config.Tools {
        // Build tool categories...
        globalItems = append(globalItems, categoryItem)
    }
    
    // AI commands to left pane
    // Scripts to left pane
    
    return globalItems, projectItems
}
```

## Step 3: Update View Rendering

### Create new view functions in view.go:
```go
func (m model) viewLeftPane() string {
    // Render global tools tree
    // Use m.globalCursor for highlighting
    // Show selection checkboxes
}

func (m model) viewRightPane() string {
    // Render projects tree  
    // Use m.projectCursor for highlighting
    // Show selection checkboxes
}

func (m model) viewInfoPane() string {
    // Render markdown content or info
    // Could use glamour for markdown rendering
    // Or simple text with formatting
}

func (m model) View() string {
    // Calculate pane dimensions based on m.width and m.height
    leftWidth := m.width / 2
    rightWidth := m.width - leftWidth
    treeHeight := m.height * 2 / 3  // Top 2/3 for trees
    infoHeight := m.height - treeHeight - 3  // Bottom 1/3 for info
    
    // Build each pane
    leftPane := m.viewLeftPane()
    rightPane := m.viewRightPane()
    infoPane := m.viewInfoPane()
    
    // Use lipgloss to create borders and layout
    // Join panes with box drawing characters
}
```

## Step 4: Update Navigation

### In update_keyboard.go:
```go
case "tab":
    // Switch between panes
    if m.activePane == paneGlobal {
        m.activePane = paneProject
    } else {
        m.activePane = paneGlobal
    }

case "up", "k":
    // Move cursor in active pane
    if m.activePane == paneGlobal {
        if m.globalCursor > 0 {
            m.globalCursor--
        }
    } else {
        if m.projectCursor > 0 {
            m.projectCursor--
        }
    }
    // Update info pane content
    m.updateInfoPane()

case "enter":
    // Launch from active pane
    var currentItem launchItem
    if m.activePane == paneGlobal {
        currentItem = m.globalTreeItems[m.globalCursor].item
    } else {
        currentItem = m.projectTreeItems[m.projectCursor].item
    }
    // Launch logic...
```

## Step 5: Add Info Content Loading

### Create info.go:
```go
func (m *model) updateInfoPane() {
    // Get current item based on active pane
    var currentItem launchItem
    if m.activePane == paneGlobal {
        if m.globalCursor < len(m.globalTreeItems) {
            currentItem = m.globalTreeItems[m.globalCursor].item
        }
    } else {
        if m.projectCursor < len(m.projectTreeItems) {
            currentItem = m.projectTreeItems[m.projectCursor].item
        }
    }
    
    // Load markdown file if specified
    if currentItem.InfoPath != "" {
        content, err := os.ReadFile(currentItem.InfoPath)
        if err == nil {
            m.infoContent = string(content)
            return
        }
    }
    
    // Otherwise show basic info
    m.infoContent = fmt.Sprintf(
        "Name: %s\nCommand: %s\nDirectory: %s\n",
        currentItem.Name,
        currentItem.Command,
        currentItem.Cwd,
    )
}
```

## Step 6: Update Config Structure

### Add to config.yaml:
```yaml
tools:
  - category: Git
    icon: 🔧
    items:
      - name: lazygit
        icon: 🎯
        command: lazygit
        description: "Terminal UI for git"
        info_file: ~/.config/tui-launcher/docs/lazygit.md
        repo: "https://github.com/jesseduffield/lazygit"

projects:
  - name: TUI Launcher
    icon: 🚀
    path: ~/projects/tui-launcher
    description: "Visual terminal launcher"
    info_file: ~/.config/tui-launcher/docs/tui-launcher.md
```

## Step 7: Add Responsive Layout (Critical for Termux!)

### Layout Modes

Add to types.go:
```go
type layoutMode int
const (
    layoutDesktop layoutMode = iota  // 3-pane (Global | Projects | Info)
    layoutCompact                    // 2-pane (Combined tree + info)
    layoutMobile                     // 1-pane (tree only, 'i' to toggle info)
)
```

### Responsive Detection

In model.go:
```go
func (m model) getLayoutMode() layoutMode {
    // Termux with keyboard open (very small)
    if m.height <= 12 {
        return layoutMobile
    }

    // Narrow terminal (portrait or phone)
    if m.width < 80 {
        return layoutCompact
    }

    // Desktop (normal terminal)
    return layoutDesktop
}

func (m model) calculateLayout() (int, int, int, int) {
    contentHeight := m.height - 3  // Header
    contentHeight -= 2              // Borders (Golden Rule #1!)

    mode := m.getLayoutMode()

    switch mode {
    case layoutDesktop:
        // 3-pane: Left | Right | Bottom
        leftWidth := m.width / 2
        rightWidth := m.width - leftWidth
        treeHeight := contentHeight * 2 / 3
        infoHeight := contentHeight - treeHeight
        return leftWidth, rightWidth, treeHeight, infoHeight

    case layoutCompact:
        // 2-pane: Combined tree on top | Info on bottom
        treeHeight := contentHeight * 3 / 4
        infoHeight := contentHeight - treeHeight
        return m.width, 0, treeHeight, infoHeight

    case layoutMobile:
        // 1-pane: Just tree, toggle info with 'i' key
        return m.width, 0, contentHeight, 0
    }
}
```

### Layout Examples

**Desktop (120x30):**
```
┌─────────────────────────────────────────┐
│ 🚀 TUI Launcher                         │
├──────────────────┬──────────────────────┤
│ Global Tools     │ Projects             │
│ ├─ Git           │ ├─ TUI Launcher      │
│ │  └─ lazygit    │ │  └─ Edit Config    │
│ └─ AI            │ └─ TKan              │
│    └─ claude     │    └─ Run TKan       │
├──────────────────┴──────────────────────┤
│ Info: lazygit                           │
│ Terminal UI for git commands            │
│ Repo: github.com/jesseduffield/lazygit  │
└─────────────────────────────────────────┘
```

**Compact (70x20) - Termux landscape:**
```
┌────────────────────────────────────┐
│ 🚀 TUI Launcher                    │
├────────────────────────────────────┤
│ Global Tools / Projects            │
│ [Tab] to switch                    │
│ ├─ Git                             │
│ │  └─ lazygit                      │
├────────────────────────────────────┤
│ Info: lazygit - TUI for git        │
└────────────────────────────────────┘
```

**Mobile (70x10) - Termux with keyboard:**
```
┌────────────────────────────────────┐
│ 🚀 Launcher [i]=info               │
├────────────────────────────────────┤
│ ├─ Git                             │
│ │  └─ lazygit                      │
│ └─ AI                              │
│    └─ claude                       │
└────────────────────────────────────┘
```

### Keyboard Shortcuts

Update for responsive modes:
```go
case "tab":
    // Switch panes in desktop mode
    // Switch between Global/Projects in compact mode
    if m.getLayoutMode() == layoutDesktop {
        // Toggle left/right pane
        if m.activePane == paneGlobal {
            m.activePane = paneProject
        } else {
            m.activePane = paneGlobal
        }
    } else if m.getLayoutMode() == layoutCompact {
        // Toggle between global tools and projects
        m.showingProjects = !m.showingProjects
    }

case "i":
    // Toggle info pane in mobile mode
    if m.getLayoutMode() == layoutMobile {
        m.showingInfo = !m.showingInfo
    }
```

### Golden Rules Applied

From the Bubbletea skill:
1. **Always account for borders** - Subtract 2 from height BEFORE rendering
2. **Never auto-wrap** - Truncate all text to prevent wrapping
3. **Match mouse to layout** - Use X coords for horizontal, Y for vertical
4. **Use proportional sizing** - No hardcoded pixel values

## Implementation Order

1. **Start with types.go** - Add new types and constants
2. **Update model.go** - Add pane fields to model struct
3. **Modify tree.go** - Split items into global/project
4. **Create basic 3-pane view** - Get layout working with placeholders
5. **Add pane navigation** - Tab to switch, arrows to navigate
6. **Implement info pane** - Load and display markdown/info
7. **Polish and test** - Borders, colors, responsive behavior

## Tips for Implementation

- Use lipgloss borders and joins for clean pane divisions
- Keep existing tree rendering logic, just call it twice (once per pane)
- Start simple - get the layout working before adding markdown rendering
- Test with different terminal sizes early
- Consider using viewport from bubbles for scrolling in info pane

This leverages your existing Bubble Tea knowledge and builds on the tree view code you already have!