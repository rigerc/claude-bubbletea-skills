// Package config provides configuration management for the application.
// It supports loading from JSON files, environment variables, and embedded defaults.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	koanfjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

var (
	// ErrInvalidConfig is returned when the configuration validation fails.
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrConfigNotFound is returned when no configuration file is found.
	ErrConfigNotFound = errors.New("configuration file not found")
)

// Config holds the application configuration.
// All fields are exported to support JSON marshaling and environment variable binding.
type Config struct {
	// LogLevel specifies the logging verbosity level.
	// Valid values: trace, debug, info, warn, error, fatal
	LogLevel string `json:"logLevel" mapstructure:"logLevel" koanf:"logLevel"`

	// Debug enables debug mode which sets log level to trace
	// and enables additional debugging features.
	Debug bool `json:"debug" mapstructure:"debug" koanf:"debug"`

	// UI contains user interface specific configuration.
	UI UIConfig `json:"ui" mapstructure:"ui" koanf:"ui"`

	// App contains general application configuration.
	App AppConfig `json:"app" mapstructure:"app" koanf:"app"`

	// Ralph contains configuration for the ralphio workflow engine.
	Ralph RalphConfig `json:"ralph" mapstructure:"ralph" koanf:"ralph"`
}

// RalphConfig contains configuration for the ralphio autonomous task execution engine.
type RalphConfig struct {
	// ProjectDir is the directory containing the tasks.json file.
	ProjectDir string `json:"projectDir" koanf:"projectDir"`

	// Agent is the adapter to use (claude, cursor, codex, opencode, kilo, pi).
	Agent string `json:"agent" koanf:"agent"`

	// AgentModel is the model to use for adapters that support model selection.
	AgentModel string `json:"agentModel" koanf:"agentModel"`

	// MaxRetries is the maximum number of times a task is retried on failure.
	MaxRetries int `json:"maxRetries" koanf:"maxRetries"`

	// RetryDelayMs is the delay in milliseconds between retries.
	RetryDelayMs int `json:"retryDelayMs" koanf:"retryDelayMs"`

	// AgentTimeoutMs is the maximum duration in milliseconds for a single agent execution.
	AgentTimeoutMs int `json:"agentTimeoutMs" koanf:"agentTimeoutMs"`

	// IterationDelayMs is the delay in milliseconds between loop iterations.
	IterationDelayMs int `json:"iterationDelayMs" koanf:"iterationDelayMs"`

	// Iterations is the maximum number of iterations before stopping.
	Iterations int `json:"iterations" koanf:"iterations"`

	// Validation contains validation command configuration.
	Validation ValidationConfig `json:"validation" koanf:"validation"`
}

// ValidationConfig contains configuration for post-task validation commands.
type ValidationConfig struct {
	// Enabled controls whether validation commands are run after each task.
	Enabled bool `json:"enabled" koanf:"enabled"`

	// Commands is the list of shell commands to run for validation.
	Commands []string `json:"commands" koanf:"commands"`

	// FailOnWarning treats warnings as failures when true.
	FailOnWarning bool `json:"failOnWarning" koanf:"failOnWarning"`
}

// UIConfig contains configuration specific to the user interface.
type UIConfig struct {
	// AltScreen runs the TUI in alternate screen mode (fullscreen).
	AltScreen bool `json:"altScreen" mapstructure:"altScreen" koanf:"altScreen"`

	// MouseEnabled enables mouse support in the TUI.
	MouseEnabled bool `json:"mouseEnabled" mapstructure:"mouseEnabled" koanf:"mouseEnabled"`

	// ThemeName specifies the color theme to use.
	ThemeName string `json:"themeName" mapstructure:"themeName" koanf:"themeName"`
}

// AppConfig contains general application configuration.
type AppConfig struct {
	// Name is the application name.
	Name string `json:"name" mapstructure:"name" koanf:"name"`

	// Version is the application version.
	Version string `json:"version" mapstructure:"version" koanf:"version"`

	// Title is the default window title.
	Title string `json:"title" mapstructure:"title" koanf:"title"`
}

// Load reads configuration from the specified file path.
// If the file does not exist, it returns ErrConfigNotFound.
// If the file exists but cannot be parsed, it returns an error.
func Load(path string) (*Config, error) {
	cfg := &Config{}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, ErrConfigNotFound
	}

	// Create koanf instance
	k := koanf.New(".")

	// Load from file
	if err := k.Load(file.Provider(path), koanfjson.Parser()); err != nil {
		return nil, fmt.Errorf("loading config from %s: %w", path, err)
	}

	// Unmarshal into config struct
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("parsing configuration: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadFromBytes loads configuration from a byte slice.
// This is useful for loading embedded default configurations.
func LoadFromBytes(data []byte) (*Config, error) {
	cfg := &Config{}

	// Create koanf instance
	k := koanf.New(".")

	// Load from bytes
	if err := k.Load(rawbytes.Provider(data), koanfjson.Parser()); err != nil {
		return nil, fmt.Errorf("loading config from bytes: %w", err)
	}

	// Unmarshal into config struct
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("parsing configuration: %w", err)
	}

	return cfg, nil
}

// Validate checks that the configuration is valid and returns an error if not.
func (c *Config) Validate() error {
	// Validate log level
	validLogLevels := map[string]bool{
		"trace": true, "debug": true, "info": true,
		"warn": true, "error": true, "fatal": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("%w: invalid log level '%s'", ErrInvalidConfig, c.LogLevel)
	}

	return nil
}

// ToJSON converts the configuration to a JSON byte slice.
// This is useful for writing the configuration to a file.
func (c *Config) ToJSON() ([]byte, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encoding configuration to JSON: %w", err)
	}
	return data, nil
}

// GetEffectiveLogLevel returns the effective log level.
// If debug mode is enabled, it returns "trace" regardless of the configured level.
func (c *Config) GetEffectiveLogLevel() string {
	if c.Debug {
		return "trace"
	}
	return c.LogLevel
}
