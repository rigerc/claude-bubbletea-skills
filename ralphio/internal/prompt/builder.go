// Package prompt builds the agent prompt for each loop iteration by reading
// the appropriate mode-specific prompt file from the project directory.
package prompt

import (
	"fmt"
	"os"
	"path/filepath"

	"ralphio/internal/plan"
)

const (
	defaultBuildPromptFile = "PROMPT_build.md"
	defaultPlanPromptFile  = "PROMPT_plan.md"
)

// Builder composes the agent prompt from project files.
type Builder struct {
	projectDir string
	buildFile  string // path to PROMPT_build.md (or override)
	planFile   string // path to PROMPT_plan.md (or override)
}

// New returns a Builder for the given project directory using default prompt filenames.
func New(projectDir string) *Builder {
	return &Builder{
		projectDir: projectDir,
		buildFile:  filepath.Join(projectDir, defaultBuildPromptFile),
		planFile:   filepath.Join(projectDir, defaultPlanPromptFile),
	}
}

// WithBuildFile overrides the path to the build-mode prompt file.
func (b *Builder) WithBuildFile(path string) *Builder {
	b.buildFile = path
	return b
}

// WithPlanFile overrides the path to the plan-mode prompt file.
func (b *Builder) WithPlanFile(path string) *Builder {
	b.planFile = path
	return b
}

// Build returns the prompt text for the given mode. It reads the
// appropriate PROMPT_*.md file from the project directory. Returns
// an error if the file cannot be read.
func (b *Builder) Build(mode plan.LoopMode) (string, error) {
	path := b.buildFile
	if mode == plan.ModePlanning {
		path = b.planFile
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading prompt file %s: %w", path, err)
	}

	return string(data), nil
}

// HasBuildPrompt reports whether PROMPT_build.md exists in the project directory.
func (b *Builder) HasBuildPrompt() bool {
	_, err := os.Stat(b.buildFile)
	return err == nil
}

// HasPlanPrompt reports whether PROMPT_plan.md exists in the project directory.
func (b *Builder) HasPlanPrompt() bool {
	_, err := os.Stat(b.planFile)
	return err == nil
}

// HasAgentsMD reports whether AGENTS.md exists in the project directory.
func (b *Builder) HasAgentsMD() bool {
	_, err := os.Stat(filepath.Join(b.projectDir, "AGENTS.md"))
	return err == nil
}
