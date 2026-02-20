// scaffold is a minimal BubbleTea v2 skeleton.
// It wires up logging, an optional Cobra CLI, and then starts the TUI.
package main

import (
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"

	"projector/cmd"
	"projector/config"
	applogger "projector/internal/logger"
	"projector/internal/ui"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Command execution failed: %v\n", err)
		os.Exit(1)
	}

	if !cmd.ShouldRunUI() {
		return
	}

	cfg := loadConfig()

	if err := cfg.ValidateProjectsDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	logOutput, cleanup, err := setupLogOutput(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}

	if err := initLogger(cfg, logOutput); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logging: %v\n", err)
		os.Exit(1)
	}

	applogger.Info().Msg("Starting scaffold")

	if err := ui.Run(ui.New(*cfg)); err != nil {
		applogger.Fatal().Err(err).Msg("UI failed")
	}
}

// setupLogOutput returns the writer to use for logging and an optional cleanup
// function that must be deferred by the caller.
func setupLogOutput(cfg *config.Config) (io.Writer, func(), error) {
	if cfg.Debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return nil, nil, fmt.Errorf("opening debug log: %w", err)
		}
		return f, func() { f.Close() }, nil
	}
	return io.Discard, nil, nil
}

// initLogger initialises the global zerolog logger.
func initLogger(cfg *config.Config, output io.Writer) error {
	format := "console"
	if os.Getenv("ENV") == "production" {
		format = "json"
	}

	if err := applogger.Init(applogger.Config{
		Level:  applogger.LogLevel(cfg.GetEffectiveLogLevel()),
		Format: format,
		Output: output,
	}); err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}

	return nil
}

// loadConfig builds the effective config following priority order:
// defaults → config file → CLI flags (only when explicitly set).
func loadConfig() *config.Config {
	cfg := config.DefaultConfig()

	if path := cmd.GetConfigFile(); path != "" {
		fileCfg, err := config.Load(path)
		if err == nil {
			cfg = fileCfg
		}
	}

	if cmd.IsDebugMode() {
		cfg.Debug = true
	}
	if cmd.WasLogLevelSet() {
		cfg.LogLevel = cmd.GetLogLevel()
	}
	if cmd.WasProjectsDirSet() {
		cfg.ProjectsDir = cmd.GetProjectsDir()
	}

	return cfg
}
