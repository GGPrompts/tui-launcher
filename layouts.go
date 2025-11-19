package main

import "fmt"

// suggestLayouts returns appropriate layout options based on pane count
// Tmux layouts automatically adapt to the number of panes
func suggestLayouts(count int) []layoutOption {
	switch count {
	case 1:
		// Single item - no layout needed
		return []layoutOption{
			{
				layout:      layoutMainVertical,
				name:        "Single Pane",
				description: "Launch in current/new pane",
				preview:     "┌────────────┐\n│            │\n│     1      │\n│            │\n└────────────┘",
			},
		}

	case 2:
		return []layoutOption{
			{
				layout:      layoutEvenHorizontal,
				name:        "Side-by-Side",
				description: "Equal width columns",
				preview:     "┌─────┬─────┐\n│  1  │  2  │\n│     │     │\n└─────┴─────┘",
			},
			{
				layout:      layoutEvenVertical,
				name:        "Top-Bottom",
				description: "Equal height rows",
				preview:     "┌───────────┐\n│     1     │\n├───────────┤\n│     2     │\n└───────────┘",
			},
			{
				layout:      layoutMainVertical,
				name:        "Main Left",
				description: "Large left, small right",
				preview:     "┌────────┬──┐\n│        │2 │\n│   1    │  │\n│        │  │\n└────────┴──┘",
			},
		}

	case 3:
		return []layoutOption{
			{
				layout:      layoutMainVertical,
				name:        "Main + Stack",
				description: "Large left, 2 stacked right",
				preview:     "┌────────┬──┐\n│        │2 │\n│   1    ├──┤\n│        │3 │\n└────────┴──┘",
			},
			{
				layout:      layoutTiled,
				name:        "Tiled",
				description: "Automatic grid (2+1)",
				preview:     "┌─────┬─────┐\n│  1  │  2  │\n├─────┴─────┤\n│     3     │\n└───────────┘",
			},
			{
				layout:      layoutEvenHorizontal,
				name:        "3 Columns",
				description: "Equal width columns",
				preview:     "┌───┬───┬───┐\n│ 1 │ 2 │ 3 │\n│   │   │   │\n└───┴───┴───┘",
			},
		}

	case 4:
		return []layoutOption{
			{
				layout:      layoutTiled,
				name:        "Quad Split",
				description: "2x2 grid",
				preview:     "┌─────┬─────┐\n│  1  │  2  │\n├─────┼─────┤\n│  3  │  4  │\n└─────┴─────┘",
			},
			{
				layout:      layoutMainVertical,
				name:        "Main + Stack",
				description: "Large left, 3 stacked right",
				preview:     "┌────────┬──┐\n│        │2 │\n│   1    ├──┤\n│        │3 │\n│        ├──┤\n│        │4 │\n└────────┴──┘",
			},
			{
				layout:      layoutEvenHorizontal,
				name:        "4 Columns",
				description: "Equal width columns",
				preview:     "┌──┬──┬──┬──┐\n│1 │2 │3 │4 │\n│  │  │  │  │\n└──┴──┴──┴──┘",
			},
		}

	case 5:
		return []layoutOption{
			{
				layout:      layoutTiled,
				name:        "Tiled Grid",
				description: "Automatic grid (3+2)",
				preview:     "┌────┬────┬────┐\n│ 1  │ 2  │ 3  │\n├────┴────┼────┤\n│    4    │ 5  │\n└─────────┴────┘",
			},
			{
				layout:      layoutMainVertical,
				name:        "Main + Stack",
				description: "Large left, 4 stacked right",
				preview:     "┌────────┬──┐\n│        │2 │\n│   1    ├──┤\n│        │3 │\n│        ├──┤\n│        │4 │\n│        ├──┤\n│        │5 │\n└────────┴──┘",
			},
		}

	case 6:
		return []layoutOption{
			{
				layout:      layoutTiled,
				name:        "Tiled Grid",
				description: "3x2 grid",
				preview:     "┌────┬────┬────┐\n│ 1  │ 2  │ 3  │\n├────┼────┼────┤\n│ 4  │ 5  │ 6  │\n└────┴────┴────┘",
			},
			{
				layout:      layoutEvenHorizontal,
				name:        "6 Columns",
				description: "Equal width columns",
				preview:     "┌─┬─┬─┬─┬─┬─┐\n│1│2│3│4│5│6│\n└─┴─┴─┴─┴─┴─┘",
			},
		}

	default:
		// 7+ panes
		return []layoutOption{
			{
				layout:      layoutTiled,
				name:        "Tiled Grid",
				description: fmt.Sprintf("Auto-arrange %d panes", count),
				preview:     "┌──┬──┬──┬──┐\n│  │  │  │  │\n├──┼──┼──┼──┤\n│  │  │  │  │\n└──┴──┴──┴──┘",
			},
			{
				layout:      layoutMainVertical,
				name:        "Main + Stack",
				description: fmt.Sprintf("1 large, %d stacked", count-1),
				preview:     "┌──────┬─┐\n│      │▪│\n│  1   │▪│\n│      │▪│\n└──────┴─┘",
			},
			{
				layout:      layoutEvenHorizontal,
				name:        "Columns",
				description: fmt.Sprintf("%d equal columns", count),
				preview:     "┌┬┬┬┬┬┬┬┐\n││││││││\n└┴┴┴┴┴┴┴┘",
			},
		}
	}
}

// getLayoutPreview returns ASCII art preview for a layout with given pane count
func getLayoutPreview(layout tmuxLayout, count int) string {
	switch layout {
	case layoutEvenHorizontal:
		switch count {
		case 2:
			return "┌─────┬─────┐\n│  1  │  2  │\n└─────┴─────┘"
		case 3:
			return "┌───┬───┬───┐\n│ 1 │ 2 │ 3 │\n└───┴───┴───┘"
		case 4:
			return "┌──┬──┬──┬──┐\n│1 │2 │3 │4 │\n└──┴──┴──┴──┘"
		default:
			return "┌─┬─┬─┬─┬─┐\n│ │ │ │ │ │\n└─┴─┴─┴─┴─┘"
		}

	case layoutEvenVertical:
		switch count {
		case 2:
			return "┌───────┐\n│   1   │\n├───────┤\n│   2   │\n└───────┘"
		case 3:
			return "┌───────┐\n│   1   │\n├───────┤\n│   2   │\n├───────┤\n│   3   │\n└───────┘"
		default:
			return "┌───────┐\n│   ▪   │\n├───────┤\n│   ▪   │\n├───────┤\n│   ▪   │\n└───────┘"
		}

	case layoutMainVertical:
		return "┌────────┬──┐\n│        │▪ │\n│   1    ├──┤\n│        │▪ │\n└────────┴──┘"

	case layoutMainHorizontal:
		return "┌──────────┐\n│    1     │\n├────┬─────┤\n│ ▪  │  ▪  │\n└────┴─────┘"

	case layoutTiled:
		switch count {
		case 2:
			return "┌─────┬─────┐\n│  1  │  2  │\n└─────┴─────┘"
		case 3:
			return "┌─────┬─────┐\n│  1  │  2  │\n├─────┴─────┤\n│     3     │\n└───────────┘"
		case 4:
			return "┌─────┬─────┐\n│  1  │  2  │\n├─────┼─────┤\n│  3  │  4  │\n└─────┴─────┘"
		case 6:
			return "┌───┬───┬───┐\n│ 1 │ 2 │ 3 │\n├───┼───┼───┤\n│ 4 │ 5 │ 6 │\n└───┴───┴───┘"
		default:
			return "┌──┬──┬──┐\n│  │  │  │\n├──┼──┼──┤\n│  │  │  │\n└──┴──┴──┘"
		}

	default:
		return "┌───────┐\n│       │\n│   ?   │\n└───────┘"
	}
}
