// ralphio is the entry point for the ralphio autonomous task execution engine.
// It wires up logging, CLI flags, configuration, the orchestrator, and the TUI.
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"

	"ralphio/cmd"
	"ralphio/config"
	"ralphio/internal/adapter"
	applogger "ralphio/internal/logger"
	"ralphio/internal/orchestrator"
	"ralphio/internal/ui"
)

func main() {
	// Execute the Cobra CLI. Subcommands (version, completion) set runUI=false
	// and exit early; the root and run commands fall through to the TUI.
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Command execution failed: %v\n", err)
		os.Exit(1)
	}

	if !cmd.ShouldRunUI() {
		return
	}

	cfg := loadConfig()

	// In TUI mode the terminal is occupied, so all logging must go to a file
	// (debug mode) or be silenced entirely (normal mode).
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

	applogger.Info().Msg("Starting ralphio")

	// Create bidirectional orchestrator channels.
	orchMsgCh := make(chan tea.Msg, 64)
	orchCmdCh := make(chan any, 16)

	// Build the adapter from config.
	adpt := adapter.NewAdapter(
		adapter.AgentType(cfg.Ralph.Agent),
		cfg.Ralph.AgentModel,
	)

	projectDir := cfg.Ralph.ProjectDir

	// Start the orchestrator in a background goroutine.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orch := orchestrator.New(projectDir, adpt, orchMsgCh, orchCmdCh)
	go orch.Run(ctx)

	// Start the TUI, wiring in the orchestrator channels.
	if err := ui.Run(
		ui.New(*cfg),
		ui.WithOrchestrator(orchMsgCh, orchCmdCh),
	); err != nil {
		applogger.Fatal().Err(err).Msg("UI failed")
	}

	// Cancel the orchestrator context after the TUI exits.
	cancel()
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
		// ErrConfigNotFound or parse error → silently fall back to defaults.
	}

	// CLI flags override file/defaults only when explicitly passed.
	if cmd.IsDebugMode() {
		cfg.Debug = true
	}
	if cmd.WasLogLevelSet() {
		cfg.LogLevel = cmd.GetLogLevel()
	}
	if cmd.WasProjectDirSet() {
		cfg.Ralph.ProjectDir = cmd.GetProjectDir()
	}
	if cmd.WasAgentSet() {
		cfg.Ralph.Agent = cmd.GetAgent()
	}
	if cmd.WasModelSet() {
		cfg.Ralph.AgentModel = cmd.GetModel()
	}

	return cfg
}
