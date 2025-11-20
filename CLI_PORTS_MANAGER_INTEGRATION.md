# CLI Ports Manager Integration

## Overview

[cli-ports-manager](https://github.com/adi-family/cli-ports-manager) is a Rust CLI tool that manages port assignments for services. It's perfect for organizing multi-service development environments.

## Why This Integration Makes Sense

**Problem:** When developing with multiple services (backend, frontend, database, cache), you need to:
- Track which ports each service uses
- Avoid port conflicts
- Remember port numbers across sessions
- Share port conventions with team

**Solution:** Use `ports-manager` to centralize port configuration, reference in launcher commands.

## Integration Patterns

### Pattern 1: Dynamic Port Assignment

Use `$(ports-manager get service)` in commands:

```yaml
projects:
  - name: Full Stack App
    icon: 🌐
    path: ~/projects/myapp
    commands:
      - name: Backend API
        icon: 🔧
        command: go run . -port $(ports-manager get backend)
        cwd: ~/projects/myapp/backend
        spawn: tmux-split-h

      - name: Frontend Dev Server
        icon: 💻
        command: npm run dev -- -p $(ports-manager get frontend)
        cwd: ~/projects/myapp/frontend
        spawn: tmux-split-v

      - name: PostgreSQL
        icon: 🐘
        command: docker run -p $(ports-manager get postgres):5432 postgres:15
        spawn: tmux-window

      - name: Redis Cache
        icon: 📦
        command: docker run -p $(ports-manager get redis):6379 redis:7
        spawn: tmux-window
```

**Benefits:**
- Centralized port management
- Consistent across team (share `~/.config/ports-manager/config.toml`)
- Auto-allocation if port not set
- Easy to change ports without editing launcher config

### Pattern 2: Info Pane Integration

Show port information in the info pane:

```yaml
projects:
  - name: Full Stack App
    commands:
      - name: Backend API
        icon: 🔧
        command: go run . -port $(ports-manager get backend)
        description: "API server on port $(ports-manager get backend)"
        info_file: ~/.config/tui-launcher/docs/backend.md
```

The info pane would show:
```
Backend API 🔧
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Command: go run . -port 8080
Port: 8080 (via ports-manager)
Directory: ~/projects/myapp/backend

API server for MyApp
Endpoints: /api/v1/*
```

### Pattern 3: Port Check Command

Add a utility command to show all ports:

```yaml
tools:
  - category: DevOps
    icon: 🛠️
    items:
      - name: Show All Ports
        icon: 🔌
        command: ports-manager list
        spawn: tmux-window

      - name: Port Conflicts
        icon: ⚠️
        command: lsof -i -P | grep LISTEN | grep $(ports-manager get backend)
        spawn: tmux-window
```

### Pattern 4: Profile with Port Dependencies

Launch entire stack with consistent ports:

```yaml
projects:
  - name: MyApp
    profiles:
      - name: Full Dev Environment
        icon: 🚀
        layout: main-vertical
        panes:
          - command: docker run -p $(ports-manager get postgres):5432 postgres:15
            cwd: ~/projects/myapp
          - command: go run . -port $(ports-manager get backend)
            cwd: ~/projects/myapp/backend
          - command: npm run dev -- -p $(ports-manager get frontend)
            cwd: ~/projects/myapp/frontend
```

One command launches entire stack with managed ports!

## Setup Workflow

### 1. Install ports-manager
```bash
cargo install cli-ports-manager
```

### 2. Configure ports for your project
```bash
ports-manager set backend 8080
ports-manager set frontend 3000
ports-manager set postgres 5432
ports-manager set redis 6379
```

Or let it auto-allocate:
```bash
# First run auto-assigns
ports-manager get backend  # Returns: 8080 (auto-assigned)
```

### 3. Update launcher config
Use `$(ports-manager get service)` in commands (see patterns above)

### 4. Share with team
```bash
# Commit ports config to repo
cp ~/.config/ports-manager/config.toml .devcontainer/ports.toml

# Team members import
ports-manager import .devcontainer/ports.toml
```

## Advanced: Port Validation

Add a command to validate ports before launch:

```yaml
tools:
  - category: DevOps
    items:
      - name: Validate Ports
        icon: ✅
        command: |
          #!/bin/bash
          echo "Checking configured ports..."
          for service in backend frontend postgres redis; do
            port=$(ports-manager get $service)
            if lsof -i :$port >/dev/null 2>&1; then
              echo "❌ Port $port ($service) already in use"
            else
              echo "✅ Port $port ($service) available"
            fi
          done
        spawn: tmux-window
```

## Example: Microservices Architecture

```yaml
projects:
  - name: Microservices Platform
    icon: 🏗️
    path: ~/projects/platform

    commands:
      # Core Services
      - name: Auth Service
        icon: 🔐
        command: go run . -port $(ports-manager get auth-service)
        cwd: ~/projects/platform/services/auth
        spawn: tmux-split-h

      - name: User Service
        icon: 👥
        command: go run . -port $(ports-manager get user-service)
        cwd: ~/projects/platform/services/users
        spawn: tmux-split-h

      - name: API Gateway
        icon: 🚪
        command: go run . -port $(ports-manager get gateway)
        cwd: ~/projects/platform/gateway
        spawn: tmux-split-h

      # Databases
      - name: Auth DB
        icon: 🗄️
        command: docker run -p $(ports-manager get auth-db):5432 -e POSTGRES_DB=auth postgres:15
        spawn: tmux-window

      - name: User DB
        icon: 🗄️
        command: docker run -p $(ports-manager get user-db):5432 -e POSTGRES_DB=users postgres:15
        spawn: tmux-window

      # Monitoring
      - name: Health Dashboard
        icon: 📊
        command: |
          echo "Service Health:"
          curl localhost:$(ports-manager get auth-service)/health
          curl localhost:$(ports-manager get user-service)/health
          curl localhost:$(ports-manager get gateway)/health
        spawn: tmux-window

    profiles:
      - name: Full Stack
        icon: 🎯
        layout: tiled
        panes:
          - command: docker run -p $(ports-manager get auth-db):5432 postgres:15
          - command: docker run -p $(ports-manager get user-db):5432 postgres:15
          - command: go run . -port $(ports-manager get auth-service)
            cwd: ~/projects/platform/services/auth
          - command: go run . -port $(ports-manager get user-service)
            cwd: ~/projects/platform/services/users
          - command: go run . -port $(ports-manager get gateway)
            cwd: ~/projects/platform/gateway
```

## Info Pane Enhancement

Show port info dynamically:

```go
// In updateInfoPane()
func (m *model) updateInfoPane() {
    // ... existing code ...

    // Extract port from command if using ports-manager
    if strings.Contains(currentItem.Command, "ports-manager get") {
        // Parse out service name
        serviceName := extractServiceName(currentItem.Command)

        // Get actual port
        cmd := exec.Command("ports-manager", "get", serviceName)
        output, _ := cmd.Output()
        port := strings.TrimSpace(string(output))

        // Add to info display
        m.infoContent += fmt.Sprintf("\nPort: %s (managed)", port)
    }
}
```

## Benefits for TUI Launcher Users

1. **Centralized config**: One place to manage all ports
2. **Team consistency**: Share port assignments via config file
3. **Conflict avoidance**: Auto-allocation prevents conflicts
4. **Flexibility**: Change ports without editing launcher config
5. **Documentation**: Ports are self-documenting in commands
6. **Validation**: Easy to check port availability before launch

## Implementation Timeline

- **v0.2.0**: Document pattern in README (users can use now!)
- **v0.3.0**: Add port extraction to info pane
- **v0.4.0**: Add validation commands to config template
- **v1.0.0**: Built-in ports-manager integration (optional)

---

**Status:** Ready to use today! No code changes needed.
**Recommendation:** Add to default config as example pattern
**Documentation:** Add to README as "Advanced: Port Management"
