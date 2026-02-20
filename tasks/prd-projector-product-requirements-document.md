# Projector — Product Requirements Document

## Overview
**Problem**: Developers lack visibility into their local projects directory—they can't see all their projects at a glance and waste time manually navigating between codebases.

**Solution**: Projector provides a terminal UI that scans a configurable directory, displays all git repositories with their status, and enables quick navigation between projects.

---

## Technical Constraints
| Constraint | Value |
|------------|-------|
| Go version | 1.23+ |
| TUI Framework | Bubble Tea v2 (`charm.land/bubbletea/v2`) |
| Styling | Lip Gloss v2 (`charm.land/lipgloss/v2`) |
| Forms/Prompts | huh v2 (`charm.land/huh/v2`) |
| CLI | Cobra (`github.com/spf13/cobra`) |
| Config | koanf v2 (`github.com/knadh/koanf/v2`) |
| Logging | zerolog (`github.com/rs/zerolog`) |
| Git operations | go-git v6 (`github.com/go-git/go-git/v6`) |
| Directory traversal | godirwalk (`github.com/karrick/godirwalk`) |
| Quality gate | `go build ./...` && `go test ./...` |

---

## Configuration

### Priority Order
1. Embedded defaults (`config.DefaultConfig()`)
2. JSON file (`--config` flag or `$HOME/.projector.json`)
3. CLI flags (only when explicitly set)

### New Config Fields
```json
{
  "projectsDir": "~/projects",
  "editor": "",
  "rememberSelection": true,
  "scanDepth": 3
}
```

| Field | Type | Description |
|-------|------|-------------|
| `projectsDir` | string | Root directory to scan for git repos. CLI `--dir` overrides. |
| `editor` | string | Editor command. Falls back to `$EDITOR`. |
| `rememberSelection` | bool | Restore cursor to last selected project on startup. |
| `scanDepth` | int | Max directory depth for repo discovery. Monorepo-aware. |

---

## User Stories

### US-1: View Project List
**As a** developer  
**I want** to see all my git repositories in a scrollable list  
**So that** I can quickly identify and select a project to work on

**Acceptance Criteria**:
- Display project name (directory name) and relative path
- Show git branch name for each repo (via go-git `r.Head()`)
- Show ahead/behind commit counts relative to upstream (via go-git reference comparison)
- Show dirty files count (via go-git `worktree.Status()`)
- Use vim-style navigation (j/k or arrows)
- Support fuzzy search filtering (/ to activate)
- Empty state message when no repos found
- Async loading with progress indicator for 50+ repos

**Commands**: `go build ./...` && `go test ./...`

---

### US-2: Open Project in Terminal
**As a** developer  
**I want** to open a selected project in my terminal (cd to directory)  
**So that** I can run commands in that project's context

**Acceptance Criteria**:
- Press Enter to select and cd to project directory
- Output shell command to stdout for shell integration: `cd /path/to/project`
- Support `--print` flag for non-interactive mode: `projector --print`
- Exit with code 0 on success, 1 on cancellation

**Commands**: `go build ./...` && `go test ./...`

---

### US-3: Open Project in Editor
**As a** developer  
**I want** to open a selected project in my preferred editor/IDE  
**So that** I can start coding immediately

**Acceptance Criteria**:
- Press 'e' to open in configured editor
- Editor determined by: config `editor` field → `$EDITOR` env var
- Spawn editor process and wait for completion
- Return to list after editor closes
- Show error in status bar if editor not configured

**Commands**: `go build ./...` && `go test ./...`

---

### US-4: Configure Projects Directory
**As a** developer  
**I want** to configure which directory Projector scans  
**So that** I can manage projects from any location

**Acceptance Criteria**:
- CLI flag: `projector --dir ~/code`
- Config file: `projectsDir` field in `~/.projector.json`
- CLI flag overrides config file setting
- Validate directory exists on startup
- Show error and exit if directory invalid

**Commands**: `go build ./...` && `go test ./...`

---

### US-5: Filter Projects with Fuzzy Search
**As a** developer  
**I want** to filter the project list by typing  
**So that** I can quickly find a specific project

**Acceptance Criteria**:
- Press '/' to enter filter mode
- Type to fuzzy-match against project names and paths
- List updates in real-time as user types
- Press ESC to clear filter and exit filter mode
- Maintain selection within filtered results
- Show "[filtering]" indicator when active

**Commands**: `go build ./...` && `go test ./...`

---

### US-6: Handle Git Status Errors Gracefully
**As a** developer  
**I want** to see warnings for repos with issues  
**So that** I know which projects need attention

**Acceptance Criteria**:
- Display warning icon (⚠) next to repos with errors
- Show error details in status bar when repo selected
- Common errors: merge conflicts, detached HEAD, corrupted .git
- Allow navigation to problematic repos
- Don't skip repos silently—always show in list
- go-git returns typed errors for graceful handling

**Commands**: `go build ./...` && `go test ./...`

---

### US-7: Remember Last Selection
**As a** developer  
**I want** Projector to remember my last selected project  
**So that** I can quickly return to where I left off

**Acceptance Criteria**:
- When `rememberSelection: true`, save last selected path on exit
- Restore cursor to that project on next launch
- Fall back to first item if saved path no longer exists
- Setting configurable via `~/.projector.json`

**Commands**: `go build ./...` && `go test ./...`

---

### US-8: Monorepo-Aware Scanning
**As a** developer working with monorepos  
**I want** Projector to respect nested `.git` directories  
**So that** I see subprojects as separate entries

**Acceptance Criteria**:
- Use godirwalk for fast directory traversal (2-10x faster than `filepath.Walk`)
- Detect `.git` at any depth up to `scanDepth`
- Stop scanning deeper when `.git` found (return `godirwalk.SkipThis`)
- Display relative path from root for nested repos
- Support git worktrees as separate entries (go-git `PlainOpen` handles worktrees)

**Commands**: `go build ./...` && `go test ./...`

---

## Non-Functional Requirements

### Performance
- Initial scan completes in <2s for 100 repos (godirwalk provides 2-10x speedup)
- Filter response time <50ms
- Memory usage <50MB for 500 repos
- No external git binary required (go-git is pure Go)

### Compatibility
- Linux, macOS, Windows (WSL)
- Terminal with 256-color support minimum
- Graceful degradation for limited color terminals
- Cross-platform path handling via go-git/godirwalk

### Accessibility
- All actions keyboard-accessible
- Screen reader friendly status messages
- High contrast mode via theme selection

---

## File Structure (New Components)

```
projector/
├── internal/
│   ├── git/
│   │   ├── scanner.go      # Repo discovery using godirwalk
│   │   ├── scanner_test.go
│   │   └── status.go       # Git status via go-git API
│   └── ui/
│       └── screens/
│           └── projects.go # Main project list screen
└── config/
    └── defaults.go         # Add new config fields
```

---

## Dependencies (Additions)

| Package | Import Path | Purpose |
|---------|-------------|---------|
| go-git | `github.com/go-git/go-git/v6` | Pure Go git operations (status, branch, ahead/behind) |
| godirwalk | `github.com/karrick/godirwalk` | Fast directory traversal with skip support |

### go-git Usage Patterns
```go
// Open repository
r, err := git.PlainOpen(path)

// Get current branch
ref, err := r.Head()
branch := ref.Name().Short()

// Get worktree status (dirty files)
w, _ := r.Worktree()
status, _ := w.Status()
dirtyCount := len(status)

// Ahead/behind counts
localRef, _ := r.Reference(ref.Name(), true)
remoteRef, _ := r.Reference(plumbing.ReferenceName("refs/remotes/origin/"+branch), true)
// Compare commits between localRef and remoteRef
```

### godirwalk Usage Patterns
```go
err := godirwalk.Walk(rootDir, &godirwalk.Options{
    Callback: func(osPathname string, de *godirwalk.Dirent) error {
        if de.IsDir() && de.Name() == ".git" {
            // Found repo, add to list
            repos = append(repos, filepath.Dir(osPathname))
            return godirwalk.SkipThis // Don't descend into .git
        }
        return nil
    },
    Unsorted: true, // Faster for large directories
    FollowSymbolicLinks: false,
})
```