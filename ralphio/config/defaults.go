// Package config provides configuration management for the application.
package config

// DefaultConfig returns a configuration with sensible default values.
// These defaults can be overridden by loading a configuration file or
// setting environment variables.
func DefaultConfig() *Config {
	return &Config{
		LogLevel: "info",
		Debug:    false,
		UI: UIConfig{
			AltScreen:    false,
			MouseEnabled: true,
			ThemeName:    "default",
		},
		App: AppConfig{
			Name:    "ralphio",
			Version: "1.0.0",
			Title:   "ralphio",
		},
		Ralph: RalphConfig{
			ProjectDir:       ".",
			Agent:            "claude",
			AgentModel:       "",
			MaxRetries:       3,
			RetryDelayMs:     5000,
			AgentTimeoutMs:   1800000,
			IterationDelayMs: 2000,
			Iterations:       10,
			Validation: ValidationConfig{
				Enabled:       false,
				Commands:      []string{},
				FailOnWarning: false,
			},
		},
	}
}

// DefaultConfigJSON returns the default configuration as a JSON byte slice.
// This can be used to create a default configuration file or as a fallback
// when no configuration file is found.
func DefaultConfigJSON() ([]byte, error) {
	return DefaultConfig().ToJSON()
}
