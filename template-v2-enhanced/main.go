// template-v2-enhanced is a production-ready scaffold for building TUI applications.
// This is the entry point that initializes logging, configuration, and starts the CLI.
package main

import (
	"fmt"
	"io"
	"os"

	appconfig "template-v2-enhanced/config"
	"template-v2-enhanced/cmd"
	applogger "template-v2-enhanced/internal/logger"
	"template-v2-enhanced/internal/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	// Get the root command
	rootCmd := cmd.GetRootCmd()

	// Execute the Cobra command first
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Command execution failed: %v\n", err)
		os.Exit(1)
	}

	// Check if we should run the TUI (don't run for subcommands)
	if !cmd.ShouldRunUI() {
		return
	}

	// IMPORTANT: Set up file logging for TUI mode BEFORE any other output
	// BubbleTea occupies the terminal, so we cannot log to stderr/stdout
	// In debug mode, log to a file; otherwise, silence logging
	var logOutput io.Writer
	var logFileCleanup func()

	if cmd.IsDebugMode() {
		// Create log file for both BubbleTea and zerolog
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
			os.Exit(1)
		}
		logOutput = f
		logFileCleanup = func() { f.Close() }
	} else {
		// Discard all logging output in non-debug mode
		logOutput = io.Discard
	}

	// Ensure log file is closed on exit
	if logFileCleanup != nil {
		defer logFileCleanup()
	}

	// Initialize logging with the appropriate output (file or discard)
	if err := initializeLogging(logOutput); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logging: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		// Config is optional, use defaults if not found
		applogger.Debug().Err(err).Msg("Using default configuration")
		cfg = appconfig.DefaultConfig()
	}

	// Log startup information (goes to file in debug mode, discarded otherwise)
	applogger.Info().
		Str("version", rootCmd.Version).
		Str("logLevel", applogger.GetLevel().String()).
		Msgf("Starting %s", cfg.App.Name)

	// Create and run the UI with config
	model := ui.New(cfg)
	if err := ui.Run(model); err != nil {
		applogger.Fatal().Err(err).Msg("UI failed")
	}
}

// initializeLogging sets up the global logger based on command flags.
// The output writer determines where logs are sent (file or discard).
func initializeLogging(output io.Writer) error {
	// Get log level from flags
	logLevel := cmd.GetLogLevel()

	// Provide default if empty
	if logLevel == "" {
		logLevel = "info"
	}

	// Override log level to trace if debug mode is enabled
	if cmd.IsDebugMode() {
		logLevel = "trace"
	}

	// Determine log format based on environment
	logFormat := "console"
	if os.Getenv("ENV") == "production" {
		logFormat = "json"
	}

	// Initialize the logger with the provided output
	cfg := applogger.Config{
		Level:  applogger.LogLevel(logLevel),
		Format: logFormat,
		Output: output,
	}

	if err := applogger.Init(cfg); err != nil {
		return fmt.Errorf("initializing logger: %w", err)
	}

	// Set up debug logging if enabled
	if cmd.IsDebugMode() {
		applogger.Trace().Msg("Debug mode enabled")
	}

	return nil
}

// loadConfiguration loads the application configuration from file or defaults.
func loadConfiguration() (*appconfig.Config, error) {
	// Check if a config file was specified via flag
	configFile := cmd.GetConfigFile()
	if configFile != "" {
		cfg, err := appconfig.Load(configFile)
		if err != nil {
			return nil, fmt.Errorf("loading config from %s: %w", configFile, err)
		}
		applogger.Info().Str("file", configFile).Msg("Loaded configuration from file")
		return cfg, nil
	}

	// Try default locations
	home, err := os.UserHomeDir()
	if err == nil {
		defaultConfigPath := home + "/.template-v2-enhanced.json"
		cfg, err := appconfig.Load(defaultConfigPath)
		if err == nil {
			applogger.Info().Str("file", defaultConfigPath).Msg("Loaded configuration from default location")
			return cfg, nil
		}
	}

	// Try local directory
	cfg, err := appconfig.Load(".template-v2-enhanced.json")
	if err == nil {
		applogger.Info().Msg("Loaded configuration from local directory")
		return cfg, nil
	}

	// No config found, return defaults
	return appconfig.DefaultConfig(), nil
}

// init is called before main and can be used for setup.
// This function is empty but kept for potential future use.
func init() {
	// Register completion command
	// The completion command is already registered in cmd/completion.go
}
