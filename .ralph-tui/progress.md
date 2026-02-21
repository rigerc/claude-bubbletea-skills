# Ralph Progress Log

This file tracks progress across iterations. Agents update this file
after each iteration and it's included in prompts for context.

## Codebase Patterns (Study These First)

*Add reusable patterns discovered during development here.*

### Configuration Priority Pattern (from US-004)

When adding new config fields with CLI flag override:

1. Add field to `Config` struct with all 3 tags: `json:"name" mapstructure:"name" koanf:"name"`
2. Set sensible default in `DefaultConfig()`
3. Add CLI flag in `cmd/root.go` with `StringVarP` for shorthand support
4. Create `GetFieldName()` and `WasFieldNameSet()` functions using `rootCmd.PersistentFlags().Changed("flag-name")`
5. In `loadConfig()`, apply CLI override only when `WasFieldNameSet()` returns true
6. Validate in `main.go` before starting TUI, exit with error message if invalid

---

## [2026-02-21] - US-004
- Implemented projects directory configuration with CLI flag and config file support
- Files changed:
  - `projector/cmd/root.go`: Added `--dir/-d` flag, `GetProjectsDir()`, `WasProjectsDirSet()`
  - `projector/config/config.go`: Added `ProjectsDir` field, `ValidateProjectsDir()`, `ExpandDir()`
  - `projector/config/defaults.go`: Added default `ProjectsDir: "~"`
  - `projector/main.go`: Added `ValidateProjectsDir()` call, CLI override in `loadConfig()`
- **Learnings:**
  - Config priority is: defaults → config file → CLI flags (only when explicitly set via `Changed()`)
  - Use `rootCmd.PersistentFlags().Changed("flag-name")` to detect explicit CLI flag usage vs default value
  - Home directory expansion (`~`) must be handled explicitly with `os.UserHomeDir()`
---



---

## Parallel Task: Configure Projects Directory (US-004)

# Ralph Progress Log

This file tracks progress across iterations. Agents update this file
after each iteration and it's included in prompts for context.

## Codebase Patterns (Study These First)

*Add reusable patterns discovered during development here.*

### Configuration Priority Pattern (from US-004)

When adding new config fields with CLI flag override:

1. Add field to `Config` struct with all 3 tags: `json:"name" mapstructure:"name" koanf:"name"`
2. Set sensible default in `DefaultConfig()`
3. Add CLI flag in `cmd/root.go` with `StringVarP` for shorthand support
4. Create `GetFieldName()` and `WasFieldNameSet()` functions using `rootCmd.PersistentFlags().Changed("flag-name")`
5. In `loadConfig()`, apply CLI override only when `WasFieldNameSet()` returns true
6. Validate in `main.go` before starting TUI, exit with error message if invalid

---

## [2026-02-21] - US-004
- Implemented projects directory configuration with CLI flag and config file support
- Files changed:
  - `projector/cmd/root.go`: Added `--dir/-d` flag, `GetProjectsDir()`, `WasProjectsDirSet()`
  - `projector/config/config.go`: Added `ProjectsDir` field, `ValidateProjectsDir()`, `ExpandDir()`
  - `projector/config/defaults.go`: Added default `ProjectsDir: "~"`
  - `projector/main.go`: Added `ValidateProjectsDir()` call, CLI override in `loadConfig()`
- **Learnings:**
  - Config priority is: defaults → config file → CLI flags (only when explicitly set via `Changed()`)
  - Use `rootCmd.PersistentFlags().Changed("flag-name")` to detect explicit CLI flag usage vs default value
  - Home directory expansion (`~`) must be handled explicitly with `os.UserHomeDir()`
---

