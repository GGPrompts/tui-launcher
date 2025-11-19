# Tmux Layouts Demo

## How It Works

When you select multiple items and press Enter, the launcher shows layout options based on the count.

---

## Example: 2 Items Selected

```
☑ 📂 TFE
☑ 💻 go run .

┌─ Choose Layout ──────────────────────────────┐
│                                               │
│  ● Side-by-Side (even-horizontal)            │
│    Equal width columns                        │
│    ┌─────┬─────┐                             │
│    │  1  │  2  │                             │
│    └─────┴─────┘                             │
│                                               │
│  ○ Top-Bottom (even-vertical)                │
│    Equal height rows                          │
│    ┌───────────┐                             │
│    │     1     │                             │
│    ├───────────┤                             │
│    │     2     │                             │
│    └───────────┘                             │
│                                               │
│  ○ Main Left (main-vertical)                 │
│    Large left, small right                    │
│    ┌────────┬──┐                             │
│    │        │2 │                             │
│    │   1    │  │                             │
│    └────────┴──┘                             │
│                                               │
│           [Launch]  [Cancel]                  │
└───────────────────────────────────────────────┘
```

---

## Example: 4 Items Selected (Quad Split!)

```
☑ 📂 TFE
☑ 💻 go run .
☑ 📊 tail -f debug.log
☑ 💹 htop

┌─ Choose Layout ──────────────────────────────┐
│                                               │
│  ● Quad Split (tiled)                        │
│    2x2 grid                                   │
│    ┌─────┬─────┐                             │
│    │  1  │  2  │                             │
│    ├─────┼─────┤                             │
│    │  3  │  4  │                             │
│    └─────┴─────┘                             │
│                                               │
│  ○ Main + Stack (main-vertical)              │
│    Large left, 3 stacked right                │
│    ┌────────┬──┐                             │
│    │        │2 │                             │
│    │   1    ├──┤                             │
│    │        │3 │                             │
│    │        ├──┤                             │
│    │        │4 │                             │
│    └────────┴──┘                             │
│                                               │
│  ○ 4 Columns (even-horizontal)               │
│    Equal width columns                        │
│    ┌──┬──┬──┬──┐                             │
│    │1 │2 │3 │4 │                             │
│    └──┴──┴──┴──┘                             │
│                                               │
│           [Launch]  [Cancel]                  │
└───────────────────────────────────────────────┘
```

---

## Example: 6 Items Selected

```
☑ 📂 TFE
☑ 💻 go run .
☑ 📊 tail -f debug.log
☑ 💹 htop
☑ 🦥 lazygit
☑ 🧪 go test ./...

┌─ Choose Layout ──────────────────────────────┐
│                                               │
│  ● Tiled Grid (tiled)                        │
│    3x2 grid                                   │
│    ┌────┬────┬────┐                          │
│    │ 1  │ 2  │ 3  │                          │
│    ├────┼────┼────┤                          │
│    │ 4  │ 5  │ 6  │                          │
│    └────┴────┴────┘                          │
│                                               │
│  ○ 6 Columns (even-horizontal)               │
│    Equal width columns                        │
│    ┌─┬─┬─┬─┬─┬─┐                             │
│    │1│2│3│4│5│6│                             │
│    └─┴─┴─┴─┴─┴─┘                             │
│                                               │
│           [Launch]  [Cancel]                  │
└───────────────────────────────────────────────┘
```

---

## How Spawning Works

### When You Press [Launch]:

```bash
# Create tmux session (if not in tmux)
tmux new-session -d -s "launcher-12345" -c ~/projects/tfe "tfe"

# Add remaining panes
tmux split-window -t launcher-12345 -c ~/projects/tfe "go run ."
tmux split-window -t launcher-12345 -c ~/projects/tfe "tail -f debug.log"
tmux split-window -t launcher-12345 -c ~/projects/tfe "htop"

# Apply selected layout
tmux select-layout -t launcher-12345 tiled

# Attach or switch to session
tmux attach -t launcher-12345  # If not in tmux
# OR
tmux switch-client -t launcher-12345  # If already in tmux
```

### Result:

```
┌─────────────┬─────────────┐
│    TFE      │  go run .   │
├─────────────┼─────────────┤
│ tail -f log │    htop     │
└─────────────┴─────────────┘
```

---

## Saved Profiles

You can also save common layouts as profiles:

```yaml
# config.yaml
projects:
  - name: TFE
    profiles:
      - name: Dev Environment
        icon: 🔧
        layout: main-vertical
        panes:
          - command: tfe
          - command: go run .
          - command: tail -f logs/debug.log
```

Press Enter on "Dev Environment" → Instant 3-pane layout!

---

## Advanced: Custom Layouts

Tmux also supports custom layout strings (for very specific arrangements):

```bash
# Custom layout string (width,height positions)
tmux select-layout "2e0e,211x54,0,0{105x54,0,0,0,105x54,106,0[105x26,106,0,1,105x27,106,27,2]}"
```

We could add a "Custom Layout" option where you paste a layout string from a working tmux session:

```bash
# In tmux, get current layout
tmux list-windows -F "#{window_layout}"
```

Then save it to a profile!

---

## Summary

- **Dynamic layouts** adapt to item count
- **Visual previews** show what you'll get
- **Arrow keys** to select layout
- **Enter** to spawn
- **Automatic** - tmux handles the math!
