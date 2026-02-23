package adapter

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// runProcess executes an agent command, streaming output lines through the
// onOutput callback. prompt is appended as the final argument. If model is
// non-empty, "--model" and model are inserted immediately before the prompt.
// Environment variables from cfg.Env are merged with os.Environ().
//
// stdout and stderr are merged so that diagnostic output from agents remains
// visible in the TUI log pane.
func runProcess(ctx context.Context, cfg CommandConfig, prompt, model string, onOutput func(string)) error {
	args := make([]string, len(cfg.Command)-1)
	copy(args, cfg.Command[1:])

	if model != "" {
		args = append(args, "--model", model)
	}
	args = append(args, prompt)

	cmd := exec.CommandContext(ctx, cfg.Command[0], args...)
	cmd.Env = buildEnv(cfg.Env)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("creating stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting agent process: %w", err)
	}

	// Merge stdout and stderr into a single reader so the scanner sees both.
	merged := io.MultiReader(stdout, stderr)
	sc := bufio.NewScanner(merged)
	for sc.Scan() {
		text := ParseStreamLine(sc.Text())
		if text != "" {
			onOutput(text)
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("agent process exited with error: %w", err)
	}
	return nil
}

// buildEnv returns the current process environment merged with extra vars.
func buildEnv(extra map[string]string) []string {
	env := os.Environ()
	for k, v := range extra {
		env = append(env, k+"="+v)
	}
	return env
}
