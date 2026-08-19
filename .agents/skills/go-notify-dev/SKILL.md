---
name: go-notify-dev
description: >-
  Use this skill when developing, refactoring, debugging, testing, or building the go-notify
  or go-notify-manager Go applications. Provides step-by-step development workflows, GTK4 event
  propagation guidelines, DBus testing commands, and Makefile verification runbooks.
---

# go-notify Development & Refactoring Skill

This skill provides comprehensive instructions for agents working on `go-notify` (the D-Bus notification daemon) and `go-notify-manager` (the desktop history manager).

---

## 1. Development & Architecture Guidelines

When adding features or fixing bugs in this project, adhere strictly to the **MVVM Architecture**:

1. **Models** (`internal/models/`):
   - Pure data structs. No methods or dependencies on GTK / SQLite / DBus.
2. **Repositories** (`internal/repository/`):
   - Define interfaces for database access.
   - Use parameterized SQL queries (`db.Query`, `db.Exec`) to prevent SQL injection.
3. **Services** (`internal/service/`, `internal/services/`):
   - `LauncherService`: Handles `OpenURL` (xdg-open, firefox, gio) and desktop file (`/usr/share/applications/*.desktop`) parsing/launching.
   - `FocusModeService`: Manages D-Bus `SetFocusMode` and system audio muting via `pactl`.
   - `DBusServer`: Implements D-Bus `org.freedesktop.Notifications`.
4. **ViewModels** (`internal/viewmodel/`):
   - `NotificationViewModel`: Coordinates card click events (Button 1: open URL/App + ActionInvoked; Button 3: NotificationClosed(2)).
   - `HistoryViewModel`: Manages search filtering, date grouping, item deletions, and focus mode state.
5. **Views** (`internal/ui/`):
   - GTK4 UI components. Do not access DB or DBus directly—communicate via ViewModels and callbacks.
   - **Crucial**: Always call `gesture.SetPropagationPhase(gtk.PhaseCapture)` on event controllers attached to container boxes to prevent child labels/boxes from swallowing touch/click events.

---

## 2. Standard Development Workflow

### Step 1: Code Edit & Architectural Verification
- Ensure UI components delegate business logic to ViewModels.
- Ensure all DB accesses pass through Repository interfaces.
- Escape string markups with `glib.MarkupEscapeText()` when creating GTK Labels.

### Step 2: D-Bus Test Verification
To test notification sending via CLI:
```bash
# Send test notification with summary and body
notify-send "Test Title" "This is a test notification message with https://example.com"

# Send notification with urgency critical
notify-send -u critical "Urgent Title" "Critical alert!"
```

### Step 3: Run Full Quality Assurance Pipeline
Execute the Makefile targets to verify the build and tests:

```bash
# 1. Run static code analysis
make vet

# 2. Run unit & UI test suites
make test

# 3. Run race condition detector
make test-race

# 4. Run performance benchmarks
make bench

# 5. Run fuzz testing (optional)
make fuzz FUZZTIME=10s

# Or run the complete check pipeline in one command:
make check-all
```

### Step 4: Install Updated Binaries
After passing all tests, update system executables in `~/.local/bin`:

```bash
make install
```

If `go-notify` is currently running in background, restart it:
```bash
pkill -f /home/june/.local/bin/go-notify || true
nohup /home/june/.local/bin/go-notify > /dev/null 2>&1 &
```

---

## 3. Verification Checklist

Before reporting completion to the user:
- [ ] `make check-all` ran with 0 failures and 0 race warnings.
- [ ] `make install` successfully updated binaries in `~/.local/bin`.
- [ ] No hardcoded tokens, passwords, or `.db` files are tracked in Git.
- [ ] Documentation in `README.md` and `.agents/` reflects any structural changes.
