# Template V2 Enhanced

A production-ready scaffold for building terminal user interface (TUI) applications using BubbleTea v2, Bubbles v2, and Lip Gloss v2.

## Features

- **BubbleTea v2** - Modern Elm Architecture for terminal UIs
- **Bubbles v2** - Pre-built components (spinner, list, textinput, textarea, etc.)
- **Lip Gloss v2** - Terminal styling and layout
- **Cobra CLI** - Full command-line interface framework with flag support
- **Zerolog** - Structured, zero-allocation logging
- **Koanf** - Configuration management with file and environment variable support
- **Shell Completions** - Bash, Zsh, Fish, and PowerShell support
- **Debug Mode** - Comprehensive debugging and tracing support

## Project Structure

```
template-v2-enhanced/
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command with flags
│   ├── version.go         # Version subcommand
│   └── completion.go      # Shell completion subcommand
├── config/                # Configuration management
│   ├── config.go         # Configuration struct and loading
│   └── defaults.go       # Default configuration values
├── internal/
│   ├── logger/           # Zerolog structured logging
│   │   └── logger.go    # Logger setup and utilities
│   └── ui/               # BubbleTea UI components
│       ├── model.go      # Main UI model
│       └── help.go       # Key bindings and help
├── assets/               # Embedded files
│   └── config.default.json  # Default configuration template
├── main.go              # Application entry point
├── go.mod
├── go.sum
└── README.md
```

## Installation

### Prerequisites

- Go 1.25.2 or later

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd template-v2-enhanced

# Download dependencies
go mod download

# Build the application
go build -o template-v2-enhanced

# Run the application
./template-v2-enhanced
```

## Usage

### Basic Usage

```bash
# Run with default settings
template-v2-enhanced

# Run with debug logging
template-v2-enhanced --debug

# Run with custom log level
template-v2-enhanced --log-level trace

# Run with custom configuration file
template-v2-enhanced --config /path/to/config.json

# Show help
template-v2-enhanced --help

# Show version
template-v2-enhanced version
```

### Command-Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to configuration file | `$HOME/.template-v2-enhanced.json` |
| `--debug` | Enable debug mode (sets log level to trace) | `false` |
| `--log-level` | Set logging level | `info` |

### Log Levels

- `trace` - Most detailed logging
- `debug` - Debug information
- `info` - General information (default)
- `warn` - Warning messages
- `error` - Error messages
- `fatal` - Fatal errors (exits application)

### Shell Completions

#### Bash

```bash
# Load for current session
source <(template-v2-enhanced completion bash)

# Load for all sessions (Linux)
template-v2-enhanced completion bash > /etc/bash_completion.d/template-v2-enhanced

# Load for all sessions (macOS)
template-v2-enhanced completion bash > /usr/local/etc/bash_completion.d/template-v2-enhanced
```

#### Zsh

```bash
# Load for current session
source <(template-v2-enhanced completion zsh)

# Load for all sessions
template-v2-enhanced completion zsh > "${fpath[1]}/_template-v2-enhanced"
```

#### Fish

```bash
# Load for current session
template-v2-enhanced completion fish | source

# Load for all sessions
template-v2-enhanced completion fish > ~/.config/fish/completions/template-v2-enhanced.fish
```

#### PowerShell

```powershell
# Load for current session
template-v2-enhanced completion powershell | Out-String | Invoke-Expression

# Load for all sessions
template-v2-enhanced completion powershell > template-v2-enhanced.ps1
# Add to your PowerShell profile
```

## Configuration

### Configuration File

The application looks for a configuration file in the following locations:

1. Path specified via `--config` flag
2. `$HOME/.template-v2-enhanced.json`
3. `./.template-v2-enhanced.json`

### Example Configuration

```json
{
  "logLevel": "info",
  "debug": false,
  "ui": {
    "altScreen": true,
    "mouseEnabled": true,
    "themeName": "default"
  },
  "app": {
    "name": "template-v2-enhanced",
    "version": "1.0.0",
    "title": "Template V2 Enhanced"
  }
}
```

### Environment Variables

Configuration can also be set via environment variables with the `APP_` prefix:

```bash
export APP_LOGLEVEL=debug
export APP_DEBUG=true
export APP_UI_ALTSCREEN=true
export APP_APP_TITLE="My App"
```

### Configuration Schema

| Field | Type | Description |
|-------|------|-------------|
| `logLevel` | string | Minimum log level (trace, debug, info, warn, error, fatal) |
| `debug` | boolean | Enable debug mode (overrides logLevel to trace) |
| `ui.altScreen` | boolean | Run in alternate screen mode (fullscreen) |
| `ui.mouseEnabled` | boolean | Enable mouse support |
| `ui.themeName` | string | Color theme name |
| `app.name` | string | Application name |
| `app.version` | string | Application version |
| `app.title` | string | Window title |

## Development

### Using as a Template

This project is designed to be used as a template for new TUI applications:

1. **Copy the project** to a new directory
2. **Rename the module** in `go.mod`
3. **Update import paths** throughout the codebase
4. **Customize the UI** in `internal/ui/model.go`
5. **Add your commands** in `cmd/`
6. **Extend configuration** in `config/config.go`

### Key Files to Modify

- `main.go` - Entry point, customize initialization logic
- `internal/ui/model.go` - Main UI model, add your components
- `config/config.go` - Configuration structure, add your fields
- `cmd/root.go` - Root command, add your flags and subcommands

### Adding a New Subcommand

```go
// In cmd/mycommand.go
package cmd

import (
    "github.com/spf13/cobra"
)

var myCommand = &cobra.Command{
    Use:   "mycommand [args]",
    Short: "Description of my command",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Your command logic here
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCommand)
}
```

### Adding Configuration Fields

```go
// In config/config.go
type Config struct {
    // ... existing fields ...
    MySetting string `json:"mySetting" mapstructure:"mySetting" koanf:"mySetting"`
}
```

### Using the Logger

```go
import applogger "template-v2-enhanced/internal/logger"

// Log at different levels
applogger.Trace().Msg("Very detailed message")
applogger.Debug().Msg("Debug information")
applogger.Info().Msg("General information")
applogger.Warn().Msg("Warning message")
applogger.Error().Msg("Error occurred")

// With structured fields
applogger.Info().
    Str("user", "john").
    Int("count", 42).
    Msg("User performed action")

// With error
applogger.Error().Err(err).Msg("Operation failed")
```

### Debugging

Enable debug mode for detailed logging:

```bash
template-v2-enhanced --debug
```

Or set the environment variable:

```bash
export DEBUG=true
template-v2-enhanced
```

Debug logging includes:
- Function entry/exit traces
- State transitions
- Event handling details
- Configuration loading

### Testing

Run tests with:

```bash
go test ./...
go test -v ./...
go test -race ./...
```

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [BubbleTea v2](https://charm.land/bubbletea/v2) | v2.0.0-rc.2 | TUI framework |
| [Bubbles v2](https://charm.land/bubbles/v2) | v2.0.0-rc.1 | UI components |
| [Lip Gloss v2](https://charm.land/lipgloss/v2) | v2.0.0-beta.3 | Terminal styling |
| [Cobra](https://github.com/spf13/cobra) | v1.8.2 | CLI framework |
| [Viper](https://github.com/spf13/viper) | v1.19.0 | Configuration management |
| [Zerolog](https://github.com/rs/zerolog) | v1.33.0 | Structured logging |
| [Koanf](https://github.com/knadh/koanf) | v2.1.2 | Configuration parsing |

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please use the GitHub issue tracker.
